package manager

import (
	"context"
	"path/filepath"

	"github.com/stashapp/stash/internal/manager/config"
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

	if err := t.TxnManager.WithTxn(ctx, func(ctx context.Context) error {
		var err error
		s, err = t.TxnManager.Scene.FindByPath(ctx, t.file.Path())
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
		CreatorUpdater:      t.TxnManager.Scene,
		Paths:               GetInstance().Paths,
		CaseSensitiveFs:     t.CaseSensitiveFs,
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

// associates captions to scene/s with the same basename
func (t *ScanTask) associateCaptions(ctx context.Context) {
	vExt := config.GetInstance().GetVideoExtensions()
	captionPath := t.file.Path()
	captionLang := scene.GetCaptionsLangFromPath(captionPath)

	relatedFiles := scene.GenerateCaptionCandidates(captionPath, vExt)
	if err := t.TxnManager.WithTxn(ctx, func(ctx context.Context) error {
		var err error
		sqb := t.TxnManager.Scene

		for _, scenePath := range relatedFiles {
			s, er := sqb.FindByPath(ctx, scenePath)

			if er != nil {
				logger.Errorf("Error searching for scene %s: %v", scenePath, er)
				continue
			}
			if s != nil { // found related Scene
				logger.Debugf("Matched captions to scene %s", s.Path)
				captions, er := sqb.GetCaptions(ctx, s.ID)
				if er == nil {
					fileExt := filepath.Ext(captionPath)
					ext := fileExt[1:]
					if !scene.IsLangInCaptions(captionLang, ext, captions) { // only update captions if language code is not present
						newCaption := &models.SceneCaption{
							LanguageCode: captionLang,
							Filename:     filepath.Base(captionPath),
							CaptionType:  ext,
						}
						captions = append(captions, newCaption)
						er = sqb.UpdateCaptions(ctx, s.ID, captions)
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
