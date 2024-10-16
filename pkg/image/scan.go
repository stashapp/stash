package image

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/plugin/hook"
	"github.com/stashapp/stash/pkg/txn"
)

var (
	ErrNotImageFile = errors.New("not an image file")
)

type ScanCreatorUpdater interface {
	FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Image, error)
	FindByFolderID(ctx context.Context, folderID models.FolderID) ([]*models.Image, error)
	FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Image, error)
	GetFiles(ctx context.Context, relatedID int) ([]models.File, error)
	GetGalleryIDs(ctx context.Context, relatedID int) ([]int, error)

	Create(ctx context.Context, newImage *models.Image, fileIDs []models.FileID) error
	UpdatePartial(ctx context.Context, id int, updatedImage models.ImagePartial) (*models.Image, error)
	AddFileID(ctx context.Context, id int, fileID models.FileID) error
}

type GalleryFinderCreator interface {
	FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Gallery, error)
	FindByFolderID(ctx context.Context, folderID models.FolderID) ([]*models.Gallery, error)
	Create(ctx context.Context, newObject *models.Gallery, fileIDs []models.FileID) error
	UpdatePartial(ctx context.Context, id int, updatedGallery models.GalleryPartial) (*models.Gallery, error)
}

type ScanConfig interface {
	GetCreateGalleriesFromFolders() bool
}

type ScanGenerator interface {
	Generate(ctx context.Context, i *models.Image, f models.File) error
}

type ScanHandler struct {
	CreatorUpdater ScanCreatorUpdater
	GalleryFinder  GalleryFinderCreator

	ScanGenerator ScanGenerator

	ScanConfig ScanConfig

	PluginCache *plugin.Cache

	Paths *paths.Paths
}

func (h *ScanHandler) validate() error {
	if h.CreatorUpdater == nil {
		return errors.New("CreatorUpdater is required")
	}
	if h.ScanGenerator == nil {
		return errors.New("ScanGenerator is required")
	}
	if h.GalleryFinder == nil {
		return errors.New("GalleryFinder is required")
	}
	if h.ScanConfig == nil {
		return errors.New("ScanConfig is required")
	}
	if h.Paths == nil {
		return errors.New("Paths is required")
	}

	return nil
}

func (h *ScanHandler) Handle(ctx context.Context, f models.File, oldFile models.File) error {
	if err := h.validate(); err != nil {
		return err
	}

	imageFile := f.Base()

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
		updateExisting := oldFile != nil

		if err := h.associateExisting(ctx, existing, imageFile, updateExisting); err != nil {
			return err
		}
	} else {
		// create a new image
		newImage := models.NewImage()
		newImage.GalleryIDs = models.NewRelatedIDs([]int{})

		logger.Infof("%s doesn't exist. Creating new image...", f.Base().Path)

		g, err := h.getGalleryToAssociate(ctx, &newImage, f)
		if err != nil {
			return err
		}

		if g != nil {
			newImage.GalleryIDs.Add(g.ID)
			logger.Infof("Adding %s to gallery %s", f.Base().Path, g.Path)
		}

		if err := h.CreatorUpdater.Create(ctx, &newImage, []models.FileID{imageFile.ID}); err != nil {
			return fmt.Errorf("creating new image: %w", err)
		}

		// update the gallery updated at timestamp if applicable
		if g != nil {
			galleryPartial := models.GalleryPartial{
				UpdatedAt: models.NewOptionalTime(newImage.UpdatedAt),
			}
			if _, err := h.GalleryFinder.UpdatePartial(ctx, g.ID, galleryPartial); err != nil {
				return fmt.Errorf("updating gallery updated at timestamp: %w", err)
			}
		}

		h.PluginCache.RegisterPostHooks(ctx, newImage.ID, hook.ImageCreatePost, nil, nil)

		existing = []*models.Image{&newImage}
	}

	// remove the old thumbnail if the checksum changed - we'll regenerate it
	if oldFile != nil {
		oldHash := oldFile.Base().Fingerprints.GetString(models.FingerprintTypeMD5)
		newHash := f.Base().Fingerprints.GetString(models.FingerprintTypeMD5)

		if oldHash != "" && newHash != "" && oldHash != newHash {
			// remove cache dir of gallery
			_ = os.Remove(h.Paths.Generated.GetThumbnailPath(oldHash, models.DefaultGthumbWidth))
		}
	}

	// do this after the commit so that generation doesn't hold up the transaction
	txn.AddPostCommitHook(ctx, func(ctx context.Context) {
		for _, s := range existing {
			if err := h.ScanGenerator.Generate(ctx, s, f); err != nil {
				// just log if cover generation fails. We can try again on rescan
				logger.Errorf("Error generating content for %s: %v", imageFile.Path, err)
			}
		}
	})

	return nil
}

func (h *ScanHandler) associateExisting(ctx context.Context, existing []*models.Image, f *models.BaseFile, updateExisting bool) error {
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

		// associate with gallery if applicable
		g, err := h.getGalleryToAssociate(ctx, i, f)
		if err != nil {
			return err
		}

		var galleryIDs *models.UpdateIDs
		changed := false
		if g != nil {
			changed = true
			galleryIDs = &models.UpdateIDs{
				IDs:  []int{g.ID},
				Mode: models.RelationshipUpdateModeAdd,
			}
		}

		if !found {
			logger.Infof("Adding %s to image %s", f.Path, i.DisplayName())

			if err := h.CreatorUpdater.AddFileID(ctx, i.ID, f.ID); err != nil {
				return fmt.Errorf("adding file to image: %w", err)
			}

			changed = true
		}

		if changed {
			// always update updated_at time
			imagePartial := models.NewImagePartial()
			imagePartial.GalleryIDs = galleryIDs

			if _, err := h.CreatorUpdater.UpdatePartial(ctx, i.ID, imagePartial); err != nil {
				return fmt.Errorf("updating image: %w", err)
			}

			if g != nil {
				galleryPartial := models.GalleryPartial{
					// set UpdatedAt directly instead of using NewGalleryPartial, to ensure
					// that the linked gallery has the same UpdatedAt time as this image
					UpdatedAt: imagePartial.UpdatedAt,
				}
				if _, err := h.GalleryFinder.UpdatePartial(ctx, g.ID, galleryPartial); err != nil {
					return fmt.Errorf("updating gallery updated at timestamp: %w", err)
				}
			}
		}

		if changed || updateExisting {
			h.PluginCache.RegisterPostHooks(ctx, i.ID, hook.ImageUpdatePost, nil, nil)
		}
	}

	return nil
}

func (h *ScanHandler) getOrCreateFolderBasedGallery(ctx context.Context, f models.File) (*models.Gallery, error) {
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
	newGallery := models.NewGallery()
	newGallery.FolderID = &folderID

	logger.Infof("Creating folder-based gallery for %s", filepath.Dir(f.Base().Path))

	if err := h.GalleryFinder.Create(ctx, &newGallery, nil); err != nil {
		return nil, fmt.Errorf("creating folder based gallery: %w", err)
	}

	h.PluginCache.RegisterPostHooks(ctx, newGallery.ID, hook.GalleryCreatePost, nil, nil)

	// it's possible that there are other images in the folder that
	// need to be added to the new gallery. Find and add them now.
	if err := h.associateFolderImages(ctx, &newGallery); err != nil {
		return nil, fmt.Errorf("associating existing folder images: %w", err)
	}

	return &newGallery, nil
}

func (h *ScanHandler) associateFolderImages(ctx context.Context, g *models.Gallery) error {
	i, err := h.CreatorUpdater.FindByFolderID(ctx, *g.FolderID)
	if err != nil {
		return fmt.Errorf("finding images in folder: %w", err)
	}

	for _, ii := range i {
		logger.Infof("Adding %s to gallery %s", ii.Path, g.Path)

		imagePartial := models.NewImagePartial()
		imagePartial.GalleryIDs = &models.UpdateIDs{
			IDs:  []int{g.ID},
			Mode: models.RelationshipUpdateModeAdd,
		}

		if _, err := h.CreatorUpdater.UpdatePartial(ctx, ii.ID, imagePartial); err != nil {
			return fmt.Errorf("updating image: %w", err)
		}
	}

	return nil
}

func (h *ScanHandler) getOrCreateZipBasedGallery(ctx context.Context, zipFile models.File) (*models.Gallery, error) {
	g, err := h.GalleryFinder.FindByFileID(ctx, zipFile.Base().ID)
	if err != nil {
		return nil, fmt.Errorf("finding zip based gallery: %w", err)
	}

	if len(g) > 0 {
		gg := g[0]
		return gg, nil
	}

	// create a new zip-based gallery
	newGallery := models.NewGallery()

	logger.Infof("%s doesn't exist. Creating new gallery...", zipFile.Base().Path)

	if err := h.GalleryFinder.Create(ctx, &newGallery, []models.FileID{zipFile.Base().ID}); err != nil {
		return nil, fmt.Errorf("creating zip-based gallery: %w", err)
	}

	h.PluginCache.RegisterPostHooks(ctx, newGallery.ID, hook.GalleryCreatePost, nil, nil)

	return &newGallery, nil
}

func (h *ScanHandler) getOrCreateGallery(ctx context.Context, f models.File) (*models.Gallery, error) {
	// don't create folder-based galleries for files in zip file
	if f.Base().ZipFile != nil {
		return h.getOrCreateZipBasedGallery(ctx, f.Base().ZipFile)
	}

	// Look for specific filename in Folder to find out if the Folder is marked to be handled differently as the setting
	folderPath := filepath.Dir(f.Base().Path)

	forceGallery := false
	if _, err := os.Stat(filepath.Join(folderPath, ".forcegallery")); err == nil {
		forceGallery = true
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("Could not test Path %s: %w", folderPath, err)
	}
	exemptGallery := false
	if _, err := os.Stat(filepath.Join(folderPath, ".nogallery")); err == nil {
		exemptGallery = true
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("Could not test Path %s: %w", folderPath, err)
	}

	if forceGallery || (h.ScanConfig.GetCreateGalleriesFromFolders() && !exemptGallery) {
		return h.getOrCreateFolderBasedGallery(ctx, f)
	}

	return nil, nil
}

func (h *ScanHandler) getGalleryToAssociate(ctx context.Context, newImage *models.Image, f models.File) (*models.Gallery, error) {
	g, err := h.getOrCreateGallery(ctx, f)
	if err != nil {
		return nil, err
	}

	if err := newImage.LoadGalleryIDs(ctx, h.CreatorUpdater); err != nil {
		return nil, err
	}

	if g != nil && !slices.Contains(newImage.GalleryIDs.List(), g.ID) {
		return g, nil
	}

	return nil, nil
}
