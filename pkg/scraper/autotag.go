package scraper

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

// autoTagScraperID is the scraper ID for the built-in AutoTag scraper
const (
	autoTagScraperID   = "builtin_autotag"
	autoTagScraperName = "Auto Tag"
)

type autotagScraper struct {
	txnManager      txn.Manager
	performerReader models.PerformerAutoTagQueryer
	studioReader    models.StudioAutoTagQueryer
	tagReader       models.TagAutoTagQueryer

	globalConfig GlobalConfig
}

func autotagMatchPerformers(ctx context.Context, path string, performerReader models.PerformerAutoTagQueryer, trimExt bool) ([]*models.ScrapedPerformer, error) {
	p, err := match.PathToPerformers(ctx, path, performerReader, nil, trimExt)
	if err != nil {
		return nil, fmt.Errorf("error matching performers: %w", err)
	}

	var ret []*models.ScrapedPerformer
	for _, pp := range p {
		id := strconv.Itoa(pp.ID)

		sp := &models.ScrapedPerformer{
			Name:     &pp.Name,
			StoredID: &id,
		}
		if pp.Gender != nil && pp.Gender.IsValid() {
			v := pp.Gender.String()
			sp.Gender = &v
		}

		ret = append(ret, sp)
	}

	return ret, nil
}

func autotagMatchStudio(ctx context.Context, path string, studioReader models.StudioAutoTagQueryer, trimExt bool) (*models.ScrapedStudio, error) {
	studio, err := match.PathToStudio(ctx, path, studioReader, nil, trimExt)
	if err != nil {
		return nil, fmt.Errorf("error matching studios: %w", err)
	}

	if studio != nil {
		id := strconv.Itoa(studio.ID)
		return &models.ScrapedStudio{
			Name:     studio.Name,
			StoredID: &id,
		}, nil
	}

	return nil, nil
}

func autotagMatchTags(ctx context.Context, path string, tagReader models.TagAutoTagQueryer, trimExt bool) ([]*models.ScrapedTag, error) {
	t, err := match.PathToTags(ctx, path, tagReader, nil, trimExt)
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

func (s autotagScraper) viaScene(ctx context.Context, _client *http.Client, scene *models.Scene) (*models.ScrapedScene, error) {
	var ret *models.ScrapedScene
	const trimExt = false

	// populate performers, studio and tags based on scene path
	if err := txn.WithReadTxn(ctx, s.txnManager, func(ctx context.Context) error {
		path := scene.Path
		if path == "" {
			return nil
		}

		performers, err := autotagMatchPerformers(ctx, path, s.performerReader, trimExt)
		if err != nil {
			return fmt.Errorf("autotag scraper viaScene: %w", err)
		}
		studio, err := autotagMatchStudio(ctx, path, s.studioReader, trimExt)
		if err != nil {
			return fmt.Errorf("autotag scraper viaScene: %w", err)
		}

		tags, err := autotagMatchTags(ctx, path, s.tagReader, trimExt)
		if err != nil {
			return fmt.Errorf("autotag scraper viaScene: %w", err)
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

func (s autotagScraper) viaGallery(ctx context.Context, _client *http.Client, gallery *models.Gallery) (*models.ScrapedGallery, error) {
	path := gallery.Path
	if path == "" {
		// not valid for non-path-based galleries
		return nil, nil
	}

	// only trim extension if gallery is file-based
	trimExt := gallery.PrimaryFileID != nil

	var ret *models.ScrapedGallery

	// populate performers, studio and tags based on scene path
	if err := txn.WithReadTxn(ctx, s.txnManager, func(ctx context.Context) error {
		path := gallery.Path
		performers, err := autotagMatchPerformers(ctx, path, s.performerReader, trimExt)
		if err != nil {
			return fmt.Errorf("autotag scraper viaGallery: %w", err)
		}
		studio, err := autotagMatchStudio(ctx, path, s.studioReader, trimExt)
		if err != nil {
			return fmt.Errorf("autotag scraper viaGallery: %w", err)
		}

		tags, err := autotagMatchTags(ctx, path, s.tagReader, trimExt)
		if err != nil {
			return fmt.Errorf("autotag scraper viaGallery: %w", err)
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

func getAutoTagScraper(repo Repository, globalConfig GlobalConfig) scraper {
	base := autotagScraper{
		txnManager:      repo.TxnManager,
		performerReader: repo.PerformerFinder,
		studioReader:    repo.StudioFinder,
		tagReader:       repo.TagFinder,
		globalConfig:    globalConfig,
	}

	return base
}
