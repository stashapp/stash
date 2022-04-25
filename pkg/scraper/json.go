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

func (s *jsonScraper) getJsonScraper() *mappedScraper {
	return s.config.JsonScrapers[s.scraper.Scraper]
}

func (s *jsonScraper) scrapeURL(ctx context.Context, url string) (string, *mappedScraper, error) {
	scraper := s.getJsonScraper()

	if scraper == nil {
		return "", nil, errors.New("json scraper with name " + s.scraper.Scraper + " not found in config")
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return "", nil, err
	}

	return doc, scraper, nil
}

func (s *jsonScraper) loadURL(ctx context.Context, url string) (string, error) {
	r, err := loadURL(ctx, url, s.client, s.config, s.globalConfig)
	if err != nil {
		return "", err
	}
	logger.Infof("loadURL (%s)\n", url)
	doc, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	docStr := string(doc)
	if !gjson.Valid(docStr) {
		return "", errors.New("not valid json")
	}

	if err == nil && s.config.DebugOptions != nil && s.config.DebugOptions.PrintHTML {
		logger.Infof("loadURL (%s) response: \n%s", url, docStr)
	}

	return docStr, err
}

func (s *jsonScraper) scrapeByURL(ctx context.Context, url string, ty ScrapeContentType) (ScrapedContent, error) {
	u := replaceURL(url, s.scraper) // allow a URL Replace for url-queries
	doc, scraper, err := s.scrapeURL(ctx, u)
	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc)
	switch ty {
	case ScrapeContentTypePerformer:
		return scraper.scrapePerformer(ctx, q)
	case ScrapeContentTypeScene:
		return scraper.scrapeScene(ctx, q)
	case ScrapeContentTypeGallery:
		return scraper.scrapeGallery(ctx, q)
	case ScrapeContentTypeMovie:
		return scraper.scrapeMovie(ctx, q)
	}

	return nil, ErrNotSupported
}

func (s *jsonScraper) scrapeByName(ctx context.Context, name string, ty ScrapeContentType) ([]ScrapedContent, error) {
	scraper := s.getJsonScraper()

	if scraper == nil {
		return nil, fmt.Errorf("%w: name %v", ErrNotFound, s.scraper.Scraper)
	}

	const placeholder = "{}"

	// replace the placeholder string with the URL-escaped name
	escapedName := url.QueryEscape(name)

	url := s.scraper.QueryURL
	url = strings.ReplaceAll(url, placeholder, escapedName)

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc)
	q.setType(SearchQuery)

	var content []ScrapedContent
	switch ty {
	case ScrapeContentTypePerformer:
		performers, err := scraper.scrapePerformers(ctx, q)
		if err != nil {
			return nil, err
		}

		for _, p := range performers {
			content = append(content, p)
		}

		return content, nil
	case ScrapeContentTypeScene:
		scenes, err := scraper.scrapeScenes(ctx, q)
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

func (s *jsonScraper) scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*ScrapedScene, error) {
	// construct the URL
	queryURL := queryURLParametersFromScene(scene)
	if s.scraper.QueryURLReplacements != nil {
		queryURL.applyReplacements(s.scraper.QueryURLReplacements)
	}
	url := queryURL.constructURL(s.scraper.QueryURL)

	scraper := s.getJsonScraper()

	if scraper == nil {
		return nil, errors.New("json scraper with name " + s.scraper.Scraper + " not found in config")
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc)
	return scraper.scrapeScene(ctx, q)
}

func (s *jsonScraper) scrapeByFragment(ctx context.Context, input Input) (ScrapedContent, error) {
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

	scraper := s.getJsonScraper()

	if scraper == nil {
		return nil, errors.New("xpath scraper with name " + s.scraper.Scraper + " not found in config")
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc)
	return scraper.scrapeScene(ctx, q)
}

func (s *jsonScraper) scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*ScrapedGallery, error) {
	// construct the URL
	queryURL := queryURLParametersFromGallery(gallery)
	if s.scraper.QueryURLReplacements != nil {
		queryURL.applyReplacements(s.scraper.QueryURLReplacements)
	}
	url := queryURL.constructURL(s.scraper.QueryURL)

	scraper := s.getJsonScraper()

	if scraper == nil {
		return nil, errors.New("json scraper with name " + s.scraper.Scraper + " not found in config")
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc)
	return scraper.scrapeGallery(ctx, q)
}

func (s *jsonScraper) getJsonQuery(doc string) *jsonQuery {
	return &jsonQuery{
		doc:     doc,
		scraper: s,
	}
}

type jsonQuery struct {
	doc       string
	scraper   *jsonScraper
	queryType QueryType
}

func (q *jsonQuery) getType() QueryType {
	return q.queryType
}

func (q *jsonQuery) setType(t QueryType) {
	q.queryType = t
}

func (q *jsonQuery) runQuery(selector string) ([]string, error) {
	value := gjson.Get(q.doc, selector)

	if !value.Exists() {
		// many possible reasons why the selector may not be in the json object
		// and not all are errors.
		// Just return nil
		return nil, nil
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

func (q *jsonQuery) subScrape(ctx context.Context, value string) mappedQuery {
	doc, err := q.scraper.loadURL(ctx, value)

	if err != nil {
		logger.Warnf("Error getting URL '%s' for sub-scraper: %s", value, err.Error())
		return nil
	}

	return q.scraper.getJsonQuery(doc)
}
