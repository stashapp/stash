package database

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type StudioStore interface {
	All(ctx context.Context) ([]*models.Studio, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, newObject *models.Studio) error
	Destroy(ctx context.Context, id int) error
	Find(ctx context.Context, id int) (*models.Studio, error)
	FindByName(ctx context.Context, name string, nocase bool) (*models.Studio, error)
	FindBySceneID(ctx context.Context, sceneID int) (*models.Studio, error)
	FindByStashID(ctx context.Context, stashID models.StashID) ([]*models.Studio, error)
	FindByStashIDStatus(ctx context.Context, hasStashID bool, stashboxEndpoint string) ([]*models.Studio, error)
	FindChildren(ctx context.Context, id int) ([]*models.Studio, error)
	FindMany(ctx context.Context, ids []int) ([]*models.Studio, error)
	GetAliases(ctx context.Context, studioID int) ([]string, error)
	GetImage(ctx context.Context, studioID int) ([]byte, error)
	GetStashIDs(ctx context.Context, studioID int) ([]models.StashID, error)
	HasImage(ctx context.Context, studioID int) (bool, error)
	Query(ctx context.Context, studioFilter *models.StudioFilterType, findFilter *models.FindFilterType) ([]*models.Studio, int, error)
	QueryCount(ctx context.Context, studioFilter *models.StudioFilterType, findFilter *models.FindFilterType) (int, error)
	QueryForAutoTag(ctx context.Context, words []string) ([]*models.Studio, error)
	Update(ctx context.Context, updatedObject *models.Studio) error
	UpdateImage(ctx context.Context, studioID int, image []byte) error
	UpdatePartial(ctx context.Context, input models.StudioPartial) (*models.Studio, error)
	GetURLs(ctx context.Context, studioID int) ([]string, error)

	// blobJoinQueryBuilder
	tagRelationshipStore
}

type tagRelationshipStore interface {
	CountByTagID(ctx context.Context, tagID int) (int, error)
	GetTagIDs(ctx context.Context, id int) ([]int, error)
}
