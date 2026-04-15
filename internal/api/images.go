package api

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/stashapp/stash/internal/static"
	"github.com/stashapp/stash/pkg/hash"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type imageBox struct {
	box   fs.FS
	files []string
}

var imageBoxExts = []string{
	".jpg",
	".jpeg",
	".png",
	".gif",
	".svg",
	".webp",
	".avif",
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
		for _, ext := range imageBoxExts {
			if strings.HasSuffix(baseName, ext) {
				ret.files = append(ret.files, path)
				break
			}
		}

		return nil
	})

	return ret, err
}

func (box *imageBox) GetRandomImageByName(name string) ([]byte, error) {
	files := box.files
	if len(files) == 0 {
		return nil, errors.New("box is empty")
	}

	index := hash.IntFromString(name) % uint64(len(files))
	img, err := box.box.Open(files[index])
	if err != nil {
		return nil, err
	}
	defer img.Close()

	return io.ReadAll(img)
}

var performerBox *imageBox
var performerBoxMale *imageBox
var performerBoxCustom *imageBox

func init() {
	var err error
	performerBox, err = newImageBox(static.Sub(static.Performer))
	if err != nil {
		panic(fmt.Sprintf("loading performer images: %v", err))
	}
	performerBoxMale, err = newImageBox(static.Sub(static.PerformerMale))
	if err != nil {
		panic(fmt.Sprintf("loading male performer images: %v", err))
	}
}

func initCustomPerformerImages(customPath string) {
	if customPath != "" {
		logger.Debugf("Loading custom performer images from %s", customPath)
		var err error
		performerBoxCustom, err = newImageBox(os.DirFS(customPath))
		if err != nil {
			logger.Warnf("error loading custom performer images from %s: %v", customPath, err)
		}
	} else {
		performerBoxCustom = nil
	}
}

func getDefaultPerformerImage(name string, gender *models.GenderEnum, sfwMode bool) []byte {
	// try the custom box first if we have one
	if performerBoxCustom != nil {
		ret, err := performerBoxCustom.GetRandomImageByName(name)
		if err == nil {
			return ret
		}
		logger.Warnf("error loading custom default performer image: %v", err)
	}

	if sfwMode {
		return static.ReadAll(static.DefaultSFWPerformerImage)
	}

	var g models.GenderEnum
	if gender != nil {
		g = *gender
	}

	var box *imageBox
	switch g {
	case models.GenderEnumFemale, models.GenderEnumTransgenderFemale:
		box = performerBox
	case models.GenderEnumMale, models.GenderEnumTransgenderMale:
		box = performerBoxMale
	default:
		box = performerBox
	}

	ret, err := box.GetRandomImageByName(name)
	if err != nil {
		logger.Warnf("error loading default performer image: %v", err)
	}
	return ret
}
