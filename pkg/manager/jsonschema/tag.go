package jsonschema

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/models"
)

type Tag struct {
	Name      string          `json:"name,omitempty"`
	Image     string          `json:"image,omitempty"`
	CreatedAt models.JSONTime `json:"created_at,omitempty"`
	UpdatedAt models.JSONTime `json:"updated_at,omitempty"`
}

func LoadTagFile(filePath string) (*Tag, error) {
	var tag Tag
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&tag)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func SaveTagFile(filePath string, tag *Tag) error {
	if tag == nil {
		return fmt.Errorf("tag must not be nil")
	}
	return marshalToFile(filePath, tag)
}
