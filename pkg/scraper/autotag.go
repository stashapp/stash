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
	repository   models.Repository
	globalConfig GlobalConfig
}

func autotagMatchPerformers(ctx context.Context, path string, performerReader models.PerformerReader) ([]*models.ScrapedPerformer, error) {
	p, err := match.PathToPerformers(ctx, path, performerReader, nil)
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

func autotagMatchStudio(ctx context.Context, path string, studioReader models.StudioReader) (*models.ScrapedStudio, error) {
	studio, err := match.PathToStudio(ctx, path, studioReader, nil)
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

func autotagMatchTags(ctx context.Context, path string, tagReader models.TagReader) ([]*models.ScrapedTag, error) {
	t, err := match.PathToTags(ctx, path, tagReader, nil)
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

	r := s.repository

	// populate performers, studio and tags based on scene path
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		path := scene.Path
		performers, err := autotagMatchPerformers(ctx, path, r.Performer)
		if err != nil {
			return fmt.Errorf("autotag scraper viaScene: %w", err)
		}
		studio, err := autotagMatchStudio(ctx, path, r.Studio)
		if err != nil {
			return fmt.Errorf("autotag scraper viaScene: %w", err)
		}

		tags, err := autotagMatchTags(ctx, path, r.Tag)
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

	var ret *ScrapedGallery

	// populate performers, studio and tags based on scene path
	r := s.repository
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		path := gallery.Path.String
		performers, err := autotagMatchPerformers(ctx, path, r.Performer)
		if err != nil {
			return fmt.Errorf("autotag scraper viaGallery: %w", err)
		}
		studio, err := autotagMatchStudio(ctx, path, r.Studio)
		if err != nil {
			return fmt.Errorf("autotag scraper viaGallery: %w", err)
		}

		tags, err := autotagMatchTags(ctx, path, r.Tag)
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

func getAutoTagScraper(txnManager models.Repository, globalConfig GlobalConfig) scraper {
	base := autotagScraper{
		repository:   txnManager,
		globalConfig: globalConfig,
	}

	return base
}
