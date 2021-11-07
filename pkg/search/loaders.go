package search

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

// loaders contain all data loaders the search system uses
type loaders struct {
	scene *models.SceneLoader
}

// newLoaders creates a new loader struct for the given transaction manager
func newLoaders(ctx context.Context, mgr models.TransactionManager) loaders {
	scene := models.NewSceneLoader(models.NewSceneLoaderConfig(ctx, mgr))
	return loaders{
		scene: scene,
	}
}
