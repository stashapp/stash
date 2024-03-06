package task

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/txn"
)

type MigrateSceneScreenshotsJob struct {
	ScreenshotsPath string
	Input           scene.MigrateSceneScreenshotsInput
	SceneRepo       scene.HashFinderCoverUpdater
	TxnManager      txn.Manager
}

func (j *MigrateSceneScreenshotsJob) Execute(ctx context.Context, progress *job.Progress) error {
	var err error
	progress.ExecuteTask("Counting files", func() {
		var count int
		count, err = j.countFiles(ctx)
		progress.SetTotal(count)
	})

	if err != nil {
		return fmt.Errorf("error counting files: %w", err)
	}

	progress.ExecuteTask("Migrating files", func() {
		err = j.migrateFiles(ctx, progress)
	})

	if job.IsCancelled(ctx) {
		logger.Info("Cancelled migrating scene screenshots")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error migrating scene screenshots: %w", err)
	}

	logger.Infof("Finished migrating scene screenshots")
	return nil
}

func (j *MigrateSceneScreenshotsJob) countFiles(ctx context.Context) (int, error) {
	f, err := os.Open(j.ScreenshotsPath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	const batchSize = 1000
	ret := 0
	files, err := f.ReadDir(batchSize)
	for err == nil && ctx.Err() == nil {
		ret += len(files)

		files, err = f.ReadDir(batchSize)
	}

	if errors.Is(err, io.EOF) {
		// end of directory
		return ret, nil
	}

	return 0, err
}

func (j *MigrateSceneScreenshotsJob) migrateFiles(ctx context.Context, progress *job.Progress) error {
	f, err := os.Open(j.ScreenshotsPath)
	if err != nil {
		return err
	}
	defer f.Close()

	m := scene.ScreenshotMigrator{
		Options:      j.Input,
		SceneUpdater: j.SceneRepo,
		TxnManager:   j.TxnManager,
	}

	const batchSize = 1000
	files, err := f.ReadDir(batchSize)
	for err == nil && ctx.Err() == nil {
		for _, f := range files {
			if ctx.Err() != nil {
				return nil
			}

			progress.ExecuteTask("Migrating file "+f.Name(), func() {
				defer progress.Increment()

				path := filepath.Join(j.ScreenshotsPath, f.Name())

				// sanity check - only process files
				if f.IsDir() {
					logger.Warnf("Skipping directory %s", path)
					return
				}

				// ignore non-jpg files
				if !strings.HasSuffix(f.Name(), ".jpg") {
					return
				}

				// ignore .thumb files
				if strings.HasSuffix(f.Name(), ".thumb.jpg") {
					return
				}

				if err := m.MigrateScreenshots(ctx, path); err != nil {
					logger.Errorf("Error migrating screenshots for %s: %v", path, err)
				}
			})
		}

		files, err = f.ReadDir(batchSize)
	}

	if errors.Is(err, io.EOF) {
		// end of directory
		return nil
	}

	return err
}
