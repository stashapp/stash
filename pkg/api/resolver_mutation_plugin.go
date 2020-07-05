package api

import (
	"context"
	"net/http"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common"
)

func (r *mutationResolver) RunPluginTask(ctx context.Context, pluginID string, taskName string, args []*models.PluginArgInput) (string, error) {
	currentUser := getCurrentUserID(ctx)

	var cookie *http.Cookie
	var err error
	if currentUser != nil {
		cookie, err = createSessionCookie(*currentUser)
		if err != nil {
			return "", err
		}
	}

	serverConnection := common.StashServerConnection{
		Scheme:        "http",
		Port:          config.GetPort(),
		SessionCookie: cookie,
		Dir:           config.GetConfigPath(),
	}

	if HasTLSConfig() {
		serverConnection.Scheme = "https"
	}

	manager.GetInstance().RunPluginTask(pluginID, taskName, args, serverConnection)
	return "todo", nil
}

func (r *mutationResolver) ReloadPlugins(ctx context.Context) (bool, error) {
	err := manager.GetInstance().PluginCache.ReloadPlugins()
	if err != nil {
		logger.Errorf("Error reading plugin configs: %s", err.Error())
	}

	return true, nil
}
