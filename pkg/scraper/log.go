package scraper

import (
	"bufio"
	"io"
	"strconv"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/scraper/common/log"
)

func (t *scrapeTask) handleStderrLine(line string, defaultLogLevel *log.Level) {
	level, l := log.DetectLogLevel(line)

	const scraperPrefix = "[Scrape] "
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
		logger.Trace(scraperPrefix, l)
	case log.DebugLevel:
		logger.Debug(scraperPrefix, l)
	case log.InfoLevel:
		logger.Info(scraperPrefix, l)
	case log.WarningLevel:
		logger.Warn(scraperPrefix, l)
	case log.ErrorLevel:
		logger.Error(scraperPrefix, l)
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

func (t *scrapeTask) handleScraperOutput(scraperOutputReader io.ReadCloser, defaultLogLevel *log.Level) {
	// pipe plugin stderr to our logging
	scanner := bufio.NewScanner(scraperOutputReader)
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

	scraperOutputReader.Close()
}

func (t *scrapeTask) handleScraperStderr(scraperOutputReader io.ReadCloser) {
	logLevel := log.LevelFromName(t.plugin.PluginErrLogLevel)
	if logLevel == nil {
		// default log level to error
		logLevel = &log.ErrorLevel
	}

	t.handlePluginOutput(scraperOutputReader, logLevel)
}
