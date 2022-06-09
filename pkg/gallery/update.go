package gallery

import (
	"context"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

type PartialUpdater interface {
	UpdatePartial(ctx context.Context, id int, updatedGallery models.GalleryPartial) (*models.Gallery, error)
}

type ImageUpdater interface {
	GetImageIDs(ctx context.Context, galleryID int) ([]int, error)
	UpdateImages(ctx context.Context, galleryID int, imageIDs []int) error
}

func UpdateFileModTime(ctx context.Context, qb PartialUpdater, id int, modTime time.Time) (*models.Gallery, error) {
	return qb.UpdatePartial(ctx, id, models.GalleryPartial{
		FileModTime: models.NewOptionalTime(modTime),
	})
}

func AddImage(ctx context.Context, qb ImageUpdater, galleryID int, imageID int) error {
	imageIDs, err := qb.GetImageIDs(ctx, galleryID)
	if err != nil {
		return err
	}

	imageIDs = intslice.IntAppendUnique(imageIDs, imageID)
	return qb.UpdateImages(ctx, galleryID, imageIDs)
}

func AddPerformer(ctx context.Context, qb PartialUpdater, o *models.Gallery, performerID int) (bool, error) {
	if !intslice.IntInclude(o.PerformerIDs, performerID) {
		if _, err := qb.UpdatePartial(ctx, o.ID, models.GalleryPartial{
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

func AddTag(ctx context.Context, qb PartialUpdater, o *models.Gallery, tagID int) (bool, error) {
	if !intslice.IntInclude(o.TagIDs, tagID) {
		if _, err := qb.UpdatePartial(ctx, o.ID, models.GalleryPartial{
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
