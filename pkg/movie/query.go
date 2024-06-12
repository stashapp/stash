package movie

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func CountByStudioID(ctx context.Context, r models.MovieQueryer, id int, depth *int) (int, error) {
	filter := &models.MovieFilterType{
		Studios: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByTagID(ctx context.Context, r models.MovieQueryer, id int, depth *int) (int, error) {
	filter := &models.MovieFilterType{
		Tags: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}
