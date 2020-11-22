package urlbuilders

import (
	"strconv"
)

type ImageURLBuilder struct {
	BaseURL string
	ImageID string
}

func NewImageURLBuilder(baseURL string, imageID int) ImageURLBuilder {
	return ImageURLBuilder{
		BaseURL: baseURL,
		ImageID: strconv.Itoa(imageID),
	}
}

func (b ImageURLBuilder) GetImageURL() string {
	return b.BaseURL + "/image/" + b.ImageID + "/image"
}

func (b ImageURLBuilder) GetThumbnailURL() string {
	return b.BaseURL + "/image/" + b.ImageID + "/thumbnail"
}
