package database

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type ImageStore interface {
	AddFileID(ctx context.Context, id int, fileID models.FileID) error
	All(ctx context.Context) ([]*models.Image, error)
	Count(ctx context.Context) (int, error)
	CountByFileID(ctx context.Context, fileID models.FileID) (int, error)
	CountByGalleryID(ctx context.Context, galleryID int) (int, error)
	CoverByGalleryID(ctx context.Context, galleryID int) (*models.Image, error)
	Create(ctx context.Context, newObject *models.Image, fileIDs []models.FileID) error
	Destroy(ctx context.Context, id int) error
	Find(ctx context.Context, id int) (*models.Image, error)
	FindByChecksum(ctx context.Context, checksum string) ([]*models.Image, error)
	FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Image, error)
	FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Image, error)
	FindByFolderID(ctx context.Context, folderID models.FolderID) ([]*models.Image, error)
	FindByGalleryID(ctx context.Context, galleryID int) ([]*models.Image, error)
	FindByGalleryIDIndex(ctx context.Context, galleryID int, index uint) (*models.Image, error)
	FindByZipFileID(ctx context.Context, zipFileID models.FileID) ([]*models.Image, error)
	FindMany(ctx context.Context, ids []int) ([]*models.Image, error)
	GetFiles(ctx context.Context, id int) ([]models.File, error)
	GetGalleryIDs(ctx context.Context, imageID int) ([]int, error)
	GetManyFileIDs(ctx context.Context, ids []int) ([][]models.FileID, error)
	GetPerformerIDs(ctx context.Context, imageID int) ([]int, error)
	GetTagIDs(ctx context.Context, imageID int) ([]int, error)
	GetURLs(ctx context.Context, imageID int) ([]string, error)
	OCount(ctx context.Context) (int, error)
	OCountByPerformerID(ctx context.Context, performerID int) (int, error)
	Query(ctx context.Context, options models.ImageQueryOptions) (*models.ImageQueryResult, error)
	QueryCount(ctx context.Context, imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) (int, error)
	RemoveFileID(ctx context.Context, id int, fileID models.FileID) error
	Size(ctx context.Context) (float64, error)
	Update(ctx context.Context, updatedObject *models.Image) error
	UpdatePartial(ctx context.Context, id int, partial models.ImagePartial) (*models.Image, error)
	UpdatePerformers(ctx context.Context, imageID int, performerIDs []int) error
	UpdateTags(ctx context.Context, imageID int, tagIDs []int) error
	OCountByStudioID(ctx context.Context, studioID int) (int, error)
	OCountStore
}

type OCountStore interface {
	DecrementOCounter(ctx context.Context, id int) (int, error)
	IncrementOCounter(ctx context.Context, id int) (int, error)
	ResetOCounter(ctx context.Context, id int) (int, error)
}
