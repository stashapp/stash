package manager

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common"
)

func (s *singleton) MakeServerConnection(cookie *http.Cookie, hasTLSConfig bool) common.StashServerConnection {
	config := s.Config
	serverConnection := common.StashServerConnection{
		Scheme:        "http",
		Port:          config.GetPort(),
		SessionCookie: cookie,
		Dir:           config.GetConfigPath(),
	}

	if hasTLSConfig {
		serverConnection.Scheme = "https"
	}

	return serverConnection
}

func (s *singleton) RunPluginTask(pluginID string, taskName string, args []*models.PluginArgInput, serverConnection common.StashServerConnection) int {
	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) {
		pluginProgress := make(chan float64)
		task, err := s.PluginCache.CreateTask(pluginID, taskName, serverConnection, args, pluginProgress)
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

		for {
			select {
			case <-done:
				return
			case p := <-pluginProgress:
				progress.SetPercent(p)
			case <-ctx.Done():
				if err := task.Stop(); err != nil {
					logger.Errorf("Error stopping plugin operation: %s", err.Error())
				}
				return
			}
		}
	})

	return s.JobManager.Add(fmt.Sprintf("Running plugin task: %s", taskName), j)
}
