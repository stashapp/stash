package api

import (
	"context"

	"github.com/stashapp/stash/internal/dlna"
	"github.com/stashapp/stash/internal/manager"
)

func (r *queryResolver) DlnaStatus(ctx context.Context) (*dlna.Status, error) {
	return manager.GetInstance().DLNAService.Status(), nil
}
