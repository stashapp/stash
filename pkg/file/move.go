package file

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

type Renamer interface {
	Rename(oldpath, newpath string) error
}

type Statter interface {
	Stat(name string) (fs.FileInfo, error)
}

type DirMakerStatRenamer interface {
	Statter
	Renamer
	Mkdir(name string, perm os.FileMode) error
	Remove(name string) error
}

type folderCreatorStatRenamerImpl struct {
	renamerRemoverImpl
	mkDirFn func(name string, perm os.FileMode) error
}

func (r folderCreatorStatRenamerImpl) Mkdir(name string, perm os.FileMode) error {
	return r.mkDirFn(name, perm)
}

type Mover struct {
	Renamer DirMakerStatRenamer
	Files   models.FileFinderUpdater
	Folders models.FolderReaderWriter

	moved          map[string]string
	foldersCreated []string
}

func NewMover(fileStore models.FileFinderUpdater, folderStore models.FolderReaderWriter) *Mover {
	return &Mover{
		Files:   fileStore,
		Folders: folderStore,
		Renamer: &folderCreatorStatRenamerImpl{
			renamerRemoverImpl: newRenamerRemoverImpl(),
			mkDirFn:            os.Mkdir,
		},
	}
}

// Move moves the file to the given folder and basename. If basename is empty, then the existing basename is used.
// Assumes that the parent folder exists in the filesystem.
func (m *Mover) Move(ctx context.Context, f models.File, folder *models.Folder, basename string) error {
	fBase := f.Base()

	// don't allow moving files in zip files
	if fBase.ZipFileID != nil {
		return fmt.Errorf("cannot move file %s, is in a zip file", fBase.Path)
	}

	if basename == "" {
		basename = fBase.Basename
	}

	// modify the database first

	oldPath := fBase.Path

	if folder.ID == fBase.ParentFolderID && (basename == "" || basename == fBase.Basename) {
		// nothing to do
		return nil
	}

	// ensure that the new path doesn't already exist
	newPath := filepath.Join(folder.Path, basename)
	if _, err := m.Renamer.Stat(newPath); !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("file %s already exists", newPath)
	}

	if err := transferZipHierarchy(ctx, m.Folders, m.Files, fBase.ID, oldPath, newPath); err != nil {
		return fmt.Errorf("moving folder hierarchy for file %s: %w", fBase.Path, err)
	}

	fBase.ParentFolderID = folder.ID
	fBase.Basename = basename
	fBase.UpdatedAt = time.Now()
	// leave ModTime as is. It may or may not be changed by this operation

	if err := m.Files.Update(ctx, f); err != nil {
		return fmt.Errorf("updating file %s: %w", oldPath, err)
	}

	// then move the file
	return m.moveFile(oldPath, newPath)
}

func (m *Mover) CreateFolderHierarchy(path string) error {
	info, err := m.Renamer.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// create the parent folder
			parentPath := filepath.Dir(path)
			if err := m.CreateFolderHierarchy(parentPath); err != nil {
				return err
			}

			// create the folder
			if err := m.Renamer.Mkdir(path, 0755); err != nil {
				return fmt.Errorf("creating folder %s: %w", path, err)
			}

			m.foldersCreated = append(m.foldersCreated, path)
		} else {
			return fmt.Errorf("getting info for %s: %w", path, err)
		}
	} else {
		if !info.IsDir() {
			return fmt.Errorf("%s is not a directory", path)
		}
	}

	return nil
}

func (m *Mover) moveFile(oldPath, newPath string) error {
	if err := m.Renamer.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("renaming file %s to %s: %w", oldPath, newPath, err)
	}

	if m.moved == nil {
		m.moved = make(map[string]string)
	}

	m.moved[newPath] = oldPath

	return nil
}

func (m *Mover) RegisterHooks(ctx context.Context) {
	txn.AddPostCommitHook(ctx, func(ctx context.Context) {
		m.commit()
	})

	txn.AddPostRollbackHook(ctx, func(ctx context.Context) {
		m.rollback()
	})
}

func (m *Mover) commit() {
	m.moved = nil
	m.foldersCreated = nil
}

func (m *Mover) rollback() {
	// move files back to their original location
	for newPath, oldPath := range m.moved {
		if err := m.Renamer.Rename(newPath, oldPath); err != nil {
			logger.Errorf("error moving file %s back to %s: %s", newPath, oldPath, err.Error())
		}
	}

	// remove folders created in reverse order
	for i := len(m.foldersCreated) - 1; i >= 0; i-- {
		folder := m.foldersCreated[i]
		if err := m.Renamer.Remove(folder); err != nil {
			logger.Errorf("error removing folder %s: %s", folder, err.Error())
		}
	}
}

// correctSubFolderHierarchy sets the path of all contained folders to be relative to the given folder.
// It does not move the folder hierarchy in the filesystem.
func correctSubFolderHierarchy(ctx context.Context, rw models.FolderReaderWriter, folder *models.Folder) error {
	folders, err := rw.FindByParentFolderID(ctx, folder.ID)
	if err != nil {
		return fmt.Errorf("finding contained folders in folder %s: %w", folder.Path, err)
	}

	folderPath := folder.Path

	for _, f := range folders {
		oldPath := f.Path
		folderBasename := filepath.Base(f.Path)
		correctPath := filepath.Join(folderPath, folderBasename)

		logger.Debugf("updating folder %s to %s", oldPath, correctPath)

		f.Path = correctPath
		if err := rw.Update(ctx, f); err != nil {
			return fmt.Errorf("updating folder path %s -> %s: %w", oldPath, f.Path, err)
		}

		// recurse
		if err := correctSubFolderHierarchy(ctx, rw, f); err != nil {
			return err
		}
	}

	return nil
}
