package models

import (
	"github.com/jmoiron/sqlx"
)

type MovieReader interface {
	Find(id int) (*Movie, error)
	FindMany(ids []int) ([]*Movie, error)
	// FindBySceneID(sceneID int) ([]*Movie, error)
	// FindByName(name string, nocase bool) (*Movie, error)
	// FindByNames(names []string, nocase bool) ([]*Movie, error)
	All() ([]*Movie, error)
	// AllSlim() ([]*Movie, error)
	// Query(movieFilter *MovieFilterType, findFilter *FindFilterType) ([]*Movie, int)
	GetFrontImage(movieID int) ([]byte, error)
	GetBackImage(movieID int) ([]byte, error)
}

type MovieWriter interface {
	// Create(newMovie Movie) (*Movie, error)
	// Update(updatedMovie MoviePartial) (*Movie, error)
	// Destroy(id string) error
	// UpdateMovieImages(movieID int, frontImage []byte, backImage []byte) error
	// DestroyMovieImages(movieID int) error
}

type MovieReaderWriter interface {
	MovieReader
	MovieWriter
}

func NewMovieReaderWriter(tx *sqlx.Tx) MovieReaderWriter {
	return &movieReaderWriter{
		tx: tx,
		qb: NewMovieQueryBuilder(),
	}
}

type movieReaderWriter struct {
	tx *sqlx.Tx
	qb MovieQueryBuilder
}

func (t *movieReaderWriter) Find(id int) (*Movie, error) {
	return t.qb.Find(id, t.tx)
}

func (t *movieReaderWriter) FindMany(ids []int) ([]*Movie, error) {
	return t.qb.FindMany(ids)
}

func (t *movieReaderWriter) All() ([]*Movie, error) {
	return t.qb.All()
}

func (t *movieReaderWriter) GetFrontImage(movieID int) ([]byte, error) {
	return t.qb.GetFrontImage(movieID, t.tx)
}

func (t *movieReaderWriter) GetBackImage(movieID int) ([]byte, error) {
	return t.qb.GetBackImage(movieID, t.tx)
}
