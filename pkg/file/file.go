// Package file provides functionality for managing, scanning and cleaning files and folders.
package file

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

// Repository provides access to storage methods for files and folders.
type Repository struct {
	TxnManager models.TxnManager

	File   models.FileReaderWriter
	Folder models.FolderReaderWriter
}

func NewRepository(repo models.Repository) Repository {
	return Repository{
		TxnManager: repo.TxnManager,
		File:       repo.File,
		Folder:     repo.Folder,
	}
}

func (r *Repository) WithTxn(ctx context.Context, fn txn.TxnFunc) error {
	return txn.WithTxn(ctx, r.TxnManager, fn)
}

func (r *Repository) WithReadTxn(ctx context.Context, fn txn.TxnFunc) error {
	return txn.WithReadTxn(ctx, r.TxnManager, fn)
}

func (r *Repository) WithDB(ctx context.Context, fn txn.TxnFunc) error {
	return txn.WithDatabase(ctx, r.TxnManager, fn)
}
