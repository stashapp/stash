package api

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindPinnedFilter(ctx context.Context, id string) (ret *models.PinnedFilter, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.PinnedFilter.Find(ctx, idInt)
		return err
	}); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return ret, err
}

func (r *queryResolver) FindPinnedFilters(ctx context.Context, mode *models.FilterMode) (ret []*models.PinnedFilter, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		if mode != nil {
			ret, err = r.repository.PinnedFilter.FindByMode(ctx, *mode)
		} else {
			ret, err = r.repository.PinnedFilter.All(ctx)
		}
		return err
	}); err != nil {
		return nil, err
	}
	return ret, err
}
