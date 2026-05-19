package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) SceneStreams(ctx context.Context, id *string) ([]*manager.SceneStreamEndpoint, error) {
	sceneID, err := strconv.Atoi(*id)
	if err != nil {
		return nil, err
	}

	// find the scene
	var scene *models.Scene
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		scene, err = r.repository.Scene.Find(ctx, sceneID)

		if scene != nil {
			err = scene.LoadPrimaryFile(ctx, r.repository.File)
		}

		return err
	}); err != nil {
		return nil, err
	}

	if scene == nil {
		return nil, fmt.Errorf("scene with id %d not found", sceneID)
	}

	config := manager.GetInstance().Config

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	userAgent, _ := ctx.Value(UserAgentCtxKey).(string)
	builder := urlbuilders.NewSceneURLBuilder(baseURL, scene)
	apiKey := config.GetAPIKey()

	endpoints, err := manager.GetSceneStreamPaths(scene, builder.GetStreamURL(apiKey), config.GetMaxStreamingTranscodeSize())
	if err != nil {
		return nil, err
	}

	if err := manager.MarkDefaultStream(ctx, endpoints, userAgent, r.repository.PlaybackDefault); err != nil {
		return nil, err
	}

	return endpoints, nil
}
