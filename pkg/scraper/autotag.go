package scraper

import (
	"context"
	"fmt"
	"net/http"
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

func autotagMatchPerformers(path string, performerReader models.PerformerReader, trimExt bool) ([]*models.ScrapedPerformer, error) {
	p, err := match.PathToPerformers(path, performerReader, nil, trimExt)
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

func autotagMatchStudio(path string, studioReader models.StudioReader, trimExt bool) (*models.ScrapedStudio, error) {
	studio, err := match.PathToStudio(path, studioReader, nil, trimExt)
	if err != nil {
		return nil, fmt.Errorf("error matching studios: %w", err)
	}

	if studio != nil {
		id := strconv.Itoa(studio.ID)
		return &models.ScrapedStudio{
			Name:     studio.Name.String,
			StoredID: &id,
		}, nil
	}

	return nil, nil
}

func autotagMatchTags(path string, tagReader models.TagReader, trimExt bool) ([]*models.ScrapedTag, error) {
	t, err := match.PathToTags(path, tagReader, nil, trimExt)
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

func (s autotagScraper) viaScene(ctx context.Context, _client *http.Client, scene *models.Scene) (*ScrapedScene, error) {
	var ret *ScrapedScene
	const trimExt = false

	// populate performers, studio and tags based on scene path
	if err := s.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		path := scene.Path
		performers, err := autotagMatchPerformers(path, r.Performer(), trimExt)
		if err != nil {
			return fmt.Errorf("autotag scraper viaScene: %w", err)
		}
		studio, err := autotagMatchStudio(path, r.Studio(), trimExt)
		if err != nil {
			return fmt.Errorf("autotag scraper viaScene: %w", err)
		}

		tags, err := autotagMatchTags(path, r.Tag(), trimExt)
		if err != nil {
			return fmt.Errorf("autotag scraper viaScene: %w", err)
		}

		if len(performers) > 0 || studio != nil || len(tags) > 0 {
			ret = &ScrapedScene{
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

func (s autotagScraper) viaGallery(ctx context.Context, _client *http.Client, gallery *models.Gallery) (*ScrapedGallery, error) {
	if !gallery.Path.Valid {
		// not valid for non-path-based galleries
		return nil, nil
	}

	// only trim extension if gallery is file-based
	trimExt := gallery.Zip

	var ret *ScrapedGallery

	// populate performers, studio and tags based on scene path
	if err := s.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		path := gallery.Path.String
		performers, err := autotagMatchPerformers(path, r.Performer(), trimExt)
		if err != nil {
			return fmt.Errorf("autotag scraper viaGallery: %w", err)
		}
		studio, err := autotagMatchStudio(path, r.Studio(), trimExt)
		if err != nil {
			return fmt.Errorf("autotag scraper viaGallery: %w", err)
		}

		tags, err := autotagMatchTags(path, r.Tag(), trimExt)
		if err != nil {
			return fmt.Errorf("autotag scraper viaGallery: %w", err)
		}

		if len(performers) > 0 || studio != nil || len(tags) > 0 {
			ret = &ScrapedGallery{
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

func (s autotagScraper) supports(ty ScrapeContentType) bool {
	switch ty {
	case ScrapeContentTypeScene:
		return true
	case ScrapeContentTypeGallery:
		return true
	}

	return false
}

func (s autotagScraper) supportsURL(url string, ty ScrapeContentType) bool {
	return false
}

func (s autotagScraper) spec() Scraper {
	supportedScrapes := []ScrapeType{
		ScrapeTypeFragment,
	}

	return Scraper{
		ID:   autoTagScraperID,
		Name: autoTagScraperName,
		Scene: &ScraperSpec{
			SupportedScrapes: supportedScrapes,
		},
		Gallery: &ScraperSpec{
			SupportedScrapes: supportedScrapes,
		},
	}
}

func getAutoTagScraper(txnManager models.TransactionManager, globalConfig GlobalConfig) scraper {
	base := autotagScraper{
		txnManager:   txnManager,
		globalConfig: globalConfig,
	}

	return base
}
