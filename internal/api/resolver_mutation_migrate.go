package api

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) MigrateSceneScreenshots(ctx context.Context, input MigrateSceneScreenshotsInput) (string, error) {
	db := manager.GetInstance().Database
	t := &migrateSceneJob{
		input:      input,
		sceneRepo:  db.Scene,
		txnManager: db,
	}
	jobID := manager.GetInstance().JobManager.Add(ctx, "Migrating scene screenshots to blobs...", t)

	return strconv.Itoa(jobID), nil
}

type migrateSceneJob struct {
	input      MigrateSceneScreenshotsInput
	sceneRepo  scene.HashFinderCoverUpdater
	txnManager txn.Manager
}

func (j *migrateSceneJob) Execute(ctx context.Context, progress *job.Progress) {
	paths := manager.GetInstance().Paths

	var err error
	progress.ExecuteTask("Counting files", func() {
		var count int
		count, err = j.countFiles(ctx, paths.Scene.Screenshots)
		progress.SetTotal(count)
	})

	if err != nil {
		logger.Errorf("Error counting files: %s", err.Error())
		return
	}

	progress.ExecuteTask("Migrating files", func() {
		err = j.migrateFiles(ctx, paths.Scene.Screenshots, progress)
	})

	if job.IsCancelled(ctx) {
		logger.Info("Cancelled migrating scene screenshots")
		return
	}

	if err != nil && !errors.Is(err, io.EOF) {
		logger.Errorf("Error reading screenshots directory: %v", err)
		return
	}

	logger.Infof("Finished migrating scene screenshots")
}

func (j *migrateSceneJob) countFiles(ctx context.Context, screenshotDir string) (int, error) {
	f, err := os.Open(screenshotDir)
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

func (j *migrateSceneJob) migrateFiles(ctx context.Context, screenshotDir string, progress *job.Progress) error {
	f, err := os.Open(screenshotDir)
	if err != nil {
		return err
	}
	defer f.Close()

	m := scene.ScreenshotMigrator{
		Options: scene.MigrateSceneScreenshotsInput{
			DeleteFiles:       utils.IsTrue(j.input.DeleteFiles),
			OverwriteExisting: utils.IsTrue(j.input.OverwriteExisting),
		},
		SceneUpdater: j.sceneRepo,
		TxnManager:   j.txnManager,
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

				path := filepath.Join(screenshotDir, f.Name())

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
