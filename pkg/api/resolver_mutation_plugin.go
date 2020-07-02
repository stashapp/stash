package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
)

func (r *mutationResolver) RunPluginTask(ctx context.Context, pluginID string, taskName string, args []*models.PluginArgInput) (string, error) {
	//return plugin.RunPluginOperation(pluginID, operationName, args)
	return "", nil
}

func (r *mutationResolver) ReloadPlugins(ctx context.Context) (bool, error) {
	err := plugin.ReloadPlugins()
	if err != nil {
		return false, err
	}

	return true, nil
}
