package api

import (
	"context"

	"github.com/stashapp/stash/internal/log"
	"github.com/stashapp/stash/internal/manager"
)

func getLogLevel(logType string) LogLevel {
	switch logType {
	case "progress":
		return LogLevelProgress
	case "trace":
		return LogLevelTrace
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn":
		return LogLevelWarning
	case "error":
		return LogLevelError
	default:
		return LogLevelDebug
	}
}

func logEntriesFromLogItems(logItems []log.LogItem) []*LogEntry {
	ret := make([]*LogEntry, len(logItems))

	for i, entry := range logItems {
		ret[i] = &LogEntry{
			Time:    entry.Time,
			Level:   getLogLevel(entry.Type),
			Message: entry.Message,
		}
	}

	return ret
}

func (r *subscriptionResolver) LoggingSubscribe(ctx context.Context) (<-chan []*LogEntry, error) {
	ret := make(chan []*LogEntry, 100)
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
