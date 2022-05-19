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
	file.Deleter

	Paths *paths.Paths
}

// MarkGeneratedFiles marks for deletion the generated files for the provided image.
func (d *FileDeleter) MarkGeneratedFiles(image *models.Image) error {
	thumbPath := d.Paths.Generated.GetThumbnailPath(image.Checksum, models.DefaultGthumbWidth)
	exists, _ := fsutil.FileExists(thumbPath)
	if exists {
		return d.Files([]string{thumbPath})
	}

	return nil
}

// Destroy destroys an image, optionally marking the file and generated files for deletion.
func Destroy(ctx context.Context, i *models.Image, destroyer Destroyer, fileDeleter *FileDeleter, deleteGenerated, deleteFile bool) error {
	// don't try to delete if the image is in a zip file
	if deleteFile && !file.IsZipPath(i.Path) {
		if err := fileDeleter.Files([]string{i.Path}); err != nil {
			return err
		}
	}

	if deleteGenerated {
		if err := fileDeleter.MarkGeneratedFiles(i); err != nil {
			return err
		}
	}

	return destroyer.Destroy(ctx, i.ID)
}
