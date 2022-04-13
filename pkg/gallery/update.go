package gallery

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

func UpdateFileModTime(ctx context.Context, qb models.GalleryWriter, id int, modTime models.NullSQLiteTimestamp) (*models.Gallery, error) {
	return qb.UpdatePartial(ctx, models.GalleryPartial{
		ID:          id,
		FileModTime: &modTime,
	})
}

func AddImage(ctx context.Context, qb models.GalleryReaderWriter, galleryID int, imageID int) error {
	imageIDs, err := qb.GetImageIDs(ctx, galleryID)
	if err != nil {
		return err
	}

	imageIDs = intslice.IntAppendUnique(imageIDs, imageID)
	return qb.UpdateImages(ctx, galleryID, imageIDs)
}

func AddPerformer(ctx context.Context, qb models.GalleryReaderWriter, id int, performerID int) (bool, error) {
	performerIDs, err := qb.GetPerformerIDs(ctx, id)
	if err != nil {
		return false, err
	}

	oldLen := len(performerIDs)
	performerIDs = intslice.IntAppendUnique(performerIDs, performerID)

	if len(performerIDs) != oldLen {
		if err := qb.UpdatePerformers(ctx, id, performerIDs); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func AddTag(ctx context.Context, qb models.GalleryReaderWriter, id int, tagID int) (bool, error) {
	tagIDs, err := qb.GetTagIDs(ctx, id)
	if err != nil {
		return false, err
	}

	oldLen := len(tagIDs)
	tagIDs = intslice.IntAppendUnique(tagIDs, tagID)

	if len(tagIDs) != oldLen {
		if err := qb.UpdateTags(ctx, id, tagIDs); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
