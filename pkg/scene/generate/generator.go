// Package generate provides functions to generate media assets from scenes.
package generate

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
)

const (
	mp4Pattern  = "*.mp4"
	webpPattern = "*.webp"
	jpgPattern  = "*.jpg"
	txtPattern  = "*.txt"
	vttPattern  = "*.vtt"
)

type Paths interface {
	TempFile(pattern string) (*os.File, error)
}

type MarkerPaths interface {
	Paths

	GetVideoPreviewPath(checksum string, seconds int) string
	GetWebpPreviewPath(checksum string, seconds int) string
	GetScreenshotPath(checksum string, seconds int) string
}

type ScenePaths interface {
	Paths

	GetVideoPreviewPath(checksum string) string
	GetWebpPreviewPath(checksum string) string

	GetSpriteImageFilePath(checksum string) string
	GetSpriteVttFilePath(checksum string) string

	GetTranscodePath(checksum string) string
}

type FFMpegConfig interface {
	GetTranscodeInputArgs() []string
	GetTranscodeOutputArgs() []string
}

type Generator struct {
	Encoder      *ffmpeg.FFMpeg
	FFMpegConfig FFMpegConfig
	LockManager  *fsutil.ReadLockManager
	MarkerPaths  MarkerPaths
	ScenePaths   ScenePaths
	Overwrite    bool
}

type generateFn func(lockCtx *fsutil.LockContext, tmpFn string) error

func (g Generator) tempFile(p Paths, pattern string) (*os.File, error) {
	tmpFile, err := p.TempFile(pattern) // tmp output in case the process ends abruptly
	if err != nil {
		return nil, fmt.Errorf("creating temporary file: %w", err)
	}
	_ = tmpFile.Close()
	return tmpFile, err
}

// generateFile performs a generate operation by generating a temporary file using p and pattern, then
// moving it to output on success.
func (g Generator) generateFile(lockCtx *fsutil.LockContext, p Paths, pattern string, output string, generateFn generateFn) error {
	tmpFile, err := g.tempFile(p, pattern) // tmp output in case the process ends abruptly
	if err != nil {
		return err
	}

	tmpFn := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpFn)
	}()

	if err := generateFn(lockCtx, tmpFn); err != nil {
		return err
	}

	// check if generated empty file
	stat, err := os.Stat(tmpFn)
	if err != nil {
		return fmt.Errorf("error getting file stat: %w", err)
	}

	if stat.Size() == 0 {
		return fmt.Errorf("ffmpeg command produced no output")
	}

	if err := fsutil.SafeMove(tmpFn, output); err != nil {
		return fmt.Errorf("moving %s to %s failed: %w", tmpFn, output, err)
	}

	return nil
}

// generateBytes performs a generate operation by generating a temporary file using p and pattern, returns the contents, then deletes it.
func (g Generator) generateBytes(lockCtx *fsutil.LockContext, p Paths, pattern string, generateFn generateFn) ([]byte, error) {
	tmpFile, err := g.tempFile(p, pattern) // tmp output in case the process ends abruptly
	if err != nil {
		return nil, err
	}

	tmpFn := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpFn)
	}()

	if err := generateFn(lockCtx, tmpFn); err != nil {
		return nil, err
	}

	defer os.Remove(tmpFn)
	return os.ReadFile(tmpFn)
}

// generate runs ffmpeg with the given args and waits for it to finish.
// Returns an error if the command fails. If the command fails, the return
// value will be of type *exec.ExitError.
func (g Generator) generate(ctx *fsutil.LockContext, args []string) error {
	cmd := g.Encoder.Command(ctx, args)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %w", err)
	}

	ctx.AttachCommand(cmd)

	if err := cmd.Wait(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitErr.Stderr = stderr.Bytes()
			err = exitErr
		}
		return fmt.Errorf("error running ffmpeg command <%s>: %w", strings.Join(args, " "), err)
	}

	return nil
}

// GenerateOutput runs ffmpeg with the given args and returns it standard output.
func (g Generator) generateOutput(lockCtx *fsutil.LockContext, args []string) ([]byte, error) {
	cmd := g.Encoder.Command(lockCtx, args)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting command: %w", err)
	}

	lockCtx.AttachCommand(cmd)

	if err := cmd.Wait(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitErr.Stderr = stderr.Bytes()
			err = exitErr
		}
		return nil, fmt.Errorf("error running ffmpeg command <%s>: %w", strings.Join(args, " "), err)
	}

	if stdout.Len() == 0 {
		return nil, fmt.Errorf("ffmpeg command produced no output: <%s>", strings.Join(args, " "))
	}

	return stdout.Bytes(), nil
}
