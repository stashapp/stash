package manager

import (
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
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

		task, err := s.PluginCache.CreateTask(pluginID, taskName, serverConnection, args)
		if err != nil {
			logger.Errorf("Error creating plugin task: %s", err.Error())
			return
		}

		err = task.Start()
		if err != nil {
			logger.Errorf("Error running plugin task: %s", err.Error())
			return
		}

		done := make(chan bool)
		go func() {
			defer close(done)
			task.Wait()
		}()

		// check for stop every five seconds
		pollingTime := time.Second * 5
		for {
			select {
			case <-done:
				return
			case <-time.After(pollingTime):
				s.Status.setProgressPercent(task.GetProgress())

				if s.Status.stopping {
					if err := task.Stop(); err != nil {
						logger.Errorf("Error stopping plugin operation: %s", err.Error())
					}
					return
				}
			}
		}
	}()
}
