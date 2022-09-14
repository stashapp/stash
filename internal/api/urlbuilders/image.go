package urlbuilders

import (
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type ImageURLBuilder struct {
	BaseURL   string
	ImageID   string
	UpdatedAt string
}

func NewImageURLBuilder(baseURL string, image *models.Image) ImageURLBuilder {
	return ImageURLBuilder{
		BaseURL:   baseURL,
		ImageID:   strconv.Itoa(image.ID),
		UpdatedAt: strconv.FormatInt(image.UpdatedAt.Unix(), 10),
	}
}

func (b ImageURLBuilder) GetImageURL() string {
	return b.BaseURL + "/image/" + b.ImageID + "/image?" + b.UpdatedAt
}

func (b ImageURLBuilder) GetThumbnailURL() string {
	return b.BaseURL + "/image/" + b.ImageID + "/thumbnail?" + b.UpdatedAt
}
