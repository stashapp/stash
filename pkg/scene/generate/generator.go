package generate

import (
	"context"
	"fmt"
	"os"

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

	GetScreenshotPath(checksum string) string
	GetThumbnailScreenshotPath(checksum string) string

	GetSpriteImageFilePath(checksum string) string
	GetSpriteVttFilePath(checksum string) string

	GetTranscodePath(checksum string) string
}

type Generator struct {
	Encoder     ffmpeg.FFMpeg
	LockManager *fsutil.ReadLockManager
	MarkerPaths MarkerPaths
	ScenePaths  ScenePaths
	Overwrite   bool
}

type generateFn func(ctx context.Context, tmpFn string) error

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
func (g Generator) generateFile(ctx context.Context, p Paths, pattern string, output string, generateFn generateFn) error {
	tmpFile, err := g.tempFile(p, pattern) // tmp output in case the process ends abruptly
	if err != nil {
		return err
	}

	tmpFn := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpFn)
	}()

	if err := generateFn(ctx, tmpFn); err != nil {
		return err
	}

	if err := fsutil.SafeMove(tmpFn, output); err != nil {
		return fmt.Errorf("moving %s to %s", tmpFn, output)
	}

	return nil
}
