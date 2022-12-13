package ffmpeg

import (
	"errors"
	"io"
	"net/http"
	"syscall"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type StreamFormat struct {
	MimeType string
	Args     func(videoFilter VideoFilter, videoOnly bool) Args
}

var (
	StreamTypeMP4 = StreamFormat{
		MimeType: MimeMp4Video,
		Args: func(videoFilter VideoFilter, videoOnly bool) (args Args) {
			args = args.VideoCodec(VideoCodecLibX264)
			args = append(args,
				"-movflags", "frag_keyframe+empty_moov",
				"-pix_fmt", "yuv420p",
				"-preset", "veryfast",
				"-crf", "25",
			)
			args = args.VideoFilter(videoFilter)
			if videoOnly {
				args = args.SkipAudio()
			} else {
				args = append(args, "-ac", "2")
			}
			args = args.Format(FormatMP4)
			return
		},
	}
	StreamTypeWEBM = StreamFormat{
		MimeType: MimeWebmVideo,
		Args: func(videoFilter VideoFilter, videoOnly bool) (args Args) {
			args = args.VideoCodec(VideoCodecVP9)
			args = append(args,
				"-pix_fmt", "yuv420p",
				"-deadline", "realtime",
				"-cpu-used", "5",
				"-row-mt", "1",
				"-crf", "30",
				"-b:v", "0",
			)
			args = args.VideoFilter(videoFilter)
			if videoOnly {
				args = args.SkipAudio()
			} else {
				args = append(args, "-ac", "2")
			}
			args = args.Format(FormatWebm)
			return
		},
	}
	StreamTypeMKV = StreamFormat{
		MimeType: MimeMkvVideo,
		Args: func(videoFilter VideoFilter, videoOnly bool) (args Args) {
			args = args.VideoCodec(VideoCodecCopy)
			if videoOnly {
				args = args.SkipAudio()
			} else {
				args = args.AudioCodec(AudioCodecLibOpus)
				args = append(args,
					"-b:a", "96k",
					"-vbr", "on",
					"-ac", "2",
				)
			}
			args = args.Format(FormatMatroska)
			return
		},
	}
)

type TranscodeOptions struct {
	StreamType StreamFormat
	VideoFile  *file.VideoFile
	Resolution string
	StartTime  float64
}

func (o TranscodeOptions) makeStreamArgs(vf *file.VideoFile, maxScale int, startTime float64) Args {
	args := Args{"-hide_banner"}
	args = args.LogLevel(LogLevelError)

	if startTime != 0 {
		args = args.Seek(startTime)
	}

	args = args.Input(vf.Path)

	videoOnly := ProbeAudioCodec(vf.AudioCodec) == MissingUnsupported

	var videoFilter VideoFilter
	videoFilter = videoFilter.ScaleMax(vf.Width, vf.Height, maxScale)

	args = append(args, o.StreamType.Args(videoFilter, videoOnly)...)

	args = args.Output("pipe:")

	return args
}

func (sm *StreamManager) ServeTranscode(w http.ResponseWriter, r *http.Request, options TranscodeOptions) {
	streamRequestCtx := NewStreamRequestContext(w, r)
	lockCtx := sm.lockManager.ReadLock(streamRequestCtx, options.VideoFile.Path)

	// hijacking and closing the connection here causes video playback to hang in Chrome
	// due to ERR_INCOMPLETE_CHUNKED_ENCODING
	// We trust that the request context will be closed, so we don't need to call Cancel on the returned context here.

	handler, err := sm.getTranscodeStream(lockCtx, options)

	if err != nil {
		logger.Errorf("[transcode] error transcoding video file: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			logger.Warnf("[transcode] error writing response: %v", err)
		}
		return
	}

	handler(w, r)
}

func (sm *StreamManager) getTranscodeStream(ctx *fsutil.LockContext, options TranscodeOptions) (http.HandlerFunc, error) {
	maxTranscodeSize := sm.config.GetMaxStreamingTranscodeSize().GetMaxResolution()
	if options.Resolution != "" {
		maxTranscodeSize = models.StreamingResolutionEnum(options.Resolution).GetMaxResolution()
	}

	args := options.makeStreamArgs(options.VideoFile, maxTranscodeSize, options.StartTime)
	cmd := sm.encoder.Command(ctx, args)

	stdout, err := cmd.StdoutPipe()
	if nil != err {
		logger.Errorf("[transcode] ffmpeg stdout not available: %v", err)
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if nil != err {
		logger.Errorf("[transcode] ffmpeg stderr not available: %v", err)
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}
	ctx.AttachCommand(cmd)

	// stderr must be consumed or the process deadlocks
	go func() {
		stderrData, _ := io.ReadAll(stderr)
		stderrString := string(stderrData)
		if len(stderrString) > 0 {
			logger.Debugf("[transcode] ffmpeg stderr: %s", stderrString)
		}
	}()

	mimeType := options.StreamType.MimeType
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mimeType)
		w.WriteHeader(http.StatusOK)

		// process killing should be handled by command context

		_, err := io.Copy(w, stdout)
		if err != nil && !errors.Is(err, syscall.EPIPE) {
			logger.Errorf("[transcode] error serving transcoded video file: %v", err)
		}

		w.(http.Flusher).Flush()
	}
	return handler, nil
}
