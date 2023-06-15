package movie

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type Queryer interface {
	Query(ctx context.Context, movieFilter *models.MovieFilterType, findFilter *models.FindFilterType) ([]*models.Movie, int, error)
}

type CountQueryer interface {
	QueryCount(ctx context.Context, movieFilter *models.MovieFilterType, findFilter *models.FindFilterType) (int, error)
}

func CountByStudioID(ctx context.Context, r CountQueryer, id int, depth *int) (int, error) {
	filter := &models.MovieFilterType{
		Studios: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}
