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

func ffmpegHelp(ffmpegPath string) (string, error) {
	cmd := stashExec.Command(ffmpegPath, "-h")
	bytes, err := cmd.CombinedOutput()
	output := string(bytes)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", fmt.Errorf("error running ffmpeg: %v", output)
		}

		return "", fmt.Errorf("error running ffmpeg: %v", err)
	}

	return output, nil
}

func ValidateFFMpeg(ffmpegPath string) error {
	_, err := ffmpegHelp(ffmpegPath)
	return err
}

func ValidateFFMpegCodecSupport(ffmpegPath string) error {
	output, err := ffmpegHelp(ffmpegPath)
	if err != nil {
		return err
	}

	var missingSupport []string

	if !strings.Contains(output, "--enable-libopus") {
		missingSupport = append(missingSupport, "libopus")
	}
	if !strings.Contains(output, "--enable-libvpx") {
		missingSupport = append(missingSupport, "libvpx")
	}
	if !strings.Contains(output, "--enable-libx264") {
		missingSupport = append(missingSupport, "libx264")
	}
	if !strings.Contains(output, "--enable-libx265") {
		missingSupport = append(missingSupport, "libx265")
	}
	if !strings.Contains(output, "--enable-libwebp") {
		missingSupport = append(missingSupport, "libwebp")
	}

	if len(missingSupport) > 0 {
		return fmt.Errorf("ffmpeg missing codec support: %v", missingSupport)
	}

	return nil
}

func LookPathFFMpeg() string {
	ret, _ := exec.LookPath(getFFMpegFilename())

	if ret != "" {
		// ensure ffmpeg has the correct flags
		if err := ValidateFFMpeg(ret); err != nil {
			logger.Warnf("ffmpeg found (%s), could not be executed: %v", ret, err)
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
			logger.Warnf("ffmpeg found (%s), could not be executed: %v", ret, err)
			ret = ""
		}
	}

	return ret
}

// ResolveFFMpeg attempts to resolve the path to the ffmpeg executable.
// It first looks in the provided path, then resolves from the environment, and finally looks in the fallback path.
// It will prefer an ffmpeg binary that has the required codec support.
// Returns an empty string if a valid ffmpeg cannot be found.
func ResolveFFMpeg(path string, fallbackPath string) string {
	var ret string
	// look in the provided path first
	pathFound := FindFFMpeg(path)
	if pathFound != "" {
		err := ValidateFFMpegCodecSupport(pathFound)
		if err == nil {
			return pathFound
		}

		logger.Warnf("ffmpeg found (%s), but it is missing required flags: %v", pathFound, err)
		ret = pathFound
	}

	// then resolve from the environment
	envFound := LookPathFFMpeg()
	if envFound != "" {
		err := ValidateFFMpegCodecSupport(envFound)
		if err == nil {
			return envFound
		}

		logger.Warnf("ffmpeg found (%s), but it is missing required flags: %v", envFound, err)
		if ret == "" {
			ret = envFound
		}
	}

	// finally, look in the fallback path
	fallbackFound := FindFFMpeg(fallbackPath)
	if fallbackFound != "" {
		err := ValidateFFMpegCodecSupport(fallbackFound)
		if err == nil {
			return fallbackFound
		}

		logger.Warnf("ffmpeg found (%s), but it is missing required flags: %v", fallbackFound, err)
		if ret == "" {
			ret = fallbackFound
		}
	}

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

	version_re := regexp.MustCompile(`ffmpeg version n?((\d+)\.(\d+)(?:\.(\d+))?)`)
	stdoutStr := stdout.String()
	match := version_re.FindStringSubmatchIndex(stdoutStr)
	if match == nil {
		return errors.New("version string malformed")
	}

	majorS := stdoutStr[match[4]:match[5]]
	minorS := stdoutStr[match[6]:match[7]]

	// patch is optional
	var patchS string
	if match[8] != -1 && match[9] != -1 {
		patchS = stdoutStr[match[8]:match[9]]
	}

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

// Gteq returns true if the version is greater than or equal to the other version.
func (v FFMpegVersion) Gteq(other FFMpegVersion) bool {
	if v.major > other.major {
		return true
	}
	if v.major == other.major && v.minor > other.minor {
		return true
	}
	if v.major == other.major && v.minor == other.minor && v.patch >= other.patch {
		return true
	}
	return false
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
