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
	scrapeByURL(ctx context.Context, url string, ty ScrapeContentType) (ScrapedContent, error)
	scrapeByName(ctx context.Context, name string, ty ScrapeContentType) ([]ScrapedContent, error)
	scrapeByFragment(ctx context.Context, input Input) (ScrapedContent, error)

	scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*ScrapedScene, error)
	scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*ScrapedGallery, error)
}

func (c config) getScraper(scraper scraperTypeConfig, client *http.Client, globalConfig GlobalConfig) scraperActionImpl {
	switch scraper.Action {
	case scraperActionScript:
		return newScriptScraper(scraper, c, globalConfig)
	case scraperActionStash:
		return newStashScraper(scraper, client, c, globalConfig)
	case scraperActionXPath:
		return newXpathScraper(scraper, client, c, globalConfig)
	case scraperActionJson:
		return newJsonScraper(scraper, client, c, globalConfig)
	}

	panic("unknown scraper action: " + scraper.Action)
}
