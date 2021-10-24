package scraper

import (
	"github.com/stashapp/stash/pkg/models"
)

func createScraperFromConfig(c config, txnManager models.TransactionManager, globalConfig GlobalConfig) scraper {
	return group{
		config:     c,
		txnManager: txnManager,
		globalConf: globalConfig,
	}
}

func (c config) spec() models.Scraper {
	ret := models.Scraper{
		ID:   c.ID,
		Name: c.Name,
	}

	performer := models.ScraperSpec{}
	if c.PerformerByName != nil {
		performer.SupportedScrapes = append(performer.SupportedScrapes, models.ScrapeTypeName)
	}
	if c.PerformerByFragment != nil {
		performer.SupportedScrapes = append(performer.SupportedScrapes, models.ScrapeTypeFragment)
	}
	if len(c.PerformerByURL) > 0 {
		performer.SupportedScrapes = append(performer.SupportedScrapes, models.ScrapeTypeURL)
		for _, v := range c.PerformerByURL {
			performer.Urls = append(performer.Urls, v.URL...)
		}
	}

	if len(performer.SupportedScrapes) > 0 {
		ret.Performer = &performer
	}

	scene := models.ScraperSpec{}
	if c.SceneByFragment != nil {
		scene.SupportedScrapes = append(scene.SupportedScrapes, models.ScrapeTypeFragment)
	}
	if c.SceneByName != nil && c.SceneByQueryFragment != nil {
		scene.SupportedScrapes = append(scene.SupportedScrapes, models.ScrapeTypeName)
	}
	if len(c.SceneByURL) > 0 {
		scene.SupportedScrapes = append(scene.SupportedScrapes, models.ScrapeTypeURL)
		for _, v := range c.SceneByURL {
			scene.Urls = append(scene.Urls, v.URL...)
		}
	}

	if len(scene.SupportedScrapes) > 0 {
		ret.Scene = &scene
	}

	gallery := models.ScraperSpec{}
	if c.GalleryByFragment != nil {
		gallery.SupportedScrapes = append(gallery.SupportedScrapes, models.ScrapeTypeFragment)
	}
	if len(c.GalleryByURL) > 0 {
		gallery.SupportedScrapes = append(gallery.SupportedScrapes, models.ScrapeTypeURL)
		for _, v := range c.GalleryByURL {
			gallery.Urls = append(gallery.Urls, v.URL...)
		}
	}

	if len(gallery.SupportedScrapes) > 0 {
		ret.Gallery = &gallery
	}

	movie := models.ScraperSpec{}
	if len(c.MovieByURL) > 0 {
		movie.SupportedScrapes = append(movie.SupportedScrapes, models.ScrapeTypeURL)
		for _, v := range c.MovieByURL {
			movie.Urls = append(movie.Urls, v.URL...)
		}
	}

	if len(movie.SupportedScrapes) > 0 {
		ret.Movie = &movie
	}

	return ret
}
