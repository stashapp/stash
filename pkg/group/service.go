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
}

type AnscestorFinder interface {
	FindInAncestors(ctx context.Context, ascestorIDs []int, ids []int) ([]int, error)
}

type Service struct {
	Repository CreatorUpdater
}
