package image

import (
	"context"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

type PartialUpdater interface {
	UpdatePartial(ctx context.Context, id int, partial models.ImagePartial) (*models.Image, error)
}

func UpdateFileModTime(ctx context.Context, qb PartialUpdater, id int, modTime time.Time) (*models.Image, error) {
	return qb.UpdatePartial(ctx, id, models.ImagePartial{
		FileModTime: models.NewOptionalTime(modTime),
	})
}

func AddPerformer(ctx context.Context, qb PartialUpdater, i *models.Image, performerID int) (bool, error) {
	if !intslice.IntInclude(i.PerformerIDs, performerID) {
		if _, err := qb.UpdatePartial(ctx, i.ID, models.ImagePartial{
			PerformerIDs: &models.UpdateIDs{
				IDs:  []int{performerID},
				Mode: models.RelationshipUpdateModeAdd,
			},
		}); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func AddTag(ctx context.Context, qb PartialUpdater, i *models.Image, tagID int) (bool, error) {
	if !intslice.IntInclude(i.TagIDs, tagID) {
		if _, err := qb.UpdatePartial(ctx, i.ID, models.ImagePartial{
			TagIDs: &models.UpdateIDs{
				IDs:  []int{tagID},
				Mode: models.RelationshipUpdateModeAdd,
			},
		}); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
