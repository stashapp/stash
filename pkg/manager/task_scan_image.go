package manager

import (
	"context"
	"database/sql"
	"path/filepath"
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

	scanner := image.Scanner{
		Scanner:            image.FileScanner(&file.FSHasher{}),
		StripFileExtension: t.StripFileExtension,
		Ctx:                t.ctx,
		TxnManager:         t.TxnManager,
		Paths:              GetInstance().Paths,
		PluginCache:        instance.PluginCache,
		MutexManager:       t.mutexManager,
	}

	var err error
	if i != nil {
		i, err = scanner.ScanExisting(i, t.file)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	} else {
		i, err = scanner.ScanNew(t.file)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		if i != nil {
			if t.zipGallery != nil {
				// associate with gallery
				if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
					return gallery.AddImage(r.Gallery(), t.zipGallery.ID, i.ID)
				}); err != nil {
					logger.Error(err.Error())
					return
				}
			} else if config.GetInstance().GetCreateGalleriesFromFolders() {
				// create gallery from folder or associate with existing gallery
				logger.Infof("Associating image %s with folder gallery", i.Path)
				var galleryID int
				var isNewGallery bool
				if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
					var err error
					galleryID, isNewGallery, err = t.associateImageWithFolderGallery(i.ID, r.Gallery())
					return err
				}); err != nil {
					logger.Error(err.Error())
					return
				}

				if isNewGallery {
					GetInstance().PluginCache.ExecutePostHooks(t.ctx, galleryID, plugin.GalleryCreatePost, nil, nil)
				}
			}
		}
	}

	if i != nil {
		t.generateThumbnail(i)
	}
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
