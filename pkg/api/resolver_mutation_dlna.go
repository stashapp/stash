package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *mutationResolver) EnableDlna(ctx context.Context, input models.EnableDLNAInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) DisableDlna(ctx context.Context, input models.DisableDLNAInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) AllowDlnaip(ctx context.Context, input models.AllowDLNAIPInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) DisallowDlnaip(ctx context.Context, input models.DisallowDLNAIPInput) (bool, error) {
	panic("not implemented")
}
