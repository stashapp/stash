package scraper

import (
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

type queryURLReplacements map[string]mappedRegexConfigs

type queryURLParameters map[string]string

func queryURLParametersFromScene(scene *models.Scene) queryURLParameters {
	ret := make(queryURLParameters)
	ret["checksum"] = scene.Checksum.String
	ret["oshash"] = scene.OSHash.String
	ret["filename"] = filepath.Base(scene.Path)
	ret["title"] = scene.Title.String
	ret["url"] = scene.URL.String
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

	if gallery.Path.Valid {
		ret["filename"] = filepath.Base(gallery.Path.String)
	}
	ret["title"] = gallery.Title.String
	ret["url"] = gallery.URL.String

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
		ret = strings.Replace(ret, "{"+k+"}", v, -1)
	}

	return ret
}
