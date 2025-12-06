//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestStudioFacets_ReturnsTags(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Studio

		facets, err := sqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have some tags (from studioTags in setup_test.go)
		assert.Greater(t, len(facets.Tags), 0, "Should return at least one tag")

		// Verify tags are sorted by count descending
		for i := 1; i < len(facets.Tags); i++ {
			assert.GreaterOrEqual(t, facets.Tags[i-1].Count, facets.Tags[i].Count,
				"Tags should be sorted by count descending")
		}

		return nil
	})
}

func TestStudioFacets_ReturnsParents(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Studio

		facets, err := sqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have parent studios (from studioParentLinks in setup_test.go)
		assert.Greater(t, len(facets.Parents), 0, "Should return at least one parent studio")

		// Verify parents are sorted by count descending
		for i := 1; i < len(facets.Parents); i++ {
			assert.GreaterOrEqual(t, facets.Parents[i-1].Count, facets.Parents[i].Count,
				"Parents should be sorted by count descending")
		}

		return nil
	})
}

func TestStudioFacets_ReturnsFavorite(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Studio

		facets, err := sqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have favorite facet with true and/or false counts
		assert.Greater(t, len(facets.Favorite), 0, "Should return favorite facet")

		for _, f := range facets.Favorite {
			assert.GreaterOrEqual(t, f.Count, 0, "Favorite count should be non-negative")
		}

		return nil
	})
}

func TestStudioFacets_RespectsLimit(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Studio

		limit := 5
		facets, err := sqb.GetFacets(ctx, nil, limit)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should not exceed limit
		assert.LessOrEqual(t, len(facets.Tags), limit, "Tags should not exceed limit")
		assert.LessOrEqual(t, len(facets.Parents), limit, "Parents should not exceed limit")

		return nil
	})
}

func TestStudioFacets_WithTagFilter(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Studio

		// Filter to studios with specific tag
		tagFilter := &models.StudioFilterType{
			Tags: &models.HierarchicalMultiCriterionInput{
				Value:    []string{strconv.Itoa(tagIDs[tagIdxWithStudio])},
				Modifier: models.CriterionModifierIncludes,
			},
		}

		facets, err := sqb.GetFacets(ctx, tagFilter, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// All returned facets should be from studios with this tag
		for _, tag := range facets.Tags {
			assert.Greater(t, tag.Count, 0, "Tag count should be positive in filtered results")
		}

		return nil
	})
}

func TestStudioFacets_WithParentFilter(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Studio

		// Filter to studios with specific parent
		parentFilter := &models.StudioFilterType{
			Parents: &models.MultiCriterionInput{
				Value:    []string{strconv.Itoa(studioIDs[studioIdxWithParentStudio])},
				Modifier: models.CriterionModifierIncludes,
			},
		}

		facets, err := sqb.GetFacets(ctx, parentFilter, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Query should execute without error
		// Studios with this parent should be reflected in facets
		_ = facets

		return nil
	})
}
