package urlbuilders

import (
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type ImageURLBuilder struct {
	BaseURL string
	ImageID string
}

func NewImageURLBuilder(baseURL string, image *models.Image) ImageURLBuilder {
	return ImageURLBuilder{
		BaseURL: baseURL,
		ImageID: strconv.Itoa(image.ID),
	}
}

func (b ImageURLBuilder) GetImageURL() string {
	return b.BaseURL + "/image/" + b.ImageID + "/image"
}

func (b ImageURLBuilder) GetThumbnailURL() string {
	return b.BaseURL + "/image/" + b.ImageID + "/thumbnail"
}
