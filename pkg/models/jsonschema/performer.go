package jsonschema

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
)

type Performer struct {
	Name          string            `json:"name,omitempty"`
	Gender        string            `json:"gender,omitempty"`
	URL           string            `json:"url,omitempty"`
	Twitter       string            `json:"twitter,omitempty"`
	Instagram     string            `json:"instagram,omitempty"`
	Birthdate     string            `json:"birthdate,omitempty"`
	Ethnicity     string            `json:"ethnicity,omitempty"`
	Country       string            `json:"country,omitempty"`
	EyeColor      string            `json:"eye_color,omitempty"`
	Height        string            `json:"height,omitempty"`
	Measurements  string            `json:"measurements,omitempty"`
	FakeTits      string            `json:"fake_tits,omitempty"`
	CareerLength  string            `json:"career_length,omitempty"`
	Tattoos       string            `json:"tattoos,omitempty"`
	Piercings     string            `json:"piercings,omitempty"`
	Aliases       string            `json:"aliases,omitempty"`
	Favorite      bool              `json:"favorite,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
	Image         string            `json:"image,omitempty"`
	CreatedAt     json.JSONTime     `json:"created_at,omitempty"`
	UpdatedAt     json.JSONTime     `json:"updated_at,omitempty"`
	Rating        int               `json:"rating,omitempty"`
	Details       string            `json:"details,omitempty"`
	DeathDate     string            `json:"death_date,omitempty"`
	HairColor     string            `json:"hair_color,omitempty"`
	Weight        int               `json:"weight,omitempty"`
	StashIDs      []*models.StashID `json:"stash_ids,omitempty"`
	IgnoreAutoTag bool              `json:"ignore_auto_tag,omitempty"`
}

func LoadPerformerFile(filePath string) (*Performer, error) {
	var performer Performer
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&performer)
	if err != nil {
		return nil, err
	}
	return &performer, nil
}

func SavePerformerFile(filePath string, performer *Performer) error {
	if performer == nil {
		return fmt.Errorf("performer must not be nil")
	}
	return marshalToFile(filePath, performer)
}
