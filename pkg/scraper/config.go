package scraper

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type config struct {
	ID   string
	path string

	// The name of the scraper. This is displayed in the UI.
	Name string `yaml:"name"`

	// Configuration for querying performers by name
	PerformerByName *scraperTypeConfig `yaml:"performerByName"`

	// Configuration for querying performers by a Performer fragment
	PerformerByFragment *scraperTypeConfig `yaml:"performerByFragment"`

	// Configuration for querying a performer by a URL
	PerformerByURL []*scrapeByURLConfig `yaml:"performerByURL"`

	// Configuration for querying scenes by a Scene fragment
	SceneByFragment *scraperTypeConfig `yaml:"sceneByFragment"`

	// Configuration for querying gallery by a Gallery fragment
	GalleryByFragment *scraperTypeConfig `yaml:"galleryByFragment"`

	// Configuration for querying scenes by name
	SceneByName *scraperTypeConfig `yaml:"sceneByName"`

	// Configuration for querying scenes by query fragment
	SceneByQueryFragment *scraperTypeConfig `yaml:"sceneByQueryFragment"`

	// Configuration for querying a scene by a URL
	SceneByURL []*scrapeByURLConfig `yaml:"sceneByURL"`

	// Configuration for querying a gallery by a URL
	GalleryByURL []*scrapeByURLConfig `yaml:"galleryByURL"`

	// Configuration for querying a movie by a URL
	MovieByURL []*scrapeByURLConfig `yaml:"movieByURL"`

	// Scraper debugging options
	DebugOptions *scraperDebugOptions `yaml:"debug"`

	// Stash server configuration
	StashServer *stashServer `yaml:"stashServer"`

	// Xpath scraping configurations
	XPathScrapers mappedScrapers `yaml:"xPathScrapers"`

	// Json scraping configurations
	JsonScrapers mappedScrapers `yaml:"jsonScrapers"`

	// Scraping driver options
	DriverOptions *scraperDriverOptions `yaml:"driver"`
}

func (c config) validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("name must not be empty")
	}

	if c.PerformerByName != nil {
		if err := c.PerformerByName.validate(); err != nil {
			return err
		}
	}

	if c.PerformerByFragment != nil {
		if err := c.PerformerByFragment.validate(); err != nil {
			return err
		}
	}

	if c.SceneByFragment != nil {
		if err := c.SceneByFragment.validate(); err != nil {
			return err
		}
	}

	for _, s := range c.PerformerByURL {
		if err := s.validate(); err != nil {
			return err
		}
	}

	for _, s := range c.SceneByURL {
		if err := s.validate(); err != nil {
			return err
		}
	}

	for _, s := range c.MovieByURL {
		if err := s.validate(); err != nil {
			return err
		}
	}

	return nil
}

type stashServer struct {
	URL string `yaml:"url"`
}

type scraperTypeConfig struct {
	Action  scraperAction `yaml:"action"`
	Script  []string      `yaml:"script,flow"`
	Scraper string        `yaml:"scraper"`

	// for xpath name scraper only
	QueryURL             string               `yaml:"queryURL"`
	QueryURLReplacements queryURLReplacements `yaml:"queryURLReplace"`
}

func (c scraperTypeConfig) validate() error {
	if !c.Action.IsValid() {
		return fmt.Errorf("%s is not a valid scraper action", c.Action)
	}

	if c.Action == scraperActionScript && len(c.Script) == 0 {
		return errors.New("script is mandatory for script scraper action")
	}

	return nil
}

type scrapeByURLConfig struct {
	scraperTypeConfig `yaml:",inline"`
	URL               []string `yaml:"url,flow"`
}

func (c scrapeByURLConfig) validate() error {
	if len(c.URL) == 0 {
		return errors.New("url is mandatory for scrape by url scrapers")
	}

	return c.scraperTypeConfig.validate()
}

func (c scrapeByURLConfig) matchesURL(url string) bool {
	for _, thisURL := range c.URL {
		if strings.Contains(url, thisURL) {
			return true
		}
	}

	return false
}

type scraperDebugOptions struct {
	PrintHTML bool `yaml:"printHTML"`
}

type scraperCookies struct {
	Name        string `yaml:"Name"`
	Value       string `yaml:"Value"`
	ValueRandom int    `yaml:"ValueRandom"`
	Domain      string `yaml:"Domain"`
	Path        string `yaml:"Path"`
}

type cookieOptions struct {
	CookieURL string            `yaml:"CookieURL"`
	Cookies   []*scraperCookies `yaml:"Cookies"`
}

type clickOptions struct {
	XPath string `yaml:"xpath"`
	Sleep int    `yaml:"sleep"`
}

type header struct {
	Key   string `yaml:"Key"`
	Value string `yaml:"Value"`
}

type scraperDriverOptions struct {
	UseCDP  bool             `yaml:"useCDP"`
	Sleep   int              `yaml:"sleep"`
	Clicks  []*clickOptions  `yaml:"clicks"`
	Cookies []*cookieOptions `yaml:"cookies"`
	Headers []*header        `yaml:"headers"`
}

func loadConfigFromYAML(id string, reader io.Reader) (*config, error) {
	ret := &config{}

	parser := yaml.NewDecoder(reader)
	parser.SetStrict(true)
	err := parser.Decode(&ret)
	if err != nil {
		return nil, err
	}

	ret.ID = id

	if err := ret.validate(); err != nil {
		return nil, err
	}

	return ret, nil
}

func loadConfigFromYAMLFile(path string) (*config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// set id to the filename
	id := filepath.Base(path)
	id = id[:strings.LastIndex(id, ".")]

	ret, err := loadConfigFromYAML(id, file)
	if err != nil {
		return nil, err
	}

	ret.path = path

	return ret, nil
}

func (c config) spec() Scraper {
	ret := Scraper{
		ID:   c.ID,
		Name: c.Name,
	}

	performer := ScraperSpec{}
	if c.PerformerByName != nil {
		performer.SupportedScrapes = append(performer.SupportedScrapes, ScrapeTypeName)
	}
	if c.PerformerByFragment != nil {
		performer.SupportedScrapes = append(performer.SupportedScrapes, ScrapeTypeFragment)
	}
	if len(c.PerformerByURL) > 0 {
		performer.SupportedScrapes = append(performer.SupportedScrapes, ScrapeTypeURL)
		for _, v := range c.PerformerByURL {
			performer.Urls = append(performer.Urls, v.URL...)
		}
	}

	if len(performer.SupportedScrapes) > 0 {
		ret.Performer = &performer
	}

	scene := ScraperSpec{}
	if c.SceneByFragment != nil {
		scene.SupportedScrapes = append(scene.SupportedScrapes, ScrapeTypeFragment)
	}
	if c.SceneByName != nil && c.SceneByQueryFragment != nil {
		scene.SupportedScrapes = append(scene.SupportedScrapes, ScrapeTypeName)
	}
	if len(c.SceneByURL) > 0 {
		scene.SupportedScrapes = append(scene.SupportedScrapes, ScrapeTypeURL)
		for _, v := range c.SceneByURL {
			scene.Urls = append(scene.Urls, v.URL...)
		}
	}

	if len(scene.SupportedScrapes) > 0 {
		ret.Scene = &scene
	}

	gallery := ScraperSpec{}
	if c.GalleryByFragment != nil {
		gallery.SupportedScrapes = append(gallery.SupportedScrapes, ScrapeTypeFragment)
	}
	if len(c.GalleryByURL) > 0 {
		gallery.SupportedScrapes = append(gallery.SupportedScrapes, ScrapeTypeURL)
		for _, v := range c.GalleryByURL {
			gallery.Urls = append(gallery.Urls, v.URL...)
		}
	}

	if len(gallery.SupportedScrapes) > 0 {
		ret.Gallery = &gallery
	}

	movie := ScraperSpec{}
	if len(c.MovieByURL) > 0 {
		movie.SupportedScrapes = append(movie.SupportedScrapes, ScrapeTypeURL)
		for _, v := range c.MovieByURL {
			movie.Urls = append(movie.Urls, v.URL...)
		}
	}

	if len(movie.SupportedScrapes) > 0 {
		ret.Movie = &movie
		ret.Group = &movie
	}

	return ret
}

func (c config) supports(ty ScrapeContentType) bool {
	switch ty {
	case ScrapeContentTypePerformer:
		return c.PerformerByName != nil || c.PerformerByFragment != nil || len(c.PerformerByURL) > 0
	case ScrapeContentTypeScene:
		return (c.SceneByName != nil && c.SceneByQueryFragment != nil) || c.SceneByFragment != nil || len(c.SceneByURL) > 0
	case ScrapeContentTypeGallery:
		return c.GalleryByFragment != nil || len(c.GalleryByURL) > 0
	case ScrapeContentTypeMovie, ScrapeContentTypeGroup:
		return len(c.MovieByURL) > 0
	}

	panic("Unhandled ScrapeContentType")
}

func (c config) matchesURL(url string, ty ScrapeContentType) bool {
	switch ty {
	case ScrapeContentTypePerformer:
		for _, scraper := range c.PerformerByURL {
			if scraper.matchesURL(url) {
				return true
			}
		}
	case ScrapeContentTypeScene:
		for _, scraper := range c.SceneByURL {
			if scraper.matchesURL(url) {
				return true
			}
		}
	case ScrapeContentTypeGallery:
		for _, scraper := range c.GalleryByURL {
			if scraper.matchesURL(url) {
				return true
			}
		}
	case ScrapeContentTypeMovie, ScrapeContentTypeGroup:
		for _, scraper := range c.MovieByURL {
			if scraper.matchesURL(url) {
				return true
			}
		}
	}

	return false
}
