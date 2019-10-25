package api

import (
	"context"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

func getLogLevel(logType string) models.LogLevel {
	if logType == "progress" {
		return models.LogLevelProgress
	} else if logType == "debug" {
		return models.LogLevelDebug
	} else if logType == "info" {
		return models.LogLevelInfo
	} else if logType == "warn" {
		return models.LogLevelWarning
	} else if logType == "error" {
		return models.LogLevelError
	}

	// default to debug
	return models.LogLevelDebug
}

func logEntriesFromLogItems(logItems []logger.LogItem) []*models.LogEntry {
	ret := make([]*models.LogEntry, len(logItems))

	for i, entry := range logItems {
		ret[i] = &models.LogEntry{
			Time:    entry.Time,
			Level:   getLogLevel(entry.Type),
			Message: entry.Message,
		}
	}

	return ret
}

func (r *subscriptionResolver) LoggingSubscribe(ctx context.Context) (<-chan []*models.LogEntry, error) {
	ret := make(chan []*models.LogEntry, 100)
	stop := make(chan int, 1)
	logSub := logger.SubscribeToLog(stop)

	go func() {
		for {
			select {
			case logEntries := <-logSub:
				ret <- logEntriesFromLogItems(logEntries)
			case <-ctx.Done():
				stop <- 0
				close(ret)
				return
			}
		}
	}()

	return ret, nil
}
