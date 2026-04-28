package audio

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/models"
)

// QueryOptions returns a AudioQueryOptions populated with the provided filters.
func QueryOptions(audioFilter *models.AudioFilterType, findFilter *models.FindFilterType, count bool) models.AudioQueryOptions {
	return models.AudioQueryOptions{
		QueryOptions: models.QueryOptions{
			FindFilter: findFilter,
			Count:      count,
		},
		AudioFilter: audioFilter,
	}
}

// QueryWithCount queries for audios, returning the audio objects and the total count.
func QueryWithCount(ctx context.Context, qb models.AudioQueryer, audioFilter *models.AudioFilterType, findFilter *models.FindFilterType) ([]*models.Audio, int, error) {
	// this was moved from the queryBuilder code
	// left here so that calling functions can reference this instead
	result, err := qb.Query(ctx, QueryOptions(audioFilter, findFilter, true))
	if err != nil {
		return nil, 0, err
	}

	audios, err := result.Resolve(ctx)
	if err != nil {
		return nil, 0, err
	}

	return audios, result.Count, nil
}

// Query queries for audios using the provided filters.
func Query(ctx context.Context, qb models.AudioQueryer, audioFilter *models.AudioFilterType, findFilter *models.FindFilterType) ([]*models.Audio, error) {
	result, err := qb.Query(ctx, QueryOptions(audioFilter, findFilter, false))
	if err != nil {
		return nil, err
	}

	audios, err := result.Resolve(ctx)
	if err != nil {
		return nil, err
	}

	return audios, nil
}

func BatchProcess(ctx context.Context, reader models.AudioQueryer, audioFilter *models.AudioFilterType, findFilter *models.FindFilterType, fn func(audio *models.Audio) error) error {
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

		audios, err := Query(ctx, reader, audioFilter, findFilter)
		if err != nil {
			return fmt.Errorf("error querying for audios: %w", err)
		}

		for _, audio := range audios {
			if err := fn(audio); err != nil {
				return err
			}
		}

		if len(audios) != batchSize {
			more = false
		} else {
			*findFilter.Page++
		}
	}

	return nil
}

// FilterFromPaths creates a AudioFilterType that filters using the provided
// paths.
func FilterFromPaths(paths []string) *models.AudioFilterType {
	ret := &models.AudioFilterType{}
	or := ret
	sep := string(filepath.Separator)

	for _, p := range paths {
		if !strings.HasSuffix(p, sep) {
			p += sep
		}

		if ret.Path == nil {
			or = ret
		} else {
			newOr := &models.AudioFilterType{}
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

func CountByStudioID(ctx context.Context, r models.AudioQueryer, id int, depth *int) (int, error) {
	filter := &models.AudioFilterType{
		Studios: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByTagID(ctx context.Context, r models.AudioQueryer, id int, depth *int) (int, error) {
	filter := &models.AudioFilterType{
		Tags: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByGroupID(ctx context.Context, r models.AudioQueryer, id int, depth *int) (int, error) {
	filter := &models.AudioFilterType{
		Groups: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}
