package manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/studio"
)

func ValidateModifyStudio(ctx context.Context, studio models.StudioPartial, qb studio.Finder) error {
	if studio.ParentID == nil || !studio.ParentID.Valid {
		return nil
	}

	// ensure there is no cyclic dependency
	thisID := studio.ID

	currentParentID := *studio.ParentID

	for currentParentID.Valid {
		if currentParentID.Int64 == int64(thisID) {
			return errors.New("studio cannot be an ancestor of itself")
		}

		currentStudio, err := qb.Find(ctx, int(currentParentID.Int64))
		if err != nil {
			return fmt.Errorf("error finding parent studio: %v", err)
		}

		currentParentID = currentStudio.ParentID
	}

	return nil
}
