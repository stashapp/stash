package jsonschema

import (
	"fmt"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models/json"
)

type GalleryChapter struct {
	Title      string        `json:"title,omitempty"`
	ImageIndex int           `json:"image_index,omitempty"`
	CreatedAt  json.JSONTime `json:"created_at,omitempty"`
	UpdatedAt  json.JSONTime `json:"updated_at,omitempty"`
}

type Gallery struct {
	ZipFiles     []string         `json:"zip_files,omitempty"`
	FolderPath   string           `json:"folder_path,omitempty"`
	Title        string           `json:"title,omitempty"`
	Code         string           `json:"code,omitempty"`
	URLs         []string         `json:"urls,omitempty"`
	Date         string           `json:"date,omitempty"`
	Details      string           `json:"details,omitempty"`
	Photographer string           `json:"photographer,omitempty"`
	Rating       int              `json:"rating,omitempty"`
	Organized    bool             `json:"organized,omitempty"`
	Chapters     []GalleryChapter `json:"chapters,omitempty"`
	Studio       string           `json:"studio,omitempty"`
	Performers   []string         `json:"performers,omitempty"`
	Tags         []string         `json:"tags,omitempty"`
	CreatedAt    json.JSONTime    `json:"created_at,omitempty"`
	UpdatedAt    json.JSONTime    `json:"updated_at,omitempty"`

	// deprecated - for import only
	URL string `json:"url,omitempty"`
}

func (s Gallery) Filename(basename string, hash string) string {
	ret := fsutil.SanitiseBasename(basename)

	if ret != "" {
		ret += "."
	}
	ret += hash

	return ret + ".json"
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

// GalleryRef is used to identify a Gallery.
// Only one field should be populated.
type GalleryRef struct {
	ZipFiles   []string `json:"zip_files,omitempty"`
	FolderPath string   `json:"folder_path,omitempty"`
	// Title is used only if FolderPath and ZipPaths is empty
	Title string `json:"title,omitempty"`
}

func (r GalleryRef) String() string {
	switch {
	case r.FolderPath != "":
		return "{ folder: " + r.FolderPath + " }"
	case len(r.ZipFiles) > 0:
		return "{ zipFiles: [" + strings.Join(r.ZipFiles, ", ") + "] }"
	default:
		return "{ title: " + r.Title + " }"
	}
}
