package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/scraper/stashbox"
)

// deprecated
func (r *queryResolver) ScrapeFreeones(ctx context.Context, performer_name string) (*models.ScrapedPerformer, error) {
	scrapedPerformer := models.ScrapedPerformerInput{
		Name: &performer_name,
	}
	return manager.GetInstance().ScraperCache.ScrapePerformer(scraper.FreeonesScraperID, scrapedPerformer)
}

// deprecated
func (r *queryResolver) ScrapeFreeonesPerformerList(ctx context.Context, query string) ([]string, error) {
	scrapedPerformers, err := manager.GetInstance().ScraperCache.ScraperPerformerQuery(scraper.FreeonesScraperID, query)

	if err != nil {
		return nil, err
	}

	var ret []string
	for _, v := range scrapedPerformers {
		if v.Name != nil {
			ret = append(ret, *v.Name)
		}
	}

	return ret, nil
}

func (r *queryResolver) ListPerformerScrapers(ctx context.Context) ([]*models.Scraper, error) {
	return manager.GetInstance().ScraperCache.ListScrapers(scraper.Performer), nil
}

func (r *queryResolver) ListSceneScrapers(ctx context.Context) ([]*models.Scraper, error) {
	return manager.GetInstance().ScraperCache.ListScrapers(scraper.Scene), nil
}

func (r *queryResolver) ListGalleryScrapers(ctx context.Context) ([]*models.Scraper, error) {
	return manager.GetInstance().ScraperCache.ListScrapers(scraper.Gallery), nil
}

func (r *queryResolver) ListMovieScrapers(ctx context.Context) ([]*models.Scraper, error) {
	return manager.GetInstance().ScraperCache.ListScrapers(scraper.Movie), nil
}

func (r *queryResolver) ScrapePerformerList(ctx context.Context, scraperID string, query string) ([]*models.ScrapedPerformer, error) {
	if query == "" {
		return nil, nil
	}

	return manager.GetInstance().ScraperCache.ScraperPerformerQuery(scraperID, query)
}

func (r *queryResolver) ScrapePerformer(ctx context.Context, scraperID string, scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	return manager.GetInstance().ScraperCache.ScrapePerformer(scraperID, scrapedPerformer)
}

func (r *queryResolver) ScrapePerformerURL(ctx context.Context, url string) (*models.ScrapedPerformer, error) {
	return manager.GetInstance().ScraperCache.ScrapePerformerURL(url)
}

func (r *queryResolver) ScrapeSceneQuery(ctx context.Context, scraperID string, query string) ([]*models.ScrapedScene, error) {
	if query == "" {
		return nil, nil
	}

	return manager.GetInstance().ScraperCache.ScrapeSceneQuery(scraperID, query)
}

func (r *queryResolver) ScrapeScene(ctx context.Context, scraperID string, scene models.SceneUpdateInput) (*models.ScrapedScene, error) {
	id, err := strconv.Atoi(scene.ID)
	if err != nil {
		return nil, err
	}

	return manager.GetInstance().ScraperCache.ScrapeScene(scraperID, id)
}

func (r *queryResolver) ScrapeSceneURL(ctx context.Context, url string) (*models.ScrapedScene, error) {
	return manager.GetInstance().ScraperCache.ScrapeSceneURL(url)
}

func (r *queryResolver) ScrapeGallery(ctx context.Context, scraperID string, gallery models.GalleryUpdateInput) (*models.ScrapedGallery, error) {
	id, err := strconv.Atoi(gallery.ID)
	if err != nil {
		return nil, err
	}

	return manager.GetInstance().ScraperCache.ScrapeGallery(scraperID, id)
}

func (r *queryResolver) ScrapeGalleryURL(ctx context.Context, url string) (*models.ScrapedGallery, error) {
	return manager.GetInstance().ScraperCache.ScrapeGalleryURL(url)
}

func (r *queryResolver) ScrapeMovieURL(ctx context.Context, url string) (*models.ScrapedMovie, error) {
	return manager.GetInstance().ScraperCache.ScrapeMovieURL(url)
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
		var singleScene *models.ScrapedScene
		var err error

		switch {
		case input.SceneID != nil:
			var sceneID int
			sceneID, err = strconv.Atoi(*input.SceneID)
			if err != nil {
				return nil, err
			}
			singleScene, err = manager.GetInstance().ScraperCache.ScrapeScene(*source.ScraperID, sceneID)
		case input.SceneInput != nil:
			singleScene, err = manager.GetInstance().ScraperCache.ScrapeSceneFragment(*source.ScraperID, *input.SceneInput)
		case input.Query != nil:
			return manager.GetInstance().ScraperCache.ScrapeSceneQuery(*source.ScraperID, *input.Query)
		default:
			err = errors.New("scene_id, scene_input or query must be set")
		}

		if err != nil {
			return nil, err
		}

		if singleScene != nil {
			return []*models.ScrapedScene{singleScene}, nil
		}

		return nil, nil
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

		return nil, errors.New("scene_id or query must be set")
	}

	return nil, errors.New("scraper_id or stash_box_index must be set")
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
			singlePerformer, err := manager.GetInstance().ScraperCache.ScrapePerformer(*source.ScraperID, *input.PerformerInput)
			if err != nil {
				return nil, err
			}

			if singlePerformer != nil {
				return []*models.ScrapedPerformer{singlePerformer}, nil
			}

			return nil, nil
		}

		if input.Query != nil {
			return manager.GetInstance().ScraperCache.ScraperPerformerQuery(*source.ScraperID, *input.Query)
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
	if source.ScraperID != nil {
		var singleGallery *models.ScrapedGallery
		var err error

		switch {
		case input.GalleryID != nil:
			var galleryID int
			galleryID, err = strconv.Atoi(*input.GalleryID)
			if err != nil {
				return nil, err
			}
			singleGallery, err = manager.GetInstance().ScraperCache.ScrapeGallery(*source.ScraperID, galleryID)
		case input.GalleryInput != nil:
			singleGallery, err = manager.GetInstance().ScraperCache.ScrapeGalleryFragment(*source.ScraperID, *input.GalleryInput)
		default:
			return nil, ErrNotImplemented
		}

		if err != nil {
			return nil, err
		}

		if singleGallery != nil {
			return []*models.ScrapedGallery{singleGallery}, nil
		}

		return nil, nil
	} else if source.StashBoxIndex != nil {
		return nil, ErrNotSupported
	}

	return nil, errors.New("scraper_id must be set")
}

func (r *queryResolver) ScrapeSingleMovie(ctx context.Context, source models.ScraperSourceInput, input models.ScrapeSingleMovieInput) ([]*models.ScrapedMovie, error) {
	return nil, ErrNotSupported
}
