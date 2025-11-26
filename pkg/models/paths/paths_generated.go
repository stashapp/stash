package paths

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

const thumbDirDepth int = 2
const thumbDirLength int = 2 // thumbDirDepth * thumbDirLength must be smaller than the length of checksum

type generatedPaths struct {
	Screenshots        string
	Thumbnails         string
	Vtt                string
	Markers            string
	Transcodes         string
	Downloads          string
	Tmp                string
	InteractiveHeatmap string
}

func newGeneratedPaths(path string) *generatedPaths {
	gp := generatedPaths{}
	gp.Screenshots = filepath.Join(path, "screenshots")
	gp.Thumbnails = filepath.Join(path, "thumbnails")
	gp.Vtt = filepath.Join(path, "vtt")
	gp.Markers = filepath.Join(path, "markers")
	gp.Transcodes = filepath.Join(path, "transcodes")
	gp.Downloads = filepath.Join(path, "download_stage")
	gp.Tmp = filepath.Join(path, "tmp")
	gp.InteractiveHeatmap = filepath.Join(path, "interactive_heatmaps")
	return &gp
}

func (gp *generatedPaths) GetTmpPath(fileName string) string {
	return filepath.Join(gp.Tmp, fileName)
}

// TempFile creates a temporary file using os.CreateTemp.
// It is the equivalent of calling os.CreateTemp using Tmp and pattern.
func (gp *generatedPaths) TempFile(pattern string) (*os.File, error) {
	if err := gp.EnsureTmpDir(); err != nil {
		logger.Warnf("Could not ensure existence of a temporary directory: %v", err)
	}
	return os.CreateTemp(gp.Tmp, pattern)
}

func (gp *generatedPaths) EnsureTmpDir() error {
	return fsutil.EnsureDir(gp.Tmp)
}

func (gp *generatedPaths) EmptyTmpDir() error {
	return fsutil.EmptyDir(gp.Tmp)
}

func (gp *generatedPaths) RemoveTmpDir() error {
	return fsutil.RemoveDir(gp.Tmp)
}

func (gp *generatedPaths) TempDir(pattern string) (string, error) {
	if err := gp.EnsureTmpDir(); err != nil {
		logger.Warnf("Could not ensure existence of a temporary directory: %v", err)
	}
	ret, err := os.MkdirTemp(gp.Tmp, pattern)
	if err != nil {
		return "", err
	}

	if err = fsutil.EmptyDir(ret); err != nil {
		logger.Warnf("could not recursively empty dir: %v", err)
	}

	return ret, nil
}

func (gp *generatedPaths) GetThumbnailPath(checksum string, width int) string {
	fname := fmt.Sprintf("%s_%d.jpg", checksum, width)
	return filepath.Join(gp.Thumbnails, fsutil.GetIntraDir(checksum, thumbDirDepth, thumbDirLength), fname)
}

func (gp *generatedPaths) GetClipPreviewPath(checksum string, width int) string {
	fname := fmt.Sprintf("%s_%d.webm", checksum, width)
	return filepath.Join(gp.Thumbnails, fsutil.GetIntraDir(checksum, thumbDirDepth, thumbDirLength), fname)
}
