package models

import (
	"github.com/jmoiron/sqlx"
)

type GalleryReader interface {
	// Find(id int) (*Gallery, error)
	FindMany(ids []int) ([]*Gallery, error)
	FindByChecksum(checksum string) (*Gallery, error)
	FindByPath(path string) (*Gallery, error)
	FindBySceneID(sceneID int) (*Gallery, error)
	// ValidGalleriesForScenePath(scenePath string) ([]*Gallery, error)
	// Count() (int, error)
	All() ([]*Gallery, error)
	// Query(galleryFilter *GalleryFilterType, findFilter *FindFilterType) ([]*Gallery, int)
}

type GalleryWriter interface {
	Create(newGallery Gallery) (*Gallery, error)
	Update(updatedGallery Gallery) (*Gallery, error)
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

func (t *galleryReaderWriter) FindByChecksum(checksum string) (*Gallery, error) {
	return t.qb.FindByChecksum(checksum, t.tx)
}

func (t *galleryReaderWriter) All() ([]*Gallery, error) {
	return t.qb.All()
}

func (t *galleryReaderWriter) FindByPath(path string) (*Gallery, error) {
	return t.qb.FindByPath(path)
}

func (t *galleryReaderWriter) FindBySceneID(sceneID int) (*Gallery, error) {
	return t.qb.FindBySceneID(sceneID, t.tx)
}

func (t *galleryReaderWriter) Create(newGallery Gallery) (*Gallery, error) {
	return t.qb.Create(newGallery, t.tx)
}

func (t *galleryReaderWriter) Update(updatedGallery Gallery) (*Gallery, error) {
	return t.qb.Update(updatedGallery, t.tx)
}
