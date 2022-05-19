package jsonschema

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
)

type SceneMarker struct {
	Title      string        `json:"title,omitempty"`
	Seconds    string        `json:"seconds,omitempty"`
	PrimaryTag string        `json:"primary_tag,omitempty"`
	Tags       []string      `json:"tags,omitempty"`
	CreatedAt  json.JSONTime `json:"created_at,omitempty"`
	UpdatedAt  json.JSONTime `json:"updated_at,omitempty"`
}

type SceneFile struct {
	ModTime    json.JSONTime `json:"mod_time,omitempty"`
	Size       string        `json:"size"`
	Duration   string        `json:"duration"`
	VideoCodec string        `json:"video_codec"`
	AudioCodec string        `json:"audio_codec"`
	Format     string        `json:"format"`
	Width      int           `json:"width"`
	Height     int           `json:"height"`
	Framerate  string        `json:"framerate"`
	Bitrate    int           `json:"bitrate"`
}

type SceneMovie struct {
	MovieName  string `json:"movieName,omitempty"`
	SceneIndex int    `json:"scene_index,omitempty"`
}

type Scene struct {
	Title      string           `json:"title,omitempty"`
	Checksum   string           `json:"checksum,omitempty"`
	OSHash     string           `json:"oshash,omitempty"`
	Phash      string           `json:"phash,omitempty"`
	Studio     string           `json:"studio,omitempty"`
	URL        string           `json:"url,omitempty"`
	Date       string           `json:"date,omitempty"`
	Rating     int              `json:"rating,omitempty"`
	Organized  bool             `json:"organized,omitempty"`
	OCounter   int              `json:"o_counter,omitempty"`
	Details    string           `json:"details,omitempty"`
	Galleries  []string         `json:"galleries,omitempty"`
	Performers []string         `json:"performers,omitempty"`
	Movies     []SceneMovie     `json:"movies,omitempty"`
	Tags       []string         `json:"tags,omitempty"`
	Markers    []SceneMarker    `json:"markers,omitempty"`
	File       *SceneFile       `json:"file,omitempty"`
	Cover      string           `json:"cover,omitempty"`
	CreatedAt  json.JSONTime    `json:"created_at,omitempty"`
	UpdatedAt  json.JSONTime    `json:"updated_at,omitempty"`
	StashIDs   []models.StashID `json:"stash_ids,omitempty"`
}

func LoadSceneFile(filePath string) (*Scene, error) {
	var scene Scene
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&scene)
	if err != nil {
		return nil, err
	}
	return &scene, nil
}

func SaveSceneFile(filePath string, scene *Scene) error {
	if scene == nil {
		return fmt.Errorf("scene must not be nil")
	}
	return marshalToFile(filePath, scene)
}
