package scene

import (
	"path/filepath"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
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

	exists, _ := utils.FileExists(markersFolder)
	if exists {
		if err := d.Dirs([]string{markersFolder}); err != nil {
			return err
		}
	}

	var files []string

	thumbPath := d.Paths.Scene.GetThumbnailScreenshotPath(sceneHash)
	exists, _ = utils.FileExists(thumbPath)
	if exists {
		files = append(files, thumbPath)
	}

	normalPath := d.Paths.Scene.GetScreenshotPath(sceneHash)
	exists, _ = utils.FileExists(normalPath)
	if exists {
		files = append(files, normalPath)
	}

	streamPreviewPath := d.Paths.Scene.GetStreamPreviewPath(sceneHash)
	exists, _ = utils.FileExists(streamPreviewPath)
	if exists {
		files = append(files, streamPreviewPath)
	}

	streamPreviewImagePath := d.Paths.Scene.GetStreamPreviewImagePath(sceneHash)
	exists, _ = utils.FileExists(streamPreviewImagePath)
	if exists {
		files = append(files, streamPreviewImagePath)
	}

	transcodePath := d.Paths.Scene.GetTranscodePath(sceneHash)
	exists, _ = utils.FileExists(transcodePath)
	if exists {
		files = append(files, transcodePath)
	}

	spritePath := d.Paths.Scene.GetSpriteImageFilePath(sceneHash)
	exists, _ = utils.FileExists(spritePath)
	if exists {
		files = append(files, spritePath)
	}

	vttPath := d.Paths.Scene.GetSpriteVttFilePath(sceneHash)
	exists, _ = utils.FileExists(vttPath)
	if exists {
		files = append(files, vttPath)
	}

	heatmapPath := d.Paths.Scene.GetInteractiveHeatmapPath(sceneHash)
	exists, _ = utils.FileExists(heatmapPath)
	if exists {
		files = append(files, heatmapPath)
	}

	return d.Files(files)
}

// MarkMarkerFiles deletes generated files for a scene marker with the
// provided scene and timestamp.
func (d *FileDeleter) MarkMarkerFiles(scene *models.Scene, seconds int) error {
	videoPath := d.Paths.SceneMarkers.GetStreamPath(scene.GetHash(d.FileNamingAlgo), seconds)
	imagePath := d.Paths.SceneMarkers.GetStreamPreviewImagePath(scene.GetHash(d.FileNamingAlgo), seconds)
	screenshotPath := d.Paths.SceneMarkers.GetStreamScreenshotPath(scene.GetHash(d.FileNamingAlgo), seconds)

	var files []string

	exists, _ := utils.FileExists(videoPath)
	if exists {
		files = append(files, videoPath)
	}

	exists, _ = utils.FileExists(imagePath)
	if exists {
		files = append(files, imagePath)
	}

	exists, _ = utils.FileExists(screenshotPath)
	if exists {
		files = append(files, screenshotPath)
	}

	return d.Files(files)
}

// Destroy deletes a scene and its associated relationships from the
// database.
func Destroy(scene *models.Scene, repo models.Repository, fileDeleter *FileDeleter, deleteGenerated, deleteFile bool) error {
	qb := repo.Scene()
	mqb := repo.SceneMarker()

	markers, err := mqb.FindBySceneID(scene.ID)
	if err != nil {
		return err
	}

	for _, m := range markers {
		if err := DestroyMarker(scene, m, mqb, fileDeleter); err != nil {
			return err
		}
	}

	if deleteFile {
		if err := fileDeleter.Files([]string{scene.Path}); err != nil {
			return err
		}

		funscriptPath := utils.GetFunscriptPath(scene.Path)
		funscriptExists, _ := utils.FileExists(funscriptPath)
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

	if err := qb.Destroy(scene.ID); err != nil {
		return err
	}

	return nil
}

// DestroyMarker deletes the scene marker from the database and returns a
// function that removes the generated files, to be executed after the
// transaction is successfully committed.
func DestroyMarker(scene *models.Scene, sceneMarker *models.SceneMarker, qb models.SceneMarkerWriter, fileDeleter *FileDeleter) error {
	if err := qb.Destroy(sceneMarker.ID); err != nil {
		return err
	}

	// delete the preview for the marker
	seconds := int(sceneMarker.Seconds)
	return fileDeleter.MarkMarkerFiles(scene, seconds)
}
