package sqlite

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type key int

const (
	txnKey key = iota + 1
	dbKey
	hookManagerKey
)

func (db *Database) WithDatabase(ctx context.Context) (context.Context, error) {
	// if we are already in a transaction or have a database already, just use it
	if tx, _ := getDBReader(ctx); tx != nil {
		return ctx, nil
	}

	return context.WithValue(ctx, dbKey, db.db), nil
}

func (db *Database) Begin(ctx context.Context) (context.Context, error) {
	if tx, _ := getTx(ctx); tx != nil {
		// log the stack trace so we can see
		logger.Error(string(debug.Stack()))

		return nil, fmt.Errorf("already in transaction")
	}

	tx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}

	hookMgr := &hookManager{}
	ctx = hookMgr.register(ctx)

	return context.WithValue(ctx, txnKey, tx), nil
}

func (db *Database) Commit(ctx context.Context) error {
	tx, err := getTx(ctx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// execute post-commit hooks
	db.executePostCommitHooks(ctx)

	return nil
}

func (db *Database) Rollback(ctx context.Context) error {
	tx, err := getTx(ctx)
	if err != nil {
		return err
	}

	if err := tx.Rollback(); err != nil {
		return err
	}

	// execute post-rollback hooks
	db.executePostRollbackHooks(ctx)

	return nil
}

func getTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, ok := ctx.Value(txnKey).(*sqlx.Tx)
	if !ok || tx == nil {
		return nil, fmt.Errorf("not in transaction")
	}
	return tx, nil
}

func getDBReader(ctx context.Context) (dbReader, error) {
	// get transaction first if present
	tx, ok := ctx.Value(txnKey).(*sqlx.Tx)
	if !ok || tx == nil {
		// try to get database if present
		db, ok := ctx.Value(dbKey).(*sqlx.DB)
		if !ok || db == nil {
			return nil, fmt.Errorf("not in transaction")
		}
		return db, nil
	}
	return tx, nil
}

func (db *Database) TxnRepository() models.Repository {
	return models.Repository{
		TxnManager:  db,
		File:        db.File,
		Folder:      db.Folder,
		Gallery:     db.Gallery,
		Image:       db.Image,
		Movie:       MovieReaderWriter,
		Performer:   PerformerReaderWriter,
		Scene:       db.Scene,
		SceneMarker: SceneMarkerReaderWriter,
		ScrapedItem: ScrapedItemReaderWriter,
		Studio:      StudioReaderWriter,
		Tag:         TagReaderWriter,
		SavedFilter: SavedFilterReaderWriter,
	}
}
