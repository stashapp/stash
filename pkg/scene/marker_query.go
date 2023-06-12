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

func MarkerCountByTagID(ctx context.Context, r MarkerCountQueryer, id int, all bool) (int, error) {
	depth := 0
	if all {
		depth = -1
	}

	filter := &models.SceneMarkerFilterType{
		Tags: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    &depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}
