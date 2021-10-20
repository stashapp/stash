package scraper

import (
	"net/http"

	"github.com/stashapp/stash/pkg/models"
)

type configSceneScraper struct {
	*configScraper
}

func (c *configSceneScraper) matchesURL(url string) bool {
	return c.config.matchesSceneURL(url)
}

func (c *configSceneScraper) scrapeByName(name string) ([]*models.ScrapedScene, error) {
	if c.config.SceneByName != nil {
		s := c.config.getScraper(*c.config.SceneByName, c.client, c.txnManager, c.globalConfig)
		return s.scrapeScenesByName(name)
	}

	return nil, nil
}

func (c *configSceneScraper) scrapeByScene(scene *models.Scene) (*models.ScrapedScene, error) {
	if c.config.SceneByFragment != nil {
		s := c.config.getScraper(*c.config.SceneByFragment, c.client, c.txnManager, c.globalConfig)
		return s.scrapeSceneByScene(scene)
	}

	return nil, nil
}

func (c *configSceneScraper) scrapeByFragment(scene models.ScrapedSceneInput) (*models.ScrapedScene, error) {
	if c.config.SceneByQueryFragment != nil {
		s := c.config.getScraper(*c.config.SceneByQueryFragment, c.client, c.txnManager, c.globalConfig)
		return s.scrapeSceneByFragment(scene)
	}

	return nil, nil
}

func (c *configSceneScraper) scrapeByURL(url string) (*models.ScrapedScene, error) {
	for _, scraper := range c.config.SceneByURL {
		if scraper.matchesURL(url) {
			s := c.config.getScraper(scraper.scraperTypeConfig, c.client, c.txnManager, c.globalConfig)
			ret, err := s.scrapeSceneByURL(url)
			if err != nil {
				return nil, err
			}

			if ret != nil {
				return ret, nil
			}
		}
	}

	return nil, nil
}

type configPerformerScraper struct {
	*configScraper
}

func (c *configPerformerScraper) matchesURL(url string) bool {
	return c.config.matchesPerformerURL(url)
}

func (c *configPerformerScraper) scrapeByName(name string) ([]*models.ScrapedPerformer, error) {
	if c.config.PerformerByName != nil {
		s := c.config.getScraper(*c.config.PerformerByName, c.client, c.txnManager, c.globalConfig)
		return s.scrapePerformersByName(name)
	}

	return nil, nil
}

func (c *configPerformerScraper) scrapeByFragment(scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	if c.config.PerformerByFragment != nil {
		s := c.config.getScraper(*c.config.PerformerByFragment, c.client, c.txnManager, c.globalConfig)
		return s.scrapePerformerByFragment(scrapedPerformer)
	}

	// try to match against URL if present
	if scrapedPerformer.URL != nil && *scrapedPerformer.URL != "" {
		return c.scrapeByURL(*scrapedPerformer.URL)
	}

	return nil, nil
}

func (c *configPerformerScraper) scrapeByURL(url string) (*models.ScrapedPerformer, error) {
	for _, scraper := range c.config.PerformerByURL {
		if scraper.matchesURL(url) {
			s := c.config.getScraper(scraper.scraperTypeConfig, c.client, c.txnManager, c.globalConfig)
			ret, err := s.scrapePerformerByURL(url)
			if err != nil {
				return nil, err
			}

			if ret != nil {
				return ret, nil
			}
		}
	}

	return nil, nil
}

type configGalleryScraper struct {
	*configScraper
}

func (c *configGalleryScraper) matchesURL(url string) bool {
	return c.config.matchesGalleryURL(url)
}

func (c *configGalleryScraper) scrapeByGallery(gallery *models.Gallery) (*models.ScrapedGallery, error) {
	if c.config.GalleryByFragment != nil {
		s := c.config.getScraper(*c.config.GalleryByFragment, c.client, c.txnManager, c.globalConfig)
		return s.scrapeGalleryByGallery(gallery)
	}

	return nil, nil
}

func (c *configGalleryScraper) scrapeByFragment(gallery models.ScrapedGalleryInput) (*models.ScrapedGallery, error) {
	if c.config.GalleryByFragment != nil {
		// TODO - this should be galleryByQueryFragment
		s := c.config.getScraper(*c.config.GalleryByFragment, c.client, c.txnManager, c.globalConfig)
		return s.scrapeGalleryByFragment(gallery)
	}

	return nil, nil
}

func (c *configGalleryScraper) scrapeByURL(url string) (*models.ScrapedGallery, error) {
	for _, scraper := range c.config.GalleryByURL {
		if scraper.matchesURL(url) {
			s := c.config.getScraper(scraper.scraperTypeConfig, c.client, c.txnManager, c.globalConfig)
			ret, err := s.scrapeGalleryByURL(url)
			if err != nil {
				return nil, err
			}

			if ret != nil {
				return ret, nil
			}
		}
	}

	return nil, nil
}

type configMovieScraper struct {
	*configScraper
}

func (c *configMovieScraper) matchesURL(url string) bool {
	return c.config.matchesMovieURL(url)
}

func (c *configMovieScraper) scrapeByURL(url string) (*models.ScrapedMovie, error) {
	for _, scraper := range c.config.MovieByURL {
		if scraper.matchesURL(url) {
			s := c.config.getScraper(scraper.scraperTypeConfig, c.client, c.txnManager, c.globalConfig)
			ret, err := s.scrapeMovieByURL(url)
			if err != nil {
				return nil, err
			}

			if ret != nil {
				return ret, nil
			}
		}
	}

	return nil, nil
}

type configScraper struct {
	config       config
	client       *http.Client
	txnManager   models.TransactionManager
	globalConfig GlobalConfig
}

func createScraperFromConfig(c config, client *http.Client, txnManager models.TransactionManager, globalConfig GlobalConfig) scraper {
	base := configScraper{
		client:       client,
		config:       c,
		txnManager:   txnManager,
		globalConfig: globalConfig,
	}

	ret := scraper{
		ID:   c.ID,
		Spec: configScraperSpec(c),
	}

	// only set fields if supported
	if c.supportsPerformers() {
		ret.Performer = &configPerformerScraper{&base}
	}
	if c.supportsGalleries() {
		ret.Gallery = &configGalleryScraper{&base}
	}
	if c.supportsMovies() {
		ret.Movie = &configMovieScraper{&base}
	}
	if c.supportsScenes() {
		ret.Scene = &configSceneScraper{&base}
	}

	return ret
}

func configScraperSpec(c config) *models.Scraper {
	ret := models.Scraper{
		ID:   c.ID,
		Name: c.Name,
	}

	performer := models.ScraperSpec{}
	if c.PerformerByName != nil {
		performer.SupportedScrapes = append(performer.SupportedScrapes, models.ScrapeTypeName)
	}
	if c.PerformerByFragment != nil {
		performer.SupportedScrapes = append(performer.SupportedScrapes, models.ScrapeTypeFragment)
	}
	if len(c.PerformerByURL) > 0 {
		performer.SupportedScrapes = append(performer.SupportedScrapes, models.ScrapeTypeURL)
		for _, v := range c.PerformerByURL {
			performer.Urls = append(performer.Urls, v.URL...)
		}
	}

	if len(performer.SupportedScrapes) > 0 {
		ret.Performer = &performer
	}

	scene := models.ScraperSpec{}
	if c.SceneByFragment != nil {
		scene.SupportedScrapes = append(scene.SupportedScrapes, models.ScrapeTypeFragment)
	}
	if c.SceneByName != nil && c.SceneByQueryFragment != nil {
		scene.SupportedScrapes = append(scene.SupportedScrapes, models.ScrapeTypeName)
	}
	if len(c.SceneByURL) > 0 {
		scene.SupportedScrapes = append(scene.SupportedScrapes, models.ScrapeTypeURL)
		for _, v := range c.SceneByURL {
			scene.Urls = append(scene.Urls, v.URL...)
		}
	}

	if len(scene.SupportedScrapes) > 0 {
		ret.Scene = &scene
	}

	gallery := models.ScraperSpec{}
	if c.GalleryByFragment != nil {
		gallery.SupportedScrapes = append(gallery.SupportedScrapes, models.ScrapeTypeFragment)
	}
	if len(c.GalleryByURL) > 0 {
		gallery.SupportedScrapes = append(gallery.SupportedScrapes, models.ScrapeTypeURL)
		for _, v := range c.GalleryByURL {
			gallery.Urls = append(gallery.Urls, v.URL...)
		}
	}

	if len(gallery.SupportedScrapes) > 0 {
		ret.Gallery = &gallery
	}

	movie := models.ScraperSpec{}
	if len(c.MovieByURL) > 0 {
		movie.SupportedScrapes = append(movie.SupportedScrapes, models.ScrapeTypeURL)
		for _, v := range c.MovieByURL {
			movie.Urls = append(movie.Urls, v.URL...)
		}
	}

	if len(movie.SupportedScrapes) > 0 {
		ret.Movie = &movie
	}

	return &ret
}
