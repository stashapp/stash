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

func TestPerformerFacets_ReturnsTags(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		pqb := db.Performer

		facets, err := pqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have some tags (from performerTags in setup_test.go)
		assert.Greater(t, len(facets.Tags), 0, "Should return at least one tag")

		// Verify tags are sorted by count descending
		for i := 1; i < len(facets.Tags); i++ {
			assert.GreaterOrEqual(t, facets.Tags[i-1].Count, facets.Tags[i].Count,
				"Tags should be sorted by count descending")
		}

		return nil
	})
}

func TestPerformerFacets_ReturnsGenders(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		pqb := db.Performer

		facets, err := pqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should return genders facet (even if empty, query should work)
		// Genders are enum-based, so they may or may not have entries depending on test data
		for _, g := range facets.Genders {
			assert.Greater(t, g.Count, 0, "Gender count should be positive")
			assert.True(t, g.Gender.IsValid(), "Gender should be valid enum value")
		}

		return nil
	})
}

func TestPerformerFacets_ReturnsFavorite(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		pqb := db.Performer

		facets, err := pqb.GetFacets(ctx, nil, 100)
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

func TestPerformerFacets_ReturnsCircumcised(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		pqb := db.Performer

		facets, err := pqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should return circumcised facet
		for _, c := range facets.Circumcised {
			assert.Greater(t, c.Count, 0, "Circumcised count should be positive")
			assert.True(t, c.Value.IsValid(), "Circumcised value should be valid enum")
		}

		return nil
	})
}

func TestPerformerFacets_ReturnsRatings(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		pqb := db.Performer

		facets, err := pqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Ratings should be sorted by rating descending
		for i := 1; i < len(facets.Ratings); i++ {
			assert.GreaterOrEqual(t, facets.Ratings[i-1].Rating, facets.Ratings[i].Rating,
				"Ratings should be sorted by rating descending")
		}

		return nil
	})
}

func TestPerformerFacets_RespectsLimit(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		pqb := db.Performer

		limit := 5
		facets, err := pqb.GetFacets(ctx, nil, limit)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should not exceed limit
		assert.LessOrEqual(t, len(facets.Tags), limit, "Tags should not exceed limit")
		assert.LessOrEqual(t, len(facets.Studios), limit, "Studios should not exceed limit")
		assert.LessOrEqual(t, len(facets.Countries), limit, "Countries should not exceed limit")

		return nil
	})
}

func TestPerformerFacets_WithTagFilter(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		pqb := db.Performer

		// Filter to performers with specific tag
		tagFilter := &models.PerformerFilterType{
			Tags: &models.HierarchicalMultiCriterionInput{
				Value:    []string{strconv.Itoa(tagIDs[tagIdxWithPerformer])},
				Modifier: models.CriterionModifierIncludes,
			},
		}

		facets, err := pqb.GetFacets(ctx, tagFilter, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// All returned facets should be from performers with this tag
		// Their count > 0 means they match the filtered set
		for _, tag := range facets.Tags {
			assert.Greater(t, tag.Count, 0, "Tag count should be positive in filtered results")
		}

		return nil
	})
}

func TestPerformerFacets_ReturnsStudios(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		pqb := db.Performer

		facets, err := pqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Studios are derived from performer appearances in scenes,
		// so may or may not have entries depending on test data
		for _, s := range facets.Studios {
			assert.Greater(t, s.Count, 0, "Studio count should be positive")
			assert.NotEmpty(t, s.ID, "Studio ID should not be empty")
		}

		return nil
	})
}

func TestPerformerFacets_ReturnsCountries(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		pqb := db.Performer

		facets, err := pqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Countries may or may not have entries depending on test data
		for _, c := range facets.Countries {
			assert.Greater(t, c.Count, 0, "Country count should be positive")
			assert.NotEmpty(t, c.ID, "Country ID should not be empty")
		}

		return nil
	})
}
