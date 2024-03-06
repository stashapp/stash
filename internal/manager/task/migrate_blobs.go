package task

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/txn"
)

type BlobStoreMigrator interface {
	Count(ctx context.Context) (int, error)
	FindBlobs(ctx context.Context, n uint, lastChecksum string) ([]string, error)
	MigrateBlob(ctx context.Context, checksum string, deleteOld bool) error
}

type Vacuumer interface {
	Vacuum(ctx context.Context) error
}

type MigrateBlobsJob struct {
	TxnManager txn.Manager
	BlobStore  BlobStoreMigrator
	Vacuumer   Vacuumer
	DeleteOld  bool
}

func (j *MigrateBlobsJob) Execute(ctx context.Context, progress *job.Progress) error {
	var (
		count int
		err   error
	)
	progress.ExecuteTask("Counting blobs", func() {
		count, err = j.countBlobs(ctx)
		progress.SetTotal(count)
	})

	if err != nil {
		return fmt.Errorf("error counting blobs: %w", err)
	}

	if count == 0 {
		logger.Infof("No blobs to migrate")
		return nil
	}

	logger.Infof("Migrating %d blobs", count)

	progress.ExecuteTask(fmt.Sprintf("Migrating %d blobs", count), func() {
		err = j.migrateBlobs(ctx, progress)
	})

	if job.IsCancelled(ctx) {
		logger.Info("Cancelled migrating blobs")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error migrating blobs: %w", err)
	}

	// run a vacuum to reclaim space
	progress.ExecuteTask("Vacuuming database", func() {
		err = j.Vacuumer.Vacuum(ctx)
		if err != nil {
			logger.Errorf("Error vacuuming database: %v", err)
		}
	})

	logger.Infof("Finished migrating blobs")
	return nil
}

func (j *MigrateBlobsJob) countBlobs(ctx context.Context) (int, error) {
	var count int
	if err := txn.WithReadTxn(ctx, j.TxnManager, func(ctx context.Context) error {
		var err error
		count, err = j.BlobStore.Count(ctx)
		return err
	}); err != nil {
		return 0, err
	}

	return count, nil
}

func (j *MigrateBlobsJob) migrateBlobs(ctx context.Context, progress *job.Progress) error {
	lastChecksum := ""
	batch, err := j.getBatch(ctx, lastChecksum)

	for len(batch) > 0 && err == nil && ctx.Err() == nil {
		for _, checksum := range batch {
			if ctx.Err() != nil {
				return nil
			}

			lastChecksum = checksum

			progress.ExecuteTask("Migrating blob "+checksum, func() {
				defer progress.Increment()

				if err := txn.WithTxn(ctx, j.TxnManager, func(ctx context.Context) error {
					return j.BlobStore.MigrateBlob(ctx, checksum, j.DeleteOld)
				}); err != nil {
					logger.Errorf("Error migrating blob %s: %v", checksum, err)
				}
			})
		}

		batch, err = j.getBatch(ctx, lastChecksum)
	}

	return err
}

func (j *MigrateBlobsJob) getBatch(ctx context.Context, lastChecksum string) ([]string, error) {
	const batchSize = 1000

	var batch []string
	err := txn.WithReadTxn(ctx, j.TxnManager, func(ctx context.Context) error {
		var err error
		batch, err = j.BlobStore.FindBlobs(ctx, batchSize, lastChecksum)
		return err
	})

	return batch, err
}
