package api

import (
	"context"
	"github.com/stashapp/stash/internal/models"
	"strconv"
)

func (r *queryResolver) FindScene(ctx context.Context, id *string, checksum *string) (*models.Scene, error) {
	qb := models.NewSceneQueryBuilder()
	idInt, _ := strconv.Atoi(*id)
	var scene *models.Scene
	var err error
	if id != nil {
		scene, err = qb.Find(idInt)
	} else if checksum != nil {
		scene, err = qb.FindByChecksum(*checksum)
	}
	return scene, err
}

func (r *queryResolver) FindScenes(ctx context.Context, scene_filter *models.SceneFilterType, scene_ids []int, filter *models.FindFilterType) (models.FindScenesResultType, error) {
	qb := models.NewSceneQueryBuilder()
	scenes, total := qb.Query(scene_filter, filter)
	return models.FindScenesResultType{
		Count: total,
		Scenes: scenes,
	}, nil
}