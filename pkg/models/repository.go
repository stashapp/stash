package models

import (
	"context"

	"github.com/stashapp/stash/pkg/txn"
)

type TxnManager struct {
	txn.Manager
	txn.DatabaseProvider
}

type Repository struct {
	TxnManager

	Blob           BlobReader
	File           FileReaderWriter
	Folder         FolderReaderWriter
	Gallery        GalleryReaderWriter
	GalleryChapter GalleryChapterReaderWriter
	Image          ImageReaderWriter
	Group          GroupReaderWriter
	Performer      PerformerReaderWriter
	Scene          SceneReaderWriter
	SceneMarker    SceneMarkerReaderWriter
	Studio         StudioReaderWriter
	Tag            TagReaderWriter
	SavedFilter    SavedFilterReaderWriter
}

func (r *TxnManager) WithTxn(ctx context.Context, fn txn.TxnFunc) error {
	return txn.WithTxn(ctx, r, fn)
}

func (r *TxnManager) WithReadTxn(ctx context.Context, fn txn.TxnFunc) error {
	return txn.WithReadTxn(ctx, r, fn)
}

func (r *TxnManager) WithDB(ctx context.Context, fn txn.TxnFunc) error {
	return txn.WithDatabase(ctx, r, fn)
}
