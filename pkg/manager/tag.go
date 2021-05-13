package manager

import (
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

func EnsureTagNameUnique(tag models.Tag, qb models.TagReader) error {
	// ensure name is unique
	sameNameTag, err := qb.FindByName(tag.Name, true)
	if err != nil {
		return err
	}

	if sameNameTag != nil && tag.ID != sameNameTag.ID {
		return fmt.Errorf("Tag with name '%s' already exists", tag.Name)
	}

	return nil
}
