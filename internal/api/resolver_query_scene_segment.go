package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindSceneSegment(ctx context.Context, id string) (*models.SceneSegment, error) {
	segmentID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	var segment *models.SceneSegment
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		segment, err = r.repository.SceneSegment.Find(ctx, segmentID)
		return err
	}); err != nil {
		return nil, err
	}

	return segment, nil
}

func (r *queryResolver) FindSceneSegments(ctx context.Context, sceneID string) ([]*models.SceneSegment, error) {
	sid, err := strconv.Atoi(sceneID)
	if err != nil {
		return nil, fmt.Errorf("converting scene id: %w", err)
	}

	var segments []*models.SceneSegment
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		segments, err = r.repository.SceneSegment.FindBySceneID(ctx, sid)
		return err
	}); err != nil {
		return nil, err
	}

	return segments, nil
}
