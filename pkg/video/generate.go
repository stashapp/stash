package video

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/stashapp/stash/pkg/ffmpeg2"
)

func Generate(f ffmpeg2.FFMpeg, ctx context.Context, args ffmpeg2.Args) error {
	cmd := f.Command(ctx, args)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		err.(*exec.ExitError).Stderr = stderr.Bytes()
		return err
	}

	return nil
}

func GenerateOutput(f ffmpeg2.FFMpeg, ctx context.Context, args ffmpeg2.Args) ([]byte, error) {
	cmd := f.Command(ctx, args)

	return cmd.Output()
}
