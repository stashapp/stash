package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

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
	builder := urlbuilders.NewSceneURLBuilder(baseURL, scene)
	expires := time.Now().Add(24 * time.Hour)
	signedURL, err := builder.GetSignedStreamURL(config.GetJWTSignKey(), expires)
	if err != nil {
		// fallback to api key
		apiKey := config.GetAPIKey()
		streamURL := builder.GetStreamURL(apiKey)
		return manager.GetSceneStreamPaths(scene, streamURL, config.GetMaxStreamingTranscodeSize())
	}

	u, err := url.Parse(signedURL)
	if err != nil {
		// fallback
		apiKey := config.GetAPIKey()
		streamURL := builder.GetStreamURL(apiKey)
		return manager.GetSceneStreamPaths(scene, streamURL, config.GetMaxStreamingTranscodeSize())
	}

	return manager.GetSceneStreamPaths(scene, u, config.GetMaxStreamingTranscodeSize())
}
