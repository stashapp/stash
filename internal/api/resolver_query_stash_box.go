package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) PendingFingerprintSubmissions(ctx context.Context, stashBoxEndpoint string) (ret []*models.FingerprintSubmission, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.FingerprintSubmission.FindByEndpoint(ctx, stashBoxEndpoint)
		return err
	}); err != nil {
		return nil, err
	}
	return ret, nil
}
