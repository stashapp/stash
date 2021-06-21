package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindSavedFilters(ctx context.Context, mode models.FilterMode) (ret []*models.SavedFilter, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.SavedFilter().FindByMode(mode)
		return err
	}); err != nil {
		return nil, err
	}
	return ret, err
}

func (r *queryResolver) FindDefaultFilter(ctx context.Context, mode models.FilterMode) (ret *models.SavedFilter, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.SavedFilter().FindDefault(mode)
		return err
	}); err != nil {
		return nil, err
	}
	return ret, err
}
