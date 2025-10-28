package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/task"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) MigrateSceneScreenshots(ctx context.Context, input MigrateSceneScreenshotsInput) (string, error) {
	mgr := manager.GetInstance()
	t := &task.MigrateSceneScreenshotsJob{
		ScreenshotsPath: manager.GetInstance().Paths.Generated.Screenshots,
		Input: scene.MigrateSceneScreenshotsInput{
			DeleteFiles:       utils.IsTrue(input.DeleteFiles),
			OverwriteExisting: utils.IsTrue(input.OverwriteExisting),
		},
		SceneRepo:  mgr.Repository.Scene,
		TxnManager: mgr.Repository.TxnManager,
	}
	jobID := mgr.JobManager.Add(ctx, "Migrating scene screenshots to blobs...", t)

	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) MigrateBlobs(ctx context.Context, input MigrateBlobsInput) (string, error) {
	mgr := manager.GetInstance()
	t := &task.MigrateBlobsJob{
		TxnManager: mgr.Database,
		BlobStore:  mgr.Database.Blobs,
		Vacuumer:   mgr.Database,
		DeleteOld:  utils.IsTrue(input.DeleteOld),
	}
	jobID := mgr.JobManager.Add(ctx, "Migrating blobs...", t)

	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) Migrate(ctx context.Context, input manager.MigrateInput) (string, error) {
	mgr := manager.GetInstance()
	t := &task.MigrateJob{
		BackupPath: input.BackupPath,
		Config:     mgr.Config,
		Database:   mgr.Database,
	}

	jobID := mgr.JobManager.Add(ctx, "Migrating database...", t)

	return strconv.Itoa(jobID), nil
}
