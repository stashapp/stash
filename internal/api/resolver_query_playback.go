package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) PlaybackDefaults(ctx context.Context) ([]*models.PlaybackDefault, error) {
	var ret []*models.PlaybackDefault
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		ret, err = r.repository.PlaybackDefault.GetAll(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
