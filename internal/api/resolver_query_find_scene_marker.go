package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindSceneMarkers(ctx context.Context, sceneMarkerFilter *models.SceneMarkerFilterType, filter *models.FindFilterType) (ret *FindSceneMarkersResultType, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		sceneMarkers, total, err := repo.SceneMarker().Query(sceneMarkerFilter, filter)
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
