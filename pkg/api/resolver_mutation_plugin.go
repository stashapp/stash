package api

import (
	"context"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
)

func (r *mutationResolver) RunPluginTask(ctx context.Context, pluginID string, taskName string, args []*models.PluginArgInput) (string, error) {
	manager.GetInstance().RunPluginTask(pluginID, taskName, args)
	return "todo", nil
}

func (r *mutationResolver) ReloadPlugins(ctx context.Context) (bool, error) {
	err := plugin.ReloadPlugins()
	if err != nil {
		return false, err
	}

	return true, nil
}
