package models

import (
	"github.com/jmoiron/sqlx"
)

type PerformerReader interface {
	// Find(id int) (*Performer, error)
	FindMany(ids []int) ([]*Performer, error)
	FindBySceneID(sceneID int) ([]*Performer, error)
	FindNamesBySceneID(sceneID int) ([]*Performer, error)
	FindByNames(names []string, nocase bool) ([]*Performer, error)
	// Count() (int, error)
	All() ([]*Performer, error)
	// AllSlim() ([]*Performer, error)
	// Query(performerFilter *PerformerFilterType, findFilter *FindFilterType) ([]*Performer, int)
	GetPerformerImage(performerID int) ([]byte, error)
}

type PerformerWriter interface {
	Create(newPerformer Performer) (*Performer, error)
	Update(updatedPerformer Performer) (*Performer, error)
	// Destroy(id string) error
	UpdatePerformerImage(performerID int, image []byte) error
	// DestroyPerformerImage(performerID int) error
}

type PerformerReaderWriter interface {
	PerformerReader
	PerformerWriter
}

func NewPerformerReaderWriter(tx *sqlx.Tx) PerformerReaderWriter {
	return &performerReaderWriter{
		tx: tx,
		qb: NewPerformerQueryBuilder(),
	}
}

type performerReaderWriter struct {
	tx *sqlx.Tx
	qb PerformerQueryBuilder
}

func (t *performerReaderWriter) FindMany(ids []int) ([]*Performer, error) {
	return t.qb.FindMany(ids)
}

func (t *performerReaderWriter) FindByNames(names []string, nocase bool) ([]*Performer, error) {
	return t.qb.FindByNames(names, t.tx, nocase)
}

func (t *performerReaderWriter) All() ([]*Performer, error) {
	return t.qb.All()
}

func (t *performerReaderWriter) GetPerformerImage(performerID int) ([]byte, error) {
	return t.qb.GetPerformerImage(performerID, t.tx)
}

func (t *performerReaderWriter) FindBySceneID(id int) ([]*Performer, error) {
	return t.qb.FindBySceneID(id, t.tx)
}

func (t *performerReaderWriter) FindNamesBySceneID(sceneID int) ([]*Performer, error) {
	return t.qb.FindNameBySceneID(sceneID, t.tx)
}

func (t *performerReaderWriter) Create(newPerformer Performer) (*Performer, error) {
	return t.qb.Create(newPerformer, t.tx)
}

func (t *performerReaderWriter) Update(updatedPerformer Performer) (*Performer, error) {
	return t.qb.Update(updatedPerformer, t.tx)
}

func (t *performerReaderWriter) UpdatePerformerImage(performerID int, image []byte) error {
	return t.qb.UpdatePerformerImage(performerID, image, t.tx)
}
