package manager

import (
	"context"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/txn"
)

type ImageReaderWriter interface {
	models.ImageReaderWriter
	image.FinderCreatorUpdater
}

type GalleryReaderWriter interface {
	models.GalleryReaderWriter
	gallery.FinderCreatorUpdater
}

type SceneReaderWriter interface {
	models.SceneReaderWriter
	scene.CreatorUpdater
}

type FileReaderWriter interface {
	file.Store
	file.Finder
	Query(ctx context.Context, options models.FileQueryOptions) (*models.FileQueryResult, error)
	GetCaptions(ctx context.Context, fileID file.ID) ([]*models.VideoCaption, error)
}

type FolderReaderWriter interface {
	file.FolderStore
	Find(ctx context.Context, id file.FolderID) (*file.Folder, error)
}

type Repository struct {
	models.TxnManager

	File        FileReaderWriter
	Folder      FolderReaderWriter
	Gallery     GalleryReaderWriter
	Image       ImageReaderWriter
	Movie       models.MovieReaderWriter
	Performer   models.PerformerReaderWriter
	Scene       SceneReaderWriter
	SceneMarker models.SceneMarkerReaderWriter
	ScrapedItem models.ScrapedItemReaderWriter
	Studio      models.StudioReaderWriter
	Tag         models.TagReaderWriter
	SavedFilter models.SavedFilterReaderWriter
}

func (r *Repository) WithTxn(ctx context.Context, fn txn.TxnFunc) error {
	return txn.WithTxn(ctx, r, fn)
}

func sqliteRepository(d *sqlite.Database) Repository {
	txnRepo := d.TxnRepository()

	return Repository{
		TxnManager:  txnRepo,
		File:        d.File,
		Folder:      d.Folder,
		Gallery:     d.Gallery,
		Image:       d.Image,
		Movie:       txnRepo.Movie,
		Performer:   txnRepo.Performer,
		Scene:       d.Scene,
		SceneMarker: txnRepo.SceneMarker,
		ScrapedItem: txnRepo.ScrapedItem,
		Studio:      txnRepo.Studio,
		Tag:         txnRepo.Tag,
		SavedFilter: txnRepo.SavedFilter,
	}
}

type SceneService interface {
	Destroy(ctx context.Context, scene *models.Scene, fileDeleter *scene.FileDeleter, deleteGenerated, deleteFile bool) error
}

type ImageService interface {
	Destroy(ctx context.Context, image *models.Image, fileDeleter *image.FileDeleter, deleteGenerated, deleteFile bool) error
}

type GalleryService interface {
	Destroy(ctx context.Context, i *models.Gallery, fileDeleter *image.FileDeleter, deleteGenerated, deleteFile bool) ([]*models.Image, error)
}
