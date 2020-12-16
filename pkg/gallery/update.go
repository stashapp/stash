package gallery

import (
	"github.com/stashapp/stash/pkg/models"	
	"github.com/stashapp/stash/pkg/utils"
)


func UpdateFileModTime(qb models.GalleryWriter, id int, modTime models.NullSQLiteTimestamp) (*models.Gallery, error) {
	return qb.UpdatePartial(models.GalleryPartial{
		ID:          id,
		FileModTime: &modTime,
	})
}

func AddImage(qb models.GalleryReaderWriter, galleryID int, imageID int) error {
	imageIDs, err := qb.GetImageIDs(galleryID)
	if err != nil {
		return err
	}

	imageIDs = utils.IntAppendUnique(imageIDs, imageID)
	return qb.UpdateImages(galleryID, imageIDs)
}