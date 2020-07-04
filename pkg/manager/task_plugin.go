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

		progress := make(chan float64)
		task, err := s.PluginCache.CreateTask(pluginID, taskName, serverConnection, args, progress)
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

			output := task.GetResult()
			if output == nil {
				logger.Debug("Plugin returned no result")
			} else {
				if output.Error != nil {
					logger.Errorf("Plugin returned error: %s", *output.Error)
				} else if output.Output != nil {
					logger.Debugf("Plugin returned: %v", output.Output)
				}
			}
		}()

		// TODO - refactor stop to use channels
		// check for stop every five seconds
		pollingTime := time.Second * 5
		stopPoller := time.Tick(pollingTime)
		for {
			select {
			case <-done:
				return
			case p := <-progress:
				s.Status.setProgressPercent(p)
			case <-stopPoller:
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
