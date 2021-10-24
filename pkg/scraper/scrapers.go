package scraper

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

var ErrMaxRedirects = errors.New("maximum number of HTTP redirects reached")
var ErrNotFound = errors.New("scraper not found")
var ErrNotSupported = errors.New("not supported")

const (
	// scrapeGetTimeout is the timeout for scraper HTTP requests. Includes transfer time.
	// We may want to bump this at some point and use local context-timeouts if more granularity
	// is needed.
	scrapeGetTimeout = time.Second * 60

	// maxIdleConnsPerHost is the maximum number of idle connections the HTTP client will
	// keep on a per-host basis.
	maxIdleConnsPerHost = 8

	// maxRedirects defines the maximum number of redirects the HTTP client will follow
	maxRedirects = 20
)

// GlobalConfig contains the global scraper options.
type GlobalConfig interface {
	GetScraperUserAgent() string
	GetScrapersPath() string
	GetScraperCDPPath() string
	GetScraperCertCheck() bool
}

func isCDPPathHTTP(c GlobalConfig) bool {
	return strings.HasPrefix(c.GetScraperCDPPath(), "http://") || strings.HasPrefix(c.GetScraperCDPPath(), "https://")
}

func isCDPPathWS(c GlobalConfig) bool {
	return strings.HasPrefix(c.GetScraperCDPPath(), "ws://")
}

// Cache stores the database of scrapers
type Cache struct {
	client       *http.Client
	scrapers     map[string]scraper // Scraper ID -> Scraper
	globalConfig GlobalConfig
	txnManager   models.TransactionManager
}

// newClient creates a scraper-local http client we use throughout the scraper subsystem.
func newClient(gc GlobalConfig) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{ // ignore insecure certificates
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: !gc.GetScraperCertCheck()},
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
		},
		Timeout: scrapeGetTimeout,
		// defaultCheckRedirect code with max changed from 10 to maxRedirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return fmt.Errorf("after %d redirects: %w", maxRedirects, ErrMaxRedirects)
			}
			return nil
		},
	}

	return client
}

// NewCache returns a new Cache loading scraper configurations from the
// scraper path provided in the global config object. It returns a new
// instance and an error if the scraper directory could not be loaded.
//
// Scraper configurations are loaded from yml files in the provided scrapers
// directory and any subdirectories.
func NewCache(globalConfig GlobalConfig, txnManager models.TransactionManager) (*Cache, error) {
	// HTTP Client setup
	client := newClient(globalConfig)

	scrapers, err := loadScrapers(globalConfig, client, txnManager)
	if err != nil {
		return nil, err
	}

	return &Cache{
		client:       client,
		globalConfig: globalConfig,
		scrapers:     scrapers,
		txnManager:   txnManager,
	}, nil
}

func loadScrapers(globalConfig GlobalConfig, client *http.Client, txnManager models.TransactionManager) (map[string]scraper, error) {
	path := globalConfig.GetScrapersPath()
	scrapers := make(map[string]scraper)

	// Add built-in scrapers
	freeOnes := getFreeonesScraper(client, txnManager, globalConfig)
	autoTag := getAutoTagScraper(txnManager, globalConfig)
	scrapers[freeOnes.spec().ID] = freeOnes
	scrapers[autoTag.spec().ID] = autoTag

	logger.Debugf("Reading scraper configs from %s", path)

	scraperFiles := []string{}
	err := utils.SymWalk(path, func(fp string, f os.FileInfo, err error) error {
		if filepath.Ext(fp) == ".yml" {
			c, err := loadConfigFromYAMLFile(fp)
			if err != nil {
				logger.Errorf("Error loading scraper %s: %v", fp, err)
			} else {
				scraper := createScraperFromConfig(*c, client, txnManager, globalConfig)
				scrapers[scraper.spec().ID] = scraper
			}
			scraperFiles = append(scraperFiles, fp)
		}
		return nil
	})

	if err != nil {
		logger.Errorf("Error reading scraper configs: %s", err.Error())
		return nil, err
	}

	return scrapers, nil
}

// ReloadScrapers clears the scraper cache and reloads from the scraper path.
// In the event of an error during loading, the cache will be left empty.
func (c *Cache) ReloadScrapers() error {
	c.scrapers = nil
	scrapers, err := loadScrapers(c.globalConfig, c.client, c.txnManager)
	if err != nil {
		return err
	}

	c.scrapers = scrapers
	return nil
}

// TODO - don't think this is needed
// UpdateConfig updates the global config for the cache. If the scraper path
// has changed, ReloadScrapers will need to be called separately.
func (c *Cache) UpdateConfig(globalConfig GlobalConfig) {
	c.globalConfig = globalConfig
}

// ListScrapers returns scrapers matching a given kind
func (c Cache) ListScrapers(k models.ScrapeContentType) []*models.Scraper {
	var ret []*models.Scraper
	for _, s := range c.scrapers {
		if s.supports(k) {
			spec := s.spec()
			ret = append(ret, &spec)
		}
	}

	return ret
}

func (c Cache) findScraper(scraperID string) scraper {
	s, ok := c.scrapers[scraperID]
	if ok {
		return s
	}

	return nil
}

func (c Cache) ScrapeName(id, query string, ty models.ScrapeContentType) ([]models.ScrapedContent, error) {
	// find scraper with the provided id
	s := c.findScraper(id)
	if s == nil {
		return nil, fmt.Errorf("scraper with id %s: %w", id, ErrNotFound)
	}
	if !s.supports(ty) {
		return nil, fmt.Errorf("scraping %v with scraper %s: %w", ty, id, ErrNotSupported)
	}

	ns, ok := s.(nameScraper)
	if !ok {
		return nil, fmt.Errorf("name-scraping with scraper %s: %w", id, ErrNotSupported)
	}

	return ns.loadByName(query, ty)
}

// ScrapeFragment uses the given fragment input to scrape
func (c Cache) ScrapeFragment(ctx context.Context, id string, input Input) (models.ScrapedContent, error) {
	s := c.findScraper(id)
	if s == nil {
		return nil, fmt.Errorf("scraper %s: %w", id, ErrNotFound)
	}

	fs, ok := s.(fragmentScraper)
	if !ok {
		return nil, fmt.Errorf("fragment scraping with scraper %s: %w", id, ErrNotSupported)
	}

	content, err := fs.loadByFragment(c.client, input)
	if err != nil {
		return nil, fmt.Errorf("fragment scraping with scraper %s: %w", id, err)
	}

	return c.postScrape(ctx, content)
}

// ScrapeURL scrapes a given url for the given content. Searches the scraper cache
// and picks the first scraper capable of scraping the given url into the desired
// content. Returns the scraped content or an error if the scrape fails.
func (c Cache) ScrapeURL(ctx context.Context, url string, ty models.ScrapeContentType) (models.ScrapedContent, error) {
	for _, s := range c.scrapers {
		if s.supportsURL(url, ty) {
			ul, ok := s.(urlScraper)
			if !ok {
				return nil, fmt.Errorf("scraper with id %s used as url scraper: %w", s.spec().ID, ErrNotSupported)
			}
			ret, err := ul.loadByURL(c.client, url, ty)
			if err != nil {
				return nil, err
			}

			if ret == nil {
				return ret, nil
			}

			return c.postScrape(ctx, ret)
		}
	}

	return nil, nil
}

func (c Cache) ScrapeID(ctx context.Context, scraperID string, id int, ty models.ScrapeContentType) (models.ScrapedContent, error) {
	s := c.findScraper(scraperID)
	if s == nil {
		return nil, fmt.Errorf("scraper %s: %w", scraperID, ErrNotFound)
	}

	if !s.supports(ty) {
		return nil, fmt.Errorf("scraper %s: scraping for %v: %w", scraperID, ty, ErrNotSupported)
	}

	var ret models.ScrapedContent
	switch ty {
	case models.ScrapeContentTypeScene:
		ss, ok := s.(sceneLoader)
		if !ok {
			return nil, fmt.Errorf("scraper with id %s used as scene scraper: %w", scraperID, ErrNotSupported)
		}

		scene, err := getScene(id, c.txnManager)
		if err != nil {
			return nil, fmt.Errorf("scraper %s: unable to load scene id %v: %w", scraperID, id, err)
		}

		ret, err = ss.loadByScene(scene)
		if err != nil {
			return nil, fmt.Errorf("scraper %s: %w", scraperID, err)
		}
	case models.ScrapeContentTypeGallery:
		gs, ok := s.(galleryLoader)
		if !ok {
			return nil, fmt.Errorf("scraper with id %s used as a gallery scraper: %w", scraperID, ErrNotSupported)
		}

		gallery, err := getGallery(id, c.txnManager)
		if err != nil {
			return nil, fmt.Errorf("scraper %s: unable to load gallery id %v: %w", scraperID, id, err)
		}

		ret, err = gs.loadByGallery(gallery)
		if err != nil {
			return nil, fmt.Errorf("scraper %s: %w", scraperID, err)
		}
	}

	return c.postScrape(ctx, ret)
}
