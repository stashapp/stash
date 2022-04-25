package plugin

import (
	"net/http"

	"github.com/stashapp/stash/pkg/plugin/common"
)

type PluginTask struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Plugin      *Plugin `json:"plugin"`
}

// Task is the interface that handles management of a single plugin task.
type Task interface {
	// Start starts the plugin task. Returns an error if task could not be
	// started or the task has already been started.
	Start() error

	// Stop instructs a running plugin task to stop and returns immediately.
	// Use Wait to subsequently wait for the task to stop.
	Stop() error

	// Wait blocks until the plugin task is complete. Returns immediately if
	// task has not been started.
	Wait()

	// GetResult returns the output of the plugin task. Returns nil if the task
	// has not completed.
	GetResult() *common.PluginOutput
}

type taskBuilder interface {
	build(task pluginTask) Task
}

type pluginTask struct {
	plugin       *Config
	operation    *OperationConfig
	input        common.PluginInput
	gqlHandler   http.Handler
	serverConfig ServerConfig

	progress chan float64
	result   *common.PluginOutput
}

func (t *pluginTask) GetResult() *common.PluginOutput {
	return t.result
}

func (t *pluginTask) createTask() Task {
	return t.plugin.Interface.getTaskBuilder().build(*t)
}
