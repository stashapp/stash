package api

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestConvertFacetCounts(t *testing.T) {
	input := []models.FacetCount{
		{ID: "1", Label: "Tag A", Count: 100},
		{ID: "2", Label: "Tag B", Count: 50},
		{ID: "3", Label: "Tag C", Count: 25},
	}

	result := convertFacetCounts(input)

	assert.Len(t, result, 3)
	assert.Equal(t, "1", result[0].ID)
	assert.Equal(t, "Tag A", result[0].Label)
	assert.Equal(t, 100, result[0].Count)
	assert.Equal(t, "2", result[1].ID)
	assert.Equal(t, "Tag B", result[1].Label)
	assert.Equal(t, 50, result[1].Count)
}

func TestConvertFacetCounts_EmptyInput(t *testing.T) {
	result := convertFacetCounts([]models.FacetCount{})
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestConvertBooleanFacetCounts(t *testing.T) {
	input := []models.BooleanFacetCount{
		{Value: true, Count: 150},
		{Value: false, Count: 350},
	}

	result := convertBooleanFacetCounts(input)

	assert.Len(t, result, 2)
	assert.Equal(t, true, result[0].Value)
	assert.Equal(t, 150, result[0].Count)
	assert.Equal(t, false, result[1].Value)
	assert.Equal(t, 350, result[1].Count)
}

func TestConvertRatingFacetCounts(t *testing.T) {
	input := []models.RatingFacetCount{
		{Rating: 100, Count: 50},
		{Rating: 80, Count: 120},
		{Rating: 60, Count: 200},
	}

	result := convertRatingFacetCounts(input)

	assert.Len(t, result, 3)
	assert.Equal(t, 100, result[0].Rating)
	assert.Equal(t, 50, result[0].Count)
	assert.Equal(t, 80, result[1].Rating)
	assert.Equal(t, 120, result[1].Count)
}

func TestConvertResolutionFacetCounts(t *testing.T) {
	input := []models.ResolutionFacetCount{
		{Resolution: models.ResolutionEnumFourK, Count: 100},
		{Resolution: models.ResolutionEnumFullHd, Count: 200},
		{Resolution: models.ResolutionEnumStandardHd, Count: 300},
	}

	result := convertResolutionFacetCounts(input)

	assert.Len(t, result, 3)
	assert.Equal(t, models.ResolutionEnumFourK, result[0].Resolution)
	assert.Equal(t, 100, result[0].Count)
}

func TestConvertOrientationFacetCounts(t *testing.T) {
	input := []models.OrientationFacetCount{
		{Orientation: models.OrientationLandscape, Count: 500},
		{Orientation: models.OrientationPortrait, Count: 100},
		{Orientation: models.OrientationSquare, Count: 50},
	}

	result := convertOrientationFacetCounts(input)

	assert.Len(t, result, 3)
	assert.Equal(t, models.OrientationLandscape, result[0].Orientation)
	assert.Equal(t, 500, result[0].Count)
}

func TestConvertGenderFacetCounts(t *testing.T) {
	input := []models.GenderFacetCount{
		{Gender: models.GenderEnumFemale, Count: 300},
		{Gender: models.GenderEnumMale, Count: 200},
	}

	result := convertGenderFacetCounts(input)

	assert.Len(t, result, 2)
	assert.Equal(t, models.GenderEnumFemale, result[0].Gender)
	assert.Equal(t, 300, result[0].Count)
}

func TestConvertCircumcisedFacetCounts(t *testing.T) {
	input := []models.CircumcisedFacetCount{
		{Value: models.CircumisedEnumCut, Count: 150},
		{Value: models.CircumisedEnumUncut, Count: 100},
	}

	result := convertCircumcisedFacetCounts(input)

	assert.Len(t, result, 2)
	assert.Equal(t, models.CircumisedEnumCut, result[0].Value)
	assert.Equal(t, 150, result[0].Count)
}

func TestConvertCaptionFacetCounts(t *testing.T) {
	input := []models.CaptionFacetCount{
		{Language: "en", Count: 500},
		{Language: "de", Count: 100},
		{Language: "fr", Count: 50},
	}

	result := convertCaptionFacetCounts(input)

	assert.Len(t, result, 3)
	assert.Equal(t, "en", result[0].Language)
	assert.Equal(t, 500, result[0].Count)
}

func TestConvertSceneFacets_Nil(t *testing.T) {
	result := convertSceneFacets(nil)
	assert.NotNil(t, result)
}

func TestConvertPerformerFacets_Nil(t *testing.T) {
	result := convertPerformerFacets(nil)
	assert.NotNil(t, result)
}

func TestConvertGalleryFacets_Nil(t *testing.T) {
	result := convertGalleryFacets(nil)
	assert.NotNil(t, result)
}

func TestConvertGroupFacets_Nil(t *testing.T) {
	result := convertGroupFacets(nil)
	assert.NotNil(t, result)
}

func TestConvertStudioFacets_Nil(t *testing.T) {
	result := convertStudioFacets(nil)
	assert.NotNil(t, result)
}

func TestConvertTagFacets_Nil(t *testing.T) {
	result := convertTagFacets(nil)
	assert.NotNil(t, result)
}

// Regression test for bug where labels were discarded
func TestConvertFacetCounts_PreservesLabels(t *testing.T) {
	input := []models.FacetCount{
		{ID: "123", Label: "My Studio Name", Count: 100},
		{ID: "456", Label: "Another Studio", Count: 50},
	}

	result := convertFacetCounts(input)

	// Verify labels are preserved, not replaced with IDs
	assert.Equal(t, "My Studio Name", result[0].Label)
	assert.Equal(t, "Another Studio", result[1].Label)
	assert.Equal(t, "123", result[0].ID)
	assert.Equal(t, "456", result[1].ID)
}

