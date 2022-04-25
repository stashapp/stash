package sqlite

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

type key int

const (
	txnKey key = iota + 1
)

func (db *Database) Begin(ctx context.Context) (context.Context, error) {
	if tx, _ := getTx(ctx); tx != nil {
		return nil, fmt.Errorf("already in transaction")
	}

	tx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}

	return context.WithValue(ctx, txnKey, tx), nil
}

func (db *Database) Commit(ctx context.Context) error {
	tx, err := getTx(ctx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *Database) Rollback(ctx context.Context) error {
	tx, err := getTx(ctx)
	if err != nil {
		return err
	}
	return tx.Rollback()
}

func getTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, ok := ctx.Value(txnKey).(*sqlx.Tx)
	if !ok || tx == nil {
		return nil, fmt.Errorf("not in transaction")
	}
	return tx, nil
}

func (db *Database) TxnRepository() models.Repository {
	return models.Repository{
		TxnManager:  db,
		Gallery:     GalleryReaderWriter,
		Image:       ImageReaderWriter,
		Movie:       MovieReaderWriter,
		Performer:   PerformerReaderWriter,
		Scene:       SceneReaderWriter,
		SceneMarker: SceneMarkerReaderWriter,
		ScrapedItem: ScrapedItemReaderWriter,
		Studio:      StudioReaderWriter,
		Tag:         TagReaderWriter,
		SavedFilter: SavedFilterReaderWriter,
	}
}
