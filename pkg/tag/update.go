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

type InvalidTagHierarchyError struct {
	Direction   string
	InvalidTag  string
	ApplyingTag string
}

func (e *InvalidTagHierarchyError) Error() string {
	return fmt.Sprintf("Cannot apply tag \"%s\" as it is already a %s of \"%s\"", e.InvalidTag, e.Direction, e.ApplyingTag)
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

func EnsureUniqueHierarchy(id int, parentIDs, childIDs []int, qb models.TagReader) error {
	allAncestors := make(map[int]*models.Tag)
	allDescendants := make(map[int]*models.Tag)
	excludeIDs := []int{id}

	parentsAncestors, err := qb.FindAllAncestors(id, excludeIDs)
	if err != nil {
		return err
	}

	for _, ancestorTag := range parentsAncestors {
		allAncestors[ancestorTag.ID] = ancestorTag
	}

	childsDescendants, err := qb.FindAllDescendants(id, excludeIDs)
	if err != nil {
		return err
	}

	for _, descendentTag := range childsDescendants {
		allDescendants[descendentTag.ID] = descendentTag
	}

	validateParent := func(testID, applyingID int) error {
		if parentTag, exists := allDescendants[testID]; exists {
			applyingTag, err := qb.Find(applyingID)

			if err != nil {
				return nil
			}

			return &InvalidTagHierarchyError{
				Direction:   "parent or ancestor",
				InvalidTag:  parentTag.Name,
				ApplyingTag: applyingTag.Name,
			}
		}

		return nil
	}

	validateChild := func(testID, applyingID int) error {
		if childTag, exists := allAncestors[testID]; exists {
			applyingTag, err := qb.Find(applyingID)

			if err != nil {
				return nil
			}

			return &InvalidTagHierarchyError{
				Direction:   "child or descendent",
				InvalidTag:  childTag.Name,
				ApplyingTag: applyingTag.Name,
			}
		}

		return nil
	}

	if parentIDs == nil {
		parentTags, err := qb.FindByChildTagID(id)
		if err != nil {
			return err
		}

		for _, parentTag := range parentTags {
			parentIDs = append(parentIDs, parentTag.ID)
		}
	}

	if childIDs == nil {
		childTags, err := qb.FindByParentTagID(id)
		if err != nil {
			return err
		}

		for _, childTag := range childTags {
			childIDs = append(childIDs, childTag.ID)
		}
	}

	for _, parentID := range parentIDs {
		if err := validateParent(parentID, id); err != nil {
			return err
		}
	}

	for _, childID := range childIDs {
		if err := validateChild(childID, id); err != nil {
			return err
		}
	}

	return nil
}

func MergeHierarchy(destination int, sources []int, qb models.TagReader) ([]int, []int, error) {
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
		parents, err := qb.FindByChildTagID(id)
		if err != nil {
			return nil, nil, err
		}

		mergedParents = addTo(mergedParents, parents)

		children, err := qb.FindByParentTagID(id)
		if err != nil {
			return nil, nil, err
		}

		mergedChildren = addTo(mergedChildren, children)
	}

	return mergedParents, mergedChildren, nil
}
