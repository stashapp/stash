package studio

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type Finder interface {
	Find(ctx context.Context, id int) (*models.Studio, error)
}

type Queryer interface {
	Query(ctx context.Context, studioFilter *models.StudioFilterType, findFilter *models.FindFilterType) ([]*models.Studio, int, error)
}

func ByName(ctx context.Context, qb Queryer, name string) (*models.Studio, error) {
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

func ByAlias(ctx context.Context, qb Queryer, alias string) (*models.Studio, error) {
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
