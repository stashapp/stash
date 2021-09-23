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
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/utils"
)

func (t *ScanTask) scanGallery() {
	var g *models.Gallery
	path := t.file.Path()
	images := 0
	scanImages := false

	if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		g, err = r.Gallery().FindByPath(path)

		if g != nil && err != nil {
			images, err = r.Image().CountByGalleryID(g.ID)
			if err != nil {
				return fmt.Errorf("error getting images for zip gallery %s: %s", path, err.Error())
			}
		}

		return err
	}); err != nil {
		logger.Error(err.Error())
		return
	}

	scanner := file.Scanner{
		Hasher:       &file.FSHasher{},
		CalculateMD5: true,
	}

	if g != nil {
		scanned, err := scanner.ScanExisting(g, t.file)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		g, scanImages, err = t.scanGalleryExisting(g, scanned)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		// scan the zip files if the gallery has no images
		scanImages = scanImages || images == 0
	} else {
		file, err := scanner.ScanNew(t.file)
		if err != nil {
			logger.Error(err.Error())
		}

		g, scanImages, err = t.scanGalleryNew(file)
		if err != nil {
			logger.Error(err.Error())
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

func (t *ScanTask) scanGalleryExisting(g *models.Gallery, scanned *file.Scanned) (ret *models.Gallery, scanImages bool, err error) {
	path := t.file.Path()
	changed := false

	if scanned.ContentsChanged() {
		logger.Infof("%s has been updated: rescanning", path)

		// TODO - do stuff here?

		g.SetFile(*scanned.New)
		changed = true
	} else if scanned.FileUpdated() {
		logger.Infof("Updated gallery file %s", path)

		g.SetFile(*scanned.New)
		changed = true
	}

	if changed {
		scanImages = true
		logger.Infof("%s has been updated: rescanning", path)

		g.UpdatedAt = models.SQLiteTimestamp{Timestamp: time.Now()}

		if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {

			// TODO - ensure no clashes of hashes

			_, err = r.Gallery().Update(*g)
			return err
		}); err != nil {
			return nil, false, err
		}

		GetInstance().PluginCache.ExecutePostHooks(t.ctx, g.ID, plugin.GalleryUpdatePost, nil, nil)
	}

	return
}

func (t *ScanTask) scanGalleryNew(file *models.File) (g *models.Gallery, scanImages bool, err error) {
	path := t.file.Path()
	isNewGallery := false
	isUpdatedGallery := false

	if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		qb := r.Gallery()

		checksum := file.Checksum

		g, _ = qb.FindByChecksum(checksum)
		if g != nil {
			exists, _ := utils.FileExists(g.Path.String)
			if !t.CaseSensitiveFs {
				// #1426 - if file exists but is a case-insensitive match for the
				// original filename, then treat it as a move
				if exists && strings.EqualFold(path, g.Path.String) {
					exists = false
				}
			}

			if exists {
				logger.Infof("%s already exists.  Duplicate of %s ", path, g.Path.String)
			} else {
				logger.Infof("%s already exists.  Updating path...", path)
				g.Path = sql.NullString{
					String: path,
					Valid:  true,
				}
				g, err = qb.Update(*g)
				if err != nil {
					return err
				}

				isUpdatedGallery = true
			}
		} else {
			// don't create gallery if it has no images
			if countImagesInZip(path) > 0 {
				currentTime := time.Now()

				newGallery := models.Gallery{
					Zip: true,
					Title: sql.NullString{
						String: utils.GetNameFromPath(path, t.StripFileExtension),
						Valid:  true,
					},
					CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
					UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
				}

				g.SetFile(*file)

				// only warn when creating the gallery
				ok, err := utils.IsZipFileUncompressed(path)
				if err == nil && !ok {
					logger.Warnf("%s is using above store (0) level compression.", path)
				}

				logger.Infof("%s doesn't exist.  Creating new item...", path)
				g, err = qb.Create(newGallery)
				if err != nil {
					return err
				}

				scanImages = true
				isNewGallery = true
			}
		}

		return nil
	}); err != nil {
		return nil, false, err
	}

	if isNewGallery {
		GetInstance().PluginCache.ExecutePostHooks(t.ctx, g.ID, plugin.GalleryCreatePost, nil, nil)
	} else if isUpdatedGallery {
		GetInstance().PluginCache.ExecutePostHooks(t.ctx, g.ID, plugin.GalleryUpdatePost, nil, nil)
	}

	scanImages = isNewGallery
	return
}

// associates a gallery to a scene with the same basename
func (t *ScanTask) associateGallery(wg *sizedwaitgroup.SizedWaitGroup) {
	path := t.file.Path()
	if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		qb := r.Gallery()
		sqb := r.Scene()
		g, err := qb.FindByPath(path)
		if err != nil {
			return err
		}

		if g == nil {
			// associate is run after scan is finished
			// should only happen if gallery is a directory or an io error occurs during hashing
			logger.Warnf("associate: gallery %s not found in DB", path)
			return nil
		}

		basename := strings.TrimSuffix(path, filepath.Ext(path))
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
					logger.Infof("associate: Gallery %s is related to scene: %d", path, scene.ID)
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
	err := walkGalleryZip(zipGallery.Path.String, func(f *zip.File) error {
		// copy this task and change the filename
		subTask := *t

		// filepath is the zip file and the internal file name, separated by a null byte
		subTask.file = file.ZipFile(zipGallery.Path.String, f)
		subTask.zipGallery = zipGallery

		// run the subtask and wait for it to complete
		subTask.Start()
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
