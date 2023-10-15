package manager

import (
	"context"

	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/txn"
)

type Repository struct {
	models.TxnManager

	File           models.FileReaderWriter
	Folder         models.FolderReaderWriter
	Gallery        models.GalleryReaderWriter
	GalleryChapter models.GalleryChapterReaderWriter
	Image          models.ImageReaderWriter
	Movie          models.MovieReaderWriter
	Performer      models.PerformerReaderWriter
	Scene          models.SceneReaderWriter
	SceneMarker    models.SceneMarkerReaderWriter
	Studio         models.StudioReaderWriter
	Tag            models.TagReaderWriter
	SavedFilter    models.SavedFilterReaderWriter
}

func (r *Repository) WithTxn(ctx context.Context, fn txn.TxnFunc) error {
	return txn.WithTxn(ctx, r, fn)
}

func (r *Repository) WithReadTxn(ctx context.Context, fn txn.TxnFunc) error {
	return txn.WithReadTxn(ctx, r, fn)
}

func (r *Repository) WithDB(ctx context.Context, fn txn.TxnFunc) error {
	return txn.WithDatabase(ctx, r, fn)
}

func sqliteRepository(d *sqlite.Database) Repository {
	txnRepo := d.TxnRepository()

	return Repository{
		TxnManager:     txnRepo,
		File:           d.File,
		Folder:         d.Folder,
		Gallery:        d.Gallery,
		GalleryChapter: txnRepo.GalleryChapter,
		Image:          d.Image,
		Movie:          txnRepo.Movie,
		Performer:      txnRepo.Performer,
		Scene:          d.Scene,
		SceneMarker:    txnRepo.SceneMarker,
		Studio:         txnRepo.Studio,
		Tag:            txnRepo.Tag,
		SavedFilter:    txnRepo.SavedFilter,
	}
}

type SceneService interface {
	Create(ctx context.Context, input *models.Scene, fileIDs []models.FileID, coverImage []byte) (*models.Scene, error)
	AssignFile(ctx context.Context, sceneID int, fileID models.FileID) error
	Merge(ctx context.Context, sourceIDs []int, destinationID int, values models.ScenePartial) error
	Destroy(ctx context.Context, scene *models.Scene, fileDeleter *scene.FileDeleter, deleteGenerated, deleteFile bool) error
}

type ImageService interface {
	Destroy(ctx context.Context, image *models.Image, fileDeleter *image.FileDeleter, deleteGenerated, deleteFile bool) error
	DestroyZipImages(ctx context.Context, zipFile models.File, fileDeleter *image.FileDeleter, deleteGenerated bool) ([]*models.Image, error)
}

type GalleryService interface {
	AddImages(ctx context.Context, g *models.Gallery, toAdd ...int) error
	RemoveImages(ctx context.Context, g *models.Gallery, toRemove ...int) error

	Destroy(ctx context.Context, i *models.Gallery, fileDeleter *image.FileDeleter, deleteGenerated, deleteFile bool) ([]*models.Image, error)

	ValidateImageGalleryChange(ctx context.Context, i *models.Image, updateIDs models.UpdateIDs) error

	Updated(ctx context.Context, galleryID int) error
}
