package scraper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"

	"golang.org/x/net/html"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type xpathScraper struct {
	scraper      scraperTypeConfig
	config       config
	globalConfig GlobalConfig
	client       *http.Client
}

func newXpathScraper(scraper scraperTypeConfig, client *http.Client, config config, globalConfig GlobalConfig) *xpathScraper {
	return &xpathScraper{
		scraper:      scraper,
		config:       config,
		globalConfig: globalConfig,
		client:       client,
	}
}

func (s *xpathScraper) getXpathScraper() *mappedScraper {
	return s.config.XPathScrapers[s.scraper.Scraper]
}

func (s *xpathScraper) scrapeURL(ctx context.Context, url string) (*html.Node, *mappedScraper, error) {
	scraper := s.getXpathScraper()

	if scraper == nil {
		return nil, nil, errors.New("xpath scraper with name " + s.scraper.Scraper + " not found in config")
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, nil, err
	}

	return doc, scraper, nil
}

func (s *xpathScraper) scrapeByURL(ctx context.Context, url string, ty ScrapeContentType) (ScrapedContent, error) {
	u := replaceURL(url, s.scraper) // allow a URL Replace for performer by URL queries
	doc, scraper, err := s.scrapeURL(ctx, u)
	if err != nil {
		return nil, err
	}

	q := s.getXPathQuery(doc)
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

func (s *xpathScraper) scrapeByName(ctx context.Context, name string, ty ScrapeContentType) ([]ScrapedContent, error) {
	scraper := s.getXpathScraper()

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

	q := s.getXPathQuery(doc)
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

func (s *xpathScraper) scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*ScrapedScene, error) {
	// construct the URL
	queryURL := queryURLParametersFromScene(scene)
	if s.scraper.QueryURLReplacements != nil {
		queryURL.applyReplacements(s.scraper.QueryURLReplacements)
	}
	url := queryURL.constructURL(s.scraper.QueryURL)

	scraper := s.getXpathScraper()

	if scraper == nil {
		return nil, errors.New("xpath scraper with name " + s.scraper.Scraper + " not found in config")
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getXPathQuery(doc)
	return scraper.scrapeScene(ctx, q)
}

func (s *xpathScraper) scrapeByFragment(ctx context.Context, input Input) (ScrapedContent, error) {
	switch {
	case input.Gallery != nil:
		return nil, fmt.Errorf("%w: cannot use an xpath scraper as a gallery fragment scraper", ErrNotSupported)
	case input.Performer != nil:
		return nil, fmt.Errorf("%w: cannot use an xpath scraper as a performer fragment scraper", ErrNotSupported)
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

	scraper := s.getXpathScraper()

	if scraper == nil {
		return nil, errors.New("xpath scraper with name " + s.scraper.Scraper + " not found in config")
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getXPathQuery(doc)
	return scraper.scrapeScene(ctx, q)
}

func (s *xpathScraper) scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*ScrapedGallery, error) {
	// construct the URL
	queryURL := queryURLParametersFromGallery(gallery)
	if s.scraper.QueryURLReplacements != nil {
		queryURL.applyReplacements(s.scraper.QueryURLReplacements)
	}
	url := queryURL.constructURL(s.scraper.QueryURL)

	scraper := s.getXpathScraper()

	if scraper == nil {
		return nil, errors.New("xpath scraper with name " + s.scraper.Scraper + " not found in config")
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getXPathQuery(doc)
	return scraper.scrapeGallery(ctx, q)
}

func (s *xpathScraper) loadURL(ctx context.Context, url string) (*html.Node, error) {
	r, err := loadURL(ctx, url, s.client, s.config, s.globalConfig)
	if err != nil {
		return nil, err
	}

	ret, err := html.Parse(r)

	if err == nil && s.config.DebugOptions != nil && s.config.DebugOptions.PrintHTML {
		var b bytes.Buffer
		if err := html.Render(&b, ret); err != nil {
			logger.Warnf("could not render HTML: %v", err)
		}
		logger.Infof("loadURL (%s) response: \n%s", url, b.String())
	}

	return ret, err
}

func (s *xpathScraper) getXPathQuery(doc *html.Node) *xpathQuery {
	return &xpathQuery{
		doc:     doc,
		scraper: s,
	}
}

type xpathQuery struct {
	doc       *html.Node
	scraper   *xpathScraper
	queryType QueryType
}

func (q *xpathQuery) getType() QueryType {
	return q.queryType
}

func (q *xpathQuery) setType(t QueryType) {
	q.queryType = t
}

func (q *xpathQuery) runQuery(selector string) ([]string, error) {
	found, err := htmlquery.QueryAll(q.doc, selector)
	if err != nil {
		return nil, fmt.Errorf("selector '%s': parse error: %v", selector, err)
	}

	var ret []string
	for _, n := range found {
		// don't add empty strings
		nodeText := q.nodeText(n)
		if nodeText != "" {
			ret = append(ret, q.nodeText(n))
		}
	}

	return ret, nil
}

func (q *xpathQuery) nodeText(n *html.Node) string {
	var ret string
	if n != nil && n.Type == html.CommentNode {
		ret = htmlquery.OutputHTML(n, true)
	} else {
		ret = htmlquery.InnerText(n)
	}

	// trim all leading and trailing whitespace
	ret = strings.TrimSpace(ret)

	// remove multiple whitespace
	re := regexp.MustCompile("  +")
	ret = re.ReplaceAllString(ret, " ")

	// TODO - make this optional
	re = regexp.MustCompile("\n")
	ret = re.ReplaceAllString(ret, "")

	return ret
}

func (q *xpathQuery) subScrape(ctx context.Context, value string) mappedQuery {
	doc, err := q.scraper.loadURL(ctx, value)

	if err != nil {
		logger.Warnf("Error getting URL '%s' for sub-scraper: %s", value, err.Error())
		return nil
	}

	return q.scraper.getXPathQuery(doc)
}
