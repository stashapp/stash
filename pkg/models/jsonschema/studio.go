package jsonschema

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
)

type Studio struct {
	Name          string            `json:"name,omitempty"`
	URL           string            `json:"url,omitempty"`
	ParentStudio  string            `json:"parent_studio,omitempty"`
	Image         string            `json:"image,omitempty"`
	CreatedAt     json.JSONTime     `json:"created_at,omitempty"`
	UpdatedAt     json.JSONTime     `json:"updated_at,omitempty"`
	Rating        int               `json:"rating,omitempty"`
	Details       string            `json:"details,omitempty"`
	Aliases       []string          `json:"aliases,omitempty"`
	StashIDs      []*models.StashID `json:"stash_ids,omitempty"`
	IgnoreAutoTag bool              `json:"ignore_auto_tag,omitempty"`
}

func LoadStudioFile(filePath string) (*Studio, error) {
	var studio Studio
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&studio)
	if err != nil {
		return nil, err
	}
	return &studio, nil
}

func SaveStudioFile(filePath string, studio *Studio) error {
	if studio == nil {
		return fmt.Errorf("studio must not be nil")
	}
	return marshalToFile(filePath, studio)
}
