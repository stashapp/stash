package ffmpeg

import (
	"fmt"
	"github.com/stashapp/stash/pkg/utils"
	"strconv"
)

type ScenePreviewChunkOptions struct {
	Time       int
	Width      int
	OutputPath string
}

func (e *Encoder) ScenePreviewVideoChunk(probeResult VideoFile, options ScenePreviewChunkOptions) {
	args := []string{
		"-v", "quiet",
		"-ss", strconv.Itoa(options.Time),
		"-t", "0.75",
		"-i", probeResult.Path,
		"-y",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-profile:v", "high",
		"-level", "4.2",
		"-preset", "veryslow",
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
		"-v", "quiet",
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
		"-v", "quiet",
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
