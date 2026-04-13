// TODO(audio): update this file

package audio

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/hook"
)

func (s *Service) Create(ctx context.Context, input models.CreateAudioInput) (*models.Audio, error) {
	// title must be set if no files are provided
	if input.Audio.Title == "" && len(input.FileIDs) == 0 {
		return nil, errors.New("title must be set if audio has no files")
	}

	now := time.Now()
	newAudio := *input.Audio
	newAudio.CreatedAt = now
	newAudio.UpdatedAt = now

	// don't pass the file ids since they may be already assigned
	// assign them afterwards
	if err := s.Repository.Create(ctx, &newAudio, nil); err != nil {
		return nil, fmt.Errorf("creating new audio: %w", err)
	}

	if len(input.CustomFields) > 0 {
		if err := s.Repository.SetCustomFields(ctx, newAudio.ID, models.CustomFieldsInput{
			Full: input.CustomFields,
		}); err != nil {
			return nil, fmt.Errorf("setting custom fields on new audio: %w", err)
		}
	}

	for _, f := range input.FileIDs {
		if err := s.AssignFile(ctx, newAudio.ID, f); err != nil {
			return nil, fmt.Errorf("assigning file %d to new audio: %w", f, err)
		}
	}

	if len(input.FileIDs) > 0 {
		// assign the primary to the first
		if _, err := s.Repository.UpdatePartial(ctx, newAudio.ID, models.AudioPartial{
			PrimaryFileID: &input.FileIDs[0],
		}); err != nil {
			return nil, fmt.Errorf("setting primary file on new audio: %w", err)
		}
	}

	// re-find the audio so that it correctly returns file-related fields
	ret, err := s.Repository.Find(ctx, newAudio.ID)
	if err != nil {
		return nil, err
	}

	if len(input.CoverImage) > 0 {
		if err := s.Repository.UpdateCover(ctx, ret.ID, input.CoverImage); err != nil {
			return nil, fmt.Errorf("setting cover on new audio: %w", err)
		}
	}

	s.PluginCache.RegisterPostHooks(ctx, ret.ID, hook.AudioCreatePost, nil, nil)

	// re-find the audio so that it correctly returns file-related fields
	return ret, nil
}
