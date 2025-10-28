package scene

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

type MigrateSceneScreenshotsInput struct {
	DeleteFiles       bool `json:"deleteFiles"`
	OverwriteExisting bool `json:"overwriteExisting"`
}

type HashFinderCoverUpdater interface {
	FindByChecksum(ctx context.Context, checksum string) ([]*models.Scene, error)
	FindByOSHash(ctx context.Context, oshash string) ([]*models.Scene, error)
	HasCover(ctx context.Context, sceneID int) (bool, error)
	UpdateCover(ctx context.Context, sceneID int, cover []byte) error
}

type ScreenshotMigrator struct {
	Options      MigrateSceneScreenshotsInput
	SceneUpdater HashFinderCoverUpdater
	TxnManager   txn.Manager
}

func (m *ScreenshotMigrator) MigrateScreenshots(ctx context.Context, screenshotPath string) error {
	// find the scene based on the screenshot path
	s, err := m.findScenes(ctx, screenshotPath)
	if err != nil {
		return fmt.Errorf("finding scenes for screenshot: %w", err)
	}

	for _, scene := range s {
		// migrate each scene in its own transaction
		if err := txn.WithTxn(ctx, m.TxnManager, func(ctx context.Context) error {
			return m.migrateSceneScreenshot(ctx, scene, screenshotPath)
		}); err != nil {
			return fmt.Errorf("migrating screenshot for scene %s: %w", scene.DisplayName(), err)
		}
	}

	// if deleteFiles is true, delete the file
	if m.Options.DeleteFiles {
		if err := os.Remove(screenshotPath); err != nil {
			// log and continue
			logger.Errorf("Error deleting screenshot file %s: %v", screenshotPath, err)
		} else {
			logger.Debugf("Deleted screenshot file %s", screenshotPath)
		}

		// also delete the thumb file
		thumbPath := strings.TrimSuffix(screenshotPath, ".jpg") + ".thumb.jpg"
		// ignore errors for thumb files
		if err := os.Remove(thumbPath); err == nil {
			logger.Debugf("Deleted thumb file %s", thumbPath)
		}
	}

	return nil
}

func (m *ScreenshotMigrator) findScenes(ctx context.Context, screenshotPath string) ([]*models.Scene, error) {
	basename := filepath.Base(screenshotPath)
	ext := filepath.Ext(basename)
	basename = basename[:len(basename)-len(ext)]

	// use the basename to determine the hash type
	algo := m.getHashType(basename)

	if algo == "" {
		// log and return
		return nil, fmt.Errorf("could not determine hash type")
	}

	// use the hash type to get the scene
	var ret []*models.Scene
	err := txn.WithReadTxn(ctx, m.TxnManager, func(ctx context.Context) error {
		var err error

		if algo == models.HashAlgorithmOshash {
			// use oshash
			ret, err = m.SceneUpdater.FindByOSHash(ctx, basename)
		} else {
			// use md5
			ret, err = m.SceneUpdater.FindByChecksum(ctx, basename)
		}

		return err
	})

	return ret, err
}

func (m *ScreenshotMigrator) getHashType(basename string) models.HashAlgorithm {
	// if the basename is 16 characters long, must be oshash
	if len(basename) == 16 {
		return models.HashAlgorithmOshash
	}

	// if its 32 characters long, must be md5
	if len(basename) == 32 {
		return models.HashAlgorithmMd5
	}

	// otherwise, it's undefined
	return ""
}

func (m *ScreenshotMigrator) migrateSceneScreenshot(ctx context.Context, scene *models.Scene, screenshotPath string) error {
	if !m.Options.OverwriteExisting {
		// check if the scene has a cover already
		hasCover, err := m.SceneUpdater.HasCover(ctx, scene.ID)
		if err != nil {
			return fmt.Errorf("checking for existing cover: %w", err)
		}

		if hasCover {
			// already has cover, just silently return
			logger.Debugf("Scene %s already has a screenshot, skipping", scene.DisplayName())
			return nil
		}
	}

	// get the data from the file
	data, err := os.ReadFile(screenshotPath)
	if err != nil {
		return fmt.Errorf("reading screenshot file: %w", err)
	}

	if err := m.SceneUpdater.UpdateCover(ctx, scene.ID, data); err != nil {
		return fmt.Errorf("updating scene screenshot: %w", err)
	}

	logger.Infof("Updated screenshot for scene %s from %s", scene.DisplayName(), filepath.Base(screenshotPath))

	return nil
}
