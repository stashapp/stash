// Package ffmpeg provides a wrapper around the ffmpeg and ffprobe executables.
package ffmpeg

import (
	"context"
	"os/exec"

	stashExec "github.com/stashapp/stash/pkg/exec"
	"github.com/stashapp/stash/pkg/logger"
)

// FFMpeg provides an interface to ffmpeg.
type FFMpeg struct {
	ffmpeg         string
	hwCodecSupport []VideoCodec
}

// Creates a new FFMpeg encoder
func NewEncoder(ffmpegPath string, ctx context.Context) *FFMpeg {
	ret := &FFMpeg{
		ffmpeg: ffmpegPath,
	}

	ret.hwCodecSupport = ret.findHWCodecs(ctx)

	logger.Info("Supported HW codecs: ")
	for _, codec := range ret.hwCodecSupport {
		logger.Info("\t", codec)
	}

	return ret
}

// Returns an exec.Cmd that can be used to run ffmpeg using args.
func (f *FFMpeg) Command(ctx context.Context, args []string) *exec.Cmd {
	return stashExec.CommandContext(ctx, string(f.ffmpeg), args...)
}
