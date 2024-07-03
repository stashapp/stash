package mocks

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stretchr/testify/mock"
)

type Database struct {
	File           *FileReaderWriter
	Folder         *FolderReaderWriter
	Gallery        *GalleryReaderWriter
	GalleryChapter *GalleryChapterReaderWriter
	Image          *ImageReaderWriter
	Group          *GroupReaderWriter
	Performer      *PerformerReaderWriter
	Scene          *SceneReaderWriter
	SceneMarker    *SceneMarkerReaderWriter
	Studio         *StudioReaderWriter
	Tag            *TagReaderWriter
	SavedFilter    *SavedFilterReaderWriter
}

func (*Database) Begin(ctx context.Context, exclusive bool) (context.Context, error) {
	return ctx, nil
}

func (*Database) WithDatabase(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (*Database) Commit(ctx context.Context) error {
	return nil
}

func (*Database) Rollback(ctx context.Context) error {
	return nil
}

func (*Database) Complete(ctx context.Context) {
}

func (*Database) AddPostCommitHook(ctx context.Context, hook txn.TxnFunc) {
}

func (*Database) AddPostRollbackHook(ctx context.Context, hook txn.TxnFunc) {
}

func (*Database) IsLocked(err error) bool {
	return false
}

func (*Database) Reset() error {
	return nil
}

func NewDatabase() *Database {
	return &Database{
		File:           &FileReaderWriter{},
		Folder:         &FolderReaderWriter{},
		Gallery:        &GalleryReaderWriter{},
		GalleryChapter: &GalleryChapterReaderWriter{},
		Image:          &ImageReaderWriter{},
		Group:          &GroupReaderWriter{},
		Performer:      &PerformerReaderWriter{},
		Scene:          &SceneReaderWriter{},
		SceneMarker:    &SceneMarkerReaderWriter{},
		Studio:         &StudioReaderWriter{},
		Tag:            &TagReaderWriter{},
		SavedFilter:    &SavedFilterReaderWriter{},
	}
}

func (db *Database) AssertExpectations(t mock.TestingT) {
	db.File.AssertExpectations(t)
	db.Folder.AssertExpectations(t)
	db.Gallery.AssertExpectations(t)
	db.GalleryChapter.AssertExpectations(t)
	db.Image.AssertExpectations(t)
	db.Group.AssertExpectations(t)
	db.Performer.AssertExpectations(t)
	db.Scene.AssertExpectations(t)
	db.SceneMarker.AssertExpectations(t)
	db.Studio.AssertExpectations(t)
	db.Tag.AssertExpectations(t)
	db.SavedFilter.AssertExpectations(t)
}

func (db *Database) Repository() models.Repository {
	return models.Repository{
		TxnManager:     db,
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
