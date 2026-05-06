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
	definition   Definition
	globalConfig GlobalConfig
	client       *http.Client
}

func (s *xpathScraper) getXpathScraper(name string) (*mappedScraper, error) {
	ret, ok := s.definition.XPathScrapers[name]
	if !ok {
		return nil, fmt.Errorf("xpath scraper with name %s not found in config", name)
	}
	return &ret, nil
}

type xpathURLScraper struct {
	xpathScraper
	definition ByURLDefinition
}

func (s *xpathURLScraper) scrapeByURL(ctx context.Context, url string, ty ScrapeContentType) (ScrapedContent, error) {
	scraper, err := s.getXpathScraper(s.definition.Scraper)
	if err != nil {
		return nil, err
	}

	doc, err := s.loadURL(ctx, url)
	if err != nil {
		return nil, err
	}

	q := s.getXPathQuery(doc, url)
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

type xpathNameScraper struct {
	xpathScraper
	definition ByNameDefinition
}

func (s *xpathNameScraper) scrapeByName(ctx context.Context, name string, ty ScrapeContentType) ([]ScrapedContent, error) {
	scraper, err := s.getXpathScraper(s.definition.Scraper)
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

	q := s.getXPathQuery(doc, url)
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

type xpathFragmentScraper struct {
	xpathScraper
	definition ByFragmentDefinition
}

func (s *xpathFragmentScraper) scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*models.ScrapedScene, error) {
	// construct the URL
	queryURL := queryURLParametersFromScene(scene)
	if s.definition.QueryURLReplacements != nil {
		queryURL.applyReplacements(s.definition.QueryURLReplacements)
	}
	url := queryURL.constructURL(s.definition.QueryURL)

	scraper, err := s.getXpathScraper(s.definition.Scraper)
	if err != nil {
		return nil, err
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getXPathQuery(doc, url)
	return scraper.scrapeScene(ctx, q)
}

func (s *xpathFragmentScraper) scrapeByFragment(ctx context.Context, input Input) (ScrapedContent, error) {
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
	if s.definition.QueryURLReplacements != nil {
		queryURL.applyReplacements(s.definition.QueryURLReplacements)
	}
	url := queryURL.constructURL(s.definition.QueryURL)

	scraper, err := s.getXpathScraper(s.definition.Scraper)
	if err != nil {
		return nil, err
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getXPathQuery(doc, url)
	return scraper.scrapeScene(ctx, q)
}

func (s *xpathFragmentScraper) scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*models.ScrapedGallery, error) {
	// construct the URL
	queryURL := queryURLParametersFromGallery(gallery)
	if s.definition.QueryURLReplacements != nil {
		queryURL.applyReplacements(s.definition.QueryURLReplacements)
	}
	url := queryURL.constructURL(s.definition.QueryURL)

	scraper, err := s.getXpathScraper(s.definition.Scraper)
	if err != nil {
		return nil, err
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getXPathQuery(doc, url)
	return scraper.scrapeGallery(ctx, q)
}

func (s *xpathFragmentScraper) scrapeImageByImage(ctx context.Context, image *models.Image) (*models.ScrapedImage, error) {
	// construct the URL
	queryURL := queryURLParametersFromImage(image)
	if s.definition.QueryURLReplacements != nil {
		queryURL.applyReplacements(s.definition.QueryURLReplacements)
	}
	url := queryURL.constructURL(s.definition.QueryURL)

	scraper, err := s.getXpathScraper(s.definition.Scraper)
	if err != nil {
		return nil, err
	}

	doc, err := s.loadURL(ctx, url)

	if err != nil {
		return nil, err
	}

	q := s.getXPathQuery(doc, url)
	return scraper.scrapeImage(ctx, q)
}

func (s *xpathScraper) loadURL(ctx context.Context, url string) (*html.Node, error) {
	r, err := loadURL(ctx, url, s.client, s.definition, s.globalConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load URL %q: %w", url, err)
	}

	ret, err := html.Parse(r)

	if err == nil && s.definition.DebugOptions != nil && s.definition.DebugOptions.PrintHTML {
		var b bytes.Buffer
		if err := html.Render(&b, ret); err != nil {
			logger.Warnf("could not render HTML: %v", err)
		}
		logger.Infof("loadURL (%s) response: \n%s", url, b.String())
	}

	return ret, err
}

func (s *xpathScraper) getXPathQuery(doc *html.Node, url string) *xpathQuery {
	return &xpathQuery{
		doc:     doc,
		scraper: s,
		url:     url,
	}
}

type xpathQuery struct {
	doc       *html.Node
	scraper   *xpathScraper
	queryType QueryType
	url       string
}

func (q *xpathQuery) getType() QueryType {
	return q.queryType
}

func (q *xpathQuery) setType(t QueryType) {
	q.queryType = t
}

func (q *xpathQuery) getURL() string {
	return q.url
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

	return q.scraper.getXPathQuery(doc, value)
}
