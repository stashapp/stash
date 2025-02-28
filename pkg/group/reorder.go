package group

import (
	"context"
	"errors"

	"github.com/stashapp/stash/pkg/models"
)

var ErrInvalidInsertIndex = errors.New("invalid insert index")

func (s *Service) ReorderSubGroups(ctx context.Context, groupID int, subGroupIDs []int, insertPointID int, insertAfter bool) error {
	// get the group
	existing, err := s.Repository.Find(ctx, groupID)
	if err != nil {
		return err
	}

	// ensure it exists
	if existing == nil {
		return models.ErrNotFound
	}

	// TODO - ensure the subgroups exist in the group

	// ensure the insert index is valid
	if insertPointID < 0 {
		return ErrInvalidInsertIndex
	}

	// reorder the subgroups
	return s.Repository.ReorderSubGroups(ctx, groupID, subGroupIDs, insertPointID, insertAfter)
}
