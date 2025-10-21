package urlbuilders

import (
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type GalleryURLBuilder struct {
	BaseURL   string
	GalleryID string
	UpdatedAt string
}

func NewGalleryURLBuilder(baseURL string, gallery *models.Gallery) GalleryURLBuilder {
	return GalleryURLBuilder{
		BaseURL:   baseURL,
		GalleryID: strconv.Itoa(gallery.ID),
		UpdatedAt: strconv.FormatInt(gallery.UpdatedAt.Unix(), 10),
	}
}

func (b GalleryURLBuilder) GetPreviewURL() string {
	return b.BaseURL + "/gallery/" + b.GalleryID + "/preview"
}

func (b GalleryURLBuilder) GetCoverURL() string {
	return b.BaseURL + "/gallery/" + b.GalleryID + "/cover?t=" + b.UpdatedAt
}
