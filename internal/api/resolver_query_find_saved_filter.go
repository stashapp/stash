package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindSavedFilter(ctx context.Context, id string) (ret *models.SavedFilter, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.SavedFilter.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}
	return ret, err
}

func (r *queryResolver) FindSavedFilters(ctx context.Context, mode *models.FilterMode) (ret []*models.SavedFilter, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		if mode != nil {
			ret, err = r.repository.SavedFilter.FindByMode(ctx, *mode)
		} else {
			ret, err = r.repository.SavedFilter.All(ctx)
		}
		return err
	}); err != nil {
		return nil, err
	}
	return ret, err
}

func (r *queryResolver) FindDefaultFilter(ctx context.Context, mode models.FilterMode) (ret *models.SavedFilter, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.SavedFilter.FindDefault(ctx, mode)
		return err
	}); err != nil {
		return nil, err
	}
	return ret, err
}
