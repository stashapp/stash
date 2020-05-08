package paths

import (
	"fmt"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/utils"
	"path/filepath"
)

type galleryPaths struct{}

const thumbDir = "gthumbs"
const thumbDirDepth int = 1
const thumbDirLength int = 2 // thumbDirDepth * thumbDirLength must be smaller than the length of checksum

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
	return filepath.Join(config.GetCachePath(), thumbDir, utils.GetIntraDir(checksum, thumbDirDepth, thumbDirLength), checksum)
}

func GetGthumbPath(checksum string, index int, width int) string {
	fname := fmt.Sprintf("%s_%d_%d.jpg", checksum, index, width)
	return filepath.Join(config.GetCachePath(), thumbDir, utils.GetIntraDir(checksum, thumbDirDepth, thumbDirLength), checksum, fname)
}

func (gp *galleryPaths) GetExtractedFilePath(checksum string, fileName string) string {
	return filepath.Join(config.GetCachePath(), checksum, fileName)
}
