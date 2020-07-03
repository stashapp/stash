package manager

import (
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/plugin/common"
)

func (s *singleton) RunPluginTask(pluginID string, taskName string, args []*models.PluginArgInput, serverConnection common.StashServerConnection) {
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
			PluginID:         pluginID,
			TaskName:         taskName,
			ServerConnection: serverConnection,
			Args:             args,
		}

		done := make(chan bool)
		go func() {
			defer close(done)
			go task.Start(&wg)
			wg.Wait()
		}()

		// check for stop every five seconds
		pollingTime := time.Second * 5
		for {
			select {
			case <-done:
				return
			case <-time.After(pollingTime):
				task.stopping = s.Status.stopping
			}
		}
	}()
}

type RunningPluginTask struct {
	PluginID         string
	TaskName         string
	ServerConnection common.StashServerConnection
	Args             []*models.PluginArgInput

	stopping bool
}

func (t *RunningPluginTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	// TODO - handle progress/stop
	opManager, err := plugin.RunPluginOperation(t.PluginID, t.TaskName, t.ServerConnection, t.Args)
	if err != nil {
		logger.Errorf("Error running plugin task: %s", err.Error())
	}

	done := make(chan bool)
	go func() {
		opManager.Wait()
		close(done)
	}()

	// check for stop every five seconds
	pollingTime := time.Second * 5
	for {
		select {
		case <-done:
			return
		case <-time.After(pollingTime):
			if t.stopping {
				if err := opManager.Stop(); err != nil {
					logger.Errorf("Error stopping plugin operation: %s", err.Error())
				}
				return
			}
		}
	}
}
