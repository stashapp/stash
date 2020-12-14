package api

import (
	"context"

	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
)

func (r *sceneMarkerResolver) Scene(ctx context.Context, obj *models.SceneMarker) (ret *models.Scene, err error) {
	if !obj.SceneID.Valid {
		panic("Invalid scene id")
	}

	if err := r.withTxn(ctx, func(r models.Repository) error {
		sceneID := int(obj.SceneID.Int64)
		ret, err = r.Scene().Find(sceneID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *sceneMarkerResolver) PrimaryTag(ctx context.Context, obj *models.SceneMarker) (ret *models.Tag, err error) {
	if err := r.withTxn(ctx, func(r models.Repository) error {
		ret, err = r.Tag().Find(obj.PrimaryTagID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, err
}

func (r *sceneMarkerResolver) Tags(ctx context.Context, obj *models.SceneMarker) (ret []*models.Tag, err error) {
	if err := r.withTxn(ctx, func(r models.Repository) error {
		ret, err = r.Tag().FindBySceneMarkerID(obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, err
}

func (r *sceneMarkerResolver) Stream(ctx context.Context, obj *models.SceneMarker) (string, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	sceneID := int(obj.SceneID.Int64)
	return urlbuilders.NewSceneURLBuilder(baseURL, sceneID).GetSceneMarkerStreamURL(obj.ID), nil
}

func (r *sceneMarkerResolver) Preview(ctx context.Context, obj *models.SceneMarker) (string, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	sceneID := int(obj.SceneID.Int64)
	return urlbuilders.NewSceneURLBuilder(baseURL, sceneID).GetSceneMarkerStreamPreviewURL(obj.ID), nil
}
