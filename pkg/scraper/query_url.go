package scraper

import (
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

type queryURLReplacements map[string]mappedRegexConfigs

type queryURLParameters map[string]string

func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func queryURLParametersFromScene(scene *models.Scene) queryURLParameters {
	ret := make(queryURLParameters)
	ret["checksum"] = stringPtrToString(scene.Checksum)
	ret["oshash"] = stringPtrToString(scene.OSHash)
	ret["filename"] = filepath.Base(scene.Path)

	if scene.Title != "" {
		ret["title"] = scene.Title
	}
	if scene.URL != "" {
		ret["url"] = scene.URL
	}
	return ret
}

func queryURLParametersFromScrapedScene(scene ScrapedSceneInput) queryURLParameters {
	ret := make(queryURLParameters)

	setField := func(field string, value *string) {
		if value != nil {
			ret[field] = *value
		}
	}

	setField("title", scene.Title)
	setField("url", scene.URL)
	setField("date", scene.Date)
	setField("details", scene.Details)
	setField("remote_site_id", scene.RemoteSiteID)
	return ret
}

func queryURLParameterFromURL(url string) queryURLParameters {
	ret := make(queryURLParameters)
	ret["url"] = url
	return ret
}

func queryURLParametersFromGallery(gallery *models.Gallery) queryURLParameters {
	ret := make(queryURLParameters)
	ret["checksum"] = gallery.Checksum

	if gallery.Path != nil {
		ret["filename"] = filepath.Base(*gallery.Path)
	}
	if gallery.Title != "" {
		ret["title"] = gallery.Title
	}

	if gallery.URL != "" {
		ret["url"] = gallery.URL
	}

	return ret
}

func (p queryURLParameters) applyReplacements(r queryURLReplacements) {
	for k, v := range p {
		rpl, found := r[k]
		if found {
			p[k] = rpl.apply(v)
		}
	}
}

func (p queryURLParameters) constructURL(url string) string {
	ret := url
	for k, v := range p {
		ret = strings.ReplaceAll(ret, "{"+k+"}", v)
	}

	return ret
}

// replaceURL does a partial URL Replace ( only url parameter is used)
func replaceURL(url string, scraperConfig scraperTypeConfig) string {
	u := url
	queryURL := queryURLParameterFromURL(u)
	if scraperConfig.QueryURLReplacements != nil {
		queryURL.applyReplacements(scraperConfig.QueryURLReplacements)
		u = queryURL.constructURL(scraperConfig.QueryURL)
	}
	return u
}
