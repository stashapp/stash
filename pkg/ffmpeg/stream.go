package ffmpeg

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type Stream struct {
	Stdout  io.ReadCloser
	Stderr  io.ReadCloser
	Process *os.Process
}

type Codec struct {
	codec     string
	format    string
	MimeType  string
	extraArgs []string
}

var CodecH264 Codec = Codec{
	codec:    "libx264",
	format:   "mp4",
	MimeType: MimeMp4,
	extraArgs: []string{
		"-movflags", "frag_keyframe",
		"-preset", "veryfast",
		"-crf", "30",
	},
}

var CodecVP8 Codec = Codec{
	codec:    "libvpx-vp9",
	format:   "webm",
	MimeType: MimeWebm,
	extraArgs: []string{
		"-deadline", "realtime",
		"-cpu-used", "5",
		"-row-mt", "1",
		"-crf", "30",
	},
}

var CodecHEVC Codec = Codec{
	codec:    "libx265",
	format:   "mp4",
	MimeType: MimeMp4,
	extraArgs: []string{
		"-movflags", "frag_keyframe",
		"-preset", "veryfast",
		"-crf", "30",
	},
}

type streamOptions struct {
	probeResult      VideoFile
	codec            Codec
	startTime        string
	maxTranscodeSize models.StreamingResolutionEnum
	videoOnly        bool
}

func (s streamOptions) getStreamArgs() []string {
	scale := calculateTranscodeScale(s.probeResult, s.maxTranscodeSize)

	args := []string{
		"-hide_banner",
		"-v", "error",
	}

	if s.startTime != "" {
		args = append(args, "-ss", s.startTime)
	}

	args = append(args,
		"-i", s.probeResult.Path,
	)

	if s.videoOnly {
		args = append(args, "-an")
	}

	args = append(args,
		"-c:v", s.codec.codec,
		"-vf", "scale="+scale,
	)

	if len(s.codec.extraArgs) > 0 {
		args = append(args, s.codec.extraArgs...)
	}

	args = append(args,
		"-b:v", "0",
		"-f", s.codec.format,
		"pipe:",
	)

	return args
}

func (e *Encoder) StreamTranscode(probeResult VideoFile, codec Codec, startTime string, maxTranscodeSize models.StreamingResolutionEnum) (*Stream, error) {
	options := streamOptions{
		probeResult:      probeResult,
		codec:            codec,
		startTime:        startTime,
		maxTranscodeSize: maxTranscodeSize,
	}

	return e.stream(probeResult, options.getStreamArgs())
}

//transcode the video, remove the audio
//in some videos where the audio codec is not supported by ffmpeg
//ffmpeg fails if you try to transcode the audio
func (e *Encoder) StreamTranscodeVideo(probeResult VideoFile, codec Codec, startTime string, maxTranscodeSize models.StreamingResolutionEnum) (*Stream, error) {
	options := streamOptions{
		probeResult:      probeResult,
		codec:            codec,
		startTime:        startTime,
		maxTranscodeSize: maxTranscodeSize,
		videoOnly:        true,
	}

	return e.stream(probeResult, options.getStreamArgs())
}

//it is very common in MKVs to have just the audio codec unsupported
//copy the video stream, transcode the audio and serve as Matroska
func (e *Encoder) StreamMkvTranscodeAudio(probeResult VideoFile, startTime string, maxTranscodeSize models.StreamingResolutionEnum) (*Stream, error) {
	args := []string{
		"-hide_banner",
		"-v", "error",
	}

	if startTime != "" {
		args = append(args, "-ss", startTime)
	}

	args = append(args,
		"-i", probeResult.Path,
		"-c:v", "copy",
		"-c:a", "libopus",
		"-b:a", "96k",
		"-vbr", "on",
		"-f", "matroska",
		"pipe:",
	)

	return e.stream(probeResult, args)
}

func (e *Encoder) stream(probeResult VideoFile, args []string) (*Stream, error) {
	cmd := exec.Command(e.Path, args...)
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
	go waitAndDeregister(probeResult.Path, cmd)

	ret := &Stream{
		Stdout:  stdout,
		Stderr:  stderr,
		Process: cmd.Process,
	}
	return ret, nil
}
