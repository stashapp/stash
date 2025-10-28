package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene/generate"
)

type GenerateCoverTask struct {
	repository   models.Repository
	Scene        models.Scene
	ScreenshotAt *float64
	Overwrite    bool
}

func (t *GenerateCoverTask) GetDescription() string {
	return fmt.Sprintf("Generating cover for %s", t.Scene.GetTitle())
}

func (t *GenerateCoverTask) Start(ctx context.Context) {
	scenePath := t.Scene.Path

	r := t.repository

	var required bool
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		required = t.required(ctx)

		return t.Scene.LoadPrimaryFile(ctx, r.File)
	}); err != nil {
		logger.Error(err)
		return
	}

	if !required {
		return
	}

	videoFile := t.Scene.Files.Primary()
	if videoFile == nil {
		return
	}

	var at float64
	if t.ScreenshotAt == nil {
		at = float64(videoFile.Duration) * 0.2
	} else {
		at = *t.ScreenshotAt
	}

	// we'll generate the screenshot, grab the generated data and set it
	// in the database.

	logger.Debugf("Creating screenshot for %s", scenePath)

	g := generate.Generator{
		Encoder:      instance.FFMpeg,
		FFMpegConfig: instance.Config,
		LockManager:  instance.ReadLockManager,
		ScenePaths:   instance.Paths.Scene,
		Overwrite:    true,
	}

	coverImageData, err := g.Screenshot(context.TODO(), videoFile.Path, videoFile.Width, videoFile.Duration, generate.ScreenshotOptions{
		At: &at,
	})
	if err != nil {
		logger.Errorf("Error generating screenshot: %v", err)
		logErrorOutput(err)
		return
	}

	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		qb := r.Scene
		scenePartial := models.NewScenePartial()

		// update the scene cover table
		if err := qb.UpdateCover(ctx, t.Scene.ID, coverImageData); err != nil {
			return fmt.Errorf("error setting screenshot: %v", err)
		}

		// update the scene with the update date
		_, err = qb.UpdatePartial(ctx, t.Scene.ID, scenePartial)
		if err != nil {
			return fmt.Errorf("error updating scene: %v", err)
		}

		return nil
	}); err != nil && ctx.Err() == nil {
		logger.Error(err.Error())
	}
}

// required returns true if the sprite needs to be generated
// assumes in a transaction
func (t *GenerateCoverTask) required(ctx context.Context) bool {
	if t.Scene.Path == "" {
		return false
	}

	if t.Overwrite {
		return true
	}

	// if the scene has a cover, then we don't need to generate it
	hasCover, err := t.repository.Scene.HasCover(ctx, t.Scene.ID)
	if err != nil {
		logger.Errorf("Error getting cover: %v", err)
		return false
	}

	return !hasCover
}
