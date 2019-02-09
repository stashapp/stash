package paths

import (
	"github.com/stashapp/stash/internal/manager/jsonschema"
	"path/filepath"
)

type galleryPaths struct {
	config *jsonschema.Config
}

func newGalleryPaths(c *jsonschema.Config) *galleryPaths {
	gp := galleryPaths{}
	gp.config = c
	return &gp
}

func (gp *galleryPaths) GetExtractedPath(checksum string) string {
	return filepath.Join(gp.config.Cache, checksum)
}

func (gp *galleryPaths) GetExtractedFilePath(checksum string, fileName string) string {
	return filepath.Join(gp.config.Cache, checksum, fileName)
}
