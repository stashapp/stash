package api

import (
	"context"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) JobStatus(ctx context.Context) (*models.MetadataUpdateStatus, error) {
	status := manager.GetInstance().Status
	ret := models.MetadataUpdateStatus{
		Progress: status.Progress,
		Status:   status.Status.String(),
		Message:  "",
	}

	return &ret, nil
}

func (r *queryResolver) SystemStatus(ctx context.Context) (*models.SystemStatus, error) {
	return manager.GetInstance().GetSystemStatus(), nil
}
