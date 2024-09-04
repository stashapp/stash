package api

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/scraper/stashbox"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

func (r *queryResolver) ScrapeURL(ctx context.Context, url string, ty scraper.ScrapeContentType) (scraper.ScrapedContent, error) {
	return r.scraperCache().ScrapeURL(ctx, url, ty)
}

func (r *queryResolver) ListScrapers(ctx context.Context, types []scraper.ScrapeContentType) ([]*scraper.Scraper, error) {
	return r.scraperCache().ListScrapers(types), nil
}

func (r *queryResolver) ScrapePerformerURL(ctx context.Context, url string) (*models.ScrapedPerformer, error) {
	content, err := r.scraperCache().ScrapeURL(ctx, url, scraper.ScrapeContentTypePerformer)
	if err != nil {
		return nil, err
	}

	return marshalScrapedPerformer(content)
}

func (r *queryResolver) ScrapeSceneQuery(ctx context.Context, scraperID string, query string) ([]*scraper.ScrapedScene, error) {
	if query == "" {
		return nil, nil
	}

	content, err := r.scraperCache().ScrapeName(ctx, scraperID, query, scraper.ScrapeContentTypeScene)
	if err != nil {
		return nil, err
	}

	ret, err := marshalScrapedScenes(content)
	if err != nil {
		return nil, err
	}

	filterSceneTags(ret)
	return ret, nil
}

func compileRegexps(patterns []string) []*regexp.Regexp {
	excludePatterns := patterns
	var excludeRegexps []*regexp.Regexp

	for _, excludePattern := range excludePatterns {
		reg, err := regexp.Compile(strings.ToLower(excludePattern))
		if err != nil {
			logger.Errorf("Invalid tag exclusion pattern: %v", err)
		} else {
			excludeRegexps = append(excludeRegexps, reg)
		}
	}

	return excludeRegexps
}

// filterSceneTags removes tags matching excluded tag patterns from the provided scraped scenes
func filterTags(excludeRegexps []*regexp.Regexp, tags []*models.ScrapedTag) (newTags []*models.ScrapedTag, ignoredTags []string) {
	if len(excludeRegexps) == 0 {
		return tags, nil
	}

	for _, t := range tags {
		ignore := false
		for _, reg := range excludeRegexps {
			if reg.MatchString(strings.ToLower(t.Name)) {
				ignore = true
				ignoredTags = sliceutil.AppendUnique(ignoredTags, t.Name)
				break
			}
		}

		if !ignore {
			newTags = append(newTags, t)
		}
	}

	return
}

// filterSceneTags removes tags matching excluded tag patterns from the provided scraped scenes
func filterSceneTags(scenes []*scraper.ScrapedScene) {
	excludeRegexps := compileRegexps(manager.GetInstance().Config.GetScraperExcludeTagPatterns())

	var ignoredTags []string

	for _, s := range scenes {
		var ignored []string
		s.Tags, ignored = filterTags(excludeRegexps, s.Tags)
		ignoredTags = sliceutil.AppendUniques(ignoredTags, ignored)
	}

	if len(ignoredTags) > 0 {
		logger.Debugf("Scraping ignored tags: %s", strings.Join(ignoredTags, ", "))
	}
}

// filterGalleryTags removes tags matching excluded tag patterns from the provided scraped galleries
func filterGalleryTags(g []*scraper.ScrapedGallery) {
	excludeRegexps := compileRegexps(manager.GetInstance().Config.GetScraperExcludeTagPatterns())

	var ignoredTags []string

	for _, s := range g {
		var ignored []string
		s.Tags, ignored = filterTags(excludeRegexps, s.Tags)
		ignoredTags = sliceutil.AppendUniques(ignoredTags, ignored)
	}

	if len(ignoredTags) > 0 {
		logger.Debugf("Scraping ignored tags: %s", strings.Join(ignoredTags, ", "))
	}
}

// filterGalleryTags removes tags matching excluded tag patterns from the provided scraped galleries
func filterPerformerTags(p []*models.ScrapedPerformer) {
	excludeRegexps := compileRegexps(manager.GetInstance().Config.GetScraperExcludeTagPatterns())

	var ignoredTags []string

	for _, s := range p {
		var ignored []string
		s.Tags, ignored = filterTags(excludeRegexps, s.Tags)
		ignoredTags = sliceutil.AppendUniques(ignoredTags, ignored)
	}

	if len(ignoredTags) > 0 {
		logger.Debugf("Scraping ignored tags: %s", strings.Join(ignoredTags, ", "))
	}
}

// filterGroupTags removes tags matching excluded tag patterns from the provided scraped movies
func filterGroupTags(p []*models.ScrapedMovie) {
	excludeRegexps := compileRegexps(manager.GetInstance().Config.GetScraperExcludeTagPatterns())

	var ignoredTags []string

	for _, s := range p {
		var ignored []string
		s.Tags, ignored = filterTags(excludeRegexps, s.Tags)
		ignoredTags = sliceutil.AppendUniques(ignoredTags, ignored)
	}

	if len(ignoredTags) > 0 {
		logger.Debugf("Scraping ignored tags: %s", strings.Join(ignoredTags, ", "))
	}
}

func (r *queryResolver) ScrapeSceneURL(ctx context.Context, url string) (*scraper.ScrapedScene, error) {
	content, err := r.scraperCache().ScrapeURL(ctx, url, scraper.ScrapeContentTypeScene)
	if err != nil {
		return nil, err
	}

	ret, err := marshalScrapedScene(content)
	if err != nil {
		return nil, err
	}

	if ret != nil {
		filterSceneTags([]*scraper.ScrapedScene{ret})
	}

	return ret, nil
}

func (r *queryResolver) ScrapeGalleryURL(ctx context.Context, url string) (*scraper.ScrapedGallery, error) {
	content, err := r.scraperCache().ScrapeURL(ctx, url, scraper.ScrapeContentTypeGallery)
	if err != nil {
		return nil, err
	}

	ret, err := marshalScrapedGallery(content)
	if err != nil {
		return nil, err
	}

	if ret != nil {
		filterGalleryTags([]*scraper.ScrapedGallery{ret})
	}

	return ret, nil
}

func (r *queryResolver) ScrapeMovieURL(ctx context.Context, url string) (*models.ScrapedMovie, error) {
	content, err := r.scraperCache().ScrapeURL(ctx, url, scraper.ScrapeContentTypeMovie)
	if err != nil {
		return nil, err
	}

	ret, err := marshalScrapedMovie(content)
	if err != nil {
		return nil, err
	}

	filterGroupTags([]*models.ScrapedMovie{ret})

	return ret, nil
}

func (r *queryResolver) ScrapeGroupURL(ctx context.Context, url string) (*models.ScrapedGroup, error) {
	content, err := r.scraperCache().ScrapeURL(ctx, url, scraper.ScrapeContentTypeMovie)
	if err != nil {
		return nil, err
	}

	ret, err := marshalScrapedMovie(content)
	if err != nil {
		return nil, err
	}

	filterGroupTags([]*models.ScrapedMovie{ret})

	// convert to scraped group
	group := &models.ScrapedGroup{
		StoredID:   ret.StoredID,
		Name:       ret.Name,
		Aliases:    ret.Aliases,
		Duration:   ret.Duration,
		Date:       ret.Date,
		Rating:     ret.Rating,
		Director:   ret.Director,
		URLs:       ret.URLs,
		Synopsis:   ret.Synopsis,
		Studio:     ret.Studio,
		Tags:       ret.Tags,
		FrontImage: ret.FrontImage,
		BackImage:  ret.BackImage,
	}

	return group, nil
}

func (r *queryResolver) ScrapeSingleScene(ctx context.Context, source scraper.Source, input ScrapeSingleSceneInput) ([]*scraper.ScrapedScene, error) {
	var ret []*scraper.ScrapedScene

	var sceneID int
	if input.SceneID != nil {
		var err error
		sceneID, err = strconv.Atoi(*input.SceneID)
		if err != nil {
			return nil, fmt.Errorf("%w: sceneID is not an integer: '%s'", ErrInput, *input.SceneID)
		}
	}

	switch {
	case source.ScraperID != nil:
		var err error
		var c scraper.ScrapedContent
		var content []scraper.ScrapedContent

		switch {
		case input.SceneID != nil:
			c, err = r.scraperCache().ScrapeID(ctx, *source.ScraperID, sceneID, scraper.ScrapeContentTypeScene)
			if c != nil {
				content = []scraper.ScrapedContent{c}
			}
		case input.SceneInput != nil:
			c, err = r.scraperCache().ScrapeFragment(ctx, *source.ScraperID, scraper.Input{Scene: input.SceneInput})
			if c != nil {
				content = []scraper.ScrapedContent{c}
			}
		case input.Query != nil:
			content, err = r.scraperCache().ScrapeName(ctx, *source.ScraperID, *input.Query, scraper.ScrapeContentTypeScene)
		default:
			err = fmt.Errorf("%w: scene_id, scene_input, or query must be set", ErrInput)
		}

		if err != nil {
			return nil, err
		}

		ret, err = marshalScrapedScenes(content)
		if err != nil {
			return nil, err
		}
	case source.StashBoxIndex != nil || source.StashBoxEndpoint != nil:
		b, err := resolveStashBox(source.StashBoxIndex, source.StashBoxEndpoint)
		if err != nil {
			return nil, err
		}

		client := r.newStashBoxClient(*b)

		switch {
		case input.SceneID != nil:
			ret, err = client.FindStashBoxSceneByFingerprints(ctx, sceneID)
		case input.Query != nil:
			ret, err = client.QueryStashBoxScene(ctx, *input.Query)
		default:
			return nil, fmt.Errorf("%w: scene_id or query must be set", ErrInput)
		}

		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%w: scraper_id or stash_box_index must be set", ErrInput)
	}

	filterSceneTags(ret)

	return ret, nil
}

func (r *queryResolver) ScrapeMultiScenes(ctx context.Context, source scraper.Source, input ScrapeMultiScenesInput) ([][]*scraper.ScrapedScene, error) {
	if source.ScraperID != nil {
		return nil, ErrNotImplemented
	} else if source.StashBoxIndex != nil || source.StashBoxEndpoint != nil {
		b, err := resolveStashBox(source.StashBoxIndex, source.StashBoxEndpoint)
		if err != nil {
			return nil, err
		}

		client := r.newStashBoxClient(*b)

		sceneIDs, err := stringslice.StringSliceToIntSlice(input.SceneIds)
		if err != nil {
			return nil, err
		}

		return client.FindStashBoxScenesByFingerprints(ctx, sceneIDs)
	}

	return nil, errors.New("scraper_id or stash_box_index must be set")
}

func (r *queryResolver) ScrapeSingleStudio(ctx context.Context, source scraper.Source, input ScrapeSingleStudioInput) ([]*models.ScrapedStudio, error) {
	if source.StashBoxIndex != nil || source.StashBoxEndpoint != nil {
		b, err := resolveStashBox(source.StashBoxIndex, source.StashBoxEndpoint)
		if err != nil {
			return nil, err
		}

		client := r.newStashBoxClient(*b)

		var ret []*models.ScrapedStudio
		out, err := client.FindStashBoxStudio(ctx, *input.Query)

		if err != nil {
			return nil, err
		} else if out != nil {
			ret = append(ret, out)
		}

		if len(ret) > 0 {
			return ret, nil
		}

		return nil, nil
	}

	return nil, errors.New("stash_box_index must be set")
}

func (r *queryResolver) ScrapeSinglePerformer(ctx context.Context, source scraper.Source, input ScrapeSinglePerformerInput) ([]*models.ScrapedPerformer, error) {
	var ret []*models.ScrapedPerformer
	switch {
	case source.ScraperID != nil:
		switch {
		case input.PerformerInput != nil:
			performer, err := r.scraperCache().ScrapeFragment(ctx, *source.ScraperID, scraper.Input{Performer: input.PerformerInput})
			if err != nil {
				return nil, err
			}

			ret, err = marshalScrapedPerformers([]scraper.ScrapedContent{performer})
			if err != nil {
				return nil, err
			}
		case input.Query != nil:
			content, err := r.scraperCache().ScrapeName(ctx, *source.ScraperID, *input.Query, scraper.ScrapeContentTypePerformer)
			if err != nil {
				return nil, err
			}

			ret, err = marshalScrapedPerformers(content)
			if err != nil {
				return nil, err
			}
		default:
			return nil, ErrNotImplemented
		}
	case source.StashBoxIndex != nil || source.StashBoxEndpoint != nil:
		b, err := resolveStashBox(source.StashBoxIndex, source.StashBoxEndpoint)
		if err != nil {
			return nil, err
		}

		client := r.newStashBoxClient(*b)

		var res []*stashbox.StashBoxPerformerQueryResult
		switch {
		case input.PerformerID != nil:
			res, err = client.FindStashBoxPerformersByNames(ctx, []string{*input.PerformerID})
		case input.Query != nil:
			res, err = client.QueryStashBoxPerformer(ctx, *input.Query)
		default:
			return nil, ErrNotImplemented
		}

		if err != nil {
			return nil, err
		}

		if len(res) > 0 {
			ret = res[0].Results
		}
	default:
		return nil, errors.New("scraper_id or stash_box_index must be set")
	}

	filterPerformerTags(ret)

	return ret, nil
}

func (r *queryResolver) ScrapeMultiPerformers(ctx context.Context, source scraper.Source, input ScrapeMultiPerformersInput) ([][]*models.ScrapedPerformer, error) {
	if source.ScraperID != nil {
		return nil, ErrNotImplemented
	} else if source.StashBoxIndex != nil || source.StashBoxEndpoint != nil {
		b, err := resolveStashBox(source.StashBoxIndex, source.StashBoxEndpoint)
		if err != nil {
			return nil, err
		}

		client := r.newStashBoxClient(*b)

		return client.FindStashBoxPerformersByPerformerNames(ctx, input.PerformerIds)
	}

	return nil, errors.New("scraper_id or stash_box_index must be set")
}

func (r *queryResolver) ScrapeSingleGallery(ctx context.Context, source scraper.Source, input ScrapeSingleGalleryInput) ([]*scraper.ScrapedGallery, error) {
	var ret []*scraper.ScrapedGallery

	if source.StashBoxIndex != nil || source.StashBoxEndpoint != nil {
		return nil, ErrNotSupported
	}

	if source.ScraperID == nil {
		return nil, fmt.Errorf("%w: scraper_id must be set", ErrInput)
	}

	var c scraper.ScrapedContent

	switch {
	case input.GalleryID != nil:
		galleryID, err := strconv.Atoi(*input.GalleryID)
		if err != nil {
			return nil, fmt.Errorf("%w: gallery id is not an integer: '%s'", ErrInput, *input.GalleryID)
		}
		c, err = r.scraperCache().ScrapeID(ctx, *source.ScraperID, galleryID, scraper.ScrapeContentTypeGallery)
		if err != nil {
			return nil, err
		}
		ret, err = marshalScrapedGalleries([]scraper.ScrapedContent{c})
		if err != nil {
			return nil, err
		}
	case input.GalleryInput != nil:
		c, err := r.scraperCache().ScrapeFragment(ctx, *source.ScraperID, scraper.Input{Gallery: input.GalleryInput})
		if err != nil {
			return nil, err
		}
		ret, err = marshalScrapedGalleries([]scraper.ScrapedContent{c})
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrNotImplemented
	}

	filterGalleryTags(ret)
	return ret, nil
}

func (r *queryResolver) ScrapeSingleMovie(ctx context.Context, source scraper.Source, input ScrapeSingleMovieInput) ([]*models.ScrapedMovie, error) {
	return nil, ErrNotSupported
}

func (r *queryResolver) ScrapeSingleGroup(ctx context.Context, source scraper.Source, input ScrapeSingleGroupInput) ([]*models.ScrapedGroup, error) {
	return nil, ErrNotSupported
}
