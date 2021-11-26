package image

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/utils"
)

type Scanner struct {
	file.Scanner

	StripFileExtension bool

	Ctx             context.Context
	CaseSensitiveFs bool
	TxnManager      models.TransactionManager
	Paths           *paths.Paths
	PluginCache     *plugin.Cache

	Image *models.Image
	IsNew bool
}

func FileScanner(hasher file.Hasher, statter file.Statter) file.Scanner {
	return file.Scanner{
		Hasher:       hasher,
		Statter:      statter,
		CalculateMD5: true,
		Done:         make(chan struct{}),
	}
}

func (scanner *Scanner) PostScan(scanned file.Scanned) error {
	if scanned.Old != nil {
		// should be an existing image
		var image *models.Image
		if err := scanner.TxnManager.WithReadTxn(scanner.Ctx, func(r models.ReaderRepository) error {
			images, err := r.Image().FindByFileID(scanned.Old.ID)
			if err != nil {
				return err
			}

			// assume only one scene for now
			if len(images) > 0 {
				image = images[0]
			}
			return err
		}); err != nil {
			logger.Error(err.Error())
			return nil
		}

		if image != nil {
			return scanner.ScanExisting(image, scanned)
		}

		// we shouldn't be able to have an existing file without an image, but
		// assuming that it's happened, treat it as a new image
	}

	// assume a new file/scene
	return scanner.ScanNew(scanned.New)
}

func (scanner *Scanner) GenerateMetadata(dest *models.File, src file.SourceFile) error {
	f, err := src.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	config, _, err := DecodeSourceImage(f)
	if err == nil {
		dest.Width = sql.NullInt64{
			Int64: int64(config.Width),
			Valid: true,
		}
		dest.Height = sql.NullInt64{
			Int64: int64(config.Height),
			Valid: true,
		}
	}

	return err
}

func (scanner *Scanner) ScanExisting(i *models.Image, scanned file.Scanned) error {
	path := scanned.New.Path
	oldChecksum := i.Checksum
	changed := false

	if scanned.ContentsChanged() {
		logger.Infof("%s has been updated: rescanning", path)

		changed = true
	} else if scanned.FileUpdated() {
		logger.Infof("Updated image file %s", path)

		changed = true
	}

	if changed {
		i.SetFile(*scanned.New)
		i.UpdatedAt = models.SQLiteTimestamp{Timestamp: time.Now()}

		var retImage *models.Image
		if err := scanner.TxnManager.WithTxn(scanner.Ctx, func(r models.Repository) error {
			if err := scanner.Scanner.ApplyChanges(r.File(), &scanned); err != nil {
				return err
			}

			var err error

			// ensure no clashes of hashes
			if scanned.New.Checksum != "" && scanned.Old.Checksum != scanned.New.Checksum {
				dupe, _ := r.Image().FindByChecksum(i.Checksum)
				if dupe != nil {
					return fmt.Errorf("MD5 for file %s is the same as that of %s", path, dupe.Path)
				}
			}

			retImage, err = r.Image().UpdateFull(*i)
			return err
		}); err != nil {
			return err
		}

		// remove the old thumbnail if the checksum changed - we'll regenerate it
		if oldChecksum != scanned.New.Checksum {
			// remove cache dir of gallery
			if err := os.Remove(scanner.Paths.Generated.GetThumbnailPath(oldChecksum, models.DefaultGthumbWidth)); err != nil {
				logger.Errorf("Error deleting thumbnail image: %s", err)
			}
		}

		scanner.PluginCache.ExecutePostHooks(scanner.Ctx, retImage.ID, plugin.ImageUpdatePost, nil, nil)
		scanner.Image = retImage
	}

	return nil
}

func (scanner *Scanner) ScanNew(f *models.File) error {
	path := f.Path
	checksum := f.Checksum

	// check for image by checksum
	var existingImage *models.Image
	if err := scanner.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		existingImage, err = r.Image().FindByChecksum(checksum)
		return err
	}); err != nil {
		return err
	}

	pathDisplayName := file.ZipPathDisplayName(path)

	if existingImage != nil {
		logger.Infof("%s already exists. Duplicate of %s ", pathDisplayName, file.ZipPathDisplayName(existingImage.Path))

		if err := scanner.TxnManager.WithTxn(scanner.Ctx, func(r models.Repository) error {
			if err := scanner.Scanner.ApplyChanges(r.File(), &file.Scanned{
				New: f,
			}); err != nil {
				return err
			}

			// link image to file
			return addFile(r.Image(), existingImage, f)
		}); err != nil {
			return err
		}

		scanner.Image = existingImage
		scanner.PluginCache.ExecutePostHooks(scanner.Ctx, existingImage.ID, plugin.ImageUpdatePost, nil, nil)
	} else {
		logger.Infof("%s doesn't exist. Creating new item...", pathDisplayName)
		currentTime := time.Now()
		newImage := models.Image{
			CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		}
		newImage.SetFile(*f)
		newImage.Title.String = GetFilename(&newImage, scanner.StripFileExtension)
		newImage.Title.Valid = true

		var retImage *models.Image

		if err := scanner.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
			if err := scanner.Scanner.ApplyChanges(r.File(), &file.Scanned{
				New: f,
			}); err != nil {
				return err
			}

			var err error
			retImage, err = r.Image().Create(newImage)
			if err != nil {
				return err
			}

			// link image to file
			return addFile(r.Image(), retImage, f)
		}); err != nil {
			return err
		}

		scanner.Image = retImage
		scanner.IsNew = true
		scanner.PluginCache.ExecutePostHooks(scanner.Ctx, retImage.ID, plugin.ImageCreatePost, nil, nil)
	}

	return nil
}

func addFile(rw models.ImageReaderWriter, i *models.Image, f *models.File) error {
	ids, err := rw.GetFileIDs(i.ID)
	if err != nil {
		return err
	}

	ids = utils.IntAppendUnique(ids, f.ID)
	return rw.UpdateFiles(i.ID, ids)
}
