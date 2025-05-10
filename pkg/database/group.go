package database

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type GroupStore interface {
	All(ctx context.Context) ([]*models.Group, error)
	Count(ctx context.Context) (int, error)
	CountByPerformerID(ctx context.Context, performerID int) (int, error)
	CountByStudioID(ctx context.Context, studioID int) (int, error)
	Create(ctx context.Context, newObject *models.Group) error
	Destroy(ctx context.Context, id int) error
	Find(ctx context.Context, id int) (*models.Group, error)
	FindByName(ctx context.Context, name string, nocase bool) (*models.Group, error)
	FindByNames(ctx context.Context, names []string, nocase bool) ([]*models.Group, error)
	FindByPerformerID(ctx context.Context, performerID int) ([]*models.Group, error)
	FindByStudioID(ctx context.Context, studioID int) ([]*models.Group, error)
	FindInAncestors(ctx context.Context, ascestorIDs []int, ids []int) ([]int, error)
	FindMany(ctx context.Context, ids []int) ([]*models.Group, error)
	FindSubGroupIDs(ctx context.Context, containingID int, ids []int) ([]int, error)
	GetBackImage(ctx context.Context, groupID int) ([]byte, error)
	GetFrontImage(ctx context.Context, groupID int) ([]byte, error)
	GetURLs(ctx context.Context, groupID int) ([]string, error)
	HasBackImage(ctx context.Context, groupID int) (bool, error)
	HasFrontImage(ctx context.Context, groupID int) (bool, error)
	Query(ctx context.Context, groupFilter *models.GroupFilterType, findFilter *models.FindFilterType) ([]*models.Group, int, error)
	QueryCount(ctx context.Context, groupFilter *models.GroupFilterType, findFilter *models.FindFilterType) (int, error)
	Update(ctx context.Context, updatedObject *models.Group) error
	UpdateBackImage(ctx context.Context, groupID int, backImage []byte) error
	UpdateFrontImage(ctx context.Context, groupID int, frontImage []byte) error
	UpdatePartial(ctx context.Context, id int, partial models.GroupPartial) (*models.Group, error)

	blobJoinQueryBuilder
	tagRelationshipStore
	groupRelationshipStore
}

type groupRelationshipStore interface {
	AddSubGroups(ctx context.Context, groupID int, subGroups []models.GroupIDDescription, insertIndex *int) error
	GetContainingGroupDescriptions(ctx context.Context, id int) ([]models.GroupIDDescription, error)
	GetSubGroupDescriptions(ctx context.Context, id int) ([]models.GroupIDDescription, error)
	RemoveSubGroups(ctx context.Context, groupID int, subGroupIDs []int) error
	ReorderSubGroups(ctx context.Context, groupID int, subGroupIDs []int, insertPointID int, insertAfter bool) error
}
