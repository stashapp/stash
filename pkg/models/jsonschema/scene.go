package jsonschema

import (
	"fmt"
	"os"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
)

type SceneFilter struct {
	Contrast    int           `json:"contrast,omitempty"`
	Brightness  int           `json:"brightness,omitempty"`
	Gamma       int           `json:"gamma,omitempty"`
	Saturate    int           `json:"saturate,omitempty"`
	HueRotate   int           `json:"hue_rotate,omitempty"`
	Warmth      int           `json:"warmth,omitempty"`
	Red         int           `json:"red,omitempty"`
	Green       int           `json:"green,omitempty"`
	Blue        int           `json:"blue,omitempty"`
	Blur        int           `json:"blur,omitempty"`
	Rotate      float64       `json:"rotate,omitempty"`
	Scale       int           `json:"scale,omitempty"`
	AspectRatio int           `json:"aspect_ratio,omitempty"`
	CreatedAt   json.JSONTime `json:"created_at,omitempty"`
	UpdatedAt   json.JSONTime `json:"updated_at,omitempty"`
}

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
	Title  string `json:"title,omitempty"`
	Code   string `json:"code,omitempty"`
	Studio string `json:"studio,omitempty"`
	// deprecated - for import only
	URL          string           `json:"url,omitempty"`
	URLs         []string         `json:"urls,omitempty"`
	Date         string           `json:"date,omitempty"`
	Rating       int              `json:"rating,omitempty"`
	Organized    bool             `json:"organized,omitempty"`
	OCounter     int              `json:"o_counter,omitempty"`
	Details      string           `json:"details,omitempty"`
	Director     string           `json:"director,omitempty"`
	Galleries    []GalleryRef     `json:"galleries,omitempty"`
	Performers   []string         `json:"performers,omitempty"`
	Movies       []SceneMovie     `json:"movies,omitempty"`
	Tags         []string         `json:"tags,omitempty"`
	Filters      []SceneFilter    `json:"filters,omitempty"`
	Markers      []SceneMarker    `json:"markers,omitempty"`
	Files        []string         `json:"files,omitempty"`
	Cover        string           `json:"cover,omitempty"`
	CreatedAt    json.JSONTime    `json:"created_at,omitempty"`
	UpdatedAt    json.JSONTime    `json:"updated_at,omitempty"`
	LastPlayedAt json.JSONTime    `json:"last_played_at,omitempty"`
	ResumeTime   float64          `json:"resume_time,omitempty"`
	PlayCount    int              `json:"play_count,omitempty"`
	PlayDuration float64          `json:"play_duration,omitempty"`
	StashIDs     []models.StashID `json:"stash_ids,omitempty"`
}

func (s Scene) Filename(id int, basename string, hash string) string {
	ret := fsutil.SanitiseBasename(s.Title)
	if ret == "" {
		ret = basename
	}

	if hash != "" {
		ret += "." + hash
	} else {
		// scenes may have no file and therefore no hash
		ret += "." + strconv.Itoa(id)
	}

	return ret + ".json"
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
