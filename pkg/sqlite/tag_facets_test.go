//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagFacets_ReturnsParents(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		tqb := db.Tag

		facets, err := tqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have parent tags (from tagParentLinks in setup_test.go)
		assert.Greater(t, len(facets.Parents), 0, "Should return at least one parent tag")

		// Verify parents are sorted by count descending
		for i := 1; i < len(facets.Parents); i++ {
			assert.GreaterOrEqual(t, facets.Parents[i-1].Count, facets.Parents[i].Count,
				"Parents should be sorted by count descending")
		}

		return nil
	})
}

func TestTagFacets_ReturnsChildren(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		tqb := db.Tag

		facets, err := tqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have child tags (from tagParentLinks in setup_test.go)
		assert.Greater(t, len(facets.Children), 0, "Should return at least one child tag")

		// Verify children are sorted by count descending
		for i := 1; i < len(facets.Children); i++ {
			assert.GreaterOrEqual(t, facets.Children[i-1].Count, facets.Children[i].Count,
				"Children should be sorted by count descending")
		}

		return nil
	})
}

func TestTagFacets_ReturnsFavorite(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		tqb := db.Tag

		facets, err := tqb.GetFacets(ctx, nil, 100)
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

func TestTagFacets_RespectsLimit(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		tqb := db.Tag

		limit := 5
		facets, err := tqb.GetFacets(ctx, nil, limit)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should not exceed limit
		assert.LessOrEqual(t, len(facets.Parents), limit, "Parents should not exceed limit")
		assert.LessOrEqual(t, len(facets.Children), limit, "Children should not exceed limit")

		return nil
	})
}

func TestTagFacets_NoFilter(t *testing.T) {
	// Note: TagStore.GetFacets doesn't use the filter parameter currently
	// This test verifies the query still works with a nil filter
	withRollbackTxn(func(ctx context.Context) error {
		tqb := db.Tag

		facets, err := tqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should return results
		assert.NotNil(t, facets, "Facets should not be nil")

		return nil
	})
}
