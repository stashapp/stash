package api

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

func (r *fingerprintSubmissionResolver) Scene(ctx context.Context, obj *models.FingerprintSubmission) (*models.Scene, error) {
	var ret *models.Scene
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		ret, err = r.repository.Scene.Find(ctx, obj.SceneID)
		return err
	}); err != nil {
		return nil, err
	}

	if ret == nil {
		return nil, fmt.Errorf("scene %d not found", obj.SceneID)
	}

	return ret, nil
}
