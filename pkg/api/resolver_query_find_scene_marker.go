package api

import (
	"context"
	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindSceneMarkers(ctx context.Context, sceneMarkerFilter *models.SceneMarkerFilterType, filter *models.FindFilterType) (*models.FindSceneMarkersResultType, error) {
	qb := models.NewSceneMarkerQueryBuilder()
	sceneMarkers, total := qb.Query(sceneMarkerFilter, filter)
	return &models.FindSceneMarkersResultType{
		Count:        total,
		SceneMarkers: sceneMarkers,
	}, nil
}
