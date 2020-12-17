package sqlite

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/models"
)

type transaction struct {
	Ctx context.Context
	tx  *sqlx.Tx
}

func (t *transaction) Begin() error {
	if t.tx != nil {
		return errors.New("transaction already begun")
	}

	var err error
	t.tx, err = database.DB.BeginTxx(t.Ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (t *transaction) Rollback() error {
	if t.tx == nil {
		return errors.New("not in transaction")
	}

	err := t.tx.Rollback()
	if err != nil {
		return err
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
		return err
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

type ReadTransaction struct {
	transaction
}

func (t *ReadTransaction) Repository() models.ReaderRepository {
	return t
}

func (t *ReadTransaction) Gallery() models.GalleryReader {
	t.ensureTx()
	return NewGalleryReaderWriter(t.tx)
}

func (t *ReadTransaction) Image() models.ImageReader {
	t.ensureTx()
	return NewImageReaderWriter(t.tx)
}

func (t *ReadTransaction) Movie() models.MovieReader {
	t.ensureTx()
	return NewMovieReaderWriter(t.tx)
}

func (t *ReadTransaction) Performer() models.PerformerReader {
	t.ensureTx()
	return NewPerformerReaderWriter(t.tx)
}

func (t *ReadTransaction) SceneMarker() models.SceneMarkerReader {
	t.ensureTx()
	return NewSceneMarkerReaderWriter(t.tx)
}

func (t *ReadTransaction) Scene() models.SceneReader {
	t.ensureTx()
	return NewSceneReaderWriter(t.tx)
}

func (t *ReadTransaction) ScrapedItem() models.ScrapedItemReader {
	t.ensureTx()
	return NewScrapedItemReaderWriter(t.tx)
}

func (t *ReadTransaction) Studio() models.StudioReader {
	t.ensureTx()
	return NewStudioReaderWriter(t.tx)
}

func (t *ReadTransaction) Tag() models.TagReader {
	t.ensureTx()
	return NewTagReaderWriter(t.tx)
}

type TransactionManager struct{}

func (t *TransactionManager) WithTxn(ctx context.Context, fn func(r models.Repository) error) error {
	return models.WithTxn(&transaction{Ctx: ctx}, fn)
}

func (t *TransactionManager) WithReadTxn(ctx context.Context, fn func(r models.ReaderRepository) error) error {
	return models.WithROTxn(&ReadTransaction{
		transaction{
			Ctx: ctx,
		},
	}, fn)
}
