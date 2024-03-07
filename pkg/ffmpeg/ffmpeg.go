// Package ffmpeg provides a wrapper around the ffmpeg and ffprobe executables.
package ffmpeg

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	stashExec "github.com/stashapp/stash/pkg/exec"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

func ValidateFFMpeg(ffmpegPath string) error {
	cmd := stashExec.Command(ffmpegPath, "-h")
	bytes, err := cmd.CombinedOutput()
	output := string(bytes)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("error running ffmpeg: %v", output)
		}

		return fmt.Errorf("error running ffmpeg: %v", err)
	}

	if !strings.Contains(output, "--enable-libopus") {
		return fmt.Errorf("ffmpeg is missing libopus support")
	}
	if !strings.Contains(output, "--enable-libvpx") {
		return fmt.Errorf("ffmpeg is missing libvpx support")
	}
	if !strings.Contains(output, "--enable-libx264") {
		return fmt.Errorf("ffmpeg is missing libx264 support")
	}
	if !strings.Contains(output, "--enable-libx265") {
		return fmt.Errorf("ffmpeg is missing libx265 support")
	}
	if !strings.Contains(output, "--enable-libwebp") {
		return fmt.Errorf("ffmpeg is missing libwebp support")
	}
	return nil
}

func FindFFMpeg(paths []string) string {
	ret, _ := exec.LookPath(getFFMpegFilename())

	if ret != "" {
		// ensure ffmpeg has the correct flags
		if err := ValidateFFMpeg(ret); err != nil {
			logger.Warnf("ffmpeg found in PATH (%s), but it is missing required flags: %v", ret, err)
			ret = ""
		}
	}

	if ret == "" {
		ret = fsutil.FindInPaths(paths, getFFMpegFilename())

		if ret != "" {
			// ensure ffmpeg has the correct flags
			if err := ValidateFFMpeg(ret); err != nil {
				logger.Warnf("ffmpeg found (%s), but it is missing required flags: %v", ret, err)
				ret = ""
			}
		}
	}

	return ret
}

// FFMpeg provides an interface to ffmpeg.
type FFMpeg struct {
	ffmpeg         string
	hwCodecSupport []VideoCodec
}

// Creates a new FFMpeg encoder
func NewEncoder(ffmpegPath string) *FFMpeg {
	ret := &FFMpeg{
		ffmpeg: ffmpegPath,
	}

	return ret
}

// Returns an exec.Cmd that can be used to run ffmpeg using args.
func (f *FFMpeg) Command(ctx context.Context, args []string) *exec.Cmd {
	return stashExec.CommandContext(ctx, string(f.ffmpeg), args...)
}

func (f *FFMpeg) Path() string {
	return f.ffmpeg
}
