package manager

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

func EnsureTagNameUnique(tag models.Tag, tx *sqlx.Tx) error {
	qb := models.NewTagQueryBuilder()

	// ensure name is unique
	sameNameTag, err := qb.FindByName(tag.Name, tx, true)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if sameNameTag != nil && tag.ID != sameNameTag.ID {
		return fmt.Errorf("Tag with name '%s' already exists", tag.Name)
	}

	return nil
}
