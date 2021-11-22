package scraper

import (
	"context"
	"net/http"

	"github.com/stashapp/stash/pkg/models"
)

type scraperAction string

const (
	scraperActionScript scraperAction = "script"
	scraperActionStash  scraperAction = "stash"
	scraperActionXPath  scraperAction = "scrapeXPath"
	scraperActionJson   scraperAction = "scrapeJson"
)

func (e scraperAction) IsValid() bool {
	switch e {
	case scraperActionScript, scraperActionStash, scraperActionXPath, scraperActionJson:
		return true
	}
	return false
}

type scraperActionImpl interface {
	scrapePerformersByName(ctx context.Context, name string) ([]*models.ScrapedPerformer, error)
	scrapePerformerByFragment(scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error)
	scrapePerformerByURL(ctx context.Context, url string) (*models.ScrapedPerformer, error)

	scrapeScenesByName(ctx context.Context, name string) ([]*models.ScrapedScene, error)
	scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*models.ScrapedScene, error)
	scrapeSceneByFragment(ctx context.Context, scene models.ScrapedSceneInput) (*models.ScrapedScene, error)
	scrapeSceneByURL(ctx context.Context, url string) (*models.ScrapedScene, error)

	scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*models.ScrapedGallery, error)
	scrapeGalleryByFragment(gallery models.ScrapedGalleryInput) (*models.ScrapedGallery, error)
	scrapeGalleryByURL(ctx context.Context, url string) (*models.ScrapedGallery, error)

	scrapeMovieByURL(ctx context.Context, url string) (*models.ScrapedMovie, error)
}

func (c config) getScraper(scraper scraperTypeConfig, client *http.Client, txnManager models.TransactionManager, globalConfig GlobalConfig) scraperActionImpl {
	switch scraper.Action {
	case scraperActionScript:
		return newScriptScraper(scraper, c, globalConfig)
	case scraperActionStash:
		return newStashScraper(scraper, client, txnManager, c, globalConfig)
	case scraperActionXPath:
		return newXpathScraper(scraper, client, txnManager, c, globalConfig)
	case scraperActionJson:
		return newJsonScraper(scraper, client, txnManager, c, globalConfig)
	}

	panic("unknown scraper action: " + scraper.Action)
}
