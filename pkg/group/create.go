package group

import (
	"context"
	"errors"

	"github.com/stashapp/stash/pkg/models"
)

var (
	ErrEmptyName     = errors.New("name cannot be empty")
	ErrHierarchyLoop = errors.New("a group cannot be contained by one of its subgroups")
)

func (s *Service) Create(ctx context.Context, group *models.Group, frontimageData []byte, backimageData []byte) error {
	r := s.Repository

	if err := s.validateCreate(ctx, group); err != nil {
		return err
	}

	err := r.Create(ctx, group)
	if err != nil {
		return err
	}

	// update image table
	if len(frontimageData) > 0 {
		if err := r.UpdateFrontImage(ctx, group.ID, frontimageData); err != nil {
			return err
		}
	}

	if len(backimageData) > 0 {
		if err := r.UpdateBackImage(ctx, group.ID, backimageData); err != nil {
			return err
		}
	}

	return nil
}
