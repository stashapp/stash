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

func setSceneImage(ctx context.Context, client *http.Client, s *ScrapedScene, globalConfig GlobalConfig) error {
	// don't try to get the image if it doesn't appear to be a URL
	if s.Image == nil || !strings.HasPrefix(*s.Image, "http") {
		// nothing to do
		return nil
	}

	img, err := getImage(ctx, *s.Image, client, globalConfig)
	if err != nil {
		return err
	}

	s.Image = img

	return nil
}

func setMovieFrontImage(ctx context.Context, client *http.Client, m *models.ScrapedMovie, globalConfig GlobalConfig) error {
	// don't try to get the image if it doesn't appear to be a URL
	if m.FrontImage == nil || !strings.HasPrefix(*m.FrontImage, "http") {
		// nothing to do
		return nil
	}

	img, err := getImage(ctx, *m.FrontImage, client, globalConfig)
	if err != nil {
		return err
	}

	m.FrontImage = img

	return nil
}

func setMovieBackImage(ctx context.Context, client *http.Client, m *models.ScrapedMovie, globalConfig GlobalConfig) error {
	// don't try to get the image if it doesn't appear to be a URL
	if m.BackImage == nil || !strings.HasPrefix(*m.BackImage, "http") {
		// nothing to do
		return nil
	}

	img, err := getImage(ctx, *m.BackImage, client, globalConfig)
	if err != nil {
		return err
	}

	m.BackImage = img

	return nil
}

func setGroupFrontImage(ctx context.Context, client *http.Client, m *models.ScrapedGroup, globalConfig GlobalConfig) error {
	// don't try to get the image if it doesn't appear to be a URL
	if m.FrontImage == nil || !strings.HasPrefix(*m.FrontImage, "http") {
		// nothing to do
		return nil
	}

	img, err := getImage(ctx, *m.FrontImage, client, globalConfig)
	if err != nil {
		return err
	}

	m.FrontImage = img

	return nil
}

func setGroupBackImage(ctx context.Context, client *http.Client, m *models.ScrapedGroup, globalConfig GlobalConfig) error {
	// don't try to get the image if it doesn't appear to be a URL
	if m.BackImage == nil || !strings.HasPrefix(*m.BackImage, "http") {
		// nothing to do
		return nil
	}

	img, err := getImage(ctx, *m.BackImage, client, globalConfig)
	if err != nil {
		return err
	}

	m.BackImage = img

	return nil
}

func getImage(ctx context.Context, url string, client *http.Client, globalConfig GlobalConfig) (*string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	userAgent := globalConfig.GetScraperUserAgent()
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	// assume is a URL for now

	// set the host of the URL as the referer
	if req.URL.Scheme != "" {
		req.Header.Set("Referer", req.URL.Scheme+"://"+req.Host+"/")
	}

	resp, err := client.Do(req)

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

func getStashPerformerImage(ctx context.Context, stashURL string, performerID string, client *http.Client, globalConfig GlobalConfig) (*string, error) {
	return getImage(ctx, stashURL+"/performer/"+performerID+"/image", client, globalConfig)
}

func getStashSceneImage(ctx context.Context, stashURL string, sceneID string, client *http.Client, globalConfig GlobalConfig) (*string, error) {
	return getImage(ctx, stashURL+"/scene/"+sceneID+"/screenshot", client, globalConfig)
}
