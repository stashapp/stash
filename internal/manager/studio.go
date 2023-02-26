package manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/studio"
)

func ValidateModifyStudio(ctx context.Context, studio models.StudioPartial, qb studio.Finder) error {
	existing, err := qb.Find(ctx, studio.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("studio with id %d not found", studio.ID)
	}

	if !studio.ParentID.Set {
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
