package api

import (
	"context"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/plugin"
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

func executePostHooks(ctx context.Context, id int, hookType plugin.HookTypeEnum, input interface{}, inputFields []string) {
	if err := manager.GetInstance().PluginCache.ExecutePostHooks(ctx, makeServerConnection(ctx), hookType, common.PostHookInput{
		ID:          id,
		Input:       input,
		InputFields: inputFields,
	}); err != nil {
		logger.Errorf("error executing post hooks: %s", err.Error())
	}
}
