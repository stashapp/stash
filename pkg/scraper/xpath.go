package scraper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	jsoniter "github.com/json-iterator/go"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"golang.org/x/net/publicsuffix"

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
	var r io.Reader

	driverOptions := s.config.DriverOptions
	if driverOptions != nil && driverOptions.UseCDP {
		// get the page using chrome dp
		resp, err := urlFromCDP(url, *driverOptions)
		if err != nil {
			return nil, err
		}
		r = strings.NewReader(resp)
		if err != nil {
			return nil, err
		}
	} else {
		// get the page using http.Client
		options := cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		}
		jar, er := cookiejar.New(&options)
		if er != nil {
			return nil, er
		}

		client := &http.Client{
			Timeout: scrapeGetTimeout,
			// defaultCheckRedirect code with max changed from 10 to 20
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 20 {
					return errors.New("stopped after 20 redirects")
				}
				return nil
			},
			Jar: jar,
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

		r, err = charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
		if err != nil {
			return nil, err
		}
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
		doc:     doc,
		scraper: s,
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
		// don't add empty strings
		nodeText := q.nodeText(n)
		if nodeText != "" {
			ret = append(ret, q.nodeText(n))
		}
	}

	return ret
}

func (q *xpathQuery) nodeText(n *html.Node) string {
	var ret string
	if n != nil && n.Type == html.CommentNode {
		ret = htmlquery.OutputHTML(n, true)
	}
	ret = htmlquery.InnerText(n)

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

func (q *xpathQuery) subScrape(value string) mappedQuery {
	doc, err := q.scraper.loadURL(value)

	if err != nil {
		logger.Warnf("Error getting URL '%s' for sub-scraper: %s", value, err.Error())
		return nil
	}

	return q.scraper.getXPathQuery(doc)
}

// func urlFromCDP uses chrome cdp and DOM to load and process the url
// if remote is set as true in the scraperConfig  it will try to use localhost:9222
// else it will look for google-chrome in path
func urlFromCDP(url string, driverOptions scraperDriverOptions) (string, error) {
	if !driverOptions.UseCDP {
		return "", fmt.Errorf("Url shouldn't be feetched through CDP")
	}

	remote := false
	sleep := 2

	if driverOptions.Remote {
		remote = true
	}

	if driverOptions.Sleep != 0 {
		sleep = driverOptions.Sleep
	}

	sleepDuration := time.Duration(sleep) * time.Second
	act := context.Background()

	if remote {
		var cancelAct context.CancelFunc
		remote, errCDP := getRemoteCDP()
		if errCDP != nil {
			return "", errCDP
		}
		act, cancelAct = chromedp.NewRemoteAllocator(context.Background(), remote)
		defer cancelAct()
	}

	ctx, cancel := chromedp.NewContext(act)
	defer cancel()

	var res string
	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.Sleep(sleepDuration),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return "", err
	}

	return res, nil
}

// func getRemoteCDP returns the complete remote address that is required to access the cdp instance
func getRemoteCDP() (string, error) {
	resp, err := http.Get("http://localhost:9222/json/version")
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	remote := result["webSocketDebuggerUrl"].(string)
	logger.Debugf("Remote cdp instance found %s", remote)
	return remote, err
}

func cdpNetwork(enable bool) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		if enable {
			network.Enable().Do(ctx)
		} else {
			network.Disable().Do(ctx)
		}
		return nil
	})
}
