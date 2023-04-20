package models

import "context"

type PinnedFilterReader interface {
	All(ctx context.Context) ([]*PinnedFilter, error)
	Find(ctx context.Context, id int) (*PinnedFilter, error)
	FindMany(ctx context.Context, ids []int, ignoreNotFound bool) ([]*PinnedFilter, error)
	FindByMode(ctx context.Context, mode FilterMode) ([]*PinnedFilter, error)
}

type PinnedFilterWriter interface {
	Create(ctx context.Context, obj PinnedFilter) (*PinnedFilter, error)
	Destroy(ctx context.Context, id int) error
}

type PinnedFilterReaderWriter interface {
	PinnedFilterReader
	PinnedFilterWriter
}
