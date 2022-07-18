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

// const mutexType = "image"

type FinderCreatorUpdater interface {
	FindByFileID(ctx context.Context, fileID file.ID) ([]*models.Image, error)
	FindByFingerprints(ctx context.Context, fp []file.Fingerprint) ([]*models.Image, error)
	Create(ctx context.Context, newImage *models.ImageCreateInput) error
	Update(ctx context.Context, updatedImage *models.Image) error
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
			CreatedAt: now,
			UpdatedAt: now,
		}

		// if the file is in a zip, then associate it with the gallery
		if imageFile.ZipFileID != nil {
			g, err := h.GalleryFinder.FindByFileID(ctx, *imageFile.ZipFileID)
			if err != nil {
				return fmt.Errorf("finding gallery for zip file id %d: %w", *imageFile.ZipFileID, err)
			}

			for _, gg := range g {
				newImage.GalleryIDs = append(newImage.GalleryIDs, gg.ID)
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
		found := false
		for _, sf := range i.Files {
			if sf.ID == f.Base().ID {
				found = true
				break
			}
		}

		if !found {
			logger.Infof("Adding %s to image %s", f.Path, i.GetTitle())
			i.Files = append(i.Files, f)

			// associate with folder-based gallery if applicable
			if h.ScanConfig.GetCreateGalleriesFromFolders() {
				if err := h.associateFolderBasedGallery(ctx, i, f); err != nil {
					return err
				}
			}

			if err := h.CreatorUpdater.Update(ctx, i); err != nil {
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

	if g != nil && !intslice.IntInclude(newImage.GalleryIDs, g.ID) {
		newImage.GalleryIDs = append(newImage.GalleryIDs, g.ID)
		logger.Infof("Adding %s to folder-based gallery %s", f.Base().Path, g.Path())
	}

	return nil
}

// type Scanner struct {
// 	file.Scanner

// 	StripFileExtension bool

// 	CaseSensitiveFs bool
// 	TxnManager      txn.Manager
// 	CreatorUpdater  FinderCreatorUpdater
// 	Paths           *paths.Paths
// 	PluginCache     *plugin.Cache
// 	MutexManager    *utils.MutexManager
// }

// func FileScanner(hasher file.Hasher) file.Scanner {
// 	return file.Scanner{
// 		Hasher:       hasher,
// 		CalculateMD5: true,
// 	}
// }

// func (scanner *Scanner) ScanExisting(ctx context.Context, existing file.FileBased, file file.SourceFile) (retImage *models.Image, err error) {
// 	scanned, err := scanner.Scanner.ScanExisting(existing, file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	i := existing.(*models.Image)

// 	path := scanned.New.Path
// 	oldChecksum := i.Checksum
// 	changed := false

// 	if scanned.ContentsChanged() {
// 		logger.Infof("%s has been updated: rescanning", path)

// 		// regenerate the file details as well
// 		if err := SetFileDetails(i); err != nil {
// 			return nil, err
// 		}

// 		changed = true
// 	} else if scanned.FileUpdated() {
// 		logger.Infof("Updated image file %s", path)

// 		changed = true
// 	}

// 	if changed {
// 		i.SetFile(*scanned.New)
// 		i.UpdatedAt = time.Now()

// 		// we are operating on a checksum now, so grab a mutex on the checksum
// 		done := make(chan struct{})
// 		scanner.MutexManager.Claim(mutexType, scanned.New.Checksum, done)

// 		if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
// 			// free the mutex once transaction is complete
// 			defer close(done)
// 			var err error

// 			// ensure no clashes of hashes
// 			if scanned.New.Checksum != "" && scanned.Old.Checksum != scanned.New.Checksum {
// 				dupe, _ := scanner.CreatorUpdater.FindByChecksum(ctx, i.Checksum)
// 				if dupe != nil {
// 					return fmt.Errorf("MD5 for file %s is the same as that of %s", path, dupe.Path)
// 				}
// 			}

// 			err = scanner.CreatorUpdater.Update(ctx, i)
// 			return err
// 		}); err != nil {
// 			return nil, err
// 		}

// 		retImage = i

// 		// remove the old thumbnail if the checksum changed - we'll regenerate it
// 		if oldChecksum != scanned.New.Checksum {
// 			// remove cache dir of gallery
// 			err = os.Remove(scanner.Paths.Generated.GetThumbnailPath(oldChecksum, models.DefaultGthumbWidth))
// 			if err != nil {
// 				logger.Errorf("Error deleting thumbnail image: %s", err)
// 			}
// 		}

// 		scanner.PluginCache.ExecutePostHooks(ctx, retImage.ID, plugin.ImageUpdatePost, nil, nil)
// 	}

// 	return
// }

// func (scanner *Scanner) ScanNew(ctx context.Context, f file.SourceFile) (retImage *models.Image, err error) {
// 	scanned, err := scanner.Scanner.ScanNew(f)
// 	if err != nil {
// 		return nil, err
// 	}

// 	path := f.Path()
// 	checksum := scanned.Checksum

// 	// grab a mutex on the checksum
// 	done := make(chan struct{})
// 	scanner.MutexManager.Claim(mutexType, checksum, done)
// 	defer close(done)

// 	// check for image by checksum
// 	var existingImage *models.Image
// 	if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
// 		var err error
// 		existingImage, err = scanner.CreatorUpdater.FindByChecksum(ctx, checksum)
// 		return err
// 	}); err != nil {
// 		return nil, err
// 	}

// 	pathDisplayName := file.ZipPathDisplayName(path)

// 	if existingImage != nil {
// 		exists := FileExists(existingImage.Path)
// 		if !scanner.CaseSensitiveFs {
// 			// #1426 - if file exists but is a case-insensitive match for the
// 			// original filename, then treat it as a move
// 			if exists && strings.EqualFold(path, existingImage.Path) {
// 				exists = false
// 			}
// 		}

// 		if exists {
// 			logger.Infof("%s already exists. Duplicate of %s ", pathDisplayName, file.ZipPathDisplayName(existingImage.Path))
// 			return nil, nil
// 		} else {
// 			logger.Infof("%s already exists. Updating path...", pathDisplayName)

// 			existingImage.Path = path
// 			if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
// 				return scanner.CreatorUpdater.Update(ctx, existingImage)
// 			}); err != nil {
// 				return nil, err
// 			}

// 			retImage = existingImage

// 			scanner.PluginCache.ExecutePostHooks(ctx, existingImage.ID, plugin.ImageUpdatePost, nil, nil)
// 		}
// 	} else {
// 		logger.Infof("%s doesn't exist. Creating new item...", pathDisplayName)
// 		currentTime := time.Now()
// 		newImage := &models.Image{
// 			CreatedAt: currentTime,
// 			UpdatedAt: currentTime,
// 		}
// 		newImage.SetFile(*scanned)
// 		fn := GetFilename(newImage, scanner.StripFileExtension)
// 		newImage.Title = fn

// 		if err := SetFileDetails(newImage); err != nil {
// 			logger.Error(err.Error())
// 			return nil, err
// 		}

// 		if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
// 			return scanner.CreatorUpdater.Create(ctx, newImage)
// 		}); err != nil {
// 			return nil, err
// 		}

// 		retImage = newImage

// 		scanner.PluginCache.ExecutePostHooks(ctx, retImage.ID, plugin.ImageCreatePost, nil, nil)
// 	}

// 	return
// }
