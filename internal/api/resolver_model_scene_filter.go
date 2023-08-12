package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *sceneFilterResolver) Scene(ctx context.Context, obj *models.SceneFilter) (ret *models.Scene, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Scene.Find(ctx, obj.SceneID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
