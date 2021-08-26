package studio

import "github.com/stashapp/stash/pkg/models"

func ByName(qb models.StudioReader, name string) (*models.Studio, error) {
	f := &models.StudioFilterType{
		Name: &models.StringCriterionInput{
			Value:    name,
			Modifier: models.CriterionModifierEquals,
		},
	}

	pp := 1
	ret, count, err := qb.Query(f, &models.FindFilterType{
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

func ByAlias(qb models.StudioReader, alias string) (*models.Studio, error) {
	f := &models.StudioFilterType{
		Aliases: &models.StringCriterionInput{
			Value:    alias,
			Modifier: models.CriterionModifierEquals,
		},
	}

	pp := 1
	ret, count, err := qb.Query(f, &models.FindFilterType{
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
