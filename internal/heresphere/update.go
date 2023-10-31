package heresphere

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/scene"
)

/*
 * Modifies the scene rating
 */
func (rs routes) updateRating(user HeresphereAuthReq, ret *scene.UpdateSet) (bool, error) {
	rating := models.Rating5To100F(*user.Rating)
	ret.Partial.Rating.Value = rating
	ret.Partial.Rating.Set = true
	return true, nil
}

/*
 * Modifies the scene PlayCount
 */
func (rs routes) updatePlayCount(ctx context.Context, scn *models.Scene, event HeresphereVideoEvent) (bool, error) {
	if per, err := getMinPlayPercent(); err == nil {
		newTime := event.Time / 1000
		file := scn.Files.Primary()

		if file != nil && newTime/file.Duration > float64(per)/100.0 {
			ret := &scene.UpdateSet{
				ID:      scn.ID,
				Partial: models.NewScenePartial(),
			}
			ret.Partial.PlayCount.Set = true
			ret.Partial.PlayCount.Value = scn.PlayCount + 1

			err := rs.withTxn(ctx, func(ctx context.Context) error {
				_, err := ret.Update(ctx, rs.SceneFinder)
				return err
			})
			return err == nil, err
		}
	}

	return false, nil
}

/*
 * Deletes the scene's primary file
 */
func (rs routes) handleDeleteScene(ctx context.Context, scn *models.Scene) (bool, error) {
	err := rs.withTxn(ctx, func(ctx context.Context) error {
		// Construct scene deletion
		deleteFile := true
		deleteGenerated := true
		input := models.ScenesDestroyInput{
			Ids:             []string{strconv.Itoa(scn.ID)},
			DeleteFile:      &deleteFile,
			DeleteGenerated: &deleteGenerated,
		}

		// Construct file deleter
		fileNamingAlgo := manager.GetInstance().Config.GetVideoFileNamingAlgorithm()
		fileDeleter := &scene.FileDeleter{
			Deleter:        file.NewDeleter(),
			FileNamingAlgo: fileNamingAlgo,
			Paths:          manager.GetInstance().Paths,
		}

		// Kill running streams
		manager.KillRunningStreams(scn, fileNamingAlgo)

		// Destroy scene
		if err := rs.SceneService.Destroy(ctx, scn, fileDeleter, deleteGenerated, deleteFile); err != nil {
			fileDeleter.Rollback()
			return err
		}

		// Commit deletion
		fileDeleter.Commit()

		// Plugin callback
		rs.HookExecutor.ExecutePostHooks(ctx, scn.ID, plugin.SceneDestroyPost, plugin.ScenesDestroyInput{
			ScenesDestroyInput: input,
			Checksum:           scn.Checksum,
			OSHash:             scn.OSHash,
			Path:               scn.Path,
		}, nil)

		return nil
	})
	return err == nil, err
}
