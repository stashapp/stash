package api

import (
	"errors"

	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindSceneMarkers(
	ctx context.Context,
	sceneMarkerFilter *models.SceneMarkerFilterType,
	savedFilterID *string,
	filter *models.FindFilterType,
	ids []string,
) (ret *FindSceneMarkersResultType, err error) {
	if sceneMarkerFilter != nil && savedFilterID != nil {
		return nil, errors.New("cannot provide both sceneMarkerFilter and saved_filter_id")
	}

	var finalFilter *models.SceneMarkerFilterType
	if savedFilterID != nil {
		finalFilter = &models.SceneMarkerFilterType{}
		var mode models.FilterMode = models.FilterModeSceneMarkers

		mergedFindFilter, err := r.resolveSavedFilter(ctx, *savedFilterID, mode, finalFilter, filter)
		if err != nil {
			return nil, err
		}
		filter = mergedFindFilter
	} else {
		finalFilter = sceneMarkerFilter
	}
	idInts, err := handleIDList(ids, "ids")
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var sceneMarkers []*models.SceneMarker
		var err error
		var total int

		if len(idInts) > 0 {
			sceneMarkers, err = r.repository.SceneMarker.FindMany(ctx, idInts)
			total = len(sceneMarkers)
		} else {
			sceneMarkers, total, err = r.repository.SceneMarker.Query(ctx, finalFilter, filter)
		}

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

func (r *queryResolver) AllSceneMarkers(ctx context.Context) (ret []*models.SceneMarker, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.SceneMarker.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
