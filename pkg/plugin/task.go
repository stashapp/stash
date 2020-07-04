package plugin

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common"
)

// PluginTaskManager is the interface that handles management of a single 
// plugin task.
type PluginTaskManager interface {
	// Start starts the plugin task. Returns an error if task could not be 
	// started.
	Start() error

	// Stop instructs a running plugin task to stop and returns immediately.
	// Use Wait to subsequently wait for the task to stop.
	Stop() error

	// Wait blocks until the plugin task is complete. Returns immediately if
	// task has not been started.
	Wait()

	// GetProgress returns the current progress of the task. Returns -1 if the
	// progress is not known, otherwise the return value will be between 0 and
	// 1, inclusively.
	GetProgress() float64

	// GetResult returns the output of the plugin task. Returns nil if the task
	// has not completed.
	GetResult() *common.PluginOutput
}

type pluginTask struct {
	Operation        *OperationConfig
	ServerConnection common.StashServerConnection
	Args             []*models.PluginArgInput

	progress float64
	result   *common.PluginOutput
}

func (t *pluginTask) GetResult() *common.PluginOutput {
	return t.result
}

func (t *pluginTask) GetProgress() float64 {
	return t.progress
}

func newPluginTask(operation *OperationConfig, args []*models.PluginArgInput, serverConnection common.StashServerConnection) pluginTask {
	return pluginTask{
		Operation:        operation,
		ServerConnection: serverConnection,
		Args:             args,
		progress:         -1,
	}
}
