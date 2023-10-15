package models

import (
	"github.com/stashapp/stash/pkg/txn"
)

type TxnManager interface {
	txn.Manager
	txn.DatabaseProvider
	Reset() error
}

type Repository struct {
	TxnManager

	File           FileReaderWriter
	Folder         FolderReaderWriter
	Gallery        GalleryReaderWriter
	GalleryChapter GalleryChapterReaderWriter
	Image          ImageReaderWriter
	Movie          MovieReaderWriter
	Performer      PerformerReaderWriter
	Scene          SceneReaderWriter
	SceneMarker    SceneMarkerReaderWriter
	Studio         StudioReaderWriter
	Tag            TagReaderWriter
	SavedFilter    SavedFilterReaderWriter
}
