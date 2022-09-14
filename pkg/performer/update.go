package performer

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type NameFinderCreator interface {
	FindByNames(ctx context.Context, names []string, nocase bool) ([]*models.Performer, error)
	Create(ctx context.Context, newPerformer models.Performer) (*models.Performer, error)
}
