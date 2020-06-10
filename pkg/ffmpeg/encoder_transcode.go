package ffmpeg

import (
	"io"
	"os"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type TranscodeOptions struct {
	OutputPath       string
	MaxTranscodeSize models.StreamingResolutionEnum
}

func calculateTranscodeScale(probeResult VideoFile, maxTranscodeSize models.StreamingResolutionEnum) string {
	maxSize := 0
	switch maxTranscodeSize {
	case models.StreamingResolutionEnumLow:
		maxSize = 240
	case models.StreamingResolutionEnumStandard:
		maxSize = 480
	case models.StreamingResolutionEnumStandardHd:
		maxSize = 720
	case models.StreamingResolutionEnumFullHd:
		maxSize = 1080
	case models.StreamingResolutionEnumFourK:
		maxSize = 2160
	}

	// get the smaller dimension of the video file
	videoSize := probeResult.Height
	if probeResult.Width < videoSize {
		videoSize = probeResult.Width
	}

	// if our streaming resolution is larger than the video dimension
	// or we are streaming the original resolution, then just set the
	// input width
	if maxSize >= videoSize || maxSize == 0 {
		return "iw:-2"
	}

	// we're setting either the width or height
	// we'll set the smaller dimesion
	if probeResult.Width > probeResult.Height {
		// set the height
		return "-2:" + strconv.Itoa(maxSize)
	}

	return strconv.Itoa(maxSize) + ":-2"
}

func (e *Encoder) Transcode(probeResult VideoFile, options TranscodeOptions) {
	scale := calculateTranscodeScale(probeResult, options.MaxTranscodeSize)
	args := []string{
		"-i", probeResult.Path,
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-profile:v", "high",
		"-level", "4.2",
		"-preset", "superfast",
		"-crf", "23",
		"-vf", "scale=" + scale,
		"-c:a", "aac",
		"-strict", "-2",
		options.OutputPath,
	}
	_, _ = e.run(probeResult, args)
}

//transcode the video, remove the audio
//in some videos where the audio codec is not supported by ffmpeg
//ffmpeg fails if you try to transcode the audio
func (e *Encoder) TranscodeVideo(probeResult VideoFile, options TranscodeOptions) {
	scale := calculateTranscodeScale(probeResult, options.MaxTranscodeSize)
	args := []string{
		"-i", probeResult.Path,
		"-an",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-profile:v", "high",
		"-level", "4.2",
		"-preset", "superfast",
		"-crf", "23",
		"-vf", "scale=" + scale,
		options.OutputPath,
	}
	_, _ = e.run(probeResult, args)
}

//copy the video stream as is, transcode audio
func (e *Encoder) TranscodeAudio(probeResult VideoFile, options TranscodeOptions) {
	args := []string{
		"-i", probeResult.Path,
		"-c:v", "copy",
		"-c:a", "aac",
		"-strict", "-2",
		options.OutputPath,
	}
	_, _ = e.run(probeResult, args)
}

//copy the video stream as is, drop audio
func (e *Encoder) CopyVideo(probeResult VideoFile, options TranscodeOptions) {
	args := []string{
		"-i", probeResult.Path,
		"-an",
		"-c:v", "copy",
		options.OutputPath,
	}
	_, _ = e.run(probeResult, args)
}

func (e *Encoder) StreamTranscode(probeResult VideoFile, codec string, startTime string, maxTranscodeSize models.StreamingResolutionEnum) (io.ReadCloser, *os.Process, error) {
	scale := calculateTranscodeScale(probeResult, maxTranscodeSize)
	args := []string{}

	if startTime != "" {
		args = append(args, "-ss", startTime)
	}

	codecArg := "libvpx-vp9"
	format := "webm"
	if codec == H264 {
		codecArg = "libx264"
		format = "mp4"
	}

	args = append(args,
		"-i", probeResult.Path,
		"-c:v", codecArg,
		"-vf", "scale="+scale,
		"-deadline", "realtime",
		"-cpu-used", "5",
		"-row-mt", "1",
		"-crf", "30",
		"-b:v", "0",
		"-f", format,
		"pipe:",
	)

	return e.stream(probeResult, args)
}

//transcode the video, remove the audio
//in some videos where the audio codec is not supported by ffmpeg
//ffmpeg fails if you try to transcode the audio
func (e *Encoder) StreamTranscodeVideo(probeResult VideoFile, codec string, startTime string, maxTranscodeSize models.StreamingResolutionEnum) (io.ReadCloser, *os.Process, error) {
	scale := calculateTranscodeScale(probeResult, maxTranscodeSize)
	args := []string{}

	if startTime != "" {
		args = append(args, "-ss", startTime)
	}

	codecArg := "libvpx-vp9"
	format := "webm"
	if codec == H264 {
		codecArg = "libx264"
		format = "mp4"
	}

	args = append(args,
		"-i", probeResult.Path,
		"-an",
		"-c:v", codecArg,
		"-vf", "scale="+scale,
		"-deadline", "realtime",
		"-cpu-used", "5",
		"-row-mt", "1",
		"-crf", "30",
		"-b:v", "0",
		"-f", format,
		"pipe:",
	)

	return e.stream(probeResult, args)
}

//it is very common in MKVs to have just the audio codec unsupported
//copy the video stream, transcode the audio and serve as Matroska
func (e *Encoder) StreamMkvTranscodeAudio(probeResult VideoFile, startTime string, maxTranscodeSize models.StreamingResolutionEnum) (io.ReadCloser, *os.Process, error) {
	args := []string{}

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
