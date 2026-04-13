// TODO(audio): update this file

package audio

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/file/video"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/plugin/hook"
	"github.com/stashapp/stash/pkg/txn"
)

var (
	ErrNotVideoFile = errors.New("not a video file")

	// fingerprint types to match with
	// only try to match by data fingerprints, _not_ perceptual fingerprints
	matchableFingerprintTypes = []string{models.FingerprintTypeOshash, models.FingerprintTypeMD5}
)

type ScanCreatorUpdater interface {
	FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Audio, error)
	FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Audio, error)
	GetFiles(ctx context.Context, relatedID int) ([]*models.VideoFile, error)

	Create(ctx context.Context, newAudio *models.Audio, fileIDs []models.FileID) error
	UpdatePartial(ctx context.Context, id int, updatedAudio models.AudioPartial) (*models.Audio, error)
	AddFileID(ctx context.Context, id int, fileID models.FileID) error
}

type ScanGalleryFinderUpdater interface {
	FindByPath(ctx context.Context, p string) ([]*models.Gallery, error)
	AddAudioIDs(ctx context.Context, galleryID int, audioIDs []int) error
}

type ScanGenerator interface {
	Generate(ctx context.Context, s *models.Audio, f *models.VideoFile) error
}

type ScanHandler struct {
	CreatorUpdater       ScanCreatorUpdater
	GalleryFinderUpdater ScanGalleryFinderUpdater

	ScanGenerator  ScanGenerator
	CaptionUpdater video.CaptionUpdater
	PluginCache    *plugin.Cache

	FileNamingAlgorithm models.HashAlgorithm
	Paths               *paths.Paths
}

func (h *ScanHandler) validate() error {
	if h.CreatorUpdater == nil {
		return errors.New("CreatorUpdater is required")
	}
	if h.ScanGenerator == nil {
		return errors.New("ScanGenerator is required")
	}
	if h.CaptionUpdater == nil {
		return errors.New("CaptionUpdater is required")
	}
	if !h.FileNamingAlgorithm.IsValid() {
		return errors.New("FileNamingAlgorithm is required")
	}
	if h.Paths == nil {
		return errors.New("Paths is required")
	}

	return nil
}

func (h *ScanHandler) Handle(ctx context.Context, f models.File, oldFile models.File) error {
	if err := h.validate(); err != nil {
		return err
	}

	videoFile, ok := f.(*models.VideoFile)
	if !ok {
		return ErrNotVideoFile
	}

	if oldFile != nil {
		if err := video.CleanCaptions(ctx, videoFile, nil, h.CaptionUpdater); err != nil {
			return fmt.Errorf("cleaning captions: %w", err)
		}
	}

	// try to match the file to a audio
	existing, err := h.CreatorUpdater.FindByFileID(ctx, f.Base().ID)
	if err != nil {
		return fmt.Errorf("finding existing audio: %w", err)
	}

	if len(existing) == 0 {
		// try also to match file by fingerprints
		existing, err = h.CreatorUpdater.FindByFingerprints(ctx, videoFile.Fingerprints.Filter(matchableFingerprintTypes...))
		if err != nil {
			return fmt.Errorf("finding existing audio by fingerprints: %w", err)
		}
	}

	if len(existing) > 0 {
		updateExisting := oldFile != nil
		if err := h.associateExisting(ctx, existing, videoFile, updateExisting); err != nil {
			return err
		}
	} else {
		// create a new audio
		newAudio := models.NewAudio()

		logger.Infof("%s doesn't exist. Creating new audio...", f.Base().Path)

		if err := h.CreatorUpdater.Create(ctx, &newAudio, []models.FileID{videoFile.ID}); err != nil {
			return fmt.Errorf("creating new audio: %w", err)
		}

		h.PluginCache.RegisterPostHooks(ctx, newAudio.ID, hook.AudioCreatePost, nil, nil)

		existing = []*models.Audio{&newAudio}
	}

	if oldFile != nil {
		// migrate hashes from the old file to the new
		oldHash := GetHash(oldFile, h.FileNamingAlgorithm)
		newHash := GetHash(f, h.FileNamingAlgorithm)

		if oldHash != "" && newHash != "" && oldHash != newHash {
			MigrateHash(h.Paths, oldHash, newHash)
		}
	}

	if err := h.associateGallery(ctx, existing, f); err != nil {
		return err
	}

	// do this after the commit so that cover generation doesn't hold up the transaction
	txn.AddPostCommitHook(ctx, func(ctx context.Context) {
		for _, s := range existing {
			if err := h.ScanGenerator.Generate(ctx, s, videoFile); err != nil {
				// just log if cover generation fails. We can try again on rescan
				logger.Errorf("Error generating content for %s: %v", videoFile.Path, err)
			}
		}
	})

	return nil
}

func (h *ScanHandler) associateExisting(ctx context.Context, existing []*models.Audio, f *models.VideoFile, updateExisting bool) error {
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

func (h *ScanHandler) associateGallery(ctx context.Context, existing []*models.Audio, f models.File) error {
	audioIDs := make([]int, len(existing))
	for i, s := range existing {
		audioIDs[i] = s.ID
	}

	path := f.Base().Path
	zipPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".zip"

	// find galleries with a file that matches
	galleries, err := h.GalleryFinderUpdater.FindByPath(ctx, zipPath)
	if err != nil {
		return err
	}

	for _, gallery := range galleries {
		// found related Audio
		logger.Infof("associate: Audio %s is related to gallery: %d", path, gallery.ID)
		if err := h.GalleryFinderUpdater.AddAudioIDs(ctx, gallery.ID, audioIDs); err != nil {
			return err
		}
	}

	return nil
}
