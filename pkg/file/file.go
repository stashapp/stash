// Package file provides functionality for managing, scanning and cleaning files and folders.
package file

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"time"

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

// ModTime returns the modification time truncated to seconds.
func ModTime(info fs.FileInfo) time.Time {
	// truncate to seconds, since we don't store beyond that in the database
	return info.ModTime().Truncate(time.Second)
}

// GetFileSize gets the size of the file, taking into account symlinks.
func GetFileSize(f models.FS, path string, info fs.FileInfo) (int64, error) {
	// #2196/#3042 - replace size with target size if file is a symlink
	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		targetInfo, err := f.Stat(path)
		if err != nil {
			return 0, fmt.Errorf("reading info for symlink %q: %w", path, err)
		}
		return targetInfo.Size(), nil
	}

	return info.Size(), nil
}
