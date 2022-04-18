package jsonschema

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/models/json"
)

type Tag struct {
	Name          string        `json:"name,omitempty"`
	Aliases       []string      `json:"aliases,omitempty"`
	Image         string        `json:"image,omitempty"`
	Parents       []string      `json:"parents,omitempty"`
	IgnoreAutoTag bool          `json:"ignore_auto_tag,omitempty"`
	CreatedAt     json.JSONTime `json:"created_at,omitempty"`
	UpdatedAt     json.JSONTime `json:"updated_at,omitempty"`
}

func LoadTagFile(filePath string) (*Tag, error) {
	var tag Tag
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
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
