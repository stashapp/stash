package scraper

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

type queryURLReplacements map[string]mappedRegexConfigs

type queryURLParameters map[string]string

func queryURLParametersFromScene(scene *models.Scene) queryURLParameters {
	ret := make(queryURLParameters)
	ret["checksum"] = scene.Checksum
	ret["oshash"] = scene.OSHash
	ret["filename"] = filepath.Base(scene.Path)
	ret["path"] = scene.Path

	// pull phash from primary file
	phashFingerprints := scene.Files.Primary().Base().Fingerprints.Filter(models.FingerprintTypePhash)
	if len(phashFingerprints) > 0 {
		ret["phash"] = phashFingerprints[0].Value()
	}

	if scene.Title != "" {
		ret["title"] = scene.Title
	}
	if len(scene.URLs.List()) > 0 {
		ret["url"] = scene.URLs.List()[0]
	}
	return ret
}

func queryURLParametersFromScrapedScene(scene models.ScrapedSceneInput) queryURLParameters {
	ret := make(queryURLParameters)

	setField := func(field string, value *string) {
		if value != nil {
			ret[field] = *value
		}
	}

	setField("title", scene.Title)
	setField("code", scene.Code)
	if len(scene.URLs) > 0 {
		setField("url", &scene.URLs[0])
	} else {
		setField("url", scene.URL)
	}
	setField("date", scene.Date)
	setField("production_date", scene.ProductionDate)
	setField("details", scene.Details)
	setField("director", scene.Director)
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
	ret["checksum"] = gallery.PrimaryChecksum()

	if gallery.Path != "" {
		ret["filename"] = filepath.Base(gallery.Path)
	}
	if gallery.Title != "" {
		ret["title"] = gallery.Title
	}

	if len(gallery.URLs.List()) > 0 {
		ret["url"] = gallery.URLs.List()[0]
	}

	return ret
}

func queryURLParametersFromImage(image *models.Image) queryURLParameters {
	ret := make(queryURLParameters)
	ret["checksum"] = image.Checksum

	if image.Path != "" {
		ret["filename"] = filepath.Base(image.Path)
	}
	if image.Title != "" {
		ret["title"] = image.Title
	}

	if len(image.URLs.List()) > 0 {
		ret["url"] = image.URLs.List()[0]
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

func (p queryURLParameters) constructURL(urlStr string) string {
	ret := urlStr
	for k, v := range p {
		ret = strings.ReplaceAll(ret, "{"+k+"}", v)
	}

	// URL-encode special characters in the query string values
	if idx := strings.IndexByte(ret, '?'); idx >= 0 {
		queryStr := ret[idx+1:]
		// Parse key=value pairs and re-encode only the values
		pairs := strings.Split(queryStr, "&")
		encodedPairs := make([]string, 0, len(pairs))
		for _, pair := range pairs {
			if eq := strings.IndexByte(pair, '='); eq >= 0 {
				key := pair[:eq]
				value := pair[eq+1:]
				// Encode the value but not the key
				encodedPairs = append(encodedPairs, key+"="+url.QueryEscape(value))
			} else {
				encodedPairs = append(encodedPairs, pair)
			}
		}
		ret = ret[:idx+1] + strings.Join(encodedPairs, "&")
	}

	return ret
}

// replaceURL does a partial URL Replace ( only url parameter is used)
func replaceURL(url string, scraperConfig ByURLDefinition) string {
	u := url
	queryURL := queryURLParameterFromURL(u)
	if scraperConfig.QueryURLReplacements != nil {
		queryURL.applyReplacements(scraperConfig.QueryURLReplacements)
		u = queryURL.constructURL(scraperConfig.QueryURL)
	}
	return u
}
