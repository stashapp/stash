package manager

import (
	"github.com/corona10/goimagehash"
	"github.com/stashapp/stash/pkg/utils"

	"image/jpeg"
	"os"
)

type PhashGenerator struct {
	Info       *GeneratorInfo
	SpritePath string
}

func NewPhashGenerator(path string) (*PhashGenerator, error) {
	exists, err := utils.FileExists(path)
	if !exists {
		return nil, err
	}

	return &PhashGenerator{
		SpritePath: path,
	}, nil
}

func (g *PhashGenerator) Generate() (*uint64, error) {
	file, err := os.Open(g.SpritePath)
	if err != nil {
		return nil, err
	}
	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}
	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return nil, err
	}
	hashValue := hash.GetHash()
	return &hashValue, nil
}
