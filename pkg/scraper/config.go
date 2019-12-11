package scraper

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/stashapp/stash/pkg/models"
)

type scraperAction string

const (
	scraperActionScript scraperAction = "script"
)

var allScraperAction = []scraperAction{
	scraperActionScript,
}

func (e scraperAction) IsValid() bool {
	switch e {
	case scraperActionScript:
		return true
	}
	return false
}

type scraperTypeConfig struct {
	Action scraperAction `yaml:"action"`
	Script []string      `yaml:"script,flow"`
}

type scrapePerformerNamesFunc func(c scraperTypeConfig, name string) ([]*models.ScrapedPerformer, error)

type performerByNameConfig struct {
	scraperTypeConfig `yaml:",inline"`
	performScrape     scrapePerformerNamesFunc
}

func (c *performerByNameConfig) resolveFn() {
	if c.Action == scraperActionScript {
		c.performScrape = scrapePerformerNamesScript
	}
}

type scrapePerformerFragmentFunc func(c scraperTypeConfig, scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error)

type performerByFragmentConfig struct {
	scraperTypeConfig `yaml:",inline"`
	performScrape     scrapePerformerFragmentFunc
}

func (c *performerByFragmentConfig) resolveFn() {
	if c.Action == scraperActionScript {
		c.performScrape = scrapePerformerFragmentScript
	}
}

type scrapePerformerByURLFunc func(c scraperTypeConfig, url string) (*models.ScrapedPerformer, error)

type scraperByURLConfig struct {
	scraperTypeConfig `yaml:",inline"`
	URL               []string `yaml:"url,flow"`
	performScrape     scrapePerformerByURLFunc
}

func (c *scraperByURLConfig) resolveFn() {
	if c.Action == scraperActionScript {
		c.performScrape = scrapePerformerURLScript
	}
}

func (s scraperByURLConfig) matchesURL(url string) bool {
	for _, thisURL := range s.URL {
		if strings.Contains(url, thisURL) {
			return true
		}
	}

	return false
}

type scraperConfig struct {
	ID                  string
	Name                string                     `yaml:"name"`
	PerformerByName     *performerByNameConfig     `yaml:"performerByName"`
	PerformerByFragment *performerByFragmentConfig `yaml:"performerByFragment"`
	PerformerByURL      []*scraperByURLConfig      `yaml:"performerByURL"`
}

func loadScraperFromYAML(path string) (*scraperConfig, error) {
	ret := &scraperConfig{}

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	parser := yaml.NewDecoder(file)
	parser.SetStrict(true)
	err = parser.Decode(&ret)
	if err != nil {
		return nil, err
	}

	// set id to the filename
	id := filepath.Base(path)
	id = id[:strings.LastIndex(id, ".")]
	ret.ID = id

	// set the scraper interface
	ret.initialiseConfigs()

	return ret, nil
}

func (c *scraperConfig) initialiseConfigs() {
	if c.PerformerByName != nil {
		c.PerformerByName.resolveFn()
	}
	if c.PerformerByFragment != nil {
		c.PerformerByFragment.resolveFn()
	}
	for _, s := range c.PerformerByURL {
		s.resolveFn()
	}
}

func (c scraperConfig) toScraper() *models.Scraper {
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

	return &ret
}

func (c scraperConfig) supportsPerformers() bool {
	return c.PerformerByName != nil || c.PerformerByFragment != nil || len(c.PerformerByURL) > 0
}

func (c scraperConfig) matchesPerformerURL(url string) bool {
	for _, scraper := range c.PerformerByURL {
		if scraper.matchesURL(url) {
			return true
		}
	}

	return false
}

func (c scraperConfig) ScrapePerformerNames(name string) ([]*models.ScrapedPerformer, error) {
	if c.PerformerByName != nil && c.PerformerByName.performScrape != nil {
		return c.PerformerByName.performScrape(c.PerformerByName.scraperTypeConfig, name)
	}

	return nil, nil
}

func (c scraperConfig) ScrapePerformer(scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	if c.PerformerByFragment != nil && c.PerformerByFragment.performScrape != nil {
		return c.PerformerByFragment.performScrape(c.PerformerByFragment.scraperTypeConfig, scrapedPerformer)
	}

	return nil, nil
}

func (c scraperConfig) ScrapePerformerURL(url string) (*models.ScrapedPerformer, error) {
	for _, scraper := range c.PerformerByURL {
		if scraper.matchesURL(url) && scraper.performScrape != nil {
			ret, err := scraper.performScrape(scraper.scraperTypeConfig, url)
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
