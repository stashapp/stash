package ffmpeg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// Generate runs ffmpeg with the given args and waits for it to finish.
// Returns an error if the command fails. If the command fails, the return
// value will be of type *exec.ExitError.
func (f *FFMpeg) Generate(ctx context.Context, args Args) error {
	cmd := f.Command(ctx, args)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %w", err)
	}

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
func (f *FFMpeg) GenerateOutput(ctx context.Context, args []string, stdin io.Reader) ([]byte, error) {
	cmd := f.Command(ctx, args)
	cmd.Stdin = stdin

	ret, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running ffmpeg command <%s>: %w", strings.Join(args, " "), err)
	}

	return ret, nil
}
