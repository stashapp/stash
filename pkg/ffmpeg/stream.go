package ffmpeg

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"

	"github.com/zencoder/go-dash/v3/mpd"
)

const (
	MimeHLS      string = "application/vnd.apple.mpegurl"
	MimeMpegTS   string = "video/MP2T"
	MimeDASH     string = "application/dash+xml"
	MimeMp4Video string = "video/mp4"
	MimeMp4Audio string = "audio/mp4"

	segmentLength = 2

	maxSegmentWait       = 15 * time.Second
	segmentCheckInterval = 100 * time.Millisecond
	monitorInterval      = 5 * time.Second

	// segment gap before counting a request as a seek and
	// restarting the transcode process at the requested segment
	segmentSeekGap = 5

	// minimum number of segments to generate
	// ahead of the currently streaming segment
	minSegmentBuffer = 3

	// maximum number of segments to generate
	// ahead of the currently streaming segment
	maxSegmentBuffer = 15

	// maximum idle time between segment requests before
	// stopping transcode and deleting cache folder
	maxIdleTime = 30 * time.Second

	// Cancel timeout for ffmpeg to prevent corrupted segments
	cancelTimeout = 2 * time.Second
)

// StreamType represents a transcode stream codec.

type StreamType string

const (
	StreamTypeDASHVideo StreamType = "dash-v"
	StreamTypeDASHAudio StreamType = "dash-a"
	StreamTypeHLS       StreamType = "hls"
	StreamTypeHLSCopy   StreamType = "hls-copy"
)

var ErrInvalidSegment = errors.New("invalid segment")

// TranscodeStreamOptions represents options for live transcoding a video file.
type TranscodeStreamOptions struct {
	Type StreamType

	Input            string
	Hash             string
	MaxTranscodeSize int

	// original video metadata
	VideoDuration float64
	VideoWidth    int
	VideoHeight   int

	// transcode the video, remove the audio
	// in some videos where the audio codec is not supported by ffmpeg
	// ffmpeg fails if you try to transcode the audio
	VideoOnly bool
}

type StreamManagerConfig interface {
	GetMaxStreamingTranscodeSize() models.StreamingResolutionEnum
}

type transcodeProcess struct {
	cmd       *exec.Cmd
	context   context.Context
	cancel    context.CancelFunc
	cancelled bool
}

type runningStream struct {
	active       bool
	lastAccessed time.Time
	segment      int
	options      *TranscodeStreamOptions
	tp           *transcodeProcess
}

type StreamManager struct {
	cacheDir    string
	encoder     FFMpeg
	config      StreamManagerConfig
	lockManager *fsutil.ReadLockManager

	context    context.Context
	cancelFunc context.CancelFunc

	runningStreams map[string]*runningStream
	streamsMutex   sync.Mutex
}

type transcodedSegment struct {
	name string
	time time.Time
}

func (c StreamType) String() string {
	switch c {
	case StreamTypeDASHVideo, StreamTypeDASHAudio:
		return "dash"
	case StreamTypeHLS:
		return "hls"
	case StreamTypeHLSCopy:
		return "hls-copy"
	}
	return ""
}

func (c StreamType) MimeType() string {
	switch c {
	case StreamTypeDASHVideo:
		return MimeMp4Video
	case StreamTypeDASHAudio:
		return MimeMp4Audio
	case StreamTypeHLS, StreamTypeHLSCopy:
		return MimeMpegTS
	}
	return ""
}

func (c StreamType) VideoCodec() VideoCodec {
	switch c {
	case StreamTypeDASHVideo, StreamTypeDASHAudio:
		return VideoCodecVP9
	case StreamTypeHLS:
		return VideoCodecLibX264
	case StreamTypeHLSCopy:
		return VideoCodecCopy
	}
	return ""
}

func (c StreamType) Args(segment int, outputDir string) []string {
	// segment length in frames
	segmentFrames := fmt.Sprint(segmentLength * 30)

	switch c {
	case StreamTypeDASHVideo, StreamTypeDASHAudio:
		return []string{
			"-pix_fmt", "yuv420p",
			"-deadline", "realtime",
			"-cpu-used", "5",
			"-row-mt", "1",
			"-crf", "30",
			"-b:v", "0",
			"-r", "30",
			"-g", segmentFrames,
			"-keyint_min", segmentFrames,
			"-copyts",
			"-avoid_negative_ts", "disabled",
			"-map", "0:v:0",
			"-f", "webm_chunk",
			"-chunk_start_index", fmt.Sprint(segment),
			"-header", filepath.Join(outputDir, "init_v.webm"),
			filepath.Join(outputDir, "%d_v.webm"),
			"-c:a", "libopus",
			"-ar", "48000",
			"-copyts",
			"-avoid_negative_ts", "disabled",
			"-map", "0:a:0",
			"-f", "webm_chunk",
			"-chunk_start_index", fmt.Sprint(segment),
			"-audio_chunk_duration", "2000",
			"-header", filepath.Join(outputDir, "init_a.webm"),
			filepath.Join(outputDir, "%d_a.webm"),
		}
	case StreamTypeHLS:
		return []string{
			"-pix_fmt", "yuv420p",
			"-preset", "veryfast",
			"-crf", "25",
			"-r", "30",
			"-g", segmentFrames,
			"-keyint_min", segmentFrames,
			"-flags", "+cgop",
			"-force_key_frames", fmt.Sprintf("expr:gte(t,n_forced*%d)", segmentLength),
			"-c:a", "aac",
			"-ac", "2",
			"-sn",
			"-copyts",
			"-avoid_negative_ts", "disabled",
			"-strict", "-2",
			"-f", "hls",
			"-start_number", fmt.Sprint(segment),
			"-hls_time", "2",
			"-hls_segment_type", "mpegts",
			"-hls_playlist_type", "vod",
			"-hls_flags", "temp_file",
			"-hls_segment_filename", filepath.Join(outputDir, "%d.ts"),
			filepath.Join(outputDir, "manifest.m3u8"),
		}
	case StreamTypeHLSCopy:
		return []string{
			"-sn",
			"-copyts",
			"-avoid_negative_ts", "disabled",
			"-strict", "-2",
			"-f", "hls",
			"-start_number", fmt.Sprint(segment),
			"-hls_time", "2",
			"-hls_segment_type", "mpegts",
			"-hls_playlist_type", "vod",
			"-hls_flags", "temp_file",
			"-hls_segment_filename", filepath.Join(outputDir, "%d.ts"),
			filepath.Join(outputDir, "manifest.m3u8"),
		}
	}
	return []string{}
}

func (c StreamType) Segment(str string) (segment int, err error) {
	switch c {
	case StreamTypeDASHVideo, StreamTypeDASHAudio:
		if str == "init" {
			segment = -1
		} else {
			segment, err = strconv.Atoi(str)
			if err != nil || segment < 0 {
				err = ErrInvalidSegment
			}
		}
	case StreamTypeHLS, StreamTypeHLSCopy:
		segment, err = strconv.Atoi(str)
		if err != nil || segment < 0 {
			err = ErrInvalidSegment
		}
	}
	return
}

func (c StreamType) FileName(segment int) string {
	switch c {
	case StreamTypeDASHVideo:
		if segment == -1 {
			return "init_v.webm"
		} else {
			return fmt.Sprintf("%d_v.webm", segment)
		}
	case StreamTypeDASHAudio:
		if segment == -1 {
			return "init_a.webm"
		} else {
			return fmt.Sprintf("%d_a.webm", segment)
		}
	case StreamTypeHLS, StreamTypeHLSCopy:
		return fmt.Sprintf("%d.ts", segment)
	}
	return ""
}

func (c StreamType) FileSegment(filename string, segment *int) bool {
	switch c {
	case StreamTypeDASHVideo:
		n, _ := fmt.Sscanf(filename, "%d_v.webm", segment)
		if n != 0 {
			return true
		}
		if filename == "init_v.webm" {
			*segment = -1
			return true
		}
	case StreamTypeDASHAudio:
		n, _ := fmt.Sscanf(filename, "%d_a.webm", segment)
		if n != 0 {
			return true
		}
		if filename == "init_a.webm" {
			*segment = -1
			return true
		}
	case StreamTypeHLS, StreamTypeHLSCopy:
		n, _ := fmt.Sscanf(filename, "%d.ts", segment)
		if n != 0 {
			return true
		}
	}
	return false
}

func (o *TranscodeStreamOptions) FileDir() string {
	resolution := o.MaxTranscodeSize
	if resolution == 0 {
		return fmt.Sprintf("%s_%s", o.Hash, o.Type)
	} else {
		return fmt.Sprintf("%s_%s_%d", o.Hash, o.Type, resolution)
	}
}

func (o *TranscodeStreamOptions) FilePath(segment int) string {
	return filepath.Join(o.FileDir(), o.Type.FileName(segment))
}

func (o *TranscodeStreamOptions) LastSegment() int {
	return int(math.Ceil(o.VideoDuration/segmentLength)) - 1
}

func NewStreamManager(cacheDir string, encoder FFMpeg, config StreamManagerConfig, lockManager *fsutil.ReadLockManager) *StreamManager {
	if cacheDir == "" {
		panic("cache directory is not set")
	}

	ctx, cancel := context.WithCancel(context.Background())

	ret := &StreamManager{
		cacheDir:       cacheDir,
		encoder:        encoder,
		config:         config,
		lockManager:    lockManager,
		context:        ctx,
		cancelFunc:     cancel,
		runningStreams: make(map[string]*runningStream),
	}

	go ret.monitorStreams()

	return ret
}

// Shutdown shuts down the stream manager, killing any running transcoding processes and removing all cached files.
func (sm *StreamManager) Shutdown() {
	sm.cancelFunc()
}

// WriteHLSManifest writes an HLS playlist manifest to w. The URLs for the segments
// are of the form {baseURL}/%d.ts{?urlQuery} where %d is the segment index.
func (sm *StreamManager) WriteHLSManifest(w io.Writer, duration float64, baseURL, urlQuery string) {
	if urlQuery != "" {
		urlQuery = "?" + urlQuery
	}

	fmt.Fprint(w, "#EXTM3U\n")

	fmt.Fprint(w, "#EXT-X-VERSION:3\n")
	fmt.Fprint(w, "#EXT-X-MEDIA-SEQUENCE:0\n")
	fmt.Fprintf(w, "#EXT-X-TARGETDURATION:%d\n", segmentLength)
	fmt.Fprint(w, "#EXT-X-PLAYLIST-TYPE:VOD\n")

	leftover := duration
	segment := 0

	for leftover > 0 {
		thisLength := float64(segmentLength)
		if leftover < thisLength {
			thisLength = leftover
		}

		fmt.Fprintf(w, "#EXTINF:%f,\n", thisLength)
		fmt.Fprintf(w, "%s/%d.ts%s\n", baseURL, segment, urlQuery)

		leftover -= thisLength
		segment++
	}

	fmt.Fprint(w, "#EXT-X-ENDLIST\n")
}

// WriteDASHManifest writes an DASH playlist manifest to w. The base URL of the playlist is set to baseURL.
func (sm *StreamManager) WriteDASHManifest(w io.Writer, duration float64, baseURL, urlQuery string) {
	mediaDuration := mpd.Duration(time.Duration(duration) * time.Second)
	m := mpd.NewMPD(mpd.DASH_PROFILE_LIVE, mediaDuration.String(), "PT4.0S")
	m.BaseURL = baseURL + "/"

	if urlQuery != "" {
		urlQuery = "?" + urlQuery
	}

	video, _ := m.AddNewAdaptationSetVideo("video/webm", "progressive", true, 1)

	_, _ = video.SetNewSegmentTemplate(2, "init_v.webm"+urlQuery, "$Number$_v.webm"+urlQuery, 0, 1)
	_, _ = video.AddNewRepresentationVideo(200000, "vp09.00.40.08", "0", "30/1", 1920, 1080)

	audio, _ := m.AddNewAdaptationSetAudio("audio/webm", true, 1, "und")
	_, _ = audio.SetNewSegmentTemplate(2, "init_a.webm"+urlQuery, "$Number$_a.webm"+urlQuery, 0, 1)
	_, _ = audio.AddNewRepresentationAudio(48000, 96000, "opus", "1")

	_ = m.Write(w)
}

func segmentExists(path string) bool {
	exists, _ := fsutil.FileExists(path)
	return exists
}

func segmentToTime(segment int) float64 {
	return float64(segment * segmentLength)
}

func (sm *StreamManager) transcodedSegments(options *TranscodeStreamOptions) map[int]*transcodedSegment {
	ret := make(map[int]*transcodedSegment)

	files, _ := os.ReadDir(filepath.Join(sm.cacheDir, options.FileDir()))

	for _, f := range files {
		segment := 0
		if !options.Type.FileSegment(f.Name(), &segment) {
			continue
		}

		fileInfo, err := f.Info()
		if err != nil {
			continue
		}
		ret[segment] = &transcodedSegment{
			name: f.Name(),
			time: fileInfo.ModTime(),
		}
	}

	return ret
}

func (sm *StreamManager) streamSegmentFunc(stream *runningStream, options *TranscodeStreamOptions, segment int) http.HandlerFunc {
	file := options.FilePath(segment)
	path := filepath.Join(sm.cacheDir, file)

	return func(w http.ResponseWriter, r *http.Request) {
		go func() {
			<-r.Context().Done()

			sm.streamsMutex.Lock()
			if stream.options == options {
				stream.active = false
				stream.lastAccessed = time.Now()
			}
			sm.streamsMutex.Unlock()
		}()

		started := time.Now()
		for {
			select {
			case <-r.Context().Done():
				return
			case <-time.After(segmentCheckInterval):
				now := time.Now()
				switch {
				case segmentExists(path):
					logger.Tracef("[transcode] streaming segment file %s", file)
					w.Header().Set("Content-Type", options.Type.MimeType())
					// Prevent caching as segments are generated on the fly
					w.Header().Add("Cache-Control", "no-cache")
					http.ServeFile(w, r, path)
					return
				case started.Add(maxSegmentWait).Before(now):
					sm.streamsMutex.Lock()
					if stream.options == options {
						stream.active = false
					}
					sm.streamsMutex.Unlock()

					logger.Warnf("[transcode] timed out waiting for segment file %s to be generated", file)
					http.Error(w, "timed out waiting for segment file to be generated", http.StatusInternalServerError)
					return
				}
			}
		}
	}
}

// StreamSegment returns a http.HandlerFunc that streams a segment.
// If the segment exists in the cache directory, then it is streamed.
// Otherwise, a transcode process will be started for the provided segment. If
// a transcode process is running already, then it will be killed before the new
// process is started.
func (sm *StreamManager) StreamSegment(options *TranscodeStreamOptions, segmentStr string) http.HandlerFunc {
	if sm.cacheDir == "" {
		logger.Error("[transcode] cannot live transcode files because cache dir is empty")
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "cannot live transcode files because cache dir is empty", http.StatusInternalServerError)
		}
	}

	if options.Hash == "" {
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "invalid hash", http.StatusBadRequest)
		}
	}

	segment, err := options.Type.Segment(segmentStr)
	// error if segment is past the end of the video
	if err != nil || segment > options.LastSegment() {
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "invalid segment", http.StatusBadRequest)
		}
	}

	onTranscodeError := func(err error) http.HandlerFunc {
		errStr := fmt.Sprintf("error starting transcode process: %v", err)
		logger.Errorf("[transcode] %s", errStr)

		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, errStr, http.StatusInternalServerError)
		}
	}

	dir := options.FileDir()

	sm.streamsMutex.Lock()

	stream := sm.runningStreams[dir]
	if stream == nil {
		stream = &runningStream{}
		sm.runningStreams[dir] = stream
	}

	stream.active = true
	stream.lastAccessed = time.Now()
	stream.options = options

	if segment != -1 {
		// this is a seek if the requested segment is more than
		// segmentSeekGap away from the previously requested segment
		segmentGap := segment - stream.segment
		if segmentGap < 0 || segmentGap > segmentSeekGap {
			logger.Debugf("[transcode] seeking stream for %s to segment #%d", dir, segment)
			stream.stopTranscode()
		}

		stream.segment = segment
	}

	// if the transcode process is already running, just wait for the segment
	if stream.tp != nil {
		sm.streamsMutex.Unlock()
		return sm.streamSegmentFunc(stream, options, segment)
	}

	segments := sm.transcodedSegments(options)

	// if any segments up to minSegmentBuffer are missing, start transcode there
	for i := segment; i <= segment+minSegmentBuffer; i++ {
		if segments[i] == nil {
			stream.stopTranscode()
			err := sm.startTranscode(stream, i)
			if err != nil {
				stream.active = false
				sm.streamsMutex.Unlock()
				return onTranscodeError(err)
			}

			sm.streamsMutex.Unlock()
			return sm.streamSegmentFunc(stream, options, segment)
		}
	}

	// segment should exist, just stream it
	sm.streamsMutex.Unlock()
	return sm.streamSegmentFunc(stream, options, segment)
}

func (sm *StreamManager) getTranscodeArgs(options *TranscodeStreamOptions, segment int, outputDir string) []string {
	var args Args
	args = append(args, "-hide_banner")
	args.LogLevel(LogLevelError)

	if segment > 0 {
		args = args.Seek(segmentToTime(segment))
	}

	args = args.Input(options.Input)

	codec := options.Type.VideoCodec()
	codecArgs := options.Type.Args(segment, outputDir)

	args = args.VideoCodec(codec)

	// don't set scale when copying video stream
	if codec != VideoCodecCopy {
		var videoFilter VideoFilter
		videoFilter = videoFilter.ScaleMax(options.VideoWidth, options.VideoHeight, options.MaxTranscodeSize)
		args = args.VideoFilter(videoFilter)
	}

	if options.VideoOnly {
		args = args.SkipAudio()
	}

	args = append(args, codecArgs...)

	return args
}

// assume lock is held
func (sm *StreamManager) startTranscode(stream *runningStream, segment int) error {
	options := stream.options

	// generate segment 0 if init segment requested
	if segment == -1 {
		segment = 0
	}

	dir := options.FileDir()
	logger.Debugf("[transcode] starting transcode for %s at segment #%d", dir, segment)

	outputDir := filepath.Join(sm.cacheDir, dir)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	lockCtx := sm.lockManager.ReadLock(sm.context, options.Input)

	args := sm.getTranscodeArgs(options, segment, outputDir)
	cmd := sm.encoder.Command(lockCtx, args)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Errorf("[transcode] ffmpeg stderr not available: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if nil != err {
		logger.Errorf("[transcode] ffmpeg stdout not available: %v", err)
	}

	logger.Tracef("[transcode] running %s", cmd)
	if err := cmd.Start(); err != nil {
		lockCtx.Cancel()
		return err
	}

	tp := &transcodeProcess{
		cmd:     cmd,
		context: lockCtx,
		cancel:  lockCtx.Cancel,
	}
	stream.tp = tp

	go func() {
		errStr, _ := io.ReadAll(stderr)
		outStr, _ := io.ReadAll(stdout)

		err := cmd.Wait()

		// don't log error if cancelled
		if !tp.cancelled && err != nil {
			e := string(errStr)
			if e == "" {
				e = string(outStr)

				if e == "" {
					e = err.Error()
				}
			}

			// error message should be in the stderr stream
			logger.Errorf("[transcode] ffmpeg error when running command <%s>: %s", strings.Join(cmd.Args, " "), e)
		}

		sm.streamsMutex.Lock()

		// make sure that cancel is called to prevent memory leaks
		tp.cancel()
		if stream.tp == tp {
			stream.tp = nil
		}

		sm.streamsMutex.Unlock()

	}()

	return nil
}

// assume lock is held
func (stream *runningStream) stopTranscode() {
	tp := stream.tp
	if tp != nil {
		go func() {
			// Windows doesn't support Interrupt
			if runtime.GOOS != "windows" {
				_ = tp.cmd.Process.Signal(os.Interrupt)
				select {
				case <-tp.context.Done():
					logger.Trace("[transcode] ffmpeg process exited cleanly")
				case <-time.After(cancelTimeout):
					logger.Warn("[transcode] ffmpeg process exited uncleanly")
				}
			}
			tp.cancel()
		}()
		tp.cancelled = true
		stream.tp = nil
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

// check for the last time a stream was accessed
// and remove anything over a certain age
func (sm *StreamManager) removeStaleFiles() {
	sm.streamsMutex.Lock()
	defer sm.streamsMutex.Unlock()

	now := time.Now()
outer:
	for dir, stream := range sm.runningStreams {
		if !stream.active && stream.lastAccessed.Add(maxIdleTime).Before(now) {
			// Stream expired. Cancel the transcode process and delete the files
			logger.Debugf("[transcode] stream for %s not accessed recently. Cancelling transcode and removing files", dir)

			stream.stopTranscode()
			sm.removeTranscodeFiles(stream.options)

			delete(sm.runningStreams, dir)
		} else if stream.tp != nil {
			// if all segments up to maxSegmentBuffer exist, stop transcode
			segment := stream.segment
			segments := sm.transcodedSegments(stream.options)
			lastSegment := stream.options.LastSegment()
			for i := segment; i < segment+maxSegmentBuffer && i <= lastSegment; i++ {
				if segments[i] == nil {
					continue outer
				}
			}

			logger.Debugf("[transcode] stopping transcode for %s, buffer is full", dir)
			stream.stopTranscode()
		}
	}
}

// stopAndRemoveAll stops all current streams and removes all cache files
func (sm *StreamManager) stopAndRemoveAll() {
	sm.streamsMutex.Lock()
	defer sm.streamsMutex.Unlock()

	for _, stream := range sm.runningStreams {
		stream.stopTranscode()
		sm.removeTranscodeFiles(stream.options)
	}

	// ensure nothing else can use the map
	sm.runningStreams = nil
}

// assume lock is held
func (sm *StreamManager) removeTranscodeFiles(options *TranscodeStreamOptions) {
	path := filepath.Join(sm.cacheDir, options.FileDir())
	if err := os.RemoveAll(path); err != nil {
		logger.Warnf("[transcode] error removing segment directory %s: %v", path, err)
	}
}
