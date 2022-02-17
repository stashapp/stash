package ffmpeg

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

const (
	hlsSegmentLength = 2

	maxSegmentWait       = 5 * time.Second
	segmentCheckInterval = 100 * time.Millisecond

	// number of segments to wait to be generated before we
	// restart the transcode process at the requested segment
	maxSegmentWaitGap = 5

	// number of segments ahead of the currently streaming segment
	// before we stop the transcode process to save CPU
	maxSegmentStopGap = 15

	maxIdleTime     = 10 * time.Second
	monitorInterval = 2 * time.Second
)

type StreamManagerConfig interface {
	GetMaxStreamingTranscodeSize() models.StreamingResolutionEnum
}

type transcodeProcess struct {
	cmd          *exec.Cmd
	path         string
	cancel       context.CancelFunc
	cancelled    bool
	hash         string
	startSegment int
}

type runningStream struct {
	running      bool
	lastAccessed time.Time
	segment      int
}

type StreamManager struct {
	cacheDir string
	encoder  FFMpeg
	ffprobe  FFProbe
	config   StreamManagerConfig

	context    context.Context
	cancelFunc context.CancelFunc

	runningTranscodes map[string]*transcodeProcess
	transcodesMutex   sync.Mutex

	runningStreams map[string]runningStream
	streamsMutex   sync.Mutex
}

func NewStreamManager(cacheDir string, encoder FFMpeg, ffprobe FFProbe, config StreamManagerConfig) *StreamManager {
	if cacheDir == "" {
		panic("cache directory is not set")
	}

	ctx, cancel := context.WithCancel(context.Background())

	ret := &StreamManager{
		cacheDir:          cacheDir,
		encoder:           encoder,
		ffprobe:           ffprobe,
		config:            config,
		context:           ctx,
		cancelFunc:        cancel,
		runningTranscodes: make(map[string]*transcodeProcess),
		runningStreams:    make(map[string]runningStream),
	}

	go ret.monitorStreams()

	return ret
}

// Shutdown shuts down the stream manager, killing any running transcoding processes and removing all cached files.
func (sm *StreamManager) Shutdown() {
	sm.cancelFunc()
}

// WriteHLSPlaylist writes a playlist manifest to w. The URLs for the segments
// are generated using urlFormat. urlFormat is expected to include a single
// %d argument, which will be populated with the segment index.
func (sm *StreamManager) WriteHLSPlaylist(duration float64, urlFormat string, w io.Writer) {
	fmt.Fprint(w, "#EXTM3U\n")
	fmt.Fprint(w, "#EXT-X-VERSION:3\n")
	fmt.Fprint(w, "#EXT-X-MEDIA-SEQUENCE:0\n")
	fmt.Fprintf(w, "#EXT-X-TARGETDURATION:%d\n", hlsSegmentLength)
	fmt.Fprint(w, "#EXT-X-PLAYLIST-TYPE:VOD\n")

	leftover := duration
	segment := 0

	for leftover > 0 {
		thisLength := float64(hlsSegmentLength)
		if leftover < thisLength {
			thisLength = leftover
		}

		fmt.Fprintf(w, "#EXTINF:%f,\n", thisLength)
		fmt.Fprintf(w, urlFormat+"\n", segment)

		leftover -= thisLength
		segment++
	}

	fmt.Fprint(w, "#EXT-X-ENDLIST\n")
}

func (sm *StreamManager) segmentDirectory(hash string) string {
	return filepath.Join(sm.cacheDir, hash)
}

func (sm *StreamManager) segmentFilename(hash string, segment int) string {
	return filepath.Join(sm.segmentDirectory(hash), fmt.Sprintf("%d.ts", segment))
}

func (sm *StreamManager) segmentExists(segmentFilename string) bool {
	exists, _ := fsutil.FileExists(segmentFilename)
	return exists
}

// lastTranscodedSegment returns the most recent segment file created. Returns -1 if no files are found.
func (sm *StreamManager) lastTranscodedSegment(hash string) int {
	files, _ := ioutil.ReadDir(sm.segmentDirectory(hash))

	var mostRecent fs.FileInfo
	for _, f := range files {
		// ignore non-ts files
		if filepath.Ext(f.Name()) != ".ts" {
			continue
		}

		if mostRecent == nil || f.ModTime().After(mostRecent.ModTime()) {
			mostRecent = f
		}
	}

	segment := -1
	if mostRecent != nil {
		_, _ = fmt.Sscanf(filepath.Base(mostRecent.Name()), "%d.ts", &segment)
	}

	return segment
}

func (sm *StreamManager) streamNotify(ctx context.Context, hash string, segment int) {
	sm.streamsMutex.Lock()
	sm.runningStreams[hash] = runningStream{
		running: true,
		segment: segment,
	}
	sm.streamsMutex.Unlock()

	go func() {
		<-ctx.Done()

		sm.streamsMutex.Lock()
		sm.runningStreams[hash] = runningStream{
			lastAccessed: time.Now(),
			segment:      segment,
		}
		sm.streamsMutex.Unlock()
	}()
}

func (sm *StreamManager) streamTSFunc(hash string, segment int) http.HandlerFunc {
	fn := sm.segmentFilename(hash, segment)
	return func(w http.ResponseWriter, r *http.Request) {
		sm.streamNotify(r.Context(), hash, segment)
		w.Header().Set("Content-Type", "video/mp2t")
		http.ServeFile(w, r, fn)
	}
}

func (sm *StreamManager) waitAndStreamTSFunc(hash string, segment int) http.HandlerFunc {
	fn := sm.segmentFilename(hash, segment)
	started := time.Now()

	logger.Debugf("waiting for segment file %q to be generated", fn)
	for {
		if sm.segmentExists(fn) {
			// TODO - may need to wait for transcode process to finish writing the file first
			return sm.streamTSFunc(hash, segment)
		}

		now := time.Now()
		if started.Add(maxSegmentWait).Before(now) {
			logger.Warnf("timed out waiting for segment file %q to be generated", fn)

			return func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "timed out waiting for segment file to be generated", http.StatusInternalServerError)
			}
		}

		time.Sleep(segmentCheckInterval)
	}
}

// StreamTS returns a http.HandlerFunc that streams a TS segment for src.
// If the segment exists in the cache directory, then it is streamed.
// Otherwise, a transcode process will be started for the provided segment. If
// a transcode process is running already, then it will be killed before the new
// process is started.
func (sm *StreamManager) StreamTS(src string, hash string, segment int, videoCodec VideoCodec) http.HandlerFunc {
	if sm.cacheDir == "" {
		logger.Error("cannot live transcode files because cache dir is empty")
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "cannot live transcode files because cache dir is empty", http.StatusInternalServerError)
		}
	}

	if hash == "" {
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "invalid hash", http.StatusBadRequest)
		}
	}

	onTranscodeError := func(err error) http.HandlerFunc {
		errStr := fmt.Sprintf("error starting transcode process: %v", err.Error())
		logger.Error(errStr)

		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, errStr, http.StatusInternalServerError)
		}
	}

	segmentFilename := sm.segmentFilename(hash, segment)

	// check if transcoded file already exists
	// TODO - may need to wait for transcode process to finish writing the file first
	// if so, return it
	if sm.segmentExists(segmentFilename) {
		return sm.streamTSFunc(hash, segment)
	}

	// check if transcoding process is already running
	// lock the mutex here to ensure we don't start multiple processes
	sm.transcodesMutex.Lock()

	tp := sm.runningTranscodes[hash]

	// if not, start one at the applicable time, wait and return stream
	if tp == nil {
		var err error
		_, err = sm.startTranscode(src, hash, segment, videoCodec)
		sm.transcodesMutex.Unlock()

		if err != nil {
			return onTranscodeError(err)
		}

		return sm.waitAndStreamTSFunc(hash, segment)
	}

	// check if transcoding process is about to transcode the necessary segment
	lastSegment := sm.lastTranscodedSegment(hash)

	if lastSegment <= segment && lastSegment+maxSegmentWaitGap >= segment {
		// if so, wait and return
		sm.transcodesMutex.Unlock()
		return sm.waitAndStreamTSFunc(hash, segment)
	}

	logger.Debugf("restarting transcode since up to segment #%d and #%d was requested", lastSegment, segment)

	// otherwise, stop the existing transcoding process, restart at the applicable time
	// wait and return stream
	sm.stopTranscode(hash)

	_, err := sm.startTranscode(src, hash, segment, videoCodec)
	sm.transcodesMutex.Unlock()

	if err != nil {
		return onTranscodeError(err)
	}
	return sm.waitAndStreamTSFunc(hash, segment)
}

func (sm *StreamManager) segmentToTime(segment int) float64 {
	return float64(segment * hlsSegmentLength)
}

func (sm *StreamManager) getTranscodeArgs(probeResult *VideoFile, outputPath string, segment int, videoCodec VideoCodec) Args {
	var args Args
	args = append(args, "-hide_banner")
	args = args.LogLevel(LogLevelError)

	if segment > 0 {
		args = args.Seek(sm.segmentToTime(segment))
		// without this ffmpeg would sometimes generate empty ts files when seeking
		args = append(args, "-noaccurate_seek")
	}

	args = args.Input(probeResult.Path)

	args = args.VideoCodec(videoCodec)

	if videoCodec != VideoCodecCopy {
		args = append(args,
			"-c:v", "libx264",
			"-pix_fmt", "yuv420p",
			"-profile:v", "high",
			"-level", "4.2",
			"-preset", "superfast",
			"-crf", "23",
			"-r", "30",
			"-g", "60",
			"-x264-params", "no-scenecut=1",
			"-force_key_frames", "0")

		// don't set scale when copying video stream
		var videoFilter VideoFilter
		videoFilter = videoFilter.ScaleMax(probeResult.Width, probeResult.Height, sm.config.GetMaxStreamingTranscodeSize().GetMaxResolution())
		args = args.VideoFilter(videoFilter)
	}

	args = args.AudioCodec(AudioCodecAAC)

	args = append(args,
		// this is needed for 5-channel ac3 files
		"-ac", "2",
		"-copyts",
		"-avoid_negative_ts", "disabled",
		"-strict", "-2",
		"-f", "hls",
		"-start_number", fmt.Sprint(segment),
		"-hls_time", "2",
		"-hls_segment_type", "mpegts",
		"-hls_playlist_type", "vod",
		"-hls_list_size", "0",
		"-hls_segment_filename", filepath.Join(outputPath, "%d.ts"),
		filepath.Join(outputPath, "playlist.m3u8"),
	)

	return args
}

// assumes mutex is held
func (sm *StreamManager) startTranscode(src string, hash string, segment int, videoCodec VideoCodec) (*transcodeProcess, error) {
	probeResult, err := sm.ffprobe.NewVideoFile(src)
	if err != nil {
		return nil, err
	}

	outputPath := sm.segmentDirectory(hash)
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(sm.context)

	args := sm.getTranscodeArgs(probeResult, outputPath, segment, videoCodec)
	cmd := sm.encoder.Command(ctx, args)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Error("FFMPEG stderr not available: " + err.Error())
	}

	stdout, err := cmd.StdoutPipe()
	if nil != err {
		logger.Error("FFMPEG stdout not available: " + err.Error())
	}

	logger.Tracef("running %s", cmd.String())
	if err := cmd.Start(); err != nil {
		cancel()
		return nil, err
	}

	// TODO - handle lock manager stuff

	p := &transcodeProcess{
		cmd:          cmd,
		path:         probeResult.Path,
		cancel:       cancel,
		hash:         hash,
		startSegment: segment,
	}
	sm.runningTranscodes[hash] = p

	// mark the stream as accessed to ensure it is not immediately cleaned up
	sm.streamsMutex.Lock()
	sm.runningStreams[hash] = runningStream{
		lastAccessed: time.Now(),
	}
	sm.streamsMutex.Unlock()

	go sm.waitAndDeregister(hash, p, stdout, stderr)

	return p, nil
}

// assumes mutex is held
func (sm *StreamManager) stopTranscode(hash string) {
	p := sm.runningTranscodes[hash]
	if p != nil {
		p.cancel()
		p.cancelled = true
		delete(sm.runningTranscodes, hash)
	}
}

func (sm *StreamManager) waitAndDeregister(hash string, p *transcodeProcess, stdout, stderr io.Reader) {
	cmd := p.cmd

	errStr, _ := io.ReadAll(stderr)
	outStr, _ := io.ReadAll(stdout)

	err := cmd.Wait()

	// make sure that cancel is called to prevent memory leaks
	p.cancel()

	// TODO - handle lock manager stuff

	// don't log error if cancelled
	if !p.cancelled && err != nil {
		e := string(errStr)
		if e == "" {
			e = string(outStr)

			if e == "" {
				e = err.Error()
			}
		}

		// error message should be in the stderr stream
		logger.Errorf("ffmpeg error when running command <%s>: %s", strings.Join(cmd.Args, " "), e)
	}

	// remove from running transcodes
	sm.transcodesMutex.Lock()
	defer sm.transcodesMutex.Unlock()

	// only delete if is the same process
	if sm.runningTranscodes[hash] == p {
		delete(sm.runningTranscodes, hash)
	}
}

// monitorStreams checks for stale streams and removes them. When the manager context
// is cancelled, stopAndRemoveAll will be called. Should be called in its own goroutine.
func (sm *StreamManager) monitorStreams() {
	for {
		select {
		case <-time.After(monitorInterval):
			sm.removeStaleFiles()
		case <-sm.context.Done():
			sm.stopAndRemoveAll()
			return
		}
	}
}

func (sm *StreamManager) removeStaleFiles() {
	// check for the last time a stream was accessed
	// remove anything over a certain age
	sm.streamsMutex.Lock()
	defer sm.streamsMutex.Unlock()

	var toRemove []string

	now := time.Now()
	for hash, stream := range sm.runningStreams {
		if !stream.running && stream.lastAccessed.Add(maxIdleTime).Before(now) {
			// Stream expired. Cancel the transcode process and delete the files
			logger.Debugf("stream for hash %q not accessed recently. Cancelling transcode and removing files", hash)
			func() {
				sm.transcodesMutex.Lock()
				defer sm.transcodesMutex.Unlock()

				sm.stopAndRemoveTranscodeFiles(hash)

				toRemove = append(toRemove, hash)
			}()
		} else {
			// check if the last transcoded file is way ahead of the current streaming one
			// if so, stop the transcode to save CPU
			lastGenerated := sm.lastTranscodedSegment(hash)
			sm.transcodesMutex.Lock()
			if sm.runningTranscodes[hash] != nil && stream.segment+maxSegmentStopGap < lastGenerated {
				logger.Debugf("stopping transcode for hash %q as last generated segment %d is too far ahead of current segment %d", hash, lastGenerated, stream.segment)
				sm.stopTranscode(hash)
			}
			sm.transcodesMutex.Unlock()
		}
	}

	for _, hash := range toRemove {
		delete(sm.runningStreams, hash)
	}
}

// stopAndRemoveAll stops all current streams and removes all cache files
func (sm *StreamManager) stopAndRemoveAll() {
	sm.streamsMutex.Lock()
	sm.transcodesMutex.Lock()
	defer sm.streamsMutex.Unlock()
	defer sm.transcodesMutex.Unlock()

	for hash := range sm.runningStreams {
		sm.stopAndRemoveTranscodeFiles(hash)
	}

	// ensure nothing else can use the map
	sm.runningStreams = nil
	sm.runningTranscodes = nil
}

// assume lock is held
func (sm *StreamManager) stopAndRemoveTranscodeFiles(hash string) {
	sm.stopTranscode(hash)

	dir := sm.segmentDirectory(hash)
	if err := os.RemoveAll(dir); err != nil {
		logger.Warnf("error removing segment directory %q: %v", dir, err)
	}
}
