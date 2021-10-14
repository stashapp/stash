package manager

import (
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

func ValidateModifyStudio(studio models.StudioPartial, qb models.StudioReader) error {
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

		currentStudio, err := qb.Find(int(currentParentID.Int64))
		if err != nil {
			return fmt.Errorf("error finding parent studio: %v", err)
		}

		currentParentID = currentStudio.ParentID
	}

	return nil
}
