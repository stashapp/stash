package image

import (
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type Queryer interface {
	Query(options models.ImageQueryOptions) (*models.ImageQueryResult, error)
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

// QueryWithCount queries for images, returning the image objects and the total count.
func QueryWithCount(qb Queryer, imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) ([]*models.Image, int, error) {
	// this was moved from the queryBuilder code
	// left here so that calling functions can reference this instead
	result, err := qb.Query(QueryOptions(imageFilter, findFilter, true))
	if err != nil {
		return nil, 0, err
	}

	images, err := result.Resolve()
	if err != nil {
		return nil, 0, err
	}

	return images, result.Count, nil
}

// Query queries for images using the provided filters.
func Query(qb Queryer, imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) ([]*models.Image, error) {
	result, err := qb.Query(QueryOptions(imageFilter, findFilter, false))
	if err != nil {
		return nil, err
	}

	images, err := result.Resolve()
	if err != nil {
		return nil, err
	}

	return images, nil
}

func CountByPerformerID(r models.ImageReader, id int) (int, error) {
	filter := &models.ImageFilterType{
		Performers: &models.MultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	return r.QueryCount(filter, nil)
}

func CountByStudioID(r models.ImageReader, id int) (int, error) {
	filter := &models.ImageFilterType{
		Studios: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	return r.QueryCount(filter, nil)
}

func CountByTagID(r models.ImageReader, id int) (int, error) {
	filter := &models.ImageFilterType{
		Tags: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	return r.QueryCount(filter, nil)
}
