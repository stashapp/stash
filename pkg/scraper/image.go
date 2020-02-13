package scraper

import (
	"io/ioutil"
	"net/http"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func setPerformerImage(p *models.ScrapedPerformer) error {
	if p == nil || p.Image == nil {
		// nothing to do
		return nil
	}

	img, err := getImage(*p.Image)
	if err != nil {
		return err
	}

	p.Image = img

	return nil
}

func setSceneImage(s *models.ScrapedScene) error {
	if s == nil || s.Image == nil {
		// nothing to do
		return nil
	}

	img, err := getImage(*s.Image)
	if err != nil {
		return err
	}

	s.Image = img

	return nil
}

func getImage(url string) (*string, error) {
	// assume is a URL for now
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// determine the image type and set the base64 type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(body)
	}

	img := "data:" + contentType + ";base64," + utils.GetBase64StringFromData(body)
	return &img, nil
}
