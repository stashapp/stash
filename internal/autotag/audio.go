package autotag

import (
	"context"
	"slices"

	"github.com/stashapp/stash/pkg/audio"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
)

type AudioFinderUpdater interface {
	models.AudioQueryer
	models.AudioUpdater
}

type AudioPerformerUpdater interface {
	models.PerformerIDLoader
	models.AudioUpdater
}

type AudioTagUpdater interface {
	models.TagIDLoader
	models.AudioUpdater
}

func getAudioFileTagger(s *models.Audio, cache *match.Cache) tagger {
	return tagger{
		ID:    s.ID,
		Type:  "audio",
		Name:  s.DisplayName(),
		Path:  s.Path,
		cache: cache,
	}
}

// AudioPerformers tags the provided audio with performers whose name matches the audio's path.
func AudioPerformers(ctx context.Context, s *models.Audio, rw AudioPerformerUpdater, performerReader models.PerformerAutoTagQueryer, cache *match.Cache) error {
	t := getAudioFileTagger(s, cache)

	return t.tagPerformers(ctx, performerReader, func(subjectID, otherID int) (bool, error) {
		if err := s.LoadPerformerIDs(ctx, rw); err != nil {
			return false, err
		}
		existing := s.PerformerIDs.List()

		if slices.Contains(existing, otherID) {
			return false, nil
		}

		if err := audio.AddPerformer(ctx, rw, s, otherID); err != nil {
			return false, err
		}

		return true, nil
	})
}

// AudioStudios tags the provided audio with the first studio whose name matches the audio's path.
//
// Audios will not be tagged if studio is already set.
func AudioStudios(ctx context.Context, s *models.Audio, rw AudioFinderUpdater, studioReader models.StudioAutoTagQueryer, cache *match.Cache) error {
	if s.StudioID != nil {
		// don't modify
		return nil
	}

	t := getAudioFileTagger(s, cache)

	return t.tagStudios(ctx, studioReader, func(subjectID, otherID int) (bool, error) {
		return addAudioStudio(ctx, rw, s, otherID)
	})
}

// AudioTags tags the provided audio with tags whose name matches the audio's path.
func AudioTags(ctx context.Context, s *models.Audio, rw AudioTagUpdater, tagReader models.TagAutoTagQueryer, cache *match.Cache) error {
	t := getAudioFileTagger(s, cache)

	return t.tagTags(ctx, tagReader, func(subjectID, otherID int) (bool, error) {
		if err := s.LoadTagIDs(ctx, rw); err != nil {
			return false, err
		}
		existing := s.TagIDs.List()

		if slices.Contains(existing, otherID) {
			return false, nil
		}

		if err := audio.AddTag(ctx, rw, s, otherID); err != nil {
			return false, err
		}

		return true, nil
	})
}
