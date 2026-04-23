// TODO(audio): update this file

package audio

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

// FileDeleter is an extension of file.Deleter that handles deletion of audio files.
type FileDeleter struct {
	*file.Deleter

	FileNamingAlgo models.HashAlgorithm
	Paths          *paths.Paths
}

// MarkGeneratedFiles marks for deletion the generated files for the provided audio.
// Generated files bypass trash and are permanently deleted since they can be regenerated.
func (d *FileDeleter) MarkGeneratedFiles(audio *models.Audio) error {
	audioHash := audio.GetHash(d.FileNamingAlgo)

	if audioHash == "" {
		return nil
	}

	markersFolder := filepath.Join(d.Paths.Generated.Markers, audioHash)

	exists, _ := fsutil.FileExists(markersFolder)
	if exists {
		if err := d.DirsWithoutTrash([]string{markersFolder}); err != nil {
			return err
		}
	}

	var files []string

	streamPreviewPath := d.Paths.Audio.GetVideoPreviewPath(audioHash)
	exists, _ = fsutil.FileExists(streamPreviewPath)
	if exists {
		files = append(files, streamPreviewPath)
	}

	streamPreviewImagePath := d.Paths.Audio.GetWebpPreviewPath(audioHash)
	exists, _ = fsutil.FileExists(streamPreviewImagePath)
	if exists {
		files = append(files, streamPreviewImagePath)
	}

	transcodePath := d.Paths.Audio.GetTranscodePath(audioHash)
	exists, _ = fsutil.FileExists(transcodePath)
	if exists {
		files = append(files, transcodePath)
	}

	spritePath := d.Paths.Audio.GetSpriteImageFilePath(audioHash)
	exists, _ = fsutil.FileExists(spritePath)
	if exists {
		files = append(files, spritePath)
	}

	vttPath := d.Paths.Audio.GetSpriteVttFilePath(audioHash)
	exists, _ = fsutil.FileExists(vttPath)
	if exists {
		files = append(files, vttPath)
	}

	return d.FilesWithoutTrash(files)
}

// Destroy deletes a audio and its associated relationships from the
// database.
func (s *Service) Destroy(ctx context.Context, audio *models.Audio, fileDeleter *FileDeleter, deleteGenerated, deleteFile, destroyFileEntry bool) error {
	mqb := s.MarkerRepository
	markers, err := mqb.FindByAudioID(ctx, audio.ID)
	if err != nil {
		return err
	}

	for _, m := range markers {
		if err := DestroyMarker(ctx, audio, m, mqb, fileDeleter); err != nil {
			return err
		}
	}

	if deleteFile {
		if err := s.deleteFiles(ctx, audio, fileDeleter); err != nil {
			return err
		}
	} else if destroyFileEntry {
		if err := s.destroyFileEntries(ctx, audio); err != nil {
			return err
		}
	}

	if deleteGenerated {
		if err := fileDeleter.MarkGeneratedFiles(audio); err != nil {
			return err
		}
	}

	if err := s.Repository.Destroy(ctx, audio.ID); err != nil {
		return err
	}

	return nil
}

// deleteFiles deletes files from the database and file system
func (s *Service) deleteFiles(ctx context.Context, audio *models.Audio, fileDeleter *FileDeleter) error {
	if err := audio.LoadFiles(ctx, s.Repository); err != nil {
		return err
	}

	for _, f := range audio.Files.List() {
		// only delete files where there is no other associated audio
		otherAudios, err := s.Repository.FindByFileID(ctx, f.ID)
		if err != nil {
			return err
		}

		if len(otherAudios) > 1 {
			// other audios associated, don't remove
			continue
		}

		const deleteFile = true
		logger.Info("Deleting audio file: ", f.Path)
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

// destroyFileEntries destroys file entries from the database without deleting
// the files from the filesystem
func (s *Service) destroyFileEntries(ctx context.Context, audio *models.Audio) error {
	if err := audio.LoadFiles(ctx, s.Repository); err != nil {
		return err
	}

	for _, f := range audio.Files.List() {
		// only destroy file entries where there is no other associated audio
		otherAudios, err := s.Repository.FindByFileID(ctx, f.ID)
		if err != nil {
			return err
		}

		if len(otherAudios) > 1 {
			// other audios associated, don't remove
			continue
		}

		const deleteFile = false
		logger.Info("Destroying audio file entry: ", f.Path)
		if err := file.Destroy(ctx, s.File, f, nil, deleteFile); err != nil {
			return err
		}
	}

	return nil
}

// DestroyMarker deletes the audio marker from the database and returns a
// function that removes the generated files, to be executed after the
// transaction is successfully committed.
func DestroyMarker(ctx context.Context, audio *models.Audio, audioMarker *models.AudioMarker, qb models.AudioMarkerDestroyer, fileDeleter *FileDeleter) error {
	if err := qb.Destroy(ctx, audioMarker.ID); err != nil {
		return err
	}

	// delete the preview for the marker
	seconds := int(audioMarker.Seconds)
	return fileDeleter.MarkMarkerFiles(audio, seconds)
}
