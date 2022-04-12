package api

import (
	"context"
	"time"

	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
)

func (r *sceneMarkerResolver) Scene(ctx context.Context, obj *models.SceneMarker) (ret *models.Scene, err error) {
	if !obj.SceneID.Valid {
		panic("Invalid scene id")
	}

	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		sceneID := int(obj.SceneID.Int64)
		ret, err = repo.Scene().Find(sceneID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *sceneMarkerResolver) PrimaryTag(ctx context.Context, obj *models.SceneMarker) (ret *models.Tag, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.Tag().Find(obj.PrimaryTagID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, err
}

func (r *sceneMarkerResolver) Tags(ctx context.Context, obj *models.SceneMarker) (ret []*models.Tag, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.Tag().FindBySceneMarkerID(obj.ID)
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

func (r *sceneMarkerResolver) Screenshot(ctx context.Context, obj *models.SceneMarker) (string, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	sceneID := int(obj.SceneID.Int64)
	return urlbuilders.NewSceneURLBuilder(baseURL, sceneID).GetSceneMarkerStreamScreenshotURL(obj.ID), nil
}

func (r *sceneMarkerResolver) CreatedAt(ctx context.Context, obj *models.SceneMarker) (*time.Time, error) {
	return &obj.CreatedAt.Timestamp, nil
}

func (r *sceneMarkerResolver) UpdatedAt(ctx context.Context, obj *models.SceneMarker) (*time.Time, error) {
	return &obj.UpdatedAt.Timestamp, nil
}
