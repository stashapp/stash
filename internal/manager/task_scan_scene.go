package manager

import (
	"context"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

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
		Screenshotter:       &instance.FFMPEG,
		VideoFileCreator:    &instance.FFProbe,
		PluginCache:         instance.PluginCache,
		MutexManager:        t.mutexManager,
		UseFileMetadata:     t.UseFileMetadata,
		DetectCaptions:      t.DetectCaptions,
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

// associates captions to scene/s with the same basename
func (t *ScanTask) associateCaptions(ctx context.Context) {
	vExt := config.GetInstance().GetVideoExtensions()
	captionPath := t.file.Path()
	captionLang := scene.GetCaptionsLangFromPath(captionPath)

	relatedFiles := scene.GenerateCaptionCandidates(captionPath, vExt)
	if err := t.TxnManager.WithTxn(ctx, func(r models.Repository) error {
		var err error
		sqb := r.Scene()

		for _, scenePath := range relatedFiles {
			s, er := sqb.FindByPath(scenePath)

			if er != nil {
				logger.Errorf("Error searching for scene %s: %v", scenePath, er)
				continue
			}
			if s != nil { // found related Scene
				logger.Debugf("Matched captions to scene %s", s.Path)
				captions, er := sqb.GetCaptions(s.ID)
				if er == nil {
					if !scene.IsLangInCaptions(captionLang, captions) { // only update captions if language code is not present
						newCaptions := scene.AddLangToCaptions(captionLang, captions)
						er = sqb.UpdateCaptions(s.ID, newCaptions)
						if er == nil {
							logger.Debugf("Updated captions for scene %s. Added %s", s.Path, captionLang)
						}
					}
				}
			}
		}
		return err
	}); err != nil {
		logger.Error(err.Error())
	}
}
