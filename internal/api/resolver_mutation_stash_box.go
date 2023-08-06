package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/scraper/stashbox"
)

func (r *mutationResolver) SubmitStashBoxFingerprints(ctx context.Context, input StashBoxFingerprintSubmissionInput) (bool, error) {
	boxes := config.GetInstance().GetStashBoxes()

	if input.StashBoxIndex < 0 || input.StashBoxIndex >= len(boxes) {
		return false, fmt.Errorf("invalid stash_box_index %d", input.StashBoxIndex)
	}

	client := stashbox.NewClient(*boxes[input.StashBoxIndex], r.stashboxRepository())

	return client.SubmitStashBoxFingerprints(ctx, input.SceneIds, boxes[input.StashBoxIndex].Endpoint)
}

func (r *mutationResolver) StashBoxBatchPerformerTag(ctx context.Context, input manager.StashBoxBatchTagInput) (string, error) {
	jobID := manager.GetInstance().StashBoxBatchPerformerTag(ctx, input)
	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) StashBoxBatchStudioTag(ctx context.Context, input manager.StashBoxBatchTagInput) (string, error) {
	jobID := manager.GetInstance().StashBoxBatchStudioTag(ctx, input)
	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) SubmitStashBoxSceneDraft(ctx context.Context, input StashBoxDraftSubmissionInput) (*string, error) {
	boxes := config.GetInstance().GetStashBoxes()

	if input.StashBoxIndex < 0 || input.StashBoxIndex >= len(boxes) {
		return nil, fmt.Errorf("invalid stash_box_index %d", input.StashBoxIndex)
	}

	client := stashbox.NewClient(*boxes[input.StashBoxIndex], r.stashboxRepository())

	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	var res *string
	err = r.withReadTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Scene
		scene, err := qb.Find(ctx, id)
		if err != nil {
			return err
		}

		if scene == nil {
			return fmt.Errorf("scene with id %d not found", id)
		}

		cover, err := qb.GetCover(ctx, id)
		if err != nil {
			logger.Errorf("Error getting scene cover: %v", err)
		}

		if err := scene.LoadURLs(ctx, r.repository.Scene); err != nil {
			return fmt.Errorf("loading scene URLs: %w", err)
		}

		res, err = client.SubmitSceneDraft(ctx, scene, boxes[input.StashBoxIndex].Endpoint, cover)
		return err
	})

	return res, err
}

func (r *mutationResolver) SubmitStashBoxPerformerDraft(ctx context.Context, input StashBoxDraftSubmissionInput) (*string, error) {
	boxes := config.GetInstance().GetStashBoxes()

	if input.StashBoxIndex < 0 || input.StashBoxIndex >= len(boxes) {
		return nil, fmt.Errorf("invalid stash_box_index %d", input.StashBoxIndex)
	}

	client := stashbox.NewClient(*boxes[input.StashBoxIndex], r.stashboxRepository())

	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	var res *string
	err = r.withReadTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Performer
		performer, err := qb.Find(ctx, id)
		if err != nil {
			return err
		}

		if performer == nil {
			return fmt.Errorf("performer with id %d not found", id)
		}

		res, err = client.SubmitPerformerDraft(ctx, performer, boxes[input.StashBoxIndex].Endpoint)
		return err
	})

	return res, err
}
