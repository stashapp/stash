package database

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type PerformerStore interface {
	All(ctx context.Context) ([]*models.Performer, error)
	Count(ctx context.Context) (int, error)
	CountByTagID(ctx context.Context, tagID int) (int, error)
	Create(ctx context.Context, newObject *models.CreatePerformerInput) error
	Destroy(ctx context.Context, id int) error
	DestroyImage(ctx context.Context, id int, blobCol string) error
	Find(ctx context.Context, id int) (*models.Performer, error)
	FindByGalleryID(ctx context.Context, galleryID int) ([]*models.Performer, error)
	FindByImageID(ctx context.Context, imageID int) ([]*models.Performer, error)
	FindByNames(ctx context.Context, names []string, nocase bool) ([]*models.Performer, error)
	FindBySceneID(ctx context.Context, sceneID int) ([]*models.Performer, error)
	FindByStashID(ctx context.Context, stashID models.StashID) ([]*models.Performer, error)
	FindByStashIDStatus(ctx context.Context, hasStashID bool, stashboxEndpoint string) ([]*models.Performer, error)
	FindMany(ctx context.Context, ids []int) ([]*models.Performer, error)
	GetAliases(ctx context.Context, performerID int) ([]string, error)
	GetImage(ctx context.Context, performerID int) ([]byte, error)
	GetStashIDs(ctx context.Context, performerID int) ([]models.StashID, error)
	GetTagIDs(ctx context.Context, id int) ([]int, error)
	GetURLs(ctx context.Context, performerID int) ([]string, error)
	HasImage(ctx context.Context, performerID int) (bool, error)
	Query(ctx context.Context, performerFilter *models.PerformerFilterType, findFilter *models.FindFilterType) ([]*models.Performer, int, error)
	QueryCount(ctx context.Context, performerFilter *models.PerformerFilterType, findFilter *models.FindFilterType) (int, error)
	QueryForAutoTag(ctx context.Context, words []string) ([]*models.Performer, error)
	Update(ctx context.Context, updatedObject *models.UpdatePerformerInput) error
	UpdateImage(ctx context.Context, performerID int, image []byte) error
	UpdatePartial(ctx context.Context, id int, partial models.PerformerPartial) (*models.Performer, error)
	// blobJoinQueryBuilder
	customFieldsStore
}
