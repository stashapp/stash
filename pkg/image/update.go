package image

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type PartialUpdater interface {
	UpdatePartial(ctx context.Context, id int, partial models.ImagePartial) (*models.Image, error)
}

func AddPerformer(ctx context.Context, qb PartialUpdater, i *models.Image, performerID int) error {
	_, err := qb.UpdatePartial(ctx, i.ID, models.ImagePartial{
		PerformerIDs: &models.UpdateIDs{
			IDs:  []int{performerID},
			Mode: models.RelationshipUpdateModeAdd,
		},
	})

	return err
}

func AddTag(ctx context.Context, qb PartialUpdater, i *models.Image, tagID int) error {
	_, err := qb.UpdatePartial(ctx, i.ID, models.ImagePartial{
		TagIDs: &models.UpdateIDs{
			IDs:  []int{tagID},
			Mode: models.RelationshipUpdateModeAdd,
		},
	})
	return err
}
