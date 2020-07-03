package manager

import (
	"sync"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
)

func (s *singleton) RunPluginTask(pluginID string, taskName string, args []*models.PluginArgInput) {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(PluginOperation)
	s.Status.indefiniteProgress()

	go func() {
		defer s.returnToIdleState()

		var wg sync.WaitGroup
		wg.Add(1)

		task := RunningPluginTask{
			PluginID: pluginID,
			TaskName: taskName,
			Args:     args,
		}
		go task.Start(&wg)
		wg.Wait()
	}()
}

type RunningPluginTask struct {
	PluginID string
	TaskName string
	Args     []*models.PluginArgInput
}

func (t *RunningPluginTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	// TODO - handle progress/stop

	err := plugin.RunPluginOperation(t.PluginID, t.TaskName, t.Args)
	if err != nil {
		logger.Errorf("Error running plugin task: %s", err.Error())
	}
}
