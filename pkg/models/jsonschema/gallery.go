package jsonschema

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/models/json"
)

type Gallery struct {
	Path        string        `json:"path,omitempty"`
	Checksum    string        `json:"checksum,omitempty"`
	Zip         bool          `json:"zip,omitempty"`
	Title       string        `json:"title,omitempty"`
	URL         string        `json:"url,omitempty"`
	Date        string        `json:"date,omitempty"`
	Details     string        `json:"details,omitempty"`
	Rating      int           `json:"rating,omitempty"`
	Organized   bool          `json:"organized,omitempty"`
	Studio      string        `json:"studio,omitempty"`
	Performers  []string      `json:"performers,omitempty"`
	Tags        []string      `json:"tags,omitempty"`
	FileModTime json.JSONTime `json:"file_mod_time,omitempty"`
	CreatedAt   json.JSONTime `json:"created_at,omitempty"`
	UpdatedAt   json.JSONTime `json:"updated_at,omitempty"`
}

func LoadGalleryFile(filePath string) (*Gallery, error) {
	var gallery Gallery
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&gallery)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

func SaveGalleryFile(filePath string, gallery *Gallery) error {
	if gallery == nil {
		return fmt.Errorf("gallery must not be nil")
	}
	return marshalToFile(filePath, gallery)
}
