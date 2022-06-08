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

	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.SavedFilter().Find(idInt)
		return err
	}); err != nil {
		return nil, err
	}
	return ret, err
}

func (r *queryResolver) FindSavedFilters(ctx context.Context, mode *models.FilterMode) (ret []*models.SavedFilter, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		if mode != nil {
			ret, err = repo.SavedFilter().FindByMode(*mode)
		} else {
			ret, err = repo.SavedFilter().All()
		}
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
