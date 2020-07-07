package manager

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

func EnsureTagNameUnique(name string, tx *sqlx.Tx) error {
	qb := models.NewTagQueryBuilder()

	// ensure name is unique
	sameNameTag, err := qb.FindByName(name, tx, true)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if sameNameTag != nil {
		return fmt.Errorf("Tag with name '%s' already exists", name)
	}

	return nil
}
