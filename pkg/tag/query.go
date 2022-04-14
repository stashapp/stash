package tag

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type Finder interface {
	Find(ctx context.Context, id int) (*models.Tag, error)
}

type Queryer interface {
	Query(ctx context.Context, tagFilter *models.TagFilterType, findFilter *models.FindFilterType) ([]*models.Tag, int, error)
}

func ByName(ctx context.Context, qb Queryer, name string) (*models.Tag, error) {
	f := &models.TagFilterType{
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

func ByAlias(ctx context.Context, qb Queryer, alias string) (*models.Tag, error) {
	f := &models.TagFilterType{
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
