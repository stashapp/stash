package heresphere

import (
	"context"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/txn"
)

/*
 * Modifies the scene rating
 */
func updateRating(user HeresphereAuthReq, ret *scene.UpdateSet) (bool, error) {
	rating := models.Rating5To100F(*user.Rating)
	ret.Partial.Rating.Value = rating
	ret.Partial.Rating.Set = true
	return true, nil
}

/*
 * Modifies the scene PlayCount
 */
func updatePlayCount(ctx context.Context, scn *models.Scene, event HeresphereVideoEvent, txnManager txn.Manager, fqb models.SceneReaderWriter) error {
	if per, err := getMinPlayPercent(); err == nil {
		newTime := event.Time / 1000
		file := scn.Files.Primary()

		// TODO: Need temporal memory, we need to track "Open" videos to do this properly
		if scn.PlayCount == 0 && file != nil && newTime/file.Duration > float64(per)/100.0 {
			ret := &scene.UpdateSet{
				ID:      scn.ID,
				Partial: models.NewScenePartial(),
			}
			ret.Partial.PlayCount.Set = true
			ret.Partial.PlayCount.Value = scn.PlayCount + 1

			if err := txn.WithTxn(ctx, txnManager, func(ctx context.Context) error {
				_, err := ret.Update(ctx, fqb)
				return err
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

/*
 * Deletes the scene's primary file
 */
func handleDeletePrimaryFile(ctx context.Context, txnManager txn.Manager, scn *models.Scene, fqb models.FileReaderWriter, fileDeleter *file.Deleter) (bool, error) {
	err := txn.WithTxn(ctx, txnManager, func(ctx context.Context) error {
		if err := scn.LoadPrimaryFile(ctx, fqb); err != nil {
			return err
		}

		ff := scn.Files.Primary()
		if ff != nil {
			if err := file.Destroy(ctx, fqb, ff, fileDeleter, true); err != nil {
				return err
			}
		}

		return nil
	})
	return err == nil, err
}
