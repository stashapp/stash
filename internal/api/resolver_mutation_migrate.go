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
	db := manager.GetInstance().Database
	t := &task.MigrateSceneScreenshotsJob{
		ScreenshotsPath: manager.GetInstance().Paths.Generated.Screenshots,
		Input: scene.MigrateSceneScreenshotsInput{
			DeleteFiles:       utils.IsTrue(input.DeleteFiles),
			OverwriteExisting: utils.IsTrue(input.OverwriteExisting),
		},
		SceneRepo:  db.Scene,
		TxnManager: db,
	}
	jobID := manager.GetInstance().JobManager.Add(ctx, "Migrating scene screenshots to blobs...", t)

	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) MigrateBlobs(ctx context.Context, input MigrateBlobsInput) (string, error) {
	db := manager.GetInstance().Database
	t := &task.MigrateBlobsJob{
		TxnManager: db,
		BlobStore:  db.Blobs,
		Vacuumer:   db,
		DeleteOld:  utils.IsTrue(input.DeleteOld),
	}
	jobID := manager.GetInstance().JobManager.Add(ctx, "Migrating blobs...", t)

	return strconv.Itoa(jobID), nil
}
