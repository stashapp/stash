package scraper

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
)

const testQueryURL = "http://host/api?file={filename}&path={path}"

func TestQueryURLParametersPath(t *testing.T) {
	tests := []struct {
		name   string
		params queryURLParameters
		want   string
	}{
		{
			"scene",
			queryURLParametersFromScene(&models.Scene{
				Path:  "/media/scenes/My Scene.mp4",
				Files: models.NewRelatedVideoFiles([]*models.VideoFile{{BaseFile: &models.BaseFile{}}}),
				URLs:  models.NewRelatedStrings([]string{}),
			}),
			"http://host/api?file=My Scene.mp4&path=/media/scenes/My Scene.mp4",
		},
		{
			"gallery",
			queryURLParametersFromGallery(&models.Gallery{
				Path:  "/media/galleries/My Gallery.zip",
				Files: models.NewRelatedFiles([]models.File{}),
				URLs:  models.NewRelatedStrings([]string{}),
			}),
			"http://host/api?file=My Gallery.zip&path=/media/galleries/My Gallery.zip",
		},
		{
			"image",
			queryURLParametersFromImage(&models.Image{
				Path: "/media/images/My Image.jpg",
				URLs: models.NewRelatedStrings([]string{}),
			}),
			"http://host/api?file=My Image.jpg&path=/media/images/My Image.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.constructURL(testQueryURL); got != tt.want {
				t.Errorf("constructURL() = %q, want %q", got, tt.want)
			}
		})
	}
}
