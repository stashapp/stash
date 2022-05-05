package manager

import (
	"context"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/scene/generate"
)

type sceneScreenshotter struct {
	g *generate.Generator
}

func (ss *sceneScreenshotter) GenerateScreenshot(ctx context.Context, probeResult *ffmpeg.VideoFile, hash string) error {
	return ss.g.Screenshot(ctx, probeResult.Path, hash, probeResult.Width, probeResult.Duration, generate.ScreenshotOptions{})
}

func (ss *sceneScreenshotter) GenerateThumbnail(ctx context.Context, probeResult *ffmpeg.VideoFile, hash string) error {
	return ss.g.Thumbnail(ctx, probeResult.Path, hash, probeResult.Duration, generate.ScreenshotOptions{})
}

func (t *ScanTask) scanScene(ctx context.Context) *models.Scene {
	logError := func(err error) *models.Scene {
		logger.Error(err.Error())
		return nil
	}

	var retScene *models.Scene
	var s *models.Scene

	if err := t.TxnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		var err error
		s, err = r.Scene().FindByPath(t.file.Path())
		return err
	}); err != nil {
		logger.Error(err.Error())
		return nil
	}

	g := &generate.Generator{
		Encoder:     instance.FFMPEG,
		LockManager: instance.ReadLockManager,
		ScenePaths:  instance.Paths.Scene,
	}

	scanner := scene.Scanner{
		Scanner:             scene.FileScanner(&file.FSHasher{}, t.fileNamingAlgorithm, t.calculateMD5),
		StripFileExtension:  t.StripFileExtension,
		FileNamingAlgorithm: t.fileNamingAlgorithm,
		TxnManager:          t.TxnManager,
		Paths:               GetInstance().Paths,
		Screenshotter: &sceneScreenshotter{
			g: g,
		},
		VideoFileCreator: &instance.FFProbe,
		PluginCache:      instance.PluginCache,
		MutexManager:     t.mutexManager,
		UseFileMetadata:  t.UseFileMetadata,
	}

	if s != nil {
		if err := scanner.ScanExisting(ctx, s, t.file); err != nil {
			return logError(err)
		}

		return nil
	}

	var err error
	retScene, err = scanner.ScanNew(ctx, t.file)
	if err != nil {
		return logError(err)
	}

	return retScene
}
