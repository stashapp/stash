package image

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

// QueryOptions returns a ImageQueryResult populated with the provided filters.
func QueryOptions(imageFilter *models.ImageFilterType, findFilter *models.FindFilterType, count bool) models.ImageQueryOptions {
	return models.ImageQueryOptions{
		QueryOptions: models.QueryOptions{
			FindFilter: findFilter,
			Count:      count,
		},
		ImageFilter: imageFilter,
	}
}

// Query queries for images using the provided filters.
func Query(ctx context.Context, qb models.ImageQueryer, imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) ([]*models.Image, error) {
	result, err := qb.Query(ctx, QueryOptions(imageFilter, findFilter, false))
	if err != nil {
		return nil, err
	}

	images, err := result.Resolve(ctx)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func CountByPerformerID(ctx context.Context, r models.ImageQueryer, id int) (int, error) {
	filter := &models.ImageFilterType{
		Performers: &models.MultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByStudioID(ctx context.Context, r models.ImageQueryer, id int, depth *int) (int, error) {
	filter := &models.ImageFilterType{
		Studios: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByTagID(ctx context.Context, r models.ImageQueryer, id int, depth *int) (int, error) {
	filter := &models.ImageFilterType{
		Tags: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func FindByGalleryID(ctx context.Context, r models.ImageQueryer, galleryID int, sortBy string, sortDir models.SortDirectionEnum) ([]*models.Image, error) {
	perPage := -1

	findFilter := models.FindFilterType{
		PerPage: &perPage,
	}

	if sortBy != "" {
		findFilter.Sort = &sortBy
	}

	if sortDir.IsValid() {
		findFilter.Direction = &sortDir
	}

	return Query(ctx, r, &models.ImageFilterType{
		Galleries: &models.MultiCriterionInput{
			Value:    []string{strconv.Itoa(galleryID)},
			Modifier: models.CriterionModifierIncludes,
		},
	}, &findFilter)
}

func FindGalleryCover(ctx context.Context, r models.ImageQueryer, galleryID int, galleryCoverRegex string) (*models.Image, error) {
	const useCoverJpg = true
	img, err := findGalleryCover(ctx, r, galleryID, useCoverJpg, galleryCoverRegex)
	if err != nil {
		return nil, err
	}

	if img != nil {
		return img, nil
	}

	// return the first image in the gallery
	return findGalleryCover(ctx, r, galleryID, !useCoverJpg, galleryCoverRegex)
}

func findGalleryCover(ctx context.Context, r models.ImageQueryer, galleryID int, useCoverJpg bool, galleryCoverRegex string) (*models.Image, error) {
	// try to find cover.jpg in the gallery
	perPage := 1
	sortBy := "path"
	sortDir := models.SortDirectionEnumAsc

	findFilter := models.FindFilterType{
		PerPage:   &perPage,
		Sort:      &sortBy,
		Direction: &sortDir,
	}

	imageFilter := &models.ImageFilterType{
		Galleries: &models.MultiCriterionInput{
			Value:    []string{strconv.Itoa(galleryID)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	if useCoverJpg {
		imageFilter.Path = &models.StringCriterionInput{
			Value:    "(?i)" + galleryCoverRegex,
			Modifier: models.CriterionModifierMatchesRegex,
		}
	}

	imgs, err := Query(ctx, r, imageFilter, &findFilter)
	if err != nil {
		return nil, err
	}

	if len(imgs) > 0 {
		return imgs[0], nil
	}

	return nil, nil
}
