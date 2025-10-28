package tag

import (
	"context"
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

type InvalidTagHierarchyError struct {
	Direction       string
	CurrentRelation string
	InvalidTag      string
	ApplyingTag     string
	TagPath         string
}

func (e *InvalidTagHierarchyError) Error() string {
	if e.ApplyingTag == "" {
		return fmt.Sprintf("cannot apply tag \"%s\" as a %s of tag as it is already %s", e.InvalidTag, e.Direction, e.CurrentRelation)
	}

	return fmt.Sprintf("cannot apply tag \"%s\" as a %s of \"%s\" as it is already %s (%s)", e.InvalidTag, e.Direction, e.ApplyingTag, e.CurrentRelation, e.TagPath)
}

// EnsureTagNameUnique returns an error if the tag name provided
// is used as a name or alias of another existing tag.
func EnsureTagNameUnique(ctx context.Context, id int, name string, qb models.TagQueryer) error {
	// ensure name is unique
	sameNameTag, err := ByName(ctx, qb, name)
	if err != nil {
		return err
	}

	if sameNameTag != nil && id != sameNameTag.ID {
		return &NameExistsError{
			Name: name,
		}
	}

	// query by alias
	sameNameTag, err = ByAlias(ctx, qb, name)
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

func EnsureAliasesUnique(ctx context.Context, id int, aliases []string, qb models.TagQueryer) error {
	for _, a := range aliases {
		if err := EnsureTagNameUnique(ctx, id, a, qb); err != nil {
			return err
		}
	}

	return nil
}

type RelationshipFinder interface {
	FindAllAncestors(ctx context.Context, tagID int, excludeIDs []int) ([]*models.TagPath, error)
	FindAllDescendants(ctx context.Context, tagID int, excludeIDs []int) ([]*models.TagPath, error)
	models.TagRelationLoader
}

func ValidateHierarchyNew(ctx context.Context, parentIDs, childIDs []int, qb RelationshipFinder) error {
	allAncestors := make(map[int]*models.TagPath)
	allDescendants := make(map[int]*models.TagPath)

	for _, parentID := range parentIDs {
		parentsAncestors, err := qb.FindAllAncestors(ctx, parentID, nil)
		if err != nil {
			return err
		}

		for _, ancestorTag := range parentsAncestors {
			allAncestors[ancestorTag.ID] = ancestorTag
		}
	}

	for _, childID := range childIDs {
		childsDescendants, err := qb.FindAllDescendants(ctx, childID, nil)
		if err != nil {
			return err
		}

		for _, descendentTag := range childsDescendants {
			allDescendants[descendentTag.ID] = descendentTag
		}
	}

	// Validate that the tag is not a parent of any of its ancestors
	validateParent := func(testID int) error {
		if parentTag, exists := allDescendants[testID]; exists {
			return &InvalidTagHierarchyError{
				Direction:       "parent",
				CurrentRelation: "a descendant",
				InvalidTag:      parentTag.Name,
				TagPath:         parentTag.Path,
			}
		}

		return nil
	}

	// Validate that the tag is not a child of any of its ancestors
	validateChild := func(testID int) error {
		if childTag, exists := allAncestors[testID]; exists {
			return &InvalidTagHierarchyError{
				Direction:       "child",
				CurrentRelation: "an ancestor",
				InvalidTag:      childTag.Name,
				TagPath:         childTag.Path,
			}
		}

		return nil
	}

	for _, parentID := range parentIDs {
		if err := validateParent(parentID); err != nil {
			return err
		}
	}

	for _, childID := range childIDs {
		if err := validateChild(childID); err != nil {
			return err
		}
	}

	return nil
}

func ValidateHierarchyExisting(ctx context.Context, tag *models.Tag, parentIDs, childIDs []int, qb RelationshipFinder) error {
	allAncestors := make(map[int]*models.TagPath)
	allDescendants := make(map[int]*models.TagPath)

	parentsAncestors, err := qb.FindAllAncestors(ctx, tag.ID, nil)
	if err != nil {
		return err
	}

	for _, ancestorTag := range parentsAncestors {
		allAncestors[ancestorTag.ID] = ancestorTag
	}

	childsDescendants, err := qb.FindAllDescendants(ctx, tag.ID, nil)
	if err != nil {
		return err
	}

	for _, descendentTag := range childsDescendants {
		allDescendants[descendentTag.ID] = descendentTag
	}

	validateParent := func(testID int) error {
		if parentTag, exists := allDescendants[testID]; exists {
			return &InvalidTagHierarchyError{
				Direction:       "parent",
				CurrentRelation: "a descendant",
				InvalidTag:      parentTag.Name,
				ApplyingTag:     tag.Name,
				TagPath:         parentTag.Path,
			}
		}

		return nil
	}

	validateChild := func(testID int) error {
		if childTag, exists := allAncestors[testID]; exists {
			return &InvalidTagHierarchyError{
				Direction:       "child",
				CurrentRelation: "an ancestor",
				InvalidTag:      childTag.Name,
				ApplyingTag:     tag.Name,
				TagPath:         childTag.Path,
			}
		}

		return nil
	}

	for _, parentID := range parentIDs {
		if err := validateParent(parentID); err != nil {
			return err
		}
	}

	for _, childID := range childIDs {
		if err := validateChild(childID); err != nil {
			return err
		}
	}

	return nil
}

func MergeHierarchy(ctx context.Context, destination int, sources []int, qb RelationshipFinder) ([]int, []int, error) {
	var mergedParents, mergedChildren []int
	allIds := append([]int{destination}, sources...)

	addTo := func(mergedItems []int, tagIDs []int) []int {
	Tags:
		for _, tagID := range tagIDs {
			// Ignore tags which are already set
			for _, existingItem := range mergedItems {
				if tagID == existingItem {
					continue Tags
				}
			}

			// Ignore tags which are being merged, as these are rolled up anyway (if A is merged into B any direct link between them can be ignored)
			for _, id := range allIds {
				if tagID == id {
					continue Tags
				}
			}

			mergedItems = append(mergedItems, tagID)
		}

		return mergedItems
	}

	for _, id := range allIds {
		parents, err := qb.GetParentIDs(ctx, id)
		if err != nil {
			return nil, nil, err
		}

		mergedParents = addTo(mergedParents, parents)

		children, err := qb.GetChildIDs(ctx, id)
		if err != nil {
			return nil, nil, err
		}

		mergedChildren = addTo(mergedChildren, children)
	}

	return mergedParents, mergedChildren, nil
}
