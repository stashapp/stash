package manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/studio"
)

// Checks to make sure that:
// 1. The studio exists locally
// 2. If the studio has a parent, it is not itself
// 3. If the studio has a parent, it exists locally and the parent does not have the studio as its parent
func ValidateModifyStudio(ctx context.Context, studio models.StudioPartial, qb studio.Finder) error {
	existing, err := qb.Find(ctx, studio.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("studio with id %d not found", studio.ID)
	}

	if !studio.ParentID.Set || studio.ParentID.Null {
		return nil
	}

	// ensure there is no cyclic dependency
	thisID := studio.ID
	currentParentID := studio.ParentID

	if currentParentID.Value == thisID {
		return errors.New("studio cannot be an ancestor of itself")
	}

	parentStudio, err := qb.Find(ctx, currentParentID.Value)
	if err != nil {
		return fmt.Errorf("error finding parent studio: %v", err)
	} else if parentStudio.ParentID == &thisID {
		return errors.New("studio is already parent studio of the new parent studio")
	}

	return nil
}
