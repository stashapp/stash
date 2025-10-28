package scene

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/hook"
)

func (s *Service) Create(ctx context.Context, input *models.Scene, fileIDs []models.FileID, coverImage []byte) (*models.Scene, error) {
	// title must be set if no files are provided
	if input.Title == "" && len(fileIDs) == 0 {
		return nil, errors.New("title must be set if scene has no files")
	}

	now := time.Now()
	newScene := *input
	newScene.CreatedAt = now
	newScene.UpdatedAt = now

	// don't pass the file ids since they may be already assigned
	// assign them afterwards
	if err := s.Repository.Create(ctx, &newScene, nil); err != nil {
		return nil, fmt.Errorf("creating new scene: %w", err)
	}

	for _, f := range fileIDs {
		if err := s.AssignFile(ctx, newScene.ID, f); err != nil {
			return nil, fmt.Errorf("assigning file %d to new scene: %w", f, err)
		}
	}

	if len(fileIDs) > 0 {
		// assign the primary to the first
		if _, err := s.Repository.UpdatePartial(ctx, newScene.ID, models.ScenePartial{
			PrimaryFileID: &fileIDs[0],
		}); err != nil {
			return nil, fmt.Errorf("setting primary file on new scene: %w", err)
		}
	}

	// re-find the scene so that it correctly returns file-related fields
	ret, err := s.Repository.Find(ctx, newScene.ID)
	if err != nil {
		return nil, err
	}

	if len(coverImage) > 0 {
		if err := s.Repository.UpdateCover(ctx, ret.ID, coverImage); err != nil {
			return nil, fmt.Errorf("setting cover on new scene: %w", err)
		}
	}

	s.PluginCache.RegisterPostHooks(ctx, ret.ID, hook.SceneCreatePost, nil, nil)

	// re-find the scene so that it correctly returns file-related fields
	return ret, nil
}
