package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *mutationResolver) SavePlaybackPreference(ctx context.Context, streamType string, quality *models.StreamingResolutionEnum) (bool, error) {
	userAgent, _ := ctx.Value(UserAgentCtxKey).(string)
	if userAgent == "" {
		return false, nil
	}

	st := models.PlaybackStreamType(streamType)
	if !st.IsValid() {
		return false, fmt.Errorf("invalid stream_type: %s", streamType)
	}

	pd := &models.PlaybackDefault{
		UserAgentPattern: userAgent,
		Priority:         100,
		StreamType:       st,
		Quality:          quality,
	}

	if err := r.withWriteTxn(ctx, func(ctx context.Context) error {
		return r.repository.PlaybackDefault.Upsert(ctx, pd)
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) PlaybackDefaultCreate(ctx context.Context, input models.PlaybackDefaultCreateInput) (*models.PlaybackDefault, error) {
	streamType := models.PlaybackStreamType(input.StreamType)
	if !streamType.IsValid() {
		return nil, fmt.Errorf("invalid stream_type: %s", input.StreamType)
	}

	pd := &models.PlaybackDefault{
		UserAgentPattern: input.UserAgentPattern,
		Priority:         input.Priority,
		StreamType:       streamType,
		Quality:          input.Quality,
	}

	if err := r.withWriteTxn(ctx, func(ctx context.Context) error {
		return r.repository.PlaybackDefault.Create(ctx, pd)
	}); err != nil {
		return nil, err
	}

	return pd, nil
}

func (r *mutationResolver) PlaybackDefaultUpdate(ctx context.Context, input models.PlaybackDefaultUpdateInput) (*models.PlaybackDefault, error) {
	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	var pd *models.PlaybackDefault
	if err := r.withWriteTxn(ctx, func(ctx context.Context) error {
		var err error
		pd, err = r.repository.PlaybackDefault.Find(ctx, id)
		if err != nil {
			return err
		}
		if pd == nil {
			return fmt.Errorf("playback default %d not found", id)
		}

		if input.UserAgentPattern != nil {
			pd.UserAgentPattern = *input.UserAgentPattern
		}
		if input.Priority != nil {
			pd.Priority = *input.Priority
		}
		if input.StreamType != nil {
			st := models.PlaybackStreamType(*input.StreamType)
			if !st.IsValid() {
				return fmt.Errorf("invalid stream_type: %s", *input.StreamType)
			}
			pd.StreamType = st
		}
		if input.Quality != nil {
			pd.Quality = input.Quality
		}

		return r.repository.PlaybackDefault.Update(ctx, pd)
	}); err != nil {
		return nil, err
	}

	return pd, nil
}

func (r *mutationResolver) PlaybackDefaultDestroy(ctx context.Context, id string) (bool, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return false, err
	}

	if err := r.withWriteTxn(ctx, func(ctx context.Context) error {
		return r.repository.PlaybackDefault.Destroy(ctx, idInt)
	}); err != nil {
		return false, err
	}

	return true, nil
}
