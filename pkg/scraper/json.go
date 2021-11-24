package scraper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/tidwall/gjson"
)

type jsonScraper struct {
	scraper      scraperTypeConfig
	config       config
	globalConfig GlobalConfig
	client       *http.Client
	txnManager   models.TransactionManager
}

func newJsonScraper(scraper scraperTypeConfig, client *http.Client, txnManager models.TransactionManager, config config, globalConfig GlobalConfig) *jsonScraper {
	return &jsonScraper{
		scraper:      scraper,
		config:       config,
		client:       client,
		globalConfig: globalConfig,
		txnManager:   txnManager,
	}
}

func (s *jsonScraper) mappedScraper() (*mappedScraper, error) {
	scraper := s.config.JsonScrapers[s.scraper.Scraper]
	if scraper == nil {
		return nil, fmt.Errorf("%w: searched for scraper name %v", ErrNotFound, s.scraper.Scraper)
	}

	return scraper, nil

}

var ErrInvalidJSON = errors.New("invalid json")

func (s *jsonScraper) loadURL(ctx context.Context, url string) (docQueryer, error) {
	r, err := loadURL(ctx, url, s.client, s.config, s.globalConfig)
	if err != nil {
		return nil, err
	}
	logger.Infof("loadURL (%s)\n", url)
	doc, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	docStr := string(doc)
	if !gjson.Valid(docStr) {
		return nil, ErrInvalidJSON
	}

	if err == nil && s.config.DebugOptions != nil && s.config.DebugOptions.PrintHTML {
		logger.Infof("loadURL (%s) response: \n%s", url, docStr)
	}

	return s.newDocQueryer(docStr), err
}

func (s *jsonScraper) scrapeByURL(ctx context.Context, url string, ty models.ScrapeContentType) (models.ScrapedContent, error) {
	u := replaceURL(url, s.scraper) // allow a URL Replace for url-queries
	scraper, err := s.mappedScraper()
	if err != nil {
		return nil, err
	}

	dq, err := s.loadURL(ctx, u)
	if err != nil {
		return nil, err
	}

	switch ty {
	case models.ScrapeContentTypePerformer:
		return scraper.scrapePerformer(ctx, dq)
	case models.ScrapeContentTypeScene:
		return scraper.scrapeScene(ctx, dq)
	case models.ScrapeContentTypeGallery:
		return scraper.scrapeGallery(ctx, dq)
	case models.ScrapeContentTypeMovie:
		return scraper.scrapeMovie(ctx, dq)
	}

	return nil, ErrNotSupported
}

func (s *jsonScraper) scrapeByName(ctx context.Context, name string, ty models.ScrapeContentType) ([]models.ScrapedContent, error) {
	scraper, err := s.mappedScraper()
	if err != nil {
		return nil, err
	}

	const placeholder = "{}"

	// replace the placeholder string with the URL-escaped name
	escapedName := url.QueryEscape(name)

	url := s.scraper.QueryURL
	url = strings.ReplaceAll(url, placeholder, escapedName)

	dq, err := s.loadURL(ctx, url)
	if err != nil {
		return nil, err
	}

	var content []models.ScrapedContent
	switch ty {
	case models.ScrapeContentTypePerformer:
		performers, err := scraper.scrapePerformers(ctx, dq)
		if err != nil {
			return nil, err
		}

		for _, p := range performers {
			content = append(content, p)
		}

		return content, nil
	case models.ScrapeContentTypeScene:
		scenes, err := scraper.scrapeScenes(ctx, dq)
		if err != nil {
			return nil, err
		}

		for _, s := range scenes {
			content = append(content, s)
		}

		return content, nil
	}

	return nil, ErrNotSupported
}

func (s *jsonScraper) scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*models.ScrapedScene, error) {
	// construct the URL
	queryURL := queryURLParametersFromScene(scene)
	if s.scraper.QueryURLReplacements != nil {
		queryURL.applyReplacements(s.scraper.QueryURLReplacements)
	}
	url := queryURL.constructURL(s.scraper.QueryURL)

	scraper, err := s.mappedScraper()
	if err != nil {
		return nil, err
	}

	dq, err := s.loadURL(ctx, url)
	if err != nil {
		return nil, err
	}

	return scraper.scrapeScene(ctx, dq)
}

func (s *jsonScraper) scrapeByFragment(ctx context.Context, input Input) (models.ScrapedContent, error) {
	switch {
	case input.Gallery != nil:
		return nil, fmt.Errorf("%w: cannot use a json scraper as a gallery fragment scraper", ErrNotSupported)
	case input.Performer != nil:
		return nil, fmt.Errorf("%w: cannot use a json scraper as a performer fragment scraper", ErrNotSupported)
	case input.Scene == nil:
		return nil, fmt.Errorf("%w: scene input is nil", ErrNotSupported)
	}

	scene := *input.Scene

	// construct the URL
	queryURL := queryURLParametersFromScrapedScene(scene)
	if s.scraper.QueryURLReplacements != nil {
		queryURL.applyReplacements(s.scraper.QueryURLReplacements)
	}
	url := queryURL.constructURL(s.scraper.QueryURL)

	scraper, err := s.mappedScraper()
	if err != nil {
		return nil, err
	}

	dq, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	return scraper.scrapeScene(ctx, dq)
}

func (s *jsonScraper) scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*models.ScrapedGallery, error) {
	// construct the URL
	queryURL := queryURLParametersFromGallery(gallery)
	if s.scraper.QueryURLReplacements != nil {
		queryURL.applyReplacements(s.scraper.QueryURLReplacements)
	}
	url := queryURL.constructURL(s.scraper.QueryURL)

	scraper, err := s.mappedScraper()
	if err != nil {
		return nil, err
	}

	dq, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	return scraper.scrapeGallery(ctx, dq)
}

func (s *jsonScraper) subScrape(ctx context.Context, value string) docQueryer {
	doc, err := s.loadURL(ctx, value)

	if err != nil {
		logger.Warnf("Error getting URL '%s' for sub-scraper: %s", value, err.Error())
		return nil
	}

	return doc
}

// newDocQueryer turns a json scraper into a queryer over doc
func (s *jsonScraper) newDocQueryer(doc string) docQueryer {
	return &jsonQuery{
		doc:     doc,
		scraper: s,
	}
}

type jsonQuery struct {
	doc     string
	scraper *jsonScraper
}

func (q *jsonQuery) docQuery(selector string) ([]string, error) {
	value := gjson.Get(q.doc, selector)

	if !value.Exists() {
		return nil, fmt.Errorf("could not find json path '%s' in json object", selector)
	}

	var ret []string
	if value.IsArray() {
		value.ForEach(func(k, v gjson.Result) bool {
			ret = append(ret, v.String())
			return true
		})
	} else {
		ret = append(ret, value.String())
	}

	return ret, nil
}

func (q *jsonQuery) subScrape(ctx context.Context, value string) docQueryer {
	return q.scraper.subScrape(ctx, value)
}
