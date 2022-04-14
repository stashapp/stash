package scene

import (
	"context"
	"path/filepath"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
)

// FileDeleter is an extension of file.Deleter that handles deletion of scene files.
type FileDeleter struct {
	file.Deleter

	FileNamingAlgo models.HashAlgorithm
	Paths          *paths.Paths
}

// MarkGeneratedFiles marks for deletion the generated files for the provided scene.
func (d *FileDeleter) MarkGeneratedFiles(scene *models.Scene) error {
	sceneHash := scene.GetHash(d.FileNamingAlgo)

	if sceneHash == "" {
		return nil
	}

	markersFolder := filepath.Join(d.Paths.Generated.Markers, sceneHash)

	exists, _ := fsutil.FileExists(markersFolder)
	if exists {
		if err := d.Dirs([]string{markersFolder}); err != nil {
			return err
		}
	}

	var files []string

	thumbPath := d.Paths.Scene.GetThumbnailScreenshotPath(sceneHash)
	exists, _ = fsutil.FileExists(thumbPath)
	if exists {
		files = append(files, thumbPath)
	}

	normalPath := d.Paths.Scene.GetScreenshotPath(sceneHash)
	exists, _ = fsutil.FileExists(normalPath)
	if exists {
		files = append(files, normalPath)
	}

	streamPreviewPath := d.Paths.Scene.GetVideoPreviewPath(sceneHash)
	exists, _ = fsutil.FileExists(streamPreviewPath)
	if exists {
		files = append(files, streamPreviewPath)
	}

	streamPreviewImagePath := d.Paths.Scene.GetWebpPreviewPath(sceneHash)
	exists, _ = fsutil.FileExists(streamPreviewImagePath)
	if exists {
		files = append(files, streamPreviewImagePath)
	}

	transcodePath := d.Paths.Scene.GetTranscodePath(sceneHash)
	exists, _ = fsutil.FileExists(transcodePath)
	if exists {
		files = append(files, transcodePath)
	}

	spritePath := d.Paths.Scene.GetSpriteImageFilePath(sceneHash)
	exists, _ = fsutil.FileExists(spritePath)
	if exists {
		files = append(files, spritePath)
	}

	vttPath := d.Paths.Scene.GetSpriteVttFilePath(sceneHash)
	exists, _ = fsutil.FileExists(vttPath)
	if exists {
		files = append(files, vttPath)
	}

	heatmapPath := d.Paths.Scene.GetInteractiveHeatmapPath(sceneHash)
	exists, _ = fsutil.FileExists(heatmapPath)
	if exists {
		files = append(files, heatmapPath)
	}

	return d.Files(files)
}

// MarkMarkerFiles deletes generated files for a scene marker with the
// provided scene and timestamp.
func (d *FileDeleter) MarkMarkerFiles(scene *models.Scene, seconds int) error {
	videoPath := d.Paths.SceneMarkers.GetVideoPreviewPath(scene.GetHash(d.FileNamingAlgo), seconds)
	imagePath := d.Paths.SceneMarkers.GetWebpPreviewPath(scene.GetHash(d.FileNamingAlgo), seconds)
	screenshotPath := d.Paths.SceneMarkers.GetScreenshotPath(scene.GetHash(d.FileNamingAlgo), seconds)

	var files []string

	exists, _ := fsutil.FileExists(videoPath)
	if exists {
		files = append(files, videoPath)
	}

	exists, _ = fsutil.FileExists(imagePath)
	if exists {
		files = append(files, imagePath)
	}

	exists, _ = fsutil.FileExists(screenshotPath)
	if exists {
		files = append(files, screenshotPath)
	}

	return d.Files(files)
}

type Destroyer interface {
	Destroy(ctx context.Context, id int) error
}

type MarkerDestroyer interface {
	FindBySceneID(ctx context.Context, sceneID int) ([]*models.SceneMarker, error)
	Destroy(ctx context.Context, id int) error
}

// Destroy deletes a scene and its associated relationships from the
// database.
func Destroy(ctx context.Context, scene *models.Scene, qb Destroyer, mqb MarkerDestroyer, fileDeleter *FileDeleter, deleteGenerated, deleteFile bool) error {
	markers, err := mqb.FindBySceneID(ctx, scene.ID)
	if err != nil {
		return err
	}

	for _, m := range markers {
		if err := DestroyMarker(ctx, scene, m, mqb, fileDeleter); err != nil {
			return err
		}
	}

	if deleteFile {
		if err := fileDeleter.Files([]string{scene.Path}); err != nil {
			return err
		}

		funscriptPath := GetFunscriptPath(scene.Path)
		funscriptExists, _ := fsutil.FileExists(funscriptPath)
		if funscriptExists {
			if err := fileDeleter.Files([]string{funscriptPath}); err != nil {
				return err
			}
		}
	}

	if deleteGenerated {
		if err := fileDeleter.MarkGeneratedFiles(scene); err != nil {
			return err
		}
	}

	if err := qb.Destroy(ctx, scene.ID); err != nil {
		return err
	}

	return nil
}

// DestroyMarker deletes the scene marker from the database and returns a
// function that removes the generated files, to be executed after the
// transaction is successfully committed.
func DestroyMarker(ctx context.Context, scene *models.Scene, sceneMarker *models.SceneMarker, qb MarkerDestroyer, fileDeleter *FileDeleter) error {
	if err := qb.Destroy(ctx, sceneMarker.ID); err != nil {
		return err
	}

	// delete the preview for the marker
	seconds := int(sceneMarker.Seconds)
	return fileDeleter.MarkMarkerFiles(scene, seconds)
}
