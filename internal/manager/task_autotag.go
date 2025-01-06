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
	repository models.Repository
	input      AutoTagMetadataInput

	cache match.Cache
}

func (j *autoTagJob) Execute(ctx context.Context, progress *job.Progress) error {
	begin := time.Now()

	input := j.input
	if j.isFileBasedAutoTag(input) {
		// doing file-based auto-tag
		j.autoTagFiles(ctx, progress, input.Paths, len(input.Performers) > 0, len(input.Studios) > 0, len(input.Tags) > 0)
	} else {
		// doing specific performer/studio/tag auto-tag
		j.autoTagSpecific(ctx, progress)
	}

	logger.Infof("Finished auto-tag after %s", time.Since(begin).String())
	return nil
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
		repository: j.repository,
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

	r := j.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		performerQuery := r.Performer
		studioQuery := r.Studio
		tagQuery := r.Tag

		const wildcard = "*"
		var err error
		if performerCount == 1 && performerIds[0] == wildcard {
			performerCount, err = performerQuery.Count(ctx)
			if err != nil {
				return fmt.Errorf("getting performer count: %v", err)
			}
		}
		if studioCount == 1 && studioIds[0] == wildcard {
			studioCount, err = studioQuery.Count(ctx)
			if err != nil {
				return fmt.Errorf("getting studio count: %v", err)
			}
		}
		if tagCount == 1 && tagIds[0] == wildcard {
			tagCount, err = tagQuery.Count(ctx)
			if err != nil {
				return fmt.Errorf("getting tag count: %v", err)
			}
		}

		return nil
	}); err != nil {
		if !job.IsCancelled(ctx) {
			logger.Errorf("auto-tag error: %v", err)
		}
		return
	}

	total := performerCount + studioCount + tagCount
	progress.SetTotal(total)

	logger.Infof("Starting auto-tag of %d performers, %d studios, %d tags", performerCount, studioCount, tagCount)

	j.autoTagPerformers(ctx, progress, input.Paths, performerIds)
	j.autoTagStudios(ctx, progress, input.Paths, studioIds)
	j.autoTagTags(ctx, progress, input.Paths, tagIds)
}

func (j *autoTagJob) autoTagPerformers(ctx context.Context, progress *job.Progress, paths []string, performerIds []string) {
	if job.IsCancelled(ctx) {
		return
	}

	r := j.repository
	tagger := autotag.Tagger{
		TxnManager: r.TxnManager,
		Cache:      &j.cache,
	}

	for _, performerId := range performerIds {
		var performers []*models.Performer

		if err := r.WithDB(ctx, func(ctx context.Context) error {
			performerQuery := r.Performer
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
					return fmt.Errorf("querying performers: %w", err)
				}
			} else {
				performerIdInt, err := strconv.Atoi(performerId)
				if err != nil {
					return fmt.Errorf("parsing performer id %s: %w", performerId, err)
				}

				performer, err := performerQuery.Find(ctx, performerIdInt)
				if err != nil {
					return fmt.Errorf("finding performer id %s: %w", performerId, err)
				}

				if performer == nil {
					return fmt.Errorf("performer with id %s not found", performerId)
				}

				if performer.IgnoreAutoTag {
					logger.Infof("Skipping performer %s because auto-tag is disabled", performer.Name)
					return nil
				}

				if err := performer.LoadAliases(ctx, r.Performer); err != nil {
					return fmt.Errorf("loading aliases for performer %d: %w", performer.ID, err)
				}
				performers = append(performers, performer)
			}

			for _, performer := range performers {
				if job.IsCancelled(ctx) {
					return nil
				}

				err := func() error {
					if err := tagger.PerformerScenes(ctx, performer, paths, r.Scene); err != nil {
						return fmt.Errorf("processing scenes: %w", err)
					}
					if err := tagger.PerformerImages(ctx, performer, paths, r.Image); err != nil {
						return fmt.Errorf("processing images: %w", err)
					}
					if err := tagger.PerformerGalleries(ctx, performer, paths, r.Gallery); err != nil {
						return fmt.Errorf("processing galleries: %w", err)
					}

					return nil
				}()

				if job.IsCancelled(ctx) {
					return nil
				}

				if err != nil {
					return fmt.Errorf("tagging performer '%s': %s", performer.Name, err.Error())
				}

				progress.Increment()
			}

			return nil
		}); err != nil {
			logger.Errorf("auto-tag error: %v", err)
		}

		if job.IsCancelled(ctx) {
			logger.Info("Stopping performer auto-tag due to user request")
			return
		}
	}
}

func (j *autoTagJob) autoTagStudios(ctx context.Context, progress *job.Progress, paths []string, studioIds []string) {
	if job.IsCancelled(ctx) {
		return
	}

	r := j.repository
	tagger := autotag.Tagger{
		TxnManager: r.TxnManager,
		Cache:      &j.cache,
	}

	for _, studioId := range studioIds {
		var studios []*models.Studio

		if err := r.WithDB(ctx, func(ctx context.Context) error {
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
					return fmt.Errorf("querying studios: %v", err)
				}
			} else {
				studioIdInt, err := strconv.Atoi(studioId)
				if err != nil {
					return fmt.Errorf("parsing studio id %s: %s", studioId, err.Error())
				}

				studio, err := studioQuery.Find(ctx, studioIdInt)
				if err != nil {
					return fmt.Errorf("finding studio id %s: %s", studioId, err.Error())
				}

				if studio == nil {
					return fmt.Errorf("studio with id %s not found", studioId)
				}

				if studio.IgnoreAutoTag {
					logger.Infof("Skipping studio %s because auto-tag is disabled", studio.Name)
					return nil
				}

				studios = append(studios, studio)
			}

			for _, studio := range studios {
				if job.IsCancelled(ctx) {
					return nil
				}

				err := func() error {
					aliases, err := r.Studio.GetAliases(ctx, studio.ID)
					if err != nil {
						return fmt.Errorf("getting studio aliases: %w", err)
					}

					if err := tagger.StudioScenes(ctx, studio, paths, aliases, r.Scene); err != nil {
						return fmt.Errorf("processing scenes: %w", err)
					}
					if err := tagger.StudioImages(ctx, studio, paths, aliases, r.Image); err != nil {
						return fmt.Errorf("processing images: %w", err)
					}
					if err := tagger.StudioGalleries(ctx, studio, paths, aliases, r.Gallery); err != nil {
						return fmt.Errorf("processing galleries: %w", err)
					}

					return nil
				}()

				if job.IsCancelled(ctx) {
					return nil
				}

				if err != nil {
					return fmt.Errorf("tagging studio '%s': %s", studio.Name, err.Error())
				}

				progress.Increment()
			}

			return nil
		}); err != nil {
			logger.Errorf("auto-tag error: %v", err)
		}

		if job.IsCancelled(ctx) {
			logger.Info("Stopping studio auto-tag due to user request")
			return
		}
	}
}

func (j *autoTagJob) autoTagTags(ctx context.Context, progress *job.Progress, paths []string, tagIds []string) {
	if job.IsCancelled(ctx) {
		return
	}

	r := j.repository
	tagger := autotag.Tagger{
		TxnManager: r.TxnManager,
		Cache:      &j.cache,
	}

	for _, tagId := range tagIds {
		var tags []*models.Tag
		if err := r.WithDB(ctx, func(ctx context.Context) error {
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
					return fmt.Errorf("querying tags: %v", err)
				}
			} else {
				tagIdInt, err := strconv.Atoi(tagId)
				if err != nil {
					return fmt.Errorf("parsing tag id %s: %s", tagId, err.Error())
				}

				tag, err := tagQuery.Find(ctx, tagIdInt)
				if err != nil {
					return fmt.Errorf("finding tag id %s: %s", tagId, err.Error())
				}

				if tag == nil {
					return fmt.Errorf("tag with id %s not found", tagId)
				}

				if tag.IgnoreAutoTag {
					logger.Infof("Skipping tag %s because auto-tag is disabled", tag.Name)
					return nil
				}

				tags = append(tags, tag)
			}

			for _, tag := range tags {
				if job.IsCancelled(ctx) {
					return nil
				}

				err := func() error {
					aliases, err := r.Tag.GetAliases(ctx, tag.ID)
					if err != nil {
						return fmt.Errorf("getting tag aliases: %w", err)
					}

					if err := tagger.TagScenes(ctx, tag, paths, aliases, r.Scene); err != nil {
						return fmt.Errorf("processing scenes: %w", err)
					}
					if err := tagger.TagImages(ctx, tag, paths, aliases, r.Image); err != nil {
						return fmt.Errorf("processing images: %w", err)
					}
					if err := tagger.TagGalleries(ctx, tag, paths, aliases, r.Gallery); err != nil {
						return fmt.Errorf("processing galleries: %w", err)
					}

					return nil
				}()

				if job.IsCancelled(ctx) {
					return nil
				}

				if err != nil {
					return fmt.Errorf("tagging tag '%s': %s", tag.Name, err.Error())
				}

				progress.Increment()
			}

			return nil
		}); err != nil {
			logger.Errorf("auto-tag error: %v", err)
		}

		if job.IsCancelled(ctx) {
			logger.Info("Stopping tag auto-tag due to user request")
			return
		}
	}
}

type autoTagFilesTask struct {
	paths      []string
	performers bool
	studios    bool
	tags       bool

	progress   *job.Progress
	repository models.Repository
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

func (t *autoTagFilesTask) getCount(ctx context.Context) (int, error) {
	r := t.repository

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
		return 0, fmt.Errorf("getting scene count: %w", err)
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
		return 0, fmt.Errorf("getting image count: %w", err)
	}

	imageCount := imageResults.Count

	_, galleryCount, err := r.Gallery.Query(ctx, t.makeGalleryFilter(), findFilter)
	if err != nil {
		return 0, fmt.Errorf("getting gallery count: %w", err)
	}

	return sceneCount + imageCount + galleryCount, nil
}

func (t *autoTagFilesTask) processScenes(ctx context.Context) {
	if job.IsCancelled(ctx) {
		return
	}

	logger.Info("Auto-tagging scenes...")

	batchSize := 1000

	findFilter := models.BatchFindFilter(batchSize)
	sceneFilter := t.makeSceneFilter()

	r := t.repository

	more := true
	for more {
		var scenes []*models.Scene
		if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
			var err error
			scenes, err = scene.Query(ctx, r.Scene, sceneFilter, findFilter)
			return err
		}); err != nil {
			if !job.IsCancelled(ctx) {
				logger.Errorf("error querying scenes for auto-tag: %w", err)
			}
			return
		}

		for _, ss := range scenes {
			if job.IsCancelled(ctx) {
				logger.Info("Stopping auto-tag due to user request")
				return
			}

			tt := autoTagSceneTask{
				repository: r,
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
}

func (t *autoTagFilesTask) processImages(ctx context.Context) {
	if job.IsCancelled(ctx) {
		return
	}

	logger.Info("Auto-tagging images...")

	batchSize := 1000

	findFilter := models.BatchFindFilter(batchSize)
	imageFilter := t.makeImageFilter()

	r := t.repository

	more := true
	for more {
		var images []*models.Image
		if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
			var err error
			images, err = image.Query(ctx, r.Image, imageFilter, findFilter)
			return err
		}); err != nil {
			if !job.IsCancelled(ctx) {
				logger.Errorf("error querying images for auto-tag: %w", err)
			}
			return
		}

		for _, ss := range images {
			if job.IsCancelled(ctx) {
				logger.Info("Stopping auto-tag due to user request")
				return
			}

			tt := autoTagImageTask{
				repository: t.repository,
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
}

func (t *autoTagFilesTask) processGalleries(ctx context.Context) {
	if job.IsCancelled(ctx) {
		return
	}

	logger.Info("Auto-tagging galleries...")

	batchSize := 1000

	findFilter := models.BatchFindFilter(batchSize)
	galleryFilter := t.makeGalleryFilter()

	r := t.repository

	more := true
	for more {
		var galleries []*models.Gallery
		if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
			var err error
			galleries, _, err = r.Gallery.Query(ctx, galleryFilter, findFilter)
			return err
		}); err != nil {
			if !job.IsCancelled(ctx) {
				logger.Errorf("error querying galleries for auto-tag: %w", err)
			}
			return
		}

		for _, ss := range galleries {
			if job.IsCancelled(ctx) {
				logger.Info("Stopping auto-tag due to user request")
				return
			}

			tt := autoTagGalleryTask{
				repository: t.repository,
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
}

func (t *autoTagFilesTask) process(ctx context.Context) {
	if err := t.repository.WithReadTxn(ctx, func(ctx context.Context) error {
		total, err := t.getCount(ctx)
		if err != nil {
			return err
		}

		t.progress.SetTotal(total)
		logger.Infof("Starting auto-tag of %d files", total)

		return nil
	}); err != nil {
		if !job.IsCancelled(ctx) {
			logger.Errorf("error getting file count for auto-tag task: %v", err)
		}
		return
	}

	t.processScenes(ctx)
	t.processImages(ctx)
	t.processGalleries(ctx)
}

type autoTagSceneTask struct {
	repository models.Repository
	scene      *models.Scene

	performers bool
	studios    bool
	tags       bool

	cache *match.Cache
}

func (t *autoTagSceneTask) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	r := t.repository
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		if t.scene.Path == "" {
			// nothing to do
			return nil
		}

		if t.performers {
			if err := autotag.ScenePerformers(ctx, t.scene, r.Scene, r.Performer, t.cache); err != nil {
				return fmt.Errorf("tagging scene performers for %s: %v", t.scene.DisplayName(), err)
			}
		}
		if t.studios {
			if err := autotag.SceneStudios(ctx, t.scene, r.Scene, r.Studio, t.cache); err != nil {
				return fmt.Errorf("tagging scene studio for %s: %v", t.scene.DisplayName(), err)
			}
		}
		if t.tags {
			if err := autotag.SceneTags(ctx, t.scene, r.Scene, r.Tag, t.cache); err != nil {
				return fmt.Errorf("tagging scene tags for %s: %v", t.scene.DisplayName(), err)
			}
		}

		return nil
	}); err != nil {
		if !job.IsCancelled(ctx) {
			logger.Errorf("auto-tag error: %v", err)
		}
	}
}

type autoTagImageTask struct {
	repository models.Repository
	image      *models.Image

	performers bool
	studios    bool
	tags       bool

	cache *match.Cache
}

func (t *autoTagImageTask) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	r := t.repository
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		if t.performers {
			if err := autotag.ImagePerformers(ctx, t.image, r.Image, r.Performer, t.cache); err != nil {
				return fmt.Errorf("tagging image performers for %s: %v", t.image.DisplayName(), err)
			}
		}
		if t.studios {
			if err := autotag.ImageStudios(ctx, t.image, r.Image, r.Studio, t.cache); err != nil {
				return fmt.Errorf("tagging image studio for %s: %v", t.image.DisplayName(), err)
			}
		}
		if t.tags {
			if err := autotag.ImageTags(ctx, t.image, r.Image, r.Tag, t.cache); err != nil {
				return fmt.Errorf("tagging image tags for %s: %v", t.image.DisplayName(), err)
			}
		}

		return nil
	}); err != nil {
		if !job.IsCancelled(ctx) {
			logger.Errorf("auto-tag error: %v", err)
		}
	}
}

type autoTagGalleryTask struct {
	repository models.Repository
	gallery    *models.Gallery

	performers bool
	studios    bool
	tags       bool

	cache *match.Cache
}

func (t *autoTagGalleryTask) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	r := t.repository
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		if t.performers {
			if err := autotag.GalleryPerformers(ctx, t.gallery, r.Gallery, r.Performer, t.cache); err != nil {
				return fmt.Errorf("tagging gallery performers for %s: %v", t.gallery.DisplayName(), err)
			}
		}
		if t.studios {
			if err := autotag.GalleryStudios(ctx, t.gallery, r.Gallery, r.Studio, t.cache); err != nil {
				return fmt.Errorf("tagging gallery studio for %s: %v", t.gallery.DisplayName(), err)
			}
		}
		if t.tags {
			if err := autotag.GalleryTags(ctx, t.gallery, r.Gallery, r.Tag, t.cache); err != nil {
				return fmt.Errorf("tagging gallery tags for %s: %v", t.gallery.DisplayName(), err)
			}
		}

		return nil
	}); err != nil {
		if !job.IsCancelled(ctx) {
			logger.Errorf("auto-tag error: %v", err)
		}
	}
}
