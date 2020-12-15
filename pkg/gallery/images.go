package gallery

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func AddImage(qb models.GalleryReaderWriter, galleryID int, imageID int) error {
	imageIDs, err := qb.GetImageIDs(galleryID)
	if err != nil {
		return err
	}

	imageIDs = utils.IntAppendUnique(imageIDs, imageID)
	return qb.UpdateImages(galleryID, imageIDs)
}
