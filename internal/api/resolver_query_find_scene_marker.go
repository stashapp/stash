package api

import (
	"context"
	"github.com/stashapp/stash/internal/models"
)

func (r *queryResolver) FindSceneMarkers(ctx context.Context, scene_marker_filter *models.SceneMarkerFilterType, filter *models.FindFilterType) (models.FindSceneMarkersResultType, error) {
	qb := models.NewSceneMarkerQueryBuilder()
	sceneMarkers, total := qb.Query(scene_marker_filter, filter)
	return models.FindSceneMarkersResultType{
		Count: total,
		SceneMarkers: sceneMarkers,
	}, nil
}
