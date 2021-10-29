package manager

import (
	"archive/zip"
	"os"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

const stashDeletePostfix = ".stashdelete"

// DeleteGeneratedImageFiles deletes generated files for the provided image.
func DeleteGeneratedImageFiles(image *models.Image) error {
	thumbPath := GetInstance().Paths.Generated.GetThumbnailPath(image.Checksum, models.DefaultGthumbWidth)
	exists, _ := utils.FileExists(thumbPath)
	if exists {
		err := os.Remove(thumbPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", thumbPath, err.Error())
		}
		return err
	}
	return nil
}

// DeleteImageFile deletes the image file from the filesystem.
func DeleteImageFile(image *models.Image) error {
	err := os.Remove(image.Path)
	if err != nil {
		logger.Warnf("Could not delete file %s: %s", image.Path, err.Error())
	}
	return err
}

// MarkImageFileForDeleteion checks if a file is writable by renaming it.
// Use DeleteMarkedImageFile afterwards to actually delete.
func MarkImageFileForDeletion(image *models.Image) error {
	// handle image deleted outside of Stash
	_, err := os.Stat(image.Path)
	if os.IsNotExist(err) {
		logger.Debugf("Could not mark for deletion: %s: %s", image.Path, err.Error())
		return nil
	}
	err = os.Rename(image.Path, image.Path+stashDeletePostfix)
	if err != nil {
		logger.Errorf("Could not mark for deletion: %s: %s", image.Path, err.Error())
	}
	return err
}

func UnmarkImageFileForDeletion(image *models.Image) {
	err := os.Rename(image.Path+stashDeletePostfix, image.Path)
	if err != nil {
		logger.Debug("Could not unmark file %s: %s", image.Path, err.Error())
	}
}

func DeleteMarkedImageFile(image *models.Image) error {
	// handle missing renamed image, probably file was already deleted
	_, err := os.Stat(image.Path + stashDeletePostfix)
	if os.IsNotExist(err) {
		err = DeleteImageFile(image)
		if err != nil {
			logger.Debugf("Could not delete file %s: %s", image.Path, err.Error())
		}
		return nil
	}
	err = os.Remove(image.Path + stashDeletePostfix)
	if err != nil {
		logger.Warnf("Could not delete file %s: %s", image.Path, err.Error())
	}
	return err
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
	err := walkGalleryZip(path, func(file *zip.File) error {
		ret++
		return nil
	})
	if err != nil {
		logger.Warnf("Error while walking gallery zip: %v", err)
	}

	return ret
}
