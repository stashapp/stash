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
