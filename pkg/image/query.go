package image

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type Queryer interface {
	Query(ctx context.Context, options models.ImageQueryOptions) (*models.ImageQueryResult, error)
}

type CountQueryer interface {
	QueryCount(ctx context.Context, imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) (int, error)
}

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
func Query(ctx context.Context, qb Queryer, imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) ([]*models.Image, error) {
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

func CountByPerformerID(ctx context.Context, r CountQueryer, id int) (int, error) {
	filter := &models.ImageFilterType{
		Performers: &models.MultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByStudioID(ctx context.Context, r CountQueryer, id int) (int, error) {
	filter := &models.ImageFilterType{
		Studios: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByTagID(ctx context.Context, r CountQueryer, id int) (int, error) {
	filter := &models.ImageFilterType{
		Tags: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func FindByGalleryID(ctx context.Context, r Queryer, galleryID int, sortBy string, sortDir models.SortDirectionEnum) ([]*models.Image, error) {
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
