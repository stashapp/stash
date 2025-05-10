package database

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type TagStore interface {
	All(ctx context.Context) ([]*models.Tag, error)
	Count(ctx context.Context) (int, error)
	CountByChildTagID(ctx context.Context, childID int) (int, error)
	CountByParentTagID(ctx context.Context, parentID int) (int, error)
	Create(ctx context.Context, newObject *models.Tag) error
	Destroy(ctx context.Context, id int) error
	Find(ctx context.Context, id int) (*models.Tag, error)
	FindAllAncestors(ctx context.Context, tagID int, excludeIDs []int) ([]*models.TagPath, error)
	FindAllDescendants(ctx context.Context, tagID int, excludeIDs []int) ([]*models.TagPath, error)
	FindByChildTagID(ctx context.Context, parentID int) ([]*models.Tag, error)
	FindByGalleryID(ctx context.Context, galleryID int) ([]*models.Tag, error)
	FindByGroupID(ctx context.Context, groupID int) ([]*models.Tag, error)
	FindByImageID(ctx context.Context, imageID int) ([]*models.Tag, error)
	FindByName(ctx context.Context, name string, nocase bool) (*models.Tag, error)
	FindByNames(ctx context.Context, names []string, nocase bool) ([]*models.Tag, error)
	FindByParentTagID(ctx context.Context, parentID int) ([]*models.Tag, error)
	FindByPerformerID(ctx context.Context, performerID int) ([]*models.Tag, error)
	FindBySceneID(ctx context.Context, sceneID int) ([]*models.Tag, error)
	FindBySceneMarkerID(ctx context.Context, sceneMarkerID int) ([]*models.Tag, error)
	FindByStudioID(ctx context.Context, studioID int) ([]*models.Tag, error)
	FindByStashID(ctx context.Context, stashID models.StashID) ([]*models.Tag, error)
	FindMany(ctx context.Context, ids []int) ([]*models.Tag, error)
	GetAliases(ctx context.Context, tagID int) ([]string, error)
	GetChildIDs(ctx context.Context, relatedID int) ([]int, error)
	GetImage(ctx context.Context, tagID int) ([]byte, error)
	GetParentIDs(ctx context.Context, relatedID int) ([]int, error)
	HasImage(ctx context.Context, tagID int) (bool, error)
	Merge(ctx context.Context, source []int, destination int) error
	Query(ctx context.Context, tagFilter *models.TagFilterType, findFilter *models.FindFilterType) ([]*models.Tag, int, error)
	QueryForAutoTag(ctx context.Context, words []string) ([]*models.Tag, error)
	Update(ctx context.Context, updatedObject *models.Tag) error
	UpdateAliases(ctx context.Context, tagID int, aliases []string) error
	UpdateChildTags(ctx context.Context, tagID int, childIDs []int) error
	UpdateImage(ctx context.Context, tagID int, image []byte) error
	UpdateParentTags(ctx context.Context, tagID int, parentIDs []int) error
	UpdatePartial(ctx context.Context, id int, partial models.TagPartial) (*models.Tag, error)

	GetStashIDs(ctx context.Context, tagID int) ([]models.StashID, error)
	UpdateStashIDs(ctx context.Context, tagID int, stashIDs []models.StashID) error

	// blobJoinQueryBuilder
}
