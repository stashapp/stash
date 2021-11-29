package api

import (
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

// marshalScrapedScenes converts ScrapedContent into ScrapedScene. If conversion fails, an
// error is returned to the caller.
func marshalScrapedScenes(content []models.ScrapedContent) ([]*models.ScrapedScene, error) {
	var ret []*models.ScrapedScene
	for _, c := range content {
		if c == nil {
			ret = append(ret, nil)
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
func marshalScrapedPerformers(content []models.ScrapedContent) ([]*models.ScrapedPerformer, error) {
	var ret []*models.ScrapedPerformer
	for _, c := range content {
		if c == nil {
			ret = append(ret, nil)
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
func marshalScrapedGalleries(content []models.ScrapedContent) ([]*models.ScrapedGallery, error) {
	var ret []*models.ScrapedGallery
	for _, c := range content {
		if c == nil {
			ret = append(ret, nil)
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

// marshalScrapedMovies converts ScrapedContent into ScrapedMovie. If conversion
// fails, an error is returned.
func marshalScrapedMovies(content []models.ScrapedContent) ([]*models.ScrapedMovie, error) {
	var ret []*models.ScrapedMovie
	for _, c := range content {
		if c == nil {
			ret = append(ret, nil)
			continue
		}

		switch m := c.(type) {
		case *models.ScrapedMovie:
			ret = append(ret, m)
		case models.ScrapedMovie:
			ret = append(ret, &m)
		default:
			return nil, fmt.Errorf("%w: cannot turn ScrapedConetnt into ScrapedMovie", models.ErrConversion)
		}
	}

	return ret, nil
}

// marshalScrapedPerformer will marshal a single performer
func marshalScrapedPerformer(content models.ScrapedContent) (*models.ScrapedPerformer, error) {
	p, err := marshalScrapedPerformers([]models.ScrapedContent{content})
	if err != nil {
		return nil, err
	}

	return p[0], nil
}

// marshalScrapedScene will marshal a single scraped scene
func marshalScrapedScene(content models.ScrapedContent) (*models.ScrapedScene, error) {
	s, err := marshalScrapedScenes([]models.ScrapedContent{content})
	if err != nil {
		return nil, err
	}

	return s[0], nil
}

// marshalScrapedGallery will marshal a single scraped gallery
func marshalScrapedGallery(content models.ScrapedContent) (*models.ScrapedGallery, error) {
	g, err := marshalScrapedGalleries([]models.ScrapedContent{content})
	if err != nil {
		return nil, err
	}

	return g[0], nil
}

// marshalScrapedMovie will marshal a single scraped movie
func marshalScrapedMovie(content models.ScrapedContent) (*models.ScrapedMovie, error) {
	m, err := marshalScrapedMovies([]models.ScrapedContent{content})
	if err != nil {
		return nil, err
	}

	return m[0], nil
}
