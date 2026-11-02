package tag

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func ByName(ctx context.Context, qb models.TagNameFinder, name string) (*models.Tag, error) {
	const nocase = true
	ret, err := qb.FindByName(ctx, name, nocase)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func ByAlias(ctx context.Context, qb models.TagNameFinder, alias string) (*models.Tag, error) {
	const nocase = true
	ret, err := qb.FindByAlias(ctx, alias, nocase)

	if err != nil {
		return nil, err
	}

	return ret, nil
}
