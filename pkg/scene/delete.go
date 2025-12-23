package scene

import (
	"context"
	"path/filepath"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/file/video"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
)

// FileDeleter is an extension of file.Deleter that handles deletion of scene files.
type FileDeleter struct {
	*file.Deleter

	FileNamingAlgo models.HashAlgorithm
	Paths          *paths.Paths
}

// MarkGeneratedFiles marks for deletion the generated files for the provided scene.
// Generated files bypass trash and are permanently deleted since they can be regenerated.
func (d *FileDeleter) MarkGeneratedFiles(scene *models.Scene) error {
	sceneHash := scene.GetHash(d.FileNamingAlgo)

	if sceneHash == "" {
		return nil
	}

	markersFolder := filepath.Join(d.Paths.Generated.Markers, sceneHash)

	exists, _ := fsutil.FileExists(markersFolder)
	if exists {
		if err := d.DirsWithoutTrash([]string{markersFolder}); err != nil {
			return err
		}
	}

	var files []string

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

	return d.FilesWithoutTrash(files)
}

// MarkMarkerFiles deletes generated files for a scene marker with the
// provided scene and timestamp.
// Generated files bypass trash and are permanently deleted since they can be regenerated.
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

	return d.FilesWithoutTrash(files)
}

// Destroy deletes a scene and its associated relationships from the
// database.
func (s *Service) Destroy(ctx context.Context, scene *models.Scene, fileDeleter *FileDeleter, deleteGenerated, deleteFile bool) error {
	mqb := s.MarkerRepository
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
		if err := s.deleteFiles(ctx, scene, fileDeleter); err != nil {
			return err
		}
	}

	if deleteGenerated {
		if err := fileDeleter.MarkGeneratedFiles(scene); err != nil {
			return err
		}
	}

	if err := s.Repository.Destroy(ctx, scene.ID); err != nil {
		return err
	}

	return nil
}

// deleteFiles deletes files from the database and file system
func (s *Service) deleteFiles(ctx context.Context, scene *models.Scene, fileDeleter *FileDeleter) error {
	if err := scene.LoadFiles(ctx, s.Repository); err != nil {
		return err
	}

	for _, f := range scene.Files.List() {
		// only delete files where there is no other associated scene
		otherScenes, err := s.Repository.FindByFileID(ctx, f.ID)
		if err != nil {
			return err
		}

		if len(otherScenes) > 1 {
			// other scenes associated, don't remove
			continue
		}

		const deleteFile = true
		logger.Info("Deleting scene file: ", f.Path)
		if err := file.Destroy(ctx, s.File, f, fileDeleter.Deleter, deleteFile); err != nil {
			return err
		}

		// don't delete files in zip archives
		if f.ZipFileID == nil {
			funscriptPath := video.GetFunscriptPath(f.Path)
			funscriptExists, _ := fsutil.FileExists(funscriptPath)
			if funscriptExists {
				if err := fileDeleter.Files([]string{funscriptPath}); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// DestroyMarker deletes the scene marker from the database and returns a
// function that removes the generated files, to be executed after the
// transaction is successfully committed.
func DestroyMarker(ctx context.Context, scene *models.Scene, sceneMarker *models.SceneMarker, qb models.SceneMarkerDestroyer, fileDeleter *FileDeleter) error {
	if err := qb.Destroy(ctx, sceneMarker.ID); err != nil {
		return err
	}

	// delete the preview for the marker
	seconds := int(sceneMarker.Seconds)
	return fileDeleter.MarkMarkerFiles(scene, seconds)
}
