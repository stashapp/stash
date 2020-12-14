package sqlite

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/models"
)

type Transaction struct {
	Ctx context.Context
	tx  *sqlx.Tx
}

func (t *Transaction) Begin() error {
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

func (t *Transaction) Rollback() error {
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

func (t *Transaction) Commit() error {
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

func (t *Transaction) Repository() models.Repository {
	return t
}

func (t *Transaction) ensureTx() {
	if t.tx == nil {
		panic("tx is nil")
	}
}

func (t *Transaction) Gallery() models.GalleryReaderWriter {
	t.ensureTx()
	return NewGalleryReaderWriter(t.tx)
}

func (t *Transaction) Image() models.ImageReaderWriter {
	t.ensureTx()
	return NewImageReaderWriter(t.tx)
}

func (t *Transaction) Join() models.JoinReaderWriter {
	t.ensureTx()
	return NewJoinReaderWriter(t.tx)
}

func (t *Transaction) Movie() models.MovieReaderWriter {
	t.ensureTx()
	return NewMovieReaderWriter(t.tx)
}

func (t *Transaction) Performer() models.PerformerReaderWriter {
	t.ensureTx()
	return NewPerformerReaderWriter(t.tx)
}

func (t *Transaction) SceneMarker() models.SceneMarkerReaderWriter {
	t.ensureTx()
	return NewSceneMarkerReaderWriter(t.tx)
}

func (t *Transaction) Scene() models.SceneReaderWriter {
	t.ensureTx()
	return NewSceneReaderWriter(t.tx)
}

func (t *Transaction) Studio() models.StudioReaderWriter {
	t.ensureTx()
	return NewStudioReaderWriter(t.tx)
}

func (t *Transaction) Tag() models.TagReaderWriter {
	t.ensureTx()
	return NewTagReaderWriter(t.tx)
}
