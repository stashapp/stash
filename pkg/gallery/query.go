package gallery

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type Queryer interface {
	Query(ctx context.Context, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) ([]*models.Gallery, int, error)
}

type CountQueryer interface {
	QueryCount(ctx context.Context, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) (int, error)
}

type ChecksumsFinder interface {
	FindByChecksums(ctx context.Context, checksums []string) ([]*models.Gallery, error)
}

func CountByPerformerID(ctx context.Context, r CountQueryer, id int) (int, error) {
	filter := &models.GalleryFilterType{
		Performers: &models.MultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByStudioID(ctx context.Context, r CountQueryer, id int) (int, error) {
	filter := &models.GalleryFilterType{
		Studios: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByTagID(ctx context.Context, r CountQueryer, id int) (int, error) {
	filter := &models.GalleryFilterType{
		Tags: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}
