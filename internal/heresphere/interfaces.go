package heresphere

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

// Repository provides access to storage methods for files and folders.
type Repository struct {
	TxnManager models.TxnManager
}

func (r *Repository) withTxn(ctx context.Context, fn txn.TxnFunc) error {
	return txn.WithTxn(ctx, r.TxnManager, fn)
}

func (r *Repository) withReadTxn(ctx context.Context, fn txn.TxnFunc) error {
	return txn.WithReadTxn(ctx, r.TxnManager, fn)
}

type sceneFinder interface {
	models.SceneGetter
	models.SceneReader
	models.SceneWriter
}

type sceneMarkerFinder interface {
	models.SceneMarkerFinder
	models.SceneMarkerCreator
	models.SceneMarkerReader
}

type tagFinder interface {
	models.TagFinder
	models.TagCreator
}

type fileFinder interface {
	models.FileFinder
	models.FileReader
	models.FileDestroyer
}

type savedfilterFinder interface {
	models.SavedFilterReader
}

type performerFinder interface {
	models.PerformerFinder
}

type galleryFinder interface {
	models.GalleryFinder
}

type movieFinder interface {
	models.MovieFinder
}

type studioFinder interface {
	models.StudioFinder
}
