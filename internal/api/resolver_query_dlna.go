package api

import (
	"context"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) DlnaStatus(ctx context.Context) (*models.DLNAStatus, error) {
	return manager.GetInstance().DLNAService.Status(), nil
}
