package ffmpeg

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	stashExec "github.com/stashapp/stash/pkg/exec"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

const CopyStreamCodec = "copy"

type Stream struct {
	Stdout   io.ReadCloser
	Process  *os.Process
	options  TranscodeStreamOptions
	mimeType string
}

func (s *Stream) Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", s.mimeType)
	w.WriteHeader(http.StatusOK)

	logger.Infof("[stream] transcoding video file to %s", s.mimeType)

	// handle if client closes the connection
	notify := r.Context().Done()
	go func() {
		<-notify
		if err := s.Process.Kill(); err != nil {
			logger.Warnf("unable to kill os process %v: %v", s.Process.Pid, err)
		}
	}()

	_, err := io.Copy(w, s.Stdout)
	if err != nil {
		logger.Errorf("[stream] error serving transcoded video file: %s", err.Error())
	}
}

type Codec struct {
	Codec     string
	format    string
	MimeType  string
	extraArgs []string
	hls       bool
}

var CodecHLS = Codec{
	Codec:    "libx264",
	format:   "mpegts",
	MimeType: MimeMpegts,
	extraArgs: []string{
		"-acodec", "aac",
		"-pix_fmt", "yuv420p",
		"-preset", "veryfast",
		"-crf", "25",
	},
	hls: true,
}

var CodecH264 = Codec{
	Codec:    "libx264",
	format:   "mp4",
	MimeType: MimeMp4,
	extraArgs: []string{
		"-movflags", "frag_keyframe+empty_moov",
		"-pix_fmt", "yuv420p",
		"-preset", "veryfast",
		"-crf", "25",
	},
}

var CodecVP9 = Codec{
	Codec:    "libvpx-vp9",
	format:   "webm",
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

var CodecVP8 = Codec{
	Codec:    "libvpx",
	format:   "webm",
	MimeType: MimeWebm,
	extraArgs: []string{
		"-deadline", "realtime",
		"-cpu-used", "5",
		"-crf", "12",
		"-b:v", "3M",
		"-pix_fmt", "yuv420p",
	},
}

var CodecHEVC = Codec{
	Codec:    "libx265",
	format:   "mp4",
	MimeType: MimeMp4,
	extraArgs: []string{
		"-movflags", "frag_keyframe",
		"-preset", "veryfast",
		"-crf", "30",
	},
}

// it is very common in MKVs to have just the audio codec unsupported
// copy the video stream, transcode the audio and serve as Matroska
var CodecMKVAudio = Codec{
	Codec:    CopyStreamCodec,
	format:   "matroska",
	MimeType: MimeMkv,
	extraArgs: []string{
		"-c:a", "libopus",
		"-b:a", "96k",
		"-vbr", "on",
	},
}

type TranscodeStreamOptions struct {
	ProbeResult      VideoFile
	Codec            Codec
	StartTime        string
	MaxTranscodeSize models.StreamingResolutionEnum
	// transcode the video, remove the audio
	// in some videos where the audio codec is not supported by ffmpeg
	// ffmpeg fails if you try to transcode the audio
	VideoOnly bool
}

func GetTranscodeStreamOptions(probeResult VideoFile, videoCodec Codec, audioCodec AudioCodec) TranscodeStreamOptions {
	options := TranscodeStreamOptions{
		ProbeResult: probeResult,
		Codec:       videoCodec,
	}

	if audioCodec == MissingUnsupported {
		// ffmpeg fails if it trys to transcode a non supported audio codec
		options.VideoOnly = true
	}

	return options
}

func (o TranscodeStreamOptions) getStreamArgs() []string {
	args := []string{
		"-hide_banner",
		"-v", "error",
	}

	if o.StartTime != "" {
		args = append(args, "-ss", o.StartTime)
	}

	if o.Codec.hls {
		// we only serve a fixed segment length
		args = append(args, "-t", strconv.Itoa(int(hlsSegmentLength)))
	}

	args = append(args,
		"-i", o.ProbeResult.Path,
	)

	if o.VideoOnly {
		args = append(args, "-an")
	}

	args = append(args,
		"-c:v", o.Codec.Codec,
	)

	// don't set scale when copying video stream
	if o.Codec.Codec != CopyStreamCodec {
		scale := calculateTranscodeScale(o.ProbeResult, o.MaxTranscodeSize)
		args = append(args,
			"-vf", "scale="+scale,
		)
	}

	if len(o.Codec.extraArgs) > 0 {
		args = append(args, o.Codec.extraArgs...)
	}

	args = append(args,
		// this is needed for 5-channel ac3 files
		"-ac", "2",
		"-f", o.Codec.format,
		"pipe:",
	)

	return args
}

func (e *Encoder) GetTranscodeStream(options TranscodeStreamOptions) (*Stream, error) {
	return e.stream(options.ProbeResult, options)
}

func (e *Encoder) stream(probeResult VideoFile, options TranscodeStreamOptions) (*Stream, error) {
	args := options.getStreamArgs()
	cmd := stashExec.Command(string(*e), args...)
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

	registerRunningEncoder(probeResult.Path, cmd.Process)
	go func() {
		if err := waitAndDeregister(probeResult.Path, cmd); err != nil {
			logger.Warnf("Error while deregistering ffmpeg stream: %v", err)
		}
	}()

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
		Process:  cmd.Process,
		options:  options,
		mimeType: options.Codec.MimeType,
	}
	return ret, nil
}
