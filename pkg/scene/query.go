package scene

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/models"
)

// QueryOptions returns a SceneQueryOptions populated with the provided filters.
func QueryOptions(sceneFilter *models.SceneFilterType, findFilter *models.FindFilterType, count bool) models.SceneQueryOptions {
	return models.SceneQueryOptions{
		QueryOptions: models.QueryOptions{
			FindFilter: findFilter,
			Count:      count,
		},
		SceneFilter: sceneFilter,
	}
}

// QueryWithCount queries for scenes, returning the scene objects and the total count.
func QueryWithCount(ctx context.Context, qb models.SceneQueryer, sceneFilter *models.SceneFilterType, findFilter *models.FindFilterType) ([]*models.Scene, int, error) {
	// this was moved from the queryBuilder code
	// left here so that calling functions can reference this instead
	result, err := qb.Query(ctx, QueryOptions(sceneFilter, findFilter, true))
	if err != nil {
		return nil, 0, err
	}

	scenes, err := result.Resolve(ctx)
	if err != nil {
		return nil, 0, err
	}

	return scenes, result.Count, nil
}

// Query queries for scenes using the provided filters.
func Query(ctx context.Context, qb models.SceneQueryer, sceneFilter *models.SceneFilterType, findFilter *models.FindFilterType) ([]*models.Scene, error) {
	result, err := qb.Query(ctx, QueryOptions(sceneFilter, findFilter, false))
	if err != nil {
		return nil, err
	}

	scenes, err := result.Resolve(ctx)
	if err != nil {
		return nil, err
	}

	return scenes, nil
}

func BatchProcess(ctx context.Context, reader models.SceneQueryer, sceneFilter *models.SceneFilterType, findFilter *models.FindFilterType, fn func(scene *models.Scene) error) error {
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

		scenes, err := Query(ctx, reader, sceneFilter, findFilter)
		if err != nil {
			return fmt.Errorf("error querying for scenes: %w", err)
		}

		for _, scene := range scenes {
			if err := fn(scene); err != nil {
				return err
			}
		}

		if len(scenes) != batchSize {
			more = false
		} else {
			*findFilter.Page++
		}
	}

	return nil
}

// FilterFromPaths creates a SceneFilterType that filters using the provided
// paths.
func FilterFromPaths(paths []string) *models.SceneFilterType {
	ret := &models.SceneFilterType{}
	or := ret
	sep := string(filepath.Separator)

	for _, p := range paths {
		if !strings.HasSuffix(p, sep) {
			p += sep
		}

		if ret.Path == nil {
			or = ret
		} else {
			newOr := &models.SceneFilterType{}
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

func CountByStudioID(ctx context.Context, r models.SceneQueryer, id int, depth *int) (int, error) {
	filter := &models.SceneFilterType{
		Studios: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByTagID(ctx context.Context, r models.SceneQueryer, id int, depth *int) (int, error) {
	filter := &models.SceneFilterType{
		Tags: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByGroupID(ctx context.Context, r models.SceneQueryer, id int, depth *int) (int, error) {
	filter := &models.SceneFilterType{
		Groups: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}
