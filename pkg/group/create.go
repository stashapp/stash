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

func (s *Service) Create(ctx context.Context, input *models.CreateGroupInput) error {
	r := s.Repository
	group := input.Group

	if err := s.validateCreate(ctx, group); err != nil {
		return err
	}

	err := r.Create(ctx, input.Group)
	if err != nil {
		return err
	}

	// set custom fields
	if len(input.CustomFields) > 0 {
		if err := r.SetCustomFields(ctx, group.ID, models.CustomFieldsInput{
			Full: input.CustomFields,
		}); err != nil {
			return err
		}
	}

	// update image table
	if len(input.FrontImageData) > 0 {
		if err := r.UpdateFrontImage(ctx, group.ID, input.FrontImageData); err != nil {
			return err
		}
	}

	if len(input.BackImageData) > 0 {
		if err := r.UpdateBackImage(ctx, group.ID, input.BackImageData); err != nil {
			return err
		}
	}

	return nil
}
