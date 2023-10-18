package heresphere

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

/*
 * Searches for favorite tag if it exists, otherwise adds it.
 * This adds a tag, which means tags must also be enabled, or it will never be written.
 */
func (rs routes) handleFavoriteTag(ctx context.Context, scn *models.Scene, user *HeresphereAuthReq, ret *scene.UpdateSet) (bool, error) {
	tagID := config.GetInstance().GetHSPFavoriteTag()

	favTag, err := func() (*models.Tag, error) {
		var tag *models.Tag
		var err error
		err = rs.withReadTxn(ctx, func(ctx context.Context) error {
			tag, err = rs.TagFinder.Find(ctx, tagID)
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

	if user.Tags == nil {
		sceneTags := rs.getVideoTags(ctx, scn)
		user.Tags = &sceneTags
	}

	if *user.IsFavorite {
		*user.Tags = append(*user.Tags, favTagVal)
	} else {
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
func (rs routes) getVideoFavorite(r *http.Request, scene *models.Scene) bool {
	tagIDs, err := func() ([]*models.Tag, error) {
		var tags []*models.Tag
		var err error
		err = rs.withReadTxn(r.Context(), func(ctx context.Context) error {
			tags, err = rs.TagFinder.FindBySceneID(ctx, scene.ID)
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
