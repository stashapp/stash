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

func (s *Service) AddSubGroups(ctx context.Context, groupID int, subGroups []models.GroupIDDescription, insertIndex *int) error {
	// get the group
	existing, err := s.Repository.Find(ctx, groupID)
	if err != nil {
		return err
	}

	// ensure it exists
	if existing == nil {
		return models.ErrNotFound
	}

	// validate the hierarchy
	d := &models.UpdateGroupDescriptions{
		Groups: subGroups,
		Mode:   models.RelationshipUpdateModeAdd,
	}
	if err := s.validateUpdateGroupHierarchy(ctx, existing, nil, d); err != nil {
		return err
	}

	// validate insert index
	if insertIndex != nil && *insertIndex < 0 {
		return ErrInvalidInsertIndex
	}

	// add the subgroups
	return s.Repository.AddSubGroups(ctx, groupID, subGroups, insertIndex)
}

func (s *Service) RemoveSubGroups(ctx context.Context, groupID int, subGroupIDs []int) error {
	// get the group
	existing, err := s.Repository.Find(ctx, groupID)
	if err != nil {
		return err
	}

	// ensure it exists
	if existing == nil {
		return models.ErrNotFound
	}

	// add the subgroups
	return s.Repository.RemoveSubGroups(ctx, groupID, subGroupIDs)
}
