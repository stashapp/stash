package scraper

import (
	"errors"
	"path/filepath"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
)

var scrapers []scraperConfig

func loadScrapers() ([]scraperConfig, error) {
	if scrapers != nil {
		return scrapers, nil
	}

	path := config.GetScrapersPath()
	scrapers = make([]scraperConfig, 0)

	logger.Debugf("Reading scraper configs from %s", path)
	scraperFiles, err := filepath.Glob(filepath.Join(path, "*.yml"))

	if err != nil {
		logger.Errorf("Error reading scraper configs: %s", err.Error())
		return nil, err
	}

	// add built-in freeones scraper
	scrapers = append(scrapers, GetFreeonesScraper())

	for _, file := range scraperFiles {
		scraper, err := loadScraperFromYAML(file)
		if err != nil {
			logger.Errorf("Error loading scraper %s: %s", file, err.Error())
		} else {
			scrapers = append(scrapers, *scraper)
		}
	}

	return scrapers, nil
}

func ListPerformerScrapers() ([]*models.Scraper, error) {
	// read scraper config files from the directory and cache
	scrapers, err := loadScrapers()

	if err != nil {
		return nil, err
	}

	var ret []*models.Scraper
	for _, s := range scrapers {
		// filter on type
		if s.supportsPerformers() {
			ret = append(ret, s.toScraper())
		}
	}

	return ret, nil
}

func findPerformerScraper(scraperID string) *scraperConfig {
	// read scraper config files from the directory and cache
	loadScrapers()

	for _, s := range scrapers {
		if s.ID == scraperID {
			return &s
		}
	}

	return nil
}

func ScrapePerformerList(scraperID string, query string) ([]*models.ScrapedPerformer, error) {
	// find scraper with the provided id
	s := findPerformerScraper(scraperID)
	if s != nil {
		return s.ScrapePerformerNames(query)
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}

func ScrapePerformer(scraperID string, scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	// find scraper with the provided id
	s := findPerformerScraper(scraperID)
	if s != nil {
		return s.ScrapePerformer(scrapedPerformer)
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}

func ScrapePerformerURL(url string) (*models.ScrapedPerformer, error) {
	for _, s := range scrapers {
		if s.matchesPerformerURL(url) {
			return s.ScrapePerformerURL(url)
		}
	}

	return nil, nil
}
