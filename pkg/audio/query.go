package audio

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"

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
