package sqlite

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type key int

const (
	txnKey key = iota + 1
	dbKey
	exclusiveKey
)

func (db *Database) WithDatabase(ctx context.Context) (context.Context, error) {
	// if we are already in a transaction or have a database already, just use it
	if tx, _ := getDBReader(ctx); tx != nil {
		return ctx, nil
	}

	return context.WithValue(ctx, dbKey, db.db), nil
}

func (db *Database) Begin(ctx context.Context, exclusive bool) (context.Context, error) {
	if tx, _ := getTx(ctx); tx != nil {
		// log the stack trace so we can see
		logger.Error(string(debug.Stack()))

		return nil, fmt.Errorf("already in transaction")
	}

	if exclusive {
		if err := db.lock(ctx); err != nil {
			return nil, err
		}
	}

	tx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		// begin failed, unlock
		if exclusive {
			db.unlock()
		}
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}

	ctx = context.WithValue(ctx, exclusiveKey, exclusive)

	return context.WithValue(ctx, txnKey, tx), nil
}

func (db *Database) Commit(ctx context.Context) error {
	tx, err := getTx(ctx)
	if err != nil {
		return err
	}

	defer db.txnComplete(ctx)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *Database) Rollback(ctx context.Context) error {
	tx, err := getTx(ctx)
	if err != nil {
		return err
	}

	defer db.txnComplete(ctx)

	if err := tx.Rollback(); err != nil {
		return err
	}

	return nil
}

func (db *Database) txnComplete(ctx context.Context) {
	if exclusive := ctx.Value(exclusiveKey).(bool); exclusive {
		db.unlock()
	}
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

func (db *Database) IsLocked(err error) bool {
	var sqliteError sqlite3.Error
	if errors.As(err, &sqliteError) {
		return sqliteError.Code == sqlite3.ErrBusy
	}
	return false
}

func (db *Database) Repository() models.Repository {
	return models.Repository{
		TxnManager:     db,
		Blob:           db.Blobs,
		File:           db.File,
		Folder:         db.Folder,
		Gallery:        db.Gallery,
		GalleryChapter: db.GalleryChapter,
		Image:          db.Image,
		Group:          db.Group,
		Performer:      db.Performer,
		Scene:          db.Scene,
		SceneMarker:    db.SceneMarker,
		Studio:         db.Studio,
		Tag:            db.Tag,
		SavedFilter:    db.SavedFilter,
	}
}
