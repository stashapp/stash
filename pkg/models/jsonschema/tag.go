package jsonschema

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
)

type Tag struct {
	Name          string           `json:"name,omitempty"`
	SortName      string           `json:"sort_name,omitempty"`
	Description   string           `json:"description,omitempty"`
	Favorite      bool             `json:"favorite,omitempty"`
	Aliases       []string         `json:"aliases,omitempty"`
	Image         string           `json:"image,omitempty"`
	Parents       []string         `json:"parents,omitempty"`
	IgnoreAutoTag bool             `json:"ignore_auto_tag,omitempty"`
	StashIDs      []models.StashID `json:"stash_ids,omitempty"`
	CreatedAt     json.JSONTime    `json:"created_at,omitempty"`
	UpdatedAt     json.JSONTime    `json:"updated_at,omitempty"`
}

func (s Tag) Filename() string {
	return fsutil.SanitiseBasename(s.Name) + ".json"
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
