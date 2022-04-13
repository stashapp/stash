package models

import "context"

type SavedFilterReader interface {
	Find(ctx context.Context, id int) (*SavedFilter, error)
	FindByMode(ctx context.Context, mode FilterMode) ([]*SavedFilter, error)
	FindDefault(ctx context.Context, mode FilterMode) (*SavedFilter, error)
}

type SavedFilterWriter interface {
	Create(ctx context.Context, obj SavedFilter) (*SavedFilter, error)
	Update(ctx context.Context, obj SavedFilter) (*SavedFilter, error)
	SetDefault(ctx context.Context, obj SavedFilter) (*SavedFilter, error)
	Destroy(ctx context.Context, id int) error
}

type SavedFilterReaderWriter interface {
	SavedFilterReader
	SavedFilterWriter
}
