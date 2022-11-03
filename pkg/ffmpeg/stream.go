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
	"sync/atomic"
	"time"

	"github.com/stashapp/stash/pkg/file"
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

	maxSegmentWait  = 15 * time.Second
	monitorInterval = 200 * time.Millisecond

	// segment gap before counting a request as a seek and
	// restarting the transcode process at the requested segment
	maxSegmentGap = 5

	// maximum number of segments to generate
	// ahead of the currently streaming segment
	maxSegmentBuffer = 15

	// maximum idle time between segment requests before
	// stopping transcode and deleting cache folder
	maxIdleTime = 30 * time.Second

	// Cancel timeout for ffmpeg to prevent corrupted segments
	cancelTimeout = 3 * time.Second
)

type StreamType string
type SegmentType string

const (
	StreamTypeDASH    StreamType = "dash"
	StreamTypeHLS     StreamType = "hls"
	StreamTypeHLSCopy StreamType = "hls-copy"

	SegmentTypeWEBMVideo SegmentType = "webm-v"
	SegmentTypeWEBMAudio SegmentType = "webm-a"
	SegmentTypeTS        SegmentType = "mpegts"
)

var ErrInvalidSegment = errors.New("invalid segment")

// TranscodeStreamOptions represents options for live transcoding a video file.
type TranscodeStreamOptions struct {
	StreamType  StreamType
	SegmentType SegmentType
	VideoFile   *file.VideoFile
	Hash        string
	Resolution  string
}

type transcodeProcess struct {
	cmd          *exec.Cmd
	context      context.Context
	cancel       context.CancelFunc
	cancelled    bool
	startSegment int
}

type waitingSegment struct {
	segmentType SegmentType
	idx         int
	file        string
	path        string
	accessed    time.Time
	available   chan error
	done        atomic.Bool
}

type runningStream struct {
	dir              string
	streamType       StreamType
	vf               *file.VideoFile
	maxTranscodeSize int
	outputDir        string

	waitingSegments []*waitingSegment
	tp              *transcodeProcess
	lastAccessed    time.Time
	lastSegment     int
}

type StreamManager struct {
	cacheDir string
	encoder  FFMpeg
	ffprobe  FFProbe

	config      StreamManagerConfig
	lockManager *fsutil.ReadLockManager

	context    context.Context
	cancelFunc context.CancelFunc

	runningStreams map[string]*runningStream
	streamsMutex   sync.Mutex
}

type StreamManagerConfig interface {
	GetMaxStreamingTranscodeSize() models.StreamingResolutionEnum
}

func (t StreamType) MainSegmentType() SegmentType {
	switch t {
	case StreamTypeDASH:
		return SegmentTypeWEBMVideo
	case StreamTypeHLS, StreamTypeHLSCopy:
		return SegmentTypeTS
	}
	return ""
}

func (t StreamType) MimeType() string {
	switch t {
	case StreamTypeDASH:
		return MimeDASH
	case StreamTypeHLS, StreamTypeHLSCopy:
		return MimeHLS
	}
	return ""
}

func (t StreamType) WriteManifest(sm *StreamManager, w io.Writer, vf *file.VideoFile, baseURL, resolution string) error {
	switch t {
	case StreamTypeDASH:
		return sm.WriteDASHManifest(w, vf, baseURL, resolution)
	case StreamTypeHLS, StreamTypeHLSCopy:
		return sm.WriteHLSManifest(w, vf, baseURL, resolution)
	}
	return fmt.Errorf("invalid stream type %s", t)
}

func (t StreamType) Args(vf *file.VideoFile, maxScale, segment int, outputDir string) (args Args) {
	args = append(args, "-hide_banner")
	args = args.LogLevel(LogLevelError)

	if segment > 0 {
		args = args.Seek(segmentToTime(segment))
	}

	args = args.Input(vf.Path)

	videoOnly := ProbeAudioCodec(vf.AudioCodec) == MissingUnsupported

	var videoFilter VideoFilter
	videoFilter = videoFilter.ScaleMax(vf.Width, vf.Height, maxScale)

	switch t {
	case StreamTypeDASH:
		args = append(args, []string{
			"-c:v", "libvpx-vp9",
			"-pix_fmt", "yuv420p",
			"-deadline", "realtime",
			"-cpu-used", "5",
			"-row-mt", "1",
			"-crf", "30",
			"-b:v", "0",
			"-force_key_frames", fmt.Sprintf("expr:gte(t,n_forced*%d)", segmentLength),
		}...)
		args = args.VideoFilter(videoFilter)
		args = append(args, []string{
			"-copyts",
			"-avoid_negative_ts", "disabled",
			"-map", "0:v:0",
			"-f", "webm_chunk",
			"-chunk_start_index", fmt.Sprint(segment),
			"-header", filepath.Join(outputDir, "init_v.webm"),
			filepath.Join(outputDir, "%d_v.webm"),
		}...)
		if !videoOnly {
			args = append(args, []string{
				"-c:a", "libopus",
				"-b:a", "96000",
				"-ar", "48000",
				"-copyts",
				"-avoid_negative_ts", "disabled",
				"-map", "0:a:0",
				"-f", "webm_chunk",
				"-chunk_start_index", fmt.Sprint(segment),
				"-audio_chunk_duration", fmt.Sprint(segmentLength * 1000),
				"-header", filepath.Join(outputDir, "init_a.webm"),
				filepath.Join(outputDir, "%d_a.webm"),
			}...)
		}
	case StreamTypeHLS:
		args = append(args, []string{
			"-c:v", "libx264",
			"-pix_fmt", "yuv420p",
			"-preset", "veryfast",
			"-crf", "25",
			"-flags", "+cgop",
			"-force_key_frames", fmt.Sprintf("expr:gte(t,n_forced*%d)", segmentLength),
		}...)
		args = args.VideoFilter(videoFilter)
		if videoOnly {
			args = append(args, "-an")
		} else {
			args = append(args, []string{
				"-c:a", "aac",
				"-ac", "2",
			}...)
		}
		args = append(args, []string{
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
		}...)
	case StreamTypeHLSCopy:
		args = append(args, []string{
			"-c:v", "copy",
		}...)
		if videoOnly {
			args = append(args, "-an")
		} else {
			args = append(args, []string{
				"-c:a", "aac",
				"-ac", "2",
			}...)
		}
		args = append(args, []string{
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
		}...)
	}

	return
}

func (t StreamType) FileDir(hash string, maxTranscodeSize int) string {
	if maxTranscodeSize == 0 {
		return fmt.Sprintf("%s_%s", hash, t)
	} else {
		return fmt.Sprintf("%s_%s_%d", hash, t, maxTranscodeSize)
	}
}

func (t SegmentType) MimeType() string {
	switch t {
	case SegmentTypeWEBMVideo:
		return MimeMp4Video
	case SegmentTypeWEBMAudio:
		return MimeMp4Audio
	case SegmentTypeTS:
		return MimeMpegTS
	}
	return ""
}

func (t SegmentType) Segment(str string) (segment int, err error) {
	switch t {
	case SegmentTypeWEBMVideo, SegmentTypeWEBMAudio:
		if str == "init" {
			segment = -1
		} else {
			segment, err = strconv.Atoi(str)
			if err != nil || segment < 0 {
				err = ErrInvalidSegment
			}
		}
	case SegmentTypeTS:
		segment, err = strconv.Atoi(str)
		if err != nil || segment < 0 {
			err = ErrInvalidSegment
		}
	}
	return
}

func (t SegmentType) FileName(segment int) string {
	switch t {
	case SegmentTypeWEBMVideo:
		if segment == -1 {
			return "init_v.webm"
		} else {
			return fmt.Sprintf("%d_v.webm", segment)
		}
	case SegmentTypeWEBMAudio:
		if segment == -1 {
			return "init_a.webm"
		} else {
			return fmt.Sprintf("%d_a.webm", segment)
		}
	case SegmentTypeTS:
		return fmt.Sprintf("%d.ts", segment)
	}
	return ""
}

func (s *runningStream) Args(segment int) Args {
	return s.streamType.Args(s.vf, s.maxTranscodeSize, segment, s.outputDir)
}

func lastSegment(vf *file.VideoFile) int {
	return int(math.Ceil(vf.Duration/segmentLength)) - 1
}

func segmentExists(path string) bool {
	exists, _ := fsutil.FileExists(path)
	return exists
}

func segmentToTime(segment int) float64 {
	return float64(segment * segmentLength)
}

func NewStreamManager(cacheDir string, encoder FFMpeg, ffprobe FFProbe, config StreamManagerConfig, lockManager *fsutil.ReadLockManager) *StreamManager {
	if cacheDir == "" {
		panic("cache directory is not set")
	}

	ctx, cancel := context.WithCancel(context.Background())

	ret := &StreamManager{
		cacheDir:       cacheDir,
		encoder:        encoder,
		ffprobe:        ffprobe,
		config:         config,
		lockManager:    lockManager,
		context:        ctx,
		cancelFunc:     cancel,
		runningStreams: make(map[string]*runningStream),
	}

	go func() {
		for {
			select {
			case <-time.After(monitorInterval):
				ret.monitorStreams()
			case <-ctx.Done():
				ret.stopAndRemoveAll()
				return
			}
		}
	}()

	return ret
}

// Shutdown shuts down the stream manager, killing any running transcoding processes and removing all cached files.
func (sm *StreamManager) Shutdown() {
	sm.cancelFunc()
}

// WriteHLSManifest writes an HLS playlist manifest to w. The URLs for the segments
// are of the form {baseURL}/%d.ts{?urlQuery} where %d is the segment index.
func (sm *StreamManager) WriteHLSManifest(w io.Writer, vf *file.VideoFile, baseURL, resolution string) error {
	probeResult, err := sm.ffprobe.NewVideoFile(vf.Path)
	if err != nil {
		return err
	}

	var urlQuery string
	if resolution != "" {
		urlQuery = fmt.Sprintf("?resolution=%s", resolution)
	}

	fmt.Fprint(w, "#EXTM3U\n")

	fmt.Fprint(w, "#EXT-X-VERSION:3\n")
	fmt.Fprint(w, "#EXT-X-MEDIA-SEQUENCE:0\n")
	fmt.Fprintf(w, "#EXT-X-TARGETDURATION:%d\n", segmentLength)
	fmt.Fprint(w, "#EXT-X-PLAYLIST-TYPE:VOD\n")

	leftover := probeResult.Duration
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

	return nil
}

// WriteDASHManifest writes an DASH playlist manifest to w. The base URL of the playlist is set to baseURL.
func (sm *StreamManager) WriteDASHManifest(w io.Writer, vf *file.VideoFile, baseURL, resolution string) error {
	probeResult, err := sm.ffprobe.NewVideoFile(vf.Path)
	if err != nil {
		return err
	}

	var framerate string
	var videoWidth int
	var videoHeight int
	videoStream := probeResult.VideoStream
	if videoStream != nil {
		framerate = videoStream.AvgFrameRate
		videoWidth = videoStream.Width
		videoHeight = videoStream.Height
	} else {
		// extract the framerate fraction from the file framerate
		// framerates 0.1% below round numbers are common,
		// attempt to infer when this is the case
		fileFramerate := vf.FrameRate
		rate1001, off1001 := math.Modf(fileFramerate * 1.001)
		var numerator int
		var denominator int
		switch {
		case off1001 < 0.005:
			numerator = int(rate1001) * 1000
			denominator = 1001
		case off1001 > 0.995:
			numerator = (int(rate1001) + 1) * 1000
			denominator = 1001
		default:
			numerator = int(fileFramerate * 1000)
			denominator = 1000
		}
		framerate = fmt.Sprintf("%d/%d", numerator, denominator)
		videoHeight = vf.Height
		videoWidth = vf.Width
	}

	var urlQuery string
	maxTranscodeSize := sm.config.GetMaxStreamingTranscodeSize().GetMaxResolution()
	if resolution != "" {
		maxTranscodeSize = models.StreamingResolutionEnum(resolution).GetMaxResolution()
		urlQuery = fmt.Sprintf("?resolution=%s", resolution)
	}
	if maxTranscodeSize != 0 {
		videoSize := videoHeight
		if videoWidth < videoSize {
			videoSize = videoWidth
		}

		if maxTranscodeSize < videoSize {
			scaleFactor := float64(maxTranscodeSize) / float64(videoSize)
			videoWidth = int(float64(videoWidth) * scaleFactor)
			videoHeight = int(float64(videoHeight) * scaleFactor)
		}
	}

	mediaDuration := mpd.Duration(time.Duration(probeResult.Duration * float64(time.Second)))
	m := mpd.NewMPD(mpd.DASH_PROFILE_LIVE, mediaDuration.String(), "PT4.0S")
	m.BaseURL = baseURL + "/"

	video, _ := m.AddNewAdaptationSetVideo("video/webm", "progressive", true, 1)

	_, _ = video.SetNewSegmentTemplate(2, "init_v.webm"+urlQuery, "$Number$_v.webm"+urlQuery, 0, 1)
	_, _ = video.AddNewRepresentationVideo(200000, "vp09.00.40.08", "0", framerate, int64(videoWidth), int64(videoHeight))

	if ProbeAudioCodec(vf.AudioCodec) != MissingUnsupported {
		audio, _ := m.AddNewAdaptationSetAudio("audio/webm", true, 1, "und")
		_, _ = audio.SetNewSegmentTemplate(2, "init_a.webm"+urlQuery, "$Number$_a.webm"+urlQuery, 0, 1)
		_, _ = audio.AddNewRepresentationAudio(48000, 96000, "opus", "1")
	}

	return m.Write(w)
}

func (sm *StreamManager) streamSegmentFunc(segment *waitingSegment) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			break
		case err := <-segment.available:
			if err == nil {
				logger.Tracef("[transcode] streaming segment file %s", segment.file)
				w.Header().Set("Content-Type", segment.segmentType.MimeType())
				// Prevent caching as segments are generated on the fly
				w.Header().Add("Cache-Control", "no-cache")
				http.ServeFile(w, r, segment.path)
			} else if !errors.Is(err, context.Canceled) {
				logger.Errorf("[transcode] %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		segment.done.Store(true)
	}
}

// StreamSegment returns a http.HandlerFunc that streams a segment.
func (sm *StreamManager) StreamSegment(options TranscodeStreamOptions, segmentStr string) http.HandlerFunc {
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

	segment, err := options.SegmentType.Segment(segmentStr)
	// error if segment is past the end of the video
	if err != nil || segment > lastSegment(options.VideoFile) {
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "invalid segment", http.StatusBadRequest)
		}
	}

	maxTranscodeSize := sm.config.GetMaxStreamingTranscodeSize().GetMaxResolution()
	if options.Resolution != "" {
		maxTranscodeSize = models.StreamingResolutionEnum(options.Resolution).GetMaxResolution()
	}

	dir := options.StreamType.FileDir(options.Hash, maxTranscodeSize)
	outputDir := filepath.Join(sm.cacheDir, dir)

	name := options.SegmentType.FileName(segment)
	file := filepath.Join(dir, name)

	sm.streamsMutex.Lock()

	stream := sm.runningStreams[dir]
	if stream == nil {
		stream = &runningStream{
			dir:              dir,
			streamType:       options.StreamType,
			vf:               options.VideoFile,
			maxTranscodeSize: maxTranscodeSize,
			outputDir:        outputDir,

			// initialize to cap 10 to avoid reallocations
			waitingSegments: make([]*waitingSegment, 0, 10),
		}
		sm.runningStreams[dir] = stream
	}

	now := time.Now()
	stream.lastAccessed = now
	if segment != -1 {
		stream.lastSegment = segment
	}

	waitingSegment := &waitingSegment{
		segmentType: options.SegmentType,
		idx:         segment,
		file:        file,
		path:        filepath.Join(sm.cacheDir, file),
		accessed:    now,
		available:   make(chan error, 1),
	}
	stream.waitingSegments = append(stream.waitingSegments, waitingSegment)

	sm.streamsMutex.Unlock()

	return sm.streamSegmentFunc(waitingSegment)
}

// assume lock is held
func (sm *StreamManager) startTranscode(stream *runningStream, segment int, done chan<- error) {
	// generate segment 0 if init segment requested
	if segment == -1 {
		segment = 0
	}

	logger.Debugf("[transcode] starting transcode for %s at segment #%d", stream.dir, segment)

	if err := os.MkdirAll(stream.outputDir, os.ModePerm); err != nil {
		done <- err
		return
	}

	lockCtx := sm.lockManager.ReadLock(sm.context, stream.vf.Path)

	args := stream.Args(segment)
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
		done <- fmt.Errorf("error starting transcode process: %w", err)
		return
	}

	tp := &transcodeProcess{
		cmd:          cmd,
		context:      lockCtx,
		cancel:       lockCtx.Cancel,
		startSegment: segment,
	}
	stream.tp = tp

	go func() {
		errStr, _ := io.ReadAll(stderr)
		outStr, _ := io.ReadAll(stdout)

		errCmd := cmd.Wait()

		var err error

		// don't log error if cancelled
		if !tp.cancelled {
			e := string(errStr)
			if e == "" {
				e = string(outStr)
			}
			if e != "" {
				err = errors.New(e)
			} else {
				err = errCmd
			}

			if err != nil {
				err = fmt.Errorf("[transcode] ffmpeg error when running command <%s>: %w", strings.Join(cmd.Args, " "), err)
			}
		}

		sm.streamsMutex.Lock()

		// make sure that cancel is called to prevent memory leaks
		tp.cancel()
		if stream.tp == tp {
			stream.tp = nil
		}

		sm.streamsMutex.Unlock()

		done <- err
	}()
}

// assume lock is held
func (sm *StreamManager) stopTranscode(stream *runningStream) {
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

func (s *waitingSegment) checkAvailable(now time.Time) bool {
	if segmentExists(s.path) {
		s.available <- nil
		return true
	} else if s.accessed.Add(maxSegmentWait).Before(now) {
		s.available <- fmt.Errorf("timed out waiting for segment file %s to be generated", s.file)
		return true
	}
	return false
}

// ensureTranscode will start a new transcode process if the transcode
// is more than maxSegmentGap behind the requested segment
func (sm *StreamManager) ensureTranscode(stream *runningStream, segment *waitingSegment) bool {
	segmentIdx := segment.idx
	tp := stream.tp
	if tp == nil {
		sm.startTranscode(stream, segmentIdx, segment.available)
		return true
	} else {
		segmentType := segment.segmentType
		bufStart := tp.startSegment
		var bufEnd int
		for i := bufStart + 1; ; i++ {
			if !segmentExists(filepath.Join(stream.outputDir, segmentType.FileName(i))) {
				bufEnd = i - 1
				break
			}
		}
		if segmentIdx < bufStart || bufEnd+maxSegmentGap < segmentIdx {
			sm.stopTranscode(stream)
			sm.startTranscode(stream, segmentIdx, segment.available)
			return true
		}
	}
	return false
}

func (sm *StreamManager) monitorStreams() {
	sm.streamsMutex.Lock()
	defer sm.streamsMutex.Unlock()

	now := time.Now()

	for dir, stream := range sm.runningStreams {
		transcodeStarted := false
		temp := stream.waitingSegments[:0]
		for _, segment := range stream.waitingSegments {
			var remove bool
			if segment.done.Load() || segment.checkAvailable(now) {
				remove = true
			} else if !transcodeStarted {
				transcodeStarted = sm.ensureTranscode(stream, segment)
			}
			if !remove {
				temp = append(temp, segment)
			}
		}
		stream.waitingSegments = temp

		if !transcodeStarted {
			if len(stream.waitingSegments) == 0 && stream.lastAccessed.Add(maxIdleTime).Before(now) {
				// Stream expired. Cancel the transcode process and delete the files
				logger.Debugf("[transcode] stream for %s not accessed recently. Cancelling transcode and removing files", dir)

				sm.stopTranscode(stream)
				sm.removeTranscodeFiles(stream)

				delete(sm.runningStreams, dir)
			}
			if stream.tp != nil {
				segmentType := stream.streamType.MainSegmentType()
				segment := stream.lastSegment
				// if all segments up to maxSegmentBuffer exist, stop transcode
				for i := segment; i < segment+maxSegmentBuffer; i++ {
					if !segmentExists(filepath.Join(stream.outputDir, segmentType.FileName(i))) {
						return
					}
				}

				logger.Debugf("[transcode] stopping transcode for %s, buffer is full", dir)
				sm.stopTranscode(stream)
			}
		}
	}
}

// assume lock is held
func (sm *StreamManager) removeTranscodeFiles(stream *runningStream) {
	path := stream.outputDir
	if err := os.RemoveAll(path); err != nil {
		logger.Warnf("[transcode] error removing segment directory %s: %v", path, err)
	}
}

// stopAndRemoveAll stops all current streams and removes all cache files
func (sm *StreamManager) stopAndRemoveAll() {
	sm.streamsMutex.Lock()
	defer sm.streamsMutex.Unlock()

	for _, stream := range sm.runningStreams {
		for _, segment := range stream.waitingSegments {
			if len(segment.available) == 0 {
				segment.available <- context.Canceled
			}
		}
		sm.stopTranscode(stream)
		sm.removeTranscodeFiles(stream)
	}

	// ensure nothing else can use the map
	sm.runningStreams = nil
}
