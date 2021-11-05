//go:generate go run -mod=vendor github.com/vektah/dataloaden SceneLoader int *github.com/stashapp/stash/pkg/models.Scene
package models

import (
	"context"
	"errors"
	"time"
)

var ErrTxn = errors.New("transaction failure")

func NewSceneLoaderConfig(ctx context.Context, mgr TransactionManager) SceneLoaderConfig {
	return SceneLoaderConfig{
		Wait:     2 * time.Millisecond,
		MaxBatch: 100,
		Fetch: func(keys []int) ([]*Scene, []error) {
			var scenes []*Scene
			var errors []error
			err := mgr.WithReadTxn(ctx, func(repo ReaderRepository) error {
				// The right implementation here would be a FindMany which return
				// all errors it finds, and optimizes the query to be a batch query:
				// either via sql IN, or via a temporary table join. At the time of
				// this writing, however, FindMany is just a wrapper around Find(),
				// so we mimic FindMany here to grab all errors and satisfy the dataloaden
				// fetch interface.
				r := repo.Scene()

				for _, k := range keys {
					s, e := r.Find(k)
					scenes = append(scenes, s)
					errors = append(errors, e)
				}

				return nil
			})

			if err != nil {
				// Txn failure, build error result set
				scenes = nil
				errors = nil
				for _ = range keys {
					scenes = append(scenes, nil)
					errors = append(errors, ErrTxn)
				}
			}

			return scenes, errors
		},
	}
}
