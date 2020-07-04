package plugin

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common"
)

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

	// GetProgress returns the current progress of the task. Returns -1 if the
	// progress is not known, otherwise the return value will be between 0 and
	// 1, inclusively.
	GetProgress() float64

	// GetResult returns the output of the plugin task. Returns nil if the task
	// has not completed.
	GetResult() *common.PluginOutput
}

type taskBuilder interface {
	build(task pluginTask) Task
}

type pluginTask struct {
	plugin           *Config
	operation        *OperationConfig
	serverConnection common.StashServerConnection
	args             []*models.PluginArgInput

	progress float64
	result   *common.PluginOutput
}

func (t *pluginTask) GetResult() *common.PluginOutput {
	return t.result
}

func (t *pluginTask) GetProgress() float64 {
	return t.progress
}

func (t *pluginTask) createTask() Task {
	return t.plugin.Interface.getTaskBuilder().build(*t)
}
