package gallery

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/plugin/hook"
)

type ScanCreatorUpdater interface {
	FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Gallery, error)
	FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Gallery, error)
	GetFiles(ctx context.Context, relatedID int) ([]models.File, error)

	Create(ctx context.Context, newGallery *models.Gallery, fileIDs []models.FileID) error
	UpdatePartial(ctx context.Context, id int, updatedGallery models.GalleryPartial) (*models.Gallery, error)
	AddFileID(ctx context.Context, id int, fileID models.FileID) error
}

type ScanSceneFinderUpdater interface {
	FindByPath(ctx context.Context, p string) ([]*models.Scene, error)
	Update(ctx context.Context, updatedScene *models.Scene) error
	AddGalleryIDs(ctx context.Context, sceneID int, galleryIDs []int) error
}

type ScanImageFinderUpdater interface {
	FindByZipFileID(ctx context.Context, zipFileID models.FileID) ([]*models.Image, error)
	UpdatePartial(ctx context.Context, id int, partial models.ImagePartial) (*models.Image, error)
}

type ScanHandler struct {
	CreatorUpdater     ScanCreatorUpdater
	SceneFinderUpdater ScanSceneFinderUpdater
	ImageFinderUpdater ScanImageFinderUpdater
	PluginCache        *plugin.Cache
}

func (h *ScanHandler) Handle(ctx context.Context, f models.File, oldFile models.File) error {
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
		// only create galleries if there is something to put in them
		// otherwise, they will be created on the fly when an image is created
		images, err := h.ImageFinderUpdater.FindByZipFileID(ctx, f.Base().ID)
		if err != nil {
			return err
		}

		if len(images) == 0 {
			// don't create an empty gallery
			return nil
		}

		// create a new gallery
		newGallery := models.NewGallery()

		logger.Infof("%s doesn't exist. Creating new gallery...", f.Base().Path)

		if err := h.CreatorUpdater.Create(ctx, &newGallery, []models.FileID{baseFile.ID}); err != nil {
			return fmt.Errorf("creating new gallery: %w", err)
		}

		h.PluginCache.RegisterPostHooks(ctx, newGallery.ID, hook.GalleryCreatePost, nil, nil)

		// associate all the images in the zip file with the gallery
		for _, i := range images {
			imagePartial := models.ImagePartial{
				GalleryIDs: &models.UpdateIDs{
					IDs:  []int{newGallery.ID},
					Mode: models.RelationshipUpdateModeAdd,
				},
				// set UpdatedAt directly instead of using NewImagePartial, to ensure
				// that the images have the same UpdatedAt time as the gallery
				UpdatedAt: models.NewOptionalTime(newGallery.UpdatedAt),
			}
			if _, err := h.ImageFinderUpdater.UpdatePartial(ctx, i.ID, imagePartial); err != nil {
				return fmt.Errorf("adding image %s to gallery: %w", i.Path, err)
			}
		}

		existing = []*models.Gallery{&newGallery}
	}

	if err := h.associateScene(ctx, existing, f); err != nil {
		return err
	}

	return nil
}

func (h *ScanHandler) associateExisting(ctx context.Context, existing []*models.Gallery, f models.File, updateExisting bool) error {
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
			h.PluginCache.RegisterPostHooks(ctx, i.ID, hook.GalleryUpdatePost, nil, nil)
		}
	}

	return nil
}

func (h *ScanHandler) associateScene(ctx context.Context, existing []*models.Gallery, f models.File) error {
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
