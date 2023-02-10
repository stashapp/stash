package ffmpeg

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"syscall"

	"github.com/stashapp/stash/pkg/logger"
)

const (
	MimeWebm   string = "video/webm"
	MimeMkv    string = "video/x-matroska"
	MimeMp4    string = "video/mp4"
	MimeHLS    string = "application/vnd.apple.mpegurl"
	MimeMpegts string = "video/MP2T"
)

// Stream represents an ongoing transcoded stream.
type Stream struct {
	Stdout   io.ReadCloser
	Cmd      *exec.Cmd
	mimeType string
}

// Serve is an http handler function that serves the stream.
func (s *Stream) Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", s.mimeType)
	w.WriteHeader(http.StatusOK)

	logger.Infof("[stream] transcoding video file to %s", s.mimeType)

	// process killing should be handled by command context

	_, err := io.Copy(w, s.Stdout)
	if err != nil && !errors.Is(err, syscall.EPIPE) {
		logger.Errorf("[stream] error serving transcoded video file: %v", err)
	}
}

// StreamFormat represents a transcode stream format.
type StreamFormat struct {
	MimeType  string
	codec     VideoCodec
	format    Format
	extraArgs []string
	hls       bool
}

var (
	StreamFormatHLS = StreamFormat{
		codec:    VideoCodecLibX264,
		format:   FormatMpegTS,
		MimeType: MimeMpegts,
		extraArgs: []string{
			"-acodec", "aac",
			"-pix_fmt", "yuv420p",
			"-preset", "veryfast",
			"-crf", "25",
		},
		hls: true,
	}

	// libx264
	StreamFormatH264 = StreamFormat{
		codec:    VideoCodecLibX264,
		format:   FormatMP4,
		MimeType: MimeMp4,
		extraArgs: []string{
			"-movflags", "frag_keyframe+empty_moov",
			"-pix_fmt", "yuv420p",
			"-preset", "veryfast",
			"-crf", "25",
		},
	}

	// NVIDIA NVENC H264
	StreamFormatN264 = StreamFormat{
		codec:    VideoCodecLibN264,
		format:   FormatMP4,
		MimeType: MimeMp4,
		extraArgs: []string{
			"-movflags", "frag_keyframe+empty_moov",
			"-pix_fmt", "yuv420p",
			"-preset", "p1",
			"-rc", "vbr",
			"-cq", "30", //TODO: Test this
		},
	}

	// Intel QSV H264
	StreamFormatI264 = StreamFormat{
		codec:    VideoCodecLibI264,
		format:   FormatMP4,
		MimeType: MimeMp4,
		extraArgs: []string{
			"-movflags", "frag_keyframe+empty_moov",
			"-pix_fmt", "yuv420p",
			"-preset", "veryfast",
			"-look_ahead", "1",
			"-global_quality", "30",
		},
	}

	// AMD AMF H264
	StreamFormatA264 = StreamFormat{
		codec:    VideoCodecLibA264,
		format:   FormatMP4,
		MimeType: MimeMp4,
		extraArgs: []string{
			"-movflags", "frag_keyframe+empty_moov",
			"-pix_fmt", "yuv420p",
			"-quality", "speed",
		},
	}

	// VAAPI H264
	StreamFormatV264 = StreamFormat{
		codec:    VideoCodecLibV264,
		format:   FormatMP4,
		MimeType: MimeMp4,
		extraArgs: []string{
			"-movflags", "frag_keyframe+empty_moov",
			"-pix_fmt", "yuv420p",
			"-quality", "50",
		},
	}

	// Raspberry Pi 4, H264
	StreamFormatR264 = StreamFormat{
		codec:    VideoCodecLibR264,
		format:   FormatMP4,
		MimeType: MimeMp4,
		extraArgs: []string{
			"-movflags", "frag_keyframe+empty_moov",
			"-pix_fmt", "yuv420p",
			"-crf", "25",
		},
	}

	// OpenMAX IL, H.264
	StreamFormatO264 = StreamFormat{
		codec:    VideoCodecLibO264,
		format:   FormatMP4,
		MimeType: MimeMp4,
		extraArgs: []string{
			"-movflags", "frag_keyframe+empty_moov",
			"-pix_fmt", "yuv420p",
			"-preset", "superfast",
			"-crf", "25",
		},
	}

	// VP9
	StreamFormatVP9 = StreamFormat{
		codec:    VideoCodecVP9,
		format:   FormatWebm,
		MimeType: MimeWebm,
		extraArgs: []string{
			"-deadline", "realtime",
			"-cpu-used", "5",
			"-row-mt", "1",
			"-crf", "30",
			"-b:v", "0",
			"-pix_fmt", "yuv420p",
		},
	}

	// Intel QSV VP9
	StreamFormatIVP9 = StreamFormat{
		codec:    VideoCodecIVP9,
		format:   FormatWebm,
		MimeType: MimeWebm,
		extraArgs: []string{
			"-pix_fmt", "yuv420p",
			"-preset", "veryfast",
			"-look_ahead", "1",
			"-global_quality", "30",
		},
	}

	// VAAPI VP9
	StreamFormatVVP9 = StreamFormat{
		codec:    VideoCodecVVP9,
		format:   FormatWebm,
		MimeType: MimeWebm,
		extraArgs: []string{
			"-pix_fmt", "yuv420p",
			"-global_quality", "30",
		},
	}

	StreamFormatVP8 = StreamFormat{
		codec:    VideoCodecVPX,
		format:   FormatWebm,
		MimeType: MimeWebm,
		extraArgs: []string{
			"-deadline", "realtime",
			"-cpu-used", "5",
			"-crf", "12",
			"-b:v", "3M",
			"-pix_fmt", "yuv420p",
		},
	}

	StreamFormatHEVC = StreamFormat{
		codec:    VideoCodecLibX265,
		format:   FormatMP4,
		MimeType: MimeMp4,
		extraArgs: []string{
			"-movflags", "frag_keyframe",
			"-preset", "veryfast",
			"-crf", "30",
		},
	}

	// it is very common in MKVs to have just the audio codec unsupported
	// copy the video stream, transcode the audio and serve as Matroska
	StreamFormatMKVAudio = StreamFormat{
		codec:    VideoCodecCopy,
		format:   FormatMatroska,
		MimeType: MimeMkv,
		extraArgs: []string{
			"-c:a", "libopus",
			"-b:a", "96k",
			"-vbr", "on",
		},
	}
)

// TranscodeStreamOptions represents options for live transcoding a video file.
type TranscodeStreamOptions struct {
	Input            string
	Codec            StreamFormat
	StartTime        float64
	MaxTranscodeSize int

	// original video dimensions
	VideoWidth  int
	VideoHeight int

	// transcode the video, remove the audio
	// in some videos where the audio codec is not supported by ffmpeg
	// ffmpeg fails if you try to transcode the audio
	VideoOnly bool

	// arguments added before the input argument
	ExtraInputArgs []string
	// arguments added before the output argument
	ExtraOutputArgs []string
}

func (o TranscodeStreamOptions) getStreamArgs() Args {
	var args Args
	args = append(args, "-hide_banner")
	args = append(args, o.ExtraInputArgs...)
	args = args.LogLevel(LogLevelError)

	if o.StartTime != 0 {
		args = args.Seek(o.StartTime)
	}

	if o.Codec.hls {
		// we only serve a fixed segment length
		args = args.Duration(hlsSegmentLength)
	}

	args = args.Input(o.Input)

	if o.VideoOnly {
		args = args.SkipAudio()
	}

	args = args.VideoCodec(o.Codec.codec)

	// don't set scale when copying video stream
	if o.Codec.codec != VideoCodecCopy {
		var videoFilter VideoFilter
		videoFilter = videoFilter.ScaleMax(o.VideoWidth, o.VideoHeight, o.MaxTranscodeSize)
		args = args.VideoFilter(videoFilter)
	}

	if len(o.Codec.extraArgs) > 0 {
		args = append(args, o.Codec.extraArgs...)
	}

	args = append(args,
		// this is needed for 5-channel ac3 files
		"-ac", "2",
	)

	args = append(args, o.ExtraOutputArgs...)

	args = args.Format(o.Codec.format)
	args = args.Output("pipe:")

	return args
}

// GetTranscodeStream starts the live transcoding process using ffmpeg and returns a stream.
func (f *FFMpeg) GetTranscodeStream(ctx context.Context, options TranscodeStreamOptions) (*Stream, error) {
	args := options.getStreamArgs()
	cmd := f.Command(ctx, args)
	logger.Debugf("Streaming via: %s", strings.Join(cmd.Args, " "))

	stdout, err := cmd.StdoutPipe()
	if nil != err {
		logger.Error("FFMPEG stdout not available: " + err.Error())
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if nil != err {
		logger.Error("FFMPEG stderr not available: " + err.Error())
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	// stderr must be consumed or the process deadlocks
	go func() {
		stderrData, _ := io.ReadAll(stderr)
		stderrString := string(stderrData)
		if len(stderrString) > 0 {
			logger.Debugf("[stream] ffmpeg stderr: %s", stderrString)
		}
	}()

	ret := &Stream{
		Stdout:   stdout,
		Cmd:      cmd,
		mimeType: options.Codec.MimeType,
	}
	return ret, nil
}
