package models

import (
	"time"
)

type GalleryChapter struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	ImageIndex int       `json:"image_index"`
	GalleryID  int       `json:"gallery_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type GalleryChapters []*GalleryChapter

func (m *GalleryChapters) Append(o interface{}) {
	*m = append(*m, o.(*GalleryChapter))
}

func (m *GalleryChapters) New() interface{} {
	return &GalleryChapter{}
}
