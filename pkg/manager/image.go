package manager

import (
	"archive/zip"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/manager/config"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
)

func walkGalleryZip(path string, walkFunc func(file *zip.File) error) error {
	readCloser, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer readCloser.Close()

	excludeImgRegex := generateRegexps(config.GetInstance().GetImageExcludes())

	for _, f := range readCloser.File {
		if f.FileInfo().IsDir() {
			continue
		}

		if strings.Contains(f.Name, "__MACOSX") {
			continue
		}

		if !isImage(f.Name) {
			continue
		}

		if matchFileRegex(file.ZipFile(path, f).Path(), excludeImgRegex) {
			continue
		}

		err := walkFunc(f)
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
