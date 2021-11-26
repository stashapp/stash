package manager

import (
	"context"
	"database/sql"
	"errors"
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

func (t *ScanTask) postScanImage(scanner *image.Scanner) {
	if scanner.Image != nil {
		i := scanner.Image

		if scanner.IsNew {
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
	if !t.GenerateThumbnails {
		return
	}

	thumbPath := GetInstance().Paths.Generated.GetThumbnailPath(i.Checksum, models.DefaultGthumbWidth)
	exists, _ := utils.FileExists(thumbPath)
	if exists {
		return
	}

	var f []*models.File
	if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		fileIDs, err := r.Image().GetFileIDs(i.ID)
		if err != nil {
			return err
		}

		if len(fileIDs) == 0 {
			return nil
		}

		f, err = r.File().Find(fileIDs[0:1])
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		logger.Errorf("error getting files for image %q: %v", i.Path, err)
		return
	}

	if len(f) == 0 {
		logger.Warnf("no files found for image %q", i.Path)
		return
	}

	reader, err := file.Open(f[0])
	if err != nil {
		logger.Errorf("error reading image %s: %s", i.Path, err.Error())
		return
	}
	defer reader.Close()

	encoder := image.NewThumbnailEncoder(instance.FFMPEG)
	data, err := encoder.GetThumbnail(reader, models.DefaultGthumbWidth)
	if err != nil {
		// don't log if image too small
		if !errors.Is(err, image.ErrTooSmall) {
			logger.Errorf("error getting thumbnail for image %s: %s", i.Path, err.Error())
		}
		return
	}

	err = utils.WriteFile(thumbPath, data)
	if err != nil {
		logger.Errorf("error writing thumbnail for image %s: %s", i.Path, err)
	}
}
