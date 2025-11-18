package studio

import (
	"context"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

var (
	ErrNameMissing       = errors.New("studio name must not be blank")
	ErrEmptyAlias        = errors.New("studio alias must not be an empty string")
	ErrStudioOwnAncestor = errors.New("studio cannot be an ancestor of itself")
)

type NameExistsError struct {
	Name string
}

func (e *NameExistsError) Error() string {
	return fmt.Sprintf("studio with name '%s' already exists", e.Name)
}

type NameUsedByAliasError struct {
	Name        string
	OtherStudio string
}

func (e *NameUsedByAliasError) Error() string {
	return fmt.Sprintf("name '%s' is used as alias for '%s'", e.Name, e.OtherStudio)
}

// EnsureStudioNameUnique returns an error if the studio name provided
// is used as a name or alias of another existing tag.
func EnsureStudioNameUnique(ctx context.Context, id int, name string, qb models.StudioQueryer) error {
	// ensure name is unique
	sameNameStudio, err := ByName(ctx, qb, name)
	if err != nil {
		return err
	}

	if sameNameStudio != nil && id != sameNameStudio.ID {
		return &NameExistsError{
			Name: name,
		}
	}

	// query by alias
	sameNameStudio, err = ByAlias(ctx, qb, name)
	if err != nil {
		return err
	}

	if sameNameStudio != nil && id != sameNameStudio.ID {
		return &NameUsedByAliasError{
			Name:        name,
			OtherStudio: sameNameStudio.Name,
		}
	}

	return nil
}

func ValidateAliases(ctx context.Context, id int, aliases []string, qb models.StudioQueryer) error {
	for _, a := range aliases {
		if err := validateName(ctx, id, a, qb); err != nil {
			if errors.Is(err, ErrNameMissing) {
				return ErrEmptyAlias
			}
			return err
		}
	}

	return nil
}

func ValidateCreate(ctx context.Context, studio models.Studio, qb models.StudioQueryer) error {
	if err := validateName(ctx, 0, studio.Name, qb); err != nil {
		return err
	}

	if studio.Aliases.Loaded() && len(studio.Aliases.List()) > 0 {
		if err := ValidateAliases(ctx, 0, studio.Aliases.List(), qb); err != nil {
			return err
		}
	}

	return nil
}

func validateName(ctx context.Context, studioID int, name string, qb models.StudioQueryer) error {
	if name == "" {
		return ErrNameMissing
	}

	if err := EnsureStudioNameUnique(ctx, studioID, name, qb); err != nil {
		return err
	}

	return nil
}

type ValidateModifyReader interface {
	models.StudioGetter
	models.StudioQueryer
	models.AliasLoader
}

// Checks to make sure that:
// 1. The studio exists locally
// 2. The studio is not its own ancestor
// 3. The studio's aliases are unique
// 4. The name is unique
func ValidateModify(ctx context.Context, s models.StudioPartial, qb ValidateModifyReader) error {
	existing, err := qb.Find(ctx, s.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("studio with id %d not found", s.ID)
	}

	newParentID := s.ParentID.Ptr()

	if newParentID != nil {
		if err := validateParent(ctx, s.ID, *newParentID, qb); err != nil {
			return err
		}
	}

	if s.Aliases != nil {
		if err := existing.LoadAliases(ctx, qb); err != nil {
			return err
		}

		effectiveAliases := s.Aliases.Apply(existing.Aliases.List())
		if err := ValidateAliases(ctx, s.ID, effectiveAliases, qb); err != nil {
			return err
		}
	}

	if s.Name.Set && s.Name.Value != existing.Name {
		if err := validateName(ctx, s.ID, s.Name.Value, qb); err != nil {
			return err
		}
	}

	return nil
}

func validateParent(ctx context.Context, studioID int, newParentID int, qb models.StudioGetter) error {
	if newParentID == studioID {
		return ErrStudioOwnAncestor
	}

	// ensure there is no cyclic dependency
	parentStudio, err := qb.Find(ctx, newParentID)
	if err != nil {
		return fmt.Errorf("error finding parent studio: %v", err)
	}

	if parentStudio == nil {
		return fmt.Errorf("studio with id %d not found", newParentID)
	}

	if parentStudio.ParentID != nil {
		return validateParent(ctx, studioID, *parentStudio.ParentID, qb)
	}

	return nil
}
