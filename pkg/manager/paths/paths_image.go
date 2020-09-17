package paths

import (
	"fmt"
	"path/filepath"

	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/utils"
)

type imagePaths struct{}

const imageThumbDir = "ithumbs"

func newImagePaths() *imagePaths {
	return &imagePaths{}
}

func (gp *imagePaths) GetExtractedPath(checksum string) string {
	return filepath.Join(config.GetCachePath(), checksum)
}

func GetImageThumbCache() string {
	return filepath.Join(config.GetCachePath(), imageThumbDir)
}

func GetImageThumbDir(checksum string) string {
	return filepath.Join(config.GetCachePath(), imageThumbDir, utils.GetIntraDir(checksum, thumbDirDepth, thumbDirLength))
}

func GetImageThumbPath(checksum string, width int) string {
	fname := fmt.Sprintf("%s_%d.jpg", checksum, width)
	return filepath.Join(config.GetCachePath(), imageThumbDir, utils.GetIntraDir(checksum, thumbDirDepth, thumbDirLength), fname)
}
