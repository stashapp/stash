package api

import (
	"context"
	"time"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
)

func (r *mutationResolver) EnableDlna(ctx context.Context, input models.EnableDLNAInput) (bool, error) {
	var duration *time.Duration
	if input.Duration != nil {
		d := time.Duration(*input.Duration) * time.Minute
		duration = &d
	}

	manager.GetInstance().DLNAService.Start(duration)
	return true, nil
}

func (r *mutationResolver) DisableDlna(ctx context.Context, input models.DisableDLNAInput) (bool, error) {
	var duration *time.Duration
	if input.Duration != nil {
		d := time.Duration(*input.Duration) * time.Minute
		duration = &d
	}

	manager.GetInstance().DLNAService.Stop(duration)
	return true, nil
}

func (r *mutationResolver) AllowDlnaip(ctx context.Context, input models.AllowDLNAIPInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) DisallowDlnaip(ctx context.Context, input models.DisallowDLNAIPInput) (bool, error) {
	panic("not implemented")
}
