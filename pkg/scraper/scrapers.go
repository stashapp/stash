package scraper

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	stash_config "github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

var ErrMaxRedirects = errors.New("maximum number of HTTP redirects reached")

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

// Cache stores scraper details.
type Cache struct {
	client       *http.Client
	scrapers     []scraper
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

func loadScrapers(globalConfig GlobalConfig, client *http.Client, txnManager models.TransactionManager) ([]scraper, error) {
	path := globalConfig.GetScrapersPath()
	scrapers := make([]scraper, 0)

	logger.Debugf("Reading scraper configs from %s", path)
	scraperFiles := []string{}
	err := utils.SymWalk(path, func(fp string, f os.FileInfo, err error) error {
		if filepath.Ext(fp) == ".yml" {
			scraperFiles = append(scraperFiles, fp)
		}
		return nil
	})

	if err != nil {
		logger.Errorf("Error reading scraper configs: %s", err.Error())
		return nil, err
	}

	// add built-in freeones scraper
	scrapers = append(scrapers, getFreeonesScraper(client, txnManager, globalConfig), getAutoTagScraper(txnManager, globalConfig))

	for _, file := range scraperFiles {
		c, err := loadConfigFromYAMLFile(file)
		if err != nil {
			logger.Errorf("Error loading scraper %s: %s", file, err.Error())
		} else {
			scraper := createScraperFromConfig(*c, client, txnManager, globalConfig)
			scrapers = append(scrapers, scraper)
		}
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

// ListPerformerScrapers returns a list of scrapers that are capable of
// scraping performers.
func (c Cache) ListPerformerScrapers() []*models.Scraper {
	var ret []*models.Scraper
	for _, s := range c.scrapers {
		// filter on type
		if s.Performer != nil {
			ret = append(ret, s.Spec)
		}
	}

	return ret
}

// ListSceneScrapers returns a list of scrapers that are capable of
// scraping scenes.
func (c Cache) ListSceneScrapers() []*models.Scraper {
	var ret []*models.Scraper
	for _, s := range c.scrapers {
		// filter on type
		if s.Scene != nil {
			ret = append(ret, s.Spec)
		}
	}

	return ret
}

// ListGalleryScrapers returns a list of scrapers that are capable of
// scraping galleries.
func (c Cache) ListGalleryScrapers() []*models.Scraper {
	var ret []*models.Scraper
	for _, s := range c.scrapers {
		// filter on type
		if s.Gallery != nil {
			ret = append(ret, s.Spec)
		}
	}

	return ret
}

// ListMovieScrapers returns a list of scrapers that are capable of
// scraping scenes.
func (c Cache) ListMovieScrapers() []*models.Scraper {
	var ret []*models.Scraper
	for _, s := range c.scrapers {
		// filter on type
		if s.Movie != nil {
			ret = append(ret, s.Spec)
		}
	}

	return ret
}

// GetScraper returns the scraper matching the provided id.
func (c Cache) GetScraper(scraperID string) *models.Scraper {
	ret := c.findScraper(scraperID)
	if ret != nil {
		return ret.Spec
	}

	return nil
}

func (c Cache) findScraper(scraperID string) *scraper {
	for _, s := range c.scrapers {
		if s.ID == scraperID {
			return &s
		}
	}

	return nil
}

// ScrapePerformerList uses the scraper with the provided ID to query for
// performers using the provided query string. It returns a list of
// scraped performer data.
func (c Cache) ScrapePerformerList(scraperID string, query string) ([]*models.ScrapedPerformer, error) {
	// find scraper with the provided id
	s := c.findScraper(scraperID)
	if s != nil && s.Performer != nil {
		return s.Performer.scrapeByName(query)
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}

// ScrapePerformer uses the scraper with the provided ID to scrape a
// performer using the provided performer fragment.
func (c Cache) ScrapePerformer(scraperID string, scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	// find scraper with the provided id
	s := c.findScraper(scraperID)
	if s != nil && s.Performer != nil {
		ret, err := s.Performer.scrapeByFragment(scrapedPerformer)
		if err != nil {
			return nil, err
		}

		if ret != nil {
			err = c.postScrapePerformer(context.TODO(), ret)
			if err != nil {
				return nil, err
			}
		}

		return ret, nil
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}

// ScrapePerformerURL uses the first scraper it finds that matches the URL
// provided to scrape a performer. If no scrapers are found that matches
// the URL, then nil is returned.
func (c Cache) ScrapePerformerURL(url string) (*models.ScrapedPerformer, error) {
	for _, s := range c.scrapers {
		if matchesURL(s.Performer, url) {
			ret, err := s.Performer.scrapeByURL(url)
			if err != nil {
				return nil, err
			}

			if ret != nil {
				err = c.postScrapePerformer(context.TODO(), ret)
				if err != nil {
					return nil, err
				}
			}

			return ret, nil
		}
	}

	return nil, nil
}

func (c Cache) postScrapePerformer(ctx context.Context, ret *models.ScrapedPerformer) error {
	if err := c.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		tqb := r.Tag()

		tags, err := postProcessTags(tqb, ret.Tags)
		if err != nil {
			return err
		}
		ret.Tags = tags

		return nil
	}); err != nil {
		return err
	}

	// post-process - set the image if applicable
	if err := setPerformerImage(ctx, c.client, ret, c.globalConfig); err != nil {
		logger.Warnf("Could not set image using URL %s: %s", *ret.Image, err.Error())
	}

	return nil
}

func (c Cache) postScrapeScenePerformer(ret *models.ScrapedPerformer) error {
	if err := c.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		tqb := r.Tag()

		tags, err := postProcessTags(tqb, ret.Tags)
		if err != nil {
			return err
		}
		ret.Tags = tags

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (c Cache) postScrapeScene(ctx context.Context, ret *models.ScrapedScene) error {
	if err := c.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		pqb := r.Performer()
		mqb := r.Movie()
		tqb := r.Tag()
		sqb := r.Studio()

		for _, p := range ret.Performers {
			if err := c.postScrapeScenePerformer(p); err != nil {
				return err
			}

			if err := match.ScrapedPerformer(pqb, p, nil); err != nil {
				return err
			}
		}

		for _, p := range ret.Movies {
			err := match.ScrapedMovie(mqb, p)
			if err != nil {
				return err
			}
		}

		tags, err := postProcessTags(tqb, ret.Tags)
		if err != nil {
			return err
		}
		ret.Tags = tags

		if ret.Studio != nil {
			err := match.ScrapedStudio(sqb, ret.Studio, nil)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	// post-process - set the image if applicable
	if err := setSceneImage(ctx, c.client, ret, c.globalConfig); err != nil {
		logger.Warnf("Could not set image using URL %s: %v", *ret.Image, err)
	}

	return nil
}

func (c Cache) postScrapeGallery(ret *models.ScrapedGallery) error {
	if err := c.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		pqb := r.Performer()
		tqb := r.Tag()
		sqb := r.Studio()

		for _, p := range ret.Performers {
			err := match.ScrapedPerformer(pqb, p, nil)
			if err != nil {
				return err
			}
		}

		tags, err := postProcessTags(tqb, ret.Tags)
		if err != nil {
			return err
		}
		ret.Tags = tags

		if ret.Studio != nil {
			err := match.ScrapedStudio(sqb, ret.Studio, nil)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// ScrapeScene uses the scraper with the provided ID to scrape a scene using existing data.
func (c Cache) ScrapeScene(scraperID string, sceneID int) (*models.ScrapedScene, error) {
	// find scraper with the provided id
	s := c.findScraper(scraperID)
	if s != nil && s.Scene != nil {
		// get scene from id
		scene, err := getScene(sceneID, c.txnManager)
		if err != nil {
			return nil, err
		}

		ret, err := s.Scene.scrapeByScene(scene)

		if err != nil {
			return nil, err
		}

		if ret != nil {
			err = c.postScrapeScene(context.TODO(), ret)
			if err != nil {
				return nil, err
			}
		}

		return ret, nil
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}

// ScrapeSceneQuery uses the scraper with the provided ID to query for
// scenes using the provided query string. It returns a list of
// scraped scene data.
func (c Cache) ScrapeSceneQuery(scraperID string, query string) ([]*models.ScrapedScene, error) {
	// find scraper with the provided id
	s := c.findScraper(scraperID)
	if s != nil && s.Scene != nil {
		return s.Scene.scrapeByName(query)
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}

// ScrapeSceneFragment uses the scraper with the provided ID to scrape a scene.
func (c Cache) ScrapeSceneFragment(scraperID string, scene models.ScrapedSceneInput) (*models.ScrapedScene, error) {
	// find scraper with the provided id
	s := c.findScraper(scraperID)
	if s != nil && s.Scene != nil {
		ret, err := s.Scene.scrapeByFragment(scene)

		if err != nil {
			return nil, err
		}

		if ret != nil {
			err = c.postScrapeScene(context.TODO(), ret)
			if err != nil {
				return nil, err
			}
		}

		return ret, nil
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}

// ScrapeSceneURL uses the first scraper it finds that matches the URL
// provided to scrape a scene. If no scrapers are found that matches
// the URL, then nil is returned.
func (c Cache) ScrapeSceneURL(url string) (*models.ScrapedScene, error) {
	for _, s := range c.scrapers {
		if matchesURL(s.Scene, url) {
			ret, err := s.Scene.scrapeByURL(url)

			if err != nil {
				return nil, err
			}

			err = c.postScrapeScene(context.TODO(), ret)
			if err != nil {
				return nil, err
			}

			return ret, nil
		}
	}

	return nil, nil
}

// ScrapeGallery uses the scraper with the provided ID to scrape a gallery using existing data.
func (c Cache) ScrapeGallery(scraperID string, galleryID int) (*models.ScrapedGallery, error) {
	s := c.findScraper(scraperID)
	if s != nil && s.Gallery != nil {
		// get gallery from id
		gallery, err := getGallery(galleryID, c.txnManager)
		if err != nil {
			return nil, err
		}

		ret, err := s.Gallery.scrapeByGallery(gallery)

		if err != nil {
			return nil, err
		}

		if ret != nil {
			err = c.postScrapeGallery(ret)
			if err != nil {
				return nil, err
			}
		}

		return ret, nil
	}

	return nil, errors.New("Scraped with ID " + scraperID + " not found")
}

// ScrapeGalleryFragment uses the scraper with the provided ID to scrape a gallery.
func (c Cache) ScrapeGalleryFragment(scraperID string, gallery models.ScrapedGalleryInput) (*models.ScrapedGallery, error) {
	s := c.findScraper(scraperID)
	if s != nil && s.Gallery != nil {
		ret, err := s.Gallery.scrapeByFragment(gallery)

		if err != nil {
			return nil, err
		}

		if ret != nil {
			err = c.postScrapeGallery(ret)
			if err != nil {
				return nil, err
			}
		}

		return ret, nil
	}

	return nil, errors.New("Scraped with ID " + scraperID + " not found")
}

// ScrapeGalleryURL uses the first scraper it finds that matches the URL
// provided to scrape a scene. If no scrapers are found that matches
// the URL, then nil is returned.
func (c Cache) ScrapeGalleryURL(url string) (*models.ScrapedGallery, error) {
	for _, s := range c.scrapers {
		if matchesURL(s.Gallery, url) {
			ret, err := s.Gallery.scrapeByURL(url)

			if err != nil {
				return nil, err
			}

			err = c.postScrapeGallery(ret)
			if err != nil {
				return nil, err
			}

			return ret, nil
		}
	}

	return nil, nil
}

// ScrapeMovieURL uses the first scraper it finds that matches the URL
// provided to scrape a movie. If no scrapers are found that matches
// the URL, then nil is returned.
func (c Cache) ScrapeMovieURL(url string) (*models.ScrapedMovie, error) {
	for _, s := range c.scrapers {
		if s.Movie != nil && matchesURL(s.Movie, url) {
			ret, err := s.Movie.scrapeByURL(url)
			if err != nil {
				return nil, err
			}

			if ret.Studio != nil {
				if err := c.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
					return match.ScrapedStudio(r.Studio(), ret.Studio, nil)
				}); err != nil {
					return nil, err
				}
			}

			// post-process - set the image if applicable
			if err := setMovieFrontImage(context.TODO(), c.client, ret, c.globalConfig); err != nil {
				logger.Warnf("Could not set front image using URL %s: %s", *ret.FrontImage, err.Error())
			}
			if err := setMovieBackImage(context.TODO(), c.client, ret, c.globalConfig); err != nil {
				logger.Warnf("Could not set back image using URL %s: %s", *ret.BackImage, err.Error())
			}

			return ret, nil
		}
	}

	return nil, nil
}

func postProcessTags(tqb models.TagReader, scrapedTags []*models.ScrapedTag) ([]*models.ScrapedTag, error) {
	var ret []*models.ScrapedTag

	excludePatterns := stash_config.GetInstance().GetScraperExcludeTagPatterns()
	var excludeRegexps []*regexp.Regexp

	for _, excludePattern := range excludePatterns {
		reg, err := regexp.Compile(strings.ToLower(excludePattern))
		if err != nil {
			logger.Errorf("Invalid tag exclusion pattern :%v", err)
		} else {
			excludeRegexps = append(excludeRegexps, reg)
		}
	}

	var ignoredTags []string
ScrapeTag:
	for _, t := range scrapedTags {
		for _, reg := range excludeRegexps {
			if reg.MatchString(strings.ToLower(t.Name)) {
				ignoredTags = append(ignoredTags, t.Name)
				continue ScrapeTag
			}
		}

		err := match.ScrapedTag(tqb, t)
		if err != nil {
			return nil, err
		}
		ret = append(ret, t)
	}

	if len(ignoredTags) > 0 {
		logger.Infof("Scraping ignored tags: %s", strings.Join(ignoredTags, ", "))
	}

	return ret, nil
}
