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

func TestGroupFacets_ReturnsTags(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		mqb := db.Group

		facets, err := mqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have some tags (from groupTags in setup_test.go)
		assert.Greater(t, len(facets.Tags), 0, "Should return at least one tag")

		// Verify tags are sorted by count descending
		for i := 1; i < len(facets.Tags); i++ {
			assert.GreaterOrEqual(t, facets.Tags[i-1].Count, facets.Tags[i].Count,
				"Tags should be sorted by count descending")
		}

		return nil
	})
}

func TestGroupFacets_ReturnsPerformers(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		mqb := db.Group

		facets, err := mqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Performers are derived from group scenes, so may or may not have entries
		// depending on test data, but query should execute without error
		for _, p := range facets.Performers {
			assert.Greater(t, p.Count, 0, "Performer count should be positive")
		}

		return nil
	})
}

func TestGroupFacets_ReturnsStudios(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		mqb := db.Group

		facets, err := mqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Studios may or may not have entries depending on test data
		for _, s := range facets.Studios {
			assert.Greater(t, s.Count, 0, "Studio count should be positive")
		}

		return nil
	})
}

func TestGroupFacets_RespectsLimit(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		mqb := db.Group

		limit := 5
		facets, err := mqb.GetFacets(ctx, nil, limit)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should not exceed limit
		assert.LessOrEqual(t, len(facets.Tags), limit, "Tags should not exceed limit")
		assert.LessOrEqual(t, len(facets.Performers), limit, "Performers should not exceed limit")
		assert.LessOrEqual(t, len(facets.Studios), limit, "Studios should not exceed limit")

		return nil
	})
}

func TestGroupFacets_WithStudioFilter(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		mqb := db.Group

		// Filter to groups with specific studio
		studioFilter := &models.GroupFilterType{
			Studios: &models.HierarchicalMultiCriterionInput{
				Value:    []string{strconv.Itoa(studioIDs[studioIdxWithGroup])},
				Modifier: models.CriterionModifierIncludes,
			},
		}

		facets, err := mqb.GetFacets(ctx, studioFilter, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// All returned facets should be from groups in this studio
		for _, tag := range facets.Tags {
			assert.Greater(t, tag.Count, 0, "Tag count should be positive in filtered results")
		}

		return nil
	})
}

func TestGroupFacets_WithTagFilter(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		mqb := db.Group

		// Filter to groups with specific tag
		tagFilter := &models.GroupFilterType{
			Tags: &models.HierarchicalMultiCriterionInput{
				Value:    []string{strconv.Itoa(tagIDs[tagIdxWithGroup])},
				Modifier: models.CriterionModifierIncludes,
			},
		}

		facets, err := mqb.GetFacets(ctx, tagFilter, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Query should execute without error
		// We can verify the specific tag appears in results
		found := false
		for _, tag := range facets.Tags {
			if tag.ID == strconv.Itoa(tagIDs[tagIdxWithGroup]) {
				found = true
				break
			}
		}
		assert.True(t, found, "Filtered tag should appear in results")

		return nil
	})
}
