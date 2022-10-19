package movie

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type NameFinderCreator interface {
	FindByName(ctx context.Context, name string, nocase bool) (*models.Movie, error)
	Create(ctx context.Context, newMovie models.Movie) (*models.Movie, error)
}
