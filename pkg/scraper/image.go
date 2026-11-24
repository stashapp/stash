package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func setPerformerImage(ctx context.Context, getter imageGetter, p *models.ScrapedPerformer) error {
	if len(p.Images) > 0 {
		getter.fetchSurfImages(ctx, p.Images)
		return nil
	}

	// backwards compatibility: we fetch the image if it's a URL and set it to the first image
	// Image is deprecated, so only do this if Images is unset
	if p.Image == nil {
		// nothing to do
		return nil
	}

	// don't try to get the image if it doesn't appear to be a URL
	if !strings.HasPrefix(*p.Image, "http") {
		p.Images = []string{*p.Image}
		return nil
	}

	img, err := getter.getImage(ctx, *p.Image)
	if err != nil {
		return err
	}

	p.Image = img
	// Image is deprecated. Use images instead
	p.Images = []string{*img}

	return nil
}

func setStudioImage(ctx context.Context, getter imageGetter, p *models.ScrapedStudio) error {
	if len(p.Images) > 0 {
		getter.fetchSurfImages(ctx, p.Images)
		return nil
	}

	// backwards compatibility: we fetch the image if it's a URL and set it to the first image
	// Image is deprecated, so only do this if Images is unset
	if p.Image == nil {
		// nothing to do
		return nil
	}

	// don't try to get the image if it doesn't appear to be a URL
	if !strings.HasPrefix(*p.Image, "http") {
		p.Images = []string{*p.Image}
		return nil
	}

	img, err := getter.getImage(ctx, *p.Image)
	if err != nil {
		return err
	}

	p.Image = img
	// Image is deprecated. Use images instead
	p.Images = []string{*img}

	return nil
}

func processImageField(ctx context.Context, imageField *string, getter imageGetter) error {
	if imageField == nil {
		return nil
	}

	// don't try to get the image if it doesn't appear to be a URL
	// this allows scrapers to return base64 data URIs directly
	if !strings.HasPrefix(*imageField, "http") {
		return nil
	}

	img, err := getter.getImage(ctx, *imageField)
	if err != nil {
		return err
	}

	*imageField = *img
	return nil
}

// fetchSurfImages downloads URL images of a Surf scraper into data URIs.
// The backend later re-fetches URL images with a plain client, which Surf
// scrapers' sites may block. Images that fail to download are left as URLs.
func (i *imageGetter) fetchSurfImages(ctx context.Context, images []string) {
	if i.surfDef == nil {
		return
	}

	for j := range images {
		if err := processImageField(ctx, &images[j], *i); err != nil {
			logger.Warnf("could not fetch image using URL %s: %v", images[j], err)
		}
	}
}

type imageGetter struct {
	client       *http.Client
	globalConfig GlobalConfig
	// surfDef is set when images must be fetched with the scraper's Surf client
	surfDef         *Definition
	requestModifier func(req *http.Request)
}

func newImageGetter(client *http.Client, globalConfig GlobalConfig, definition *Definition) imageGetter {
	ret := imageGetter{
		client:       client,
		globalConfig: globalConfig,
	}

	if definition != nil && definition.usesSurf() {
		ret.surfDef = definition
	}

	return ret
}

func (i *imageGetter) getImage(ctx context.Context, url string) (*string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := i.client
	if i.surfDef != nil {
		// Surf sets its own browser User-Agent and the scrape's cookie jar
		client, _, err = surfClient(ctx, *i.surfDef, i.globalConfig)
		if err != nil {
			return nil, err
		}
	} else if userAgent := i.globalConfig.GetScraperUserAgent(); userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	// assume is a URL for now

	// set the host of the URL as the referer
	if req.URL.Scheme != "" {
		req.Header.Set("Referer", req.URL.Scheme+"://"+req.Host+"/")
	}

	if i.requestModifier != nil {
		i.requestModifier(req)
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http error %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
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

func getStashPerformerImage(ctx context.Context, stashURL string, performerID string, imageGetter imageGetter) (*string, error) {
	return imageGetter.getImage(ctx, stashURL+"/performer/"+performerID+"/image")
}

func getStashSceneImage(ctx context.Context, stashURL string, sceneID string, imageGetter imageGetter) (*string, error) {
	return imageGetter.getImage(ctx, stashURL+"/scene/"+sceneID+"/screenshot")
}
