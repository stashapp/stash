package performer

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func CountByStudioID(ctx context.Context, r models.PerformerQueryer, id int, depth *int) (int, error) {
	filter := &models.PerformerFilterType{
		Studios: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByTagID(ctx context.Context, r models.PerformerQueryer, id int, depth *int) (int, error) {
	filter := &models.PerformerFilterType{
		Tags: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByAppearsWith(ctx context.Context, r models.PerformerQueryer, id int) (int, error) {
	filter := &models.PerformerFilterType{
		Performers: &models.MultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}
