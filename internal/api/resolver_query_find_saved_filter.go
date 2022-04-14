package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindSavedFilters(ctx context.Context, mode models.FilterMode) (ret []*models.SavedFilter, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.SavedFilter.FindByMode(ctx, mode)
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
