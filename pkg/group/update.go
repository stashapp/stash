package group

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type ImageInput struct {
	Image []byte
	Set   bool
}

func (s *Service) UpdatePartial(ctx context.Context, id int, updatedGroup models.GroupPartial, frontImage ImageInput, backImage ImageInput) (*models.Group, error) {
	if err := s.validateUpdate(ctx, id, updatedGroup); err != nil {
		return nil, err
	}

	r := s.Repository

	group, err := r.UpdatePartial(ctx, id, updatedGroup)
	if err != nil {
		return nil, err
	}

	// update image table
	if frontImage.Set {
		if err := r.UpdateFrontImage(ctx, id, frontImage.Image); err != nil {
			return nil, err
		}
	}

	if backImage.Set {
		if err := r.UpdateBackImage(ctx, id, backImage.Image); err != nil {
			return nil, err
		}
	}

	return group, nil
}
