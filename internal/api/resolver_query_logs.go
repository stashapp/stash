package api

import (
	"context"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) Logs(ctx context.Context, minLevel *models.LogLevel) ([]*LogEntry, error) {
	logger := manager.GetInstance().Logger
	level := ""
	if minLevel != nil {
		level = minLevel.String()
	}
	logCache := logger.GetLogCache(level)

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
