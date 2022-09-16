package image

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

var (
	ErrNotImageFile = errors.New("not an image file")
)

type FinderCreatorUpdater interface {
	FindByFileID(ctx context.Context, fileID file.ID) ([]*models.Image, error)
	FindByFingerprints(ctx context.Context, fp []file.Fingerprint) ([]*models.Image, error)
	Create(ctx context.Context, newImage *models.ImageCreateInput) error
	UpdatePartial(ctx context.Context, id int, updatedImage models.ImagePartial) (*models.Image, error)
	AddFileID(ctx context.Context, id int, fileID file.ID) error
	models.GalleryIDLoader
	models.ImageFileLoader
}

type GalleryFinderCreator interface {
	FindByFileID(ctx context.Context, fileID file.ID) ([]*models.Gallery, error)
	FindByFolderID(ctx context.Context, folderID file.FolderID) ([]*models.Gallery, error)
	Create(ctx context.Context, newObject *models.Gallery, fileIDs []file.ID) error
}

type ScanConfig interface {
	GetCreateGalleriesFromFolders() bool
	IsGenerateThumbnails() bool
}

type ScanHandler struct {
	CreatorUpdater FinderCreatorUpdater
	GalleryFinder  GalleryFinderCreator

	ThumbnailGenerator ThumbnailGenerator

	ScanConfig ScanConfig

	PluginCache *plugin.Cache
}

func (h *ScanHandler) validate() error {
	if h.CreatorUpdater == nil {
		return errors.New("CreatorUpdater is required")
	}
	if h.GalleryFinder == nil {
		return errors.New("GalleryFinder is required")
	}
	if h.ScanConfig == nil {
		return errors.New("ScanConfig is required")
	}

	return nil
}

func (h *ScanHandler) Handle(ctx context.Context, f file.File) error {
	if err := h.validate(); err != nil {
		return err
	}

	imageFile, ok := f.(*file.ImageFile)
	if !ok {
		return ErrNotImageFile
	}

	// try to match the file to an image
	existing, err := h.CreatorUpdater.FindByFileID(ctx, imageFile.ID)
	if err != nil {
		return fmt.Errorf("finding existing image: %w", err)
	}

	if len(existing) == 0 {
		// try also to match file by fingerprints
		existing, err = h.CreatorUpdater.FindByFingerprints(ctx, imageFile.Fingerprints)
		if err != nil {
			return fmt.Errorf("finding existing image by fingerprints: %w", err)
		}
	}

	if len(existing) > 0 {
		if err := h.associateExisting(ctx, existing, imageFile); err != nil {
			return err
		}
	} else {
		// create a new image
		now := time.Now()
		newImage := &models.Image{
			CreatedAt:  now,
			UpdatedAt:  now,
			GalleryIDs: models.NewRelatedIDs([]int{}),
		}

		// if the file is in a zip, then associate it with the gallery
		if imageFile.ZipFileID != nil {
			g, err := h.GalleryFinder.FindByFileID(ctx, *imageFile.ZipFileID)
			if err != nil {
				return fmt.Errorf("finding gallery for zip file id %d: %w", *imageFile.ZipFileID, err)
			}

			for _, gg := range g {
				newImage.GalleryIDs.Add(gg.ID)
			}
		} else if h.ScanConfig.GetCreateGalleriesFromFolders() {
			if err := h.associateFolderBasedGallery(ctx, newImage, imageFile); err != nil {
				return err
			}
		}

		logger.Infof("%s doesn't exist. Creating new image...", f.Base().Path)

		if err := h.CreatorUpdater.Create(ctx, &models.ImageCreateInput{
			Image:   newImage,
			FileIDs: []file.ID{imageFile.ID},
		}); err != nil {
			return fmt.Errorf("creating new image: %w", err)
		}

		h.PluginCache.ExecutePostHooks(ctx, newImage.ID, plugin.ImageCreatePost, nil, nil)

		existing = []*models.Image{newImage}
	}

	if h.ScanConfig.IsGenerateThumbnails() {
		for _, s := range existing {
			if err := h.ThumbnailGenerator.GenerateThumbnail(ctx, s, imageFile); err != nil {
				// just log if cover generation fails. We can try again on rescan
				logger.Errorf("Error generating thumbnail for %s: %v", imageFile.Path, err)
			}
		}
	}

	return nil
}

func (h *ScanHandler) associateExisting(ctx context.Context, existing []*models.Image, f *file.ImageFile) error {
	for _, i := range existing {
		if err := i.LoadFiles(ctx, h.CreatorUpdater); err != nil {
			return err
		}

		found := false
		for _, sf := range i.Files.List() {
			if sf.ID == f.Base().ID {
				found = true
				break
			}
		}

		if !found {
			logger.Infof("Adding %s to image %s", f.Path, i.DisplayName())

			// associate with folder-based gallery if applicable
			if h.ScanConfig.GetCreateGalleriesFromFolders() {
				if err := h.associateFolderBasedGallery(ctx, i, f); err != nil {
					return err
				}
			}

			if err := h.CreatorUpdater.AddFileID(ctx, i.ID, f.ID); err != nil {
				return fmt.Errorf("adding file to image: %w", err)
			}
			// update updated_at time
			if _, err := h.CreatorUpdater.UpdatePartial(ctx, i.ID, models.NewImagePartial()); err != nil {
				return fmt.Errorf("updating image: %w", err)
			}
		}
	}

	return nil
}

func (h *ScanHandler) getOrCreateFolderBasedGallery(ctx context.Context, f file.File) (*models.Gallery, error) {
	// don't create folder-based galleries for files in zip file
	if f.Base().ZipFileID != nil {
		return nil, nil
	}

	folderID := f.Base().ParentFolderID
	g, err := h.GalleryFinder.FindByFolderID(ctx, folderID)
	if err != nil {
		return nil, fmt.Errorf("finding folder based gallery: %w", err)
	}

	if len(g) > 0 {
		gg := g[0]
		return gg, nil
	}

	// create a new folder-based gallery
	now := time.Now()
	newGallery := &models.Gallery{
		FolderID:  &folderID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	logger.Infof("Creating folder-based gallery for %s", filepath.Dir(f.Base().Path))
	if err := h.GalleryFinder.Create(ctx, newGallery, nil); err != nil {
		return nil, fmt.Errorf("creating folder based gallery: %w", err)
	}

	return newGallery, nil
}

func (h *ScanHandler) associateFolderBasedGallery(ctx context.Context, newImage *models.Image, f file.File) error {
	g, err := h.getOrCreateFolderBasedGallery(ctx, f)
	if err != nil {
		return err
	}

	if err := newImage.LoadGalleryIDs(ctx, h.CreatorUpdater); err != nil {
		return err
	}

	if g != nil && !intslice.IntInclude(newImage.GalleryIDs.List(), g.ID) {
		newImage.GalleryIDs.Add(g.ID)
		logger.Infof("Adding %s to folder-based gallery %s", f.Base().Path, g.Path)
	}

	return nil
}
