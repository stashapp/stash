package heresphere

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/txn"
)

/*
 * Searches for favorite tag if it exists, otherwise adds it.
 * This adds a tag, which means tags must also be enabled, or it will never be written.
 */
func handleFavoriteTag(ctx context.Context, rs Routes, scn *models.Scene, user HeresphereAuthReq, txnManager txn.Manager, ret *scene.UpdateSet) (bool, error) {
	if err := txn.WithReadTxn(ctx, txnManager, func(ctx context.Context) error {
		var err error
		var favTag *models.Tag

		tagId := config.GetInstance().GetHSPFavoriteTag()
		if favTag, err = rs.Repository.Tag.Find(ctx, tagId); err == nil {
			favTagVal := HeresphereVideoTag{Name: fmt.Sprintf("Tag:%s", favTag.Name)}

			// Do the old switcheroo to figure out how to add the tag
			if *user.IsFavorite {
				if user.Tags == nil {
					user.Tags = &[]HeresphereVideoTag{favTagVal}
				} else {
					*user.Tags = append(*user.Tags, favTagVal)
				}
			} else {
				if user.Tags == nil {
					sceneTags := getVideoTags(rs, ctx, scn)
					user.Tags = &sceneTags
				}

				for i, tag := range *user.Tags {
					if tag.Name == favTagVal.Name {
						*user.Tags = append((*user.Tags)[:i], (*user.Tags)[i+1:]...)
						break
					}
				}
			}

			return nil
		}

		return err
	}); err != nil {
		return false, err
	}

	return true, nil
}

/*
 * This auxiliary function searches for the "favorite" tag
 */
func getVideoFavorite(rs Routes, r *http.Request, scene *models.Scene) bool {
	found := false
	favTag := config.GetInstance().GetHSPFavoriteTag()

	txn.WithReadTxn(r.Context(), rs.TxnManager, func(ctx context.Context) error {
		tag_ids, err := rs.Repository.Tag.FindBySceneID(ctx, scene.ID)
		if err == nil {
			for _, tag := range tag_ids {
				if tag.ID == favTag {
					found = true
					return nil
				}
			}
		}
		return err
	})

	return found
}
