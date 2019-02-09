package api

import (
	"context"
	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/internal/models"
	"strconv"
)

func (r *sceneMarkerResolver) ID(ctx context.Context, obj *models.SceneMarker) (string, error) {
	return strconv.Itoa(obj.ID), nil
}

func (r *sceneMarkerResolver) Scene(ctx context.Context, obj *models.SceneMarker) (models.Scene, error) {
	if !obj.SceneID.Valid {
		panic("Invalid scene id")
	}
	qb := models.NewSceneQueryBuilder()
	sceneID := int(obj.SceneID.Int64)
	scene, err := qb.Find(sceneID)
	return *scene, err
}

func (r *sceneMarkerResolver) PrimaryTag(ctx context.Context, obj *models.SceneMarker) (models.Tag, error) {
	qb := models.NewTagQueryBuilder()
	if !obj.PrimaryTagID.Valid {
		panic("TODO no primary tag id")
	}
	tag, err := qb.Find(int(obj.PrimaryTagID.Int64), nil) // TODO make primary tag id not null in DB
	return *tag, err
}

func (r *sceneMarkerResolver) Tags(ctx context.Context, obj *models.SceneMarker) ([]models.Tag, error) {
	qb := models.NewTagQueryBuilder()
	return qb.FindBySceneMarkerID(obj.ID, nil)
}

func (r *sceneMarkerResolver) Stream(ctx context.Context, obj *models.SceneMarker) (string, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	sceneID := int(obj.SceneID.Int64)
	return urlbuilders.NewSceneURLBuilder(baseURL, sceneID).GetSceneMarkerStreamUrl(obj.ID), nil
}

func (r *sceneMarkerResolver) Preview(ctx context.Context, obj *models.SceneMarker) (string, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	sceneID := int(obj.SceneID.Int64)
	return urlbuilders.NewSceneURLBuilder(baseURL, sceneID).GetSceneMarkerStreamPreviewUrl(obj.ID), nil
}