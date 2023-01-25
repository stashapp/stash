package performer

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type NameFinderCreator interface {
	FindByNames(ctx context.Context, names []string, nocase bool) ([]*models.Performer, error)
	Query(ctx context.Context, performerFilter *models.PerformerFilterType, findFilter *models.FindFilterType) ([]*models.Performer, int, error)
	Create(ctx context.Context, newPerformer *models.Performer) error
}
