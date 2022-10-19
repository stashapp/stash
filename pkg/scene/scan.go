package scene

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/plugin"
)

var (
	ErrNotVideoFile = errors.New("not a video file")
)

type CreatorUpdater interface {
	FindByFileID(ctx context.Context, fileID file.ID) ([]*models.Scene, error)
	FindByFingerprints(ctx context.Context, fp []file.Fingerprint) ([]*models.Scene, error)
	Create(ctx context.Context, newScene *models.Scene, fileIDs []file.ID) error
	UpdatePartial(ctx context.Context, id int, updatedScene models.ScenePartial) (*models.Scene, error)
	AddFileID(ctx context.Context, id int, fileID file.ID) error
	models.VideoFileLoader
}

type ScanGenerator interface {
	Generate(ctx context.Context, s *models.Scene, f *file.VideoFile) error
}

type ScanHandler struct {
	CreatorUpdater CreatorUpdater

	CoverGenerator CoverGenerator
	ScanGenerator  ScanGenerator
	PluginCache    *plugin.Cache

	FileNamingAlgorithm models.HashAlgorithm
	Paths               *paths.Paths
}

func (h *ScanHandler) validate() error {
	if h.CreatorUpdater == nil {
		return errors.New("CreatorUpdater is required")
	}
	if h.CoverGenerator == nil {
		return errors.New("CoverGenerator is required")
	}
	if h.ScanGenerator == nil {
		return errors.New("ScanGenerator is required")
	}
	if !h.FileNamingAlgorithm.IsValid() {
		return errors.New("FileNamingAlgorithm is required")
	}
	if h.Paths == nil {
		return errors.New("Paths is required")
	}

	return nil
}

func (h *ScanHandler) Handle(ctx context.Context, f file.File, oldFile file.File) error {
	if err := h.validate(); err != nil {
		return err
	}

	videoFile, ok := f.(*file.VideoFile)
	if !ok {
		return ErrNotVideoFile
	}

	// try to match the file to a scene
	existing, err := h.CreatorUpdater.FindByFileID(ctx, f.Base().ID)
	if err != nil {
		return fmt.Errorf("finding existing scene: %w", err)
	}

	if len(existing) == 0 {
		// try also to match file by fingerprints
		existing, err = h.CreatorUpdater.FindByFingerprints(ctx, videoFile.Fingerprints)
		if err != nil {
			return fmt.Errorf("finding existing scene by fingerprints: %w", err)
		}
	}

	if len(existing) > 0 {
		updateExisting := oldFile != nil
		if err := h.associateExisting(ctx, existing, videoFile, updateExisting); err != nil {
			return err
		}
	} else {
		// create a new scene
		now := time.Now()
		newScene := &models.Scene{
			CreatedAt: now,
			UpdatedAt: now,
		}

		logger.Infof("%s doesn't exist. Creating new scene...", f.Base().Path)

		if err := h.CreatorUpdater.Create(ctx, newScene, []file.ID{videoFile.ID}); err != nil {
			return fmt.Errorf("creating new scene: %w", err)
		}

		h.PluginCache.RegisterPostHooks(ctx, newScene.ID, plugin.SceneCreatePost, nil, nil)

		existing = []*models.Scene{newScene}
	}

	if oldFile != nil {
		// migrate hashes from the old file to the new
		oldHash := GetHash(oldFile, h.FileNamingAlgorithm)
		newHash := GetHash(f, h.FileNamingAlgorithm)

		if oldHash != "" && newHash != "" && oldHash != newHash {
			MigrateHash(h.Paths, oldHash, newHash)
		}
	}

	for _, s := range existing {
		if err := h.CoverGenerator.GenerateCover(ctx, s, videoFile); err != nil {
			// just log if cover generation fails. We can try again on rescan
			logger.Errorf("Error generating cover for %s: %v", videoFile.Path, err)
		}

		if err := h.ScanGenerator.Generate(ctx, s, videoFile); err != nil {
			// just log if cover generation fails. We can try again on rescan
			logger.Errorf("Error generating content for %s: %v", videoFile.Path, err)
		}
	}

	return nil
}

func (h *ScanHandler) associateExisting(ctx context.Context, existing []*models.Scene, f *file.VideoFile, updateExisting bool) error {
	for _, s := range existing {
		if err := s.LoadFiles(ctx, h.CreatorUpdater); err != nil {
			return err
		}

		found := false
		for _, sf := range s.Files.List() {
			if sf.ID == f.ID {
				found = true
				break
			}
		}

		if !found {
			logger.Infof("Adding %s to scene %s", f.Path, s.DisplayName())

			if err := h.CreatorUpdater.AddFileID(ctx, s.ID, f.ID); err != nil {
				return fmt.Errorf("adding file to scene: %w", err)
			}

			// update updated_at time
			if _, err := h.CreatorUpdater.UpdatePartial(ctx, s.ID, models.NewScenePartial()); err != nil {
				return fmt.Errorf("updating scene: %w", err)
			}
		}

		if !found || updateExisting {
			h.PluginCache.RegisterPostHooks(ctx, s.ID, plugin.SceneUpdatePost, nil, nil)
		}
	}

	return nil
}
