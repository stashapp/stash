package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindSceneMarkers(ctx context.Context, sceneMarkerFilter *models.SceneMarkerFilterType, filter *models.FindFilterType) (ret *FindSceneMarkersResultType, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		sceneMarkers, total, err := r.repository.SceneMarker.Query(ctx, sceneMarkerFilter, filter)
		if err != nil {
			return err
		}
		ret = &FindSceneMarkersResultType{
			Count:        total,
			SceneMarkers: sceneMarkers,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
