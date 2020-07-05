package plugin

import (
	"bufio"
	"io"
	"strconv"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/plugin/common/log"
)

func (t *pluginTask) handleStderrLine(line string, defaultLogLevel *log.Level) {
	level, l := log.DetectLogLevel(line)

	const pluginPrefix = "[Plugin] "
	// if no log level, just output to info
	if level == nil {
		if defaultLogLevel != nil {
			level = defaultLogLevel
		} else {
			level = &log.InfoLevel
		}
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
	case log.ProgressLevel:
		progress, err := strconv.ParseFloat(l, 64)
		if err != nil {
			logger.Errorf("Error parsing progress value '%s': %s", l, err.Error())
		} else {
			// only pass progress through if channel present
			if t.progress != nil {
				// don't block on this
				select {
				case t.progress <- progress:
				default:
				}
			}
		}
	}
}

func (t *pluginTask) handlePluginOutput(pluginOutputReader io.ReadCloser, defaultLogLevel *log.Level) {
	// pipe plugin stderr to our logging
	scanner := bufio.NewScanner(pluginOutputReader)
	for scanner.Scan() {
		str := scanner.Text()
		if str != "" {
			t.handleStderrLine(str, defaultLogLevel)
		}
	}

	str := scanner.Text()
	if str != "" {
		t.handleStderrLine(str, defaultLogLevel)
	}

	pluginOutputReader.Close()
}

func (t *pluginTask) handlePluginStderr(pluginOutputReader io.ReadCloser) {
	logLevel := log.LevelFromName(t.plugin.PluginErrLogLevel)
	if logLevel == nil {
		// default log level to error
		logLevel = &log.ErrorLevel
	}

	t.handlePluginOutput(pluginOutputReader, logLevel)
}
