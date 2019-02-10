package urlbuilders

import "strconv"

type galleryURLBuilder struct {
	BaseURL string
	GalleryID string
}

func NewGalleryURLBuilder(baseURL string, galleryID int) galleryURLBuilder {
	return galleryURLBuilder{
		BaseURL: baseURL,
		GalleryID: strconv.Itoa(galleryID),
	}
}

func (b galleryURLBuilder) GetGalleryImageUrl(fileIndex int) string {
	return b.BaseURL + "/gallery/" + b.GalleryID + "/" + strconv.Itoa(fileIndex)
}
