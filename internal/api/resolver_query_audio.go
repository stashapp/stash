package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) AudioStreams(ctx context.Context, id *string) ([]*manager.AudioStreamEndpoint, error) {
	audioID, err := strconv.Atoi(*id)
	if err != nil {
		return nil, err
	}

	// find the audio
	var audio *models.Audio
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		audio, err = r.repository.Audio.Find(ctx, audioID)

		if audio != nil {
			err = audio.LoadPrimaryFile(ctx, r.repository.File)
		}

		return err
	}); err != nil {
		return nil, err
	}

	if audio == nil {
		return nil, fmt.Errorf("audio with id %d not found", audioID)
	}

	config := manager.GetInstance().Config

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewAudioURLBuilder(baseURL, audio)
	apiKey := config.GetAPIKey()

	return manager.GetAudioStreamPaths(audio, builder.GetStreamURL(apiKey))
}
