package plugin

import (
	"io"

	"github.com/stashapp/stash/pkg/logger"
)

func (t *pluginTask) handlePluginStderr(pluginOutputReader io.ReadCloser) {
	logLevel := logger.PluginLogLevelFromName(t.plugin.PluginErrLogLevel)
	if logLevel == nil {
		// default log level to error
		logLevel = &logger.ErrorLevel
	}

	const pluginPrefix = "[Plugin] "

	lgr := logger.PluginLogger{
		Prefix:          pluginPrefix,
		DefaultLogLevel: logLevel,
		ProgressChan:    t.progress,
	}

	lgr.HandlePluginStdErr(pluginOutputReader)
}
