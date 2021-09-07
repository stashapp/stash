package manager

import (
	"archive/zip"
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/remeh/sizedwaitgroup"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/utils"
)

func (t *ScanTask) scanGallery() {
	var g *models.Gallery
	images := 0
	scanImages := false

	if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		g, err = r.Gallery().FindByPath(t.FilePath)

		if g != nil && err != nil {
			images, err = r.Image().CountByGalleryID(g.ID)
			if err != nil {
				return fmt.Errorf("error getting images for zip gallery %s: %s", t.FilePath, err.Error())
			}
		}

		return err
	}); err != nil {
		logger.Error(err.Error())
		return
	}

	fileModTime, err := t.getFileModTime()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if g != nil {
		// We already have this item in the database, keep going

		// if file mod time is not set, set it now
		if !g.FileModTime.Valid {
			// we will also need to rescan the zip contents
			scanImages = true
			logger.Infof("setting file modification time on %s", t.FilePath)

			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				qb := r.Gallery()
				if _, err := gallery.UpdateFileModTime(qb, g.ID, models.NullSQLiteTimestamp{
					Timestamp: fileModTime,
					Valid:     true,
				}); err != nil {
					return err
				}

				// update our copy of the gallery
				var err error
				g, err = qb.Find(g.ID)
				return err
			}); err != nil {
				logger.Error(err.Error())
				return
			}
		}

		// if the mod time of the zip file is different than that of the associated
		// gallery, then recalculate the checksum
		modified := t.isFileModified(fileModTime, g.FileModTime)
		if modified {
			scanImages = true
			logger.Infof("%s has been updated: rescanning", t.FilePath)

			// update the checksum and the modification time
			checksum, err := t.calculateChecksum()
			if err != nil {
				logger.Error(err.Error())
				return
			}

			currentTime := time.Now()
			galleryPartial := models.GalleryPartial{
				ID:       g.ID,
				Checksum: &checksum,
				FileModTime: &models.NullSQLiteTimestamp{
					Timestamp: fileModTime,
					Valid:     true,
				},
				UpdatedAt: &models.SQLiteTimestamp{Timestamp: currentTime},
			}

			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				_, err := r.Gallery().UpdatePartial(galleryPartial)
				return err
			}); err != nil {
				logger.Error(err.Error())
				return
			}
		}

		// scan the zip files if the gallery has no images
		scanImages = scanImages || images == 0
	} else {
		// Ignore directories.
		if isDir, _ := utils.DirExists(t.FilePath); isDir {
			return
		}

		checksum, err := t.calculateChecksum()
		if err != nil {
			logger.Error(err.Error())
			return
		}

		if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
			qb := r.Gallery()
			g, _ = qb.FindByChecksum(checksum)
			if g != nil {
				exists, _ := utils.FileExists(g.Path.String)
				if !t.CaseSensitiveFs {
					// #1426 - if file exists but is a case-insensitive match for the
					// original filename, then treat it as a move
					if exists && strings.EqualFold(t.FilePath, g.Path.String) {
						exists = false
					}
				}

				if exists {
					logger.Infof("%s already exists.  Duplicate of %s ", t.FilePath, g.Path.String)
				} else {
					logger.Infof("%s already exists.  Updating path...", t.FilePath)
					g.Path = sql.NullString{
						String: t.FilePath,
						Valid:  true,
					}
					g, err = qb.Update(*g)
					if err != nil {
						return err
					}

					GetInstance().PluginCache.ExecutePostHooks(t.ctx, g.ID, plugin.GalleryUpdatePost, nil, nil)
				}
			} else {
				currentTime := time.Now()

				newGallery := models.Gallery{
					Checksum: checksum,
					Zip:      true,
					Path: sql.NullString{
						String: t.FilePath,
						Valid:  true,
					},
					FileModTime: models.NullSQLiteTimestamp{
						Timestamp: fileModTime,
						Valid:     true,
					},
					Title: sql.NullString{
						String: utils.GetNameFromPath(t.FilePath, t.StripFileExtension),
						Valid:  true,
					},
					CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
					UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
				}

				// don't create gallery if it has no images
				if countImagesInZip(t.FilePath) > 0 {
					// only warn when creating the gallery
					ok, err := utils.IsZipFileUncompressed(t.FilePath)
					if err == nil && !ok {
						logger.Warnf("%s is using above store (0) level compression.", t.FilePath)
					}

					logger.Infof("%s doesn't exist.  Creating new item...", t.FilePath)
					g, err = qb.Create(newGallery)
					if err != nil {
						return err
					}
					scanImages = true

					GetInstance().PluginCache.ExecutePostHooks(t.ctx, g.ID, plugin.GalleryCreatePost, nil, nil)
				}
			}

			return nil
		}); err != nil {
			logger.Error(err.Error())
			return
		}
	}

	if g != nil {
		if scanImages {
			t.scanZipImages(g)
		} else {
			// in case thumbnails have been deleted, regenerate them
			t.regenerateZipImages(g)
		}
	}
}

// associates a gallery to a scene with the same basename
func (t *ScanTask) associateGallery(wg *sizedwaitgroup.SizedWaitGroup) {
	if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		qb := r.Gallery()
		sqb := r.Scene()
		g, err := qb.FindByPath(t.FilePath)
		if err != nil {
			return err
		}

		if g == nil {
			// associate is run after scan is finished
			// should only happen if gallery is a directory or an io error occurs during hashing
			logger.Warnf("associate: gallery %s not found in DB", t.FilePath)
			return nil
		}

		basename := strings.TrimSuffix(t.FilePath, filepath.Ext(t.FilePath))
		var relatedFiles []string
		vExt := config.GetInstance().GetVideoExtensions()
		// make a list of media files that can be related to the gallery
		for _, ext := range vExt {
			related := basename + "." + ext
			// exclude gallery extensions from the related files
			if !isGallery(related) {
				relatedFiles = append(relatedFiles, related)
			}
		}
		for _, scenePath := range relatedFiles {
			scene, _ := sqb.FindByPath(scenePath)
			// found related Scene
			if scene != nil {
				sceneGalleries, _ := sqb.FindByGalleryID(g.ID) // check if gallery is already associated to the scene
				isAssoc := false
				for _, sg := range sceneGalleries {
					if scene.ID == sg.ID {
						isAssoc = true
						break
					}
				}
				if !isAssoc {
					logger.Infof("associate: Gallery %s is related to scene: %d", t.FilePath, scene.ID)
					if err := sqb.UpdateGalleries(scene.ID, []int{g.ID}); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
	wg.Done()
}

func (t *ScanTask) scanZipImages(zipGallery *models.Gallery) {
	err := walkGalleryZip(zipGallery.Path.String, func(file *zip.File) error {
		// copy this task and change the filename
		subTask := *t

		// filepath is the zip file and the internal file name, separated by a null byte
		subTask.FilePath = image.ZipFilename(zipGallery.Path.String, file.Name)
		subTask.zipGallery = zipGallery

		// run the subtask and wait for it to complete
		iwg := sizedwaitgroup.New(1)
		iwg.Add()
		subTask.Start(&iwg)
		return nil
	})
	if err != nil {
		logger.Warnf("failed to scan zip file images for %s: %s", zipGallery.Path.String, err.Error())
	}
}

func (t *ScanTask) regenerateZipImages(zipGallery *models.Gallery) {
	var images []*models.Image
	if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		iqb := r.Image()

		var err error
		images, err = iqb.FindByGalleryID(zipGallery.ID)
		return err
	}); err != nil {
		logger.Warnf("failed to find gallery images: %s", err.Error())
		return
	}

	for _, img := range images {
		t.generateThumbnail(img)
	}
}
