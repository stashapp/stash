package api

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

func (r *mutationResolver) SavePinnedFilter(ctx context.Context, input SavePinnedFilterInput) (ret *models.PinnedFilter, err error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("name must be non-empty")
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		f := models.PinnedFilter{
			Mode: input.Mode,
			Name: input.Name,
		}
		ret, err = r.repository.PinnedFilter.Create(ctx, f)
		return err
	}); err != nil {
		return nil, err
	}
	return ret, err
}

func (r *mutationResolver) DestroyPinnedFilter(ctx context.Context, input DestroyPinnedFilterInput) (bool, error) {
	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.repository.PinnedFilter.Destroy(ctx, id)
	}); err != nil {
		return false, err
	}

	return true, nil
}
