package database

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type SavedFilterStore interface {
	All(ctx context.Context) ([]*models.SavedFilter, error)
	Create(ctx context.Context, newObject *models.SavedFilter) error
	Destroy(ctx context.Context, id int) error
	Find(ctx context.Context, id int) (*models.SavedFilter, error)
	FindByMode(ctx context.Context, mode models.FilterMode) ([]*models.SavedFilter, error)
	FindMany(ctx context.Context, ids []int, ignoreNotFound bool) ([]*models.SavedFilter, error)
	Update(ctx context.Context, updatedObject *models.SavedFilter) error
}
