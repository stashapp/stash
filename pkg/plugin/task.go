package plugin

import (
	"sync"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common"
)

type PluginTask struct {
	Operation        *PluginOperationConfig
	ServerConnection common.StashServerConnection
	Args             []*models.PluginArgInput

	WaitGroup sync.WaitGroup

	result common.PluginOutput
}

func (t *PluginTask) GetResult() common.PluginOutput {
	return t.result
}

func newPluginTask(operation *PluginOperationConfig, args []*models.PluginArgInput, serverConnection common.StashServerConnection) PluginTask {
	return PluginTask{
		Operation:        operation,
		ServerConnection: serverConnection,
		Args:             args,
	}
}

type PluginTaskManager interface {
	Start() error
	Stop() error
	Wait()
	GetResult() common.PluginOutput
}
