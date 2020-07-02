package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
)

func (r *queryResolver) Plugins(ctx context.Context) ([]*models.Plugin, error) {
	return plugin.ListPlugins()
}

func (r *queryResolver) PluginTasks(ctx context.Context) ([]*models.PluginTask, error) {
	return plugin.ListPluginTasks()
}
