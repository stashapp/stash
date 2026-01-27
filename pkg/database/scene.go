package database

import (
	"context"
	"time"

	"github.com/stashapp/stash/pkg/models"
)

type SceneStore interface {
	AddFileID(ctx context.Context, id int, fileID models.FileID) error
	AddGalleryIDs(ctx context.Context, sceneID int, galleryIDs []int) error
	All(ctx context.Context) ([]*models.Scene, error)
	AssignFiles(ctx context.Context, sceneID int, fileIDs []models.FileID) error
	Count(ctx context.Context) (int, error)
	CountByFileID(ctx context.Context, fileID models.FileID) (int, error)
	CountByPerformerID(ctx context.Context, performerID int) (int, error)
	CountByStudioID(ctx context.Context, studioID int) (int, error)
	CountMissingChecksum(ctx context.Context) (int, error)
	CountMissingOSHash(ctx context.Context) (int, error)
	Create(ctx context.Context, newObject *models.Scene, fileIDs []models.FileID) error
	Destroy(ctx context.Context, id int) error
	Duration(ctx context.Context) (float64, error)
	Find(ctx context.Context, id int) (*models.Scene, error)
	FindByChecksum(ctx context.Context, checksum string) ([]*models.Scene, error)
	FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Scene, error)
	FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Scene, error)
	FindByGalleryID(ctx context.Context, galleryID int) ([]*models.Scene, error)
	FindByGroupID(ctx context.Context, groupID int) ([]*models.Scene, error)
	FindByIDs(ctx context.Context, ids []int) ([]*models.Scene, error)
	FindByOSHash(ctx context.Context, oshash string) ([]*models.Scene, error)
	FindByPath(ctx context.Context, p string) ([]*models.Scene, error)
	FindByPerformerID(ctx context.Context, performerID int) ([]*models.Scene, error)
	FindByPrimaryFileID(ctx context.Context, fileID models.FileID) ([]*models.Scene, error)
	FindDuplicates(ctx context.Context, distance int, durationDiff float64) ([][]*models.Scene, error)
	FindMany(ctx context.Context, ids []int) ([]*models.Scene, error)
	GetCover(ctx context.Context, sceneID int) ([]byte, error)
	GetFiles(ctx context.Context, id int) ([]*models.VideoFile, error)
	GetGalleryIDs(ctx context.Context, id int) ([]int, error)
	GetGroups(ctx context.Context, id int) (ret []models.GroupsScenes, err error)
	GetManyFileIDs(ctx context.Context, ids []int) ([][]models.FileID, error)
	GetPerformerIDs(ctx context.Context, id int) ([]int, error)
	GetStashIDs(ctx context.Context, sceneID int) ([]models.StashID, error)
	GetTagIDs(ctx context.Context, id int) ([]int, error)
	GetURLs(ctx context.Context, sceneID int) ([]string, error)
	HasCover(ctx context.Context, sceneID int) (bool, error)
	OCountByPerformerID(ctx context.Context, performerID int) (int, error)
	PlayDuration(ctx context.Context) (float64, error)
	Query(ctx context.Context, options models.SceneQueryOptions) (*models.SceneQueryResult, error)
	QueryCount(ctx context.Context, sceneFilter *models.SceneFilterType, findFilter *models.FindFilterType) (int, error)
	ResetActivity(ctx context.Context, id int, resetResume bool, resetDuration bool) (bool, error)
	SaveActivity(ctx context.Context, id int, resumeTime *float64, playDuration *float64) (bool, error)
	Size(ctx context.Context) (float64, error)
	Update(ctx context.Context, updatedObject *models.Scene) error
	UpdateCover(ctx context.Context, sceneID int, image []byte) error
	UpdatePartial(ctx context.Context, id int, partial models.ScenePartial) (*models.Scene, error)
	Wall(ctx context.Context, q *string) ([]*models.Scene, error)
	OCountByGroupID(ctx context.Context, groupID int) (int, error)
	OCountByStudioID(ctx context.Context, studioID int) (int, error)
	blobJoinQueryBuilder
	oDateManager
	viewDateManager
}

type blobJoinQueryBuilder interface {
	DestroyImage(ctx context.Context, id int, blobCol string) error
	GetImage(ctx context.Context, id int, blobCol string) ([]byte, error)
	HasImage(ctx context.Context, id int, blobCol string) (bool, error)
	UpdateImage(ctx context.Context, id int, blobCol string, image []byte) error
}

type oDateManager interface {
	AddO(ctx context.Context, id int, dates []time.Time) ([]time.Time, error)
	DeleteO(ctx context.Context, id int, dates []time.Time) ([]time.Time, error)
	GetAllOCount(ctx context.Context) (int, error)
	GetManyOCount(ctx context.Context, ids []int) ([]int, error)
	GetManyODates(ctx context.Context, ids []int) ([][]time.Time, error)
	GetOCount(ctx context.Context, id int) (int, error)
	GetODates(ctx context.Context, id int) ([]time.Time, error)
	GetUniqueOCount(ctx context.Context) (int, error)
	ResetO(ctx context.Context, id int) (int, error)
}

type viewDateManager interface {
	AddViews(ctx context.Context, id int, dates []time.Time) ([]time.Time, error)
	CountAllViews(ctx context.Context) (int, error)
	CountUniqueViews(ctx context.Context) (int, error)
	CountViews(ctx context.Context, id int) (int, error)
	DeleteAllViews(ctx context.Context, id int) (int, error)
	DeleteViews(ctx context.Context, id int, dates []time.Time) ([]time.Time, error)
	GetManyLastViewed(ctx context.Context, ids []int) ([]*time.Time, error)
	GetManyViewCount(ctx context.Context, ids []int) ([]int, error)
	GetManyViewDates(ctx context.Context, ids []int) ([][]time.Time, error)
	GetViewDates(ctx context.Context, id int) ([]time.Time, error)
	LastView(ctx context.Context, id int) (*time.Time, error)
}
