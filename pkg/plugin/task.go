package plugin

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common"
)

type PluginTaskManager interface {
	Start() error
	Stop() error
	Wait()
	GetProgress() float64
	GetResult() *common.PluginOutput
}

type PluginTask struct {
	Operation        *PluginOperationConfig
	ServerConnection common.StashServerConnection
	Args             []*models.PluginArgInput

	progress float64
	result   *common.PluginOutput
}

func (t *PluginTask) GetResult() *common.PluginOutput {
	return t.result
}

func (t *PluginTask) GetProgress() float64 {
	return t.progress
}

func newPluginTask(operation *PluginOperationConfig, args []*models.PluginArgInput, serverConnection common.StashServerConnection) PluginTask {
	return PluginTask{
		Operation:        operation,
		ServerConnection: serverConnection,
		Args:             args,
		progress:         -1,
	}
}
