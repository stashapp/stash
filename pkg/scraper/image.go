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

	// assume is a URL for now
	resp, err := http.Get(*p.Image)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	img := utils.GetBase64StringFromData(body)
	p.Image = &img

	return nil
}
