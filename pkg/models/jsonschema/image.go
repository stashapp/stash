package jsonschema

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models/json"
)

type Image struct {
	Title  string `json:"title,omitempty"`
	Code   string `json:"code,omitempty"`
	Studio string `json:"studio,omitempty"`
	Rating int    `json:"rating,omitempty"`

	// deprecated - for import only
	URL string `json:"url,omitempty"`

	URLs         []string      `json:"urls,omitempty"`
	Date         string        `json:"date,omitempty"`
	Details      string        `json:"details,omitempty"`
	Photographer string        `json:"photographer,omitempty"`
	Organized    bool          `json:"organized,omitempty"`
	OCounter     int           `json:"o_counter,omitempty"`
	Galleries    []GalleryRef  `json:"galleries,omitempty"`
	Performers   []string      `json:"performers,omitempty"`
	Tags         []string      `json:"tags,omitempty"`
	Files        []string      `json:"files,omitempty"`
	CreatedAt    json.JSONTime `json:"created_at,omitempty"`
	UpdatedAt    json.JSONTime `json:"updated_at,omitempty"`
}

func (s Image) Filename(basename string, hash string) string {
	ret := fsutil.SanitiseBasename(s.Title)
	if ret == "" {
		ret = basename
	}

	if hash != "" {
		ret += "." + hash
	}

	return ret + ".json"
}

func LoadImageFile(filePath string) (*Image, error) {
	var image Image
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&image)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func SaveImageFile(filePath string, image *Image) error {
	if image == nil {
		return fmt.Errorf("image must not be nil")
	}
	return marshalToFile(filePath, image)
}
