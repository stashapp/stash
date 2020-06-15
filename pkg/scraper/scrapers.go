package scraper

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

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
	scraperFiles := []string{}
	err := filepath.Walk(path, func(fp string, f os.FileInfo, err error) error {
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
	scrapers = append(scrapers, GetFreeonesScraper())

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

func ReloadScrapers() error {
	scrapers = nil
	_, err := loadScrapers()
	return err
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

func ListSceneScrapers() ([]*models.Scraper, error) {
	// read scraper config files from the directory and cache
	scrapers, err := loadScrapers()

	if err != nil {
		return nil, err
	}

	var ret []*models.Scraper
	for _, s := range scrapers {
		// filter on type
		if s.supportsScenes() {
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
		ret, err := s.ScrapePerformer(scrapedPerformer)
		if err != nil {
			return nil, err
		}

		// post-process - set the image if applicable
		if err := setPerformerImage(ret); err != nil {
			logger.Warnf("Could not set image using URL %s: %s", *ret.Image, err.Error())
		}

		return ret, nil
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}

func ScrapePerformerURL(url string) (*models.ScrapedPerformer, error) {
	for _, s := range scrapers {
		if s.matchesPerformerURL(url) {
			ret, err := s.ScrapePerformerURL(url)
			if err != nil {
				return nil, err
			}

			// post-process - set the image if applicable
			if err := setPerformerImage(ret); err != nil {
				logger.Warnf("Could not set image using URL %s: %s", *ret.Image, err.Error())
			}

			return ret, nil
		}
	}

	return nil, nil
}

func matchPerformer(p *models.ScrapedScenePerformer) error {
	qb := models.NewPerformerQueryBuilder()

	performers, err := qb.FindByNames([]string{p.Name}, nil, true)

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

	studio, err := qb.FindByName(s.Name, nil, true)

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
func matchMovie(m *models.ScrapedSceneMovie) error {
	qb := models.NewMovieQueryBuilder()

	movies, err := qb.FindByNames([]string{m.Name}, nil, true)

	if err != nil {
		return err
	}

	if len(movies) != 1 {
		// ignore - cannot match
		return nil
	}

	id := strconv.Itoa(movies[0].ID)
	m.ID = &id
	return nil
}

func matchTag(s *models.ScrapedSceneTag) error {
	qb := models.NewTagQueryBuilder()

	tag, err := qb.FindByName(s.Name, nil, true)

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

func postScrapeScene(ret *models.ScrapedScene) error {
	for _, p := range ret.Performers {
		err := matchPerformer(p)
		if err != nil {
			return err
		}
	}

	for _, p := range ret.Movies {
		err := matchMovie(p)
		if err != nil {
			return err
		}
	}

	for _, t := range ret.Tags {
		err := matchTag(t)
		if err != nil {
			return err
		}
	}

	if ret.Studio != nil {
		err := matchStudio(ret.Studio)
		if err != nil {
			return err
		}
	}

	// post-process - set the image if applicable
	if err := setSceneImage(ret); err != nil {
		logger.Warnf("Could not set image using URL %s: %s", *ret.Image, err.Error())
	}

	return nil
}

func ScrapeScene(scraperID string, scene models.SceneUpdateInput) (*models.ScrapedScene, error) {
	// find scraper with the provided id
	s := findScraper(scraperID)
	if s != nil {
		ret, err := s.ScrapeScene(scene)

		if err != nil {
			return nil, err
		}

		if ret != nil {
			err = postScrapeScene(ret)
			if err != nil {
				return nil, err
			}
		}

		return ret, nil
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}

func ScrapeSceneURL(url string) (*models.ScrapedScene, error) {
	for _, s := range scrapers {
		if s.matchesSceneURL(url) {
			ret, err := s.ScrapeSceneURL(url)

			if err != nil {
				return nil, err
			}

			err = postScrapeScene(ret)
			if err != nil {
				return nil, err
			}

			return ret, nil
		}
	}

	return nil, nil
}
