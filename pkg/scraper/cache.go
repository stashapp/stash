package scraper

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
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
	GetProxy() string
	GetScraperExcludeTagPatterns() []string
}

func isCDPPathHTTP(c GlobalConfig) bool {
	return strings.HasPrefix(c.GetScraperCDPPath(), "http://") || strings.HasPrefix(c.GetScraperCDPPath(), "https://")
}

func isCDPPathWS(c GlobalConfig) bool {
	return strings.HasPrefix(c.GetScraperCDPPath(), "ws://")
}

type SceneFinder interface {
	models.SceneGetter
	models.URLLoader
	models.VideoFileLoader
}

type PerformerFinder interface {
	models.PerformerAutoTagQueryer
	match.PerformerFinder
}

type StudioFinder interface {
	models.StudioAutoTagQueryer
	FindByStashID(ctx context.Context, stashID models.StashID) ([]*models.Studio, error)
}

type TagFinder interface {
	models.TagGetter
	models.TagAutoTagQueryer
}

type GalleryFinder interface {
	models.GalleryGetter
	models.FileLoader
	models.URLLoader
}

type ImageFinder interface {
	models.ImageGetter
	models.FileLoader
	models.URLLoader
}

type Repository struct {
	TxnManager models.TxnManager

	SceneFinder     SceneFinder
	GalleryFinder   GalleryFinder
	ImageFinder     ImageFinder
	TagFinder       TagFinder
	PerformerFinder PerformerFinder
	GroupFinder     match.GroupNamesFinder
	StudioFinder    StudioFinder
}

func NewRepository(repo models.Repository) Repository {
	return Repository{
		TxnManager:      repo.TxnManager,
		SceneFinder:     repo.Scene,
		GalleryFinder:   repo.Gallery,
		ImageFinder:     repo.Image,
		TagFinder:       repo.Tag,
		PerformerFinder: repo.Performer,
		GroupFinder:     repo.Group,
		StudioFinder:    repo.Studio,
	}
}

func (r *Repository) WithReadTxn(ctx context.Context, fn txn.TxnFunc) error {
	return txn.WithReadTxn(ctx, r.TxnManager, fn)
}

// Cache stores the database of scrapers
type Cache struct {
	client       *http.Client
	scrapers     map[string]scraper // Scraper ID -> Scraper
	globalConfig GlobalConfig

	repository Repository
}

// newClient creates a scraper-local http client we use throughout the scraper subsystem.
func newClient(gc GlobalConfig) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{ // ignore insecure certificates
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: !gc.GetScraperCertCheck()},
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
			Proxy:               http.ProxyFromEnvironment,
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

// NewCache returns a new Cache.
//
// Scraper configurations are loaded from yml files in the scrapers
// directory in the config and any subdirectories.
//
// Does not load scrapers. Scrapers will need to be
// loaded explicitly using ReloadScrapers.
func NewCache(globalConfig GlobalConfig, repo Repository) *Cache {
	// HTTP Client setup
	client := newClient(globalConfig)

	return &Cache{
		client:       client,
		globalConfig: globalConfig,
		repository:   repo,
	}
}

// ReloadScrapers clears the scraper cache and reloads from the scraper path.
// If a scraper cannot be loaded, an error is logged and the scraper is skipped.
func (c *Cache) ReloadScrapers() {
	path := c.globalConfig.GetScrapersPath()
	scrapers := make(map[string]scraper)

	// Add built-in scrapers
	freeOnes := getFreeonesScraper(c.globalConfig)
	autoTag := getAutoTagScraper(c.repository, c.globalConfig)
	scrapers[freeOnes.spec().ID] = freeOnes
	scrapers[autoTag.spec().ID] = autoTag

	logger.Debugf("Reading scraper configs from %s", path)

	err := fsutil.SymWalk(path, func(fp string, f os.FileInfo, err error) error {
		if filepath.Ext(fp) == ".yml" {
			conf, err := loadConfigFromYAMLFile(fp)
			if err != nil {
				logger.Errorf("Error loading scraper %s: %v", fp, err)
			} else {
				scraper := newGroupScraper(*conf, c.globalConfig)
				scrapers[scraper.spec().ID] = scraper
			}
		}
		return nil
	})

	if err != nil {
		logger.Errorf("Error reading scraper configs: %v", err)
	}

	c.scrapers = scrapers
}

// ListScrapers lists scrapers matching one of the given types.
// Returns a list of scrapers, sorted by their name.
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
		return strings.ToLower(ret[i].Name) < strings.ToLower(ret[j].Name)
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

func (c Cache) compileExcludeTagPatterns() []*regexp.Regexp {
	return CompileExclusionRegexps(c.globalConfig.GetScraperExcludeTagPatterns())
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

	content, err := ns.viaName(ctx, c.client, query, ty)
	if err != nil {
		return nil, fmt.Errorf("error while name scraping with scraper %s: %w", id, err)
	}

	pp := postScraper{
		Cache:        c,
		excludeTagRE: c.compileExcludeTagPatterns(),
	}
	if err := c.repository.WithReadTxn(ctx, func(ctx context.Context) error {
		for i, cc := range content {
			content[i], err = pp.postScrape(ctx, cc)
			if err != nil {
				return fmt.Errorf("error while post-scraping with scraper %s: %w", id, err)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	LogIgnoredTags(pp.ignoredTags)

	return content, nil
}

// ScrapeFragment uses the given fragment input to scrape
func (c Cache) ScrapeFragment(ctx context.Context, id string, input Input) (ScrapedContent, error) {
	// set the deprecated URL field if it's not set
	input.populateURL()

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

	return c.postScrapeSingle(ctx, content)
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

			return c.postScrapeSingle(ctx, ret)
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

		scene, err := c.getScene(ctx, id)
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

		gallery, err := c.getGallery(ctx, id)
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

	case ScrapeContentTypeImage:
		is, ok := s.(imageScraper)
		if !ok {
			return nil, fmt.Errorf("%w: cannot use scraper %s as a image scraper", ErrNotSupported, scraperID)
		}

		scene, err := c.getImage(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("scraper %s: unable to load image id %v: %w", scraperID, id, err)
		}

		// don't assign nil concrete pointer to ret interface, otherwise nil
		// detection is harder
		scraped, err := is.viaImage(ctx, c.client, scene)
		if err != nil {
			return nil, fmt.Errorf("scraper %s: %w", scraperID, err)
		}

		if scraped != nil {
			ret = scraped
		}
	}

	return c.postScrapeSingle(ctx, ret)
}

func (c Cache) getScene(ctx context.Context, sceneID int) (*models.Scene, error) {
	var ret *models.Scene
	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		qb := r.SceneFinder

		var err error
		ret, err = qb.Find(ctx, sceneID)
		if err != nil {
			return err
		}

		if ret == nil {
			return fmt.Errorf("scene with id %d not found", sceneID)
		}

		if err := ret.LoadURLs(ctx, qb); err != nil {
			return err
		}

		if err := ret.LoadFiles(ctx, qb); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c Cache) getGallery(ctx context.Context, galleryID int) (*models.Gallery, error) {
	var ret *models.Gallery
	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		qb := r.GalleryFinder

		var err error
		ret, err = qb.Find(ctx, galleryID)
		if err != nil {
			return err
		}

		if ret == nil {
			return fmt.Errorf("gallery with id %d not found", galleryID)
		}

		if err := ret.LoadURLs(ctx, qb); err != nil {
			return err
		}

		if err := ret.LoadFiles(ctx, qb); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c Cache) getImage(ctx context.Context, imageID int) (*models.Image, error) {
	var ret *models.Image
	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		qb := r.ImageFinder

		var err error
		ret, err = qb.Find(ctx, imageID)
		if err != nil {
			return err
		}

		if ret == nil {
			return fmt.Errorf("image with id %d not found", imageID)
		}

		err = ret.LoadFiles(ctx, qb)
		if err != nil {
			return err
		}

		return ret.LoadURLs(ctx, qb)
	}); err != nil {
		return nil, err
	}
	return ret, nil
}
