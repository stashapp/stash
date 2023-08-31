package heresphere

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/txn"
)

/*
 * Searches for favorite tag if it exists, otherwise adds it.
 * This adds a tag, which means tags must also be enabled, or it will never be written.
 */
func handleFavoriteTag(ctx context.Context, rs Routes, scn *models.Scene, user *HeresphereAuthReq, txnManager txn.Manager, ret *scene.UpdateSet) (bool, error) {
	tagID := config.GetInstance().GetHSPFavoriteTag()

	favTag, err := func() (*models.Tag, error) {
		var tag *models.Tag
		var err error
		err = txn.WithReadTxn(ctx, txnManager, func(ctx context.Context) error {
			tag, err = rs.Repository.Tag.Find(ctx, tagID)
			return err
		})
		return tag, err
	}()

	if err != nil {
		logger.Errorf("Heresphere handleFavoriteTag Tag.Find error: %s\n", err.Error())
		return false, err
	}

	if favTag == nil {
		return false, nil
	}

	favTagVal := HeresphereVideoTag{Name: fmt.Sprintf("Tag:%s", favTag.Name)}

	if *user.IsFavorite {
		if user.Tags == nil {
			user.Tags = &[]HeresphereVideoTag{favTagVal}
		} else {
			*user.Tags = append(*user.Tags, favTagVal)
		}
	} else {
		if user.Tags == nil {
			sceneTags := getVideoTags(ctx, rs, scn)
			user.Tags = &sceneTags
		}

		for i, tag := range *user.Tags {
			if tag.Name == favTagVal.Name {
				*user.Tags = append((*user.Tags)[:i], (*user.Tags)[i+1:]...)
				break
			}
		}
	}

	return true, nil
}

/*
 * This auxiliary function searches for the "favorite" tag
 */
func getVideoFavorite(rs Routes, r *http.Request, scene *models.Scene) bool {
	tagIDs, err := func() ([]*models.Tag, error) {
		var tags []*models.Tag
		var err error
		err = txn.WithReadTxn(r.Context(), rs.TxnManager, func(ctx context.Context) error {
			tags, err = rs.Repository.Tag.FindBySceneID(ctx, scene.ID)
			return err
		})
		return tags, err
	}()

	if err != nil {
		logger.Errorf("Heresphere getVideoFavorite error: %s\n", err.Error())
		return false
	}

	favTag := config.GetInstance().GetHSPFavoriteTag()
	for _, tag := range tagIDs {
		if tag.ID == favTag {
			return true
		}
	}

	return false
}
