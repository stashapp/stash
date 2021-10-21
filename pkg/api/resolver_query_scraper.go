package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/scraper/stashbox"
)

var ErrInput = errors.New("input error")

func (r *queryResolver) ScrapeURL(ctx context.Context, url string, ty models.ScrapeContentType) (models.ScrapedContent, error) {
	return r.scraperCache.ScrapeURL(ctx, url, ty)
}

// deprecated
func (r *queryResolver) ScrapeFreeones(ctx context.Context, performer_name string) (*models.ScrapedPerformer, error) {
	scrapedPerformer := models.ScrapedPerformerInput{
		Name: &performer_name,
	}

	content, err := r.scraperCache.ScrapeFragment(ctx, scraper.FreeonesScraperID, scraper.Input{Performer: &scrapedPerformer})
	if err != nil {
		return nil, err
	}

	return marshalScrapedPerformer(content)
}

// deprecated
func (r *queryResolver) ScrapeFreeonesPerformerList(ctx context.Context, query string) ([]string, error) {
	scrapedPerformers, err := r.scraperCache.ScrapeByName(scraper.FreeonesScraperID, query, models.ScrapeContentTypePerformer)

	if err != nil {
		return nil, err
	}

	var ret []string
	for _, v := range scrapedPerformers {
		if p, ok := v.(models.ScrapedPerformer); ok {
			if p.Name != nil {
				ret = append(ret, *p.Name)
			}
		} else {
			logger.Errorf("Internal Server Error: could not convert scraped content into a performer")
		}
	}

	return ret, nil
}

func (r *queryResolver) ListScrapers(ctx context.Context, ty models.ScrapeContentType) ([]*models.Scraper, error) {
	return r.scraperCache.ListScrapers(ty), nil
}

func (r *queryResolver) ListPerformerScrapers(ctx context.Context) ([]*models.Scraper, error) {
	return r.scraperCache.ListScrapers(models.ScrapeContentTypePerformer), nil
}

func (r *queryResolver) ListSceneScrapers(ctx context.Context) ([]*models.Scraper, error) {
	return r.scraperCache.ListScrapers(models.ScrapeContentTypeScene), nil
}

func (r *queryResolver) ListGalleryScrapers(ctx context.Context) ([]*models.Scraper, error) {
	return r.scraperCache.ListScrapers(models.ScrapeContentTypeGallery), nil
}

func (r *queryResolver) ListMovieScrapers(ctx context.Context) ([]*models.Scraper, error) {
	return r.scraperCache.ListScrapers(models.ScrapeContentTypeMovie), nil
}

func (r *queryResolver) ScrapePerformerList(ctx context.Context, scraperID string, query string) ([]*models.ScrapedPerformer, error) {
	if query == "" {
		return nil, nil
	}

	content, err := r.scraperCache.ScrapeByName(scraperID, query, models.ScrapeContentTypePerformer)
	if err != nil {
		return nil, err
	}

	return marshalScrapedPerformers(content)
}

func (r *queryResolver) ScrapePerformer(ctx context.Context, scraperID string, scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	content, err := r.scraperCache.ScrapeFragment(ctx, scraperID, scraper.Input{Performer: &scrapedPerformer})
	if err != nil {
		return nil, err
	}
	return marshalScrapedPerformer(content)
}

func (r *queryResolver) ScrapePerformerURL(ctx context.Context, url string) (*models.ScrapedPerformer, error) {
	ret, err := r.scraperCache.ScrapeURL(ctx, url, models.ScrapeContentTypePerformer)
	if err != nil {
		return nil, err
	}

	if p, ok := ret.(models.ScrapedPerformer); ok {
		return &p, err
	}

	return nil, ErrInternalUnreachable
}

func (r *queryResolver) ScrapeSceneQuery(ctx context.Context, scraperID string, query string) ([]*models.ScrapedScene, error) {
	if query == "" {
		return nil, nil
	}

	content, err := r.scraperCache.ScrapeByName(scraperID, query, models.ScrapeContentTypeScene)
	if err != nil {
		return nil, err
	}

	return marshalScrapedScenes(content)
}

func (r *queryResolver) ScrapeScene(ctx context.Context, scraperID string, scene models.SceneUpdateInput) (*models.ScrapedScene, error) {
	id, err := strconv.Atoi(scene.ID)
	if err != nil {
		return nil, fmt.Errorf("scene ID input %s: err", scene.ID)
	}

	content, err := r.scraperCache.ScrapeID(ctx, scraperID, id, models.ScrapeContentTypeScene)
	if err != nil {
		return nil, err
	}

	return marshalScrapedScene(content)
}

func (r *queryResolver) ScrapeSceneURL(ctx context.Context, url string) (*models.ScrapedScene, error) {
	ret, err := r.scraperCache.ScrapeURL(ctx, url, models.ScrapeContentTypePerformer)
	if err != nil {
		return nil, err
	}
	p, ok := ret.(models.ScrapedScene)
	if ok {
		return &p, err
	}

	return nil, ErrInternalUnreachable
}

func (r *queryResolver) ScrapeGallery(ctx context.Context, scraperID string, gallery models.GalleryUpdateInput) (*models.ScrapedGallery, error) {
	id, err := strconv.Atoi(gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("gallery id input %s: %w", gallery.ID, err)
	}

	content, err := r.scraperCache.ScrapeID(ctx, scraperID, id, models.ScrapeContentTypeGallery)
	if err != nil {
		return nil, err
	}

	return marshalScrapedGallery(content)
}

func (r *queryResolver) ScrapeGalleryURL(ctx context.Context, url string) (*models.ScrapedGallery, error) {
	ret, err := r.scraperCache.ScrapeURL(ctx, url, models.ScrapeContentTypePerformer)
	if err != nil {
		return nil, err
	}
	p, ok := ret.(models.ScrapedGallery)
	if ok {
		return &p, err
	}

	return nil, ErrInternalUnreachable
}

func (r *queryResolver) ScrapeMovieURL(ctx context.Context, url string) (*models.ScrapedMovie, error) {
	ret, err := r.scraperCache.ScrapeURL(ctx, url, models.ScrapeContentTypePerformer)
	if err != nil {
		return nil, err
	}
	p, ok := ret.(models.ScrapedMovie)
	if ok {
		return &p, err
	}

	return nil, ErrInternalUnreachable
}

func (r *queryResolver) QueryStashBoxScene(ctx context.Context, input models.StashBoxSceneQueryInput) ([]*models.ScrapedScene, error) {
	boxes := config.GetInstance().GetStashBoxes()

	if input.StashBoxIndex < 0 || input.StashBoxIndex >= len(boxes) {
		return nil, fmt.Errorf("invalid stash_box_index %d", input.StashBoxIndex)
	}

	client := stashbox.NewClient(*boxes[input.StashBoxIndex], r.txnManager)

	if len(input.SceneIds) > 0 {
		return client.FindStashBoxScenesByFingerprintsFlat(input.SceneIds)
	}

	if input.Q != nil {
		return client.QueryStashBoxScene(ctx, *input.Q)
	}

	return nil, nil
}

func (r *queryResolver) QueryStashBoxPerformer(ctx context.Context, input models.StashBoxPerformerQueryInput) ([]*models.StashBoxPerformerQueryResult, error) {
	boxes := config.GetInstance().GetStashBoxes()

	if input.StashBoxIndex < 0 || input.StashBoxIndex >= len(boxes) {
		return nil, fmt.Errorf("invalid stash_box_index %d", input.StashBoxIndex)
	}

	client := stashbox.NewClient(*boxes[input.StashBoxIndex], r.txnManager)

	if len(input.PerformerIds) > 0 {
		return client.FindStashBoxPerformersByNames(input.PerformerIds)
	}

	if input.Q != nil {
		return client.QueryStashBoxPerformer(*input.Q)
	}

	return nil, nil
}

func (r *queryResolver) getStashBoxClient(index int) (*stashbox.Client, error) {
	boxes := config.GetInstance().GetStashBoxes()

	if index < 0 || index >= len(boxes) {
		return nil, fmt.Errorf("invalid stash_box_index %d", index)
	}

	return stashbox.NewClient(*boxes[index], r.txnManager), nil
}

func (r *queryResolver) ScrapeSingleScene(ctx context.Context, source models.ScraperSourceInput, input models.ScrapeSingleSceneInput) ([]*models.ScrapedScene, error) {
	if source.ScraperID != nil {
		var c models.ScrapedContent
		var content []models.ScrapedContent
		var err error

		switch {
		case input.SceneID != nil:
			var sceneID int
			sceneID, err = strconv.Atoi(*input.SceneID)
			if err != nil {
				return nil, fmt.Errorf("scraper %s: converting input %s: %w", *source.ScraperID, *input.SceneID, err)
			}
			c, err = r.scraperCache.ScrapeID(ctx, *source.ScraperID, sceneID, models.ScrapeContentTypeScene)
			content = []models.ScrapedContent{c}
		case input.SceneInput != nil:
			c, err = r.scraperCache.ScrapeFragment(ctx, *source.ScraperID, scraper.Input{Scene: input.SceneInput})
			content = []models.ScrapedContent{c}
		case input.Query != nil:
			content, err = r.scraperCache.ScrapeByName(*source.ScraperID, *input.Query, models.ScrapeContentTypeScene)
		default:
			err = fmt.Errorf("%w: scene_id, scene_input or query must be set", ErrInput)
		}

		if err != nil {
			return nil, err
		}

		return marshalScrapedScenes(content)
	} else if source.StashBoxIndex != nil {
		client, err := r.getStashBoxClient(*source.StashBoxIndex)
		if err != nil {
			return nil, err
		}

		if input.SceneID != nil {
			return client.FindStashBoxScenesByFingerprintsFlat([]string{*input.SceneID})
		} else if input.Query != nil {
			return client.QueryStashBoxScene(ctx, *input.Query)
		}

		return nil, fmt.Errorf("%w: scene_id or query must be set", ErrInput)
	}

	return nil, fmt.Errorf("%w: scraper_id or stash_box_index must be set", ErrInput)
}

func (r *queryResolver) ScrapeMultiScenes(ctx context.Context, source models.ScraperSourceInput, input models.ScrapeMultiScenesInput) ([][]*models.ScrapedScene, error) {
	if source.ScraperID != nil {
		return nil, ErrNotImplemented
	} else if source.StashBoxIndex != nil {
		client, err := r.getStashBoxClient(*source.StashBoxIndex)
		if err != nil {
			return nil, err
		}

		return client.FindStashBoxScenesByFingerprints(input.SceneIds)
	}

	return nil, errors.New("scraper_id or stash_box_index must be set")
}

func (r *queryResolver) ScrapeSinglePerformer(ctx context.Context, source models.ScraperSourceInput, input models.ScrapeSinglePerformerInput) ([]*models.ScrapedPerformer, error) {
	if source.ScraperID != nil {
		if input.PerformerInput != nil {
			performer, err := r.scraperCache.ScrapeFragment(ctx, *source.ScraperID, scraper.Input{Performer: input.PerformerInput})
			if err != nil {
				return nil, err
			}

			return marshalScrapedPerformers([]models.ScrapedContent{performer})
		}

		if input.Query != nil {
			content, err := r.scraperCache.ScrapeByName(*source.ScraperID, *input.Query, models.ScrapeContentTypePerformer)
			if err != nil {
				return nil, err
			}

			return marshalScrapedPerformers(content)
		}

		return nil, ErrNotImplemented
	} else if source.StashBoxIndex != nil {
		client, err := r.getStashBoxClient(*source.StashBoxIndex)
		if err != nil {
			return nil, err
		}

		var ret []*models.StashBoxPerformerQueryResult
		switch {
		case input.PerformerID != nil:
			ret, err = client.FindStashBoxPerformersByNames([]string{*input.PerformerID})
		case input.Query != nil:
			ret, err = client.QueryStashBoxPerformer(*input.Query)
		default:
			return nil, ErrNotImplemented
		}

		if err != nil {
			return nil, err
		}

		if len(ret) > 0 {
			return ret[0].Results, nil
		}

		return nil, nil
	}

	return nil, errors.New("scraper_id or stash_box_index must be set")
}

func (r *queryResolver) ScrapeMultiPerformers(ctx context.Context, source models.ScraperSourceInput, input models.ScrapeMultiPerformersInput) ([][]*models.ScrapedPerformer, error) {
	if source.ScraperID != nil {
		return nil, ErrNotImplemented
	} else if source.StashBoxIndex != nil {
		client, err := r.getStashBoxClient(*source.StashBoxIndex)
		if err != nil {
			return nil, err
		}

		return client.FindStashBoxPerformersByPerformerNames(input.PerformerIds)
	}

	return nil, errors.New("scraper_id or stash_box_index must be set")
}

func (r *queryResolver) ScrapeSingleGallery(ctx context.Context, source models.ScraperSourceInput, input models.ScrapeSingleGalleryInput) ([]*models.ScrapedGallery, error) {
	if source.StashBoxIndex != nil {
		return nil, ErrNotSupported
	}

	if source.ScraperID == nil {
		return nil, fmt.Errorf("%w: scraper_id must be set", ErrInput)
	}

	var c models.ScrapedContent

	switch {
	case input.GalleryID != nil:
		galleryID, err := strconv.Atoi(*input.GalleryID)
		if err != nil {
			return nil, fmt.Errorf("scraper %s: converting gallery id input %s: %w", *source.ScraperID, *input.GalleryID, err)
		}
		c, err = r.scraperCache.ScrapeID(ctx, *source.ScraperID, galleryID, models.ScrapeContentTypeGallery)
		if err != nil {
			return nil, err
		}
		return marshalScrapedGalleries([]models.ScrapedContent{c})
	case input.GalleryInput != nil:
		c, err := r.scraperCache.ScrapeFragment(ctx, *source.ScraperID, scraper.Input{Gallery: input.GalleryInput})
		if err != nil {
			return nil, err
		}
		return marshalScrapedGalleries([]models.ScrapedContent{c})
	default:
		return nil, ErrNotImplemented
	}
}

func (r *queryResolver) ScrapeSingleMovie(ctx context.Context, source models.ScraperSourceInput, input models.ScrapeSingleMovieInput) ([]*models.ScrapedMovie, error) {
	return nil, ErrNotSupported
}
