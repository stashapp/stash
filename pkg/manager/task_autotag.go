package manager

import (
	"context"
	"path/filepath"
	"strings"
	"sync"

	"github.com/stashapp/stash/pkg/autotag"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type autoTagFilesTask struct {
	paths      []string
	performers bool
	studios    bool
	tags       bool

	txnManager models.TransactionManager
	status     *TaskStatus
}

func (t *autoTagFilesTask) makeSceneFilter() *models.SceneFilterType {
	ret := &models.SceneFilterType{}
	or := ret
	sep := string(filepath.Separator)

	for _, p := range t.paths {
		if !strings.HasSuffix(p, sep) {
			p = p + sep
		}

		if ret.Path == nil {
			or = ret
		} else {
			newOr := &models.SceneFilterType{}
			or.Or = newOr
			or = newOr
		}

		or.Path = &models.StringCriterionInput{
			Modifier: models.CriterionModifierEquals,
			Value:    p + "%",
		}
	}

	organized := false
	ret.Organized = &organized

	return ret
}

func (t *autoTagFilesTask) makeImageFilter() *models.ImageFilterType {
	ret := &models.ImageFilterType{}
	or := ret
	sep := string(filepath.Separator)

	for _, p := range t.paths {
		if !strings.HasSuffix(p, sep) {
			p = p + sep
		}

		if ret.Path == nil {
			or = ret
		} else {
			newOr := &models.ImageFilterType{}
			or.Or = newOr
			or = newOr
		}

		or.Path = &models.StringCriterionInput{
			Modifier: models.CriterionModifierEquals,
			Value:    p + "%",
		}
	}

	organized := false
	ret.Organized = &organized

	return ret
}

func (t *autoTagFilesTask) makeGalleryFilter() *models.GalleryFilterType {
	ret := &models.GalleryFilterType{}
	or := ret
	sep := string(filepath.Separator)

	for _, p := range t.paths {
		if !strings.HasSuffix(p, sep) {
			p = p + sep
		}

		if ret.Path == nil {
			or = ret
		} else {
			newOr := &models.GalleryFilterType{}
			or.Or = newOr
			or = newOr
		}

		or.Path = &models.StringCriterionInput{
			Modifier: models.CriterionModifierEquals,
			Value:    p + "%",
		}
	}

	organized := false
	ret.Organized = &organized
	zip := true
	ret.IsZip = &zip

	return ret
}

func (t *autoTagFilesTask) getCount(r models.ReaderRepository) (int, error) {
	pp := 0
	findFilter := &models.FindFilterType{
		PerPage: &pp,
	}

	_, sceneCount, err := r.Scene().Query(t.makeSceneFilter(), findFilter)
	if err != nil {
		return 0, err
	}

	_, imageCount, err := r.Image().Query(t.makeImageFilter(), findFilter)
	if err != nil {
		return 0, err
	}

	_, galleryCount, err := r.Gallery().Query(t.makeGalleryFilter(), findFilter)
	if err != nil {
		return 0, err
	}

	return sceneCount + imageCount + galleryCount, nil
}

func (t *autoTagFilesTask) batchFindFilter(batchSize int) *models.FindFilterType {
	page := 1
	return &models.FindFilterType{
		PerPage: &batchSize,
		Page:    &page,
	}
}

func (t *autoTagFilesTask) processScenes(r models.ReaderRepository) error {
	if t.status.stopping {
		return nil
	}

	batchSize := 1000

	findFilter := t.batchFindFilter(batchSize)
	sceneFilter := t.makeSceneFilter()

	more := true
	for more {
		scenes, _, err := r.Scene().Query(sceneFilter, findFilter)
		if err != nil {
			return err
		}

		for _, ss := range scenes {
			if t.status.stopping {
				return nil
			}

			tt := autoTagSceneTask{
				txnManager: t.txnManager,
				scene:      ss,
				performers: t.performers,
				studios:    t.studios,
				tags:       t.tags,
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go tt.Start(&wg)
			wg.Wait()

			t.status.incrementProgress()
		}

		if len(scenes) != batchSize {
			more = false
		} else {
			*findFilter.Page++
		}
	}

	return nil
}

func (t *autoTagFilesTask) processImages(r models.ReaderRepository) error {
	if t.status.stopping {
		return nil
	}

	batchSize := 1000

	findFilter := t.batchFindFilter(batchSize)
	imageFilter := t.makeImageFilter()

	more := true
	for more {
		images, _, err := r.Image().Query(imageFilter, findFilter)
		if err != nil {
			return err
		}

		for _, ss := range images {
			if t.status.stopping {
				return nil
			}

			tt := autoTagImageTask{
				txnManager: t.txnManager,
				image:      ss,
				performers: t.performers,
				studios:    t.studios,
				tags:       t.tags,
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go tt.Start(&wg)
			wg.Wait()

			t.status.incrementProgress()
		}

		if len(images) != batchSize {
			more = false
		} else {
			*findFilter.Page++
		}
	}

	return nil
}

func (t *autoTagFilesTask) processGalleries(r models.ReaderRepository) error {
	if t.status.stopping {
		return nil
	}

	batchSize := 1000

	findFilter := t.batchFindFilter(batchSize)
	galleryFilter := t.makeGalleryFilter()

	more := true
	for more {
		galleries, _, err := r.Gallery().Query(galleryFilter, findFilter)
		if err != nil {
			return err
		}

		for _, ss := range galleries {
			if t.status.stopping {
				return nil
			}

			tt := autoTagGalleryTask{
				txnManager: t.txnManager,
				gallery:    ss,
				performers: t.performers,
				studios:    t.studios,
				tags:       t.tags,
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go tt.Start(&wg)
			wg.Wait()

			t.status.incrementProgress()
		}

		if len(galleries) != batchSize {
			more = false
		} else {
			*findFilter.Page++
		}
	}

	return nil
}

func (t *autoTagFilesTask) process() {
	if err := t.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		total, err := t.getCount(r)
		if err != nil {
			return err
		}

		t.status.total = total

		logger.Infof("Starting autotag of %d files", total)

		if err := t.processScenes(r); err != nil {
			return err
		}

		if err := t.processImages(r); err != nil {
			return err
		}

		if err := t.processGalleries(r); err != nil {
			return err
		}

		if t.status.stopping {
			logger.Info("Stopping due to user request")
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}

	logger.Info("Finished autotag")
}

type autoTagSceneTask struct {
	txnManager models.TransactionManager
	scene      *models.Scene

	performers bool
	studios    bool
	tags       bool
}

func (t *autoTagSceneTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	if err := t.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		if t.performers {
			if err := autotag.ScenePerformers(t.scene, r.Scene(), r.Performer()); err != nil {
				return err
			}
		}
		if t.studios {
			if err := autotag.SceneStudios(t.scene, r.Scene(), r.Studio()); err != nil {
				return err
			}
		}
		if t.tags {
			if err := autotag.SceneTags(t.scene, r.Scene(), r.Tag()); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}

type autoTagImageTask struct {
	txnManager models.TransactionManager
	image      *models.Image

	performers bool
	studios    bool
	tags       bool
}

func (t *autoTagImageTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	if err := t.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		if t.performers {
			if err := autotag.ImagePerformers(t.image, r.Image(), r.Performer()); err != nil {
				return err
			}
		}
		if t.studios {
			if err := autotag.ImageStudios(t.image, r.Image(), r.Studio()); err != nil {
				return err
			}
		}
		if t.tags {
			if err := autotag.ImageTags(t.image, r.Image(), r.Tag()); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}

type autoTagGalleryTask struct {
	txnManager models.TransactionManager
	gallery    *models.Gallery

	performers bool
	studios    bool
	tags       bool
}

func (t *autoTagGalleryTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	if err := t.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		if t.performers {
			if err := autotag.GalleryPerformers(t.gallery, r.Gallery(), r.Performer()); err != nil {
				return err
			}
		}
		if t.studios {
			if err := autotag.GalleryStudios(t.gallery, r.Gallery(), r.Studio()); err != nil {
				return err
			}
		}
		if t.tags {
			if err := autotag.GalleryTags(t.gallery, r.Gallery(), r.Tag()); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}
