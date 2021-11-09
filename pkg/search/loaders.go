package search

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

// loaders contain all data loaders the search system uses
type loaders struct {
	mgr       models.TransactionManager
	scene     *models.SceneLoader
	performer *models.PerformerLoader
}

// newLoaders creates a new loader struct for the given transaction manager
func newLoaders(ctx context.Context, mgr models.TransactionManager) *loaders {
	scene := models.NewSceneLoader(models.NewSceneLoaderConfig(ctx, mgr))
	performer := models.NewPerformerLoader(models.NewPerformerLoaderConfig(ctx, mgr))

	return &loaders{
		mgr:       mgr,
		scene:     scene,
		performer: performer,
	}
}

func (l *loaders) reset(ctx context.Context) {
	// We only reset the large loaders. The volume is often on either Scene
	// and or Image. So they are the large data sets. Other tables tend to
	// be relatively small in size, so we can just keep them around for the
	// loaders run.
	l.scene = models.NewSceneLoader(models.NewSceneLoaderConfig(ctx, l.mgr))
}
