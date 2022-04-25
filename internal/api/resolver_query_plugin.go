package api

import (
	"context"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/plugin"
)

func (r *queryResolver) Plugins(ctx context.Context) ([]*plugin.Plugin, error) {
	return manager.GetInstance().PluginCache.ListPlugins(), nil
}

func (r *queryResolver) PluginTasks(ctx context.Context) ([]*plugin.PluginTask, error) {
	return manager.GetInstance().PluginCache.ListPluginTasks(), nil
}
