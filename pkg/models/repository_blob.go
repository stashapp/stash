package models

import "context"

// BlobReader provides methods to get files by ID.
type BlobReader interface {
	EntryExists(ctx context.Context, checksum string) (bool, error)
}
