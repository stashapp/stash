package file

import (
	"fmt"
	"os"
)

const deleteFileSuffix = ".delete"

// RenamerRemover provides access to the Rename and Remove functions.
type RenamerRemover interface {
	Rename(oldpath, newpath string) error
	Remove(name string) error
	RemoveAll(path string) error
}

type renamerRemoverImpl struct {
	RenameFn    func(oldpath, newpath string) error
	RemoveFn    func(name string) error
	RemoveAllFn func(path string) error
}

func (r renamerRemoverImpl) Rename(oldpath, newpath string) error {
	return r.RenameFn(oldpath, newpath)
}

func (r renamerRemoverImpl) Remove(name string) error {
	return r.RemoveFn(name)
}

func (r renamerRemoverImpl) RemoveAll(name string) error {
	return r.RemoveAllFn(name)
}

// Deleter is used to safely delete files from the filesystem.
// During a transaction, files are marked for deletion using the Mark method.
// This will rename the files to be deleted. If the transaction is rolled back,
// then the files are restored to their original state. If the transaction is
// committed, the marked files are then deleted from the filesystem.
type Deleter struct {
	RenamerRemover RenamerRemover
	files          []string
	dirs           []string
}

func NewDeleter() *Deleter {
	return &Deleter{
		RenamerRemover: renamerRemoverImpl{
			RenameFn:    os.Rename,
			RemoveFn:    os.Remove,
			RemoveAllFn: os.RemoveAll,
		},
	}
}

// Files designates files to be deleted. Each file marked will be renamed to add
// a `.delete` suffix. An error is returned if a file could not be renamed.
// Note that if an error is returned, then some files may be left renamed.
// Abort should be called to restore marked files if this function returns an
// error.
func (d *Deleter) Files(paths []string) error {
	for _, p := range paths {
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
		if err := d.renameForDelete(p); err != nil {
			return fmt.Errorf("marking directory %q for deletion: %w", p, err)
		}
		d.dirs = append(d.dirs, p)
	}

	return nil
}

// Abort tries to rename all marked files and directories back to their
// original names. It returns a slice of errors for each file/directory
// that fails. All files will be attempted regardless of any errors occurred.
// All files and directories will be unmarked.
func (d *Deleter) Abort() []error {
	var ret []error
	for _, f := range append(d.files, d.dirs...) {
		if err := d.renameForRestore(f); err != nil {
			err = fmt.Errorf("restoring %q: %w", f, err)
			ret = append(ret, err)
		}
	}

	d.files = nil
	d.dirs = nil

	return ret
}

// Complete deletes all files marked for deletion and clears the marked list.
// It returns a slice of errors for each file that failed to be deleted. All
// files will be attempted, regardless of the errors encountered.
func (d *Deleter) Complete() []error {
	var ret []error
	for _, f := range d.files {
		if err := d.RenamerRemover.Remove(f + deleteFileSuffix); err != nil {
			err = fmt.Errorf("deleting file %q: %w", f+deleteFileSuffix, err)
			ret = append(ret, err)
		}
	}

	for _, f := range d.dirs {
		if err := d.RenamerRemover.RemoveAll(f + deleteFileSuffix); err != nil {
			err = fmt.Errorf("deleting directory %q: %w", f+deleteFileSuffix, err)
			ret = append(ret, err)
		}
	}

	d.files = nil
	d.dirs = nil

	return ret
}

func (d *Deleter) renameForDelete(path string) error {
	return d.RenamerRemover.Rename(path, path+deleteFileSuffix)
}

func (d *Deleter) renameForRestore(path string) error {
	return d.RenamerRemover.Rename(path+deleteFileSuffix, path)
}
