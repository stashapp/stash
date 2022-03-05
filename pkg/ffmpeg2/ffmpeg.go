package ffmpeg2

import (
	"context"
	"os/exec"

	stashExec "github.com/stashapp/stash/pkg/exec"
)

type FFMpeg string

func (f *FFMpeg) Command(ctx context.Context, args []string) *exec.Cmd {
	return stashExec.CommandContext(ctx, string(*f), args...)
}
