package tag

import (
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type NameExistsError struct {
	Name string
}

func (e *NameExistsError) Error() string {
	return fmt.Sprintf("tag with name '%s' already exists", e.Name)
}

type NameUsedByAliasError struct {
	Name     string
	OtherTag string
}

func (e *NameUsedByAliasError) Error() string {
	return fmt.Sprintf("name '%s' is used as alias for '%s'", e.Name, e.OtherTag)
}

// EnsureTagNameUnique returns an error if the tag name provided
// is used as a name or alias of another existing tag.
func EnsureTagNameUnique(id int, name string, qb models.TagReader) error {
	// ensure name is unique
	sameNameTag, err := ByName(qb, name)
	if err != nil {
		return err
	}

	if sameNameTag != nil && id != sameNameTag.ID {
		return &NameExistsError{
			Name: name,
		}
	}

	// query by alias
	sameNameTag, err = ByAlias(qb, name)
	if err != nil {
		return err
	}

	if sameNameTag != nil && id != sameNameTag.ID {
		return &NameUsedByAliasError{
			Name:     name,
			OtherTag: sameNameTag.Name,
		}
	}

	return nil
}

func EnsureAliasesUnique(id int, aliases []string, qb models.TagReader) error {
	for _, a := range aliases {
		if err := EnsureTagNameUnique(id, a, qb); err != nil {
			return err
		}
	}

	return nil
}
