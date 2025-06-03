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

func CountByGroupID(ctx context.Context, r models.PerformerQueryer, id int, depth *int) (int, error) {
	filter := &models.PerformerFilterType{
		Groups: &models.HierarchicalMultiCriterionInput{
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

func ByAlias(ctx context.Context, r models.PerformerQueryer, alias string) ([]*models.Performer, error) {
	f := &models.PerformerFilterType{
		Aliases: &models.StringCriterionInput{
			Value:    alias,
			Modifier: models.CriterionModifierEquals,
		},
	}

	ret, count, err := r.Query(ctx, f, nil)

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return ret, nil
	}

	return nil, nil
}
