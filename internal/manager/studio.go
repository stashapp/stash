package manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/studio"
)

func ValidateModifyStudio(ctx context.Context, studioID int, studio models.StudioPartial, qb studio.Finder) error {
	if studio.ParentID.Ptr() == nil {
		return nil
	}

	// ensure there is no cyclic dependency
	currentParentID := studio.ParentID.Ptr()

	for currentParentID != nil {
		if *currentParentID == studioID {
			return errors.New("studio cannot be an ancestor of itself")
		}

		currentStudio, err := qb.Find(ctx, *currentParentID)
		if err != nil {
			return fmt.Errorf("error finding parent studio: %v", err)
		}

		if currentStudio == nil {
			return fmt.Errorf("studio with id %d not found", *currentParentID)
		}

		currentParentID = currentStudio.ParentID
	}

	return nil
}
