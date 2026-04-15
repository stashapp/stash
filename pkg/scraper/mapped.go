package scraper

import (
	"context"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type mappedQuery interface {
	runQuery(selector string) ([]string, error)
	getType() QueryType
	setType(QueryType)
	subScrape(ctx context.Context, value string) mappedQuery
	getURL() string
}

type mappedScrapers map[string]mappedScraper

type mappedScraper struct {
	Common    commonMappedConfig            `yaml:"common"`
	Scene     *mappedSceneScraperConfig     `yaml:"scene"`
	Gallery   *mappedGalleryScraperConfig   `yaml:"gallery"`
	Image     *mappedImageScraperConfig     `yaml:"image"`
	Performer *mappedPerformerScraperConfig `yaml:"performer"`
	Group     *mappedMovieScraperConfig     `yaml:"group"`

	// deprecated
	Movie *mappedMovieScraperConfig `yaml:"movie"`
}

func urlsIsMulti(key string) bool {
	return key == "URLs"
}

func (s mappedScraper) scrapePerformer(ctx context.Context, q mappedQuery) (*models.ScrapedPerformer, error) {
	var ret *models.ScrapedPerformer

	performerMap := s.Performer
	if performerMap == nil {
		return nil, nil
	}

	performerTagsMap := performerMap.Tags

	results := performerMap.process(ctx, q, s.Common, urlsIsMulti)

	// now apply the tags
	var tagResults mappedResults

	if performerTagsMap != nil {
		logger.Debug(`Processing performer tags:`)
		tagResults = performerTagsMap.process(ctx, q, s.Common, nil)
	}

	if len(results) == 0 {
		return nil, nil
	}

	if len(results) > 0 {
		ret = results[0].scrapedPerformer()
		ret.Tags = tagResults.scrapedTags()
	}

	return ret, nil
}

func (s mappedScraper) scrapePerformers(ctx context.Context, q mappedQuery) ([]*models.ScrapedPerformer, error) {
	performerMap := s.Performer
	if performerMap == nil {
		return nil, nil
	}

	// isMulti is nil because it will behave incorrect when scraping multiple performers
	results := performerMap.process(ctx, q, s.Common, nil)
	return results.scrapedPerformers(), nil
}

// processSceneRelationships sets the relationships on the models.ScrapedScene. It returns true if any relationships were set.
func (s mappedScraper) processSceneRelationships(ctx context.Context, q mappedQuery, resultIndex int, ret *models.ScrapedScene) bool {
	sceneScraperConfig := s.Scene

	scenePerformersMap := sceneScraperConfig.Performers
	sceneTagsMap := sceneScraperConfig.Tags
	sceneStudioMap := sceneScraperConfig.Studio
	sceneMoviesMap := sceneScraperConfig.Movies
	sceneGroupsMap := sceneScraperConfig.Groups

	ret.Performers = s.processPerformers(ctx, scenePerformersMap, q)

	if sceneTagsMap != nil {
		logger.Debug(`Processing scene tags:`)

		ret.Tags = sceneTagsMap.process(ctx, q, s.Common, nil).scrapedTags()
	}

	if sceneStudioMap != nil {
		logger.Debug(`Processing scene studio:`)
		studioResults := sceneStudioMap.process(ctx, q, s.Common, nil)

		if len(studioResults) > 0 && resultIndex < len(studioResults) {
			// when doing a `search` scrape get the related studio
			studio := studioResults[resultIndex].scrapedStudio()
			ret.Studio = studio
		}
	}

	if sceneMoviesMap != nil {
		logger.Debug(`Processing scene movies:`)
		ret.Movies = sceneMoviesMap.process(ctx, q, s.Common, nil).scrapedMovies()
	}

	if sceneGroupsMap != nil {
		logger.Debug(`Processing scene groups:`)
		ret.Groups = sceneGroupsMap.process(ctx, q, s.Common, nil).scrapedGroups()
	}

	return len(ret.Performers) > 0 || len(ret.Tags) > 0 || ret.Studio != nil || len(ret.Movies) > 0 || len(ret.Groups) > 0
}

func (s mappedScraper) processPerformers(ctx context.Context, performersMap mappedPerformerScraperConfig, q mappedQuery) []*models.ScrapedPerformer {
	var ret []*models.ScrapedPerformer

	// now apply the performers and tags
	if performersMap.mappedConfig != nil {
		logger.Debug(`Processing performers:`)
		// isMulti is nil because it will behave incorrect when scraping multiple performers
		performerResults := performersMap.process(ctx, q, s.Common, nil)

		scenePerformerTagsMap := performersMap.Tags

		// process performer tags once
		var performerTagResults mappedResults
		if scenePerformerTagsMap != nil {
			performerTagResults = scenePerformerTagsMap.process(ctx, q, s.Common, nil)
		}

		for _, p := range performerResults {
			performer := p.scrapedPerformer()

			for _, p := range performerTagResults {
				tag := p.scrapedTag()
				performer.Tags = append(performer.Tags, tag)
			}

			ret = append(ret, performer)
		}
	}

	return ret
}

func (s mappedScraper) scrapeScenes(ctx context.Context, q mappedQuery) ([]*models.ScrapedScene, error) {
	var ret []*models.ScrapedScene

	sceneScraperConfig := s.Scene
	sceneMap := sceneScraperConfig.mappedConfig
	if sceneMap == nil {
		return nil, nil
	}

	logger.Debug(`Processing scenes:`)
	// urlsIsMulti is nil because it will behave incorrect when scraping multiple scenes
	results := sceneMap.process(ctx, q, s.Common, nil)
	for i, r := range results {
		logger.Debug(`Processing scene:`)

		thisScene := r.scrapedScene()
		s.processSceneRelationships(ctx, q, i, thisScene)
		ret = append(ret, thisScene)
	}

	return ret, nil
}

func (s mappedScraper) scrapeScene(ctx context.Context, q mappedQuery) (*models.ScrapedScene, error) {
	sceneScraperConfig := s.Scene
	if sceneScraperConfig == nil {
		return nil, nil
	}

	sceneMap := sceneScraperConfig.mappedConfig

	logger.Debug(`Processing scene:`)
	results := sceneMap.process(ctx, q, s.Common, urlsIsMulti)

	var ret *models.ScrapedScene
	if len(results) > 0 {
		ret = results[0].scrapedScene()
	}
	hasRelationships := s.processSceneRelationships(ctx, q, 0, ret)

	// #3953 - process only returns results if the non-relationship fields are
	// populated
	// only return if we have results or relationships
	if len(results) > 0 || hasRelationships {
		return ret, nil
	}

	return nil, nil
}

func (s mappedScraper) scrapeImage(ctx context.Context, q mappedQuery) (*models.ScrapedImage, error) {
	var ret models.ScrapedImage

	imageScraperConfig := s.Image
	if imageScraperConfig == nil {
		return nil, nil
	}

	imageMap := imageScraperConfig.mappedConfig

	imagePerformersMap := imageScraperConfig.Performers
	imageTagsMap := imageScraperConfig.Tags
	imageStudioMap := imageScraperConfig.Studio

	logger.Debug(`Processing image:`)
	results := imageMap.process(ctx, q, s.Common, urlsIsMulti)

	if len(results) > 0 {
		ret = *results[0].scrapedImage()
	}

	// now apply the performers and tags
	if imagePerformersMap != nil {
		logger.Debug(`Processing image performers:`)
		ret.Performers = imagePerformersMap.process(ctx, q, s.Common, nil).scrapedPerformers()
	}

	if imageTagsMap != nil {
		logger.Debug(`Processing image tags:`)
		ret.Tags = imageTagsMap.process(ctx, q, s.Common, nil).scrapedTags()
	}

	if imageStudioMap != nil {
		logger.Debug(`Processing image studio:`)
		studioResults := imageStudioMap.process(ctx, q, s.Common, nil)

		if len(studioResults) > 0 {
			ret.Studio = studioResults[0].scrapedStudio()
		}
	}

	// if no basic fields are populated, and no relationships, then return nil
	if len(results) == 0 && len(ret.Performers) == 0 && len(ret.Tags) == 0 && ret.Studio == nil {
		return nil, nil
	}

	return &ret, nil
}

func (s mappedScraper) scrapeGallery(ctx context.Context, q mappedQuery) (*models.ScrapedGallery, error) {
	var ret models.ScrapedGallery

	galleryScraperConfig := s.Gallery
	if galleryScraperConfig == nil {
		return nil, nil
	}

	galleryMap := galleryScraperConfig.mappedConfig

	galleryPerformersMap := galleryScraperConfig.Performers
	galleryTagsMap := galleryScraperConfig.Tags
	galleryStudioMap := galleryScraperConfig.Studio

	logger.Debug(`Processing gallery:`)
	results := galleryMap.process(ctx, q, s.Common, urlsIsMulti)

	if len(results) > 0 {
		ret = *results[0].scrapedGallery()
	}

	// now apply the performers and tags
	if galleryPerformersMap != nil {
		logger.Debug(`Processing gallery performers:`)
		performerResults := galleryPerformersMap.process(ctx, q, s.Common, urlsIsMulti)

		ret.Performers = performerResults.scrapedPerformers()
	}

	if galleryTagsMap != nil {
		logger.Debug(`Processing gallery tags:`)
		tagResults := galleryTagsMap.process(ctx, q, s.Common, nil)
		ret.Tags = tagResults.scrapedTags()
	}

	if galleryStudioMap != nil {
		logger.Debug(`Processing gallery studio:`)
		studioResults := galleryStudioMap.process(ctx, q, s.Common, nil)

		if len(studioResults) > 0 {
			ret.Studio = studioResults[0].scrapedStudio()
		}
	}

	// if no basic fields are populated, and no relationships, then return nil
	if len(results) == 0 && len(ret.Performers) == 0 && len(ret.Tags) == 0 && ret.Studio == nil {
		return nil, nil
	}

	return &ret, nil
}

func (s mappedScraper) scrapeGroup(ctx context.Context, q mappedQuery) (*models.ScrapedGroup, error) {
	var ret models.ScrapedGroup

	// try group scraper first, falling back to movie
	groupScraperConfig := s.Group

	if groupScraperConfig == nil {
		groupScraperConfig = s.Movie
	}
	if groupScraperConfig == nil {
		return nil, nil
	}

	groupMap := groupScraperConfig.mappedConfig

	groupStudioMap := groupScraperConfig.Studio
	groupTagsMap := groupScraperConfig.Tags

	results := groupMap.process(ctx, q, s.Common, urlsIsMulti)

	if len(results) > 0 {
		ret = *results[0].scrapedGroup()
	}

	if groupStudioMap != nil {
		logger.Debug(`Processing group studio:`)
		studioResults := groupStudioMap.process(ctx, q, s.Common, nil)

		if len(studioResults) > 0 {
			ret.Studio = studioResults[0].scrapedStudio()
		}
	}

	// now apply the tags
	if groupTagsMap != nil {
		logger.Debug(`Processing group tags:`)
		tagResults := groupTagsMap.process(ctx, q, s.Common, nil)

		ret.Tags = tagResults.scrapedTags()
	}

	if len(results) == 0 && ret.Studio == nil && len(ret.Tags) == 0 {
		return nil, nil
	}

	return &ret, nil
}
