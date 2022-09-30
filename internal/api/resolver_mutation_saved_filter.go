package api

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

func (r *mutationResolver) SaveFilter(ctx context.Context, input SaveFilterInput) (ret *models.SavedFilter, err error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("name must be non-empty")
	}

	var id *int
	if input.ID != nil {
		idv, err := strconv.Atoi(*input.ID)
		if err != nil {
			return nil, err
		}
		id = &idv
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		f := models.SavedFilter{
			Mode:   input.Mode,
			Name:   input.Name,
			Filter: input.Filter,
		}
		if id == nil {
			ret, err = r.repository.SavedFilter.Create(ctx, f)
		} else {
			f.ID = *id
			ret, err = r.repository.SavedFilter.Update(ctx, f)
		}
		return err
	}); err != nil {
		return nil, err
	}
	return ret, err
}

func (r *mutationResolver) DestroySavedFilter(ctx context.Context, input DestroyFilterInput) (bool, error) {
	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.repository.SavedFilter.Destroy(ctx, id)
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SetDefaultFilter(ctx context.Context, input SetDefaultFilterInput) (bool, error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.SavedFilter

		if input.Filter == nil {
			// clearing
			def, err := qb.FindDefault(ctx, input.Mode)
			if err != nil {
				return err
			}

			if def != nil {
				return qb.Destroy(ctx, def.ID)
			}

			return nil
		}

		_, err := qb.SetDefault(ctx, models.SavedFilter{
			Mode:   input.Mode,
			Filter: *input.Filter,
		})

		return err
	}); err != nil {
		return false, err
	}

	return true, nil
}
