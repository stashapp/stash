package manager

import (
	"archive/zip"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/remeh/sizedwaitgroup"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

func (t *ScanTask) scanGallery(ctx context.Context) {
	var g *models.Gallery
	path := t.file.Path()
	images := 0
	scanImages := false

	if err := t.TxnManager.WithTxn(ctx, func(ctx context.Context) error {
		var err error
		g, err = t.TxnManager.Gallery.FindByPath(ctx, path)

		if g != nil && err == nil {
			images, err = t.TxnManager.Image.CountByGalleryID(ctx, g.ID)
			if err != nil {
				return fmt.Errorf("error getting images for zip gallery %s: %s", path, err.Error())
			}
		}

		return err
	}); err != nil {
		logger.Error(err.Error())
		return
	}

	scanner := gallery.Scanner{
		Scanner:            gallery.FileScanner(&file.FSHasher{}),
		ImageExtensions:    instance.Config.GetImageExtensions(),
		StripFileExtension: t.StripFileExtension,
		CaseSensitiveFs:    t.CaseSensitiveFs,
		CreatorUpdater:     t.TxnManager.Gallery,
		Paths:              instance.Paths,
		PluginCache:        instance.PluginCache,
		MutexManager:       t.mutexManager,
	}

	var err error
	if g != nil {
		g, scanImages, err = scanner.ScanExisting(ctx, g, t.file)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		// scan the zip files if the gallery has no images
		scanImages = scanImages || images == 0
	} else {
		g, scanImages, err = scanner.ScanNew(ctx, t.file)
		if err != nil {
			logger.Error(err.Error())
		}
	}

	if g != nil {
		if scanImages {
			t.scanZipImages(ctx, g)
		} else {
			// in case thumbnails have been deleted, regenerate them
			t.regenerateZipImages(ctx, g)
		}
	}
}

// associates a gallery to a scene with the same basename
func (t *ScanTask) associateGallery(ctx context.Context, wg *sizedwaitgroup.SizedWaitGroup) {
	path := t.file.Path()
	if err := t.TxnManager.WithTxn(ctx, func(ctx context.Context) error {
		r := t.TxnManager
		qb := r.Gallery
		sqb := r.Scene
		g, err := qb.FindByPath(ctx, path)
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
			scene, _ := sqb.FindByPath(ctx, scenePath)
			// found related Scene
			if scene != nil {
				sceneGalleries, _ := sqb.FindByGalleryID(ctx, g.ID) // check if gallery is already associated to the scene
				isAssoc := false
				for _, sg := range sceneGalleries {
					if scene.ID == sg.ID {
						isAssoc = true
						break
					}
				}
				if !isAssoc {
					logger.Infof("associate: Gallery %s is related to scene: %d", path, scene.ID)
					if _, err := sqb.UpdatePartial(ctx, scene.ID, models.ScenePartial{
						GalleryIDs: &models.UpdateIDs{
							IDs:  []int{g.ID},
							Mode: models.RelationshipUpdateModeAdd,
						},
					}); err != nil {
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

func (t *ScanTask) scanZipImages(ctx context.Context, zipGallery *models.Gallery) {
	err := walkGalleryZip(*zipGallery.Path, func(f *zip.File) error {
		// copy this task and change the filename
		subTask := *t

		// filepath is the zip file and the internal file name, separated by a null byte
		subTask.file = file.ZipFile(*zipGallery.Path, f)
		subTask.zipGallery = zipGallery

		// run the subtask and wait for it to complete
		subTask.Start(ctx)
		return nil
	})
	if err != nil {
		logger.Warnf("failed to scan zip file images for %s: %s", *zipGallery.Path, err.Error())
	}
}

func (t *ScanTask) regenerateZipImages(ctx context.Context, zipGallery *models.Gallery) {
	var images []*models.Image
	if err := t.TxnManager.WithTxn(ctx, func(ctx context.Context) error {
		iqb := t.TxnManager.Image

		var err error
		images, err = iqb.FindByGalleryID(ctx, zipGallery.ID)
		return err
	}); err != nil {
		logger.Warnf("failed to find gallery images: %s", err.Error())
		return
	}

	for _, img := range images {
		t.generateThumbnail(img)
	}
}
