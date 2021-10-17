package scraper

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

// Timeout to get the image. Includes transfer time. May want to make this
// configurable at some point.
const imageGetTimeout = time.Second * 30

func setPerformerImage(ctx context.Context, p *models.ScrapedPerformer, globalConfig GlobalConfig) error {
	if p == nil || p.Image == nil || !strings.HasPrefix(*p.Image, "http") {
		// nothing to do
		return nil
	}

	img, err := getImage(ctx, *p.Image, globalConfig)
	if err != nil {
		return err
	}

	p.Image = img
	// Image is deprecated. Use images instead
	p.Images = []string{*img}

	return nil
}

func setSceneImage(ctx context.Context, s *models.ScrapedScene, globalConfig GlobalConfig) error {
	// don't try to get the image if it doesn't appear to be a URL
	if s == nil || s.Image == nil || !strings.HasPrefix(*s.Image, "http") {
		// nothing to do
		return nil
	}

	img, err := getImage(ctx, *s.Image, globalConfig)
	if err != nil {
		return err
	}

	s.Image = img

	return nil
}

func setMovieFrontImage(ctx context.Context, m *models.ScrapedMovie, globalConfig GlobalConfig) error {
	// don't try to get the image if it doesn't appear to be a URL
	if m == nil || m.FrontImage == nil || !strings.HasPrefix(*m.FrontImage, "http") {
		// nothing to do
		return nil
	}

	img, err := getImage(ctx, *m.FrontImage, globalConfig)
	if err != nil {
		return err
	}

	m.FrontImage = img

	return nil
}

func setMovieBackImage(ctx context.Context, m *models.ScrapedMovie, globalConfig GlobalConfig) error {
	// don't try to get the image if it doesn't appear to be a URL
	if m == nil || m.BackImage == nil || !strings.HasPrefix(*m.BackImage, "http") {
		// nothing to do
		return nil
	}

	img, err := getImage(ctx, *m.BackImage, globalConfig)
	if err != nil {
		return err
	}

	m.BackImage = img

	return nil
}

func getImage(ctx context.Context, url string, globalConfig GlobalConfig) (*string, error) {
	client := &http.Client{
		Transport: &http.Transport{ // ignore insecure certificates
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !globalConfig.GetScraperCertCheck()}},
		Timeout: imageGetTimeout,
	}

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

func getStashPerformerImage(ctx context.Context, stashURL string, performerID string, globalConfig GlobalConfig) (*string, error) {
	return getImage(ctx, stashURL+"/performer/"+performerID+"/image", globalConfig)
}

func getStashSceneImage(ctx context.Context, stashURL string, sceneID string, globalConfig GlobalConfig) (*string, error) {
	return getImage(ctx, stashURL+"/scene/"+sceneID+"/screenshot", globalConfig)
}
