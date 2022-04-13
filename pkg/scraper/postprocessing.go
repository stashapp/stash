package scraper

import (
	"context"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
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
	case *models.ScrapedMovie:
		if v != nil {
			return c.postScrapeMovie(ctx, *v)
		}
	case models.ScrapedMovie:
		return c.postScrapeMovie(ctx, v)
	}

	// If nothing matches, pass the content through
	return content, nil
}

func (c Cache) postScrapePerformer(ctx context.Context, p models.ScrapedPerformer) (ScrapedContent, error) {
	r := c.txnManager
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		tqb := r.Tag

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

	return p, nil
}

func (c Cache) postScrapeMovie(ctx context.Context, m models.ScrapedMovie) (ScrapedContent, error) {
	r := c.txnManager
	if m.Studio != nil {
		if err := c.txnManager.WithTxn(ctx, func(ctx context.Context) error {
			return match.ScrapedStudio(ctx, r.Studio, m.Studio, nil)
		}); err != nil {
			return nil, err
		}
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

func (c Cache) postScrapeScenePerformer(ctx context.Context, p models.ScrapedPerformer) error {
	r := c.txnManager
	if err := c.txnManager.WithTxn(ctx, func(ctx context.Context) error {
		tqb := r.Tag

		tags, err := postProcessTags(ctx, tqb, p.Tags)
		if err != nil {
			return err
		}
		p.Tags = tags

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (c Cache) postScrapeScene(ctx context.Context, scene ScrapedScene) (ScrapedContent, error) {
	r := c.txnManager
	if err := c.txnManager.WithTxn(ctx, func(ctx context.Context) error {
		pqb := r.Performer
		mqb := r.Movie
		tqb := r.Tag
		sqb := r.Studio

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
			err := match.ScrapedMovie(ctx, mqb, p)
			if err != nil {
				return err
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
	r := c.txnManager
	if err := c.txnManager.WithTxn(ctx, func(ctx context.Context) error {
		pqb := r.Performer
		tqb := r.Tag
		sqb := r.Studio

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

func postProcessTags(ctx context.Context, tqb models.TagReader, scrapedTags []*models.ScrapedTag) ([]*models.ScrapedTag, error) {
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
