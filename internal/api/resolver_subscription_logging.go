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

func logEntriesFromLogItems(logItems []log.LogItem) []*models.LogEntry {
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
	logger := manager.GetInstance().Logger
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
