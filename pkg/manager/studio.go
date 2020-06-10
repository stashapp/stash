package manager

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

func ValidateModifyStudio(studio models.StudioPartial, tx *sqlx.Tx) error {
	if studio.ParentID == nil || !studio.ParentID.Valid {
		return nil
	}

	// ensure there is no cyclic dependency
	thisID := studio.ID
	qb := models.NewStudioQueryBuilder()

	currentParentID := *studio.ParentID

	for currentParentID.Valid {
		if currentParentID.Int64 == int64(thisID) {
			return errors.New("studio cannot be an ancestor of itself")
		}

		currentStudio, err := qb.Find(int(currentParentID.Int64), tx)
		if err != nil {
			return fmt.Errorf("error finding parent studio: %s", err.Error())
		}

		currentParentID = currentStudio.ParentID
	}

	return nil
}
