package ffmpeg

import (
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/utils"
)

type ScenePreviewChunkOptions struct {
	Time       int
	Width      int
	OutputPath string
}

func (e *Encoder) ScenePreviewVideoChunk(probeResult VideoFile, options ScenePreviewChunkOptions) {
	args := []string{
		"-v", "error",
		"-ss", strconv.Itoa(options.Time),
		"-i", probeResult.Path,
		"-t", "0.75",
		"-max_muxing_queue_size", "1024", // https://trac.ffmpeg.org/ticket/6375
		"-y",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-profile:v", "high",
		"-level", "4.2",
		"-preset", "slow",
		"-crf", "21",
		"-threads", "4",
		"-vf", fmt.Sprintf("scale=%v:-2", options.Width),
		"-c:a", "aac",
		"-b:a", "128k",
		"-strict", "-2",
		options.OutputPath,
	}
	_, _ = e.run(probeResult, args)
}

func (e *Encoder) ScenePreviewVideoChunkCombine(probeResult VideoFile, concatFilePath string, outputPath string) {
	args := []string{
		"-v", "error",
		"-f", "concat",
		"-i", utils.FixWindowsPath(concatFilePath),
		"-y",
		"-c", "copy",
		outputPath,
	}
	_, _ = e.run(probeResult, args)
}

func (e *Encoder) ScenePreviewVideoToImage(probeResult VideoFile, width int, videoPreviewPath string, outputPath string) error {
	args := []string{
		"-v", "error",
		"-i", videoPreviewPath,
		"-y",
		"-c:v", "libwebp",
		"-lossless", "1",
		"-q:v", "70",
		"-compression_level", "6",
		"-preset", "default",
		"-loop", "0",
		"-threads", "4",
		"-vf", fmt.Sprintf("scale=%v:-2,fps=12", width),
		"-an",
		outputPath,
	}
	_, err := e.run(probeResult, args)
	return err
}
