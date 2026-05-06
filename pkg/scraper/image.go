package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func setPerformerImage(ctx context.Context, client *http.Client, p *models.ScrapedPerformer, globalConfig GlobalConfig) error {
	// backwards compatibility: we fetch the image if it's a URL and set it to the first image
	// Image is deprecated, so only do this if Images is unset
	if p.Image == nil || len(p.Images) > 0 {
		// nothing to do
		return nil
	}

	// don't try to get the image if it doesn't appear to be a URL
	if !strings.HasPrefix(*p.Image, "http") {
		p.Images = []string{*p.Image}
		return nil
	}

	img, err := getImage(ctx, *p.Image, client, globalConfig)
	if err != nil {
		return err
	}

	p.Image = img
	// Image is deprecated. Use images instead
	p.Images = []string{*img}

	return nil
}

func setStudioImage(ctx context.Context, client *http.Client, p *models.ScrapedStudio, globalConfig GlobalConfig) error {
	// backwards compatibility: we fetch the image if it's a URL and set it to the first image
	// Image is deprecated, so only do this if Images is unset
	if p.Image == nil || len(p.Images) > 0 {
		// nothing to do
		return nil
	}

	// don't try to get the image if it doesn't appear to be a URL
	if !strings.HasPrefix(*p.Image, "http") {
		p.Images = []string{*p.Image}
		return nil
	}

	img, err := getImage(ctx, *p.Image, client, globalConfig)
	if err != nil {
		return err
	}

	p.Image = img
	// Image is deprecated. Use images instead
	p.Images = []string{*img}

	return nil
}

func processImageField(ctx context.Context, imageField *string, client *http.Client, globalConfig GlobalConfig) error {
	if imageField == nil {
		return nil
	}

	// don't try to get the image if it doesn't appear to be a URL
	// this allows scrapers to return base64 data URIs directly
	if !strings.HasPrefix(*imageField, "http") {
		return nil
	}

	img, err := getImage(ctx, *imageField, client, globalConfig)
	if err != nil {
		return err
	}

	*imageField = *img
	return nil
}

type imageGetter struct {
	client          *http.Client
	globalConfig    GlobalConfig
	requestModifier func(req *http.Request)
}

func (i *imageGetter) getImage(ctx context.Context, url string) (*string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	userAgent := i.globalConfig.GetScraperUserAgent()
	if userAgent != "" {
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

	resp, err := i.client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http error %d", resp.StatusCode)
	}

	defer resp.Body.Close()

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

func getImage(ctx context.Context, url string, client *http.Client, globalConfig GlobalConfig) (*string, error) {
	g := imageGetter{
		client:       client,
		globalConfig: globalConfig,
	}

	return g.getImage(ctx, url)
}

func getStashPerformerImage(ctx context.Context, stashURL string, performerID string, imageGetter imageGetter) (*string, error) {
	return imageGetter.getImage(ctx, stashURL+"/performer/"+performerID+"/image")
}

func getStashSceneImage(ctx context.Context, stashURL string, sceneID string, imageGetter imageGetter) (*string, error) {
	return imageGetter.getImage(ctx, stashURL+"/scene/"+sceneID+"/screenshot")
}
