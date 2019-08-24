package api

import (
	"math/rand"

	"github.com/gobuffalo/packr/v2"
)

var performerBox *packr.Box

func initialiseImages() {
	performerBox = packr.New("Performer Box", "../../static/performer")
}

func getRandomPerformerImage() ([]byte, error) {
	imageFiles := performerBox.List()
	index := rand.Intn(len(imageFiles))
	return performerBox.Find(imageFiles[index])
}
