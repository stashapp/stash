package image

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

const mutexType = "image"

type FinderCreatorUpdater interface {
	FindByChecksum(ctx context.Context, checksum string) (*models.Image, error)
	Create(ctx context.Context, newImage models.Image) (*models.Image, error)
	UpdateFull(ctx context.Context, updatedImage models.Image) (*models.Image, error)
	Update(ctx context.Context, updatedImage models.ImagePartial) (*models.Image, error)
}

type Scanner struct {
	file.Scanner

	StripFileExtension bool

	CaseSensitiveFs bool
	TxnManager      txn.Manager
	CreatorUpdater  FinderCreatorUpdater
	Paths           *paths.Paths
	PluginCache     *plugin.Cache
	MutexManager    *utils.MutexManager
}

func FileScanner(hasher file.Hasher) file.Scanner {
	return file.Scanner{
		Hasher:       hasher,
		CalculateMD5: true,
	}
}

func (scanner *Scanner) ScanExisting(ctx context.Context, existing file.FileBased, file file.SourceFile) (retImage *models.Image, err error) {
	scanned, err := scanner.Scanner.ScanExisting(existing, file)
	if err != nil {
		return nil, err
	}

	i := existing.(*models.Image)

	path := scanned.New.Path
	oldChecksum := i.Checksum
	changed := false

	if scanned.ContentsChanged() {
		logger.Infof("%s has been updated: rescanning", path)

		// regenerate the file details as well
		if err := SetFileDetails(i); err != nil {
			return nil, err
		}

		changed = true
	} else if scanned.FileUpdated() {
		logger.Infof("Updated image file %s", path)

		changed = true
	}

	if changed {
		i.SetFile(*scanned.New)
		i.UpdatedAt = models.SQLiteTimestamp{Timestamp: time.Now()}

		// we are operating on a checksum now, so grab a mutex on the checksum
		done := make(chan struct{})
		scanner.MutexManager.Claim(mutexType, scanned.New.Checksum, done)

		if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
			// free the mutex once transaction is complete
			defer close(done)
			var err error

			// ensure no clashes of hashes
			if scanned.New.Checksum != "" && scanned.Old.Checksum != scanned.New.Checksum {
				dupe, _ := scanner.CreatorUpdater.FindByChecksum(ctx, i.Checksum)
				if dupe != nil {
					return fmt.Errorf("MD5 for file %s is the same as that of %s", path, dupe.Path)
				}
			}

			retImage, err = scanner.CreatorUpdater.UpdateFull(ctx, *i)
			return err
		}); err != nil {
			return nil, err
		}

		// remove the old thumbnail if the checksum changed - we'll regenerate it
		if oldChecksum != scanned.New.Checksum {
			// remove cache dir of gallery
			err = os.Remove(scanner.Paths.Generated.GetThumbnailPath(oldChecksum, models.DefaultGthumbWidth))
			if err != nil {
				logger.Errorf("Error deleting thumbnail image: %s", err)
			}
		}

		scanner.PluginCache.ExecutePostHooks(ctx, retImage.ID, plugin.ImageUpdatePost, nil, nil)
	}

	return
}

func (scanner *Scanner) ScanNew(ctx context.Context, f file.SourceFile) (retImage *models.Image, err error) {
	scanned, err := scanner.Scanner.ScanNew(f)
	if err != nil {
		return nil, err
	}

	path := f.Path()
	checksum := scanned.Checksum

	// grab a mutex on the checksum
	done := make(chan struct{})
	scanner.MutexManager.Claim(mutexType, checksum, done)
	defer close(done)

	// check for image by checksum
	var existingImage *models.Image
	if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
		var err error
		existingImage, err = scanner.CreatorUpdater.FindByChecksum(ctx, checksum)
		return err
	}); err != nil {
		return nil, err
	}

	pathDisplayName := file.ZipPathDisplayName(path)

	if existingImage != nil {
		exists := FileExists(existingImage.Path)
		if !scanner.CaseSensitiveFs {
			// #1426 - if file exists but is a case-insensitive match for the
			// original filename, then treat it as a move
			if exists && strings.EqualFold(path, existingImage.Path) {
				exists = false
			}
		}

		if exists {
			logger.Infof("%s already exists. Duplicate of %s ", pathDisplayName, file.ZipPathDisplayName(existingImage.Path))
			return nil, nil
		} else {
			logger.Infof("%s already exists. Updating path...", pathDisplayName)
			imagePartial := models.ImagePartial{
				ID:   existingImage.ID,
				Path: &path,
			}

			if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
				retImage, err = scanner.CreatorUpdater.Update(ctx, imagePartial)
				return err
			}); err != nil {
				return nil, err
			}

			scanner.PluginCache.ExecutePostHooks(ctx, existingImage.ID, plugin.ImageUpdatePost, nil, nil)
		}
	} else {
		logger.Infof("%s doesn't exist. Creating new item...", pathDisplayName)
		currentTime := time.Now()
		newImage := models.Image{
			CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		}
		newImage.SetFile(*scanned)
		newImage.Title.String = GetFilename(&newImage, scanner.StripFileExtension)
		newImage.Title.Valid = true

		if err := SetFileDetails(&newImage); err != nil {
			logger.Error(err.Error())
			return nil, err
		}

		if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
			var err error
			retImage, err = scanner.CreatorUpdater.Create(ctx, newImage)
			return err
		}); err != nil {
			return nil, err
		}

		scanner.PluginCache.ExecutePostHooks(ctx, retImage.ID, plugin.ImageCreatePost, nil, nil)
	}

	return
}
