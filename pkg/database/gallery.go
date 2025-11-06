package database

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type GalleryStore interface {
	AddFileID(ctx context.Context, id int, fileID models.FileID) error
	AddImages(ctx context.Context, galleryID int, imageIDs ...int) error
	All(ctx context.Context) ([]*models.Gallery, error)
	Count(ctx context.Context) (int, error)
	CountByFileID(ctx context.Context, fileID models.FileID) (int, error)
	CountByImageID(ctx context.Context, imageID int) (int, error)
	Create(ctx context.Context, newObject *models.Gallery, fileIDs []models.FileID) error
	Destroy(ctx context.Context, id int) error
	Find(ctx context.Context, id int) (*models.Gallery, error)
	FindByChecksum(ctx context.Context, checksum string) ([]*models.Gallery, error)
	FindByChecksums(ctx context.Context, checksums []string) ([]*models.Gallery, error)
	FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Gallery, error)
	FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Gallery, error)
	FindByFolderID(ctx context.Context, folderID models.FolderID) ([]*models.Gallery, error)
	FindByImageID(ctx context.Context, imageID int) ([]*models.Gallery, error)
	FindByPath(ctx context.Context, p string) ([]*models.Gallery, error)
	FindBySceneID(ctx context.Context, sceneID int) ([]*models.Gallery, error)
	FindMany(ctx context.Context, ids []int) ([]*models.Gallery, error)
	FindUserGalleryByTitle(ctx context.Context, title string) ([]*models.Gallery, error)
	GetFiles(ctx context.Context, id int) ([]models.File, error)
	GetImageIDs(ctx context.Context, galleryID int) ([]int, error)
	GetManyFileIDs(ctx context.Context, ids []int) ([][]models.FileID, error)
	GetPerformerIDs(ctx context.Context, id int) ([]int, error)
	GetSceneIDs(ctx context.Context, id int) ([]int, error)
	GetTagIDs(ctx context.Context, id int) ([]int, error)
	GetURLs(ctx context.Context, galleryID int) ([]string, error)
	Query(ctx context.Context, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) ([]*models.Gallery, int, error)
	QueryCount(ctx context.Context, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) (int, error)
	RemoveImages(ctx context.Context, galleryID int, imageIDs ...int) error
	ResetCover(ctx context.Context, galleryID int) error
	SetCover(ctx context.Context, galleryID int, coverImageID int) error
	Update(ctx context.Context, updatedObject *models.Gallery) error
	UpdateImages(ctx context.Context, galleryID int, imageIDs []int) error
	UpdatePartial(ctx context.Context, id int, partial models.GalleryPartial) (*models.Gallery, error)
}
