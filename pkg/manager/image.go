package manager

import (
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

// DestroyImage deletes an image and its associated relationships from the
// database.
func DestroyImage(imageID int, tx *sqlx.Tx) error {
	qb := models.NewImageQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()

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

	if err := qb.Destroy(strconv.Itoa(imageID), tx); err != nil {
		return err
	}

	return nil
}

// DeleteGeneratedImageFiles deletes generated files for the provided image.
func DeleteGeneratedImageFiles(image *models.Image) {
	// thumbPath := GetInstance().Paths.Scene.GetThumbnailScreenshotPath(sceneHash)
	// exists, _ = utils.FileExists(thumbPath)
	// if exists {
	// 	err := os.Remove(thumbPath)
	// 	if err != nil {
	// 		logger.Warnf("Could not delete file %s: %s", thumbPath, err.Error())
	// 	}
	// }
}

// DeleteImageFile deletes the image file from the filesystem.
func DeleteImageFile(image *models.Image) {
	err := os.Remove(image.Path)
	if err != nil {
		logger.Warnf("Could not delete file %s: %s", image.Path, err.Error())
	}
}
