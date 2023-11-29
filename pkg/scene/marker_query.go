package scene

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func MarkerCountByTagID(ctx context.Context, r models.SceneMarkerQueryer, id int, depth *int) (int, error) {
	filter := &models.SceneMarkerFilterType{
		Tags: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}
