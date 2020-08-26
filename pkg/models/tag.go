package models

import (
	"github.com/jmoiron/sqlx"
)

type TagReader interface {
	// Find(id int) (*Tag, error)
	// FindBySceneID(sceneID int) ([]*Tag, error)
	// FindBySceneMarkerID(sceneMarkerID int) ([]*Tag, error)
	// FindByName(name string, nocase bool) (*Tag, error)
	// FindByNames(names []string, nocase bool) ([]*Tag, error)
	// Count() (int, error)
	// All() ([]*Tag, error)
	// AllSlim() ([]*Tag, error)
	// Query(tagFilter *TagFilterType, findFilter *FindFilterType) ([]*Tag, int, error)
	GetTagImage(tagID int) ([]byte, error)
}

type TagWriter interface {
	// Create(newTag Tag) (*Tag, error)
	// Update(updatedTag Tag) (*Tag, error)
	// Destroy(id string) error
	// UpdateTagImage(tagID int, image []byte) error
	// DestroyTagImage(tagID int) error
}

type TagReaderWriter interface {
	TagReader
	TagWriter
}

func NewTagReaderWriter(tx *sqlx.Tx) TagReaderWriter {
	return &tagReaderWriter{
		tx: tx,
		qb: NewTagQueryBuilder(),
	}
}

type tagReaderWriter struct {
	tx *sqlx.Tx
	qb TagQueryBuilder
}

func (t *tagReaderWriter) GetTagImage(tagID int) ([]byte, error) {
	return t.qb.GetTagImage(tagID, t.tx)
}
