package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

const defaultFacetLimit = 100

func (r *queryResolver) SceneFacets(
	ctx context.Context,
	sceneFilter *models.SceneFilterType,
	limit *int,
	includePerformerTags *bool,
	includeCaptions *bool,
) (*SceneFacetsResult, error) {
	effectiveLimit := defaultFacetLimit
	if limit != nil && *limit > 0 {
		effectiveLimit = *limit
	}

	// Build options for expensive facets (default to false for lazy loading)
	options := models.SceneFacetOptions{
		IncludePerformerTags: includePerformerTags != nil && *includePerformerTags,
		IncludeCaptions:      includeCaptions != nil && *includeCaptions,
	}

	var result *SceneFacetsResult

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		facets, err := r.repository.Scene.GetFacets(ctx, sceneFilter, effectiveLimit, options)
		if err != nil {
			return err
		}

		result = convertSceneFacets(facets)
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) PerformerFacets(
	ctx context.Context,
	performerFilter *models.PerformerFilterType,
	limit *int,
) (*PerformerFacetsResult, error) {
	effectiveLimit := defaultFacetLimit
	if limit != nil && *limit > 0 {
		effectiveLimit = *limit
	}

	var result *PerformerFacetsResult

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		facets, err := r.repository.Performer.GetFacets(ctx, performerFilter, effectiveLimit)
		if err != nil {
			return err
		}

		result = convertPerformerFacets(facets)
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) GalleryFacets(
	ctx context.Context,
	galleryFilter *models.GalleryFilterType,
	limit *int,
) (*GalleryFacetsResult, error) {
	effectiveLimit := defaultFacetLimit
	if limit != nil && *limit > 0 {
		effectiveLimit = *limit
	}

	var result *GalleryFacetsResult

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		facets, err := r.repository.Gallery.GetFacets(ctx, galleryFilter, effectiveLimit)
		if err != nil {
			return err
		}

		result = convertGalleryFacets(facets)
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) GroupFacets(
	ctx context.Context,
	groupFilter *models.GroupFilterType,
	limit *int,
) (*GroupFacetsResult, error) {
	effectiveLimit := defaultFacetLimit
	if limit != nil && *limit > 0 {
		effectiveLimit = *limit
	}

	var result *GroupFacetsResult

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		facets, err := r.repository.Group.GetFacets(ctx, groupFilter, effectiveLimit)
		if err != nil {
			return err
		}

		result = convertGroupFacets(facets)
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) StudioFacets(
	ctx context.Context,
	studioFilter *models.StudioFilterType,
	limit *int,
) (*StudioFacetsResult, error) {
	effectiveLimit := defaultFacetLimit
	if limit != nil && *limit > 0 {
		effectiveLimit = *limit
	}

	var result *StudioFacetsResult

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		facets, err := r.repository.Studio.GetFacets(ctx, studioFilter, effectiveLimit)
		if err != nil {
			return err
		}

		result = convertStudioFacets(facets)
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) TagFacets(
	ctx context.Context,
	tagFilter *models.TagFilterType,
	limit *int,
) (*TagFacetsResult, error) {
	effectiveLimit := defaultFacetLimit
	if limit != nil && *limit > 0 {
		effectiveLimit = *limit
	}

	var result *TagFacetsResult

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		facets, err := r.repository.Tag.GetFacets(ctx, tagFilter, effectiveLimit)
		if err != nil {
			return err
		}

		result = convertTagFacets(facets)
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

// Conversion functions

func convertSceneFacets(f *models.SceneFacets) *SceneFacetsResult {
	if f == nil {
		return &SceneFacetsResult{}
	}

	return &SceneFacetsResult{
		Tags:          convertFacetCounts(f.Tags),
		Performers:    convertFacetCounts(f.Performers),
		Studios:       convertFacetCounts(f.Studios),
		Groups:        convertFacetCounts(f.Groups),
		PerformerTags: convertFacetCounts(f.PerformerTags),
		Resolutions:   convertResolutionFacetCounts(f.Resolutions),
		Orientations:  convertOrientationFacetCounts(f.Orientations),
		Organized:     convertBooleanFacetCounts(f.Organized),
		Interactive:   convertBooleanFacetCounts(f.Interactive),
		Ratings:       convertRatingFacetCounts(f.Ratings),
		Captions:      convertCaptionFacetCounts(f.Captions),
	}
}

func convertPerformerFacets(f *models.PerformerFacets) *PerformerFacetsResult {
	if f == nil {
		return &PerformerFacetsResult{}
	}

	return &PerformerFacetsResult{
		Tags:        convertFacetCounts(f.Tags),
		Studios:     convertFacetCounts(f.Studios),
		Genders:     convertGenderFacetCounts(f.Genders),
		Countries:   convertFacetCounts(f.Countries),
		Circumcised: convertCircumcisedFacetCounts(f.Circumcised),
		Favorite:    convertBooleanFacetCounts(f.Favorite),
		Ratings:     convertRatingFacetCounts(f.Ratings),
	}
}

func convertGalleryFacets(f *models.GalleryFacets) *GalleryFacetsResult {
	if f == nil {
		return &GalleryFacetsResult{}
	}

	return &GalleryFacetsResult{
		Tags:       convertFacetCounts(f.Tags),
		Performers: convertFacetCounts(f.Performers),
		Studios:    convertFacetCounts(f.Studios),
		Organized:  convertBooleanFacetCounts(f.Organized),
		Ratings:    convertRatingFacetCounts(f.Ratings),
	}
}

func convertGroupFacets(f *models.GroupFacets) *GroupFacetsResult {
	if f == nil {
		return &GroupFacetsResult{}
	}

	return &GroupFacetsResult{
		Tags:       convertFacetCounts(f.Tags),
		Performers: convertFacetCounts(f.Performers),
		Studios:    convertFacetCounts(f.Studios),
	}
}

func convertStudioFacets(f *models.StudioFacets) *StudioFacetsResult {
	if f == nil {
		return &StudioFacetsResult{}
	}

	return &StudioFacetsResult{
		Tags:     convertFacetCounts(f.Tags),
		Parents:  convertFacetCounts(f.Parents),
		Favorite: convertBooleanFacetCounts(f.Favorite),
	}
}

func convertTagFacets(f *models.TagFacets) *TagFacetsResult {
	if f == nil {
		return &TagFacetsResult{}
	}

	return &TagFacetsResult{
		Parents:  convertFacetCounts(f.Parents),
		Children: convertFacetCounts(f.Children),
		Favorite: convertBooleanFacetCounts(f.Favorite),
	}
}

func convertFacetCounts(counts []models.FacetCount) []*FacetCount {
	result := make([]*FacetCount, len(counts))
	for i, c := range counts {
		result[i] = &FacetCount{
			ID:    c.ID,
			Label: c.Label,
			Count: c.Count,
		}
	}
	return result
}

func convertResolutionFacetCounts(counts []models.ResolutionFacetCount) []*ResolutionFacetCount {
	result := make([]*ResolutionFacetCount, len(counts))
	for i, c := range counts {
		result[i] = &ResolutionFacetCount{
			Resolution: c.Resolution,
			Count:      c.Count,
		}
	}
	return result
}

func convertOrientationFacetCounts(counts []models.OrientationFacetCount) []*OrientationFacetCount {
	result := make([]*OrientationFacetCount, len(counts))
	for i, c := range counts {
		result[i] = &OrientationFacetCount{
			Orientation: c.Orientation,
			Count:       c.Count,
		}
	}
	return result
}

func convertGenderFacetCounts(counts []models.GenderFacetCount) []*GenderFacetCount {
	result := make([]*GenderFacetCount, len(counts))
	for i, c := range counts {
		result[i] = &GenderFacetCount{
			Gender: c.Gender,
			Count:  c.Count,
		}
	}
	return result
}

func convertBooleanFacetCounts(counts []models.BooleanFacetCount) []*BooleanFacetCount {
	result := make([]*BooleanFacetCount, len(counts))
	for i, c := range counts {
		result[i] = &BooleanFacetCount{
			Value: c.Value,
			Count: c.Count,
		}
	}
	return result
}

func convertRatingFacetCounts(counts []models.RatingFacetCount) []*RatingFacetCount {
	result := make([]*RatingFacetCount, len(counts))
	for i, c := range counts {
		result[i] = &RatingFacetCount{
			Rating: c.Rating,
			Count:  c.Count,
		}
	}
	return result
}

func convertCircumcisedFacetCounts(counts []models.CircumcisedFacetCount) []*CircumcisedFacetCount {
	result := make([]*CircumcisedFacetCount, len(counts))
	for i, c := range counts {
		result[i] = &CircumcisedFacetCount{
			Value: c.Value,
			Count: c.Count,
		}
	}
	return result
}

func convertCaptionFacetCounts(counts []models.CaptionFacetCount) []*CaptionFacetCount {
	result := make([]*CaptionFacetCount, len(counts))
	for i, c := range counts {
		result[i] = &CaptionFacetCount{
			Language: c.Language,
			Count:    c.Count,
		}
	}
	return result
}
