package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/file/video"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GenerateEmbeddedCaptionsTask struct {
	repository models.Repository
	Scene      models.Scene
	Overwrite  bool
}

func (t *GenerateEmbeddedCaptionsTask) GetDescription() string {
	return fmt.Sprintf("Generating embedded captions for %s", t.Scene.Path)
}

func (t *GenerateEmbeddedCaptionsTask) Start(ctx context.Context) {
	if !t.required() {
		return
	}

	videoFile := t.Scene.Files.Primary()
	var captions []*models.VideoCaption

	r := t.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		var err error
		captions, err = r.File.GetCaptions(ctx, videoFile.ID)
		return err
	}); err != nil {
		if ctx.Err() == nil {
			logger.Errorf("Error getting captions for %s: %v", videoFile.Path, err)
		}
		return
	}

	extractor := video.EmbeddedCaptionExtractor{
		FFProbe: instance.FFProbe,
		FFMpeg:  instance.FFMpeg,
	}

	generatedCaptions, err := extractor.Extract(ctx, videoFile, instance.Paths.Generated.Captions, videoFile.ID.String(), captions, t.Overwrite)
	if err != nil {
		if ctx.Err() == nil {
			logger.Errorf("Error extracting embedded captions for %s: %v", videoFile.Path, err)
			logErrorOutput(err)
		}
		return
	}

	if len(generatedCaptions) == 0 {
		return
	}

	mergedCaptions := video.MergeGeneratedCaptions(captions, generatedCaptions)
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		return r.File.UpdateCaptions(ctx, videoFile.ID, mergedCaptions)
	}); err != nil && ctx.Err() == nil {
		logger.Errorf("Error updating embedded captions for %s: %v", videoFile.Path, err)
	}
}

func (t *GenerateEmbeddedCaptionsTask) required() bool {
	if t.Scene.Path == "" || !t.Scene.Files.PrimaryLoaded() {
		return false
	}

	videoFile := t.Scene.Files.Primary()
	return videoFile != nil && videoFile.ZipFileID == nil
}
