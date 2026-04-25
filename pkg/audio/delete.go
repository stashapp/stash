// TODO(audio): update this file

package audio

import (
	"context"
	"path/filepath"

	"github.com/stashapp/stash/pkg/file"
	file_audio "github.com/stashapp/stash/pkg/file/audio"
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

	// TODO(future|audio generated files): add paths here

	return d.FilesWithoutTrash(files)
}

// Destroy deletes a audio and its associated relationships from the
// database.
func (s *Service) Destroy(ctx context.Context, audio *models.Audio, fileDeleter *FileDeleter, deleteGenerated, deleteFile, destroyFileEntry bool) error {
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
			funscriptPath := file_audio.GetFunscriptPath(f.Path)
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
