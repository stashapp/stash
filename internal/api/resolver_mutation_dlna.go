package api

import (
	"context"
	"time"

	"github.com/stashapp/stash/internal/manager"
)

func (r *mutationResolver) EnableDlna(ctx context.Context, input EnableDLNAInput) (bool, error) {
	err := manager.GetInstance().DLNAService.Start(parseMinutes(input.Duration))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) DisableDlna(ctx context.Context, input DisableDLNAInput) (bool, error) {
	manager.GetInstance().DLNAService.Stop(parseMinutes(input.Duration))
	return true, nil
}

func (r *mutationResolver) AddTempDlnaip(ctx context.Context, input AddTempDLNAIPInput) (bool, error) {
	manager.GetInstance().DLNAService.AddTempDLNAIP(input.Address, parseMinutes(input.Duration))
	return true, nil
}

func (r *mutationResolver) RemoveTempDlnaip(ctx context.Context, input RemoveTempDLNAIPInput) (bool, error) {
	ret := manager.GetInstance().DLNAService.RemoveTempDLNAIP(input.Address)
	return ret, nil
}

func parseMinutes(minutes *int) *time.Duration {
	var ret *time.Duration
	if minutes != nil {
		d := time.Duration(*minutes) * time.Minute
		ret = &d
	}

	return ret
}
