package api

import (
	"context"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/sliceutil"
)

func (r *mutationResolver) RunPluginTask(ctx context.Context, pluginID string, taskName string, args []*plugin.PluginArgInput) (string, error) {
	m := manager.GetInstance()
	m.RunPluginTask(ctx, pluginID, taskName, args)
	return "todo", nil
}

func (r *mutationResolver) ReloadPlugins(ctx context.Context) (bool, error) {
	manager.GetInstance().RefreshPluginCache()
	return true, nil
}

func (r *mutationResolver) SetPluginsEnabled(ctx context.Context, enabledMap map[string]bool) (bool, error) {
	c := config.GetInstance()

	existingDisabled := c.GetDisabledPlugins()
	var newDisabled []string

	// remove plugins that are no longer disabled
	for _, disabledID := range existingDisabled {
		if enabled, found := enabledMap[disabledID]; !enabled || !found {
			newDisabled = append(newDisabled, disabledID)
		}
	}

	// add plugins that are newly disabled
	for pluginID, enabled := range enabledMap {
		if !enabled {
			newDisabled = sliceutil.AppendUnique(newDisabled, pluginID)
		}
	}

	c.Set(config.DisabledPlugins, newDisabled)

	if err := c.Write(); err != nil {
		return false, err
	}

	return true, nil
}
