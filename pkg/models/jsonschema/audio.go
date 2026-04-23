package jsonschema

import (
	"fmt"
	"os"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models/json"
)

type AudioFile struct {
	ModTime    json.JSONTime `json:"mod_time,omitempty"`
	Size       string        `json:"size"`
	Duration   string        `json:"duration"`
	AudioCodec string        `json:"audio_codec"`
	Format     string        `json:"format"`
	Width      int           `json:"width"`
	Height     int           `json:"height"`
	Samplerate string        `json:"samplerate"`
	Bitrate    int           `json:"bitrate"`
}

type AudioGroup struct {
	GroupName  string `json:"movieName,omitempty"`
	AudioIndex int    `json:"audio_index,omitempty"`
}

type Audio struct {
	Title  string `json:"title,omitempty"`
	Code   string `json:"code,omitempty"`
	Studio string `json:"studio,omitempty"`

	// deprecated - for import only
	URL string `json:"url,omitempty"`

	URLs      []string `json:"urls,omitempty"`
	Date      string   `json:"date,omitempty"`
	Rating    int      `json:"rating,omitempty"`
	Organized bool     `json:"organized,omitempty"`

	// deprecated - for import only
	OCounter int `json:"o_counter,omitempty"`

	Details    string       `json:"details,omitempty"`
	Galleries  []GalleryRef `json:"galleries,omitempty"`
	Performers []string     `json:"performers,omitempty"`
	Groups     []AudioGroup `json:"movies,omitempty"`
	Tags       []string     `json:"tags,omitempty"`
	// Markers    []AudioMarker `json:"markers,omitempty"`
	Files     []string      `json:"files,omitempty"`
	Cover     string        `json:"cover,omitempty"`
	CreatedAt json.JSONTime `json:"created_at,omitempty"`
	UpdatedAt json.JSONTime `json:"updated_at,omitempty"`

	// deprecated - for import only
	LastPlayedAt json.JSONTime `json:"last_played_at,omitempty"`

	ResumeTime float64 `json:"resume_time,omitempty"`

	// deprecated - for import only
	PlayCount int `json:"play_count,omitempty"`

	PlayHistory []json.JSONTime `json:"play_history,omitempty"`
	OHistory    []json.JSONTime `json:"o_history,omitempty"`

	PlayDuration float64 `json:"play_duration,omitempty"`

	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

func (s Audio) Filename(id int, basename string, hash string) string {
	ret := fsutil.SanitiseBasename(s.Title)
	if ret == "" {
		ret = basename
	}

	if hash != "" {
		ret += "." + hash
	} else {
		// audios may have no file and therefore no hash
		ret += "." + strconv.Itoa(id)
	}

	return ret + ".json"
}

func LoadAudioFile(filePath string) (*Audio, error) {
	var audio Audio
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&audio)
	if err != nil {
		return nil, err
	}
	return &audio, nil
}

func SaveAudioFile(filePath string, audio *Audio) error {
	if audio == nil {
		return fmt.Errorf("audio must not be nil")
	}
	return marshalToFile(filePath, audio)
}
