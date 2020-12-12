package api

import (
	"context"
	"errors"
	"strconv"

	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

func (r *queryResolver) SceneStreams(ctx context.Context, id *string) ([]*models.SceneStreamEndpoint, error) {
	// find the scene
	qb := sqlite.NewSceneQueryBuilder()
	idInt, _ := strconv.Atoi(*id)
	scene, err := qb.Find(idInt)

	if err != nil {
		return nil, err
	}

	if scene == nil {
		return nil, errors.New("nil scene")
	}

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewSceneURLBuilder(baseURL, scene.ID)

	return manager.GetSceneStreamPaths(scene, builder.GetStreamURL())
}
