package scraper

import (
	"bytes"
	"context"
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
	txnManager   models.TransactionManager
}

func newXpathScraper(scraper scraperTypeConfig, client *http.Client, txnManager models.TransactionManager, config config, globalConfig GlobalConfig) *xpathScraper {
	return &xpathScraper{
		scraper:      scraper,
		config:       config,
		globalConfig: globalConfig,
		client:       client,
		txnManager:   txnManager,
	}
}

func (s *xpathScraper) mappedScraper() (*mappedScraper, error) {
	scraper := s.config.XPathScrapers[s.scraper.Scraper]
	if scraper == nil {
		return nil, fmt.Errorf("%w: searched for xpath scraper %v", ErrNotFound, s.scraper.Scraper)
	}

	return scraper, nil
}

func (s *xpathScraper) scrapeByURL(ctx context.Context, url string, ty models.ScrapeContentType) (models.ScrapedContent, error) {
	u := replaceURL(url, s.scraper) // allow a URL Replace for performer by URL queries
	scraper, err := s.mappedScraper()
	if err != nil {
		return nil, err
	}

	doc, err := s.loadURL(ctx, u)
	if err != nil {
		return nil, err
	}

	q := s.newDocQueryer(doc)
	switch ty {
	case models.ScrapeContentTypePerformer:
		return scraper.scrapePerformer(ctx, q)
	case models.ScrapeContentTypeScene:
		return scraper.scrapeScene(ctx, q)
	case models.ScrapeContentTypeGallery:
		return scraper.scrapeGallery(ctx, q)
	case models.ScrapeContentTypeMovie:
		return scraper.scrapeMovie(ctx, q)
	}

	return nil, ErrNotSupported
}

func (s *xpathScraper) scrapeByName(ctx context.Context, name string, ty models.ScrapeContentType) ([]models.ScrapedContent, error) {
	scraper, err := s.mappedScraper()
	if err != nil {
		return nil, err
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

	q := s.newDocQueryer(doc)

	var content []models.ScrapedContent
	switch ty {
	case models.ScrapeContentTypePerformer:
		performers, err := scraper.scrapePerformers(ctx, q)
		if err != nil {
			return nil, err
		}
		for _, p := range performers {
			content = append(content, p)
		}

		return content, nil
	case models.ScrapeContentTypeScene:
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

func (s *xpathScraper) scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*models.ScrapedScene, error) {
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

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.newDocQueryer(doc)
	return scraper.scrapeScene(ctx, q)
}

func (s *xpathScraper) scrapeByFragment(ctx context.Context, input Input) (models.ScrapedContent, error) {
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

	scraper, err := s.mappedScraper()
	if err != nil {
		return nil, err
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.newDocQueryer(doc)
	return scraper.scrapeScene(ctx, q)
}

func (s *xpathScraper) scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*models.ScrapedGallery, error) {
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

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.newDocQueryer(doc)
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

// newDocQueryer turns an xpath scraper into a queryer over doc
func (s *xpathScraper) newDocQueryer(doc *html.Node) docQueryer {
	return &xpathQuery{
		doc:     doc,
		scraper: s,
	}
}

func (s *xpathScraper) subScrape(ctx context.Context, value string) docQueryer {
	doc, err := s.loadURL(ctx, value)

	if err != nil {
		logger.Warnf("Error getting URL '%s' for sub-scraper: %s", value, err.Error())
		return nil
	}

	return s.newDocQueryer(doc)
}

type xpathQuery struct {
	doc     *html.Node
	scraper *xpathScraper
}

func (q *xpathQuery) docQuery(selector string) ([]string, error) {
	found, err := htmlquery.QueryAll(q.doc, selector)
	if err != nil {
		return nil, fmt.Errorf("selector '%s': parse error: %v", selector, err)
	}

	var ret []string
	for _, n := range found {
		// don't add empty strings
		nt := nodeText(n)
		if nt != "" {
			ret = append(ret, nt)
		}
	}

	return ret, nil
}

var (
	stripWhiteSpace = regexp.MustCompile("  +")
	stripNewLine    = regexp.MustCompile("\n")
)

func nodeText(n *html.Node) string {
	var ret string
	if n != nil && n.Type == html.CommentNode {
		ret = htmlquery.OutputHTML(n, true)
	} else {
		ret = htmlquery.InnerText(n)
	}

	// trim all leading and trailing whitespace
	ret = strings.TrimSpace(ret)

	// remove multiple whitespace
	ret = stripWhiteSpace.ReplaceAllString(ret, " ")

	// TODO - make this optional
	ret = stripNewLine.ReplaceAllString(ret, "")

	return ret
}

func (q *xpathQuery) subScrape(ctx context.Context, value string) docQueryer {
	return q.scraper.subScrape(ctx, value)
}
