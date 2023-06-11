package urlbuilders

import (
	"strconv"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models"
)

type ImageURLBuilder struct {
	BaseURL   string
	ImageID   string
	Checksum  string
	UpdatedAt string
}

func NewImageURLBuilder(baseURL string, image *models.Image) ImageURLBuilder {
	return ImageURLBuilder{
		BaseURL:   baseURL,
		ImageID:   strconv.Itoa(image.ID),
		Checksum:  image.Checksum,
		UpdatedAt: strconv.FormatInt(image.UpdatedAt.Unix(), 10),
	}
}

func (b ImageURLBuilder) GetImageURL() string {
	return b.BaseURL + "/image/" + b.ImageID + "/image?t=" + b.UpdatedAt
}

func (b ImageURLBuilder) GetThumbnailURL() string {
	return b.BaseURL + "/image/" + b.ImageID + "/thumbnail?t=" + b.UpdatedAt
}

func (b ImageURLBuilder) GetPreviewURL() string {
	if exists, err := fsutil.FileExists(manager.GetInstance().Paths.Generated.GetClipPreviewPath(b.Checksum, models.DefaultGthumbWidth)); exists && err == nil {
		return b.BaseURL + "/image/" + b.ImageID + "/preview?" + b.UpdatedAt
	} else {
		return ""
	}
}
