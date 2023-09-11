package image

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func AddPerformer(ctx context.Context, qb models.ImageUpdater, i *models.Image, performerID int) error {
	imagePartial := models.NewImagePartial()
	imagePartial.PerformerIDs = &models.UpdateIDs{
		IDs:  []int{performerID},
		Mode: models.RelationshipUpdateModeAdd,
	}
	_, err := qb.UpdatePartial(ctx, i.ID, imagePartial)
	return err
}

func AddTag(ctx context.Context, qb models.ImageUpdater, i *models.Image, tagID int) error {
	imagePartial := models.NewImagePartial()
	imagePartial.TagIDs = &models.UpdateIDs{
		IDs:  []int{tagID},
		Mode: models.RelationshipUpdateModeAdd,
	}
	_, err := qb.UpdatePartial(ctx, i.ID, imagePartial)
	return err
}
