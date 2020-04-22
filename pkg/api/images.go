package api

import (
	"math/rand"

	"github.com/gobuffalo/packr/v2"
)

var performerBox *packr.Box
var performerBoxMale *packr.Box

func initialiseImages() {
	performerBox = packr.New("Performer Box", "../../static/performer")
	performerBoxMale = packr.New("Male Performer Box", "../../static/performer_male")
}

func getRandomPerformerImage(gender string) ([]byte, error) {
	var box *packr.Box
	switch gender {
	case "FEMALE":
		box = performerBox
	case "MALE":
		box = performerBoxMale
	default:
		box = performerBox

	}
	imageFiles := box.List()
	index := rand.Intn(len(imageFiles))
	return box.Find(imageFiles[index])
}
