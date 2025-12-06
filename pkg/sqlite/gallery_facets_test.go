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

func TestGalleryFacets_ReturnsTags(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		gqb := db.Gallery

		facets, err := gqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have some tags (from galleryTags in setup_test.go)
		assert.Greater(t, len(facets.Tags), 0, "Should return at least one tag")

		// Verify tags are sorted by count descending
		for i := 1; i < len(facets.Tags); i++ {
			assert.GreaterOrEqual(t, facets.Tags[i-1].Count, facets.Tags[i].Count,
				"Tags should be sorted by count descending")
		}

		return nil
	})
}

func TestGalleryFacets_ReturnsPerformers(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		gqb := db.Gallery

		facets, err := gqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have some performers (from galleryPerformers in setup_test.go)
		assert.Greater(t, len(facets.Performers), 0, "Should return at least one performer")

		// Verify performers are sorted by count descending
		for i := 1; i < len(facets.Performers); i++ {
			assert.GreaterOrEqual(t, facets.Performers[i-1].Count, facets.Performers[i].Count,
				"Performers should be sorted by count descending")
		}

		return nil
	})
}

func TestGalleryFacets_ReturnsStudios(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		gqb := db.Gallery

		facets, err := gqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have some studios (from galleryStudios in setup_test.go)
		assert.Greater(t, len(facets.Studios), 0, "Should return at least one studio")

		// Verify studios are sorted by count descending
		for i := 1; i < len(facets.Studios); i++ {
			assert.GreaterOrEqual(t, facets.Studios[i-1].Count, facets.Studios[i].Count,
				"Studios should be sorted by count descending")
		}

		return nil
	})
}

func TestGalleryFacets_ReturnsOrganized(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		gqb := db.Gallery

		facets, err := gqb.GetFacets(ctx, nil, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have organized facet with true and/or false counts
		assert.Greater(t, len(facets.Organized), 0, "Should return organized facet")

		for _, o := range facets.Organized {
			assert.GreaterOrEqual(t, o.Count, 0, "Organized count should be non-negative")
		}

		return nil
	})
}

func TestGalleryFacets_ReturnsRatings(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		gqb := db.Gallery

		facets, err := gqb.GetFacets(ctx, nil, 100)
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

func TestGalleryFacets_RespectsLimit(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		gqb := db.Gallery

		limit := 5
		facets, err := gqb.GetFacets(ctx, nil, limit)
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

func TestGalleryFacets_WithStudioFilter(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		gqb := db.Gallery

		// Filter to galleries with specific studio
		studioFilter := &models.GalleryFilterType{
			Studios: &models.HierarchicalMultiCriterionInput{
				Value:    []string{strconv.Itoa(studioIDs[studioIdxWithGallery])},
				Modifier: models.CriterionModifierIncludes,
			},
		}

		facets, err := gqb.GetFacets(ctx, studioFilter, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// All returned performers should be from galleries in this studio
		for _, p := range facets.Performers {
			assert.Greater(t, p.Count, 0, "Performer count should be positive in filtered results")
		}

		return nil
	})
}

func TestGalleryFacets_WithStudioExcludeFilter(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		gqb := db.Gallery

		// Exclude a specific studio
		studioFilter := &models.GalleryFilterType{
			Studios: &models.HierarchicalMultiCriterionInput{
				Value:    []string{},
				Excludes: []string{strconv.Itoa(studioIDs[studioIdxWithGallery])},
				Modifier: models.CriterionModifierIncludes,
			},
		}

		facets, err := gqb.GetFacets(ctx, studioFilter, 100)
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// The excluded studio should not appear in studio facets
		for _, s := range facets.Studios {
			assert.NotEqual(t, strconv.Itoa(studioIDs[studioIdxWithGallery]), s.ID,
				"Excluded studio should not appear in facets")
		}

		return nil
	})
}
