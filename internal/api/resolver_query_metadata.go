package api

import (
	"context"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) SystemStatus(ctx context.Context) (*models.SystemStatus, error) {
	return manager.GetInstance().GetSystemStatus(), nil
}
