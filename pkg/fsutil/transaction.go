package fsutil

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/logger"
)

const (
	deleteFileSuffix  = ".delete"
	changedFileSuffix = ".orig"
)

type Writer interface {
	WriteFile(path string, file []byte) error
}

type Deleter interface {
	DeleteFiles(paths []string) error
}

type WriterDeleter interface {
	Writer
	Deleter
}

// TxnFS provides access to the filesystem operations that are used by the FSTransaction object.
type TxnFS interface {
	Stat(name string) (fs.FileInfo, error)

	WriteFile(name string, data []byte, perm fs.FileMode) error

	Rename(oldpath, newpath string) error
	Remove(name string) error
	RemoveAll(path string) error
}

type txnFSImpl struct {
	RenameFn    func(oldpath, newpath string) error
	RemoveFn    func(name string) error
	RemoveAllFn func(path string) error
	StatFn      func(path string) (fs.FileInfo, error)
	WriteFileFn func(name string, data []byte, perm fs.FileMode) error
}

func (r txnFSImpl) Rename(oldpath, newpath string) error {
	return r.RenameFn(oldpath, newpath)
}

func (r txnFSImpl) Remove(name string) error {
	return r.RemoveFn(name)
}

func (r txnFSImpl) RemoveAll(path string) error {
	return r.RemoveAllFn(path)
}

func (r txnFSImpl) Stat(path string) (fs.FileInfo, error) {
	return r.StatFn(path)
}

func (r txnFSImpl) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return r.WriteFileFn(name, data, perm)
}

// FSTransaction is used to safely make changes to the filesystem within a transaction.
// Changes are made to the filesystem and can then be committed or rolled back.
// NOTE: does not currently handle created directories.
type FSTransaction struct {
	TxnFS TxnFS

	removedFiles []string
	removedDirs  []string
	addedFiles   []string
	changedFiles []string
}

func NewFSTransaction() *FSTransaction {
	return &FSTransaction{
		TxnFS: txnFSImpl{
			RenameFn:    os.Rename,
			RemoveFn:    os.Remove,
			RemoveAllFn: os.RemoveAll,
			StatFn:      os.Stat,
			WriteFileFn: os.WriteFile,
		},
	}
}

func (d *FSTransaction) WriteFile(path string, file []byte) error {
	// check if the file already exists
	_, err := d.TxnFS.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		// file does not exist, so we can just write it
		pathErr := EnsureDirAll(filepath.Dir(path))
		if pathErr != nil {
			return fmt.Errorf("cannot ensure path %s", pathErr)
		}

		err = d.TxnFS.WriteFile(path, file, 0755)
		if err != nil {
			return err
		}
		d.addedFiles = append(d.addedFiles, path)
		return nil
	}

	// file already exists, so we need to write it to a temp file
	if err = d.renameForChanged(path); err != nil {
		return fmt.Errorf("renaming file %q: %w", path, err)
	}

	d.changedFiles = append(d.changedFiles, path)
	return d.TxnFS.WriteFile(path, file, 0755)
}

// Files designates files to be deleted. Each file marked will be renamed to add
// a `.delete` suffix. An error is returned if a file could not be renamed.
// Note that if an error is returned, then some files may be left renamed.
// Abort should be called to restore marked files if this function returns an
// error.
func (d *FSTransaction) DeleteFiles(paths []string) error {
	for _, p := range paths {
		// fail silently if the file does not exist
		if _, err := d.TxnFS.Stat(p); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				logger.Warnf("File %q does not exist and therefore cannot be deleted. Ignoring.", p)
				continue
			}

			return fmt.Errorf("check file %q exists: %w", p, err)
		}

		if err := d.renameForDelete(p); err != nil {
			return fmt.Errorf("marking file %q for deletion: %w", p, err)
		}
		d.removedFiles = append(d.removedFiles, p)
	}

	return nil
}

// Dirs designates directories to be deleted. Each directory marked will be renamed to add
// a `.delete` suffix. An error is returned if a directory could not be renamed.
// Note that if an error is returned, then some directories may be left renamed.
// Abort should be called to restore marked files/directories if this function returns an
// error.
func (d *FSTransaction) DeleteDirs(paths []string) error {
	for _, p := range paths {
		// fail silently if the file does not exist
		if _, err := d.TxnFS.Stat(p); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				logger.Warnf("Directory %q does not exist and therefore cannot be deleted. Ignoring.", p)
				continue
			}

			return fmt.Errorf("check directory %q exists: %w", p, err)
		}

		if err := d.renameForDelete(p); err != nil {
			return fmt.Errorf("marking directory %q for deletion: %w", p, err)
		}
		d.removedDirs = append(d.removedDirs, p)
	}

	return nil
}

// Rollback tries to rename all marked files and directories back to their
// original names and clears the marked list. Any errors encountered are
// logged. All files will be attempted regardless of any errors occurred.
func (d *FSTransaction) Rollback() {
	for _, f := range d.changedFiles {
		if err := d.renameForRestoreChanged(f); err != nil {
			logger.Warnf("Error restoring %q: %v", f, err)
		}
	}

	for _, f := range append(d.removedFiles, d.removedDirs...) {
		if err := d.renameForRestoreDelete(f); err != nil {
			logger.Warnf("Error restoring %q: %v", f, err)
		}
	}

	d.reset()
}

// Commit deletes all files marked for deletion and clears the marked list.
// Any errors encountered are logged. All files will be attempted, regardless
// of the errors encountered.
func (d *FSTransaction) Commit() {
	for _, f := range d.changedFiles {
		if err := d.TxnFS.Remove(f + changedFileSuffix); err != nil {
			logger.Warnf("Error deleting file %q: %v", f+changedFileSuffix, err)
		}
	}

	for _, f := range d.removedFiles {
		if err := d.TxnFS.Remove(f + deleteFileSuffix); err != nil {
			logger.Warnf("Error deleting file %q: %v", f+deleteFileSuffix, err)
		}
	}

	for _, f := range d.removedDirs {
		if err := d.TxnFS.RemoveAll(f + deleteFileSuffix); err != nil {
			logger.Warnf("Error deleting directory %q: %v", f+deleteFileSuffix, err)
		}
	}

	d.reset()
}

func (d *FSTransaction) reset() {
	d.removedFiles = nil
	d.removedDirs = nil
	d.changedFiles = nil
	d.addedFiles = nil
}

func (d *FSTransaction) renameForDelete(path string) error {
	return d.TxnFS.Rename(path, path+deleteFileSuffix)
}

func (d *FSTransaction) renameForRestoreDelete(path string) error {
	return d.TxnFS.Rename(path+deleteFileSuffix, path)
}

func (d *FSTransaction) renameForChanged(path string) error {
	return d.TxnFS.Rename(path, path+changedFileSuffix)
}

func (d *FSTransaction) renameForRestoreChanged(path string) error {
	return d.TxnFS.Rename(path+changedFileSuffix, path)
}
