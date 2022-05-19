package manager

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/stashapp/stash/internal/identify"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/scraper/stashbox"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/txn"
)

var ErrInput = errors.New("invalid request input")

type IdentifyJob struct {
	postHookExecutor identify.SceneUpdatePostHookExecutor
	input            identify.Options

	stashBoxes []*models.StashBox
	progress   *job.Progress
}

func CreateIdentifyJob(input identify.Options) *IdentifyJob {
	return &IdentifyJob{
		postHookExecutor: instance.PluginCache,
		input:            input,
		stashBoxes:       instance.Config.GetStashBoxes(),
	}
}

func (j *IdentifyJob) Execute(ctx context.Context, progress *job.Progress) {
	j.progress = progress

	// if no sources provided - just return
	if len(j.input.Sources) == 0 {
		return
	}

	sources, err := j.getSources()
	if err != nil {
		logger.Error(err)
		return
	}

	// if scene ids provided, use those
	// otherwise, batch query for all scenes - ordering by path
	if err := txn.WithTxn(ctx, instance.Repository, func(ctx context.Context) error {
		if len(j.input.SceneIDs) == 0 {
			return j.identifyAllScenes(ctx, sources)
		}

		sceneIDs, err := stringslice.StringSliceToIntSlice(j.input.SceneIDs)
		if err != nil {
			return fmt.Errorf("invalid scene IDs: %w", err)
		}

		progress.SetTotal(len(sceneIDs))
		for _, id := range sceneIDs {
			if job.IsCancelled(ctx) {
				break
			}

			// find the scene
			var err error
			scene, err := instance.Repository.Scene.Find(ctx, id)
			if err != nil {
				return fmt.Errorf("error finding scene with id %d: %w", id, err)
			}

			if scene == nil {
				return fmt.Errorf("%w: scene with id %d", models.ErrNotFound, id)
			}

			j.identifyScene(ctx, scene, sources)
		}

		return nil
	}); err != nil {
		logger.Errorf("Error encountered while identifying scenes: %v", err)
	}
}

func (j *IdentifyJob) identifyAllScenes(ctx context.Context, sources []identify.ScraperSource) error {
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
	countResult, err := instance.Repository.Scene.Query(ctx, models.SceneQueryOptions{
		QueryOptions: models.QueryOptions{
			FindFilter: findFilter,
			Count:      true,
		},
		SceneFilter: sceneFilter,
	})
	if err != nil {
		return fmt.Errorf("error getting scene count: %w", err)
	}

	j.progress.SetTotal(countResult.Count)

	return scene.BatchProcess(ctx, instance.Repository.Scene, sceneFilter, findFilter, func(scene *models.Scene) error {
		if job.IsCancelled(ctx) {
			return nil
		}

		j.identifyScene(ctx, scene, sources)
		return nil
	})
}

func (j *IdentifyJob) identifyScene(ctx context.Context, s *models.Scene, sources []identify.ScraperSource) {
	if job.IsCancelled(ctx) {
		return
	}

	var taskError error
	j.progress.ExecuteTask("Identifying "+s.Path, func() {
		task := identify.SceneIdentifier{
			SceneReaderUpdater: instance.Repository.Scene,
			StudioCreator:      instance.Repository.Studio,
			PerformerCreator:   instance.Repository.Performer,
			TagCreator:         instance.Repository.Tag,

			DefaultOptions: j.input.Options,
			Sources:        sources,
			ScreenshotSetter: &scene.PathsScreenshotSetter{
				Paths:               instance.Paths,
				FileNamingAlgorithm: instance.Config.GetVideoFileNamingAlgorithm(),
			},
			SceneUpdatePostHookExecutor: j.postHookExecutor,
		}

		taskError = task.Identify(ctx, instance.Repository, s)
	})

	if taskError != nil {
		logger.Errorf("Error encountered identifying %s: %v", s.Path, taskError)
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
			src = identify.ScraperSource{
				Name: "stash-box: " + stashBox.Endpoint,
				Scraper: stashboxSource{
					stashbox.NewClient(*stashBox, instance.Repository, stashbox.Repository{
						Scene:     instance.Repository.Scene,
						Performer: instance.Repository.Performer,
						Tag:       instance.Repository.Tag,
						Studio:    instance.Repository.Studio,
					}),
					stashBox.Endpoint,
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
}

func (s stashboxSource) ScrapeScene(ctx context.Context, sceneID int) (*scraper.ScrapedScene, error) {
	results, err := s.FindStashBoxSceneByFingerprints(ctx, sceneID)
	if err != nil {
		return nil, fmt.Errorf("error querying stash-box using scene ID %d: %w", sceneID, err)
	}

	if len(results) > 0 {
		return results[0], nil
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

func (s scraperSource) ScrapeScene(ctx context.Context, sceneID int) (*scraper.ScrapedScene, error) {
	content, err := s.cache.ScrapeID(ctx, s.scraperID, sceneID, scraper.ScrapeContentTypeScene)
	if err != nil {
		return nil, err
	}

	// don't try to convert nil return value
	if content == nil {
		return nil, nil
	}

	if scene, ok := content.(scraper.ScrapedScene); ok {
		return &scene, nil
	}

	return nil, errors.New("could not convert content to scene")
}

func (s scraperSource) String() string {
	return fmt.Sprintf("scraper %s", s.scraperID)
}
