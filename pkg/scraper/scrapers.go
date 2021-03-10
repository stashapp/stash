package scraper

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

// GlobalConfig contains the global scraper options.
type GlobalConfig struct {
	// User Agent used when scraping using http.
	UserAgent string

	// Path (file or remote address) to a Chrome CDP instance.
	CDPPath string
	Path    string
}

func (c GlobalConfig) isCDPPathHTTP() bool {
	return strings.HasPrefix(c.CDPPath, "http://") || strings.HasPrefix(c.CDPPath, "https://")
}

func (c GlobalConfig) isCDPPathWS() bool {
	return strings.HasPrefix(c.CDPPath, "ws://")
}

// Cache stores scraper details.
type Cache struct {
	scrapers     []config
	globalConfig GlobalConfig
	txnManager   models.TransactionManager
}

// NewCache returns a new Cache loading scraper configurations from the
// scraper path provided in the global config object. It returns a new
// instance and an error if the scraper directory could not be loaded.
//
// Scraper configurations are loaded from yml files in the provided scrapers
// directory and any subdirectories.
func NewCache(globalConfig GlobalConfig, txnManager models.TransactionManager) (*Cache, error) {
	scrapers, err := loadScrapers(globalConfig.Path)
	if err != nil {
		return nil, err
	}

	return &Cache{
		globalConfig: globalConfig,
		scrapers:     scrapers,
		txnManager:   txnManager,
	}, nil
}

func loadScrapers(path string) ([]config, error) {
	scrapers := make([]config, 0)

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
	scrapers = append(scrapers, getFreeonesScraper())

	for _, file := range scraperFiles {
		scraper, err := loadScraperFromYAMLFile(file)
		if err != nil {
			logger.Errorf("Error loading scraper %s: %s", file, err.Error())
		} else {
			scrapers = append(scrapers, *scraper)
		}
	}

	return scrapers, nil
}

// ReloadScrapers clears the scraper cache and reloads from the scraper path.
// In the event of an error during loading, the cache will be left empty.
func (c *Cache) ReloadScrapers() error {
	c.scrapers = nil
	scrapers, err := loadScrapers(c.globalConfig.Path)
	if err != nil {
		return err
	}

	c.scrapers = scrapers
	return nil
}

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
		if s.supportsPerformers() {
			ret = append(ret, s.toScraper())
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
		if s.supportsScenes() {
			ret = append(ret, s.toScraper())
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
		if s.supportsGalleries() {
			ret = append(ret, s.toScraper())
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
		if s.supportsMovies() {
			ret = append(ret, s.toScraper())
		}
	}

	return ret
}

func (c Cache) findScraper(scraperID string) *config {
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
	if s != nil {
		return s.ScrapePerformerNames(query, c.txnManager, c.globalConfig)
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}

// ScrapePerformer uses the scraper with the provided ID to scrape a
// performer using the provided performer fragment.
func (c Cache) ScrapePerformer(scraperID string, scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	// find scraper with the provided id
	s := c.findScraper(scraperID)
	if s != nil {
		ret, err := s.ScrapePerformer(scrapedPerformer, c.txnManager, c.globalConfig)
		if err != nil {
			return nil, err
		}

		// post-process - set the image if applicable
		if err := setPerformerImage(ret, c.globalConfig); err != nil {
			logger.Warnf("Could not set image using URL %s: %s", *ret.Image, err.Error())
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
		if s.matchesPerformerURL(url) {
			ret, err := s.ScrapePerformerURL(url, c.txnManager, c.globalConfig)
			if err != nil {
				return nil, err
			}

			if ret != nil {
				err = c.postScrapePerformer(ret)
				if err != nil {
					return nil, err
				}
			}

			return ret, nil
		}
	}

	return nil, nil
}

func (c Cache) postScrapePerformer(ret *models.ScrapedPerformer) error {
	if err := c.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		tqb := r.Tag()

		for _, t := range ret.Tags {
			err := MatchScrapedSceneTag(tqb, t)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	// post-process - set the image if applicable
	if err := setPerformerImage(ret, c.globalConfig); err != nil {
		logger.Warnf("Could not set image using URL %s: %s", *ret.Image, err.Error())
	}

	return nil
}

func (c Cache) postScrapeScenePerformer(ret *models.ScrapedScenePerformer) error {
	if err := c.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		tqb := r.Tag()

		for _, t := range ret.Tags {
			err := MatchScrapedSceneTag(tqb, t)
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

func (c Cache) postScrapeScene(ret *models.ScrapedScene) error {
	if err := c.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		pqb := r.Performer()
		mqb := r.Movie()
		tqb := r.Tag()
		sqb := r.Studio()

		for _, p := range ret.Performers {
			if err := c.postScrapeScenePerformer(p); err != nil {
				return err
			}

			if err := MatchScrapedScenePerformer(pqb, p); err != nil {
				return err
			}
		}

		for _, p := range ret.Movies {
			err := MatchScrapedSceneMovie(mqb, p)
			if err != nil {
				return err
			}
		}

		for _, t := range ret.Tags {
			err := MatchScrapedSceneTag(tqb, t)
			if err != nil {
				return err
			}
		}

		if ret.Studio != nil {
			err := MatchScrapedSceneStudio(sqb, ret.Studio)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	// post-process - set the image if applicable
	if err := setSceneImage(ret, c.globalConfig); err != nil {
		logger.Warnf("Could not set image using URL %s: %s", *ret.Image, err.Error())
	}

	return nil
}

func (c Cache) postScrapeGallery(ret *models.ScrapedGallery) error {
	if err := c.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		pqb := r.Performer()
		tqb := r.Tag()
		sqb := r.Studio()

		for _, p := range ret.Performers {
			err := MatchScrapedScenePerformer(pqb, p)
			if err != nil {
				return err
			}
		}

		for _, t := range ret.Tags {
			err := MatchScrapedSceneTag(tqb, t)
			if err != nil {
				return err
			}
		}

		if ret.Studio != nil {
			err := MatchScrapedSceneStudio(sqb, ret.Studio)
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

// ScrapeScene uses the scraper with the provided ID to scrape a scene.
func (c Cache) ScrapeScene(scraperID string, scene models.SceneUpdateInput) (*models.ScrapedScene, error) {
	// find scraper with the provided id
	s := c.findScraper(scraperID)
	if s != nil {
		ret, err := s.ScrapeScene(scene, c.txnManager, c.globalConfig)

		if err != nil {
			return nil, err
		}

		if ret != nil {
			err = c.postScrapeScene(ret)
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
		if s.matchesSceneURL(url) {
			ret, err := s.ScrapeSceneURL(url, c.txnManager, c.globalConfig)

			if err != nil {
				return nil, err
			}

			err = c.postScrapeScene(ret)
			if err != nil {
				return nil, err
			}

			return ret, nil
		}
	}

	return nil, nil
}

// ScrapeGallery uses the scraper with the provided ID to scrape a scene.
func (c Cache) ScrapeGallery(scraperID string, gallery models.GalleryUpdateInput) (*models.ScrapedGallery, error) {
	s := c.findScraper(scraperID)
	if s != nil {
		ret, err := s.ScrapeGallery(gallery, c.txnManager, c.globalConfig)

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
		if s.matchesGalleryURL(url) {
			ret, err := s.ScrapeGalleryURL(url, c.txnManager, c.globalConfig)

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

func matchMovieStudio(qb models.StudioReader, s *models.ScrapedMovieStudio) error {
	studio, err := qb.FindByName(s.Name, true)

	if err != nil {
		return err
	}

	if studio == nil {
		// ignore - cannot match
		return nil
	}

	id := strconv.Itoa(studio.ID)
	s.ID = &id
	return nil
}

// ScrapeMovieURL uses the first scraper it finds that matches the URL
// provided to scrape a movie. If no scrapers are found that matches
// the URL, then nil is returned.
func (c Cache) ScrapeMovieURL(url string) (*models.ScrapedMovie, error) {
	for _, s := range c.scrapers {
		if s.matchesMovieURL(url) {
			ret, err := s.ScrapeMovieURL(url, c.txnManager, c.globalConfig)
			if err != nil {
				return nil, err
			}

			if ret.Studio != nil {
				if err := c.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
					return matchMovieStudio(r.Studio(), ret.Studio)
				}); err != nil {
					return nil, err
				}
			}

			// post-process - set the image if applicable
			if err := setMovieFrontImage(ret, c.globalConfig); err != nil {
				logger.Warnf("Could not set front image using URL %s: %s", *ret.FrontImage, err.Error())
			}
			if err := setMovieBackImage(ret, c.globalConfig); err != nil {
				logger.Warnf("Could not set back image using URL %s: %s", *ret.BackImage, err.Error())
			}

			return ret, nil
		}
	}

	return nil, nil
}
