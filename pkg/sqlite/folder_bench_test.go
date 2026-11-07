//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

// BenchmarkFolderFindByPath_Empty measures the cost of a FindByPath miss on an
// empty folders table — the dominant per-directory cost on a cold scan.
// Wrap in WithReadTxn to match CheckFolder's exact call path.
func BenchmarkFolderFindByPath_Empty(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = txn.WithReadTxn(ctx, db, func(ctx context.Context) error {
			_, err := db.Folder.FindByPath(ctx, "/nonexistent/benchmark/path", true)
			return err
		})
	}
}

// BenchmarkFolderFindByPath_Populated measures the cost of a FindByPath hit on
// a table with ~1000 folders — the warm-scan case.
func BenchmarkFolderFindByPath_Populated(b *testing.B) {
	ctx := context.Background()

	// Pre-populate 1000 folders in a write transaction, then roll back after
	// the benchmark so we don't pollute other tests.
	const folderCount = 1000
	paths := make([]string, folderCount)
	for i := range paths {
		paths[i] = fmt.Sprintf("/bench/library/dir%04d", i)
	}

	// Insert outside the timed region.
	if err := withTxn(func(ctx context.Context) error {
		for _, p := range paths {
			f := models.Folder{
				Path:     p,
				DirEntry: models.DirEntry{ModTime: time.Now()},
			}
			if err := db.Folder.Create(ctx, &f); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		b.Fatalf("setup: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = txn.WithReadTxn(ctx, db, func(ctx context.Context) error {
			_, err := db.Folder.FindByPath(ctx, paths[i%folderCount], true)
			return err
		})
	}
	b.StopTimer()

	// Clean up inserted rows.
	_ = withTxn(func(ctx context.Context) error {
		for _, p := range paths {
			f, err := db.Folder.FindByPath(ctx, p, true)
			if err != nil || f == nil {
				continue
			}
			if err := db.Folder.Destroy(ctx, f.ID); err != nil {
				return err
			}
		}
		return nil
	})
}
