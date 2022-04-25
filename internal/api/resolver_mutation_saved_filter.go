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

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		f := models.SavedFilter{
			Mode:   input.Mode,
			Name:   input.Name,
			Filter: input.Filter,
		}
		if id == nil {
			ret, err = repo.SavedFilter().Create(f)
		} else {
			f.ID = *id
			ret, err = repo.SavedFilter().Update(f)
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

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		return repo.SavedFilter().Destroy(id)
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SetDefaultFilter(ctx context.Context, input SetDefaultFilterInput) (bool, error) {
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.SavedFilter()

		if input.Filter == nil {
			// clearing
			def, err := qb.FindDefault(input.Mode)
			if err != nil {
				return err
			}

			if def != nil {
				return qb.Destroy(def.ID)
			}

			return nil
		}

		_, err := qb.SetDefault(models.SavedFilter{
			Mode:   input.Mode,
			Filter: *input.Filter,
		})

		return err
	}); err != nil {
		return false, err
	}

	return true, nil
}
