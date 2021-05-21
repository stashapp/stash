package tag

import (
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

// EnsureTagNameUnique returns an error if the tag name provided
// is used as a name or alias of another existing tag.
func EnsureTagNameUnique(id int, name string, qb models.TagReader) error {
	// ensure name is unique
	sameNameTag, err := ByName(qb, name)
	if err != nil {
		return err
	}

	if sameNameTag != nil && id != sameNameTag.ID {
		return fmt.Errorf("tag with name '%s' already exists", name)
	}

	// query by alias
	sameNameTag, err = ByAlias(qb, name)
	if err != nil {
		return err
	}

	if sameNameTag != nil && id != sameNameTag.ID {
		return fmt.Errorf("name '%s' is used as alias for '%s'", name, sameNameTag.Name)
	}

	return nil
}
