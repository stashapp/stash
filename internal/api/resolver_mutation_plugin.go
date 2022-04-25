package api

import (
	"context"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/plugin"
)

func (r *mutationResolver) RunPluginTask(ctx context.Context, pluginID string, taskName string, args []*plugin.PluginArgInput) (string, error) {
	m := manager.GetInstance()
	m.RunPluginTask(ctx, pluginID, taskName, args)
	return "todo", nil
}

func (r *mutationResolver) ReloadPlugins(ctx context.Context) (bool, error) {
	err := manager.GetInstance().PluginCache.LoadPlugins()
	if err != nil {
		logger.Errorf("Error reading plugin configs: %v", err)
	}

	return true, nil
}
