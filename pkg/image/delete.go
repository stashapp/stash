package image

import (
	"context"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
)

type Destroyer interface {
	Destroy(ctx context.Context, id int) error
}

// FileDeleter is an extension of file.Deleter that handles deletion of image files.
type FileDeleter struct {
	*file.Deleter

	Paths *paths.Paths
}

// MarkGeneratedFiles marks for deletion the generated files for the provided image.
func (d *FileDeleter) MarkGeneratedFiles(image *models.Image) error {
	thumbPath := d.Paths.Generated.GetThumbnailPath(image.Checksum(), models.DefaultGthumbWidth)
	exists, _ := fsutil.FileExists(thumbPath)
	if exists {
		return d.Files([]string{thumbPath})
	}

	return nil
}

// Destroy destroys an image, optionally marking the file and generated files for deletion.
func (s *Service) Destroy(ctx context.Context, i *models.Image, fileDeleter *FileDeleter, deleteGenerated, deleteFile bool) error {
	if deleteFile {
		if err := s.deleteFiles(ctx, i, fileDeleter); err != nil {
			return err
		}
	}

	if deleteGenerated {
		if err := fileDeleter.MarkGeneratedFiles(i); err != nil {
			return err
		}
	}

	return s.Repository.Destroy(ctx, i.ID)
}

func (s *Service) deleteFiles(ctx context.Context, i *models.Image, fileDeleter *FileDeleter) error {
	for _, f := range i.Files {
		// don't delete files in zip archives
		if f.ZipFileID != nil {
			continue
		}

		// only delete files where there is no other associated image
		otherImages, err := s.Repository.FindByFileID(ctx, f.ID)
		if err != nil {
			return err
		}

		if len(otherImages) > 1 {
			// other image associated, don't remove
			continue
		}

		const deleteFile = true
		if err := file.Destroy(ctx, s.File, f, fileDeleter.Deleter, deleteFile); err != nil {
			return err
		}
	}

	return nil
}
