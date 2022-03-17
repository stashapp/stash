package api

import (
	"context"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) Logs(ctx context.Context) ([]*models.LogEntry, error) {
	logger := manager.GetInstance().Logger
	logCache := logger.GetLogCache()
	ret := make([]*models.LogEntry, len(logCache))

	for i, entry := range logCache {
		ret[i] = &models.LogEntry{
			Time:    entry.Time,
			Level:   getLogLevel(entry.Type),
			Message: entry.Message,
		}
	}

	return ret, nil
}
