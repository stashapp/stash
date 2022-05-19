package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/scraper/stashbox"
)

func (r *Resolver) stashboxRepository() stashbox.Repository {
	return stashbox.Repository{
		Scene:     r.repository.Scene,
		Performer: r.repository.Performer,
		Tag:       r.repository.Tag,
		Studio:    r.repository.Studio,
	}
}

func (r *mutationResolver) SubmitStashBoxFingerprints(ctx context.Context, input StashBoxFingerprintSubmissionInput) (bool, error) {
	boxes := config.GetInstance().GetStashBoxes()

	if input.StashBoxIndex < 0 || input.StashBoxIndex >= len(boxes) {
		return false, fmt.Errorf("invalid stash_box_index %d", input.StashBoxIndex)
	}

	client := stashbox.NewClient(*boxes[input.StashBoxIndex], r.txnManager, r.stashboxRepository())

	return client.SubmitStashBoxFingerprints(ctx, input.SceneIds, boxes[input.StashBoxIndex].Endpoint)
}

func (r *mutationResolver) StashBoxBatchPerformerTag(ctx context.Context, input manager.StashBoxBatchPerformerTagInput) (string, error) {
	jobID := manager.GetInstance().StashBoxBatchPerformerTag(ctx, input)
	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) SubmitStashBoxSceneDraft(ctx context.Context, input StashBoxDraftSubmissionInput) (*string, error) {
	boxes := config.GetInstance().GetStashBoxes()

	if input.StashBoxIndex < 0 || input.StashBoxIndex >= len(boxes) {
		return nil, fmt.Errorf("invalid stash_box_index %d", input.StashBoxIndex)
	}

	client := stashbox.NewClient(*boxes[input.StashBoxIndex], r.txnManager, r.stashboxRepository())

	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	var res *string
	err = r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Scene
		scene, err := qb.Find(ctx, id)
		if err != nil {
			return err
		}
		filepath := manager.GetInstance().Paths.Scene.GetScreenshotPath(scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm()))

		res, err = client.SubmitSceneDraft(ctx, id, boxes[input.StashBoxIndex].Endpoint, filepath)
		return err
	})

	return res, err
}

func (r *mutationResolver) SubmitStashBoxPerformerDraft(ctx context.Context, input StashBoxDraftSubmissionInput) (*string, error) {
	boxes := config.GetInstance().GetStashBoxes()

	if input.StashBoxIndex < 0 || input.StashBoxIndex >= len(boxes) {
		return nil, fmt.Errorf("invalid stash_box_index %d", input.StashBoxIndex)
	}

	client := stashbox.NewClient(*boxes[input.StashBoxIndex], r.txnManager, r.stashboxRepository())

	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	var res *string
	err = r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Performer
		performer, err := qb.Find(ctx, id)
		if err != nil {
			return err
		}

		res, err = client.SubmitPerformerDraft(ctx, performer, boxes[input.StashBoxIndex].Endpoint)
		return err
	})

	return res, err
}
