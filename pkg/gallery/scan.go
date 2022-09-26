package gallery

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
)

type FinderCreatorUpdater interface {
	Finder
	Create(ctx context.Context, newGallery *models.Gallery, fileIDs []file.ID) error
	UpdatePartial(ctx context.Context, id int, updatedGallery models.GalleryPartial) (*models.Gallery, error)
	AddFileID(ctx context.Context, id int, fileID file.ID) error
	models.FileLoader
}

type SceneFinderUpdater interface {
	FindByPath(ctx context.Context, p string) ([]*models.Scene, error)
	Update(ctx context.Context, updatedScene *models.Scene) error
	AddGalleryIDs(ctx context.Context, sceneID int, galleryIDs []int) error
}

type ScanHandler struct {
	CreatorUpdater     FullCreatorUpdater
	SceneFinderUpdater SceneFinderUpdater

	PluginCache *plugin.Cache
}

func (h *ScanHandler) Handle(ctx context.Context, f file.File, oldFile file.File) error {
	baseFile := f.Base()

	// try to match the file to a gallery
	existing, err := h.CreatorUpdater.FindByFileID(ctx, f.Base().ID)
	if err != nil {
		return fmt.Errorf("finding existing gallery: %w", err)
	}

	if len(existing) == 0 {
		// try also to match file by fingerprints
		existing, err = h.CreatorUpdater.FindByFingerprints(ctx, baseFile.Fingerprints)
		if err != nil {
			return fmt.Errorf("finding existing gallery by fingerprints: %w", err)
		}
	}

	if len(existing) > 0 {
		updateExisting := oldFile != nil
		if err := h.associateExisting(ctx, existing, f, updateExisting); err != nil {
			return err
		}
	} else {
		// create a new gallery
		now := time.Now()
		newGallery := &models.Gallery{
			CreatedAt: now,
			UpdatedAt: now,
		}

		logger.Infof("%s doesn't exist. Creating new gallery...", f.Base().Path)

		if err := h.CreatorUpdater.Create(ctx, newGallery, []file.ID{baseFile.ID}); err != nil {
			return fmt.Errorf("creating new gallery: %w", err)
		}

		h.PluginCache.RegisterPostHooks(ctx, newGallery.ID, plugin.GalleryCreatePost, nil, nil)

		existing = []*models.Gallery{newGallery}
	}

	if err := h.associateScene(ctx, existing, f); err != nil {
		return err
	}

	return nil
}

func (h *ScanHandler) associateExisting(ctx context.Context, existing []*models.Gallery, f file.File, updateExisting bool) error {
	for _, i := range existing {
		if err := i.LoadFiles(ctx, h.CreatorUpdater); err != nil {
			return err
		}

		found := false
		for _, sf := range i.Files.List() {
			if sf.Base().ID == f.Base().ID {
				found = true
				break
			}
		}

		if !found {
			logger.Infof("Adding %s to gallery %s", f.Base().Path, i.DisplayName())

			if err := h.CreatorUpdater.AddFileID(ctx, i.ID, f.Base().ID); err != nil {
				return fmt.Errorf("adding file to gallery: %w", err)
			}
			// update updated_at time
			if _, err := h.CreatorUpdater.UpdatePartial(ctx, i.ID, models.NewGalleryPartial()); err != nil {
				return fmt.Errorf("updating gallery: %w", err)
			}
		}

		if !found || updateExisting {
			h.PluginCache.RegisterPostHooks(ctx, i.ID, plugin.GalleryUpdatePost, nil, nil)
		}
	}

	return nil
}

func (h *ScanHandler) associateScene(ctx context.Context, existing []*models.Gallery, f file.File) error {
	galleryIDs := make([]int, len(existing))
	for i, g := range existing {
		galleryIDs[i] = g.ID
	}

	path := f.Base().Path
	withoutExt := strings.TrimSuffix(path, filepath.Ext(path)) + ".*"

	// find scenes with a file that matches
	scenes, err := h.SceneFinderUpdater.FindByPath(ctx, withoutExt)
	if err != nil {
		return err
	}

	for _, scene := range scenes {
		// found related Scene
		logger.Infof("associate: Gallery %s is related to scene: %d", path, scene.ID)
		if err := h.SceneFinderUpdater.AddGalleryIDs(ctx, scene.ID, galleryIDs); err != nil {
			return err
		}
	}

	return nil
}
