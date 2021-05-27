package api

import (
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/plugin"
)

func pluginCache() *plugin.Cache {
	return manager.GetInstance().PluginCache
}
