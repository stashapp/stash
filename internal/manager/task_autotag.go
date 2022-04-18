package manager

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/stashapp/stash/internal/autotag"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

type autoTagJob struct {
	txnManager models.Repository
	input      AutoTagMetadataInput

	cache match.Cache
}

func (j *autoTagJob) Execute(ctx context.Context, progress *job.Progress) {
	begin := time.Now()

	input := j.input
	if j.isFileBasedAutoTag(input) {
		// doing file-based auto-tag
		j.autoTagFiles(ctx, progress, input.Paths, len(input.Performers) > 0, len(input.Studios) > 0, len(input.Tags) > 0)
	} else {
		// doing specific performer/studio/tag auto-tag
		j.autoTagSpecific(ctx, progress)
	}

	logger.Infof("Finished autotag after %s", time.Since(begin).String())
}

func (j *autoTagJob) isFileBasedAutoTag(input AutoTagMetadataInput) bool {
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
		progress:   progress,
		txnManager: j.txnManager,
		cache:      &j.cache,
	}

	t.process(ctx)
}

func (j *autoTagJob) autoTagSpecific(ctx context.Context, progress *job.Progress) {
	input := j.input
	performerIds := input.Performers
	studioIds := input.Studios
	tagIds := input.Tags

	performerCount := len(performerIds)
	studioCount := len(studioIds)
	tagCount := len(tagIds)

	if err := j.txnManager.WithTxn(ctx, func(ctx context.Context) error {
		r := j.txnManager
		performerQuery := r.Performer
		studioQuery := r.Studio
		tagQuery := r.Tag

		const wildcard = "*"
		var err error
		if performerCount == 1 && performerIds[0] == wildcard {
			performerCount, err = performerQuery.Count(ctx)
			if err != nil {
				return fmt.Errorf("error getting performer count: %v", err)
			}
		}
		if studioCount == 1 && studioIds[0] == wildcard {
			studioCount, err = studioQuery.Count(ctx)
			if err != nil {
				return fmt.Errorf("error getting studio count: %v", err)
			}
		}
		if tagCount == 1 && tagIds[0] == wildcard {
			tagCount, err = tagQuery.Count(ctx)
			if err != nil {
				return fmt.Errorf("error getting tag count: %v", err)
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
}

func (j *autoTagJob) autoTagPerformers(ctx context.Context, progress *job.Progress, paths []string, performerIds []string) {
	if job.IsCancelled(ctx) {
		return
	}

	for _, performerId := range performerIds {
		var performers []*models.Performer

		if err := j.txnManager.WithTxn(ctx, func(ctx context.Context) error {
			performerQuery := j.txnManager.Performer
			ignoreAutoTag := false
			perPage := -1

			if performerId == "*" {
				var err error
				performers, _, err = performerQuery.Query(ctx, &models.PerformerFilterType{
					IgnoreAutoTag: &ignoreAutoTag,
				}, &models.FindFilterType{
					PerPage: &perPage,
				})
				if err != nil {
					return fmt.Errorf("error querying performers: %v", err)
				}
			} else {
				performerIdInt, err := strconv.Atoi(performerId)
				if err != nil {
					return fmt.Errorf("error parsing performer id %s: %s", performerId, err.Error())
				}

				performer, err := performerQuery.Find(ctx, performerIdInt)
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

				if err := j.txnManager.WithTxn(ctx, func(ctx context.Context) error {
					r := j.txnManager
					if err := autotag.PerformerScenes(ctx, performer, paths, r.Scene, &j.cache); err != nil {
						return err
					}
					if err := autotag.PerformerImages(ctx, performer, paths, r.Image, &j.cache); err != nil {
						return err
					}
					if err := autotag.PerformerGalleries(ctx, performer, paths, r.Gallery, &j.cache); err != nil {
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

	r := j.txnManager

	for _, studioId := range studioIds {
		var studios []*models.Studio

		if err := r.WithTxn(ctx, func(ctx context.Context) error {
			studioQuery := r.Studio
			ignoreAutoTag := false
			perPage := -1
			if studioId == "*" {
				var err error
				studios, _, err = studioQuery.Query(ctx, &models.StudioFilterType{
					IgnoreAutoTag: &ignoreAutoTag,
				}, &models.FindFilterType{
					PerPage: &perPage,
				})
				if err != nil {
					return fmt.Errorf("error querying studios: %v", err)
				}
			} else {
				studioIdInt, err := strconv.Atoi(studioId)
				if err != nil {
					return fmt.Errorf("error parsing studio id %s: %s", studioId, err.Error())
				}

				studio, err := studioQuery.Find(ctx, studioIdInt)
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

				if err := j.txnManager.WithTxn(ctx, func(ctx context.Context) error {
					aliases, err := r.Studio.GetAliases(ctx, studio.ID)
					if err != nil {
						return err
					}

					if err := autotag.StudioScenes(ctx, studio, paths, aliases, r.Scene, &j.cache); err != nil {
						return err
					}
					if err := autotag.StudioImages(ctx, studio, paths, aliases, r.Image, &j.cache); err != nil {
						return err
					}
					if err := autotag.StudioGalleries(ctx, studio, paths, aliases, r.Gallery, &j.cache); err != nil {
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

	r := j.txnManager

	for _, tagId := range tagIds {
		var tags []*models.Tag
		if err := j.txnManager.WithTxn(ctx, func(ctx context.Context) error {
			tagQuery := r.Tag
			ignoreAutoTag := false
			perPage := -1
			if tagId == "*" {
				var err error
				tags, _, err = tagQuery.Query(ctx, &models.TagFilterType{
					IgnoreAutoTag: &ignoreAutoTag,
				}, &models.FindFilterType{
					PerPage: &perPage,
				})
				if err != nil {
					return fmt.Errorf("error querying tags: %v", err)
				}
			} else {
				tagIdInt, err := strconv.Atoi(tagId)
				if err != nil {
					return fmt.Errorf("error parsing tag id %s: %s", tagId, err.Error())
				}

				tag, err := tagQuery.Find(ctx, tagIdInt)
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

				if err := j.txnManager.WithTxn(ctx, func(ctx context.Context) error {
					aliases, err := r.Tag.GetAliases(ctx, tag.ID)
					if err != nil {
						return err
					}

					if err := autotag.TagScenes(ctx, tag, paths, aliases, r.Scene, &j.cache); err != nil {
						return err
					}
					if err := autotag.TagImages(ctx, tag, paths, aliases, r.Image, &j.cache); err != nil {
						return err
					}
					if err := autotag.TagGalleries(ctx, tag, paths, aliases, r.Gallery, &j.cache); err != nil {
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

	progress   *job.Progress
	txnManager models.Repository
	cache      *match.Cache
}

func (t *autoTagFilesTask) makeSceneFilter() *models.SceneFilterType {
	ret := scene.FilterFromPaths(t.paths)

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
			p += sep
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
			p += sep
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

func (t *autoTagFilesTask) getCount(ctx context.Context, r models.Repository) (int, error) {
	pp := 0
	findFilter := &models.FindFilterType{
		PerPage: &pp,
	}

	sceneResults, err := r.Scene.Query(ctx, models.SceneQueryOptions{
		QueryOptions: models.QueryOptions{
			FindFilter: findFilter,
			Count:      true,
		},
		SceneFilter: t.makeSceneFilter(),
	})
	if err != nil {
		return 0, err
	}

	sceneCount := sceneResults.Count

	imageResults, err := r.Image.Query(ctx, models.ImageQueryOptions{
		QueryOptions: models.QueryOptions{
			FindFilter: findFilter,
			Count:      true,
		},
		ImageFilter: t.makeImageFilter(),
	})
	if err != nil {
		return 0, err
	}

	imageCount := imageResults.Count

	_, galleryCount, err := r.Gallery.Query(ctx, t.makeGalleryFilter(), findFilter)
	if err != nil {
		return 0, err
	}

	return sceneCount + imageCount + galleryCount, nil
}

func (t *autoTagFilesTask) processScenes(ctx context.Context, r models.Repository) error {
	if job.IsCancelled(ctx) {
		return nil
	}

	batchSize := 1000

	findFilter := models.BatchFindFilter(batchSize)
	sceneFilter := t.makeSceneFilter()

	more := true
	for more {
		scenes, err := scene.Query(ctx, r.Scene, sceneFilter, findFilter)
		if err != nil {
			return err
		}

		for _, ss := range scenes {
			if job.IsCancelled(ctx) {
				return nil
			}

			tt := autoTagSceneTask{
				txnManager: t.txnManager,
				scene:      ss,
				performers: t.performers,
				studios:    t.studios,
				tags:       t.tags,
				cache:      t.cache,
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go tt.Start(ctx, &wg)
			wg.Wait()

			t.progress.Increment()
		}

		if len(scenes) != batchSize {
			more = false
		} else {
			*findFilter.Page++

			if *findFilter.Page%10 == 1 {
				logger.Infof("Processed %d scenes...", (*findFilter.Page-1)*batchSize)
			}
		}
	}

	return nil
}

func (t *autoTagFilesTask) processImages(ctx context.Context, r models.Repository) error {
	if job.IsCancelled(ctx) {
		return nil
	}

	batchSize := 1000

	findFilter := models.BatchFindFilter(batchSize)
	imageFilter := t.makeImageFilter()

	more := true
	for more {
		images, err := image.Query(ctx, r.Image, imageFilter, findFilter)
		if err != nil {
			return err
		}

		for _, ss := range images {
			if job.IsCancelled(ctx) {
				return nil
			}

			tt := autoTagImageTask{
				txnManager: t.txnManager,
				image:      ss,
				performers: t.performers,
				studios:    t.studios,
				tags:       t.tags,
				cache:      t.cache,
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go tt.Start(ctx, &wg)
			wg.Wait()

			t.progress.Increment()
		}

		if len(images) != batchSize {
			more = false
		} else {
			*findFilter.Page++

			if *findFilter.Page%10 == 1 {
				logger.Infof("Processed %d images...", (*findFilter.Page-1)*batchSize)
			}
		}
	}

	return nil
}

func (t *autoTagFilesTask) processGalleries(ctx context.Context, r models.Repository) error {
	if job.IsCancelled(ctx) {
		return nil
	}

	batchSize := 1000

	findFilter := models.BatchFindFilter(batchSize)
	galleryFilter := t.makeGalleryFilter()

	more := true
	for more {
		galleries, _, err := r.Gallery.Query(ctx, galleryFilter, findFilter)
		if err != nil {
			return err
		}

		for _, ss := range galleries {
			if job.IsCancelled(ctx) {
				return nil
			}

			tt := autoTagGalleryTask{
				txnManager: t.txnManager,
				gallery:    ss,
				performers: t.performers,
				studios:    t.studios,
				tags:       t.tags,
				cache:      t.cache,
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go tt.Start(ctx, &wg)
			wg.Wait()

			t.progress.Increment()
		}

		if len(galleries) != batchSize {
			more = false
		} else {
			*findFilter.Page++

			if *findFilter.Page%10 == 1 {
				logger.Infof("Processed %d galleries...", (*findFilter.Page-1)*batchSize)
			}
		}
	}

	return nil
}

func (t *autoTagFilesTask) process(ctx context.Context) {
	r := t.txnManager
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		total, err := t.getCount(ctx, t.txnManager)
		if err != nil {
			return err
		}

		t.progress.SetTotal(total)

		logger.Infof("Starting autotag of %d files", total)

		logger.Info("Autotagging scenes...")
		if err := t.processScenes(ctx, r); err != nil {
			return err
		}

		logger.Info("Autotagging images...")
		if err := t.processImages(ctx, r); err != nil {
			return err
		}

		logger.Info("Autotagging galleries...")
		if err := t.processGalleries(ctx, r); err != nil {
			return err
		}

		if job.IsCancelled(ctx) {
			logger.Info("Stopping due to user request")
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}

type autoTagSceneTask struct {
	txnManager models.Repository
	scene      *models.Scene

	performers bool
	studios    bool
	tags       bool

	cache *match.Cache
}

func (t *autoTagSceneTask) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	r := t.txnManager
	if err := t.txnManager.WithTxn(ctx, func(ctx context.Context) error {
		if t.performers {
			if err := autotag.ScenePerformers(ctx, t.scene, r.Scene, r.Performer, t.cache); err != nil {
				return fmt.Errorf("error tagging scene performers for %s: %v", t.scene.Path, err)
			}
		}
		if t.studios {
			if err := autotag.SceneStudios(ctx, t.scene, r.Scene, r.Studio, t.cache); err != nil {
				return fmt.Errorf("error tagging scene studio for %s: %v", t.scene.Path, err)
			}
		}
		if t.tags {
			if err := autotag.SceneTags(ctx, t.scene, r.Scene, r.Tag, t.cache); err != nil {
				return fmt.Errorf("error tagging scene tags for %s: %v", t.scene.Path, err)
			}
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}

type autoTagImageTask struct {
	txnManager models.Repository
	image      *models.Image

	performers bool
	studios    bool
	tags       bool

	cache *match.Cache
}

func (t *autoTagImageTask) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	r := t.txnManager
	if err := t.txnManager.WithTxn(ctx, func(ctx context.Context) error {
		if t.performers {
			if err := autotag.ImagePerformers(ctx, t.image, r.Image, r.Performer, t.cache); err != nil {
				return fmt.Errorf("error tagging image performers for %s: %v", t.image.Path, err)
			}
		}
		if t.studios {
			if err := autotag.ImageStudios(ctx, t.image, r.Image, r.Studio, t.cache); err != nil {
				return fmt.Errorf("error tagging image studio for %s: %v", t.image.Path, err)
			}
		}
		if t.tags {
			if err := autotag.ImageTags(ctx, t.image, r.Image, r.Tag, t.cache); err != nil {
				return fmt.Errorf("error tagging image tags for %s: %v", t.image.Path, err)
			}
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}

type autoTagGalleryTask struct {
	txnManager models.Repository
	gallery    *models.Gallery

	performers bool
	studios    bool
	tags       bool

	cache *match.Cache
}

func (t *autoTagGalleryTask) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	r := t.txnManager
	if err := t.txnManager.WithTxn(ctx, func(ctx context.Context) error {
		if t.performers {
			if err := autotag.GalleryPerformers(ctx, t.gallery, r.Gallery, r.Performer, t.cache); err != nil {
				return fmt.Errorf("error tagging gallery performers for %s: %v", t.gallery.Path.String, err)
			}
		}
		if t.studios {
			if err := autotag.GalleryStudios(ctx, t.gallery, r.Gallery, r.Studio, t.cache); err != nil {
				return fmt.Errorf("error tagging gallery studio for %s: %v", t.gallery.Path.String, err)
			}
		}
		if t.tags {
			if err := autotag.GalleryTags(ctx, t.gallery, r.Gallery, r.Tag, t.cache); err != nil {
				return fmt.Errorf("error tagging gallery tags for %s: %v", t.gallery.Path.String, err)
			}
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}
