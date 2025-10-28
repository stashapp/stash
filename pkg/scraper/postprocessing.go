package scraper

import (
	"context"
	"regexp"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/utils"
)

// postScrape handles post-processing of scraped content. If the content
// requires post-processing, this function fans out to the given content
// type and post-processes it.
func (c Cache) postScrape(ctx context.Context, content ScrapedContent, excludeTagRE []*regexp.Regexp) (_ ScrapedContent, ignoredTags []string, err error) {
	// Analyze the concrete type, call the right post-processing function
	switch v := content.(type) {
	case *models.ScrapedPerformer:
		if v != nil {
			return c.postScrapePerformer(ctx, *v, excludeTagRE)
		}
	case models.ScrapedPerformer:
		return c.postScrapePerformer(ctx, v, excludeTagRE)
	case *models.ScrapedScene:
		if v != nil {
			return c.postScrapeScene(ctx, *v, excludeTagRE)
		}
	case models.ScrapedScene:
		return c.postScrapeScene(ctx, v, excludeTagRE)
	case *models.ScrapedGallery:
		if v != nil {
			return c.postScrapeGallery(ctx, *v, excludeTagRE)
		}
	case models.ScrapedGallery:
		return c.postScrapeGallery(ctx, v, excludeTagRE)
	case *models.ScrapedImage:
		if v != nil {
			return c.postScrapeImage(ctx, *v, excludeTagRE)
		}
	case models.ScrapedImage:
		return c.postScrapeImage(ctx, v, excludeTagRE)
	case *models.ScrapedMovie:
		if v != nil {
			return c.postScrapeMovie(ctx, *v, excludeTagRE)
		}
	case models.ScrapedMovie:
		return c.postScrapeMovie(ctx, v, excludeTagRE)
	case *models.ScrapedGroup:
		if v != nil {
			return c.postScrapeGroup(ctx, *v, excludeTagRE)
		}
	case models.ScrapedGroup:
		return c.postScrapeGroup(ctx, v, excludeTagRE)
	}

	// If nothing matches, pass the content through
	return content, nil, nil
}

// postScrapeSingle handles post-processing of a single scraped content item.
// This is a convenience function that includes logging the ignored tags, as opposed to logging them in the caller.
func (c Cache) postScrapeSingle(ctx context.Context, content ScrapedContent) (ScrapedContent, error) {
	ret, ignoredTags, err := c.postScrape(ctx, content, c.compileExcludeTagPatterns())
	if err != nil {
		return nil, err
	}

	LogIgnoredTags(ignoredTags)
	return ret, nil
}

func (c Cache) postScrapePerformer(ctx context.Context, p models.ScrapedPerformer, excludeTagRE []*regexp.Regexp) (_ ScrapedContent, ignoredTags []string, err error) {
	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		tqb := r.TagFinder

		tags, err := postProcessTags(ctx, tqb, p.Tags)
		if err != nil {
			return err
		}
		p.Tags, ignoredTags = FilterTags(excludeTagRE, tags)

		return nil
	}); err != nil {
		return nil, nil, err
	}

	// post-process - set the image if applicable
	if err := setPerformerImage(ctx, c.client, &p, c.globalConfig); err != nil {
		logger.Warnf("Could not set image using URL %s: %s", *p.Image, err.Error())
	}

	p.Country = resolveCountryName(p.Country)

	// populate URL/URLs
	// if URLs are provided, only use those
	if len(p.URLs) > 0 {
		p.URL = &p.URLs[0]
	} else {
		urls := []string{}
		if p.URL != nil {
			urls = append(urls, *p.URL)
		}
		if p.Twitter != nil && *p.Twitter != "" {
			// handle twitter profile names
			u := utils.URLFromHandle(*p.Twitter, "https://twitter.com")
			urls = append(urls, u)
		}
		if p.Instagram != nil && *p.Instagram != "" {
			// handle instagram profile names
			u := utils.URLFromHandle(*p.Instagram, "https://instagram.com")
			urls = append(urls, u)
		}

		if len(urls) > 0 {
			p.URLs = urls
		}
	}

	return p, ignoredTags, nil
}

func (c Cache) postScrapeMovie(ctx context.Context, m models.ScrapedMovie, excludeTagRE []*regexp.Regexp) (_ ScrapedContent, ignoredTags []string, err error) {
	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		tqb := r.TagFinder
		tags, err := postProcessTags(ctx, tqb, m.Tags)
		if err != nil {
			return err
		}
		m.Tags, ignoredTags = FilterTags(excludeTagRE, tags)

		if m.Studio != nil {
			if err := match.ScrapedStudio(ctx, r.StudioFinder, m.Studio, ""); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, nil, err
	}

	// post-process - set the image if applicable
	if err := setMovieFrontImage(ctx, c.client, &m, c.globalConfig); err != nil {
		logger.Warnf("could not set front image using URL %s: %v", *m.FrontImage, err)
	}
	if err := setMovieBackImage(ctx, c.client, &m, c.globalConfig); err != nil {
		logger.Warnf("could not set back image using URL %s: %v", *m.BackImage, err)
	}

	return m, ignoredTags, nil
}

func (c Cache) postScrapeGroup(ctx context.Context, m models.ScrapedGroup, excludeTagRE []*regexp.Regexp) (_ ScrapedContent, ignoredTags []string, err error) {
	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		tqb := r.TagFinder
		tags, err := postProcessTags(ctx, tqb, m.Tags)
		if err != nil {
			return err
		}
		m.Tags, ignoredTags = FilterTags(excludeTagRE, tags)

		if m.Studio != nil {
			if err := match.ScrapedStudio(ctx, r.StudioFinder, m.Studio, ""); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, nil, err
	}

	// post-process - set the image if applicable
	if err := setGroupFrontImage(ctx, c.client, &m, c.globalConfig); err != nil {
		logger.Warnf("could not set front image using URL %s: %v", *m.FrontImage, err)
	}
	if err := setGroupBackImage(ctx, c.client, &m, c.globalConfig); err != nil {
		logger.Warnf("could not set back image using URL %s: %v", *m.BackImage, err)
	}

	return m, ignoredTags, nil
}

func (c Cache) postScrapeScenePerformer(ctx context.Context, p models.ScrapedPerformer, excludeTagRE []*regexp.Regexp) (ignoredTags []string, err error) {
	tqb := c.repository.TagFinder

	tags, err := postProcessTags(ctx, tqb, p.Tags)
	if err != nil {
		return nil, err
	}
	p.Tags = tags
	p.Tags, ignoredTags = FilterTags(excludeTagRE, tags)

	p.Country = resolveCountryName(p.Country)

	return ignoredTags, nil
}

func (c Cache) postScrapeScene(ctx context.Context, scene models.ScrapedScene, excludeTagRE []*regexp.Regexp) (_ ScrapedContent, ignoredTags []string, err error) {
	// set the URL/URLs field
	if scene.URL == nil && len(scene.URLs) > 0 {
		scene.URL = &scene.URLs[0]
	}
	if scene.URL != nil && len(scene.URLs) == 0 {
		scene.URLs = []string{*scene.URL}
	}

	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		pqb := r.PerformerFinder
		gqb := r.GroupFinder
		tqb := r.TagFinder
		sqb := r.StudioFinder

		for _, p := range scene.Performers {
			if p == nil {
				continue
			}

			thisIgnoredTags, err := c.postScrapeScenePerformer(ctx, *p, excludeTagRE)
			if err != nil {
				return err
			}

			if err := match.ScrapedPerformer(ctx, pqb, p, ""); err != nil {
				return err
			}

			ignoredTags = sliceutil.AppendUniques(ignoredTags, thisIgnoredTags)
		}

		for _, p := range scene.Movies {
			matchedID, err := match.ScrapedGroup(ctx, gqb, p.StoredID, p.Name)
			if err != nil {
				return err
			}

			if matchedID != nil {
				p.StoredID = matchedID
			}
		}

		for _, p := range scene.Groups {
			matchedID, err := match.ScrapedGroup(ctx, gqb, p.StoredID, p.Name)
			if err != nil {
				return err
			}

			if matchedID != nil {
				p.StoredID = matchedID
			}
		}

		// HACK - if movies was returned but not groups, add the groups from the movies
		// if groups was returned but not movies, add the movies from the groups for backward compatibility
		if len(scene.Movies) > 0 && len(scene.Groups) == 0 {
			for _, m := range scene.Movies {
				g := m.ScrapedGroup()
				scene.Groups = append(scene.Groups, &g)
			}
		} else if len(scene.Groups) > 0 && len(scene.Movies) == 0 {
			for _, g := range scene.Groups {
				m := g.ScrapedMovie()
				scene.Movies = append(scene.Movies, &m)
			}
		}

		tags, err := postProcessTags(ctx, tqb, scene.Tags)
		if err != nil {
			return err
		}
		scene.Tags, ignoredTags = FilterTags(excludeTagRE, tags)

		if scene.Studio != nil {
			err := match.ScrapedStudio(ctx, sqb, scene.Studio, "")
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, nil, err
	}

	// post-process - set the image if applicable
	if err := setSceneImage(ctx, c.client, &scene, c.globalConfig); err != nil {
		logger.Warnf("Could not set image using URL %s: %v", *scene.Image, err)
	}

	return scene, ignoredTags, nil
}

func (c Cache) postScrapeGallery(ctx context.Context, g models.ScrapedGallery, excludeTagRE []*regexp.Regexp) (_ ScrapedContent, ignoredTags []string, err error) {
	// set the URL/URLs field
	if g.URL == nil && len(g.URLs) > 0 {
		g.URL = &g.URLs[0]
	}
	if g.URL != nil && len(g.URLs) == 0 {
		g.URLs = []string{*g.URL}
	}

	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		pqb := r.PerformerFinder
		tqb := r.TagFinder
		sqb := r.StudioFinder

		for _, p := range g.Performers {
			err := match.ScrapedPerformer(ctx, pqb, p, "")
			if err != nil {
				return err
			}
		}

		tags, err := postProcessTags(ctx, tqb, g.Tags)
		if err != nil {
			return err
		}
		g.Tags, ignoredTags = FilterTags(excludeTagRE, tags)

		if g.Studio != nil {
			err := match.ScrapedStudio(ctx, sqb, g.Studio, "")
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, nil, err
	}

	return g, ignoredTags, nil
}

func (c Cache) postScrapeImage(ctx context.Context, image models.ScrapedImage, excludeTagRE []*regexp.Regexp) (_ ScrapedContent, ignoredTags []string, err error) {
	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		pqb := r.PerformerFinder
		tqb := r.TagFinder
		sqb := r.StudioFinder

		for _, p := range image.Performers {
			if err := match.ScrapedPerformer(ctx, pqb, p, ""); err != nil {
				return err
			}
		}

		tags, err := postProcessTags(ctx, tqb, image.Tags)
		if err != nil {
			return err
		}

		image.Tags, ignoredTags = FilterTags(excludeTagRE, tags)

		if image.Studio != nil {
			err := match.ScrapedStudio(ctx, sqb, image.Studio, "")
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, nil, err
	}

	return image, ignoredTags, nil
}
