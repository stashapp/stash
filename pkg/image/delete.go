package image

import (
	"fmt"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type Destroyer interface {
	GetFileIDs(id int) ([]int, error)
	Destroy(id int) error
}

// FileDeleter is an extension of file.Deleter that handles deletion of image files.
type FileDeleter struct {
	file.Deleter

	Paths *paths.Paths
}

// MarkGeneratedFiles marks for deletion the generated files for the provided image.
func (d *FileDeleter) MarkGeneratedFiles(image *models.Image) error {
	thumbPath := d.Paths.Generated.GetThumbnailPath(image.Checksum, models.DefaultGthumbWidth)
	exists, _ := utils.FileExists(thumbPath)
	if exists {
		return d.Files([]string{thumbPath})
	}

	return nil
}

// Destroy destroys an image, optionally marking the file and generated files for deletion.
func Destroy(i *models.Image, repo models.Repository, fileDeleter *FileDeleter, deleteGenerated, deleteFile bool) error {
	iqb := repo.Image()
	fqb := repo.File()

	// destroy associated files - this assumes that a file can only be
	// associated with a single scene
	fileIDs, err := iqb.GetFileIDs(i.ID)
	if err != nil {
		return fmt.Errorf("getting related file ids for id %d: %w", i.ID, err)
	}

	files, err := fqb.Find(fileIDs)
	if err != nil {
		return fmt.Errorf("getting image files: %w", err)
	}

	for _, f := range files {
		if err := fqb.Destroy(f.ID); err != nil {
			return err
		}

		// don't try to delete if the image is in a zip file
		if deleteFile && !f.ZipFileID.Valid {
			if err := fileDeleter.Files([]string{f.Path}); err != nil {
				return err
			}
		}
	}

	if deleteGenerated {
		if err := fileDeleter.MarkGeneratedFiles(i); err != nil {
			return err
		}
	}

	return iqb.Destroy(i.ID)
}
