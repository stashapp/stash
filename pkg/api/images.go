package api

import (
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/utils"
)

type imageBox struct {
	box   *packr.Box
	files []string
}

func newImageBox(box *packr.Box) *imageBox {
	return &imageBox{
		box:   box,
		files: box.List(),
	}
}

var performerBox *imageBox
var performerBoxMale *imageBox
var performerBoxCustom *imageBox

func initialiseImages() {
	performerBox = newImageBox(packr.New("Performer Box", "../../static/performer"))
	performerBoxMale = newImageBox(packr.New("Male Performer Box", "../../static/performer_male"))
	initialiseCustomImages()
}

func initialiseCustomImages() {
	customPath := config.GetInstance().GetCustomPerformerImageLocation()
	if customPath != "" {
		logger.Debugf("Loading custom performer images from %s", customPath)
		// We need to set performerBoxCustom at runtime, as this is a custom path, and store it in a pointer.
		performerBoxCustom = newImageBox(packr.Folder(customPath))
	} else {
		performerBoxCustom = nil
	}
}

func getRandomPerformerImageUsingName(name, gender, customPath string) ([]byte, error) {
	var box *imageBox

	// If we have a custom path, we should return a new box in the given path.
	if performerBoxCustom != nil && len(performerBoxCustom.files) > 0 {
		box = performerBoxCustom
	}

	if box == nil {
		switch strings.ToUpper(gender) {
		case "FEMALE":
			box = performerBox
		case "MALE":
			box = performerBoxMale
		default:
			box = performerBox
		}
	}

	imageFiles := box.files
	index := utils.IntFromString(name) % uint64(len(imageFiles))
	return box.box.Find(imageFiles[index])
}
