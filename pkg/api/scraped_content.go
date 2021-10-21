package api

import (
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

// ErrConversion signifies conversion errors
var ErrConversion = errors.New("conversion error")

// marshalScrapedScenes converts ScrapedContent into ScrapedScene. If conversion fails, an
// error is returned to the caller.
func marshalScrapedScenes(content []models.ScrapedContent) ([]*models.ScrapedScene, error) {
	var ret []*models.ScrapedScene
	for _, c := range content {
		if s, ok := c.(*models.ScrapedScene); ok {
			ret = append(ret, s)
		} else {
			return nil, fmt.Errorf("cannot turn ScrapedContent into ScrapedScene: %w", ErrConversion)
		}
	}

	return ret, nil
}

// marshalScrapedPerformers converts ScrapedContent into ScrapedPerformer. If conversion
// fails, an error is returned to the caller.
func marshalScrapedPerformers(content []models.ScrapedContent) ([]*models.ScrapedPerformer, error) {
	var ret []*models.ScrapedPerformer
	for _, c := range content {
		if s, ok := c.(*models.ScrapedPerformer); ok {
			ret = append(ret, s)
		} else {
			return nil, fmt.Errorf("cannot turn ScrapedContent into ScrapedPerformer: %w", ErrConversion)
		}
	}

	return ret, nil
}

// marshalScrapedGalleries converts ScrapedContent into ScrapedGallery. If
// conversion fails, an error is returned.
func marshalScrapedGalleries(content []models.ScrapedContent) ([]*models.ScrapedGallery, error) {
	var ret []*models.ScrapedGallery
	for _, c := range content {
		if s, ok := c.(*models.ScrapedGallery); ok {
			ret = append(ret, s)
		} else {
			return nil, fmt.Errorf("cannot turn ScrapedContent into ScrapedGallery: %w", ErrConversion)
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
