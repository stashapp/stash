package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/performer"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/scraper/stashbox"
)

func (r *mutationResolver) SubmitStashBoxFingerprints(ctx context.Context, input models.StashBoxFingerprintSubmissionInput) (bool, error) {
	boxes := config.GetInstance().GetStashBoxes()

	if input.StashBoxIndex < 0 || input.StashBoxIndex >= len(boxes) {
		return false, fmt.Errorf("invalid stash_box_index %d", input.StashBoxIndex)
	}

	client := stashbox.NewClient(*boxes[input.StashBoxIndex])
	client.TxnManager = r.txnManager

	return client.SubmitStashBoxFingerprints(ctx, input.SceneIds, boxes[input.StashBoxIndex].Endpoint)
}

func (r *mutationResolver) StashBoxBatchPerformerTag(ctx context.Context, input models.StashBoxBatchPerformerTagInput) (string, error) {
	jobID := manager.GetInstance().StashBoxBatchPerformerTag(ctx, input)
	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) SubmitStashBoxSceneDraft(ctx context.Context, input models.StashBoxDraftSubmissionInput) (*string, error) {
	boxes := config.GetInstance().GetStashBoxes()

	if input.StashBoxIndex < 0 || input.StashBoxIndex >= len(boxes) {
		return nil, fmt.Errorf("invalid stash_box_index %d", input.StashBoxIndex)
	}

	client := stashbox.NewClient(*boxes[input.StashBoxIndex])
	client.TxnManager = r.txnManager
	client.Scenes = &scene.FileCoverGetter{
		PathGetter: manager.GetInstance().Paths.Scene,
	}

	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	res, err := client.SubmitSceneDraft(ctx, id, boxes[input.StashBoxIndex].Endpoint)

	return res, err
}

func (r *mutationResolver) SubmitStashBoxPerformerDraft(ctx context.Context, input models.StashBoxDraftSubmissionInput) (*string, error) {
	boxes := config.GetInstance().GetStashBoxes()

	if input.StashBoxIndex < 0 || input.StashBoxIndex >= len(boxes) {
		return nil, fmt.Errorf("invalid stash_box_index %d", input.StashBoxIndex)
	}

	client := stashbox.NewClient(*boxes[input.StashBoxIndex])
	client.TxnManager = r.txnManager
	client.Performers = &performer.FileImageGetter{
		PathGetter: manager.GetInstance().Paths.Performer,
	}

	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	var res *string
	err = r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		qb := repo.Performer()
		performer, err := qb.Find(id)
		if err != nil {
			return err
		}

		res, err = client.SubmitPerformerDraft(ctx, performer, boxes[input.StashBoxIndex].Endpoint)
		return err
	})

	return res, err
}
