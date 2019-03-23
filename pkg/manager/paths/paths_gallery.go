package paths

import (
	"github.com/stashapp/stash/pkg/manager/config"
	"path/filepath"
)

type galleryPaths struct{}

func newGalleryPaths() *galleryPaths {
	return &galleryPaths{}
}

func (gp *galleryPaths) GetExtractedPath(checksum string) string {
	return filepath.Join(config.GetCachePath(), checksum)
}

func (gp *galleryPaths) GetExtractedFilePath(checksum string, fileName string) string {
	return filepath.Join(config.GetCachePath(), checksum, fileName)
}
