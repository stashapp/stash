package paths

import (
	"fmt"
	"github.com/stashapp/stash/pkg/manager/config"
	"path/filepath"
)

type galleryPaths struct{}

const thumbDir = "gthumbs"

func newGalleryPaths() *galleryPaths {
	return &galleryPaths{}
}

func (gp *galleryPaths) GetExtractedPath(checksum string) string {
	return filepath.Join(config.GetCachePath(), checksum)
}

func GetGthumbCache() string {
	return filepath.Join(config.GetCachePath(), thumbDir)
}

func GetGthumbDir(checksum string) string {
	return filepath.Join(config.GetCachePath(), thumbDir, checksum)
}

func GetGthumbPath(checksum string, index int, width int) string {
	fname := fmt.Sprintf("%s_%d_%d.jpg", checksum, index, width)
	return filepath.Join(config.GetCachePath(), thumbDir, checksum, fname)
}

func (gp *galleryPaths) GetExtractedFilePath(checksum string, fileName string) string {
	return filepath.Join(config.GetCachePath(), checksum, fileName)
}
