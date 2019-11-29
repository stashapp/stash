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

func findScraper(scraperID string) *scraperConfig {
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
	s := findScraper(scraperID)
	if s != nil {
		return s.ScrapePerformerNames(query)
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}

func ScrapePerformer(scraperID string, scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	// find scraper with the provided id
	s := findScraper(scraperID)
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

func ScrapeScene(scraperID string, scene models.SceneUpdateInput) (*models.ScrapedScene, error) {
	// find scraper with the provided id
	s := findScraper(scraperID)
	if s != nil {
		return s.ScrapeScene(scene)
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}

func matchPerformer(p *models.ScrapedScenePerformer) error {
	qb := models.NewPerformerQueryBuilder()

	performers, err := qb.FindByNames([]string{p.Name}, nil)

	if err != nil {
		return err
	}

	if len(performers) != 1 {
		// ignore - cannot match
		return nil
	}

	id := strconv.Itoa(performers[0].ID)
	p.ID = &id
	return nil
}

func matchStudio(s *models.ScrapedSceneStudio) error {
	qb := models.NewStudioQueryBuilder()

	studio, err := qb.FindByName(s.Name, nil)

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

func matchTag(s *models.ScrapedSceneTag) error {
	qb := models.NewTagQueryBuilder()

	tag, err := qb.FindByName(s.Name, nil)

	if err != nil {
		return err
	}

	if tag == nil {
		// ignore - cannot match
		return nil
	}

	id := strconv.Itoa(tag.ID)
	s.ID = &id
	return nil
}

func ScrapeSceneURL(url string) (*models.ScrapedScene, error) {
	// find scraper that matches the url given
	s := findScraperURL(url, models.ScraperTypeScene)
	if s != nil {
		ret, err := s.ScrapeSceneURL(url)

		if err != nil {
			return nil, err
		}

		for _, p := range ret.Performers {
			err = matchPerformer(p)
			if err != nil {
				return nil, err
			}
		}

		for _, t := range ret.Tags {
			err = matchTag(t)
			if err != nil {
				return nil, err
			}
		}

		if ret.Studio != nil {
			err = matchStudio(ret.Studio)
			if err != nil {
				return nil, err
			}
		}

		return ret, nil
	}

	return nil, nil
}
