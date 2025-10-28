package tag

import (
	"context"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

var (
	ErrNameMissing = errors.New("tag name must not be blank")
)

type NotFoundError struct {
	id int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("tag with id %d not found", e.id)
}

func ValidateCreate(ctx context.Context, tag models.Tag, qb models.TagReader) error {
	if tag.Name == "" {
		return ErrNameMissing
	}

	if err := EnsureTagNameUnique(ctx, 0, tag.Name, qb); err != nil {
		return err
	}

	if tag.Aliases.Loaded() {
		if err := EnsureAliasesUnique(ctx, tag.ID, tag.Aliases.List(), qb); err != nil {
			return err
		}
	}

	if len(tag.ParentIDs.List()) > 0 || len(tag.ChildIDs.List()) > 0 {
		if err := ValidateHierarchyNew(ctx, tag.ParentIDs.List(), tag.ChildIDs.List(), qb); err != nil {
			return err
		}
	}

	return nil
}

func ValidateUpdate(ctx context.Context, id int, partial models.TagPartial, qb models.TagReader) error {
	existing, err := qb.Find(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return &NotFoundError{id}
	}

	if partial.Name.Set {
		if partial.Name.Value == "" {
			return ErrNameMissing
		}

		if err := EnsureTagNameUnique(ctx, id, partial.Name.Value, qb); err != nil {
			return err
		}
	}

	if partial.Aliases != nil {
		if err := existing.LoadAliases(ctx, qb); err != nil {
			return err
		}

		if err := EnsureAliasesUnique(ctx, id, partial.Aliases.Apply(existing.Aliases.List()), qb); err != nil {
			return err
		}
	}

	if partial.ParentIDs != nil || partial.ChildIDs != nil {
		if err := existing.LoadParentIDs(ctx, qb); err != nil {
			return err
		}

		if err := existing.LoadChildIDs(ctx, qb); err != nil {
			return err
		}

		parentIDs := partial.ParentIDs
		if parentIDs == nil {
			parentIDs = &models.UpdateIDs{IDs: existing.ParentIDs.List(), Mode: models.RelationshipUpdateModeSet}
		}

		childIDs := partial.ChildIDs
		if childIDs == nil {
			childIDs = &models.UpdateIDs{IDs: existing.ChildIDs.List(), Mode: models.RelationshipUpdateModeSet}
		}

		if err := ValidateHierarchyExisting(ctx, existing, parentIDs.Apply(existing.ParentIDs.List()), childIDs.Apply(existing.ChildIDs.List()), qb); err != nil {
			return err
		}
	}

	return nil
}
