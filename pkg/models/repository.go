package models

import (
	"context"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/txn"
)

type TxnManager interface {
	txn.Manager
	txn.DatabaseProvider
	Reset() error
}

type Repository struct {
	TxnManager

	File        file.Store
	Folder      file.FolderStore
	Gallery     GalleryReaderWriter
	Image       ImageReaderWriter
	Movie       MovieReaderWriter
	Performer   PerformerReaderWriter
	Scene       SceneReaderWriter
	SceneMarker SceneMarkerReaderWriter
	ScrapedItem ScrapedItemReaderWriter
	Studio      StudioReaderWriter
	Tag         TagReaderWriter
	SavedFilter SavedFilterReaderWriter
}

func (r *Repository) WithTxn(ctx context.Context, fn txn.TxnFunc) error {
	return txn.WithTxn(ctx, r, fn)
}
