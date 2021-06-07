package api

import (
	"math/rand"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/stashapp/stash/pkg/utils"
)

var performerBox *packr.Box
var accessiblePerformerBoxFemale *packr.Box
var performerBoxMale *packr.Box

func initialiseImages() {
	performerBox = packr.New("Performer Box", "../../static/performer")
	accessiblePerformerBoxFemale = packr.New("Accessible Performer Box", "../../static/performer_female_accessible")
	performerBoxMale = packr.New("Male Performer Box", "../../static/performer_male")
}

func getRandomPerformerImage(gender string) ([]byte, error) {
	var box *packr.Box
	switch strings.ToUpper(gender) {
	case "FEMALE":
		box = performerBox
	case "FEMALE_ACCESSIBLE":
		box = accessiblePerformerBoxFemale
	case "MALE":
		box = performerBoxMale
	default:
		box = performerBox

	}
	imageFiles := box.List()
	index := rand.Intn(len(imageFiles))
	return box.Find(imageFiles[index])
}

func getRandomPerformerImageUsingName(name, gender string) ([]byte, error) {
	var box *packr.Box
	switch strings.ToUpper(gender) {
	case "FEMALE":
		box = performerBox
	case "FEMALE_ACCESSIBLE":
		box = accessiblePerformerBoxFemale
	case "MALE":
		box = performerBoxMale
	default:
		box = performerBox

	}
	imageFiles := box.List()
	index := utils.IntFromString(name) % uint64(len(imageFiles))
	return box.Find(imageFiles[index])
}
