package studio

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func ByName(ctx context.Context, qb models.StudioQueryer, name string) (*models.Studio, error) {
	f := &models.StudioFilterType{
		Name: &models.StringCriterionInput{
			Value:    name,
			Modifier: models.CriterionModifierEquals,
		},
	}

	pp := 1
	ret, count, err := qb.Query(ctx, f, &models.FindFilterType{
		PerPage: &pp,
	})

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return ret[0], nil
	}

	return nil, nil
}

func ByAlias(ctx context.Context, qb models.StudioQueryer, alias string) (*models.Studio, error) {
	f := &models.StudioFilterType{
		Aliases: &models.StringCriterionInput{
			Value:    alias,
			Modifier: models.CriterionModifierEquals,
		},
	}

	pp := 1
	ret, count, err := qb.Query(ctx, f, &models.FindFilterType{
		PerPage: &pp,
	})

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return ret[0], nil
	}

	return nil, nil
}

func CountByTagID(ctx context.Context, qb models.StudioQueryer, id int, depth *int) (int, error) {
	filter := &models.StudioFilterType{
		Tags: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return qb.QueryCount(ctx, filter, nil)
}
