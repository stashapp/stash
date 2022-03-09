package manager

import (
	"context"

	"github.com/stashapp/stash/internal/video/encoder"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

type sceneScreenshotter struct {
	ff ffmpeg.FFMpeg
}

func (ss *sceneScreenshotter) Screenshot(input string, options encoder.ScreenshotOptions) error {
	return encoder.Screenshot(ss.ff, input, options)
}

func (t *ScanTask) scanScene() *models.Scene {
	logError := func(err error) *models.Scene {
		logger.Error(err.Error())
		return nil
	}

	var retScene *models.Scene
	var s *models.Scene

	if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		s, err = r.Scene().FindByPath(t.file.Path())
		return err
	}); err != nil {
		logger.Error(err.Error())
		return nil
	}

	scanner := scene.Scanner{
		Scanner:             scene.FileScanner(&file.FSHasher{}, t.fileNamingAlgorithm, t.calculateMD5),
		StripFileExtension:  t.StripFileExtension,
		FileNamingAlgorithm: t.fileNamingAlgorithm,
		Ctx:                 t.ctx,
		TxnManager:          t.TxnManager,
		Paths:               GetInstance().Paths,
		Screenshotter: &sceneScreenshotter{
			ff: instance.FFMPEG,
		},
		VideoFileCreator: &instance.FFProbe,
		PluginCache:      instance.PluginCache,
		MutexManager:     t.mutexManager,
		UseFileMetadata:  t.UseFileMetadata,
	}

	if s != nil {
		if err := scanner.ScanExisting(s, t.file); err != nil {
			return logError(err)
		}

		return nil
	}

	var err error
	retScene, err = scanner.ScanNew(t.file)
	if err != nil {
		return logError(err)
	}

	return retScene
}
