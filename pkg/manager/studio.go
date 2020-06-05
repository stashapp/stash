package manager

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

func ValidateModifyStudio(studio models.Studio, tx *sqlx.Tx) error {
	// ensure there is no cyclic dependency
	thisID := studio.ID
	qb := models.NewStudioQueryBuilder()

	currentStudio := &studio
	for currentStudio.ParentID.Valid {
		if currentStudio.ParentID.Int64 == int64(thisID) {
			return errors.New("studio cannot be an ancestor of itself")
		}

		var err error
		currentStudio, err = qb.Find(int(currentStudio.ParentID.Int64), tx)
		if err != nil {
			return fmt.Errorf("error finding parent studio: %s", err.Error())
		}
	}

	return nil
}
