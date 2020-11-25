package models

import (
	"github.com/jmoiron/sqlx"
)

type ImageReader interface {
	// Find(id int) (*Image, error)
	FindMany(ids []int) ([]*Image, error)
	FindByChecksum(checksum string) (*Image, error)
	FindByGalleryID(galleryID int) ([]*Image, error)
	// FindByPath(path string) (*Image, error)
	// FindByPerformerID(performerID int) ([]*Image, error)
	// CountByPerformerID(performerID int) (int, error)
	// FindByStudioID(studioID int) ([]*Image, error)
	// Count() (int, error)
	// SizeCount() (string, error)
	// CountByStudioID(studioID int) (int, error)
	// CountByTagID(tagID int) (int, error)
	All() ([]*Image, error)
	// Query(imageFilter *ImageFilterType, findFilter *FindFilterType) ([]*Image, int)
}

type ImageWriter interface {
	Create(newImage Image) (*Image, error)
	Update(updatedImage ImagePartial) (*Image, error)
	UpdateFull(updatedImage Image) (*Image, error)
	// IncrementOCounter(id int) (int, error)
	// DecrementOCounter(id int) (int, error)
	// ResetOCounter(id int) (int, error)
	// Destroy(id string) error
}

type ImageReaderWriter interface {
	ImageReader
	ImageWriter
}

func NewImageReaderWriter(tx *sqlx.Tx) ImageReaderWriter {
	return &imageReaderWriter{
		tx: tx,
		qb: NewImageQueryBuilder(),
	}
}

type imageReaderWriter struct {
	tx *sqlx.Tx
	qb ImageQueryBuilder
}

func (t *imageReaderWriter) FindMany(ids []int) ([]*Image, error) {
	return t.qb.FindMany(ids)
}

func (t *imageReaderWriter) FindByChecksum(checksum string) (*Image, error) {
	return t.qb.FindByChecksum(checksum)
}

func (t *imageReaderWriter) FindByGalleryID(galleryID int) ([]*Image, error) {
	return t.qb.FindByGalleryID(galleryID)
}

func (t *imageReaderWriter) All() ([]*Image, error) {
	return t.qb.All()
}

func (t *imageReaderWriter) Create(newImage Image) (*Image, error) {
	return t.qb.Create(newImage, t.tx)
}

func (t *imageReaderWriter) Update(updatedImage ImagePartial) (*Image, error) {
	return t.qb.Update(updatedImage, t.tx)
}

func (t *imageReaderWriter) UpdateFull(updatedImage Image) (*Image, error) {
	return t.qb.UpdateFull(updatedImage, t.tx)
}
