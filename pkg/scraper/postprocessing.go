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

type postScraper struct {
	Cache
	excludeTagRE []*regexp.Regexp

	// ignoredTags is a list of tags that were ignored during post-processing
	ignoredTags []string
}

// postScrape handles post-processing of scraped content. If the content
// requires post-processing, this function fans out to the given content
// type and post-processes it.
// Assumes called within a read transaction.
func (c *postScraper) postScrape(ctx context.Context, content ScrapedContent) (_ ScrapedContent, err error) {
	const related = false

	// Analyze the concrete type, call the right post-processing function
	switch v := content.(type) {
	case *models.ScrapedPerformer:
		if v != nil {
			return c.postScrapePerformer(ctx, *v, related)
		}
	case models.ScrapedPerformer:
		return c.postScrapePerformer(ctx, v, related)
	case *models.ScrapedScene:
		if v != nil {
			return c.postScrapeScene(ctx, *v)
		}
	case models.ScrapedScene:
		return c.postScrapeScene(ctx, v)
	case *models.ScrapedGallery:
		if v != nil {
			return c.postScrapeGallery(ctx, *v)
		}
	case models.ScrapedGallery:
		return c.postScrapeGallery(ctx, v)
	case *models.ScrapedImage:
		if v != nil {
			return c.postScrapeImage(ctx, *v)
		}
	case models.ScrapedImage:
		return c.postScrapeImage(ctx, v)
	case *models.ScrapedMovie:
		if v != nil {
			return c.postScrapeMovie(ctx, *v, related)
		}
	case models.ScrapedMovie:
		return c.postScrapeMovie(ctx, v, related)
	case *models.ScrapedGroup:
		if v != nil {
			return c.postScrapeGroup(ctx, *v, related)
		}
	case models.ScrapedGroup:
		return c.postScrapeGroup(ctx, v, related)
	}

	// If nothing matches, pass the content through
	return content, nil
}

func (c *postScraper) filterTags(tags []*models.ScrapedTag) []*models.ScrapedTag {
	var ret []*models.ScrapedTag
	var thisIgnoredTags []string
	ret, thisIgnoredTags = FilterTags(c.excludeTagRE, tags)
	c.ignoredTags = sliceutil.AppendUniques(c.ignoredTags, thisIgnoredTags)

	return ret
}

func (c *postScraper) postScrapePerformer(ctx context.Context, p models.ScrapedPerformer, related bool) (_ ScrapedContent, err error) {
	r := c.repository
	tqb := r.TagFinder

	tags, err := postProcessTags(ctx, tqb, p.Tags)
	if err != nil {
		return nil, err
	}

	p.Tags = c.filterTags(tags)

	// post-process - set the image if applicable
	// don't set image for related performers to avoid excessive network calls
	if !related {
		if err := setPerformerImage(ctx, c.client, &p, c.globalConfig); err != nil {
			logger.Warnf("Could not set image using URL %s: %s", *p.Image, err.Error())
		}
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

	return p, nil
}

func (c *postScraper) postScrapeMovie(ctx context.Context, m models.ScrapedMovie, related bool) (_ ScrapedContent, err error) {
	r := c.repository
	tqb := r.TagFinder
	tags, err := postProcessTags(ctx, tqb, m.Tags)
	if err != nil {
		return nil, err
	}
	m.Tags = c.filterTags(tags)

	if m.Studio != nil {
		if err := match.ScrapedStudio(ctx, r.StudioFinder, m.Studio, ""); err != nil {
			return nil, err
		}
	}

	// populate URL/URLs
	// if URLs are provided, only use those
	if len(m.URLs) > 0 {
		m.URL = &m.URLs[0]
	} else {
		urls := []string{}
		if m.URL != nil {
			urls = append(urls, *m.URL)
		}

		if len(urls) > 0 {
			m.URLs = urls
		}
	}

	// post-process - set the image if applicable
	// don't set images for related movies to avoid excessive network calls
	if !related {
		if err := processImageField(ctx, m.FrontImage, c.client, c.globalConfig); err != nil {
			logger.Warnf("could not set front image using URL %s: %v", *m.FrontImage, err)
		}
		if err := processImageField(ctx, m.BackImage, c.client, c.globalConfig); err != nil {
			logger.Warnf("could not set back image using URL %s: %v", *m.BackImage, err)
		}
	}

	return m, nil
}

func (c *postScraper) postScrapeGroup(ctx context.Context, m models.ScrapedGroup, related bool) (_ ScrapedContent, err error) {
	r := c.repository
	tqb := r.TagFinder
	tags, err := postProcessTags(ctx, tqb, m.Tags)
	if err != nil {
		return nil, err
	}
	m.Tags = c.filterTags(tags)

	if m.Studio != nil {
		if err := match.ScrapedStudio(ctx, r.StudioFinder, m.Studio, ""); err != nil {
			return nil, err
		}
	}

	// populate URL/URLs
	// if URLs are provided, only use those
	if len(m.URLs) > 0 {
		m.URL = &m.URLs[0]
	} else {
		urls := []string{}
		if m.URL != nil {
			urls = append(urls, *m.URL)
		}

		if len(urls) > 0 {
			m.URLs = urls
		}
	}

	// post-process - set the image if applicable
	// don't set images for related groups to avoid excessive network calls
	if !related {
		if err := processImageField(ctx, m.FrontImage, c.client, c.globalConfig); err != nil {
			logger.Warnf("could not set front image using URL %s: %v", *m.FrontImage, err)
		}
		if err := processImageField(ctx, m.BackImage, c.client, c.globalConfig); err != nil {
			logger.Warnf("could not set back image using URL %s: %v", *m.BackImage, err)
		}
	}

	return m, nil
}

// postScrapeRelatedPerformers post-processes a list of performers.
// It modifies the performers in place.
func (c *postScraper) postScrapeRelatedPerformers(ctx context.Context, items []*models.ScrapedPerformer) error {
	for _, p := range items {
		if p == nil {
			continue
		}

		const related = true
		sc, err := c.postScrapePerformer(ctx, *p, related)
		if err != nil {
			return err
		}
		newP := sc.(models.ScrapedPerformer)
		*p = newP

		if err := match.ScrapedPerformer(ctx, c.repository.PerformerFinder, p, ""); err != nil {
			return err
		}
	}
	return nil
}

func (c *postScraper) postScrapeRelatedMovies(ctx context.Context, items []*models.ScrapedMovie) error {
	for _, p := range items {
		const related = true
		sc, err := c.postScrapeMovie(ctx, *p, related)
		if err != nil {
			return err
		}
		newP := sc.(models.ScrapedMovie)
		*p = newP

		matchedID, err := match.ScrapedGroup(ctx, c.repository.GroupFinder, p.StoredID, p.Name)
		if err != nil {
			return err
		}

		if matchedID != nil {
			p.StoredID = matchedID
		}
	}

	return nil
}

func (c *postScraper) postScrapeRelatedGroups(ctx context.Context, items []*models.ScrapedGroup) error {
	for _, p := range items {
		const related = true
		sc, err := c.postScrapeGroup(ctx, *p, related)
		if err != nil {
			return err
		}
		newP := sc.(models.ScrapedGroup)
		*p = newP

		matchedID, err := match.ScrapedGroup(ctx, c.repository.GroupFinder, p.StoredID, p.Name)
		if err != nil {
			return err
		}

		if matchedID != nil {
			p.StoredID = matchedID
		}
	}

	return nil
}

func (c *postScraper) postScrapeStudio(ctx context.Context, s models.ScrapedStudio, related bool) (_ ScrapedContent, err error) {
	r := c.repository
	tqb := r.TagFinder

	tags, err := postProcessTags(ctx, tqb, s.Tags)
	if err != nil {
		return nil, err
	}

	s.Tags = c.filterTags(tags)

	// post-process - set the image if applicable
	// don't set image for related studios to avoid excessive network calls
	if !related {
		if err := setStudioImage(ctx, c.client, &s, c.globalConfig); err != nil {
			logger.Warnf("Could not set image using URL %s: %s", *s.Image, err.Error())
		}
	}

	// populate URL/URLs
	// if URLs are provided, only use those
	if len(s.URLs) > 0 {
		s.URL = &s.URLs[0]
	} else {
		urls := []string{}
		if s.URL != nil {
			urls = append(urls, *s.URL)
		}

		if len(urls) > 0 {
			s.URLs = urls
		}
	}

	return s, nil
}

func (c *postScraper) postScrapeRelatedStudio(ctx context.Context, s *models.ScrapedStudio) error {
	if s == nil {
		return nil
	}

	const related = true
	sc, err := c.postScrapeStudio(ctx, *s, related)
	if err != nil {
		return err
	}
	newS := sc.(models.ScrapedStudio)
	*s = newS

	if err = match.ScrapedStudio(ctx, c.repository.StudioFinder, s, ""); err != nil {
		return err
	}

	return nil
}

func (c *postScraper) postScrapeScene(ctx context.Context, scene models.ScrapedScene) (_ ScrapedContent, err error) {
	// set the URL/URLs field
	if scene.URL == nil && len(scene.URLs) > 0 {
		scene.URL = &scene.URLs[0]
	}
	if scene.URL != nil && len(scene.URLs) == 0 {
		scene.URLs = []string{*scene.URL}
	}

	r := c.repository
	tqb := r.TagFinder

	if err = c.postScrapeRelatedPerformers(ctx, scene.Performers); err != nil {
		return nil, err
	}

	if err = c.postScrapeRelatedMovies(ctx, scene.Movies); err != nil {
		return nil, err
	}

	if err = c.postScrapeRelatedGroups(ctx, scene.Groups); err != nil {
		return nil, err
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
		return nil, err
	}
	scene.Tags = c.filterTags(tags)

	if err := c.postScrapeRelatedStudio(ctx, scene.Studio); err != nil {
		return nil, err
	}

	// post-process - set the image if applicable
	if err := processImageField(ctx, scene.Image, c.client, c.globalConfig); err != nil {
		logger.Warnf("Could not set image using URL %s: %v", *scene.Image, err)
	}

	return scene, nil
}

func (c *postScraper) postScrapeGallery(ctx context.Context, g models.ScrapedGallery) (_ ScrapedContent, err error) {
	// set the URL/URLs field
	if g.URL == nil && len(g.URLs) > 0 {
		g.URL = &g.URLs[0]
	}
	if g.URL != nil && len(g.URLs) == 0 {
		g.URLs = []string{*g.URL}
	}

	r := c.repository
	tqb := r.TagFinder

	if err = c.postScrapeRelatedPerformers(ctx, g.Performers); err != nil {
		return nil, err
	}

	tags, err := postProcessTags(ctx, tqb, g.Tags)
	if err != nil {
		return nil, err
	}
	g.Tags = c.filterTags(tags)

	if err := c.postScrapeRelatedStudio(ctx, g.Studio); err != nil {
		return nil, err
	}

	return g, nil
}

func (c *postScraper) postScrapeImage(ctx context.Context, image models.ScrapedImage) (_ ScrapedContent, err error) {
	r := c.repository
	tqb := r.TagFinder

	if err = c.postScrapeRelatedPerformers(ctx, image.Performers); err != nil {
		return nil, err
	}

	tags, err := postProcessTags(ctx, tqb, image.Tags)
	if err != nil {
		return nil, err
	}

	image.Tags = c.filterTags(tags)

	if err := c.postScrapeRelatedStudio(ctx, image.Studio); err != nil {
		return nil, err
	}

	return image, nil
}

// postScrapeSingle handles post-processing of a single scraped content item.
// This is a convenience function that includes logging the ignored tags, as opposed to logging them in the caller.
func (c Cache) postScrapeSingle(ctx context.Context, content ScrapedContent) (ret ScrapedContent, err error) {
	pp := postScraper{
		Cache:        c,
		excludeTagRE: c.compileExcludeTagPatterns(),
	}

	if err := c.repository.WithReadTxn(ctx, func(ctx context.Context) error {
		ret, err = pp.postScrape(ctx, content)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	LogIgnoredTags(pp.ignoredTags)
	return ret, nil
}
