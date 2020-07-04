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

		task := runningPluginTask{
			pluginCache:      s.PluginCache,
			status:           &s.Status,
			pluginID:         pluginID,
			taskName:         taskName,
			serverConnection: serverConnection,
			args:             args,
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

type runningPluginTask struct {
	pluginCache *plugin.PluginCache
	status      *TaskStatus

	pluginID         string
	taskName         string
	serverConnection common.StashServerConnection
	args             []*models.PluginArgInput

	stopping bool
}

func (t *runningPluginTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	opManager, err := t.pluginCache.RunPluginOperation(t.pluginID, t.taskName, t.serverConnection, t.args)
	if err != nil {
		logger.Errorf("Error running plugin task: %s", err.Error())
		return
	}

	done := make(chan bool)
	go func() {
		opManager.Wait()
		close(done)
	}()

	// check for stop/progress every five seconds
	pollingTime := time.Second * 5
	for {
		select {
		case <-done:
			return
		case <-time.After(pollingTime):
			t.status.setProgressPercent(opManager.GetProgress())

			if t.stopping {
				if err := opManager.Stop(); err != nil {
					logger.Errorf("Error stopping plugin operation: %s", err.Error())
				}
				return
			}
		}
	}
}
