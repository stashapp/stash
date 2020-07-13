package scraper

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

// Timeout for the scrape http request. Includes transfer time. May want to make this
// configurable at some point.
const scrapeGetTimeout = time.Second * 30

type xpathScraper struct {
	scraper      scraperTypeConfig
	config       config
	globalConfig GlobalConfig
}

func newXpathScraper(scraper scraperTypeConfig, config config, globalConfig GlobalConfig) *xpathScraper {
	return &xpathScraper{
		scraper:      scraper,
		config:       config,
		globalConfig: globalConfig,
	}
}

func (s *xpathScraper) getXpathScraper() *mappedScraper {
	return s.config.XPathScrapers[s.scraper.Scraper]
}

func (s *xpathScraper) scrapeURL(url string) (*html.Node, *mappedScraper, error) {
	scraper := s.getXpathScraper()

	if scraper == nil {
		return nil, nil, errors.New("xpath scraper with name " + s.scraper.Scraper + " not found in config")
	}

	doc, err := s.loadURL(url)

	if err != nil {
		return nil, nil, err
	}

	return doc, scraper, nil
}

func (s *xpathScraper) scrapePerformerByURL(url string) (*models.ScrapedPerformer, error) {
	doc, scraper, err := s.scrapeURL(url)
	if err != nil {
		return nil, err
	}

	q := s.getXPathQuery(doc)
	return scraper.scrapePerformer(q)
}

func (s *xpathScraper) scrapeSceneByURL(url string) (*models.ScrapedScene, error) {
	doc, scraper, err := s.scrapeURL(url)
	if err != nil {
		return nil, err
	}

	q := s.getXPathQuery(doc)
	return scraper.scrapeScene(q)
}

func (s *xpathScraper) scrapePerformersByName(name string) ([]*models.ScrapedPerformer, error) {
	scraper := s.getXpathScraper()

	if scraper == nil {
		return nil, errors.New("xpath scraper with name " + s.scraper.Scraper + " not found in config")
	}

	const placeholder = "{}"

	// replace the placeholder string with the URL-escaped name
	escapedName := url.QueryEscape(name)

	url := s.scraper.QueryURL
	url = strings.Replace(url, placeholder, escapedName, -1)

	doc, err := s.loadURL(url)

	if err != nil {
		return nil, err
	}

	q := s.getXPathQuery(doc)
	return scraper.scrapePerformers(q)
}

func (s *xpathScraper) scrapePerformerByFragment(scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	return nil, errors.New("scrapePerformerByFragment not supported for xpath scraper")
}

func (s *xpathScraper) scrapeSceneByFragment(scene models.SceneUpdateInput) (*models.ScrapedScene, error) {
	return nil, errors.New("scrapeSceneByFragment not supported for xpath scraper")
}

func (s *xpathScraper) loadURL(url string) (*html.Node, error) {
	client := &http.Client{
		Timeout: scrapeGetTimeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	userAgent := s.globalConfig.UserAgent
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	ret, err := html.Parse(r)

	if err == nil && s.config.DebugOptions != nil && s.config.DebugOptions.PrintHTML {
		var b bytes.Buffer
		html.Render(&b, ret)
		logger.Infof("loadURL (%s) response: \n%s", url, b.String())
	}

	return ret, err
}

func (s *xpathScraper) getXPathQuery(doc *html.Node) *xpathQuery {
	return &xpathQuery{
		doc: doc,
	}
}

type xpathQuery struct {
	doc     *html.Node
	scraper *xpathScraper
}

func (q *xpathQuery) runQuery(selector string) []string {
	found, err := htmlquery.QueryAll(q.doc, selector)
	if err != nil {
		logger.Warnf("Error parsing xpath expression '%s': %s", selector, err.Error())
		return nil
	}

	var ret []string
	for _, n := range found {
		ret = append(ret, nodeText(n))
	}

	return ret
}

func (q *xpathQuery) subScrape(value string) mappedQuery {
	doc, err := q.scraper.loadURL(value)

	if err != nil {
		logger.Warnf("Error getting URL '%s' for sub-scraper: %s", value, err.Error())
		return nil
	}

	return q.scraper.getXPathQuery(doc)
}

func commonPostProcess(value string) string {
	value = strings.TrimSpace(value)

	// remove multiple whitespace and end lines
	re := regexp.MustCompile("\n")
	value = re.ReplaceAllString(value, "")
	re = regexp.MustCompile("  +")
	value = re.ReplaceAllString(value, " ")

	return value
}

// func replaceLines replaces all newlines ("\n") with alert ("\a")
func replaceLines(value string) string {
	re := regexp.MustCompile("\a")         // \a shouldn't exist in the string
	value = re.ReplaceAllString(value, "") // remove it
	re = regexp.MustCompile("\n")          // replace newlines with (\a)'s so that they don't get removed by commonPostprocess
	value = re.ReplaceAllString(value, "\a")

	return value
}

// func restoreLines replaces all alerts ("\a") with newlines ("\n")
func restoreLines(value string) string {
	re := regexp.MustCompile("\a")
	value = re.ReplaceAllString(value, "\n")

	return value
}

func nodeText(n *html.Node) string {
	if n != nil && n.Type == html.CommentNode {
		return htmlquery.OutputHTML(n, true)
	}
	return htmlquery.InnerText(n)
}
