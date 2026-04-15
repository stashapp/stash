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
	definition   Definition
	globalConfig GlobalConfig
	client       *http.Client
}

func (s *jsonScraper) getJsonScraper(name string) (*mappedScraper, error) {
	ret, ok := s.definition.JsonScrapers[name]
	if !ok {
		return nil, fmt.Errorf("json scraper with name %s not found in config", name)
	}

	return &ret, nil
}

func (s *jsonScraper) loadURL(ctx context.Context, url string) (string, error) {
	r, err := loadURL(ctx, url, s.client, s.definition, s.globalConfig)
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

	if s.definition.DebugOptions != nil && s.definition.DebugOptions.PrintHTML {
		logger.Infof("loadURL (%s) response: \n%s", url, docStr)
	}

	return docStr, err
}

type jsonURLScraper struct {
	jsonScraper
	definition ByURLDefinition
}

func (s *jsonURLScraper) scrapeByURL(ctx context.Context, url string, ty ScrapeContentType) (ScrapedContent, error) {
	scraper, err := s.getJsonScraper(s.definition.Scraper)
	if err != nil {
		return nil, err
	}

	doc, err := s.loadURL(ctx, url)
	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc, url)
	// if these just return the return values from scraper.scrape* functions then
	// it ends up returning ScrapedContent(nil) rather than nil
	switch ty {
	case ScrapeContentTypePerformer:
		ret, err := scraper.scrapePerformer(ctx, q)
		if err != nil || ret == nil {
			return nil, err
		}
		return ret, nil
	case ScrapeContentTypeScene:
		ret, err := scraper.scrapeScene(ctx, q)
		if err != nil || ret == nil {
			return nil, err
		}
		return ret, nil
	case ScrapeContentTypeGallery:
		ret, err := scraper.scrapeGallery(ctx, q)
		if err != nil || ret == nil {
			return nil, err
		}
		return ret, nil
	case ScrapeContentTypeImage:
		ret, err := scraper.scrapeImage(ctx, q)
		if err != nil || ret == nil {
			return nil, err
		}
		return ret, nil
	case ScrapeContentTypeMovie, ScrapeContentTypeGroup:
		ret, err := scraper.scrapeGroup(ctx, q)
		if err != nil || ret == nil {
			return nil, err
		}
		return ret, nil
	}

	return nil, ErrNotSupported
}

type jsonNameScraper struct {
	jsonScraper
	definition ByNameDefinition
}

func (s *jsonNameScraper) scrapeByName(ctx context.Context, name string, ty ScrapeContentType) ([]ScrapedContent, error) {
	scraper, err := s.getJsonScraper(s.definition.Scraper)
	if err != nil {
		return nil, err
	}

	const placeholder = "{}"

	// replace the placeholder string with the URL-escaped name
	escapedName := url.QueryEscape(name)

	url := s.definition.QueryURL
	url = strings.ReplaceAll(url, placeholder, escapedName)

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc, url)
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

type jsonFragmentScraper struct {
	jsonScraper
	definition ByFragmentDefinition
}

func (s *jsonFragmentScraper) scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*models.ScrapedScene, error) {
	// construct the URL
	queryURL := queryURLParametersFromScene(scene)
	if s.definition.QueryURLReplacements != nil {
		queryURL.applyReplacements(s.definition.QueryURLReplacements)
	}
	url := queryURL.constructURL(s.definition.QueryURL)

	scraper, err := s.getJsonScraper(s.definition.Scraper)
	if err != nil {
		return nil, err
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc, url)
	return scraper.scrapeScene(ctx, q)
}

func (s *jsonFragmentScraper) scrapeByFragment(ctx context.Context, input Input) (ScrapedContent, error) {
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
	if s.definition.QueryURLReplacements != nil {
		queryURL.applyReplacements(s.definition.QueryURLReplacements)
	}
	url := queryURL.constructURL(s.definition.QueryURL)

	scraper, err := s.getJsonScraper(s.definition.Scraper)
	if err != nil {
		return nil, err
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc, url)
	return scraper.scrapeScene(ctx, q)
}

func (s *jsonFragmentScraper) scrapeImageByImage(ctx context.Context, image *models.Image) (*models.ScrapedImage, error) {
	// construct the URL
	queryURL := queryURLParametersFromImage(image)
	if s.definition.QueryURLReplacements != nil {
		queryURL.applyReplacements(s.definition.QueryURLReplacements)
	}
	url := queryURL.constructURL(s.definition.QueryURL)

	scraper, err := s.getJsonScraper(s.definition.Scraper)
	if err != nil {
		return nil, err
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc, url)
	return scraper.scrapeImage(ctx, q)
}

func (s *jsonFragmentScraper) scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*models.ScrapedGallery, error) {
	// construct the URL
	queryURL := queryURLParametersFromGallery(gallery)
	if s.definition.QueryURLReplacements != nil {
		queryURL.applyReplacements(s.definition.QueryURLReplacements)
	}
	url := queryURL.constructURL(s.definition.QueryURL)

	scraper, err := s.getJsonScraper(s.definition.Scraper)
	if err != nil {
		return nil, err
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getJsonQuery(doc, url)
	return scraper.scrapeGallery(ctx, q)
}

func (s *jsonScraper) getJsonQuery(doc string, url string) *jsonQuery {
	return &jsonQuery{
		doc:     doc,
		scraper: s,
		url:     url,
	}
}

type jsonQuery struct {
	doc       string
	scraper   *jsonScraper
	queryType QueryType
	url       string
}

func (q *jsonQuery) getType() QueryType {
	return q.queryType
}

func (q *jsonQuery) setType(t QueryType) {
	q.queryType = t
}

func (q *jsonQuery) getURL() string {
	return q.url
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

	return q.scraper.getJsonQuery(doc, value)
}
