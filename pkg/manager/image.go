package manager

import (
	"archive/zip"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
)

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
