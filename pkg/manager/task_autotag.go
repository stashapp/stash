package manager

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/stashapp/stash/pkg/autotag"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type autoTagJob struct {
	txnManager models.TransactionManager
	input      models.AutoTagMetadataInput
}

func (j *autoTagJob) Execute(ctx context.Context, progress *job.Progress) {
	input := j.input
	if j.isFileBasedAutoTag(input) {
		// doing file-based auto-tag
		j.autoTagFiles(ctx, progress, input.Paths, len(input.Performers) > 0, len(input.Studios) > 0, len(input.Tags) > 0)
	} else {
		// doing specific performer/studio/tag auto-tag
		j.autoTagSpecific(ctx, progress)
	}
}

func (j *autoTagJob) isFileBasedAutoTag(input models.AutoTagMetadataInput) bool {
	const wildcard = "*"
	performerIds := input.Performers
	studioIds := input.Studios
	tagIds := input.Tags

	return (len(performerIds) == 0 || performerIds[0] == wildcard) && (len(studioIds) == 0 || studioIds[0] == wildcard) && (len(tagIds) == 0 || tagIds[0] == wildcard)
}

func (j *autoTagJob) autoTagFiles(ctx context.Context, progress *job.Progress, paths []string, performers, studios, tags bool) {
	t := autoTagFilesTask{
		paths:      paths,
		performers: performers,
		studios:    studios,
		tags:       tags,
		ctx:        ctx,
		progress:   progress,
		txnManager: j.txnManager,
	}

	t.process()
}

func (j *autoTagJob) autoTagSpecific(ctx context.Context, progress *job.Progress) {
	input := j.input
	performerIds := input.Performers
	studioIds := input.Studios
	tagIds := input.Tags

	performerCount := len(performerIds)
	studioCount := len(studioIds)
	tagCount := len(tagIds)

	if err := j.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		performerQuery := r.Performer()
		studioQuery := r.Studio()
		tagQuery := r.Tag()

		const wildcard = "*"
		var err error
		if performerCount == 1 && performerIds[0] == wildcard {
			performerCount, err = performerQuery.Count()
			if err != nil {
				return fmt.Errorf("error getting performer count: %s", err.Error())
			}
		}
		if studioCount == 1 && studioIds[0] == wildcard {
			studioCount, err = studioQuery.Count()
			if err != nil {
				return fmt.Errorf("error getting studio count: %s", err.Error())
			}
		}
		if tagCount == 1 && tagIds[0] == wildcard {
			tagCount, err = tagQuery.Count()
			if err != nil {
				return fmt.Errorf("error getting tag count: %s", err.Error())
			}
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
		return
	}

	total := performerCount + studioCount + tagCount
	progress.SetTotal(total)

	logger.Infof("Starting autotag of %d performers, %d studios, %d tags", performerCount, studioCount, tagCount)

	j.autoTagPerformers(ctx, progress, input.Paths, performerIds)
	j.autoTagStudios(ctx, progress, input.Paths, studioIds)
	j.autoTagTags(ctx, progress, input.Paths, tagIds)

	logger.Info("Finished autotag")
}

func (j *autoTagJob) autoTagPerformers(ctx context.Context, progress *job.Progress, paths []string, performerIds []string) {
	if job.IsCancelled(ctx) {
		return
	}

	for _, performerId := range performerIds {
		var performers []*models.Performer

		if err := j.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			performerQuery := r.Performer()

			if performerId == "*" {
				var err error
				performers, err = performerQuery.All()
				if err != nil {
					return fmt.Errorf("error querying performers: %s", err.Error())
				}
			} else {
				performerIdInt, err := strconv.Atoi(performerId)
				if err != nil {
					return fmt.Errorf("error parsing performer id %s: %s", performerId, err.Error())
				}

				performer, err := performerQuery.Find(performerIdInt)
				if err != nil {
					return fmt.Errorf("error finding performer id %s: %s", performerId, err.Error())
				}

				if performer == nil {
					return fmt.Errorf("performer with id %s not found", performerId)
				}
				performers = append(performers, performer)
			}

			for _, performer := range performers {
				if job.IsCancelled(ctx) {
					logger.Info("Stopping due to user request")
					return nil
				}

				if err := j.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
					if err := autotag.PerformerScenes(performer, paths, r.Scene()); err != nil {
						return err
					}
					if err := autotag.PerformerImages(performer, paths, r.Image()); err != nil {
						return err
					}
					if err := autotag.PerformerGalleries(performer, paths, r.Gallery()); err != nil {
						return err
					}

					return nil
				}); err != nil {
					return fmt.Errorf("error auto-tagging performer '%s': %s", performer.Name.String, err.Error())
				}

				progress.Increment()
			}

			return nil
		}); err != nil {
			logger.Error(err.Error())
			continue
		}
	}
}

func (j *autoTagJob) autoTagStudios(ctx context.Context, progress *job.Progress, paths []string, studioIds []string) {
	if job.IsCancelled(ctx) {
		return
	}

	for _, studioId := range studioIds {
		var studios []*models.Studio

		if err := j.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			studioQuery := r.Studio()
			if studioId == "*" {
				var err error
				studios, err = studioQuery.All()
				if err != nil {
					return fmt.Errorf("error querying studios: %s", err.Error())
				}
			} else {
				studioIdInt, err := strconv.Atoi(studioId)
				if err != nil {
					return fmt.Errorf("error parsing studio id %s: %s", studioId, err.Error())
				}

				studio, err := studioQuery.Find(studioIdInt)
				if err != nil {
					return fmt.Errorf("error finding studio id %s: %s", studioId, err.Error())
				}

				if studio == nil {
					return fmt.Errorf("studio with id %s not found", studioId)
				}

				studios = append(studios, studio)
			}

			for _, studio := range studios {
				if job.IsCancelled(ctx) {
					logger.Info("Stopping due to user request")
					return nil
				}

				if err := j.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
					aliases, err := r.Studio().GetAliases(studio.ID)
					if err != nil {
						return err
					}

					if err := autotag.StudioScenes(studio, paths, aliases, r.Scene()); err != nil {
						return err
					}
					if err := autotag.StudioImages(studio, paths, aliases, r.Image()); err != nil {
						return err
					}
					if err := autotag.StudioGalleries(studio, paths, aliases, r.Gallery()); err != nil {
						return err
					}

					return nil
				}); err != nil {
					return fmt.Errorf("error auto-tagging studio '%s': %s", studio.Name.String, err.Error())
				}

				progress.Increment()
			}

			return nil
		}); err != nil {
			logger.Error(err.Error())
			continue
		}
	}
}

func (j *autoTagJob) autoTagTags(ctx context.Context, progress *job.Progress, paths []string, tagIds []string) {
	if job.IsCancelled(ctx) {
		return
	}

	for _, tagId := range tagIds {
		var tags []*models.Tag
		if err := j.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			tagQuery := r.Tag()
			if tagId == "*" {
				var err error
				tags, err = tagQuery.All()
				if err != nil {
					return fmt.Errorf("error querying tags: %s", err.Error())
				}
			} else {
				tagIdInt, err := strconv.Atoi(tagId)
				if err != nil {
					return fmt.Errorf("error parsing tag id %s: %s", tagId, err.Error())
				}

				tag, err := tagQuery.Find(tagIdInt)
				if err != nil {
					return fmt.Errorf("error finding tag id %s: %s", tagId, err.Error())
				}
				tags = append(tags, tag)
			}

			for _, tag := range tags {
				if job.IsCancelled(ctx) {
					logger.Info("Stopping due to user request")
					return nil
				}

				if err := j.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
					aliases, err := r.Tag().GetAliases(tag.ID)
					if err != nil {
						return err
					}

					if err := autotag.TagScenes(tag, paths, aliases, r.Scene()); err != nil {
						return err
					}
					if err := autotag.TagImages(tag, paths, aliases, r.Image()); err != nil {
						return err
					}
					if err := autotag.TagGalleries(tag, paths, aliases, r.Gallery()); err != nil {
						return err
					}

					return nil
				}); err != nil {
					return fmt.Errorf("error auto-tagging tag '%s': %s", tag.Name, err.Error())
				}

				progress.Increment()
			}

			return nil
		}); err != nil {
			logger.Error(err.Error())
			continue
		}
	}
}

type autoTagFilesTask struct {
	paths      []string
	performers bool
	studios    bool
	tags       bool

	ctx        context.Context
	progress   *job.Progress
	txnManager models.TransactionManager
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

	if len(t.paths) == 0 {
		ret.Path = &models.StringCriterionInput{
			Modifier: models.CriterionModifierNotNull,
		}
	}

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
	if job.IsCancelled(t.ctx) {
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
			if job.IsCancelled(t.ctx) {
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

			t.progress.Increment()
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
	if job.IsCancelled(t.ctx) {
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
			if job.IsCancelled(t.ctx) {
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

			t.progress.Increment()
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
	if job.IsCancelled(t.ctx) {
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
			if job.IsCancelled(t.ctx) {
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

			t.progress.Increment()
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

		t.progress.SetTotal(total)

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

		if job.IsCancelled(t.ctx) {
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
				return fmt.Errorf("error tagging scene performers for %s: %v", t.scene.Path, err)
			}
		}
		if t.studios {
			if err := autotag.SceneStudios(t.scene, r.Scene(), r.Studio()); err != nil {
				return fmt.Errorf("error tagging scene studio for %s: %v", t.scene.Path, err)
			}
		}
		if t.tags {
			if err := autotag.SceneTags(t.scene, r.Scene(), r.Tag()); err != nil {
				return fmt.Errorf("error tagging scene tags for %s: %v", t.scene.Path, err)
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
				return fmt.Errorf("error tagging image performers for %s: %v", t.image.Path, err)
			}
		}
		if t.studios {
			if err := autotag.ImageStudios(t.image, r.Image(), r.Studio()); err != nil {
				return fmt.Errorf("error tagging image studio for %s: %v", t.image.Path, err)
			}
		}
		if t.tags {
			if err := autotag.ImageTags(t.image, r.Image(), r.Tag()); err != nil {
				return fmt.Errorf("error tagging image tags for %s: %v", t.image.Path, err)
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
				return fmt.Errorf("error tagging gallery performers for %s: %v", t.gallery.Path.String, err)
			}
		}
		if t.studios {
			if err := autotag.GalleryStudios(t.gallery, r.Gallery(), r.Studio()); err != nil {
				return fmt.Errorf("error tagging gallery studio for %s: %v", t.gallery.Path.String, err)
			}
		}
		if t.tags {
			if err := autotag.GalleryTags(t.gallery, r.Gallery(), r.Tag()); err != nil {
				return fmt.Errorf("error tagging gallery tags for %s: %v", t.gallery.Path.String, err)
			}
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}
