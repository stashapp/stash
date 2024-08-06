package group

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type CreatorUpdater interface {
	models.GroupGetter
	models.GroupCreator
	models.GroupUpdater

	models.ContainingGroupLoader
	models.SubGroupLoader

	AnscestorFinder
	SubGroupReorderer
}

type AnscestorFinder interface {
	FindInAncestors(ctx context.Context, ascestorIDs []int, ids []int) ([]int, error)
}

type SubGroupReorderer interface {
	ReorderSubGroups(ctx context.Context, groupID int, subGroupIDs []int, insertID int, insertAfter bool) error
}

type Service struct {
	Repository CreatorUpdater
}
