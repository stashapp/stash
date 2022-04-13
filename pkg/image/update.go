package image

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

func UpdateFileModTime(ctx context.Context, qb models.ImageWriter, id int, modTime models.NullSQLiteTimestamp) (*models.Image, error) {
	return qb.Update(ctx, models.ImagePartial{
		ID:          id,
		FileModTime: &modTime,
	})
}

func AddPerformer(ctx context.Context, qb models.ImageReaderWriter, id int, performerID int) (bool, error) {
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

func AddTag(ctx context.Context, qb models.ImageReaderWriter, id int, tagID int) (bool, error) {
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
