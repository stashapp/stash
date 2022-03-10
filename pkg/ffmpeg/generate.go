package ffmpeg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func Generate(ctx context.Context, f FFMpeg, args Args) error {
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

func GenerateOutput(ctx context.Context, f FFMpeg, args Args) ([]byte, error) {
	cmd := f.Command(ctx, args)

	return cmd.Output()
}
