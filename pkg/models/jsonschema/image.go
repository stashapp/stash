package jsonschema

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/models/json"
)

type ImageFile struct {
	ModTime json.JSONTime `json:"mod_time,omitempty"`
	Size    int64         `json:"size"`
	Width   int           `json:"width"`
	Height  int           `json:"height"`
}

type Image struct {
	Title      string        `json:"title,omitempty"`
	Checksum   string        `json:"checksum,omitempty"`
	Studio     string        `json:"studio,omitempty"`
	Rating     int           `json:"rating,omitempty"`
	Organized  bool          `json:"organized,omitempty"`
	OCounter   int           `json:"o_counter,omitempty"`
	Galleries  []string      `json:"galleries,omitempty"`
	Performers []string      `json:"performers,omitempty"`
	Tags       []string      `json:"tags,omitempty"`
	File       *ImageFile    `json:"file,omitempty"`
	CreatedAt  json.JSONTime `json:"created_at,omitempty"`
	UpdatedAt  json.JSONTime `json:"updated_at,omitempty"`
}

func LoadImageFile(filePath string) (*Image, error) {
	var image Image
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&image)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func SaveImageFile(filePath string, image *Image) error {
	if image == nil {
		return fmt.Errorf("image must not be nil")
	}
	return marshalToFile(filePath, image)
}
