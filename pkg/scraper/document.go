package scraper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
)

type docScraperType int

const (
	jsonScraper docScraperType = iota
	xpathScraper
)

// docScraper is the type of document scrapers. These work by constructing
// a URL, loading the URL and then turn the content into a document. They
// then scrape the document.
type docScraper struct {
	ty           docScraperType
	scraper      scraperTypeConfig
	config       config
	globalConfig GlobalConfig
	client       *http.Client
	txnManager   models.TransactionManager
}

func newJsonScraper(scraper scraperTypeConfig, client *http.Client, txnManager models.TransactionManager, config config, globalConfig GlobalConfig) *docScraper {
	return &docScraper{
		ty:           jsonScraper,
		scraper:      scraper,
		config:       config,
		client:       client,
		globalConfig: globalConfig,
		txnManager:   txnManager,
	}
}

func newXpathScraper(scraper scraperTypeConfig, client *http.Client, txnManager models.TransactionManager, config config, globalConfig GlobalConfig) *docScraper {
	return &docScraper{
		ty:           xpathScraper,
		scraper:      scraper,
		config:       config,
		globalConfig: globalConfig,
		client:       client,
		txnManager:   txnManager,
	}
}

func (s *docScraper) mappedScraper() (*mappedScraper, error) {
	var scraper *mappedScraper
	switch s.ty {
	case jsonScraper:
		scraper = s.config.JsonScrapers[s.scraper.Scraper]
	case xpathScraper:
		scraper = s.config.XPathScrapers[s.scraper.Scraper]
	}

	if scraper == nil {
		return nil, fmt.Errorf("%w: searched for scraper name %v", ErrNotFound, s.scraper.Scraper)
	}

	return scraper, nil

}

var ErrInvalidJSON = errors.New("invalid json")

func (s *docScraper) loadURL(ctx context.Context, url string) (docQueryer, error) {
	r, err := loadURL(ctx, url, s.client, s.config, s.globalConfig)
	if err != nil {
		return nil, err
	}

	switch s.ty {
	case jsonScraper:
		q, err := newJsonDocQueryer(s, r)
		if err == nil && s.config.DebugOptions != nil && s.config.DebugOptions.PrintHTML {
			logger.Infof("loadURL (%s) response: \n%s", url, q)
		}

		return q, err
	case xpathScraper:
		q, err := newXpathDocQueryer(s, r)
		if err == nil && s.config.DebugOptions != nil && s.config.DebugOptions.PrintHTML {
			var b bytes.Buffer
			if err := q.Render(&b); err != nil {
				logger.Warnf("could not render HTML: %v", err)
			}
			logger.Infof("loadURL (%s) response: \n%s", url, b.String())
		}

		return q, err
	}

	panic("unknown docScraperType")
}

func (s *docScraper) scrapeByURL(ctx context.Context, url string, ty models.ScrapeContentType) (models.ScrapedContent, error) {
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

func (s *docScraper) scrapeByName(ctx context.Context, name string, ty models.ScrapeContentType) ([]models.ScrapedContent, error) {
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

func (s *docScraper) scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*models.ScrapedScene, error) {
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

func (s *docScraper) scrapeByFragment(ctx context.Context, input Input) (models.ScrapedContent, error) {
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

func (s *docScraper) scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*models.ScrapedGallery, error) {
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

func (s *docScraper) subScrape(ctx context.Context, value string) docQueryer {
	doc, err := s.loadURL(ctx, value)

	if err != nil {
		logger.Warnf("Error getting URL '%s' for sub-scraper: %s", value, err.Error())
		return nil
	}

	return doc
}

type jsonQuery struct {
	doc     string
	scraper *docScraper
}

func (jq *jsonQuery) String() string {
	if jq == nil {
		return "<nil>"
	}

	return jq.doc
}

// newJsonDocQueryer turns a json scraper into a queryer over doc
func newJsonDocQueryer(s *docScraper, r io.Reader) (*jsonQuery, error) {
	doc, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	docStr := string(doc)
	if !gjson.Valid(docStr) {
		return nil, ErrInvalidJSON
	}

	jq := jsonQuery{
		doc:     docStr,
		scraper: s,
	}

	return &jq, nil
}

func (jq *jsonQuery) docQuery(selector string) ([]string, error) {
	value := gjson.Get(jq.doc, selector)

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

func (jq *jsonQuery) subScrape(ctx context.Context, value string) docQueryer {
	return jq.scraper.subScrape(ctx, value)
}

type xpathQuery struct {
	doc     *html.Node
	scraper *docScraper
}

// newXpathDocQueryer turns an xpath scraper into a queryer over doc
func newXpathDocQueryer(s *docScraper, r io.Reader) (*xpathQuery, error) {
	doc, err := html.Parse(r)

	xq := xpathQuery{
		doc:     doc,
		scraper: s,
	}

	return &xq, err
}

var ErrNilDocument = errors.New("(X)HTML document is <nil>")

func (xq *xpathQuery) Render(w io.Writer) error {
	if xq == nil {
		return ErrNilDocument
	}

	return html.Render(w, xq.doc)
}

func (xq *xpathQuery) docQuery(selector string) ([]string, error) {
	found, err := htmlquery.QueryAll(xq.doc, selector)
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

func (xq *xpathQuery) subScrape(ctx context.Context, value string) docQueryer {
	return xq.scraper.subScrape(ctx, value)
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
