package manager

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	fileModTime, err := image.GetFileModTime(path)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if i != nil {
		// if file mod time is not set, set it now
		if !i.FileModTime.Valid {
			logger.Infof("setting file modification time on %s", path)

			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				qb := r.Image()
				if _, err := image.UpdateFileModTime(qb, i.ID, models.NullSQLiteTimestamp{
					Timestamp: fileModTime,
					Valid:     true,
				}); err != nil {
					return err
				}

				// update our copy of the gallery
				var err error
				i, err = qb.Find(i.ID)
				return err
			}); err != nil {
				logger.Error(err.Error())
				return
			}
		}

		// if the mod time of the file is different than that of the associated
		// image, then recalculate the checksum and regenerate the thumbnail
		modified := t.isFileModified(fileModTime, i.FileModTime)
		if modified {
			i, err = t.rescanImage(i, fileModTime)
			if err != nil {
				logger.Error(err.Error())
				return
			}
		}

		// We already have this item in the database
		// check for thumbnails
		t.generateThumbnail(i)
	} else {
		// Ignore directories.
		if isDir, _ := utils.DirExists(path); isDir {
			return
		}

		var checksum string

		logger.Infof("%s not found.  Calculating checksum...", path)
		checksum, err = t.calculateImageChecksum()
		if err != nil {
			logger.Errorf("error calculating checksum for %s: %s", path, err.Error())
			return
		}

		// check for scene by checksum and oshash - MD5 should be
		// redundant, but check both
		if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			var err error
			i, err = r.Image().FindByChecksum(checksum)
			return err
		}); err != nil {
			logger.Error(err.Error())
			return
		}

		if i != nil {
			exists := image.FileExists(i.Path)
			if !t.CaseSensitiveFs {
				// #1426 - if file exists but is a case-insensitive match for the
				// original filename, then treat it as a move
				if exists && strings.EqualFold(path, i.Path) {
					exists = false
				}
			}

			if exists {
				logger.Infof("%s already exists.  Duplicate of %s ", image.PathDisplayName(path), image.PathDisplayName(i.Path))
			} else {
				logger.Infof("%s already exists.  Updating path...", image.PathDisplayName(path))
				imagePartial := models.ImagePartial{
					ID:   i.ID,
					Path: &path,
				}

				if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
					_, err := r.Image().Update(imagePartial)
					return err
				}); err != nil {
					logger.Error(err.Error())
					return
				}

				GetInstance().PluginCache.ExecutePostHooks(t.ctx, i.ID, plugin.ImageUpdatePost, nil, nil)
			}
		} else {
			logger.Infof("%s doesn't exist.  Creating new item...", image.PathDisplayName(path))
			currentTime := time.Now()
			newImage := models.Image{
				Checksum: checksum,
				Path:     path,
				FileModTime: models.NullSQLiteTimestamp{
					Timestamp: fileModTime,
					Valid:     true,
				},
				CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
				UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
			}
			newImage.Title.String = image.GetFilename(&newImage, t.StripFileExtension)
			newImage.Title.Valid = true

			if err := image.SetFileDetails(&newImage); err != nil {
				logger.Error(err.Error())
				return
			}

			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				var err error
				i, err = r.Image().Create(newImage)
				return err
			}); err != nil {
				logger.Error(err.Error())
				return
			}

			GetInstance().PluginCache.ExecutePostHooks(t.ctx, i.ID, plugin.ImageCreatePost, nil, nil)
		}

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
			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				return t.associateImageWithFolderGallery(i.ID, r.Gallery())
			}); err != nil {
				logger.Error(err.Error())
				return
			}
		}
	}

	if i != nil {
		t.generateThumbnail(i)
	}
}

func (t *ScanTask) rescanImage(i *models.Image, fileModTime time.Time) (*models.Image, error) {
	path := t.file.Path()
	logger.Infof("%s has been updated: rescanning", path)

	oldChecksum := i.Checksum

	// update the checksum and the modification time
	checksum, err := t.calculateImageChecksum()
	if err != nil {
		return nil, err
	}

	// regenerate the file details as well
	fileDetails, err := image.GetFileDetails(path)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	imagePartial := models.ImagePartial{
		ID:       i.ID,
		Checksum: &checksum,
		Width:    &fileDetails.Width,
		Height:   &fileDetails.Height,
		Size:     &fileDetails.Size,
		FileModTime: &models.NullSQLiteTimestamp{
			Timestamp: fileModTime,
			Valid:     true,
		},
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: currentTime},
	}

	var ret *models.Image
	if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		var err error
		ret, err = r.Image().Update(imagePartial)
		return err
	}); err != nil {
		return nil, err
	}

	// remove the old thumbnail if the checksum changed - we'll regenerate it
	if oldChecksum != checksum {
		err = os.Remove(GetInstance().Paths.Generated.GetThumbnailPath(oldChecksum, models.DefaultGthumbWidth)) // remove cache dir of gallery
		if err != nil {
			logger.Errorf("Error deleting thumbnail image: %s", err)
		}
	}

	GetInstance().PluginCache.ExecutePostHooks(t.ctx, ret.ID, plugin.ImageUpdatePost, nil, nil)

	return ret, nil
}

func (t *ScanTask) associateImageWithFolderGallery(imageID int, qb models.GalleryReaderWriter) error {
	// find a gallery with the path specified
	path := filepath.Dir(t.file.Path())
	g, err := qb.FindByPath(path)
	if err != nil {
		return err
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
			return err
		}
	}

	// associate image with gallery
	err = gallery.AddImage(qb, g.ID, imageID)
	return err
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
