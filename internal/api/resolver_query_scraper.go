package api

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

func (r *queryResolver) ScrapeURL(ctx context.Context, url string, ty scraper.ScrapeContentType) (scraper.ScrapedContent, error) {
	return r.scraperCache().ScrapeURL(ctx, url, ty)
}

func (r *queryResolver) ListScrapers(ctx context.Context, types []scraper.ScrapeContentType) ([]*scraper.Scraper, error) {
	return r.scraperCache().ListScrapers(types), nil
}

func (r *queryResolver) ScrapePerformerURL(ctx context.Context, url string) (*models.ScrapedPerformer, error) {
	content, err := r.scraperCache().ScrapeURL(ctx, url, scraper.ScrapeContentTypePerformer)
	if err != nil {
		return nil, err
	}

	return marshalScrapedPerformer(content)
}

func (r *queryResolver) ScrapeSceneQuery(ctx context.Context, scraperID string, query string) ([]*models.ScrapedScene, error) {
	if query == "" {
		return nil, nil
	}

	content, err := r.scraperCache().ScrapeName(ctx, scraperID, query, scraper.ScrapeContentTypeScene)
	if err != nil {
		return nil, err
	}

	ret, err := marshalScrapedScenes(content)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) ScrapeSceneURL(ctx context.Context, url string) (*models.ScrapedScene, error) {
	content, err := r.scraperCache().ScrapeURL(ctx, url, scraper.ScrapeContentTypeScene)
	if err != nil {
		return nil, err
	}

	ret, err := marshalScrapedScene(content)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) ScrapeGalleryURL(ctx context.Context, url string) (*models.ScrapedGallery, error) {
	content, err := r.scraperCache().ScrapeURL(ctx, url, scraper.ScrapeContentTypeGallery)
	if err != nil {
		return nil, err
	}

	ret, err := marshalScrapedGallery(content)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) ScrapeImageURL(ctx context.Context, url string) (*models.ScrapedImage, error) {
	content, err := r.scraperCache().ScrapeURL(ctx, url, scraper.ScrapeContentTypeImage)
	if err != nil {
		return nil, err
	}

	return marshalScrapedImage(content)
}

func (r *queryResolver) ScrapeMovieURL(ctx context.Context, url string) (*models.ScrapedMovie, error) {
	content, err := r.scraperCache().ScrapeURL(ctx, url, scraper.ScrapeContentTypeMovie)
	if err != nil {
		return nil, err
	}

	ret, err := marshalScrapedMovie(content)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) ScrapeGroupURL(ctx context.Context, url string) (*models.ScrapedGroup, error) {
	content, err := r.scraperCache().ScrapeURL(ctx, url, scraper.ScrapeContentTypeGroup)
	if err != nil {
		return nil, err
	}

	ret, err := marshalScrapedGroup(content)
	if err != nil {
		return nil, err
	}

	// convert to scraped group
	group := &models.ScrapedGroup{
		StoredID:   ret.StoredID,
		Name:       ret.Name,
		Aliases:    ret.Aliases,
		Duration:   ret.Duration,
		Date:       ret.Date,
		Rating:     ret.Rating,
		Director:   ret.Director,
		URLs:       ret.URLs,
		Synopsis:   ret.Synopsis,
		Studio:     ret.Studio,
		Tags:       ret.Tags,
		FrontImage: ret.FrontImage,
		BackImage:  ret.BackImage,
	}

	return group, nil
}

func (r *queryResolver) ScrapeSingleScene(ctx context.Context, source scraper.Source, input ScrapeSingleSceneInput) ([]*models.ScrapedScene, error) {
	var ret []*models.ScrapedScene

	var sceneID int
	if input.SceneID != nil {
		var err error
		sceneID, err = strconv.Atoi(*input.SceneID)
		if err != nil {
			return nil, fmt.Errorf("%w: sceneID is not an integer: '%s'", ErrInput, *input.SceneID)
		}
	}

	switch {
	case source.ScraperID != nil:
		var err error
		var c scraper.ScrapedContent
		var content []scraper.ScrapedContent

		switch {
		case input.SceneID != nil:
			c, err = r.scraperCache().ScrapeID(ctx, *source.ScraperID, sceneID, scraper.ScrapeContentTypeScene)
			if c != nil {
				content = []scraper.ScrapedContent{c}
			}
		case input.SceneInput != nil:
			c, err = r.scraperCache().ScrapeFragment(ctx, *source.ScraperID, scraper.Input{Scene: input.SceneInput})
			if c != nil {
				content = []scraper.ScrapedContent{c}
			}
		case input.Query != nil:
			content, err = r.scraperCache().ScrapeName(ctx, *source.ScraperID, *input.Query, scraper.ScrapeContentTypeScene)
		default:
			err = fmt.Errorf("%w: scene_id, scene_input, or query must be set", ErrInput)
		}

		if err != nil {
			return nil, err
		}

		ret, err = marshalScrapedScenes(content)
		if err != nil {
			return nil, err
		}
	case source.StashBoxIndex != nil || source.StashBoxEndpoint != nil:
		b, err := resolveStashBox(source.StashBoxIndex, source.StashBoxEndpoint)
		if err != nil {
			return nil, err
		}

		client := r.newStashBoxClient(*b)

		switch {
		case input.SceneID != nil:
			var fps []models.Fingerprints
			fps, err = r.getScenesFingerprints(ctx, []int{sceneID})
			if err != nil {
				return nil, err
			}
			ret, err = client.FindSceneByFingerprints(ctx, fps[0])
		case input.Query != nil:
			ret, err = client.QueryScene(ctx, *input.Query)
		default:
			return nil, fmt.Errorf("%w: scene_id or query must be set", ErrInput)
		}

		if err != nil {
			return nil, err
		}

		// TODO - this should happen after any scene is scraped
		if err := r.matchScenesRelationships(ctx, ret, b.Endpoint); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%w: scraper_id or stash_box_index must be set", ErrInput)
	}

	for i := range ret {
		slices.SortFunc(ret[i].Tags, models.ScrapedTagSortFunction)
	}

	return ret, nil
}

func (r *queryResolver) ScrapeMultiScenes(ctx context.Context, source scraper.Source, input ScrapeMultiScenesInput) ([][]*models.ScrapedScene, error) {
	if source.ScraperID != nil {
		return nil, ErrNotImplemented
	} else if source.StashBoxIndex != nil || source.StashBoxEndpoint != nil {
		b, err := resolveStashBox(source.StashBoxIndex, source.StashBoxEndpoint)
		if err != nil {
			return nil, err
		}

		client := r.newStashBoxClient(*b)

		sceneIDs, err := stringslice.StringSliceToIntSlice(input.SceneIds)
		if err != nil {
			return nil, err
		}

		fps, err := r.getScenesFingerprints(ctx, sceneIDs)
		if err != nil {
			return nil, err
		}

		ret, err := client.FindScenesByFingerprints(ctx, fps)
		if err != nil {
			return nil, err
		}

		// match relationships - this mutates the existing scenes so we can
		// just flatten the slice and pass it in
		flat := sliceutil.Flatten(ret)

		if err := r.matchScenesRelationships(ctx, flat, b.Endpoint); err != nil {
			return nil, err
		}

		return ret, nil
	}

	return nil, errors.New("scraper_id or stash_box_index must be set")
}

func (r *queryResolver) getScenesFingerprints(ctx context.Context, ids []int) ([]models.Fingerprints, error) {
	fingerprints := make([]models.Fingerprints, len(ids))

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Scene

		for i, sceneID := range ids {
			scene, err := qb.Find(ctx, sceneID)
			if err != nil {
				return err
			}

			if scene == nil {
				return fmt.Errorf("scene with id %d not found", sceneID)
			}

			if err := scene.LoadFiles(ctx, qb); err != nil {
				return err
			}

			var sceneFPs models.Fingerprints

			for _, f := range scene.Files.List() {
				sceneFPs = append(sceneFPs, f.Fingerprints...)
			}

			fingerprints[i] = sceneFPs
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return fingerprints, nil
}

// matchSceneRelationships accepts scraped scenes and attempts to match its relationships to existing stash models.
func (r *queryResolver) matchScenesRelationships(ctx context.Context, ss []*models.ScrapedScene, endpoint string) error {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		matcher := match.SceneRelationships{
			PerformerFinder: r.repository.Performer,
			TagFinder:       r.repository.Tag,
			StudioFinder:    r.repository.Studio,
		}

		for _, s := range ss {
			if err := matcher.MatchRelationships(ctx, s, endpoint); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *queryResolver) ScrapeSingleStudio(ctx context.Context, source scraper.Source, input ScrapeSingleStudioInput) ([]*models.ScrapedStudio, error) {
	if source.StashBoxIndex != nil || source.StashBoxEndpoint != nil {
		b, err := resolveStashBox(source.StashBoxIndex, source.StashBoxEndpoint)
		if err != nil {
			return nil, err
		}

		client := r.newStashBoxClient(*b)

		var ret []*models.ScrapedStudio
		out, err := client.FindStudio(ctx, *input.Query)

		if err != nil {
			return nil, err
		} else if out != nil {
			ret = append(ret, out)
		}

		if len(ret) > 0 {
			if err := r.withReadTxn(ctx, func(ctx context.Context) error {
				for _, studio := range ret {
					if err := match.ScrapedStudioHierarchy(ctx, r.repository.Studio, studio, b.Endpoint); err != nil {
						return err
					}
				}

				return nil
			}); err != nil {
				return nil, err
			}
			return ret, nil
		}

		return nil, nil
	}

	return nil, errors.New("stash_box_endpoint must be set")
}

func (r *queryResolver) ScrapeSingleTag(ctx context.Context, source scraper.Source, input ScrapeSingleTagInput) ([]*models.ScrapedTag, error) {
	if source.StashBoxIndex != nil || source.StashBoxEndpoint != nil {
		b, err := resolveStashBox(source.StashBoxIndex, source.StashBoxEndpoint)
		if err != nil {
			return nil, err
		}

		client := r.newStashBoxClient(*b)

		var ret []*models.ScrapedTag
		out, err := client.QueryTag(ctx, *input.Query)

		if err != nil {
			return nil, err
		} else if out != nil {
			ret = append(ret, out...)
		}

		if len(ret) > 0 {
			if err := r.withReadTxn(ctx, func(ctx context.Context) error {
				for _, tag := range ret {
					if err := match.ScrapedTag(ctx, r.repository.Tag, tag, b.Endpoint); err != nil {
						return err
					}
				}

				return nil
			}); err != nil {
				return nil, err
			}
			return ret, nil
		}

		return nil, nil
	}

	return nil, errors.New("stash_box_endpoint must be set")
}

func (r *queryResolver) ScrapeSinglePerformer(ctx context.Context, source scraper.Source, input ScrapeSinglePerformerInput) ([]*models.ScrapedPerformer, error) {
	var ret []*models.ScrapedPerformer
	switch {
	case source.ScraperID != nil:
		switch {
		case input.PerformerInput != nil:
			performer, err := r.scraperCache().ScrapeFragment(ctx, *source.ScraperID, scraper.Input{Performer: input.PerformerInput})
			if err != nil {
				return nil, err
			}

			ret, err = marshalScrapedPerformers([]scraper.ScrapedContent{performer})
			if err != nil {
				return nil, err
			}
		case input.Query != nil:
			content, err := r.scraperCache().ScrapeName(ctx, *source.ScraperID, *input.Query, scraper.ScrapeContentTypePerformer)
			if err != nil {
				return nil, err
			}

			ret, err = marshalScrapedPerformers(content)
			if err != nil {
				return nil, err
			}
		default:
			return nil, ErrNotImplemented
		}
	case source.StashBoxIndex != nil || source.StashBoxEndpoint != nil:
		b, err := resolveStashBox(source.StashBoxIndex, source.StashBoxEndpoint)
		if err != nil {
			return nil, err
		}

		client := r.newStashBoxClient(*b)

		var query string
		switch {
		case input.PerformerID != nil:
			names, err := r.findPerformerNames(ctx, []string{*input.PerformerID})
			if err != nil {
				return nil, err
			}

			query = names[0]
		case input.Query != nil:
			query = *input.Query
		default:
			return nil, ErrNotImplemented
		}

		if query == "" {
			return nil, nil
		}
		ret, err = client.QueryPerformer(ctx, query)

		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("scraper_id or stash_box_index must be set")
	}

	return ret, nil
}

func (r *queryResolver) ScrapeMultiPerformers(ctx context.Context, source scraper.Source, input ScrapeMultiPerformersInput) ([][]*models.ScrapedPerformer, error) {
	if source.ScraperID != nil {
		return nil, ErrNotImplemented
	} else if source.StashBoxIndex != nil || source.StashBoxEndpoint != nil {
		names, err := r.findPerformerNames(ctx, input.PerformerIds)
		if err != nil {
			return nil, err
		}

		b, err := resolveStashBox(source.StashBoxIndex, source.StashBoxEndpoint)
		if err != nil {
			return nil, err
		}

		client := r.newStashBoxClient(*b)

		return client.QueryPerformers(ctx, names)
	}

	return nil, errors.New("scraper_id or stash_box_index must be set")
}

func (r *queryResolver) findPerformerNames(ctx context.Context, performerIDs []string) ([]string, error) {
	ids, err := stringslice.StringSliceToIntSlice(performerIDs)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(ids))

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		p, err := r.repository.Performer.FindMany(ctx, ids)
		if err != nil {
			return err
		}

		for i, pp := range p {
			names[i] = pp.Name
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return names, nil
}

func (r *queryResolver) ScrapeSingleGallery(ctx context.Context, source scraper.Source, input ScrapeSingleGalleryInput) ([]*models.ScrapedGallery, error) {
	var ret []*models.ScrapedGallery

	if source.StashBoxIndex != nil || source.StashBoxEndpoint != nil {
		return nil, ErrNotSupported
	}

	if source.ScraperID == nil {
		return nil, fmt.Errorf("%w: scraper_id must be set", ErrInput)
	}

	var c scraper.ScrapedContent

	switch {
	case input.GalleryID != nil:
		galleryID, err := strconv.Atoi(*input.GalleryID)
		if err != nil {
			return nil, fmt.Errorf("%w: gallery id is not an integer: '%s'", ErrInput, *input.GalleryID)
		}
		c, err = r.scraperCache().ScrapeID(ctx, *source.ScraperID, galleryID, scraper.ScrapeContentTypeGallery)
		if err != nil {
			return nil, err
		}
		ret, err = marshalScrapedGalleries([]scraper.ScrapedContent{c})
		if err != nil {
			return nil, err
		}
	case input.GalleryInput != nil:
		c, err := r.scraperCache().ScrapeFragment(ctx, *source.ScraperID, scraper.Input{Gallery: input.GalleryInput})
		if err != nil {
			return nil, err
		}
		ret, err = marshalScrapedGalleries([]scraper.ScrapedContent{c})
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrNotImplemented
	}

	return ret, nil
}

func (r *queryResolver) ScrapeSingleImage(ctx context.Context, source scraper.Source, input ScrapeSingleImageInput) ([]*models.ScrapedImage, error) {
	if source.StashBoxIndex != nil {
		return nil, ErrNotSupported
	}

	if source.ScraperID == nil {
		return nil, fmt.Errorf("%w: scraper_id must be set", ErrInput)
	}

	var c scraper.ScrapedContent

	switch {
	case input.ImageID != nil:
		imageID, err := strconv.Atoi(*input.ImageID)
		if err != nil {
			return nil, fmt.Errorf("%w: image id is not an integer: '%s'", ErrInput, *input.ImageID)
		}
		c, err = r.scraperCache().ScrapeID(ctx, *source.ScraperID, imageID, scraper.ScrapeContentTypeImage)
		if err != nil {
			return nil, err
		}
		return marshalScrapedImages([]scraper.ScrapedContent{c})
	case input.ImageInput != nil:
		c, err := r.scraperCache().ScrapeFragment(ctx, *source.ScraperID, scraper.Input{Image: input.ImageInput})
		if err != nil {
			return nil, err
		}
		return marshalScrapedImages([]scraper.ScrapedContent{c})
	default:
		return nil, ErrNotImplemented
	}
}

func (r *queryResolver) ScrapeSingleMovie(ctx context.Context, source scraper.Source, input ScrapeSingleMovieInput) ([]*models.ScrapedMovie, error) {
	return nil, ErrNotSupported
}

func (r *queryResolver) ScrapeSingleGroup(ctx context.Context, source scraper.Source, input ScrapeSingleGroupInput) ([]*models.ScrapedGroup, error) {
	return nil, ErrNotSupported
}
