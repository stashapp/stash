package manager

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/utils"
)

func (t *ScanTask) scanImage() {
	var i *models.Image
	path := t.file.Path()

	if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		i, err = r.Image().FindByPath(path)
		return err
	}); err != nil {
		logger.Error(err.Error())
		return
	}

	scanner := file.Scanner{
		Hasher:       &file.FSHasher{},
		CalculateMD5: true,
	}

	if i != nil {
		scanned, err := scanner.ScanExisting(i, t.file)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		i, err = t.scanImageExisting(i, scanned)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	} else {
		file, err := scanner.ScanNew(t.file)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		i, err = t.scanImageNew(file)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}

	if i != nil {
		t.generateThumbnail(i)
	}
}

func (t *ScanTask) scanImageExisting(i *models.Image, scanned *file.Scanned) (retImage *models.Image, err error) {
	path := t.file.Path()
	oldChecksum := i.Checksum
	changed := false

	if scanned.ContentsChanged() {
		logger.Infof("%s has been updated: rescanning", path)

		// regenerate the file details as well
		if err := image.SetFileDetails(i); err != nil {
			return nil, err
		}

		changed = true
	} else if scanned.FileUpdated() {
		logger.Infof("Updated scene file %s", path)

		changed = true
	}

	if changed {
		i.SetFile(*scanned.New)
		i.UpdatedAt = models.SQLiteTimestamp{Timestamp: time.Now()}

		if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
			var err error
			retImage, err = r.Image().UpdateFull(*i)
			return err
		}); err != nil {
			return nil, err
		}

		// remove the old thumbnail if the checksum changed - we'll regenerate it
		if oldChecksum != scanned.New.Checksum {
			err = os.Remove(GetInstance().Paths.Generated.GetThumbnailPath(oldChecksum, models.DefaultGthumbWidth)) // remove cache dir of gallery
			if err != nil {
				logger.Errorf("Error deleting thumbnail image: %s", err)
			}
		}

		GetInstance().PluginCache.ExecutePostHooks(t.ctx, retImage.ID, plugin.ImageUpdatePost, nil, nil)
	}

	return
}

func (t *ScanTask) scanImageNew(file *models.File) (retImage *models.Image, err error) {
	path := t.file.Path()
	checksum := file.Checksum

	// check for image by checksum
	var existingImage *models.Image
	if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		existingImage, err = r.Image().FindByChecksum(checksum)
		return err
	}); err != nil {
		return nil, err
	}

	if existingImage != nil {
		exists := image.FileExists(existingImage.Path)
		if !t.CaseSensitiveFs {
			// #1426 - if file exists but is a case-insensitive match for the
			// original filename, then treat it as a move
			if exists && strings.EqualFold(path, existingImage.Path) {
				exists = false
			}
		}

		if exists {
			logger.Infof("%s already exists.  Duplicate of %s ", image.PathDisplayName(path), image.PathDisplayName(existingImage.Path))
			return nil, nil
		} else {
			logger.Infof("%s already exists.  Updating path...", image.PathDisplayName(path))
			imagePartial := models.ImagePartial{
				ID:   existingImage.ID,
				Path: &path,
			}

			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				retImage, err = r.Image().Update(imagePartial)
				return err
			}); err != nil {
				return nil, err
			}

			GetInstance().PluginCache.ExecutePostHooks(t.ctx, existingImage.ID, plugin.ImageUpdatePost, nil, nil)
		}
	} else {
		logger.Infof("%s doesn't exist.  Creating new item...", image.PathDisplayName(path))
		currentTime := time.Now()
		newImage := models.Image{
			CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		}
		newImage.SetFile(*file)
		newImage.Title.String = image.GetFilename(&newImage, t.StripFileExtension)
		newImage.Title.Valid = true

		if err := image.SetFileDetails(&newImage); err != nil {
			logger.Error(err.Error())
			return nil, err
		}

		if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
			var err error
			retImage, err = r.Image().Create(newImage)
			return err
		}); err != nil {
			return nil, err
		}

		GetInstance().PluginCache.ExecutePostHooks(t.ctx, retImage.ID, plugin.ImageCreatePost, nil, nil)
	}

	if retImage != nil {
		if t.zipGallery != nil {
			// associate with gallery
			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				return gallery.AddImage(r.Gallery(), t.zipGallery.ID, retImage.ID)
			}); err != nil {
				return nil, err
			}
		} else if config.GetInstance().GetCreateGalleriesFromFolders() {
			// create gallery from folder or associate with existing gallery
			logger.Infof("Associating image %s with folder gallery", retImage.Path)
			var galleryID int
			var isNewGallery bool
			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				var err error
				galleryID, isNewGallery, err = t.associateImageWithFolderGallery(retImage.ID, r.Gallery())
				return err
			}); err != nil {
				return nil, err
			}

			if isNewGallery {
				GetInstance().PluginCache.ExecutePostHooks(t.ctx, galleryID, plugin.GalleryCreatePost, nil, nil)
			}
		}
	}

	return
}

func (t *ScanTask) associateImageWithFolderGallery(imageID int, qb models.GalleryReaderWriter) (galleryID int, isNew bool, err error) {
	// find a gallery with the path specified
	path := filepath.Dir(t.file.Path())
	var g *models.Gallery
	g, err = qb.FindByPath(path)
	if err != nil {
		return
	}

	if g == nil {
		checksum := utils.MD5FromString(path)

		// create the gallery
		currentTime := time.Now()

		newGallery := models.Gallery{
			Checksum: checksum,
			Path: sql.NullString{
				String: path,
				Valid:  true,
			},
			CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
			Title: sql.NullString{
				String: utils.GetNameFromPath(path, false),
				Valid:  true,
			},
		}

		logger.Infof("Creating gallery for folder %s", path)
		g, err = qb.Create(newGallery)
		if err != nil {
			return 0, false, err
		}

		isNew = true
	}

	// associate image with gallery
	err = gallery.AddImage(qb, g.ID, imageID)
	galleryID = g.ID
	return
}

func (t *ScanTask) generateThumbnail(i *models.Image) {
	thumbPath := GetInstance().Paths.Generated.GetThumbnailPath(i.Checksum, models.DefaultGthumbWidth)
	exists, _ := utils.FileExists(thumbPath)
	if exists {
		return
	}

	srcImage, err := image.GetSourceImage(i)
	if err != nil {
		logger.Errorf("error reading image %s: %s", i.Path, err.Error())
		return
	}

	if image.ThumbnailNeeded(srcImage, models.DefaultGthumbWidth) {
		data, err := image.GetThumbnail(srcImage, models.DefaultGthumbWidth)
		if err != nil {
			logger.Errorf("error getting thumbnail for image %s: %s", i.Path, err.Error())
			return
		}

		err = utils.WriteFile(thumbPath, data)
		if err != nil {
			logger.Errorf("error writing thumbnail for image %s: %s", i.Path, err)
		}
	}
}
