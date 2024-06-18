package jsonschema

import (
	"fmt"
	"io"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

type StringOrStringList []string

func (s *StringOrStringList) UnmarshalJSON(data []byte) error {
	var stringList []string
	var stringVal string

	err := jsoniter.Unmarshal(data, &stringList)
	if err == nil {
		*s = stringList
		return nil
	}

	err = jsoniter.Unmarshal(data, &stringVal)
	if err == nil {
		*s = stringslice.FromString(stringVal, ",")
		return nil
	}

	return err
}

type Performer struct {
	Name           string   `json:"name,omitempty"`
	Disambiguation string   `json:"disambiguation,omitempty"`
	Gender         string   `json:"gender,omitempty"`
	URLs           []string `json:"urls,omitempty"`
	Birthdate      string   `json:"birthdate,omitempty"`
	Ethnicity      string   `json:"ethnicity,omitempty"`
	Country        string   `json:"country,omitempty"`
	EyeColor       string   `json:"eye_color,omitempty"`
	// this should be int, but keeping string for backwards compatibility
	Height        string             `json:"height,omitempty"`
	Measurements  string             `json:"measurements,omitempty"`
	FakeTits      string             `json:"fake_tits,omitempty"`
	PenisLength   float64            `json:"penis_length,omitempty"`
	Circumcised   string             `json:"circumcised,omitempty"`
	CareerLength  string             `json:"career_length,omitempty"`
	Tattoos       string             `json:"tattoos,omitempty"`
	Piercings     string             `json:"piercings,omitempty"`
	Aliases       StringOrStringList `json:"aliases,omitempty"`
	Favorite      bool               `json:"favorite,omitempty"`
	Tags          []string           `json:"tags,omitempty"`
	Image         string             `json:"image,omitempty"`
	CreatedAt     json.JSONTime      `json:"created_at,omitempty"`
	UpdatedAt     json.JSONTime      `json:"updated_at,omitempty"`
	Rating        int                `json:"rating,omitempty"`
	Details       string             `json:"details,omitempty"`
	DeathDate     string             `json:"death_date,omitempty"`
	HairColor     string             `json:"hair_color,omitempty"`
	Weight        int                `json:"weight,omitempty"`
	StashIDs      []models.StashID   `json:"stash_ids,omitempty"`
	IgnoreAutoTag bool               `json:"ignore_auto_tag,omitempty"`

	// deprecated - for import only
	URL       string `json:"url,omitempty"`
	Twitter   string `json:"twitter,omitempty"`
	Instagram string `json:"instagram,omitempty"`
}

func (s Performer) Filename() string {
	name := s.Name
	if s.Disambiguation != "" {
		name += "_" + s.Disambiguation
	}
	return fsutil.SanitiseBasename(name) + ".json"
}

func LoadPerformerFile(filePath string) (*Performer, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return loadPerformer(file)
}

func loadPerformer(r io.ReadSeeker) (*Performer, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonParser := json.NewDecoder(r)

	var performer Performer
	if err := jsonParser.Decode(&performer); err != nil {
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
