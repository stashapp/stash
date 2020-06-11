package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) IsSceneStreamable(ctx context.Context, id *string, supportedVideoCodecs []string) (bool, error) {
	// find the scene
	qb := models.NewSceneQueryBuilder()
	idInt, _ := strconv.Atoi(*id)
	scene, err := qb.Find(idInt)

	if err != nil {
		return false, err
	}

	return manager.IsStreamable(scene, supportedVideoCodecs)
}
