package manager

import (
	"context"
	"database/sql"
	"errors"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
)

func (t *ScanTask) scanImage(ctx context.Context) {
	var i *models.Image
	path := t.file.Path()

	if err := t.TxnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
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
		TxnManager:         t.TxnManager,
		Paths:              GetInstance().Paths,
		PluginCache:        instance.PluginCache,
		MutexManager:       t.mutexManager,
	}

	var err error
	if i != nil {
		i, err = scanner.ScanExisting(ctx, i, t.file)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	} else {
		i, err = scanner.ScanNew(ctx, t.file)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		if i != nil {
			if t.zipGallery != nil {
				// associate with gallery
				if err := t.TxnManager.WithTxn(ctx, func(r models.Repository) error {
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
				if err := t.TxnManager.WithTxn(ctx, func(r models.Repository) error {
					var err error
					galleryID, isNewGallery, err = t.associateImageWithFolderGallery(i.ID, r.Gallery())
					return err
				}); err != nil {
					logger.Error(err.Error())
					return
				}

				if isNewGallery {
					GetInstance().PluginCache.ExecutePostHooks(ctx, galleryID, plugin.GalleryCreatePost, nil, nil)
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
		checksum := md5.FromString(path)

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
				String: fsutil.GetNameFromPath(path, false),
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
	exists, _ := fsutil.FileExists(thumbPath)
	if exists {
		return
	}

	config, _, err := image.DecodeSourceImage(i)
	if err != nil {
		logger.Errorf("error reading image %s: %s", i.Path, err.Error())
		return
	}

	if config.Height > models.DefaultGthumbWidth || config.Width > models.DefaultGthumbWidth {
		encoder := image.NewThumbnailEncoder(instance.FFMPEG)
		data, err := encoder.GetThumbnail(i, models.DefaultGthumbWidth)

		if err != nil {
			// don't log for animated images
			if !errors.Is(err, image.ErrNotSupportedForThumbnail) {
				logger.Errorf("error getting thumbnail for image %s: %s", i.Path, err.Error())

				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					logger.Errorf("stderr: %s", string(exitErr.Stderr))
				}
			}
			return
		}

		err = fsutil.WriteFile(thumbPath, data)
		if err != nil {
			logger.Errorf("error writing thumbnail for image %s: %s", i.Path, err)
		}
	}
}
