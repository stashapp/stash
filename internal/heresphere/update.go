package heresphere

import (
	"context"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
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
func (rs routes) handleDeletePrimaryFile(ctx context.Context, scn *models.Scene, fileDeleter *file.Deleter) (bool, error) {
	err := rs.withTxn(ctx, func(ctx context.Context) error {
		if err := scn.LoadPrimaryFile(ctx, rs.FileFinder); err != nil {
			return err
		}

		ff := scn.Files.Primary()
		if ff != nil {
			if err := file.Destroy(ctx, rs.FileFinder, ff, fileDeleter, true); err != nil {
				return err
			}
		}

		return nil
	})
	return err == nil, err
}
