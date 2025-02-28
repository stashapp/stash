package ffmpeg

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"syscall"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type StreamFormat struct {
	MimeType string
	Args     func(codec VideoCodec, videoFilter VideoFilter, videoOnly bool) Args
}

func CodecInit(codec VideoCodec) (args Args) {
	args = args.VideoCodec(codec)

	switch codec {
	// CPU Codecs
	case VideoCodecLibX264:
		args = append(args,
			"-pix_fmt", "yuv420p",
			"-preset", "veryfast",
			"-crf", "25",
			"-sc_threshold", "0",
		)
	case VideoCodecVP9:
		args = append(args,
			"-pix_fmt", "yuv420p",
			"-deadline", "realtime",
			"-cpu-used", "5",
			"-row-mt", "1",
			"-crf", "30",
			"-b:v", "0",
		)
	// HW Codecs
	case VideoCodecN264:
		args = append(args,
			"-rc", "vbr",
			"-cq", "15",
		)
	case VideoCodecN264H:
		args = append(args,
			"-profile", "p7",
			"-tune", "hq",
			"-profile", "high",
			"-rc", "vbr",
			"-rc-lookahead", "60",
			"-surfaces", "64",
			"-spatial-aq", "1",
			"-aq-strength", "15",
			"-cq", "15",
			"-coder", "cabac",
			"-b_ref_mode", "middle",
		)
	case VideoCodecI264, VideoCodecIVP9:
		args = append(args,
			"-global_quality", "20",
			"-preset", "faster",
		)
	case VideoCodecI264C:
		args = append(args,
			"-q", "20",
			"-preset", "faster",
		)
	case VideoCodecV264, VideoCodecVVP9:
		args = append(args,
			"-qp", "20",
		)
	case VideoCodecA264:
		args = append(args,
			"-quality", "speed",
		)
	case VideoCodecM264:
		args = append(args,
			"-realtime", "1",
		)
	case VideoCodecO264:
		args = append(args,
			"-preset", "superfast",
			"-crf", "25",
		)
	}

	return args
}

var (
	StreamTypeMP4 = StreamFormat{
		MimeType: MimeMp4Video,
		Args: func(codec VideoCodec, videoFilter VideoFilter, videoOnly bool) (args Args) {
			args = CodecInit(codec)
			args = append(args, "-movflags", "frag_keyframe+empty_moov")
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
		Args: func(codec VideoCodec, videoFilter VideoFilter, videoOnly bool) (args Args) {
			args = CodecInit(codec)
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
		Args: func(codec VideoCodec, videoFilter VideoFilter, videoOnly bool) (args Args) {
			args = CodecInit(codec)
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
	VideoFile  *models.VideoFile
	Resolution string
	StartTime  float64
}

func (o TranscodeOptions) FileGetCodec(sm *StreamManager, maxTranscodeSize int) (codec VideoCodec) {
	needsResize := false

	if maxTranscodeSize != 0 {
		if o.VideoFile.Width > o.VideoFile.Height {
			needsResize = o.VideoFile.Width > maxTranscodeSize
		} else {
			needsResize = o.VideoFile.Height > maxTranscodeSize
		}
	}

	switch o.StreamType.MimeType {
	case MimeMp4Video:
		if !needsResize && o.VideoFile.VideoCodec == H264 {
			return VideoCodecCopy
		}
		codec = VideoCodecLibX264
		if hwcodec := sm.encoder.hwCodecMP4Compatible(); hwcodec != nil && sm.config.GetTranscodeHardwareAcceleration() {
			codec = *hwcodec
		}
	case MimeWebmVideo:
		if !needsResize && (o.VideoFile.VideoCodec == Vp8 || o.VideoFile.VideoCodec == Vp9) {
			return VideoCodecCopy
		}
		codec = VideoCodecVP9
		if hwcodec := sm.encoder.hwCodecWEBMCompatible(); hwcodec != nil && sm.config.GetTranscodeHardwareAcceleration() {
			codec = *hwcodec
		}
	case MimeMkvVideo:
		codec = VideoCodecCopy
	}

	return codec
}

func (o TranscodeOptions) makeStreamArgs(sm *StreamManager) Args {
	maxTranscodeSize := sm.config.GetMaxStreamingTranscodeSize().GetMaxResolution()
	if o.Resolution != "" {
		maxTranscodeSize = models.StreamingResolutionEnum(o.Resolution).GetMaxResolution()
	}
	extraInputArgs := sm.config.GetLiveTranscodeInputArgs()
	extraOutputArgs := sm.config.GetLiveTranscodeOutputArgs()

	args := Args{"-hide_banner"}
	args = args.LogLevel(LogLevelError)

	codec := o.FileGetCodec(sm, maxTranscodeSize)

	fullhw := sm.config.GetTranscodeHardwareAcceleration() && sm.encoder.hwCanFullHWTranscode(sm.context, codec, o.VideoFile, maxTranscodeSize)
	args = sm.encoder.hwDeviceInit(args, codec, fullhw)
	args = append(args, extraInputArgs...)

	if o.StartTime != 0 {
		args = args.Seek(o.StartTime)
	}

	args = args.Input(o.VideoFile.Path)

	videoOnly := ProbeAudioCodec(o.VideoFile.AudioCodec) == MissingUnsupported

	videoFilter := sm.encoder.hwMaxResFilter(codec, o.VideoFile, maxTranscodeSize, fullhw)

	args = append(args, o.StreamType.Args(codec, videoFilter, videoOnly)...)

	args = append(args, extraOutputArgs...)

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
		// don't log context canceled errors
		if !errors.Is(err, context.Canceled) {
			logger.Errorf("[transcode] error transcoding video file: %v", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			logger.Warnf("[transcode] error writing response: %v", err)
		}
		return
	}

	handler(w, r)
}

func (sm *StreamManager) getTranscodeStream(ctx *fsutil.LockContext, options TranscodeOptions) (http.HandlerFunc, error) {
	args := options.makeStreamArgs(sm)
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
		errStr, _ := io.ReadAll(stderr)

		errCmd := cmd.Wait()

		var err error

		e := string(errStr)
		if e != "" {
			err = errors.New(e)
		} else {
			err = errCmd
		}

		// ignore ExitErrors, the process is always forcibly killed
		var exitError *exec.ExitError
		if err != nil && !errors.As(err, &exitError) {
			logger.Errorf("[transcode] ffmpeg error when running command <%s>: %v", strings.Join(cmd.Args, " "), err)
		}
	}()

	mimeType := options.StreamType.MimeType
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", mimeType)
		w.WriteHeader(http.StatusOK)

		// process killing should be handled by command context

		_, err := io.Copy(w, stdout)
		if err != nil && !errors.Is(err, syscall.EPIPE) && !errors.Is(err, syscall.ECONNRESET) {
			logger.Errorf("[transcode] error serving transcoded video file: %v", err)
		}

		w.(http.Flusher).Flush()
	}
	return handler, nil
}
