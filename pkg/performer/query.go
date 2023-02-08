package performer

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type Queryer interface {
	Query(ctx context.Context, performerFilter *models.PerformerFilterType, findFilter *models.FindFilterType) ([]*models.Performer, int, error)
}

type CountQueryer interface {
	QueryCount(ctx context.Context, galleryFilter *models.PerformerFilterType, findFilter *models.FindFilterType) (int, error)
}

func CountByStudioID(ctx context.Context, r CountQueryer, id int) (int, error) {
	filter := &models.PerformerFilterType{
		Studios: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}
