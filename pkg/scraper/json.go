package scraper

import (
	"context"
	"errors"
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

func (s *jsonScraper) scrapePerformerByURL(ctx context.Context, url string) (*models.ScrapedPerformer, error) {
	u := replaceURL(url, s.scraper) // allow a URL Replace for performer by URL queries
	doc, scraper, err := s.scrapeURL(ctx, u)
	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc)
	return scraper.scrapePerformer(q)
}

func (s *jsonScraper) scrapeSceneByURL(ctx context.Context, url string) (*models.ScrapedScene, error) {
	u := replaceURL(url, s.scraper) // allow a URL Replace for scene by URL queries
	doc, scraper, err := s.scrapeURL(ctx, u)
	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc)
	return scraper.scrapeScene(q)
}

func (s *jsonScraper) scrapeGalleryByURL(ctx context.Context, url string) (*models.ScrapedGallery, error) {
	u := replaceURL(url, s.scraper) // allow a URL Replace for gallery by URL queries
	doc, scraper, err := s.scrapeURL(ctx, u)
	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc)
	return scraper.scrapeGallery(q)
}

func (s *jsonScraper) scrapeMovieByURL(ctx context.Context, url string) (*models.ScrapedMovie, error) {
	u := replaceURL(url, s.scraper) // allow a URL Replace for movie by URL queries
	doc, scraper, err := s.scrapeURL(ctx, u)
	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc)
	return scraper.scrapeMovie(q)
}

func (s *jsonScraper) scrapePerformersByName(ctx context.Context, name string) ([]*models.ScrapedPerformer, error) {
	scraper := s.getJsonScraper()

	if scraper == nil {
		return nil, errors.New("json scraper with name " + s.scraper.Scraper + " not found in config")
	}

	const placeholder = "{}"

	// replace the placeholder string with the URL-escaped name
	escapedName := url.QueryEscape(name)

	url := s.scraper.QueryURL
	url = strings.ReplaceAll(url, placeholder, escapedName)

	doc, err := s.loadURL(context.TODO(), url)

	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc)
	return scraper.scrapePerformers(q)
}

func (s *jsonScraper) scrapePerformerByFragment(scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	return nil, errors.New("scrapePerformerByFragment not supported for json scraper")
}

func (s *jsonScraper) scrapeScenesByName(ctx context.Context, name string) ([]*models.ScrapedScene, error) {
	scraper := s.getJsonScraper()

	if scraper == nil {
		return nil, errors.New("json scraper with name " + s.scraper.Scraper + " not found in config")
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
	return scraper.scrapeScenes(q)
}

func (s *jsonScraper) scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*models.ScrapedScene, error) {
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
	return scraper.scrapeScene(q)
}

func (s *jsonScraper) scrapeSceneByFragment(ctx context.Context, scene models.ScrapedSceneInput) (*models.ScrapedScene, error) {
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
	return scraper.scrapeScene(q)
}

func (s *jsonScraper) scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*models.ScrapedGallery, error) {
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
	return scraper.scrapeGallery(q)
}

func (s *jsonScraper) scrapeGalleryByFragment(gallery models.ScrapedGalleryInput) (*models.ScrapedGallery, error) {
	return nil, errors.New("scrapeGalleryByFragment not supported for json scraper")
}

func (s *jsonScraper) getJsonQuery(doc string) *jsonQuery {
	return &jsonQuery{
		doc:     doc,
		scraper: s,
	}
}

type jsonQuery struct {
	doc     string
	scraper *jsonScraper
}

func (q *jsonQuery) runQuery(selector string) []string {
	value := gjson.Get(q.doc, selector)

	if !value.Exists() {
		logger.Warnf("Could not find json path '%s' in json object", selector)
		return nil
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

	return ret
}

func (q *jsonQuery) subScrape(value string) mappedQuery {
	doc, err := q.scraper.loadURL(context.TODO(), value)

	if err != nil {
		logger.Warnf("Error getting URL '%s' for sub-scraper: %s", value, err.Error())
		return nil
	}

	return q.scraper.getJsonQuery(doc)
}
