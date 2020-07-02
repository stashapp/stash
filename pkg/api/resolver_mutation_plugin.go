package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
)

func (r *mutationResolver) RunPluginOperation(ctx context.Context, pluginID string, operationName string, args []*models.OperationArgInput) (*string, error) {
	// TODO - route to task if necessary
	err := plugin.RunPluginOperation(pluginID, operationName, args)
	return nil, err
}

func (r *mutationResolver) ReloadPlugins(ctx context.Context) (bool, error) {
	err := plugin.ReloadPlugins()
	if err != nil {
		return false, err
	}

	return true, nil
}
