package scene

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type MarkerQueryer interface {
	Query(ctx context.Context, sceneMarkerFilter *models.SceneMarkerFilterType, findFilter *models.FindFilterType) ([]*models.SceneMarker, int, error)
}

type MarkerCountQueryer interface {
	QueryCount(ctx context.Context, sceneMarkerFilter *models.SceneMarkerFilterType, findFilter *models.FindFilterType) (int, error)
}

func MarkerCountByTagID(ctx context.Context, r MarkerCountQueryer, id int, depth *int) (int, error) {
	filter := &models.SceneMarkerFilterType{
		Tags: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}
