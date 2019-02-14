package urlbuilders

import "strconv"

type GalleryURLBuilder struct {
	BaseURL   string
	GalleryID string
}

func NewGalleryURLBuilder(baseURL string, galleryID int) GalleryURLBuilder {
	return GalleryURLBuilder{
		BaseURL:   baseURL,
		GalleryID: strconv.Itoa(galleryID),
	}
}

func (b GalleryURLBuilder) GetGalleryImageURL(fileIndex int) string {
	return b.BaseURL + "/gallery/" + b.GalleryID + "/" + strconv.Itoa(fileIndex)
}
