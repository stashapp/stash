package file

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

const deleteFileSuffix = ".delete"

// RenamerRemover provides access to the Rename and Remove functions.
type RenamerRemover interface {
	Renamer
	Remove(name string) error
	RemoveAll(path string) error
	Statter
}

type renamerRemoverImpl struct {
	RenameFn    func(oldpath, newpath string) error
	RemoveFn    func(name string) error
	RemoveAllFn func(path string) error
	StatFn      func(path string) (fs.FileInfo, error)
}

func (r renamerRemoverImpl) Rename(oldpath, newpath string) error {
	return r.RenameFn(oldpath, newpath)
}

func (r renamerRemoverImpl) Remove(name string) error {
	return r.RemoveFn(name)
}

func (r renamerRemoverImpl) RemoveAll(path string) error {
	return r.RemoveAllFn(path)
}

func (r renamerRemoverImpl) Stat(path string) (fs.FileInfo, error) {
	return r.StatFn(path)
}

func newRenamerRemoverImpl() renamerRemoverImpl {
	return renamerRemoverImpl{
		// use fsutil.SafeMove to support cross-device moves
		RenameFn:    fsutil.SafeMove,
		RemoveFn:    os.Remove,
		RemoveAllFn: os.RemoveAll,
		StatFn:      os.Stat,
	}
}

// Deleter is used to safely delete files and directories from the filesystem.
// During a transaction, files and directories are marked for deletion using
// the Files and Dirs methods. This will rename the files/directories to be
// deleted. If the transaction is rolled back, then the files/directories can
// be restored to their original state with the Abort method. If the
// transaction is committed, the marked files are then deleted from the
// filesystem using the Complete method.
type Deleter struct {
	RenamerRemover RenamerRemover
	files          []string
	dirs           []string
}

func NewDeleter() *Deleter {
	return &Deleter{
		RenamerRemover: newRenamerRemoverImpl(),
	}
}

// RegisterHooks registers post-commit and post-rollback hooks.
func (d *Deleter) RegisterHooks(ctx context.Context) {
	txn.AddPostCommitHook(ctx, func(ctx context.Context) {
		d.Commit()
	})

	txn.AddPostRollbackHook(ctx, func(ctx context.Context) {
		d.Rollback()
	})
}

// Files designates files to be deleted. Each file marked will be renamed to add
// a `.delete` suffix. An error is returned if a file could not be renamed.
// Note that if an error is returned, then some files may be left renamed.
// Abort should be called to restore marked files if this function returns an
// error.
func (d *Deleter) Files(paths []string) error {
	for _, p := range paths {
		// fail silently if the file does not exist
		if _, err := d.RenamerRemover.Stat(p); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				logger.Warnf("File %q does not exist and therefore cannot be deleted. Ignoring.", p)
				continue
			}

			return fmt.Errorf("check file %q exists: %w", p, err)
		}

		if err := d.renameForDelete(p); err != nil {
			return fmt.Errorf("marking file %q for deletion: %w", p, err)
		}
		d.files = append(d.files, p)
	}

	return nil
}

// Dirs designates directories to be deleted. Each directory marked will be renamed to add
// a `.delete` suffix. An error is returned if a directory could not be renamed.
// Note that if an error is returned, then some directories may be left renamed.
// Abort should be called to restore marked files/directories if this function returns an
// error.
func (d *Deleter) Dirs(paths []string) error {
	for _, p := range paths {
		// fail silently if the file does not exist
		if _, err := d.RenamerRemover.Stat(p); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				logger.Warnf("Directory %q does not exist and therefore cannot be deleted. Ignoring.", p)
				continue
			}

			return fmt.Errorf("check directory %q exists: %w", p, err)
		}

		if err := d.renameForDelete(p); err != nil {
			return fmt.Errorf("marking directory %q for deletion: %w", p, err)
		}
		d.dirs = append(d.dirs, p)
	}

	return nil
}

// Rollback tries to rename all marked files and directories back to their
// original names and clears the marked list. Any errors encountered are
// logged. All files will be attempted regardless of any errors occurred.
func (d *Deleter) Rollback() {
	for _, f := range append(d.files, d.dirs...) {
		if err := d.renameForRestore(f); err != nil {
			logger.Warnf("Error restoring %q: %v", f, err)
		}
	}

	d.files = nil
	d.dirs = nil
}

// Commit deletes all files marked for deletion and clears the marked list.
// Any errors encountered are logged. All files will be attempted, regardless
// of the errors encountered.
func (d *Deleter) Commit() {
	for _, f := range d.files {
		if err := d.RenamerRemover.Remove(f + deleteFileSuffix); err != nil {
			logger.Warnf("Error deleting file %q: %v", f+deleteFileSuffix, err)
		}
	}

	for _, f := range d.dirs {
		if err := d.RenamerRemover.RemoveAll(f + deleteFileSuffix); err != nil {
			logger.Warnf("Error deleting directory %q: %v", f+deleteFileSuffix, err)
		}
	}

	d.files = nil
	d.dirs = nil
}

func (d *Deleter) renameForDelete(path string) error {
	return d.RenamerRemover.Rename(path, path+deleteFileSuffix)
}

func (d *Deleter) renameForRestore(path string) error {
	return d.RenamerRemover.Rename(path+deleteFileSuffix, path)
}

func Destroy(ctx context.Context, destroyer models.FileDestroyer, f models.File, fileDeleter *Deleter, deleteFile bool) error {
	if err := destroyer.Destroy(ctx, f.Base().ID); err != nil {
		return err
	}

	// don't delete files in zip files
	if deleteFile && f.Base().ZipFileID == nil {
		if err := fileDeleter.Files([]string{f.Base().Path}); err != nil {
			return err
		}
	}

	return nil
}

type ZipDestroyer struct {
	FileDestroyer   models.FileFinderDestroyer
	FolderDestroyer models.FolderFinderDestroyer
}

func (d *ZipDestroyer) DestroyZip(ctx context.Context, f models.File, fileDeleter *Deleter, deleteFile bool) error {
	// destroy contained files
	files, err := d.FileDestroyer.FindByZipFileID(ctx, f.Base().ID)
	if err != nil {
		return err
	}

	for _, ff := range files {
		if err := d.FileDestroyer.Destroy(ctx, ff.Base().ID); err != nil {
			return err
		}
	}

	// destroy contained folders
	folders, err := d.FolderDestroyer.FindByZipFileID(ctx, f.Base().ID)
	if err != nil {
		return err
	}

	for _, ff := range folders {
		if err := d.FolderDestroyer.Destroy(ctx, ff.ID); err != nil {
			return err
		}
	}

	if err := d.FileDestroyer.Destroy(ctx, f.Base().ID); err != nil {
		return err
	}

	if deleteFile {
		if err := fileDeleter.Files([]string{f.Base().Path}); err != nil {
			return err
		}
	}

	return nil
}
