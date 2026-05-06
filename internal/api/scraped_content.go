package api

import (
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
)

// marshalScrapedScenes converts ScrapedContent into ScrapedScene. If conversion fails, an
// error is returned to the caller.
func marshalScrapedScenes(content []scraper.ScrapedContent) ([]*models.ScrapedScene, error) {
	var ret []*models.ScrapedScene
	for _, c := range content {
		if c == nil {
			// graphql schema requires scenes to be non-nil
			continue
		}

		switch s := c.(type) {
		case *models.ScrapedScene:
			ret = append(ret, s)
		case models.ScrapedScene:
			ret = append(ret, &s)
		default:
			return nil, fmt.Errorf("%w: cannot turn ScrapedContent into ScrapedScene", models.ErrConversion)
		}
	}

	return ret, nil
}

// marshalScrapedPerformers converts ScrapedContent into ScrapedPerformer. If conversion
// fails, an error is returned to the caller.
func marshalScrapedPerformers(content []scraper.ScrapedContent) ([]*models.ScrapedPerformer, error) {
	var ret []*models.ScrapedPerformer
	for _, c := range content {
		if c == nil {
			// graphql schema requires performers to be non-nil
			continue
		}

		switch p := c.(type) {
		case *models.ScrapedPerformer:
			ret = append(ret, p)
		case models.ScrapedPerformer:
			ret = append(ret, &p)
		default:
			return nil, fmt.Errorf("%w: cannot turn ScrapedContent into ScrapedPerformer", models.ErrConversion)
		}
	}

	return ret, nil
}

// marshalScrapedGalleries converts ScrapedContent into ScrapedGallery. If
// conversion fails, an error is returned.
func marshalScrapedGalleries(content []scraper.ScrapedContent) ([]*models.ScrapedGallery, error) {
	var ret []*models.ScrapedGallery
	for _, c := range content {
		if c == nil {
			// graphql schema requires galleries to be non-nil
			continue
		}

		switch g := c.(type) {
		case *models.ScrapedGallery:
			ret = append(ret, g)
		case models.ScrapedGallery:
			ret = append(ret, &g)
		default:
			return nil, fmt.Errorf("%w: cannot turn ScrapedContent into ScrapedGallery", models.ErrConversion)
		}
	}

	return ret, nil
}

func marshalScrapedImages(content []scraper.ScrapedContent) ([]*models.ScrapedImage, error) {
	var ret []*models.ScrapedImage
	for _, c := range content {
		if c == nil {
			// graphql schema requires images to be non-nil
			continue
		}

		switch g := c.(type) {
		case *models.ScrapedImage:
			ret = append(ret, g)
		case models.ScrapedImage:
			ret = append(ret, &g)
		default:
			return nil, fmt.Errorf("%w: cannot turn ScrapedContent into ScrapedImage", models.ErrConversion)
		}
	}

	return ret, nil
}

// marshalScrapedMovies converts ScrapedContent into ScrapedMovie. If conversion
// fails, an error is returned.
func marshalScrapedMovies(content []scraper.ScrapedContent) ([]*models.ScrapedMovie, error) {
	var ret []*models.ScrapedMovie
	for _, c := range content {
		if c == nil {
			// graphql schema requires movies to be non-nil
			continue
		}

		switch m := c.(type) {
		case *models.ScrapedMovie:
			ret = append(ret, m)
		case models.ScrapedMovie:
			ret = append(ret, &m)
		default:
			return nil, fmt.Errorf("%w: cannot turn ScrapedContent into ScrapedMovie", models.ErrConversion)
		}
	}

	return ret, nil
}

// marshalScrapedMovies converts ScrapedContent into ScrapedMovie. If conversion
// fails, an error is returned.
func marshalScrapedGroups(content []scraper.ScrapedContent) ([]*models.ScrapedGroup, error) {
	var ret []*models.ScrapedGroup
	for _, c := range content {
		if c == nil {
			// graphql schema requires groups to be non-nil
			continue
		}

		switch m := c.(type) {
		case *models.ScrapedGroup:
			ret = append(ret, m)
		case models.ScrapedGroup:
			ret = append(ret, &m)
		// it's possible that a scraper returns models.ScrapedMovie
		case *models.ScrapedMovie:
			g := m.ScrapedGroup()
			ret = append(ret, &g)
		case models.ScrapedMovie:
			g := m.ScrapedGroup()
			ret = append(ret, &g)
		default:
			return nil, fmt.Errorf("%w: cannot turn ScrapedContent into ScrapedGroup", models.ErrConversion)
		}
	}

	return ret, nil
}

// marshalScrapedPerformer will marshal a single performer
func marshalScrapedPerformer(content scraper.ScrapedContent) (*models.ScrapedPerformer, error) {
	p, err := marshalScrapedPerformers([]scraper.ScrapedContent{content})
	if err != nil {
		return nil, err
	}

	return p[0], nil
}

// marshalScrapedScene will marshal a single scraped scene
func marshalScrapedScene(content scraper.ScrapedContent) (*models.ScrapedScene, error) {
	s, err := marshalScrapedScenes([]scraper.ScrapedContent{content})
	if err != nil {
		return nil, err
	}

	return s[0], nil
}

// marshalScrapedGallery will marshal a single scraped gallery
func marshalScrapedGallery(content scraper.ScrapedContent) (*models.ScrapedGallery, error) {
	g, err := marshalScrapedGalleries([]scraper.ScrapedContent{content})
	if err != nil {
		return nil, err
	}

	return g[0], nil
}

// marshalScrapedImage will marshal a single scraped image
func marshalScrapedImage(content scraper.ScrapedContent) (*models.ScrapedImage, error) {
	g, err := marshalScrapedImages([]scraper.ScrapedContent{content})
	if err != nil {
		return nil, err
	}

	return g[0], nil
}

// marshalScrapedMovie will marshal a single scraped movie
func marshalScrapedMovie(content scraper.ScrapedContent) (*models.ScrapedMovie, error) {
	m, err := marshalScrapedMovies([]scraper.ScrapedContent{content})
	if err != nil {
		return nil, err
	}

	return m[0], nil
}

// marshalScrapedMovie will marshal a single scraped movie
func marshalScrapedGroup(content scraper.ScrapedContent) (*models.ScrapedGroup, error) {
	m, err := marshalScrapedGroups([]scraper.ScrapedContent{content})
	if err != nil {
		return nil, err
	}

	return m[0], nil
}
