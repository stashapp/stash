package generate

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
	"github.com/stashapp/stash/pkg/fsutil"
)

// ExtractSubtitle extracts the embedded subtitle stream at streamIndex from inputPath
// as WebVTT and writes it atomically to outputPath.
func (g Generator) ExtractSubtitle(ctx context.Context, inputPath string, streamIndex int, outputPath string) error {
	if g.Encoder == nil {
		return fmt.Errorf("ffmpeg encoder not configured")
	}

	if err := fsutil.EnsureDir(filepath.Dir(outputPath)); err != nil {
		return fmt.Errorf("ensuring caption cache directory: %w", err)
	}

	var lockCtx *fsutil.LockContext
	if g.LockManager != nil {
		lockCtx = g.LockManager.ReadLock(ctx, inputPath)
		defer lockCtx.Cancel()
	} else {
		lockCtx = &fsutil.LockContext{Context: ctx}
	}

	return g.generateFile(lockCtx, g.ScenePaths, vttPattern, outputPath, func(lockCtx *fsutil.LockContext, tmpFn string) error {
		args := transcoder.SubtitleExtract(inputPath, streamIndex, tmpFn)
		return g.generate(lockCtx, args)
	})
}
