package mocks

import (
	context "context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

type TxnManager struct{}

func (*TxnManager) Begin(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (*TxnManager) WithDatabase(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (*TxnManager) Commit(ctx context.Context) error {
	return nil
}

func (*TxnManager) Rollback(ctx context.Context) error {
	return nil
}

func (*TxnManager) AddPostCommitHook(ctx context.Context, hook txn.TxnFunc) {
}

func (*TxnManager) AddPostRollbackHook(ctx context.Context, hook txn.TxnFunc) {
}

func (*TxnManager) Reset() error {
	return nil
}

func NewTxnRepository() models.Repository {
	return models.Repository{
		TxnManager:  &TxnManager{},
		Gallery:     &GalleryReaderWriter{},
		Image:       &ImageReaderWriter{},
		Movie:       &MovieReaderWriter{},
		Performer:   &PerformerReaderWriter{},
		Scene:       &SceneReaderWriter{},
		SceneMarker: &SceneMarkerReaderWriter{},
		ScrapedItem: &ScrapedItemReaderWriter{},
		Studio:      &StudioReaderWriter{},
		Tag:         &TagReaderWriter{},
		SavedFilter: &SavedFilterReaderWriter{},
	}
}
