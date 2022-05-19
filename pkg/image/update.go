package image

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

type PartialUpdater interface {
	Update(ctx context.Context, updatedImage models.ImagePartial) (*models.Image, error)
}

type PerformerUpdater interface {
	GetPerformerIDs(ctx context.Context, imageID int) ([]int, error)
	UpdatePerformers(ctx context.Context, imageID int, performerIDs []int) error
}

type TagUpdater interface {
	GetTagIDs(ctx context.Context, imageID int) ([]int, error)
	UpdateTags(ctx context.Context, imageID int, tagIDs []int) error
}

func UpdateFileModTime(ctx context.Context, qb PartialUpdater, id int, modTime models.NullSQLiteTimestamp) (*models.Image, error) {
	return qb.Update(ctx, models.ImagePartial{
		ID:          id,
		FileModTime: &modTime,
	})
}

func AddPerformer(ctx context.Context, qb PerformerUpdater, id int, performerID int) (bool, error) {
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

func AddTag(ctx context.Context, qb TagUpdater, id int, tagID int) (bool, error) {
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
