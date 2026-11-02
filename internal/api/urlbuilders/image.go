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

// cacheBuster returns the value used to bust client-side caches. The image
// content is immutable for a given file, so the checksum of the primary file
// is preferred: this keeps the URL stable across metadata edits (rating,
// o-counter, tags, etc.) so the browser can reuse its cached copy. When no
// primary file is present the checksum is empty, in which case we fall back to
// the updated timestamp.
func (b ImageURLBuilder) cacheBuster() string {
	if b.Checksum != "" {
		return b.Checksum
	}
	return b.UpdatedAt
}

func (b ImageURLBuilder) GetImageURL() string {
	return b.BaseURL + "/image/" + b.ImageID + "/image?t=" + b.cacheBuster()
}

func (b ImageURLBuilder) GetThumbnailURL() string {
	return b.BaseURL + "/image/" + b.ImageID + "/thumbnail?t=" + b.cacheBuster()
}

func (b ImageURLBuilder) GetPreviewURL() string {
	if exists, err := fsutil.FileExists(manager.GetInstance().Paths.Generated.GetClipPreviewPath(b.Checksum, models.DefaultGthumbWidth)); exists && err == nil {
		return b.BaseURL + "/image/" + b.ImageID + "/preview?t=" + b.cacheBuster()
	} else {
		return ""
	}
}
