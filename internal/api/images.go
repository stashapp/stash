package api

import (
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/internal/static"
	"github.com/stashapp/stash/pkg/hash"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type imageBox struct {
	box   fs.FS
	files []string
}

var imageExtensions = []string{
	".jpg",
	".jpeg",
	".png",
	".gif",
	".svg",
	".webp",
}

func newImageBox(box fs.FS) (*imageBox, error) {
	ret := &imageBox{
		box: box,
	}

	err := fs.WalkDir(box, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		baseName := strings.ToLower(d.Name())
		for _, ext := range imageExtensions {
			if strings.HasSuffix(baseName, ext) {
				ret.files = append(ret.files, path)
				break
			}
		}

		return nil
	})

	return ret, err
}

var performerBox *imageBox
var performerBoxMale *imageBox
var performerBoxCustom *imageBox

func initialiseImages() {
	var err error
	performerBox, err = newImageBox(&static.Performer)
	if err != nil {
		logger.Warnf("error loading performer images: %v", err)
	}
	performerBoxMale, err = newImageBox(&static.PerformerMale)
	if err != nil {
		logger.Warnf("error loading male performer images: %v", err)
	}
	initialiseCustomImages()
}

func initialiseCustomImages() {
	customPath := config.GetInstance().GetCustomPerformerImageLocation()
	if customPath != "" {
		logger.Debugf("Loading custom performer images from %s", customPath)
		// We need to set performerBoxCustom at runtime, as this is a custom path, and store it in a pointer.
		var err error
		performerBoxCustom, err = newImageBox(os.DirFS(customPath))
		if err != nil {
			logger.Warnf("error loading custom performer from %s: %v", customPath, err)
		}
	} else {
		performerBoxCustom = nil
	}
}

func getRandomPerformerImageUsingName(name string, gender *models.GenderEnum, customPath string) ([]byte, error) {
	var box *imageBox

	// If we have a custom path, we should return a new box in the given path.
	if performerBoxCustom != nil && len(performerBoxCustom.files) > 0 {
		box = performerBoxCustom
	}

	var g models.GenderEnum
	if gender != nil {
		g = *gender
	}

	if box == nil {
		switch g {
		case models.GenderEnumFemale, models.GenderEnumTransgenderFemale:
			box = performerBox
		case models.GenderEnumMale, models.GenderEnumTransgenderMale:
			box = performerBoxMale
		default:
			box = performerBox
		}
	}

	imageFiles := box.files
	index := hash.IntFromString(name) % uint64(len(imageFiles))
	img, err := box.box.Open(imageFiles[index])
	if err != nil {
		return nil, err
	}
	defer img.Close()

	return io.ReadAll(img)
}
