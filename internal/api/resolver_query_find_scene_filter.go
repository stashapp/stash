package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindSceneFilters(ctx context.Context, sceneFilterFilter *models.SceneFilterFilterType, filter *models.FindFilterType) (ret *FindSceneFiltersResultType, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		sceneFilters, total, err := r.repository.SceneFilter.Query(ctx, sceneFilterFilter, filter)
		if err != nil {
			return err
		}
		ret = &FindSceneFiltersResultType{
			Count:        total,
			SceneFilters: sceneFilters,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllSceneFilters(ctx context.Context) (ret []*models.SceneFilter, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.SceneFilter.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
