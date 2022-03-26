package api

import (
	"context"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) Plugins(ctx context.Context) ([]*models.Plugin, error) {
	return manager.GetInstance().PluginCache.ListPlugins(), nil
}

func (r *queryResolver) PluginTasks(ctx context.Context) ([]*models.PluginTask, error) {
	return manager.GetInstance().PluginCache.ListPluginTasks(), nil
}
