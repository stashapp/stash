package gallery

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type PartialUpdater interface {
	UpdatePartial(ctx context.Context, id int, updatedGallery models.GalleryPartial) (*models.Gallery, error)
}

type ImageUpdater interface {
	GetImageIDs(ctx context.Context, galleryID int) ([]int, error)
	AddImages(ctx context.Context, galleryID int, imageIDs ...int) error
	RemoveImages(ctx context.Context, galleryID int, imageIDs ...int) error
}

// AddImages adds images to the provided gallery.
// It returns an error if the gallery does not support adding images, or if
// the operation fails.
func (s *Service) AddImages(ctx context.Context, g *models.Gallery, toAdd ...int) error {
	if err := validateContentChange(g); err != nil {
		return err
	}

	return s.Repository.AddImages(ctx, g.ID, toAdd...)
}

// RemoveImages removes images from the provided gallery.
// It does not validate if the images are part of the gallery.
// It returns an error if the gallery does not support removing images, or if
// the operation fails.
func (s *Service) RemoveImages(ctx context.Context, g *models.Gallery, toRemove ...int) error {
	if err := validateContentChange(g); err != nil {
		return err
	}

	return s.Repository.RemoveImages(ctx, g.ID, toRemove...)
}

func AddPerformer(ctx context.Context, qb PartialUpdater, o *models.Gallery, performerID int) error {
	_, err := qb.UpdatePartial(ctx, o.ID, models.GalleryPartial{
		PerformerIDs: &models.UpdateIDs{
			IDs:  []int{performerID},
			Mode: models.RelationshipUpdateModeAdd,
		},
	})
	return err
}

func AddTag(ctx context.Context, qb PartialUpdater, o *models.Gallery, tagID int) error {
	_, err := qb.UpdatePartial(ctx, o.ID, models.GalleryPartial{
		TagIDs: &models.UpdateIDs{
			IDs:  []int{tagID},
			Mode: models.RelationshipUpdateModeAdd,
		},
	})
	return err
}
