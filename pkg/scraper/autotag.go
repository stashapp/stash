package scraper

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
)

// autoTagScraperID is the scraper ID for the built-in AutoTag scraper
const (
	autoTagScraperID   = "builtin_autotag"
	autoTagScraperName = "Auto Tag"
)

type autotagScraper struct {
	txnManager   models.TransactionManager
	globalConfig GlobalConfig
}

func (s *autotagScraper) matchPerformers(path string, performerReader models.PerformerReader) ([]*models.ScrapedPerformer, error) {
	p, err := match.PathToPerformers(path, performerReader)
	if err != nil {
		return nil, fmt.Errorf("error matching performers: %w", err)
	}

	var ret []*models.ScrapedPerformer
	for _, pp := range p {
		id := strconv.Itoa(pp.ID)

		sp := &models.ScrapedPerformer{
			Name:     &pp.Name.String,
			StoredID: &id,
		}
		if pp.Gender.Valid {
			sp.Gender = &pp.Gender.String
		}

		ret = append(ret, sp)
	}

	return ret, nil
}

func (s *autotagScraper) matchStudio(path string, studioReader models.StudioReader) (*models.ScrapedStudio, error) {
	st, err := match.PathToStudios(path, studioReader)
	if err != nil {
		return nil, fmt.Errorf("error matching studios: %w", err)
	}

	if len(st) > 0 {
		id := strconv.Itoa(st[0].ID)
		return &models.ScrapedStudio{
			Name:     st[0].Name.String,
			StoredID: &id,
		}, nil
	}

	return nil, nil
}

func (s *autotagScraper) matchTags(path string, tagReader models.TagReader) ([]*models.ScrapedTag, error) {
	t, err := match.PathToTags(path, tagReader)
	if err != nil {
		return nil, fmt.Errorf("error matching tags: %w", err)
	}

	var ret []*models.ScrapedTag
	for _, tt := range t {
		id := strconv.Itoa(tt.ID)

		st := &models.ScrapedTag{
			Name:     tt.Name,
			StoredID: &id,
		}

		ret = append(ret, st)
	}

	return ret, nil
}

type autotagSceneScraper struct {
	*autotagScraper
}

func (c *autotagSceneScraper) scrapeByName(name string) ([]*models.ScrapedScene, error) {
	return nil, ErrNotSupported
}

func (c *autotagSceneScraper) scrapeByScene(scene *models.Scene) (*models.ScrapedScene, error) {
	var ret *models.ScrapedScene

	// populate performers, studio and tags based on scene path
	if err := c.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		path := scene.Path
		performers, err := c.matchPerformers(path, r.Performer())
		if err != nil {
			return err
		}
		studio, err := c.matchStudio(path, r.Studio())
		if err != nil {
			return err
		}

		tags, err := c.matchTags(path, r.Tag())
		if err != nil {
			return err
		}

		if len(performers) > 0 || studio != nil || len(tags) > 0 {
			ret = &models.ScrapedScene{
				Performers: performers,
				Studio:     studio,
				Tags:       tags,
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *autotagSceneScraper) scrapeByFragment(scene models.ScrapedSceneInput) (*models.ScrapedScene, error) {
	return nil, ErrNotSupported
}

func (c *autotagSceneScraper) scrapeByURL(url string) (*models.ScrapedScene, error) {
	return nil, ErrNotSupported
}

type autotagGalleryScraper struct {
	*autotagScraper
}

func (c *autotagGalleryScraper) scrapeByGallery(gallery *models.Gallery) (*models.ScrapedGallery, error) {
	if !gallery.Path.Valid {
		// not valid for non-path-based galleries
		return nil, nil
	}

	var ret *models.ScrapedGallery

	// populate performers, studio and tags based on scene path
	if err := c.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		path := gallery.Path.String
		performers, err := c.matchPerformers(path, r.Performer())
		if err != nil {
			return err
		}
		studio, err := c.matchStudio(path, r.Studio())
		if err != nil {
			return err
		}

		tags, err := c.matchTags(path, r.Tag())
		if err != nil {
			return err
		}

		if len(performers) > 0 || studio != nil || len(tags) > 0 {
			ret = &models.ScrapedGallery{
				Performers: performers,
				Studio:     studio,
				Tags:       tags,
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *autotagGalleryScraper) scrapeByFragment(gallery models.ScrapedGalleryInput) (*models.ScrapedGallery, error) {
	return nil, ErrNotSupported
}

func (c *autotagGalleryScraper) scrapeByURL(url string) (*models.ScrapedGallery, error) {
	return nil, ErrNotSupported
}

func getAutoTagScraper(txnManager models.TransactionManager, globalConfig GlobalConfig) scraper {
	base := autotagScraper{
		txnManager:   txnManager,
		globalConfig: globalConfig,
	}

	supportedScrapes := []models.ScrapeType{
		models.ScrapeTypeFragment,
	}

	return scraper_s{
		Spec: &models.Scraper{
			ID:   autoTagScraperID,
			Name: autoTagScraperName,
			Scene: &models.ScraperSpec{
				SupportedScrapes: supportedScrapes,
			},
			Gallery: &models.ScraperSpec{
				SupportedScrapes: supportedScrapes,
			},
		},
		scene:   &autotagSceneScraper{&base},
		gallery: &autotagGalleryScraper{&base},
	}
}
