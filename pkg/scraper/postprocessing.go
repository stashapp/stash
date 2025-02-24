package scraper

import (
	"context"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

// postScrape handles post-processing of scraped content. If the content
// requires post-processing, this function fans out to the given content
// type and post-processes it.
func (c Cache) postScrape(ctx context.Context, content ScrapedContent) (ScrapedContent, error) {
	// Analyze the concrete type, call the right post-processing function
	switch v := content.(type) {
	case *models.ScrapedPerformer:
		if v != nil {
			return c.postScrapePerformer(ctx, *v)
		}
	case models.ScrapedPerformer:
		return c.postScrapePerformer(ctx, v)
	case *ScrapedScene:
		if v != nil {
			return c.postScrapeScene(ctx, *v)
		}
	case ScrapedScene:
		return c.postScrapeScene(ctx, v)
	case *ScrapedGallery:
		if v != nil {
			return c.postScrapeGallery(ctx, *v)
		}
	case ScrapedGallery:
		return c.postScrapeGallery(ctx, v)
	case *ScrapedImage:
		if v != nil {
			return c.postScrapeImage(ctx, *v)
		}
	case ScrapedImage:
		return c.postScrapeImage(ctx, v)
	case *models.ScrapedMovie:
		if v != nil {
			return c.postScrapeMovie(ctx, *v)
		}
	case models.ScrapedMovie:
		return c.postScrapeMovie(ctx, v)
	case *models.ScrapedGroup:
		if v != nil {
			return c.postScrapeGroup(ctx, *v)
		}
	case models.ScrapedGroup:
		return c.postScrapeGroup(ctx, v)
	}

	// If nothing matches, pass the content through
	return content, nil
}

func (c Cache) postScrapePerformer(ctx context.Context, p models.ScrapedPerformer) (ScrapedContent, error) {
	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		tqb := r.TagFinder

		tags, err := postProcessTags(ctx, tqb, p.Tags)
		if err != nil {
			return err
		}
		p.Tags = tags

		return nil
	}); err != nil {
		return nil, err
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

	return p, nil
}

func (c Cache) postScrapeMovie(ctx context.Context, m models.ScrapedMovie) (ScrapedContent, error) {
	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		tqb := r.TagFinder
		tags, err := postProcessTags(ctx, tqb, m.Tags)
		if err != nil {
			return err
		}
		m.Tags = tags

		if m.Studio != nil {
			if err := match.ScrapedStudio(ctx, r.StudioFinder, m.Studio, nil); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// post-process - set the image if applicable
	if err := setMovieFrontImage(ctx, c.client, &m, c.globalConfig); err != nil {
		logger.Warnf("could not set front image using URL %s: %v", *m.FrontImage, err)
	}
	if err := setMovieBackImage(ctx, c.client, &m, c.globalConfig); err != nil {
		logger.Warnf("could not set back image using URL %s: %v", *m.BackImage, err)
	}

	return m, nil
}

func (c Cache) postScrapeGroup(ctx context.Context, m models.ScrapedGroup) (ScrapedContent, error) {
	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		tqb := r.TagFinder
		tags, err := postProcessTags(ctx, tqb, m.Tags)
		if err != nil {
			return err
		}
		m.Tags = tags

		if m.Studio != nil {
			if err := match.ScrapedStudio(ctx, r.StudioFinder, m.Studio, nil); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// post-process - set the image if applicable
	if err := setGroupFrontImage(ctx, c.client, &m, c.globalConfig); err != nil {
		logger.Warnf("could not set front image using URL %s: %v", *m.FrontImage, err)
	}
	if err := setGroupBackImage(ctx, c.client, &m, c.globalConfig); err != nil {
		logger.Warnf("could not set back image using URL %s: %v", *m.BackImage, err)
	}

	return m, nil
}

func (c Cache) postScrapeScenePerformer(ctx context.Context, p models.ScrapedPerformer) error {
	tqb := c.repository.TagFinder

	tags, err := postProcessTags(ctx, tqb, p.Tags)
	if err != nil {
		return err
	}
	p.Tags = tags

	p.Country = resolveCountryName(p.Country)

	return nil
}

func (c Cache) postScrapeScene(ctx context.Context, scene ScrapedScene) (ScrapedContent, error) {
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

			if err := c.postScrapeScenePerformer(ctx, *p); err != nil {
				return err
			}

			if err := match.ScrapedPerformer(ctx, pqb, p, nil); err != nil {
				return err
			}
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
		scene.Tags = tags

		if scene.Studio != nil {
			err := match.ScrapedStudio(ctx, sqb, scene.Studio, nil)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// post-process - set the image if applicable
	if err := setSceneImage(ctx, c.client, &scene, c.globalConfig); err != nil {
		logger.Warnf("Could not set image using URL %s: %v", *scene.Image, err)
	}

	return scene, nil
}

func (c Cache) postScrapeGallery(ctx context.Context, g ScrapedGallery) (ScrapedContent, error) {
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
			err := match.ScrapedPerformer(ctx, pqb, p, nil)
			if err != nil {
				return err
			}
		}

		tags, err := postProcessTags(ctx, tqb, g.Tags)
		if err != nil {
			return err
		}
		g.Tags = tags

		if g.Studio != nil {
			err := match.ScrapedStudio(ctx, sqb, g.Studio, nil)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return g, nil
}

func (c Cache) postScrapeImage(ctx context.Context, image ScrapedImage) (ScrapedContent, error) {
	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		pqb := r.PerformerFinder
		tqb := r.TagFinder
		sqb := r.StudioFinder

		for _, p := range image.Performers {
			if err := match.ScrapedPerformer(ctx, pqb, p, nil); err != nil {
				return err
			}
		}

		tags, err := postProcessTags(ctx, tqb, image.Tags)
		if err != nil {
			return err
		}
		image.Tags = tags

		if image.Studio != nil {
			err := match.ScrapedStudio(ctx, sqb, image.Studio, nil)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return image, nil
}

func postProcessTags(ctx context.Context, tqb models.TagQueryer, scrapedTags []*models.ScrapedTag) ([]*models.ScrapedTag, error) {
	var ret []*models.ScrapedTag

	for _, t := range scrapedTags {
		err := match.ScrapedTag(ctx, tqb, t)
		if err != nil {
			return nil, err
		}
		ret = append(ret, t)
	}

	return ret, nil
}
