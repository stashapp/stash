package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/models"
)

type dbi interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type transaction struct {
	Ctx context.Context
	tx  *sqlx.Tx
}

func (t *transaction) Begin() error {
	if t.tx != nil {
		return errors.New("transaction already begun")
	}

	if err := database.Ready(); err != nil {
		return err
	}

	var err error
	t.tx, err = database.DB.BeginTxx(t.Ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	return nil
}

func (t *transaction) Rollback() error {
	if t.tx == nil {
		return errors.New("not in transaction")
	}

	err := t.tx.Rollback()
	if err != nil {
		return fmt.Errorf("error rolling back transaction: %v", err)
	}
	t.tx = nil

	return nil
}

func (t *transaction) Commit() error {
	if t.tx == nil {
		return errors.New("not in transaction")
	}

	err := t.tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}
	t.tx = nil

	return nil
}

func (t *transaction) Repository() models.Repository {
	return t
}

func (t *transaction) ensureTx() {
	if t.tx == nil {
		panic("tx is nil")
	}
}

func (t *transaction) Gallery() models.GalleryReaderWriter {
	t.ensureTx()
	return NewGalleryReaderWriter(t.tx)
}

func (t *transaction) Image() models.ImageReaderWriter {
	t.ensureTx()
	return NewImageReaderWriter(t.tx)
}

func (t *transaction) Movie() models.MovieReaderWriter {
	t.ensureTx()
	return NewMovieReaderWriter(t.tx)
}

func (t *transaction) Performer() models.PerformerReaderWriter {
	t.ensureTx()
	return NewPerformerReaderWriter(t.tx)
}

func (t *transaction) SceneMarker() models.SceneMarkerReaderWriter {
	t.ensureTx()
	return NewSceneMarkerReaderWriter(t.tx)
}

func (t *transaction) Scene() models.SceneReaderWriter {
	t.ensureTx()
	return NewSceneReaderWriter(t.tx)
}

func (t *transaction) ScrapedItem() models.ScrapedItemReaderWriter {
	t.ensureTx()
	return NewScrapedItemReaderWriter(t.tx)
}

func (t *transaction) Studio() models.StudioReaderWriter {
	t.ensureTx()
	return NewStudioReaderWriter(t.tx)
}

func (t *transaction) Tag() models.TagReaderWriter {
	t.ensureTx()
	return NewTagReaderWriter(t.tx)
}

func (t *transaction) SavedFilter() models.SavedFilterReaderWriter {
	t.ensureTx()
	return NewSavedFilterReaderWriter(t.tx)
}

type ReadTransaction struct{}

func (t *ReadTransaction) Begin() error {
	if err := database.Ready(); err != nil {
		return err
	}

	return nil
}

func (t *ReadTransaction) Rollback() error {
	return nil
}

func (t *ReadTransaction) Commit() error {
	return nil
}

func (t *ReadTransaction) Repository() models.ReaderRepository {
	return t
}

func (t *ReadTransaction) Gallery() models.GalleryReader {
	return NewGalleryReaderWriter(database.DB)
}

func (t *ReadTransaction) Image() models.ImageReader {
	return NewImageReaderWriter(database.DB)
}

func (t *ReadTransaction) Movie() models.MovieReader {
	return NewMovieReaderWriter(database.DB)
}

func (t *ReadTransaction) Performer() models.PerformerReader {
	return NewPerformerReaderWriter(database.DB)
}

func (t *ReadTransaction) SceneMarker() models.SceneMarkerReader {
	return NewSceneMarkerReaderWriter(database.DB)
}

func (t *ReadTransaction) Scene() models.SceneReader {
	return NewSceneReaderWriter(database.DB)
}

func (t *ReadTransaction) ScrapedItem() models.ScrapedItemReader {
	return NewScrapedItemReaderWriter(database.DB)
}

func (t *ReadTransaction) Studio() models.StudioReader {
	return NewStudioReaderWriter(database.DB)
}

func (t *ReadTransaction) Tag() models.TagReader {
	return NewTagReaderWriter(database.DB)
}

func (t *ReadTransaction) SavedFilter() models.SavedFilterReader {
	return NewSavedFilterReaderWriter(database.DB)
}

type TransactionManager struct {
}

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{}
}

func (t *TransactionManager) WithTxn(ctx context.Context, fn func(r models.Repository) error) error {
	database.WriteMu.Lock()
	defer database.WriteMu.Unlock()
	return models.WithTxn(&transaction{Ctx: ctx}, fn)
}

func (t *TransactionManager) WithReadTxn(ctx context.Context, fn func(r models.ReaderRepository) error) error {
	return models.WithROTxn(&ReadTransaction{}, fn)
}
