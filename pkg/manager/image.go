package manager

import (
	"archive/zip"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/utils"
)

// DestroyImage deletes an image and its associated relationships from the
// database.
func DestroyImage(imageID int, tx *sqlx.Tx) error {
	qb := sqlite.NewImageQueryBuilder()
	jqb := sqlite.NewJoinsQueryBuilder()

	_, err := qb.Find(imageID)
	if err != nil {
		return err
	}

	if err := jqb.DestroyImagesTags(imageID, tx); err != nil {
		return err
	}

	if err := jqb.DestroyPerformersImages(imageID, tx); err != nil {
		return err
	}

	if err := jqb.DestroyImageGalleries(imageID, tx); err != nil {
		return err
	}

	if err := qb.Destroy(imageID, tx); err != nil {
		return err
	}

	return nil
}

// DeleteGeneratedImageFiles deletes generated files for the provided image.
func DeleteGeneratedImageFiles(image *models.Image) {
	thumbPath := GetInstance().Paths.Generated.GetThumbnailPath(image.Checksum, models.DefaultGthumbWidth)
	exists, _ := utils.FileExists(thumbPath)
	if exists {
		err := os.Remove(thumbPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", thumbPath, err.Error())
		}
	}
}

// DeleteImageFile deletes the image file from the filesystem.
func DeleteImageFile(image *models.Image) {
	err := os.Remove(image.Path)
	if err != nil {
		logger.Warnf("Could not delete file %s: %s", image.Path, err.Error())
	}
}

func walkGalleryZip(path string, walkFunc func(file *zip.File) error) error {
	readCloser, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer readCloser.Close()

	for _, file := range readCloser.File {
		if file.FileInfo().IsDir() {
			continue
		}

		if strings.Contains(file.Name, "__MACOSX") {
			continue
		}

		if !isImage(file.Name) {
			continue
		}

		err := walkFunc(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func countImagesInZip(path string) int {
	ret := 0
	walkGalleryZip(path, func(file *zip.File) error {
		ret++
		return nil
	})

	return ret
}
