package manager

import (
	"os"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

func DeleteGalleryFile(gallery *models.Gallery) {
	if gallery.Path != nil {
		err := os.Remove(*gallery.Path)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", *gallery.Path, err.Error())
		}
	}
}
