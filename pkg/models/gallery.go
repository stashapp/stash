package models

import (
	"github.com/jmoiron/sqlx"
)

type GalleryReader interface {
	// Find(id int) (*Gallery, error)
	FindMany(ids []int) ([]*Gallery, error)
	// FindByChecksum(checksum string) (*Gallery, error)
	// FindByPath(path string) (*Gallery, error)
	FindBySceneID(sceneID int) (*Gallery, error)
	// ValidGalleriesForScenePath(scenePath string) ([]*Gallery, error)
	// Count() (int, error)
	All() ([]*Gallery, error)
	// Query(galleryFilter *GalleryFilterType, findFilter *FindFilterType) ([]*Gallery, int)
}

type GalleryWriter interface {
	// Create(newGallery Gallery) (*Gallery, error)
	// Update(updatedGallery Gallery) (*Gallery, error)
	// Destroy(id int) error
	// ClearGalleryId(sceneID int) error
}

type GalleryReaderWriter interface {
	GalleryReader
	GalleryWriter
}

func NewGalleryReaderWriter(tx *sqlx.Tx) GalleryReaderWriter {
	return &galleryReaderWriter{
		tx: tx,
		qb: NewGalleryQueryBuilder(),
	}
}

type galleryReaderWriter struct {
	tx *sqlx.Tx
	qb GalleryQueryBuilder
}

func (t *galleryReaderWriter) FindMany(ids []int) ([]*Gallery, error) {
	return t.qb.FindMany(ids)
}

func (t *galleryReaderWriter) All() ([]*Gallery, error) {
	return t.qb.All()
}

func (t *galleryReaderWriter) FindBySceneID(sceneID int) (*Gallery, error) {
	return t.qb.FindBySceneID(sceneID, t.tx)
}
