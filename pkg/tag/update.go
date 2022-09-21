package tag

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type NameFinderCreator interface {
	FindByNames(ctx context.Context, names []string, nocase bool) ([]*models.Tag, error)
	Create(ctx context.Context, newTag models.Tag) (*models.Tag, error)
}

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
	return fmt.Sprintf("cannot apply tag \"%s\" as a %s of \"%s\" as it is already %s (%s)", e.InvalidTag, e.Direction, e.ApplyingTag, e.CurrentRelation, e.TagPath)
}

// EnsureTagNameUnique returns an error if the tag name provided
// is used as a name or alias of another existing tag.
func EnsureTagNameUnique(ctx context.Context, id int, name string, qb Queryer) error {
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

func EnsureAliasesUnique(ctx context.Context, id int, aliases []string, qb Queryer) error {
	for _, a := range aliases {
		if err := EnsureTagNameUnique(ctx, id, a, qb); err != nil {
			return err
		}
	}

	return nil
}

type RelationshipGetter interface {
	FindAllAncestors(ctx context.Context, tagID int, excludeIDs []int) ([]*models.TagPath, error)
	FindAllDescendants(ctx context.Context, tagID int, excludeIDs []int) ([]*models.TagPath, error)
	FindByChildTagID(ctx context.Context, childID int) ([]*models.Tag, error)
	FindByParentTagID(ctx context.Context, parentID int) ([]*models.Tag, error)
}

func ValidateHierarchy(ctx context.Context, tag *models.Tag, parentIDs, childIDs []int, qb RelationshipGetter) error {
	id := tag.ID
	allAncestors := make(map[int]*models.TagPath)
	allDescendants := make(map[int]*models.TagPath)

	parentsAncestors, err := qb.FindAllAncestors(ctx, id, nil)
	if err != nil {
		return err
	}

	for _, ancestorTag := range parentsAncestors {
		allAncestors[ancestorTag.ID] = ancestorTag
	}

	childsDescendants, err := qb.FindAllDescendants(ctx, id, nil)
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

	if parentIDs == nil {
		parentTags, err := qb.FindByChildTagID(ctx, id)
		if err != nil {
			return err
		}

		for _, parentTag := range parentTags {
			parentIDs = append(parentIDs, parentTag.ID)
		}
	}

	if childIDs == nil {
		childTags, err := qb.FindByParentTagID(ctx, id)
		if err != nil {
			return err
		}

		for _, childTag := range childTags {
			childIDs = append(childIDs, childTag.ID)
		}
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

func MergeHierarchy(ctx context.Context, destination int, sources []int, qb RelationshipGetter) ([]int, []int, error) {
	var mergedParents, mergedChildren []int
	allIds := append([]int{destination}, sources...)

	addTo := func(mergedItems []int, tags []*models.Tag) []int {
	Tags:
		for _, tag := range tags {
			// Ignore tags which are already set
			for _, existingItem := range mergedItems {
				if tag.ID == existingItem {
					continue Tags
				}
			}

			// Ignore tags which are being merged, as these are rolled up anyway (if A is merged into B any direct link between them can be ignored)
			for _, id := range allIds {
				if tag.ID == id {
					continue Tags
				}
			}

			mergedItems = append(mergedItems, tag.ID)
		}

		return mergedItems
	}

	for _, id := range allIds {
		parents, err := qb.FindByChildTagID(ctx, id)
		if err != nil {
			return nil, nil, err
		}

		mergedParents = addTo(mergedParents, parents)

		children, err := qb.FindByParentTagID(ctx, id)
		if err != nil {
			return nil, nil, err
		}

		mergedChildren = addTo(mergedChildren, children)
	}

	return mergedParents, mergedChildren, nil
}
