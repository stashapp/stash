package scraper

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

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
	GetPythonPath() string
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
				return fmt.Errorf("%w: gave up after %d redirects", ErrMaxRedirects, maxRedirects)
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

	scrapers, err := loadScrapers(globalConfig, txnManager)
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

func loadScrapers(globalConfig GlobalConfig, txnManager models.TransactionManager) (map[string]scraper, error) {
	path := globalConfig.GetScrapersPath()
	scrapers := make(map[string]scraper)

	// Add built-in scrapers
	freeOnes := getFreeonesScraper(txnManager, globalConfig)
	autoTag := getAutoTagScraper(txnManager, globalConfig)
	scrapers[freeOnes.spec().ID] = freeOnes
	scrapers[autoTag.spec().ID] = autoTag

	logger.Debugf("Reading scraper configs from %s", path)

	scraperFiles := []string{}
	err := fsutil.SymWalk(path, func(fp string, f os.FileInfo, err error) error {
		if filepath.Ext(fp) == ".yml" {
			c, err := loadConfigFromYAMLFile(fp)
			if err != nil {
				logger.Errorf("Error loading scraper %s: %v", fp, err)
			} else {
				scraper := newGroupScraper(*c, txnManager, globalConfig)
				scrapers[scraper.spec().ID] = scraper
			}
			scraperFiles = append(scraperFiles, fp)
		}
		return nil
	})

	if err != nil {
		logger.Errorf("Error reading scraper configs: %v", err)
		return nil, err
	}

	return scrapers, nil
}

// ReloadScrapers clears the scraper cache and reloads from the scraper path.
// In the event of an error during loading, the cache will be left empty.
func (c *Cache) ReloadScrapers() error {
	c.scrapers = nil
	scrapers, err := loadScrapers(c.globalConfig, c.txnManager)
	if err != nil {
		return err
	}

	c.scrapers = scrapers
	return nil
}

// ListScrapers lists scrapers matching one of the given types.
// Returns a list of scrapers, sorted by their ID.
func (c Cache) ListScrapers(tys []ScrapeContentType) []*Scraper {
	var ret []*Scraper
	for _, s := range c.scrapers {
		for _, t := range tys {
			if s.supports(t) {
				spec := s.spec()
				ret = append(ret, &spec)
				break
			}
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].ID < ret[j].ID
	})

	return ret
}

// GetScraper returns the scraper matching the provided id.
func (c Cache) GetScraper(scraperID string) *Scraper {
	s := c.findScraper(scraperID)
	if s != nil {
		spec := s.spec()
		return &spec
	}

	return nil
}

func (c Cache) findScraper(scraperID string) scraper {
	s, ok := c.scrapers[scraperID]
	if ok {
		return s
	}

	return nil
}

func (c Cache) ScrapeName(ctx context.Context, id, query string, ty ScrapeContentType) ([]ScrapedContent, error) {
	// find scraper with the provided id
	s := c.findScraper(id)
	if s == nil {
		return nil, fmt.Errorf("%w: id %s", ErrNotFound, id)
	}
	if !s.supports(ty) {
		return nil, fmt.Errorf("%w: cannot use scraper %s as a %v scraper", ErrNotSupported, id, ty)
	}

	ns, ok := s.(nameScraper)
	if !ok {
		return nil, fmt.Errorf("%w: cannot use scraper %s to scrape by name", ErrNotSupported, id)
	}

	return ns.viaName(ctx, c.client, query, ty)
}

// ScrapeFragment uses the given fragment input to scrape
func (c Cache) ScrapeFragment(ctx context.Context, id string, input Input) (ScrapedContent, error) {
	s := c.findScraper(id)
	if s == nil {
		return nil, fmt.Errorf("%w: id %s", ErrNotFound, id)
	}

	fs, ok := s.(fragmentScraper)
	if !ok {
		return nil, fmt.Errorf("%w: cannot use scraper %s as a fragment scraper", ErrNotSupported, id)
	}

	content, err := fs.viaFragment(ctx, c.client, input)
	if err != nil {
		return nil, fmt.Errorf("error while fragment scraping with scraper %s: %w", id, err)
	}

	return c.postScrape(ctx, content)
}

// ScrapeURL scrapes a given url for the given content. Searches the scraper cache
// and picks the first scraper capable of scraping the given url into the desired
// content. Returns the scraped content or an error if the scrape fails.
func (c Cache) ScrapeURL(ctx context.Context, url string, ty ScrapeContentType) (ScrapedContent, error) {
	for _, s := range c.scrapers {
		if s.supportsURL(url, ty) {
			ul, ok := s.(urlScraper)
			if !ok {
				return nil, fmt.Errorf("%w: cannot use scraper %s as an url scraper", ErrNotSupported, s.spec().ID)
			}
			ret, err := ul.viaURL(ctx, c.client, url, ty)
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

func (c Cache) ScrapeID(ctx context.Context, scraperID string, id int, ty ScrapeContentType) (ScrapedContent, error) {
	s := c.findScraper(scraperID)
	if s == nil {
		return nil, fmt.Errorf("%w: id %s", ErrNotFound, scraperID)
	}

	if !s.supports(ty) {
		return nil, fmt.Errorf("%w: cannot use scraper %s to scrape %v content", ErrNotSupported, scraperID, ty)
	}

	var ret ScrapedContent
	switch ty {
	case ScrapeContentTypeScene:
		ss, ok := s.(sceneScraper)
		if !ok {
			return nil, fmt.Errorf("%w: cannot use scraper %s as a scene scraper", ErrNotSupported, scraperID)
		}

		scene, err := getScene(ctx, id, c.txnManager)
		if err != nil {
			return nil, fmt.Errorf("scraper %s: unable to load scene id %v: %w", scraperID, id, err)
		}

		// don't assign nil concrete pointer to ret interface, otherwise nil
		// detection is harder
		scraped, err := ss.viaScene(ctx, c.client, scene)
		if err != nil {
			return nil, fmt.Errorf("scraper %s: %w", scraperID, err)
		}

		if scraped != nil {
			ret = scraped
		}
	case ScrapeContentTypeGallery:
		gs, ok := s.(galleryScraper)
		if !ok {
			return nil, fmt.Errorf("%w: cannot use scraper %s as a gallery scraper", ErrNotSupported, scraperID)
		}

		gallery, err := getGallery(ctx, id, c.txnManager)
		if err != nil {
			return nil, fmt.Errorf("scraper %s: unable to load gallery id %v: %w", scraperID, id, err)
		}

		// don't assign nil concrete pointer to ret interface, otherwise nil
		// detection is harder
		scraped, err := gs.viaGallery(ctx, c.client, gallery)
		if err != nil {
			return nil, fmt.Errorf("scraper %s: %w", scraperID, err)
		}

		if scraped != nil {
			ret = scraped
		}
	}

	return c.postScrape(ctx, ret)
}
