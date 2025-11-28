package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/stashbox"
)

func (r *mutationResolver) SubmitStashBoxFingerprints(ctx context.Context, input StashBoxFingerprintSubmissionInput) (bool, error) {
	b, err := resolveStashBox(input.StashBoxIndex, input.StashBoxEndpoint)
	if err != nil {
		return false, err
	}

	ids, err := stringslice.StringSliceToIntSlice(input.SceneIds)
	if err != nil {
		return false, err
	}

	client := r.newStashBoxClient(*b)

	var scenes []*models.Scene

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		scenes, err = r.sceneService.FindByIDs(ctx, ids, scene.LoadStashIDs, scene.LoadFiles)
		return err
	}); err != nil {
		return false, err
	}

	return client.SubmitFingerprints(ctx, scenes)
}

func (r *mutationResolver) StashBoxBatchPerformerTag(ctx context.Context, input manager.StashBoxBatchTagInput) (string, error) {
	b, err := resolveStashBoxBatchTagInput(input.Endpoint, input.StashBoxEndpoint) //nolint:staticcheck
	if err != nil {
		return "", err
	}

	jobID := manager.GetInstance().StashBoxBatchPerformerTag(ctx, b, input)
	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) StashBoxBatchStudioTag(ctx context.Context, input manager.StashBoxBatchTagInput) (string, error) {
	b, err := resolveStashBoxBatchTagInput(input.Endpoint, input.StashBoxEndpoint) //nolint:staticcheck
	if err != nil {
		return "", err
	}

	jobID := manager.GetInstance().StashBoxBatchStudioTag(ctx, b, input)
	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) SubmitStashBoxSceneDraft(ctx context.Context, input StashBoxDraftSubmissionInput) (*string, error) {
	b, err := resolveStashBox(input.StashBoxIndex, input.StashBoxEndpoint)
	if err != nil {
		return nil, err
	}

	client := r.newStashBoxClient(*b)

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

		draft, err := r.makeSceneDraft(ctx, scene, cover)
		if err != nil {
			return err
		}

		res, err = client.SubmitSceneDraft(ctx, *draft)
		return err
	})

	return res, err
}

func (r *mutationResolver) makeSceneDraft(ctx context.Context, s *models.Scene, cover []byte) (*stashbox.SceneDraft, error) {
	if err := s.LoadURLs(ctx, r.repository.Scene); err != nil {
		return nil, fmt.Errorf("loading scene URLs: %w", err)
	}

	if err := s.LoadStashIDs(ctx, r.repository.Scene); err != nil {
		return nil, err
	}

	draft := &stashbox.SceneDraft{
		Scene: s,
	}

	pqb := r.repository.Performer
	sqb := r.repository.Studio

	if s.StudioID != nil {
		var err error
		draft.Studio, err = sqb.Find(ctx, *s.StudioID)
		if err != nil {
			return nil, err
		}
		if draft.Studio == nil {
			return nil, fmt.Errorf("studio with id %d not found", *s.StudioID)
		}

		if err := draft.Studio.LoadStashIDs(ctx, r.repository.Studio); err != nil {
			return nil, err
		}
	}

	// submit all file fingerprints
	if err := s.LoadFiles(ctx, r.repository.Scene); err != nil {
		return nil, err
	}

	scenePerformers, err := pqb.FindBySceneID(ctx, s.ID)
	if err != nil {
		return nil, err
	}

	for _, p := range scenePerformers {
		if err := p.LoadStashIDs(ctx, pqb); err != nil {
			return nil, err
		}
	}
	draft.Performers = scenePerformers

	draft.Tags, err = r.repository.Tag.FindBySceneID(ctx, s.ID)
	if err != nil {
		return nil, err
	}

	// Load StashIDs for tags
	tqb := r.repository.Tag
	for _, t := range draft.Tags {
		if err := t.LoadStashIDs(ctx, tqb); err != nil {
			return nil, err
		}
	}

	draft.Cover = cover

	return draft, nil
}

func (r *mutationResolver) SubmitStashBoxPerformerDraft(ctx context.Context, input StashBoxDraftSubmissionInput) (*string, error) {
	b, err := resolveStashBox(input.StashBoxIndex, input.StashBoxEndpoint)
	if err != nil {
		return nil, err
	}

	client := r.newStashBoxClient(*b)

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

		pqb := r.repository.Performer
		if err := performer.LoadAliases(ctx, pqb); err != nil {
			return err
		}

		if err := performer.LoadURLs(ctx, pqb); err != nil {
			return err
		}

		if err := performer.LoadStashIDs(ctx, pqb); err != nil {
			return err
		}

		img, _ := pqb.GetImage(ctx, performer.ID)

		res, err = client.SubmitPerformerDraft(ctx, performer, img)
		return err
	})

	return res, err
}
