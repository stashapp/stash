package api

import (
	"context"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
)

func (r *mutationResolver) RunPluginTask(ctx context.Context, pluginID string, taskName string, args []*models.PluginArgInput) (string, error) {
	manager.GetInstance().RunPluginTask(pluginID, taskName, args, makeServerConnection(ctx))
	return "todo", nil
}

func (r *mutationResolver) ReloadPlugins(ctx context.Context) (bool, error) {
	err := manager.GetInstance().PluginCache.LoadPlugins()
	if err != nil {
		logger.Errorf("Error reading plugin configs: %s", err.Error())
	}

	return true, nil
}
