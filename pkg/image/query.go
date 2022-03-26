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

func FindByGalleryID(r models.ImageReader, galleryID int, sortBy string, sortDir models.SortDirectionEnum) ([]*models.Image, error) {
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

	return Query(r, &models.ImageFilterType{
		Galleries: &models.MultiCriterionInput{
			Value:    []string{strconv.Itoa(galleryID)},
			Modifier: models.CriterionModifierIncludes,
		},
	}, &findFilter)
}
