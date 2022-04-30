package manager

import (
	"os"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

func DeleteGalleryFile(gallery *models.Gallery) {
	path := gallery.Path()
	if path != "" {
		err := os.Remove(path)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", path, err.Error())
		}
	}
}
