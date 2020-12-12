package api

import (
	"context"

	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

func (r *sceneMarkerResolver) Scene(ctx context.Context, obj *models.SceneMarker) (*models.Scene, error) {
	if !obj.SceneID.Valid {
		panic("Invalid scene id")
	}
	qb := sqlite.NewSceneQueryBuilder()
	sceneID := int(obj.SceneID.Int64)
	scene, err := qb.Find(sceneID)
	return scene, err
}

func (r *sceneMarkerResolver) PrimaryTag(ctx context.Context, obj *models.SceneMarker) (*models.Tag, error) {
	qb := sqlite.NewTagQueryBuilder()
	tag, err := qb.Find(obj.PrimaryTagID, nil)
	return tag, err
}

func (r *sceneMarkerResolver) Tags(ctx context.Context, obj *models.SceneMarker) ([]*models.Tag, error) {
	qb := sqlite.NewTagQueryBuilder()
	return qb.FindBySceneMarkerID(obj.ID, nil)
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
