package ffmpeg

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type ScenePreviewChunkOptions struct {
	StartTime  float64
	Duration   float64
	Width      int
	OutputPath string
	Audio      bool
}

func (e *Encoder) ScenePreviewVideoChunk(probeResult VideoFile, options ScenePreviewChunkOptions, preset string, fallback bool) error {
	var fastSeek float64
	var slowSeek float64
	fallbackMinSlowSeek := 20.0

	args := []string{
		"-v", "error",
	}

	argsAudio := []string{
		"-c:a", "aac",
		"-b:a", "128k",
	}

	if !options.Audio {
		argsAudio = []string{
			"-an",
		}
	}

	// Non-fallback: enable xerror.
	// "-xerror" causes ffmpeg to fail on warnings, often the preview is fine but could be broken.
	if !fallback {
		args = append(args, "-xerror")
		fastSeek = options.StartTime
		slowSeek = 0
	} else {
		// In fallback mode, disable "-xerror" and try a combination of fast/slow seek instead of just fastseek
		// Commonly with avi/wmv ffmpeg doesn't seem to always predict the right start point to begin decoding when
		// using fast seek. If you force ffmpeg to decode more, it avoids the "blocky green artifact" issue.
		if options.StartTime > fallbackMinSlowSeek {
			// Handle seeks longer than fallbackMinSlowSeek with fast/slow seeks
			// Allow for at least fallbackMinSlowSeek seconds of slow seek
			fastSeek = options.StartTime - fallbackMinSlowSeek
			slowSeek = fallbackMinSlowSeek
		} else {
			// Handle seeks shorter than fallbackMinSlowSeek with only slow seeks.
			slowSeek = options.StartTime
			fastSeek = 0
		}
	}

	if fastSeek > 0 {
		args = append(args, "-ss")
		args = append(args, strconv.FormatFloat(fastSeek, 'f', 2, 64))
	}

	args = append(args, "-i")
	args = append(args, probeResult.Path)

	if slowSeek > 0 {
		args = append(args, "-ss")
		args = append(args, strconv.FormatFloat(slowSeek, 'f', 2, 64))
	}

	args2 := []string{
		"-t", strconv.FormatFloat(options.Duration, 'f', 2, 64),
		"-max_muxing_queue_size", "1024", // https://trac.ffmpeg.org/ticket/6375
		"-y",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-profile:v", "high",
		"-level", "4.2",
		"-preset", preset,
		"-crf", "21",
		"-threads", "4",
		"-vf", fmt.Sprintf("scale=%v:-2", options.Width),
		"-strict", "-2",
	}

	args = append(args, args2...)
	args = append(args, argsAudio...)
	args = append(args, options.OutputPath)

	_, err := e.run(probeResult.Path, args, nil)
	return err
}

// fixWindowsPath replaces \ with / in the given path because the \ isn't recognized as valid on windows ffmpeg
func fixWindowsPath(str string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(str, `\`, "/")
	}
	return str
}

func (e *Encoder) ScenePreviewVideoChunkCombine(probeResult VideoFile, concatFilePath string, outputPath string) error {
	args := []string{
		"-v", "error",
		"-f", "concat",
		"-i", fixWindowsPath(concatFilePath),
		"-y",
		"-c", "copy",
		outputPath,
	}
	_, err := e.run(probeResult.Path, args, nil)
	return err
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
	_, err := e.run(probeResult.Path, args, nil)
	return err
}
