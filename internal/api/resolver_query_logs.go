package api

import (
	"context"

	"github.com/stashapp/stash/internal/manager"
)

func (r *queryResolver) Logs(ctx context.Context) ([]*LogEntry, error) {
	logger := manager.GetInstance().Logger
	logCache := logger.GetLogCache()
	ret := make([]*LogEntry, len(logCache))

	for i, entry := range logCache {
		ret[i] = &LogEntry{
			Time:    entry.Time,
			Level:   getLogLevel(entry.Type),
			Message: entry.Message,
		}
	}

	return ret, nil
}
