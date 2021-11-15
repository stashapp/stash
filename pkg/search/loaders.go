package search

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

// loaders contain all data loaders the search system uses. The general
// strategy is to reset the loaders once in a while if processing large
// batch jobs to keep the cache somewhat reasonable. However, we only
// reset data which is likely to grow.
//
// As an example, scenes contains performers, but nothing contains scenes.
// They are a leaf in the topological sorting. Hence, scenes are a prime
// candidate for cleaning once in a while: one used, it isn't going to be
// reused, ever.
//
// As another example, performers are likely to be a small table in the
// database. Just caching it over the batch processing seems correct as it
// drops database query pressure by quite a lot.
type loaders struct {
	mgr       models.TransactionManager
	scene     *models.SceneLoader
	performer *models.PerformerLoader
	tag       *models.TagLoader
	studio    *models.StudioLoader
}

// newLoaders creates a new loader struct for the given transaction manager
func newLoaders(ctx context.Context, mgr models.TransactionManager) *loaders {
	scene := models.NewSceneLoader(models.NewSceneLoaderConfig(ctx, mgr))
	performer := models.NewPerformerLoader(models.NewPerformerLoaderConfig(ctx, mgr))
	tag := models.NewTagLoader(models.NewTagLoaderConfig(ctx, mgr))
	studio := models.NewStudioLoader(models.NewStudioLoaderConfig(ctx, mgr))

	return &loaders{
		mgr:       mgr,
		scene:     scene,
		performer: performer,
		tag:       tag,
		studio:    studio,
	}
}

func (l *loaders) reset(ctx context.Context) {
	// We only reset the large loaders. The volume is often on either Scene
	// and or Image. So they are the large data sets. Other tables tend to
	// be relatively small in size, so we can just keep them around for the
	// loaders run.
	l.scene = models.NewSceneLoader(models.NewSceneLoaderConfig(ctx, l.mgr))
}
