package group

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
)

var ErrEmptyAlias = errors.New("group alias must not be an empty string")

type DuplicateAliasError struct {
	Alias string
}

func (e *DuplicateAliasError) Error() string {
	return fmt.Sprintf("group contains duplicate alias '%s'", e.Alias)
}

func (s *Service) validateCreate(ctx context.Context, group *models.Group) error {
	if err := validateName(group.Name); err != nil {
		return err
	}

	if err := ValidateAliases(group.Name, group.Aliases); err != nil {
		return err
	}

	containingIDs := group.ContainingGroups.IDs()
	subIDs := group.SubGroups.IDs()

	if err := s.validateGroupHierarchy(ctx, containingIDs, subIDs); err != nil {
		return err
	}

	return nil
}

func (s *Service) validateUpdate(ctx context.Context, id int, partial models.GroupPartial) error {
	// get the existing group - ensure it exists
	existing, err := s.Repository.Find(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return models.ErrNotFound
	}

	if partial.Name.Set {
		if err := validateName(partial.Name.Value); err != nil {
			return err
		}
	}

	if err := existing.LoadAliases(ctx, s.Repository); err != nil {
		return err
	}
	if err := ValidateUpdateAliases(*existing, partial.Name, partial.Aliases); err != nil {
		return err
	}

	if err := s.validateUpdateGroupHierarchy(ctx, existing, partial.ContainingGroups, partial.SubGroups); err != nil {
		return err
	}

	return nil
}

func validateName(n string) error {
	// ensure name is not empty
	if strings.TrimSpace(n) == "" {
		return ErrEmptyName
	}

	return nil
}

// ValidateAliases returns an error if an alias is empty, duplicated, or the
// same as the group name. The comparison is case-insensitive.
func ValidateAliases(name string, aliases models.RelatedStrings) error {
	if !aliases.Loaded() {
		return nil
	}

	m := make(map[string]bool)
	m[strings.ToLower(name)] = true

	for _, alias := range aliases.List() {
		if strings.TrimSpace(alias) == "" {
			return ErrEmptyAlias
		}

		aliasL := strings.ToLower(alias)
		if m[aliasL] {
			return &DuplicateAliasError{alias}
		}
		m[aliasL] = true
	}

	return nil
}

func ValidateUpdateAliases(existing models.Group, name models.OptionalString, aliases *models.UpdateStrings) error {
	// if neither name nor aliases is set, don't check anything
	if !name.Set && aliases == nil {
		return nil
	}

	newName := existing.Name
	if name.Set {
		newName = name.Value
	}

	// If aliases is nil, we're only changing the name - check existing aliases against new name
	if aliases == nil {
		return ValidateAliases(newName, existing.Aliases)
	}

	newAliases := aliases.Apply(existing.Aliases.List())

	return ValidateAliases(newName, models.NewRelatedStrings(newAliases))
}

func (s *Service) validateGroupHierarchy(ctx context.Context, containingIDs []int, subIDs []int) error {
	// only need to validate if both are non-empty
	if len(containingIDs) == 0 || len(subIDs) == 0 {
		return nil
	}

	// ensure none of the containing groups are in the sub groups
	found, err := s.Repository.FindInAncestors(ctx, containingIDs, subIDs)
	if err != nil {
		return err
	}

	if len(found) > 0 {
		return ErrHierarchyLoop
	}

	return nil
}

func (s *Service) validateUpdateGroupHierarchy(ctx context.Context, existing *models.Group, containingGroups *models.UpdateGroupDescriptions, subGroups *models.UpdateGroupDescriptions) error {
	// no need to validate if there are no changes
	if containingGroups == nil && subGroups == nil {
		return nil
	}

	if err := existing.LoadContainingGroupIDs(ctx, s.Repository); err != nil {
		return err
	}
	existingContainingGroups := existing.ContainingGroups.List()

	if err := existing.LoadSubGroupIDs(ctx, s.Repository); err != nil {
		return err
	}
	existingSubGroups := existing.SubGroups.List()

	effectiveContainingGroups := existingContainingGroups
	if containingGroups != nil {
		effectiveContainingGroups = containingGroups.Apply(existingContainingGroups)
	}

	effectiveSubGroups := existingSubGroups
	if subGroups != nil {
		effectiveSubGroups = subGroups.Apply(existingSubGroups)
	}

	containingIDs := idsFromGroupDescriptions(effectiveContainingGroups)
	subIDs := idsFromGroupDescriptions(effectiveSubGroups)

	// ensure we haven't set the group as a subgroup of itself
	if slices.Contains(containingIDs, existing.ID) || slices.Contains(subIDs, existing.ID) {
		return ErrHierarchyLoop
	}

	return s.validateGroupHierarchy(ctx, containingIDs, subIDs)
}

func idsFromGroupDescriptions(v []models.GroupIDDescription) []int {
	return sliceutil.Map(v, func(g models.GroupIDDescription) int { return g.GroupID })
}
