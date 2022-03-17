package plugin

import (
	"fmt"
	"io"

	"github.com/stashapp/stash/pkg/logger"
)

func (t *pluginTask) handlePluginStderr(name string, pluginOutputReader io.ReadCloser) {
	logLevel := logger.PluginLogLevelFromName(t.plugin.PluginErrLogLevel)
	if logLevel == nil {
		// default log level to error
		logLevel = &logger.ErrorLevel
	}

	const pluginPrefix = "[Plugin / %s] "

	lgr := logger.PluginLogger{
		Logger:          logger.Logger,
		Prefix:          fmt.Sprintf(pluginPrefix, name),
		DefaultLogLevel: logLevel,
		ProgressChan:    t.progress,
	}

	lgr.ReadLogMessages(pluginOutputReader)
}
