package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/stashbox"
)

func (r *mutationResolver) SubmitStashBoxFingerprints(ctx context.Context, input StashBoxFingerprintSubmissionInput) (bool, error) {
	// New format: use fingerprints field with explicit stash-box scene IDs and votes
	if len(input.Fingerprints) > 0 {
		return r.submitFingerprintsNew(ctx, input.Fingerprints)
	}

	// Legacy format: use scene_ids and look up stash_ids from scenes
	b, err := resolveStashBox(nil, input.StashBoxEndpoint)
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

func (r *mutationResolver) submitFingerprintsNew(ctx context.Context, submissions []*FingerprintSubmissionInput) (bool, error) {
	// Group submissions by endpoint
	byEndpoint := make(map[string][]*FingerprintSubmissionInput)
	for _, s := range submissions {
		byEndpoint[s.StashBoxEndpoint] = append(byEndpoint[s.StashBoxEndpoint], s)
	}

	for endpoint, endpointSubmissions := range byEndpoint {
		b, err := resolveStashBox(nil, &endpoint)
		if err != nil {
			return false, err
		}

		// Collect all scene IDs for this endpoint
		sceneIDSet := make(map[string]struct{})
		for _, s := range endpointSubmissions {
			sceneIDSet[s.SceneID] = struct{}{}
		}

		sceneIDs := make([]string, 0, len(sceneIDSet))
		for id := range sceneIDSet {
			sceneIDs = append(sceneIDs, id)
		}

		ids, err := stringslice.StringSliceToIntSlice(sceneIDs)
		if err != nil {
			return false, err
		}

		var scenes []*models.Scene
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			scenes, err = r.sceneService.FindByIDs(ctx, ids, scene.LoadFiles)
			return err
		}); err != nil {
			return false, err
		}

		// Build a map of scene ID to scene for quick lookup
		sceneMap := make(map[int]*models.Scene)
		for _, s := range scenes {
			sceneMap[s.ID] = s
		}

		client := r.newStashBoxClient(*b)

		// Submit each fingerprint with its vote
		for _, sub := range endpointSubmissions {
			sceneID, err := strconv.Atoi(sub.SceneID)
			if err != nil {
				return false, fmt.Errorf("invalid scene ID %s: %w", sub.SceneID, err)
			}

			s, ok := sceneMap[sceneID]
			if !ok {
				return false, fmt.Errorf("scene %d not found", sceneID)
			}

			vote := stashbox.FingerprintVote(sub.Vote)
			if err := client.SubmitFingerprintsWithVote(ctx, s, sub.StashBoxSceneID, vote); err != nil {
				return false, err
			}
		}
	}

	return true, nil
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

func (r *mutationResolver) StashBoxBatchTagTag(ctx context.Context, input manager.StashBoxBatchTagInput) (string, error) {
	b, err := resolveStashBoxBatchTagInput(input.Endpoint, input.StashBoxEndpoint) //nolint:staticcheck
	if err != nil {
		return "", err
	}

	jobID := manager.GetInstance().StashBoxBatchTagTag(ctx, b, input)
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

func (r *mutationResolver) QueueFingerprintSubmission(ctx context.Context, input QueueFingerprintInput) (bool, error) {
	sceneID, err := strconv.Atoi(input.SceneID)
	if err != nil {
		return false, fmt.Errorf("invalid scene ID: %w", err)
	}

	submission := &models.FingerprintSubmission{
		Endpoint:  input.Endpoint,
		StashID:   input.StashID,
		SceneID:   sceneID,
		Vote:      models.FingerprintVote(input.Vote),
		CreatedAt: time.Now(),
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.repository.FingerprintSubmission.Create(ctx, submission)
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) RemoveFingerprintSubmission(ctx context.Context, input RemoveFingerprintInput) (bool, error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.repository.FingerprintSubmission.Delete(ctx, input.Endpoint, input.StashID)
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SubmitFingerprintSubmissions(ctx context.Context, endpoint string) (bool, error) {
	b, err := resolveStashBox(nil, &endpoint)
	if err != nil {
		return false, err
	}

	var submissions []*models.FingerprintSubmission
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		submissions, err = r.repository.FingerprintSubmission.FindByEndpoint(ctx, endpoint)
		return err
	}); err != nil {
		return false, err
	}

	if len(submissions) == 0 {
		return true, nil
	}

	// Collect all scene IDs
	sceneIDSet := make(map[int]struct{})
	for _, s := range submissions {
		sceneIDSet[s.SceneID] = struct{}{}
	}

	sceneIDs := make([]int, 0, len(sceneIDSet))
	for id := range sceneIDSet {
		sceneIDs = append(sceneIDs, id)
	}

	var scenes []*models.Scene
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		scenes, err = r.sceneService.FindByIDs(ctx, sceneIDs, scene.LoadFiles)
		return err
	}); err != nil {
		return false, err
	}

	// Build a map of scene ID to scene
	sceneMap := make(map[int]*models.Scene)
	for _, s := range scenes {
		sceneMap[s.ID] = s
	}

	client := r.newStashBoxClient(*b)

	// Submit each fingerprint and track successful submissions
	var successfulSubmissions []*models.FingerprintSubmission
	for _, sub := range submissions {
		s, ok := sceneMap[sub.SceneID]
		if !ok {
			logger.Warnf("Scene %d not found for fingerprint submission, skipping", sub.SceneID)
			continue
		}

		vote := stashbox.FingerprintVote(sub.Vote)
		if err := client.SubmitFingerprintsWithVote(ctx, s, sub.StashID, vote); err != nil {
			logger.Warnf("Failed to submit fingerprint for scene %d: %v", sub.SceneID, err)
			continue
		}

		successfulSubmissions = append(successfulSubmissions, sub)
	}

	// Delete successful submissions from the queue
	if len(successfulSubmissions) > 0 {
		if err := r.withTxn(ctx, func(ctx context.Context) error {
			for _, sub := range successfulSubmissions {
				if err := r.repository.FingerprintSubmission.Delete(ctx, sub.Endpoint, sub.StashID); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return false, err
		}
	}

	return true, nil
}
