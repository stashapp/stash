package gallery

import (
	"context"
	"fmt"
	"github.com/stashapp/stash/pkg/job"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

// QueryOptions returns a GalleryQueryOptions populated with the provided filters.
func QueryOptions(galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType, count bool) models.GalleryQueryOptions {
	return models.GalleryQueryOptions{
		QueryOptions: models.QueryOptions{
			FindFilter: findFilter,
			Count:      count,
		},
		GalleryFilter: galleryFilter,
	}
}

func CountByPerformerID(ctx context.Context, r models.GalleryQueryer, id int) (int, error) {
	filter := &models.GalleryFilterType{
		Performers: &models.MultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByStudioID(ctx context.Context, r models.GalleryQueryer, id int, depth *int) (int, error) {
	filter := &models.GalleryFilterType{
		Studios: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByTagID(ctx context.Context, r models.GalleryQueryer, id int, depth *int) (int, error) {
	filter := &models.GalleryFilterType{
		Tags: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

// Query queries for galleries using the provided filters.
func Query(ctx context.Context, qb models.GalleryQueryer, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) ([]*models.Gallery, error) {
	result, err := qb.Query(ctx, QueryOptions(galleryFilter, findFilter, false))
	if err != nil {
		return nil, err
	}

	galleries, err := result.Resolve(ctx)
	if err != nil {
		return nil, err
	}

	return galleries, nil
}

func BatchProcess(ctx context.Context, reader models.GalleryQueryer, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType, fn func(gallery *models.Gallery) error) error {
	const batchSize = 1000

	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	page := 1
	perPage := batchSize
	findFilter.Page = &page
	findFilter.PerPage = &perPage

	for more := true; more; {
		if job.IsCancelled(ctx) {
			return nil
		}

		galleries, err := Query(ctx, reader, galleryFilter, findFilter)
		if err != nil {
			return fmt.Errorf("error querying for galleries: %w", err)
		}

		for _, gallery := range galleries {
			if err := fn(gallery); err != nil {
				return err
			}
		}

		if len(galleries) != batchSize {
			more = false
		} else {
			*findFilter.Page++
		}
	}

	return nil
}

// FilterFromPaths creates a GalleryFilterType that filters using the provided
// paths.
func FilterFromPaths(paths []string) *models.GalleryFilterType {
	ret := &models.GalleryFilterType{}
	or := ret
	sep := string(filepath.Separator)

	for _, p := range paths {
		if !strings.HasSuffix(p, sep) {
			p += sep
		}

		if ret.Path == nil {
			or = ret
		} else {
			newOr := &models.GalleryFilterType{}
			or.Or = newOr
			or = newOr
		}

		or.Path = &models.StringCriterionInput{
			Modifier: models.CriterionModifierEquals,
			Value:    p + "%",
		}
	}

	return ret
}
