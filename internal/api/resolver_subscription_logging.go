package api

import (
	"context"

	"github.com/stashapp/stash/internal/log"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
)

func getLogLevel(logType string) models.LogLevel {
	switch logType {
	case "progress":
		return models.LogLevelProgress
	case "trace":
		return models.LogLevelTrace
	case "debug":
		return models.LogLevelDebug
	case "info":
		return models.LogLevelInfo
	case "warn":
		return models.LogLevelWarning
	case "error":
		return models.LogLevelError
	default:
		return models.LogLevelDebug
	}
}

func logEntriesFromLogItems(logItems []log.LogItem, minLevel *models.LogLevel) []*LogEntry {
	ret := make([]*LogEntry, 0, len(logItems))

	for _, entry := range logItems {
		level := getLogLevel(entry.Type)
		if minLevel != nil && level < *minLevel {
			continue
		}
		ret = append(ret, &LogEntry{
			Time:    entry.Time,
			Level:   level,
			Message: entry.Message,
		})
	}

	return ret
}

func (r *subscriptionResolver) LoggingSubscribe(ctx context.Context, minLevel *models.LogLevel) (<-chan []*LogEntry, error) {
	ret := make(chan []*LogEntry, 100)
	stop := make(chan int, 1)
	logger := manager.GetInstance().Logger
	logSub := logger.SubscribeToLog(stop)

	go func() {
		for {
			select {
			case logEntries := <-logSub:
				ret <- logEntriesFromLogItems(logEntries, minLevel)
			case <-ctx.Done():
				stop <- 0
				close(ret)
				return
			}
		}
	}()

	return ret, nil
}
