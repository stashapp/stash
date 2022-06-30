package scraper

import (
	"context"
	"fmt"
	"github.com/stashapp/stash/pkg/logger"
	"io"
	"net/http"
	pkg_neturl "net/url"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func setPerformerImage(ctx context.Context, client *http.Client, p *models.ScrapedPerformer, globalConfig GlobalConfig) error {
	if p.Images != nil && len(p.Images) > 0 {
		for i := 0; i < len(p.Images); i++ {
			if strings.HasPrefix(p.Images[i], "http") {
				img, err := getImage(ctx, p.Images[i], client, globalConfig)
				if err != nil {
					logger.Warnf("Could not set image using URL %s: %s", p.Images[i], err.Error())
					return err
				}

				p.Images[i] = *img
				// Image is deprecated. Use images instead
			}
		}
	}

	if p.Image != nil && strings.HasPrefix(*p.Image, "http") {
		img, err := getImage(ctx, *p.Image, client, globalConfig)
		if err != nil {
			return err
		}

		p.Image = img
		// Image is deprecated. Use images instead
		p.Images = append([]string{*img}, p.Images...)
	}

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

func getImage(ctx context.Context, url string, client *http.Client, globalConfig GlobalConfig) (*string, error) {
	if strings.HasPrefix(url, "https%3A") || strings.HasPrefix(url, "http%3A") {
		urlDecoded, err := pkg_neturl.PathUnescape(url)
		if err != nil {
			return nil, err
		}
		url = urlDecoded
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

func getStashPerformerImage(ctx context.Context, stashURL string, performerID string, client *http.Client, globalConfig GlobalConfig) (*string, error) {
	return getImage(ctx, stashURL+"/performer/"+performerID+"/image", client, globalConfig)
}

func getStashSceneImage(ctx context.Context, stashURL string, sceneID string, client *http.Client, globalConfig GlobalConfig) (*string, error) {
	return getImage(ctx, stashURL+"/scene/"+sceneID+"/screenshot", client, globalConfig)
}
