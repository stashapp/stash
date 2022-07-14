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
	const fromZip = false
	return s.destroyImage(ctx, i, fileDeleter, deleteGenerated, deleteFile, fromZip)
}

// DestroyZipImages destroys all images in zip, optionally marking the files and generated files for deletion.
// Returns a slice of images that were destroyed.
func (s *Service) DestroyZipImages(ctx context.Context, zipFile file.File, fileDeleter *FileDeleter, deleteGenerated bool) ([]*models.Image, error) {
	// TODO - we currently destroy associated files so that they will be rescanned.
	// A better way would be to keep the file entries in the database, and recreate
	// associated objects during the scan process if there are none already.

	var imgsDestroyed []*models.Image

	imgs, err := s.Repository.FindByZipFileID(ctx, zipFile.Base().ID)
	if err != nil {
		return nil, err
	}

	for _, img := range imgs {
		const deleteFileInZip = false
		const fromZip = true
		if err := s.destroyImage(ctx, img, fileDeleter, deleteGenerated, deleteFileInZip, fromZip); err != nil {
			return nil, err
		}

		imgsDestroyed = append(imgsDestroyed, img)
	}

	return imgsDestroyed, nil
}

// Destroy destroys an image, optionally marking the file and generated files for deletion.
func (s *Service) destroyImage(ctx context.Context, i *models.Image, fileDeleter *FileDeleter, deleteGenerated, deleteFile bool, fromZip bool) error {
	// TODO - we currently destroy associated files so that they will be rescanned.
	// A better way would be to keep the file entries in the database, and recreate
	// associated objects during the scan process if there are none already.

	if err := s.destroyFiles(ctx, i, fileDeleter, deleteFile, fromZip); err != nil {
		return err
	}

	if deleteGenerated {
		if err := fileDeleter.MarkGeneratedFiles(i); err != nil {
			return err
		}
	}

	return s.Repository.Destroy(ctx, i.ID)
}

func (s *Service) destroyFiles(ctx context.Context, i *models.Image, fileDeleter *FileDeleter, deleteFile bool, fromZip bool) error {
	for _, f := range i.Files {
		// only delete files where there is no other associated image
		otherImages, err := s.Repository.FindByFileID(ctx, f.ID)
		if err != nil {
			return err
		}

		if len(otherImages) > 1 {
			// other image associated, don't remove
			continue
		}

		// don't delete files in zip archives
		if fromZip || f.ZipFileID == nil {
			if err := file.Destroy(ctx, s.File, f, fileDeleter.Deleter, deleteFile); err != nil {
				return err
			}
		}
	}

	return nil
}
