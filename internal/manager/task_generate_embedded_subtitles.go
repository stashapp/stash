package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/file/video"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

// ExtractEmbeddedSubtitlesTask extracts embedded subtitles from a video file.
type ExtractEmbeddedSubtitlesTask struct {
	repository          models.Repository
	Scene               models.Scene
	fileNamingAlgorithm models.HashAlgorithm
	VideoFile           *models.VideoFile
}

func (t *ExtractEmbeddedSubtitlesTask) GetDescription() string {
	return fmt.Sprintf("Extracting embedded subtitles for %s", t.VideoFile.Path)
}

func (t *ExtractEmbeddedSubtitlesTask) Start(ctx context.Context) {
	config := instance.Config

	ffmpegPath := GetInstance().FFMpeg.Path()
	ffprobePath := GetInstance().FFProbe.Path()
	generatedPath := config.GetGeneratedPath()
	videoFile := t.VideoFile

	// Get the scene hash
	sceneHash := t.Scene.GetHash(t.fileNamingAlgorithm)
	logger.Debugf("Extracting embedded subtitles for scene with hash: %s", sceneHash)

	logger.Infof("Extracting embedded subtitles from %s", videoFile.Path)

	// Extract embedded subtitles (pass the scene hash)
	extractedCaptions, err := video.ExtractEmbeddedSubtitles(ctx, videoFile, generatedPath, ffmpegPath, ffprobePath, sceneHash)
	if err != nil {
		logger.Errorf("Error extracting embedded subtitles from %s: %v", videoFile.Path, err)
		return
	}

	if len(extractedCaptions) == 0 {
		logger.Infof("No embedded subtitles found in %s", videoFile.Path)
		return
	}

	// Associate the extracted captions with the video file
	if err := t.repository.WithTxn(ctx, func(ctx context.Context) error {
		fileID := videoFile.ID

		// Get existing captions for the file
		existingCaptions, err := t.repository.File.GetCaptions(ctx, fileID)
		if err != nil {
			return fmt.Errorf("error getting existing captions: %v", err)
		}

		// Merge existing captions with newly extracted ones
		// (avoiding duplicates based on language code and caption type)
		allCaptions := make([]*models.VideoCaption, len(existingCaptions))
		copy(allCaptions, existingCaptions)
		for _, caption := range extractedCaptions {
			duplicate := false
			for _, existing := range existingCaptions {
				if caption.LanguageCode == existing.LanguageCode && caption.CaptionType == existing.CaptionType {
					// Skip if this language already exists
					duplicate = true
					break
				}
			}

			if !duplicate {
				// Only save the filename (don't save the absolute path)
				// This way, the FileDeleter can find and delete the subtitle files when the scene is deleted
				allCaptions = append(allCaptions, caption)
				logger.Debugf("Added subtitle: lang=%s, type=%s, filename=%s",
					caption.LanguageCode, caption.CaptionType, caption.Filename)
			}
		}

		// Update the captions in the database
		if err := t.repository.File.UpdateCaptions(ctx, fileID, allCaptions); err != nil {
			return fmt.Errorf("error updating captions: %v", err)
		}

		return nil
	}); err != nil {
		logger.Errorf("Error associating extracted captions with video file: %v", err)
		return
	}

	logger.Infof("Successfully extracted and associated %d subtitle tracks from %s", len(extractedCaptions), videoFile.Path)
}
