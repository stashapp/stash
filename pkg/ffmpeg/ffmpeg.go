// Package ffmpeg provides a wrapper around the ffmpeg and ffprobe executables.
package ffmpeg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
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

func LookPathFFMpeg() string {
	ret, _ := exec.LookPath(getFFMpegFilename())

	if ret != "" {
		// ensure ffmpeg has the correct flags
		if err := ValidateFFMpeg(ret); err != nil {
			logger.Warnf("ffmpeg found in PATH (%s), but it is missing required flags: %v", ret, err)
			ret = ""
		}
	}

	return ret
}

func FindFFMpeg(path string) string {
	ret := fsutil.FindInPaths([]string{path}, getFFMpegFilename())

	if ret != "" {
		// ensure ffmpeg has the correct flags
		if err := ValidateFFMpeg(ret); err != nil {
			logger.Warnf("ffmpeg found (%s), but it is missing required flags: %v", ret, err)
			ret = ""
		}
	}

	return ret
}

// ResolveFFMpeg attempts to resolve the path to the ffmpeg executable.
// It first looks in the provided path, then resolves from the environment, and finally looks in the fallback path.
// Returns an empty string if a valid ffmpeg cannot be found.
func ResolveFFMpeg(path string, fallbackPath string) string {
	// look in the provided path first
	ret := FindFFMpeg(path)
	if ret != "" {
		return ret
	}

	// then resolve from the environment
	ret = LookPathFFMpeg()
	if ret != "" {
		return ret
	}

	// finally, look in the fallback path
	ret = FindFFMpeg(fallbackPath)
	return ret
}

func (f *FFMpeg) getVersion() error {
	var args Args
	args = append(args, "-version")
	cmd := f.Command(context.Background(), args)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var err error
	if err = cmd.Run(); err != nil {
		return err
	}

	version_re := regexp.MustCompile(`ffmpeg version ((\d+)\.(\d+)\.(\d+))`)
	stdoutStr := stdout.String()
	match := version_re.FindStringSubmatchIndex(stdoutStr)
	if match == nil {
		return errors.New("version string malformed")
	}

	majorS := stdoutStr[match[4]:match[5]]
	minorS := stdoutStr[match[6]:match[7]]
	patchS := stdoutStr[match[8]:match[9]]
	if i, err := strconv.Atoi(majorS); err == nil {
		f.version.major = i
	}
	if i, err := strconv.Atoi(minorS); err == nil {
		f.version.minor = i
	}
	if i, err := strconv.Atoi(patchS); err == nil {
		f.version.patch = i
	}
	logger.Debugf("FFMpeg version %d.%d.%d detected", f.version.major, f.version.minor, f.version.patch)

	return nil
}

// FFMpeg version params
type FFMpegVersion struct {
	major int
	minor int
	patch int
}

// FFMpeg provides an interface to ffmpeg.
type FFMpeg struct {
	ffmpeg         string
	version        FFMpegVersion
	hwCodecSupport []VideoCodec
}

// Creates a new FFMpeg encoder
func NewEncoder(ffmpegPath string) *FFMpeg {
	ret := &FFMpeg{
		ffmpeg: ffmpegPath,
	}
	if err := ret.getVersion(); err != nil {
		logger.Warnf("FFMpeg version not detected %v", err)
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
