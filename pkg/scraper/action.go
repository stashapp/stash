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

type urlScraperActionImpl interface {
	scrapeByURL(ctx context.Context, url string, ty ScrapeContentType) (ScrapedContent, error)
}

func (c Definition) getURLScraper(def ByURLDefinition, client *http.Client, globalConfig GlobalConfig) urlScraperActionImpl {
	switch def.Action {
	case scraperActionScript:
		return &scriptURLScraper{
			scriptScraper: scriptScraper{
				definition:   c,
				globalConfig: globalConfig,
			},
			definition: def,
		}
	case scraperActionStash:
		return newStashScraper(client, c, globalConfig)
	case scraperActionXPath:
		return &xpathURLScraper{
			xpathScraper: xpathScraper{
				definition:   c,
				globalConfig: globalConfig,
				client:       client,
			},
			definition: def,
		}
	case scraperActionJson:
		return &jsonURLScraper{
			jsonScraper: jsonScraper{
				definition:   c,
				globalConfig: globalConfig,
				client:       client,
			},
			definition: def,
		}
	}

	panic("unknown scraper action: " + def.Action)
}

type nameScraperActionImpl interface {
	scrapeByName(ctx context.Context, name string, ty ScrapeContentType) ([]ScrapedContent, error)
}

func (c Definition) getNameScraper(def ByNameDefinition, client *http.Client, globalConfig GlobalConfig) nameScraperActionImpl {
	switch def.Action {
	case scraperActionScript:
		return &scriptNameScraper{
			scriptScraper: scriptScraper{
				definition:   c,
				globalConfig: globalConfig,
			},
			definition: def,
		}
	case scraperActionStash:
		return newStashScraper(client, c, globalConfig)
	case scraperActionXPath:
		return &xpathNameScraper{
			xpathScraper: xpathScraper{
				definition:   c,
				globalConfig: globalConfig,
				client:       client,
			},
			definition: def,
		}
	case scraperActionJson:
		return &jsonNameScraper{
			jsonScraper: jsonScraper{
				definition:   c,
				globalConfig: globalConfig,
				client:       client,
			},
			definition: def,
		}
	}

	panic("unknown scraper action: " + def.Action)
}

type fragmentScraperActionImpl interface {
	scrapeByFragment(ctx context.Context, input Input) (ScrapedContent, error)

	scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*models.ScrapedScene, error)
	scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*models.ScrapedGallery, error)
	scrapeImageByImage(ctx context.Context, image *models.Image) (*models.ScrapedImage, error)
}

func (c Definition) getFragmentScraper(actionDef ByFragmentDefinition, client *http.Client, globalConfig GlobalConfig) fragmentScraperActionImpl {
	switch actionDef.Action {
	case scraperActionScript:
		return &scriptFragmentScraper{
			scriptScraper: scriptScraper{
				definition:   c,
				globalConfig: globalConfig,
			},
			definition: actionDef,
		}
	case scraperActionStash:
		return newStashScraper(client, c, globalConfig)
	case scraperActionXPath:
		return &xpathFragmentScraper{
			xpathScraper: xpathScraper{
				definition:   c,
				globalConfig: globalConfig,
				client:       client,
			},
			definition: actionDef,
		}
	case scraperActionJson:
		return &jsonFragmentScraper{
			jsonScraper: jsonScraper{
				definition:   c,
				globalConfig: globalConfig,
				client:       client,
			},
			definition: actionDef,
		}
	}

	panic("unknown scraper action: " + actionDef.Action)
}
