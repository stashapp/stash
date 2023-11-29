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

func NewGalleryChapter() GalleryChapter {
	currentTime := time.Now()
	return GalleryChapter{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}

// GalleryChapterPartial represents part of a GalleryChapter object.
// It is used to update the database entry.
type GalleryChapterPartial struct {
	Title      OptionalString
	ImageIndex OptionalInt
	GalleryID  OptionalInt
	CreatedAt  OptionalTime
	UpdatedAt  OptionalTime
}

func NewGalleryChapterPartial() GalleryChapterPartial {
	currentTime := time.Now()
	return GalleryChapterPartial{
		UpdatedAt: NewOptionalTime(currentTime),
	}
}
