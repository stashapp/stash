package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
)

func (r *mutationResolver) RunPluginOperation(ctx context.Context, pluginID string, operationName string, args []*models.OperationArgInput) (*models.OperationResult, error) {
	return plugin.RunPluginOperation(pluginID, operationName, args)
}

func (r *mutationResolver) RunPluginOperationJob(ctx context.Context, pluginID string, operationName string, args []*models.OperationArgInput) (string, error) {
	// TODO
	return "", nil
}

func (r *mutationResolver) ReloadPlugins(ctx context.Context) (bool, error) {
	err := plugin.ReloadPlugins()
	if err != nil {
		return false, err
	}

	return true, nil
}
