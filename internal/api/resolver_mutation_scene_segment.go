package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

func (r *mutationResolver) SceneSegmentCreate(ctx context.Context, input SceneSegmentCreateInput) (*models.SceneSegment, error) {
	sceneID, err := strconv.Atoi(input.SceneID)
	if err != nil {
		return nil, fmt.Errorf("converting scene id: %w", err)
	}

	// Populate a new scene segment from the input
	newSegment := models.NewSceneSegment()

	newSegment.Title = strings.TrimSpace(input.Title)
	newSegment.SceneID = sceneID
	newSegment.StartSeconds = input.StartSeconds
	newSegment.EndSeconds = input.EndSeconds

	// Validate
	if err := newSegment.Validate(); err != nil {
		return nil, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.SceneSegment

		return qb.Create(ctx, &newSegment)
	}); err != nil {
		return nil, err
	}

	return &newSegment, nil
}

func (r *mutationResolver) SceneSegmentUpdate(ctx context.Context, input SceneSegmentUpdateInput) (*models.SceneSegment, error) {
	segmentID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate scene segment from the input
	updatedSegment := models.SceneSegmentPartial{}

	updatedSegment.Title = translator.optionalString(input.Title, "title")
	updatedSegment.StartSeconds = translator.optionalFloat64(input.StartSeconds, "start_seconds")
	updatedSegment.EndSeconds = translator.optionalFloat64(input.EndSeconds, "end_seconds")

	var ret *models.SceneSegment
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.SceneSegment

		ret, err = qb.UpdatePartial(ctx, segmentID, updatedSegment)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) SceneSegmentDestroy(ctx context.Context, id string) (bool, error) {
	segmentID, err := strconv.Atoi(id)
	if err != nil {
		return false, fmt.Errorf("converting id: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.repository.SceneSegment.Destroy(ctx, segmentID)
	}); err != nil {
		return false, err
	}

	return true, nil
}
