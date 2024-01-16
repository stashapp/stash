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

func CountByPerformerIDStudioID(ctx context.Context, r models.MovieQueryer, performerID int, studioID int, depth *int) (int, error) {
	filter := &models.MovieFilterType{
		Studios: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(studioID)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
		Performers: &models.MultiCriterionInput{
			Value:    []string{strconv.Itoa(performerID)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}
