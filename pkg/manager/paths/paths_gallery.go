package paths

import (
	"fmt"
	"github.com/stashapp/stash/pkg/manager/config"
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

//return a string that can be added to filepath.Join to implement directory depth
func getIntraDir(checksum string) string {
	if thumbDirDepth == 0 {
		return "" // filepath.Join ignores empty elements
	}
	intraDir := checksum[0:thumbDirLength]
	for i := 1; i < thumbDirDepth; i++ {
		intraDir = filepath.Join(intraDir, checksum[thumbDirLength*i:thumbDirLength*(i+1)])
	}
	return intraDir
}

func GetGthumbDir(checksum string) string {
	return filepath.Join(config.GetCachePath(), thumbDir, getIntraDir(checksum), checksum)
}

func GetGthumbPath(checksum string, index int, width int) string {
	fname := fmt.Sprintf("%s_%d_%d.jpg", checksum, index, width)
	return filepath.Join(config.GetCachePath(), thumbDir, getIntraDir(checksum), checksum, fname)
}

func (gp *galleryPaths) GetExtractedFilePath(checksum string, fileName string) string {
	return filepath.Join(config.GetCachePath(), checksum, fileName)
}
