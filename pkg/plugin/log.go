package plugin

import (
	"bufio"
	"io"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/plugin/common/log"
)

func handleStderrLine(line string) {
	level, l := log.DetectLogLevel(line)

	const pluginPrefix = "[Plugin] "
	// if no log level, just output to info
	if level == nil {
		logger.Infof("%s %s", pluginPrefix, l)
		return
	}

	switch *level {
	case log.TraceLevel:
		logger.Trace(pluginPrefix, l)
	case log.DebugLevel:
		logger.Debug(pluginPrefix, l)
	case log.InfoLevel:
		logger.Info(pluginPrefix, l)
	case log.WarningLevel:
		logger.Warn(pluginPrefix, l)
	case log.ErrorLevel:
		logger.Error(pluginPrefix, l)
	}
}

func handleStderr(pluginErrReader io.ReadCloser) {
	// pipe plugin stderr to our logging
	scanner := bufio.NewScanner(pluginErrReader)
	for scanner.Scan() {
		str := scanner.Text()
		if str != "" {
			// TODO - support progress
			handleStderrLine(str)
		}
	}

	str := scanner.Text()
	if str != "" {
		// TODO - support progress
		handleStderrLine(str)
	}

	pluginErrReader.Close()
}
