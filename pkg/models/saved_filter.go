package models

import "context"

type SavedFilterReader interface {
	All(ctx context.Context) ([]*SavedFilter, error)
	Find(ctx context.Context, id int) (*SavedFilter, error)
	FindMany(ctx context.Context, ids []int, ignoreNotFound bool) ([]*SavedFilter, error)
	FindByMode(ctx context.Context, mode FilterMode) ([]*SavedFilter, error)
}

type SavedFilterWriter interface {
	Create(ctx context.Context, obj *SavedFilter) error
	Update(ctx context.Context, obj *SavedFilter) error
	Destroy(ctx context.Context, id int) error
}

type SavedFilterReaderWriter interface {
	SavedFilterReader
	SavedFilterWriter
}
