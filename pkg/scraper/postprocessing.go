package scraper

import (
	"context"
	"regexp"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	stash_config "github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
)

// postScrape handles post-processing of scraped content. If the content
// requires post-processing, this function fans out to the given content
// type and post-processes it.
func (c Cache) postScrape(ctx context.Context, content models.ScrapedContent) (models.ScrapedContent, error) {
	// Analyze the concrete type, call the right post-processing function
	switch v := content.(type) {
	case models.ScrapedPerformer:
		return c.postScrapePerformer(ctx, &v)
	case models.ScrapedScene:
		return c.postScrapeScene(ctx, &v)
	case models.ScrapedGallery:
		return c.postScrapeGallery(ctx, &v)
	case models.ScrapedMovie:
		return c.postScrapeMovie(ctx, &v)
	}

	// If nothing matches, pass the content through
	return content, nil
}

func (c Cache) postScrapePerformer(ctx context.Context, ret *models.ScrapedPerformer) (models.ScrapedContent, error) {
	if err := c.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		tqb := r.Tag()

		tags, err := postProcessTags(tqb, ret.Tags)
		if err != nil {
			return err
		}
		ret.Tags = tags

		return nil
	}); err != nil {
		return nil, err
	}

	// post-process - set the image if applicable
	if err := setPerformerImage(ctx, c.client, ret, c.globalConfig); err != nil {
		logger.Warnf("Could not set image using URL %s: %s", *ret.Image, err.Error())
	}

	return ret, nil
}

func (c Cache) postScrapeMovie(ctx context.Context, ret *models.ScrapedMovie) (models.ScrapedContent, error) {
	if ret.Studio != nil {
		if err := c.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
			return match.ScrapedStudio(r.Studio(), ret.Studio)
		}); err != nil {
			return nil, err
		}
	}

	// post-process - set the image if applicable
	if err := setMovieFrontImage(ctx, c.client, ret, c.globalConfig); err != nil {
		logger.Warnf("could not set front image using URL %s: %v", *ret.FrontImage, err)
	}
	if err := setMovieBackImage(ctx, c.client, ret, c.globalConfig); err != nil {
		logger.Warnf("could not set back image using URL %s: %v", *ret.BackImage, err)
	}

	return ret, nil
}

func (c Cache) postScrapeScenePerformer(ret *models.ScrapedPerformer) error {
	if err := c.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		tqb := r.Tag()

		tags, err := postProcessTags(tqb, ret.Tags)
		if err != nil {
			return err
		}
		ret.Tags = tags

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (c Cache) postScrapeScene(ctx context.Context, ret *models.ScrapedScene) (models.ScrapedContent, error) {
	if err := c.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		pqb := r.Performer()
		mqb := r.Movie()
		tqb := r.Tag()
		sqb := r.Studio()

		for _, p := range ret.Performers {
			if err := c.postScrapeScenePerformer(p); err != nil {
				return err
			}

			if err := match.ScrapedPerformer(pqb, p); err != nil {
				return err
			}
		}

		for _, p := range ret.Movies {
			err := match.ScrapedMovie(mqb, p)
			if err != nil {
				return err
			}
		}

		tags, err := postProcessTags(tqb, ret.Tags)
		if err != nil {
			return err
		}
		ret.Tags = tags

		if ret.Studio != nil {
			err := match.ScrapedStudio(sqb, ret.Studio)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// post-process - set the image if applicable
	if err := setSceneImage(ctx, c.client, ret, c.globalConfig); err != nil {
		logger.Warnf("Could not set image using URL %s: %v", *ret.Image, err)
	}

	return ret, nil
}

func (c Cache) postScrapeGallery(ctx context.Context, ret *models.ScrapedGallery) (models.ScrapedContent, error) {
	if err := c.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		pqb := r.Performer()
		tqb := r.Tag()
		sqb := r.Studio()

		for _, p := range ret.Performers {
			err := match.ScrapedPerformer(pqb, p)
			if err != nil {
				return err
			}
		}

		tags, err := postProcessTags(tqb, ret.Tags)
		if err != nil {
			return err
		}
		ret.Tags = tags

		if ret.Studio != nil {
			err := match.ScrapedStudio(sqb, ret.Studio)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func postProcessTags(tqb models.TagReader, scrapedTags []*models.ScrapedTag) ([]*models.ScrapedTag, error) {
	var ret []*models.ScrapedTag

	excludePatterns := stash_config.GetInstance().GetScraperExcludeTagPatterns()
	var excludeRegexps []*regexp.Regexp

	for _, excludePattern := range excludePatterns {
		reg, err := regexp.Compile(strings.ToLower(excludePattern))
		if err != nil {
			logger.Errorf("Invalid tag exclusion pattern :%v", err)
		} else {
			excludeRegexps = append(excludeRegexps, reg)
		}
	}

	var ignoredTags []string
ScrapeTag:
	for _, t := range scrapedTags {
		for _, reg := range excludeRegexps {
			if reg.MatchString(strings.ToLower(t.Name)) {
				ignoredTags = append(ignoredTags, t.Name)
				continue ScrapeTag
			}
		}

		err := match.ScrapedTag(tqb, t)
		if err != nil {
			return nil, err
		}
		ret = append(ret, t)
	}

	if len(ignoredTags) > 0 {
		logger.Infof("Scraping ignored tags: %s", strings.Join(ignoredTags, ", "))
	}

	return ret, nil
}
