package file

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/txn"
)

type Renamer interface {
	Rename(oldpath, newpath string) error
}

type Mover struct {
	Renamer Renamer
	Updater Updater

	moved map[string]string
}

func (m *Mover) Move(ctx context.Context, f File, folder Folder, basename string) error {
	// don't allow moving files in zip files
	if f.Base().ZipFileID != nil {
		return fmt.Errorf("cannot move file %s in zip file", f.Base().Path)
	}

	// modify the database first
	fBase := f.Base()
	oldPath := fBase.Path
	fBase.ParentFolderID = folder.ID
	fBase.Basename = basename

	if err := m.Updater.Update(ctx, f); err != nil {
		return fmt.Errorf("updating file %s: %w", oldPath, err)
	}

	// then move the file
	newPath := filepath.Join(folder.Path, basename)
	return m.moveFile(oldPath, newPath)
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

func (m *Mover) RegisterHooks(ctx context.Context, mgr txn.Manager) {
	txn.AddPostCommitHook(ctx, func(ctx context.Context) {
		m.commit()
	})

	txn.AddPostRollbackHook(ctx, func(ctx context.Context) {
		m.rollback()
	})
}

func (m *Mover) commit() {
	m.moved = nil
}

func (m *Mover) rollback() {
	// move files back to their original location
	for newPath, oldPath := range m.moved {
		if err := m.Renamer.Rename(newPath, oldPath); err != nil {
			logger.Errorf("error moving file %s back to %s: %s", newPath, oldPath, err.Error())
		}
	}
}
