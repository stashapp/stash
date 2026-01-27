package sqlite

import (
	"github.com/stashapp/stash/pkg/database"
)

const defaultBatchSize = database.DefaultBatchSize

func batchExec[T any](ids []T, batchSize int, fn func(batch []T) error) error {
	return database.BatchExec(ids, batchSize, fn)
}
