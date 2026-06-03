package audio

import (
	"context"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/plugin/hook"
)

var (
	ErrNotAudioFile = errors.New("not a audio file")

	// fingerprint types to match with
	// only try to match by data fingerprints, _not_ perceptual fingerprints
	matchableFingerprintTypes = []string{models.FingerprintTypeMD5}
)

type ScanCreatorUpdater interface {
	FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Audio, error)
	FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Audio, error)
	GetFiles(ctx context.Context, relatedID int) ([]*models.AudioFile, error)

	Create(ctx context.Context, newAudio *models.Audio, fileIDs []models.FileID) error
	UpdatePartial(ctx context.Context, id int, updatedAudio models.AudioPartial) (*models.Audio, error)
	AddFileID(ctx context.Context, id int, fileID models.FileID) error
}

type ScanGalleryFinderUpdater interface {
	FindByPath(ctx context.Context, p string) ([]*models.Gallery, error)
	AddAudioIDs(ctx context.Context, galleryID int, audioIDs []int) error
}

type ScanGenerator interface {
	Generate(ctx context.Context, s *models.Audio, f *models.AudioFile) error
}

type ScanHandler struct {
	CreatorUpdater ScanCreatorUpdater

	PluginCache *plugin.Cache
	Paths       *paths.Paths
}

func (h *ScanHandler) validate() error {
	if h.CreatorUpdater == nil {
		return errors.New("internal error:CreatorUpdater is required")
	}
	if h.Paths == nil {
		return errors.New("internal error:Paths is required")
	}

	return nil
}

func (h *ScanHandler) Handle(ctx context.Context, f models.File, oldFile models.File) error {
	if err := h.validate(); err != nil {
		return err
	}

	AudioFile, ok := f.(*models.AudioFile)
	if !ok {
		return ErrNotAudioFile
	}

	// try to match the file to a audio
	existing, err := h.CreatorUpdater.FindByFileID(ctx, f.Base().ID)
	if err != nil {
		return fmt.Errorf("finding existing audio: %w", err)
	}

	if len(existing) == 0 {
		// try also to match file by fingerprints
		existing, err = h.CreatorUpdater.FindByFingerprints(ctx, AudioFile.Fingerprints.Filter(matchableFingerprintTypes...))
		if err != nil {
			return fmt.Errorf("finding existing audio by fingerprints: %w", err)
		}
	}

	if len(existing) > 0 {
		updateExisting := oldFile != nil
		if err := h.associateExisting(ctx, existing, AudioFile, updateExisting); err != nil {
			return err
		}
	} else {
		// create a new audio
		newAudio := models.NewAudio()

		logger.Infof("%s doesn't exist. Creating new audio...", f.Base().Path)

		if err := h.CreatorUpdater.Create(ctx, &newAudio, []models.FileID{AudioFile.ID}); err != nil {
			return fmt.Errorf("creating new audio: %w", err)
		}

		h.PluginCache.RegisterPostHooks(ctx, newAudio.ID, hook.AudioCreatePost, nil, nil)
	}

	if oldFile != nil {
		// migrate hashes from the old file to the new
		oldHash := GetHash(oldFile, models.HashAlgorithmMd5)
		newHash := GetHash(f, models.HashAlgorithmMd5)

		if oldHash != "" && newHash != "" && oldHash != newHash {
			MigrateHash(h.Paths, oldHash, newHash)
		}
	}

	return nil
}

func (h *ScanHandler) associateExisting(ctx context.Context, existing []*models.Audio, f *models.AudioFile, updateExisting bool) error {
	for _, s := range existing {
		if err := s.LoadFiles(ctx, h.CreatorUpdater); err != nil {
			return err
		}

		found := false
		for _, sf := range s.Files.List() {
			if sf.ID == f.ID {
				found = true
				break
			}
		}

		if !found {
			logger.Infof("Adding %s to audio %s", f.Path, s.DisplayName())

			if err := h.CreatorUpdater.AddFileID(ctx, s.ID, f.ID); err != nil {
				return fmt.Errorf("adding file to audio: %w", err)
			}
		}

		if !found || updateExisting {
			// update updated_at time when file association or content changes
			audioPartial := models.NewAudioPartial()
			if _, err := h.CreatorUpdater.UpdatePartial(ctx, s.ID, audioPartial); err != nil {
				return fmt.Errorf("updating audio: %w", err)
			}

			h.PluginCache.RegisterPostHooks(ctx, s.ID, hook.AudioUpdatePost, nil, nil)
		}
	}

	return nil
}
