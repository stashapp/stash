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

func TestSceneFacets_ReturnsPerformers(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Scene

		// Get facets with no filter - should return performers sorted by count
		facets, err := sqb.GetFacets(ctx, nil, 100, models.SceneFacetOptions{})
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have some performers
		assert.Greater(t, len(facets.Performers), 0, "Should return at least one performer")

		// Verify performers are sorted by count descending
		for i := 1; i < len(facets.Performers); i++ {
			assert.GreaterOrEqual(t, facets.Performers[i-1].Count, facets.Performers[i].Count,
				"Performers should be sorted by count descending")
		}

		return nil
	})
}

func TestSceneFacets_ReturnsStudios(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Scene

		facets, err := sqb.GetFacets(ctx, nil, 100, models.SceneFacetOptions{})
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have some studios
		assert.Greater(t, len(facets.Studios), 0, "Should return at least one studio")

		// Verify studios are sorted by count descending
		for i := 1; i < len(facets.Studios); i++ {
			assert.GreaterOrEqual(t, facets.Studios[i-1].Count, facets.Studios[i].Count,
				"Studios should be sorted by count descending")
		}

		return nil
	})
}

func TestSceneFacets_ReturnsTags(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Scene

		facets, err := sqb.GetFacets(ctx, nil, 100, models.SceneFacetOptions{})
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have some tags
		assert.Greater(t, len(facets.Tags), 0, "Should return at least one tag")

		// Verify tags are sorted by count descending
		for i := 1; i < len(facets.Tags); i++ {
			assert.GreaterOrEqual(t, facets.Tags[i-1].Count, facets.Tags[i].Count,
				"Tags should be sorted by count descending")
		}

		return nil
	})
}

func TestSceneFacets_WithStudioFilter(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Scene

		// Filter to specific studio
		studioFilter := &models.SceneFilterType{
			Studios: &models.HierarchicalMultiCriterionInput{
				Value:    []string{strconv.Itoa(studioIDs[studioIdxWithScene])},
				Modifier: models.CriterionModifierIncludes,
			},
		}

		facets, err := sqb.GetFacets(ctx, studioFilter, 100, models.SceneFacetOptions{})
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// All returned performers should be from scenes in this studio
		// (their count > 0 means they have scenes in the filtered set)
		for _, p := range facets.Performers {
			assert.Greater(t, p.Count, 0, "Performer count should be positive in filtered results")
		}

		return nil
	})
}

func TestSceneFacets_WithStudioExcludeFilter(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Scene

		// Exclude a specific studio - using excludes array with INCLUDES modifier
		// This is the pattern that was causing issues
		studioFilter := &models.SceneFilterType{
			Studios: &models.HierarchicalMultiCriterionInput{
				Value:    []string{},
				Excludes: []string{strconv.Itoa(studioIDs[studioIdxWithScene])},
				Modifier: models.CriterionModifierIncludes,
			},
		}

		facets, err := sqb.GetFacets(ctx, studioFilter, 100, models.SceneFacetOptions{})
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should still return results (scenes not from the excluded studio)
		// The specific studio we excluded should not appear in studio facets
		for _, s := range facets.Studios {
			assert.NotEqual(t, strconv.Itoa(studioIDs[studioIdxWithScene]), s.ID,
				"Excluded studio should not appear in facets")
		}

		return nil
	})
}

func TestSceneFacets_RespectsLimit(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Scene

		limit := 5
		facets, err := sqb.GetFacets(ctx, nil, limit, models.SceneFacetOptions{})
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should not exceed limit
		assert.LessOrEqual(t, len(facets.Performers), limit, "Performers should not exceed limit")
		assert.LessOrEqual(t, len(facets.Tags), limit, "Tags should not exceed limit")
		assert.LessOrEqual(t, len(facets.Studios), limit, "Studios should not exceed limit")
		assert.LessOrEqual(t, len(facets.Groups), limit, "Groups should not exceed limit")

		return nil
	})
}

func TestSceneFacets_LazyLoadPerformerTags(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Scene

		// Test with performer tags excluded (lazy loading off)
		facetsWithout, err := sqb.GetFacets(ctx, nil, 100, models.SceneFacetOptions{
			IncludePerformerTags: false,
			IncludeCaptions:      false,
		})
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Performer tags should be empty when not requested
		assert.Equal(t, 0, len(facetsWithout.PerformerTags),
			"PerformerTags should be empty when not requested")

		// Test with performer tags included
		facetsWith, err := sqb.GetFacets(ctx, nil, 100, models.SceneFacetOptions{
			IncludePerformerTags: true,
			IncludeCaptions:      false,
		})
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Performer tags may or may not have results depending on test data,
		// but the query should execute without error
		_ = facetsWith

		return nil
	})
}

func TestSceneFacets_ReturnsBooleanFacets(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Scene

		facets, err := sqb.GetFacets(ctx, nil, 100, models.SceneFacetOptions{})
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have organized facet with true and false counts
		assert.Greater(t, len(facets.Organized), 0, "Should return organized facet")

		// Check that boolean values have counts
		for _, o := range facets.Organized {
			assert.GreaterOrEqual(t, o.Count, 0, "Organized count should be non-negative")
		}

		return nil
	})
}

func TestSceneFacets_ReturnsRatings(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Scene

		facets, err := sqb.GetFacets(ctx, nil, 100, models.SceneFacetOptions{})
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

func TestSceneFacets_ReturnsGroups(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Scene

		facets, err := sqb.GetFacets(ctx, nil, 100, models.SceneFacetOptions{})
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Groups may or may not have entries depending on test data
		// but the query should execute without error
		for _, g := range facets.Groups {
			assert.Greater(t, g.Count, 0, "Group count should be positive")
			assert.NotEmpty(t, g.ID, "Group ID should not be empty")
		}

		return nil
	})
}

func TestSceneFacets_ReturnsResolutions(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Scene

		facets, err := sqb.GetFacets(ctx, nil, 100, models.SceneFacetOptions{})
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Resolution facets should have valid enum values
		for _, r := range facets.Resolutions {
			assert.Greater(t, r.Count, 0, "Resolution count should be positive")
			assert.True(t, r.Resolution.IsValid(), "Resolution should be valid enum value")
		}

		return nil
	})
}

func TestSceneFacets_ReturnsOrientations(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Scene

		facets, err := sqb.GetFacets(ctx, nil, 100, models.SceneFacetOptions{})
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Orientation facets should have valid enum values
		for _, o := range facets.Orientations {
			assert.Greater(t, o.Count, 0, "Orientation count should be positive")
			assert.True(t, o.Orientation.IsValid(), "Orientation should be valid enum value")
		}

		return nil
	})
}

func TestSceneFacets_ReturnsInteractive(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Scene

		facets, err := sqb.GetFacets(ctx, nil, 100, models.SceneFacetOptions{})
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Should have interactive facet with true and/or false counts
		assert.Greater(t, len(facets.Interactive), 0, "Should return interactive facet")

		for _, i := range facets.Interactive {
			assert.GreaterOrEqual(t, i.Count, 0, "Interactive count should be non-negative")
		}

		return nil
	})
}

func TestSceneFacets_LazyLoadCaptions(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		sqb := db.Scene

		// Test with captions excluded (lazy loading off)
		facetsWithout, err := sqb.GetFacets(ctx, nil, 100, models.SceneFacetOptions{
			IncludePerformerTags: false,
			IncludeCaptions:      false,
		})
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Captions should be empty when not requested
		assert.Equal(t, 0, len(facetsWithout.Captions),
			"Captions should be empty when not requested")

		// Test with captions included
		facetsWith, err := sqb.GetFacets(ctx, nil, 100, models.SceneFacetOptions{
			IncludePerformerTags: false,
			IncludeCaptions:      true,
		})
		if err != nil {
			t.Errorf("Error getting facets: %s", err.Error())
			return nil
		}

		// Captions may or may not have results depending on test data,
		// but the query should execute without error
		_ = facetsWith

		return nil
	})
}
