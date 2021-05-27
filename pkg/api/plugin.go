package api

import (
	"context"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/plugin/common"
)

func makeServerConnection(ctx context.Context) common.StashServerConnection {
	cookie := makePluginCookie(ctx)

	config := manager.GetInstance().Config
	serverConnection := common.StashServerConnection{
		Scheme:        "http",
		Port:          config.GetPort(),
		SessionCookie: cookie,
		Dir:           config.GetConfigPath(),
	}

	if HasTLSConfig() {
		serverConnection.Scheme = "https"
	}

	return serverConnection
}
