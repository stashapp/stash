package gallery

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
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

	imageIDs = intslice.IntAppendUnique(imageIDs, imageID)
	return qb.UpdateImages(galleryID, imageIDs)
}

func AddPerformer(qb models.GalleryReaderWriter, id int, performerID int) (bool, error) {
	performerIDs, err := qb.GetPerformerIDs(id)
	if err != nil {
		return false, err
	}

	oldLen := len(performerIDs)
	performerIDs = intslice.IntAppendUnique(performerIDs, performerID)

	if len(performerIDs) != oldLen {
		if err := qb.UpdatePerformers(id, performerIDs); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func AddTag(qb models.GalleryReaderWriter, id int, tagID int) (bool, error) {
	tagIDs, err := qb.GetTagIDs(id)
	if err != nil {
		return false, err
	}

	oldLen := len(tagIDs)
	tagIDs = intslice.IntAppendUnique(tagIDs, tagID)

	if len(tagIDs) != oldLen {
		if err := qb.UpdateTags(id, tagIDs); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
