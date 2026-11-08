package manager

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/stashapp/stash/internal/identify"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/stashbox"
	"github.com/stashapp/stash/pkg/txn"
)

var ErrInput = errors.New("invalid request input")

type IdentifyJob struct {
	postHookExecutor        identify.SceneUpdatePostHookExecutor
	galleryPostHookExecutor identify.GalleryUpdatePostHookExecutor
	input                   identify.Options

	stashBoxes       []*models.StashBox
	progress         *job.Progress
	progressTotalSet bool
	repository       models.Repository

	sourcesFn        func() ([]identify.ScraperSource, error)
	gallerySourcesFn func() ([]identify.GalleryScraperSource, error)

	identifySceneFn   func(ctx context.Context, s *models.Scene, sources []identify.ScraperSource)
	identifyGalleryFn func(ctx context.Context, g *models.Gallery, sources []identify.GalleryScraperSource)
}

func CreateIdentifyJob(input identify.Options) *IdentifyJob {
	j := &IdentifyJob{
		postHookExecutor:        instance.PluginCache,
		galleryPostHookExecutor: instance.PluginCache,
		input:                   input,
		stashBoxes:              instance.Config.GetStashBoxes(),
		repository:              instance.Repository,
	}
	j.sourcesFn = j.getSources
	j.gallerySourcesFn = j.getGallerySources
	j.identifySceneFn = j.identifyScene
	j.identifyGalleryFn = j.identifyGallery
	return j
}

func (j *IdentifyJob) Execute(ctx context.Context, progress *job.Progress) error {
	j.progress = progress
	j.progressTotalSet = false

	if len(j.input.Sources) == 0 {
		return nil
	}

	// run gallery identification if gallery IDs given or in identify-all mode
	if len(j.input.GalleryIDs) > 0 {
		if err := j.executeGalleryIDs(ctx); err != nil {
			return err
		}
	} else if len(j.input.SceneIDs) == 0 {
		gallerySources, err := j.gallerySourcesFn()
		if err != nil {
			return err
		}

		if len(gallerySources) > 0 {
			if err := j.identifyAllGalleries(ctx, gallerySources); err != nil {
				return err
			}
		}
	}

	// run scene identification if scene IDs given, in identify-all mode,
	// or no gallery IDs were explicitly provided (legacy behavior)
	if len(j.input.SceneIDs) > 0 || len(j.input.GalleryIDs) == 0 {
		sources, err := j.sourcesFn()
		if err != nil {
			return err
		}

		if len(sources) > 0 {
			r := j.repository
			if err := r.WithDB(ctx, func(ctx context.Context) error {
				if len(j.input.SceneIDs) == 0 {
					return j.identifyAllScenes(ctx, sources)
				}

				sceneIDs, err := stringslice.StringSliceToIntSlice(j.input.SceneIDs)
				if err != nil {
					return fmt.Errorf("invalid scene IDs: %w", err)
				}

				j.addProgressTotal(len(sceneIDs))
				for _, id := range sceneIDs {
					if job.IsCancelled(ctx) {
						break
					}

					scene, err := r.Scene.Find(ctx, id)
					if err != nil {
						return fmt.Errorf("finding scene id %d: %w", id, err)
					}

					if scene == nil {
						return fmt.Errorf("scene with id %d not found", id)
					}

					j.identifySceneFn(ctx, scene, sources)
				}

				return nil
			}); err != nil {
				return fmt.Errorf("error encountered while identifying scenes: %w", err)
			}
		}
	}

	return nil
}

func (j *IdentifyJob) addProgressTotal(total int) {
	if j.progressTotalSet {
		j.progress.AddTotal(total)
		return
	}

	j.progress.SetTotal(total)
	j.progressTotalSet = true
}

func (j *IdentifyJob) executeGalleryIDs(ctx context.Context) error {
	r := j.repository
	gallerySources, err := j.gallerySourcesFn()
	if err != nil {
		return err
	}

	if len(gallerySources) == 0 {
		return fmt.Errorf("no gallery sources available")
	}

	if err := r.WithDB(ctx, func(ctx context.Context) error {
		galleryIDs, err := stringslice.StringSliceToIntSlice(j.input.GalleryIDs)
		if err != nil {
			return fmt.Errorf("invalid gallery IDs: %w", err)
		}

		j.addProgressTotal(len(galleryIDs))
		for _, id := range galleryIDs {
			if job.IsCancelled(ctx) {
				break
			}

			g, err := r.Gallery.Find(ctx, id)
			if err != nil {
				return fmt.Errorf("finding gallery id %d: %w", id, err)
			}

			if g == nil {
				return fmt.Errorf("gallery with id %d not found", id)
			}

			j.identifyGalleryFn(ctx, g, gallerySources)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error encountered while identifying galleries: %w", err)
	}

	return nil
}

func (j *IdentifyJob) identifyAllScenes(ctx context.Context, sources []identify.ScraperSource) error {
	r := j.repository

	// exclude organised
	organised := false
	sceneFilter := scene.FilterFromPaths(j.input.Paths)
	sceneFilter.Organized = &organised

	sort := "path"
	findFilter := &models.FindFilterType{
		Sort: &sort,
	}

	// get the count
	pp := 0
	findFilter.PerPage = &pp
	countResult, err := r.Scene.Query(ctx, models.SceneQueryOptions{
		QueryOptions: models.QueryOptions{
			FindFilter: findFilter,
			Count:      true,
		},
		SceneFilter: sceneFilter,
	})
	if err != nil {
		return fmt.Errorf("error getting scene count: %w", err)
	}

	j.addProgressTotal(countResult.Count)

	return scene.BatchProcess(ctx, r.Scene, sceneFilter, findFilter, func(scene *models.Scene) error {
		if job.IsCancelled(ctx) {
			return nil
		}

		j.identifySceneFn(ctx, scene, sources)
		return nil
	})
}

func (j *IdentifyJob) identifyAllGalleries(ctx context.Context, sources []identify.GalleryScraperSource) error {
	r := j.repository

	// exclude organised
	organised := false
	galleryFilter := gallery.PathsFilter(j.input.Paths)
	if galleryFilter == nil {
		galleryFilter = &models.GalleryFilterType{}
	}
	galleryFilter.Organized = &organised

	sort := "path"
	findFilter := &models.FindFilterType{
		Sort: &sort,
	}

	// get the count
	pp := 0
	findFilter.PerPage = &pp
	_, countResult, err := r.Gallery.Query(ctx, galleryFilter, findFilter)
	if err != nil {
		return fmt.Errorf("error getting gallery count: %w", err)
	}

	j.addProgressTotal(countResult)

	return j.batchProcessGalleries(ctx, r.Gallery, galleryFilter, findFilter, func(g *models.Gallery) error {
		if job.IsCancelled(ctx) {
			return nil
		}

		j.identifyGalleryFn(ctx, g, sources)
		return nil
	})
}

func (j *IdentifyJob) batchProcessGalleries(ctx context.Context, reader models.GalleryQueryer, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType, fn func(gallery *models.Gallery) error) error {
	const batchSize = 1000

	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	page := 1
	perPage := batchSize
	findFilter.Page = &page
	findFilter.PerPage = &perPage

	for more := true; more; {
		if job.IsCancelled(ctx) {
			return nil
		}

		galleries, _, err := reader.Query(ctx, galleryFilter, findFilter)
		if err != nil {
			return fmt.Errorf("error querying for galleries: %w", err)
		}

		for _, g := range galleries {
			if err := fn(g); err != nil {
				return err
			}
		}

		if len(galleries) != batchSize {
			more = false
		} else {
			*findFilter.Page++
		}
	}

	return nil
}

func (j *IdentifyJob) identifyScene(ctx context.Context, s *models.Scene, sources []identify.ScraperSource) {
	if job.IsCancelled(ctx) {
		return
	}

	var taskError error
	j.progress.ExecuteTask("Identifying "+s.Path, func() {
		r := j.repository
		task := identify.SceneIdentifier{
			TxnManager:         r.TxnManager,
			SceneReaderUpdater: r.Scene,
			StudioReaderWriter: r.Studio,
			PerformerCreator:   r.Performer,
			TagFinderCreator:   r.Tag,

			DefaultOptions:              j.input.Options,
			Sources:                     sources,
			SceneUpdatePostHookExecutor: j.postHookExecutor,
		}

		taskError = task.Identify(ctx, s)
	})

	if taskError != nil {
		logger.Errorf("Error encountered identifying %s: %v", s.Path, taskError)
	}

	j.progress.Increment()
}

func (j *IdentifyJob) identifyGallery(ctx context.Context, g *models.Gallery, sources []identify.GalleryScraperSource) {
	if job.IsCancelled(ctx) {
		return
	}

	var taskError error
	j.progress.ExecuteTask("Identifying "+g.DisplayName(), func() {
		r := j.repository
		task := identify.GalleryIdentifier{
			TxnManager:           r.TxnManager,
			GalleryReaderUpdater: r.Gallery,
			StudioReaderWriter:   r.Studio,
			PerformerCreator:     r.Performer,
			TagFinderCreator:     r.Tag,

			DefaultOptions:   j.input.Options,
			Sources:          sources,
			PostHookExecutor: j.galleryPostHookExecutor,
		}

		taskError = task.Identify(ctx, g)
	})

	if taskError != nil {
		logger.Errorf("Error encountered identifying %s: %v", g.DisplayName(), taskError)
	}

	j.progress.Increment()
}

func (j *IdentifyJob) getSources() ([]identify.ScraperSource, error) {
	var ret []identify.ScraperSource
	for _, source := range j.input.Sources {
		// get scraper source
		stashBox, err := j.getStashBox(source.Source)
		if err != nil {
			return nil, err
		}

		var src identify.ScraperSource
		if stashBox != nil {
			matcher := match.SceneRelationships{
				PerformerFinder: instance.Repository.Performer,
				TagFinder:       instance.Repository.Tag,
				StudioFinder:    instance.Repository.Studio,
			}

			src = identify.ScraperSource{
				Name: "stash-box: " + stashBox.Endpoint,
				Scraper: stashboxSource{
					Client:                 stashbox.NewClient(*stashBox, stashbox.ExcludeTagPatterns(instance.Config.GetScraperExcludeTagPatterns())),
					endpoint:               stashBox.Endpoint,
					txnManager:             instance.Repository.TxnManager,
					sceneFingerprintGetter: instance.SceneService,
					matcher:                matcher,
				},
				RemoteSite: stashBox.Endpoint,
			}
		} else {
			scraperID := *source.Source.ScraperID
			s := instance.ScraperCache.GetScraper(scraperID)
			if s == nil {
				return nil, fmt.Errorf("%w: scraper with id %q", models.ErrNotFound, scraperID)
			}
			src = identify.ScraperSource{
				Name: s.Name,
				Scraper: scraperSource{
					cache:     instance.ScraperCache,
					scraperID: scraperID,
				},
			}
		}

		src.Options = source.Options
		ret = append(ret, src)
	}

	return ret, nil
}

func (j *IdentifyJob) getGallerySources() ([]identify.GalleryScraperSource, error) {
	var ret []identify.GalleryScraperSource
	for _, source := range j.input.Sources {
		// skip stash-box sources for galleries
		stashBox, err := j.getStashBox(source.Source)
		if err != nil {
			return nil, err
		}

		if stashBox != nil {
			logger.Warnf("Skipping stash-box source %s for gallery identify: stash-box gallery scraping not supported", stashBox.Endpoint)
			continue
		}

		// must be a scraper
		if source.Source.ScraperID == nil {
			return nil, fmt.Errorf("source must have scraper_id for gallery identify")
		}
		scraperID := *source.Source.ScraperID
		s := instance.ScraperCache.GetScraper(scraperID)
		if s == nil {
			return nil, fmt.Errorf("%w: scraper with id %q", models.ErrNotFound, scraperID)
		}

		if !supportsGalleryFragment(s) {
			logger.Warnf("Skipping scraper %s for gallery identify: gallery fragment scraping not supported", s.Name)
			continue
		}

		src := identify.GalleryScraperSource{
			Name: s.Name,
			Scraper: galleryScraperSource{
				cache:     instance.ScraperCache,
				scraperID: scraperID,
			},
		}

		src.Options = source.Options
		ret = append(ret, src)
	}

	return ret, nil
}

func supportsGalleryFragment(s *scraper.Scraper) bool {
	if s == nil || s.Gallery == nil {
		return false
	}

	return slices.Contains(s.Gallery.SupportedScrapes, scraper.ScrapeTypeFragment)
}

func (j *IdentifyJob) getStashBox(src *scraper.Source) (*models.StashBox, error) {
	if src.ScraperID != nil {
		return nil, nil
	}

	// must be stash-box
	if src.StashBoxIndex == nil && src.StashBoxEndpoint == nil {
		return nil, fmt.Errorf("%w: stash_box_index or stash_box_endpoint or scraper_id must be set", ErrInput)
	}

	return resolveStashBox(j.stashBoxes, *src)
}

func resolveStashBox(sb []*models.StashBox, source scraper.Source) (*models.StashBox, error) {
	if source.StashBoxIndex != nil {
		index := source.StashBoxIndex
		if *index < 0 || *index >= len(sb) {
			return nil, fmt.Errorf("%w: invalid stash_box_index: %d", models.ErrScraperSource, index)
		}

		return sb[*index], nil
	}

	if source.StashBoxEndpoint != nil {
		var ret *models.StashBox
		endpoint := *source.StashBoxEndpoint
		for _, b := range sb {
			if strings.EqualFold(endpoint, b.Endpoint) {
				ret = b
			}
		}

		if ret == nil {
			return nil, fmt.Errorf(`%w: stash-box with endpoint "%s"`, models.ErrNotFound, endpoint)
		}

		return ret, nil
	}

	// neither stash-box inputs were provided, so assume it is a scraper

	return nil, nil
}

type stashboxSource struct {
	*stashbox.Client
	endpoint string

	txnManager             models.TxnManager
	sceneFingerprintGetter sceneFingerprintGetter
	matcher                match.SceneRelationships
}

type sceneFingerprintGetter interface {
	GetScenesFingerprints(ctx context.Context, ids []int) ([]models.Fingerprints, error)
}

func (s stashboxSource) ScrapeScenes(ctx context.Context, sceneID int) ([]*models.ScrapedScene, error) {
	var fps []models.Fingerprints
	if err := txn.WithReadTxn(ctx, s.txnManager, func(ctx context.Context) error {
		var err error
		fps, err = s.sceneFingerprintGetter.GetScenesFingerprints(ctx, []int{sceneID})
		return err
	}); err != nil {
		return nil, fmt.Errorf("error getting scene fingerprints: %w", err)
	}

	results, err := s.FindSceneByFingerprints(ctx, fps[0])
	if err != nil {
		return nil, fmt.Errorf("error querying stash-box using scene ID %d: %w", sceneID, err)
	}

	if err := txn.WithReadTxn(ctx, s.txnManager, func(ctx context.Context) error {
		for _, ret := range results {
			if err := s.matcher.MatchRelationships(ctx, ret, s.endpoint); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("error matching scene relationships: %w", err)
	}

	if len(results) > 0 {
		return results, nil
	}

	return nil, nil
}

func (s stashboxSource) String() string {
	return fmt.Sprintf("stash-box %s", s.endpoint)
}

type scraperSource struct {
	cache     *scraper.Cache
	scraperID string
}

func (s scraperSource) ScrapeScenes(ctx context.Context, sceneID int) ([]*models.ScrapedScene, error) {
	content, err := s.cache.ScrapeID(ctx, s.scraperID, sceneID, scraper.ScrapeContentTypeScene)
	if err != nil {
		return nil, err
	}

	// don't try to convert nil return value
	if content == nil {
		return nil, nil
	}

	if sceneResult, ok := content.(models.ScrapedScene); ok {
		return []*models.ScrapedScene{&sceneResult}, nil
	}

	return nil, errors.New("could not convert content to scene")
}

func (s scraperSource) String() string {
	return fmt.Sprintf("scraper %s", s.scraperID)
}

type galleryScraperSource struct {
	cache     *scraper.Cache
	scraperID string
}

func (s galleryScraperSource) ScrapeGalleries(ctx context.Context, galleryID int) ([]*models.ScrapedGallery, error) {
	content, err := s.cache.ScrapeID(ctx, s.scraperID, galleryID, scraper.ScrapeContentTypeGallery)
	if err != nil {
		return nil, err
	}

	// don't try to convert nil return value
	if content == nil {
		return nil, nil
	}

	if galleryResult, ok := content.(models.ScrapedGallery); ok {
		return []*models.ScrapedGallery{&galleryResult}, nil
	}

	return nil, errors.New("could not convert content to gallery")
}

func (s galleryScraperSource) String() string {
	return fmt.Sprintf("scraper %s", s.scraperID)
}
