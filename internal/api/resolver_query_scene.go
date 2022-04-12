package api

import (
	"context"
	"errors"
	"strconv"

	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) SceneStreams(ctx context.Context, id *string) ([]*manager.SceneStreamEndpoint, error) {
	// find the scene
	var scene *models.Scene
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		idInt, _ := strconv.Atoi(*id)
		var err error
		scene, err = repo.Scene().Find(idInt)
		return err
	}); err != nil {
		return nil, err
	}

	if scene == nil {
		return nil, errors.New("nil scene")
	}

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewSceneURLBuilder(baseURL, scene.ID)

	return manager.GetSceneStreamPaths(scene, builder.GetStreamURL(), config.GetInstance().GetMaxStreamingTranscodeSize())
}
