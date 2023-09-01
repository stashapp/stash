package image

import (
	"context"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
)

// FileDeleter is an extension of file.Deleter that handles deletion of image files.
type FileDeleter struct {
	*file.Deleter

	Paths *paths.Paths
}

// MarkGeneratedFiles marks for deletion the generated files for the provided image.
func (d *FileDeleter) MarkGeneratedFiles(image *models.Image) error {
	var files []string
	thumbPath := d.Paths.Generated.GetThumbnailPath(image.Checksum, models.DefaultGthumbWidth)
	exists, _ := fsutil.FileExists(thumbPath)
	if exists {
		files = append(files, thumbPath)
	}
	prevPath := d.Paths.Generated.GetClipPreviewPath(image.Checksum, models.DefaultGthumbWidth)
	exists, _ = fsutil.FileExists(prevPath)
	if exists {
		files = append(files, prevPath)
	}

	return d.Files(files)
}

// Destroy destroys an image, optionally marking the file and generated files for deletion.
func (s *Service) Destroy(ctx context.Context, i *models.Image, fileDeleter *FileDeleter, deleteGenerated, deleteFile bool) error {
	return s.destroyImage(ctx, i, fileDeleter, deleteGenerated, deleteFile)
}

// DestroyZipImages destroys all images in zip, optionally marking the files and generated files for deletion.
// Returns a slice of images that were destroyed.
func (s *Service) DestroyZipImages(ctx context.Context, zipFile models.File, fileDeleter *FileDeleter, deleteGenerated bool) ([]*models.Image, error) {
	var imgsDestroyed []*models.Image

	imgs, err := s.Repository.FindByZipFileID(ctx, zipFile.Base().ID)
	if err != nil {
		return nil, err
	}

	for _, img := range imgs {
		if err := img.LoadFiles(ctx, s.Repository); err != nil {
			return nil, err
		}

		const deleteFileInZip = false
		if err := s.destroyImage(ctx, img, fileDeleter, deleteGenerated, deleteFileInZip); err != nil {
			return nil, err
		}

		imgsDestroyed = append(imgsDestroyed, img)
	}

	return imgsDestroyed, nil
}

// Destroy destroys an image, optionally marking the file and generated files for deletion.
func (s *Service) destroyImage(ctx context.Context, i *models.Image, fileDeleter *FileDeleter, deleteGenerated, deleteFile bool) error {
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

// deleteFiles deletes files for the image from the database and file system, if they are not in use by other images
func (s *Service) deleteFiles(ctx context.Context, i *models.Image, fileDeleter *FileDeleter) error {
	if err := i.LoadFiles(ctx, s.Repository); err != nil {
		return err
	}

	for _, f := range i.Files.List() {
		// only delete files where there is no other associated image
		otherImages, err := s.Repository.FindByFileID(ctx, f.Base().ID)
		if err != nil {
			return err
		}

		if len(otherImages) > 1 {
			// other image associated, don't remove
			continue
		}

		// don't delete files in zip archives
		const deleteFile = true
		if f.Base().ZipFileID == nil {
			logger.Info("Deleting image file: ", f.Base().Path)
			if err := file.Destroy(ctx, s.File, f, fileDeleter.Deleter, deleteFile); err != nil {
				return err
			}
		}
	}

	return nil
}
