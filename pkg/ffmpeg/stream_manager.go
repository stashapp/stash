package ffmpeg

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

const (
	hlsSegmentLength = 2

	maxSegmentWait       = 10 * time.Second
	segmentCheckInterval = 100 * time.Millisecond

	// number of segments to wait to be generated before we
	// restart the transcode process at the requested segment
	maxSegmentWaitGap = 3

	// number of segments ahead of the currently streaming segment
	// before we stop the transcode process to save CPU
	maxSegmentStopGap = 30

	// number of segments ahead of the currently streaming segment
	// before we restart the transcode process
	maxSegmentRestartGap = 15

	// number of segments ahead of the start segment the steam must
	// be before we consider cleaning up the transcode process
	minSegmentTranscode = 5

	// time to after segment generation delay before serving
	segmentCreationDelay = 1 * time.Second

	maxIdleTime     = 60 * time.Second
	monitorInterval = 2 * time.Second

	// Cancel timeout for ffmpeg so there are no corrupted segments
	cancelTimeout = 2 * time.Second
)

type StreamManagerConfig interface {
	GetMaxStreamingTranscodeSize() models.StreamingResolutionEnum
}

type transcodeProcess struct {
	cmd          *exec.Cmd
	context      context.Context
	path         string
	cancel       context.CancelFunc
	cancelled    bool
	hash         string
	resolution   string
	startSegment int
}

type runningStream struct {
	running      bool
	lastAccessed time.Time
	segment      int
}

type StreamManager struct {
	cacheDir    string
	ffmpeg      FFMpeg
	ffprobe     FFProbe
	config      StreamManagerConfig
	lockManager *fsutil.ReadLockManager

	context    context.Context
	cancelFunc context.CancelFunc

	runningTranscodes map[string]*transcodeProcess
	transcodesMutex   sync.Mutex

	runningStreams map[string]runningStream
	streamsMutex   sync.Mutex
}

func NewStreamManager(cacheDir string, ffmpeg FFMpeg, ffprobe FFProbe, config StreamManagerConfig, lockManager *fsutil.ReadLockManager) *StreamManager {
	if cacheDir == "" {
		panic("cache directory is not set")
	}

	ctx, cancel := context.WithCancel(context.Background())

	ret := &StreamManager{
		cacheDir:          cacheDir,
		ffmpeg:            ffmpeg,
		ffprobe:           ffprobe,
		config:            config,
		lockManager:       lockManager,
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

func (sm *StreamManager) segmentDirectory(hashResolution string) string {
	return filepath.Join(sm.cacheDir, hashResolution)
}

func (sm *StreamManager) segmentFilename(hashResolution string, segment int) string {
	return filepath.Join(sm.segmentDirectory(hashResolution), fmt.Sprintf("%d.ts", segment))
}

func (sm *StreamManager) segmentExists(segmentFilename string) bool {
	exists, _ := fsutil.FileExists(segmentFilename)
	return exists
}

// lastTranscodedSegment returns the most recent segment file created. Returns -1 if no files are found.
func (sm *StreamManager) lastTranscodedSegment(hashResolution string) int {
	files, _ := os.ReadDir(sm.segmentDirectory(hashResolution))

	var mostRecent fs.FileInfo
	for _, f := range files {
		// ignore non-ts files
		if filepath.Ext(f.Name()) != ".ts" {
			continue
		}

		info, _ := f.Info()
		if info == nil {
			continue
		}

		if mostRecent == nil || info.ModTime().After(mostRecent.ModTime()) {
			mostRecent = info
		}
	}

	segment := -1
	if mostRecent != nil {
		_, _ = fmt.Sscanf(filepath.Base(mostRecent.Name()), "%d.ts", &segment)
	}

	return segment
}

func (sm *StreamManager) streamNotify(ctx context.Context, hashResolution string, segment int) {
	sm.streamsMutex.Lock()
	defer sm.streamsMutex.Unlock()
	sm.runningStreams[hashResolution] = runningStream{
		running: true,
		segment: segment,
	}

	go func() {
		<-ctx.Done()

		sm.streamsMutex.Lock()
		defer sm.streamsMutex.Unlock()
		sm.runningStreams[hashResolution] = runningStream{
			lastAccessed: time.Now(),
			segment:      segment,
		}
	}()
}

func (sm *StreamManager) streamTSFunc(hashResolution string, segment int) http.HandlerFunc {
	fn := sm.segmentFilename(hashResolution, segment)
	return func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		for {
			select {
			case <-r.Context().Done():
				logger.Trace("exiting streamTSFunc because context is done")
				return
			case <-time.After(segmentCheckInterval):
				// notify that we're still waiting so it doesn't get cleaned up
				sm.streamNotify(r.Context(), hashResolution, segment)
				now := time.Now()
				switch {
				case sm.segmentExists(fn):
					fi, err := os.Stat(fn)
					switch {
					case err != nil:
						logger.Warnf("error getting file info for %s: %s", fn, err)
					case now.Add(segmentCreationDelay).After(fi.ModTime()):
						logger.Tracef("streaming segment %d hashResolution %s", segment, hashResolution)
						sm.streamNotify(r.Context(), hashResolution, segment)
						w.Header().Set("Content-Type", "video/mp2t")
						http.ServeFile(w, r, fn)
						return
					default:
						logger.Tracef("segment %d hashResolution %s is still being created", segment, hashResolution)
					}
				case started.Add(maxSegmentWait).Before(now):
					logger.Warnf("timed out waiting for segment file %q to be generated", fn)
					http.Error(w, "timed out waiting for segment file to be generated", http.StatusInternalServerError)
					return
				}
			}
		}
	}
}

// StreamTS returns a http.HandlerFunc that streams a TS segment for src.
// If the segment exists in the cache directory, then it is streamed.
// Otherwise, a transcode process will be started for the provided segment. If
// a transcode process is running already, then it will be killed before the new
// process is started.
// resolution is the frame height as a string
// cannot use copy because keyframe interval is unable to be determined prior to encoding
func (sm *StreamManager) StreamTS(src string, hash string, segment int, resolution string) http.HandlerFunc {
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

	segmentFilename := sm.segmentFilename(hash+resolution, segment)

	tp := sm.runningTranscodes[hash+resolution]
	transcodeStartSegment := segment
	// check if transcoded file already exists
	// TODO - may need to wait for transcode process to finish writing the file first
	// if so, return it
	if sm.segmentExists(segmentFilename) {
		lastTranscodedSegment := sm.lastTranscodedSegment(hash + resolution)
		// Don't need to start transcode process yet
		if tp != nil || lastTranscodedSegment >= segment+maxSegmentRestartGap {
			return sm.streamTSFunc(hash+resolution, segment)
		} else {
			logger.Debugf("restarting transcode since last transcoded segment %d is close to requested segment %d", lastTranscodedSegment, segment)
			// reencode last segment to ensure it lines up with the previous segment
			transcodeStartSegment = lastTranscodedSegment
		}
	}

	// check if transcoding process is already running
	// lock the mutex here to ensure we don't start multiple processes
	sm.transcodesMutex.Lock()
	defer sm.transcodesMutex.Unlock()

	if tp != nil {
		// check if transcoding process is about to transcode the necessary segment
		lastSegment := sm.lastTranscodedSegment(hash + resolution)

		if (lastSegment <= segment && lastSegment+maxSegmentWaitGap >= segment) || tp.startSegment == segment {
			// if so, wait and return
			return sm.streamTSFunc(hash+resolution, segment)
		}

		logger.Debugf("restarting transcode since up to segment #%d and #%d was requested", lastSegment, segment)

		// otherwise, stop the existing transcoding process, restart at the applicable time
		sm.stopTranscode(hash + resolution)
	}

	// no transcode processes exist now, so start a new one
	// start one at the applicable time, wait and return stream
	_, err := sm.startTranscode(src, hash, resolution, transcodeStartSegment)

	if err != nil {
		return onTranscodeError(err)
	}

	return sm.streamTSFunc(hash+resolution, segment)
}

func (sm *StreamManager) segmentToTime(segment int) string {
	return fmt.Sprint(segment * hlsSegmentLength)
}

func (sm *StreamManager) getTranscodeArgs(probeResult *VideoFile, outputPath string, segment int, resolution string) []string {
	args := []string{
		"-loglevel", "error",
	}

	if segment > 0 {
		args = append(args,
			"-ss", sm.segmentToTime(segment),
			// without this ffmpeg would sometimes generate empty ts files when seeking
			"-noaccurate_seek",
		)
	}

	args = append(args, "-i", probeResult.Path)

	scale := fmt.Sprintf("-2:%s", resolution)
	args = append(args,
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-profile:v", "high",
		"-level", "4.2",
		"-preset", "veryfast",
		"-crf", "23",
		"-r", "30",
		"-g", "60",
		"-x264-params", "no-scenecut=1",
		"-force_key_frames", fmt.Sprintf("expr:gte(t,n_forced*%d)", 1),
		"-vf", "scale="+scale)

	args = append(args,
		"-c:a", "aac",
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
		"-hls_flags", "temp_file",
		filepath.Join(outputPath, "playlist.m3u8"),
	)

	return args
}

// assumes mutex is held
func (sm *StreamManager) startTranscode(src string, hash string, resolution string, segment int) (*transcodeProcess, error) {
	probeResult, err := sm.ffprobe.NewVideoFile(src)
	if err != nil {
		return nil, err
	}

	outputPath := sm.segmentDirectory(hash + resolution)
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return nil, err
	}

	lockCtx := sm.lockManager.ReadLock(sm.context, probeResult.Path)

	args := sm.getTranscodeArgs(probeResult, outputPath, segment, resolution)
	cmd := sm.ffmpeg.Command(lockCtx, args)

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
		lockCtx.Cancel()
		return nil, err
	}

	p := &transcodeProcess{
		cmd:          cmd,
		context:      lockCtx,
		path:         probeResult.Path,
		cancel:       lockCtx.Cancel,
		hash:         hash,
		resolution:   resolution,
		startSegment: segment,
	}
	sm.runningTranscodes[hash+resolution] = p

	// mark the stream as accessed to ensure it is not immediately cleaned up
	sm.streamsMutex.Lock()
	defer sm.streamsMutex.Unlock()
	sm.runningStreams[hash+resolution] = runningStream{
		lastAccessed: time.Now(),
		segment:      segment,
	}

	go sm.waitAndDeregister(hash+resolution, p, stdout, stderr)

	return p, nil
}

// assumes mutex is held
func (sm *StreamManager) stopTranscode(hashResolution string) {
	p := sm.runningTranscodes[hashResolution]
	if p != nil {
		p.cancelled = true
		// Windows doesn't support Interrupt
		if runtime.GOOS != "windows" {
			_ = p.cmd.Process.Signal(os.Interrupt)
			select {
			case <-p.context.Done():
				logger.Debug("ffmpeg process exited cleanly")
				break
			case <-time.After(cancelTimeout):
				logger.Warn("ffmpeg process exited uncleanly")
				break
			}
		}
		p.cancel()
		delete(sm.runningTranscodes, hashResolution)
	}
}

func (sm *StreamManager) waitAndDeregister(hashResolution string, p *transcodeProcess, stdout, stderr io.Reader) {
	cmd := p.cmd

	errStr, _ := io.ReadAll(stderr)
	outStr, _ := io.ReadAll(stdout)

	err := cmd.Wait()

	// make sure that cancel is called to prevent memory leaks
	p.cancel()

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
	if sm.runningTranscodes[hashResolution] == p {
		delete(sm.runningTranscodes, hashResolution)
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
	// lock both streams and transcodes mutex to prevent deadlock when
	// starting a transcode (uses both streams and transcodes) and removing stale files
	// must call transcode mutex before streams mutex as transcode mutex is locked first
	// and will result in a deadlock if a request is made while removing files
	sm.transcodesMutex.Lock()
	defer sm.transcodesMutex.Unlock()
	sm.streamsMutex.Lock()
	defer sm.streamsMutex.Unlock()

	var toRemove []string

	now := time.Now()
	for hashResolution, stream := range sm.runningStreams {
		if !stream.running && stream.lastAccessed.Add(maxIdleTime).Before(now) {
			// Stream expired. Cancel the transcode process and delete the files
			logger.Debugf("stream for hashResolution %q not accessed recently. Cancelling transcode and removing files", hashResolution)

			sm.stopAndRemoveTranscodeFiles(hashResolution)

			toRemove = append(toRemove, hashResolution)
		} else {
			// check if the last transcoded file is way ahead of the current streaming one
			// if so, stop the transcode to save CPU
			lastGenerated := sm.lastTranscodedSegment(hashResolution)

			// prevent stoppping the stream if we just scrubbed to it
			if sm.runningTranscodes[hashResolution] != nil &&
				sm.runningTranscodes[hashResolution].startSegment+minSegmentTranscode < stream.segment &&
				stream.segment+maxSegmentStopGap < lastGenerated {
				logger.Debugf("stopping transcode for hashResolution %q as last generated segment %d is too far ahead of current segment %d", hashResolution, lastGenerated, stream.segment)
				sm.stopTranscode(hashResolution)
			}
		}
	}

	for _, hashResolution := range toRemove {
		delete(sm.runningStreams, hashResolution)
	}
}

// stopAndRemoveAll stops all current streams and removes all cache files
func (sm *StreamManager) stopAndRemoveAll() {
	sm.transcodesMutex.Lock()
	defer sm.transcodesMutex.Unlock()
	sm.streamsMutex.Lock()
	defer sm.streamsMutex.Unlock()

	for hashResolution := range sm.runningStreams {
		sm.stopAndRemoveTranscodeFiles(hashResolution)
	}

	// ensure nothing else can use the map
	sm.runningStreams = nil
	sm.runningTranscodes = nil
}

// assume lock is held
func (sm *StreamManager) stopAndRemoveTranscodeFiles(hashResolution string) {
	sm.stopTranscode(hashResolution)

	dir := sm.segmentDirectory(hashResolution)
	if err := os.RemoveAll(dir); err != nil {
		logger.Warnf("error removing segment directory %q: %v", dir, err)
	}
}
