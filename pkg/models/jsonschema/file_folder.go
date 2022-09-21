package jsonschema

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models/json"
)

const (
	DirEntryTypeFolder = "folder"
	DirEntryTypeVideo  = "video"
	DirEntryTypeImage  = "image"
	DirEntryTypeFile   = "file"
)

type DirEntry interface {
	IsFile() bool
	Filename() string
	DirEntry() *BaseDirEntry
}

type BaseDirEntry struct {
	ZipFile string        `json:"zip_file,omitempty"`
	ModTime json.JSONTime `json:"mod_time"`

	Type string `json:"type,omitempty"`

	Path string `json:"path,omitempty"`

	CreatedAt json.JSONTime `json:"created_at,omitempty"`
	UpdatedAt json.JSONTime `json:"updated_at,omitempty"`
}

func (f *BaseDirEntry) DirEntry() *BaseDirEntry {
	return f
}

func (f *BaseDirEntry) IsFile() bool {
	return false
}

func (f *BaseDirEntry) Filename() string {
	// prefix with the path depth so that we can import lower-level files/folders first
	depth := strings.Count(f.Path, string(filepath.Separator))

	// hash the full path for a unique filename
	hash := md5.FromString(f.Path)

	basename := filepath.Base(f.Path)

	return fmt.Sprintf("%02x.%s.%s.json", depth, basename, hash)
}

type BaseFile struct {
	BaseDirEntry

	Fingerprints []Fingerprint `json:"fingerprints,omitempty"`
	Size         int64         `json:"size"`
}

func (f *BaseFile) IsFile() bool {
	return true
}

type Fingerprint struct {
	Type        string      `json:"type,omitempty"`
	Fingerprint interface{} `json:"fingerprint,omitempty"`
}

type VideoFile struct {
	*BaseFile
	Format     string  `json:"format,omitempty"`
	Width      int     `json:"width,omitempty"`
	Height     int     `json:"height,omitempty"`
	Duration   float64 `json:"duration,omitempty"`
	VideoCodec string  `json:"video_codec,omitempty"`
	AudioCodec string  `json:"audio_codec,omitempty"`
	FrameRate  float64 `json:"frame_rate,omitempty"`
	BitRate    int64   `json:"bitrate,omitempty"`

	Interactive      bool `json:"interactive,omitempty"`
	InteractiveSpeed *int `json:"interactive_speed,omitempty"`
}

type ImageFile struct {
	*BaseFile
	Format string `json:"format,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

func LoadFileFile(filePath string) (DirEntry, error) {
	r, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonParser := json.NewDecoder(bytes.NewReader(data))

	var bf BaseDirEntry
	if err := jsonParser.Decode(&bf); err != nil {
		return nil, err
	}

	jsonParser = json.NewDecoder(bytes.NewReader(data))

	switch bf.Type {
	case DirEntryTypeFolder:
		return &bf, nil
	case DirEntryTypeVideo:
		var vf VideoFile
		if err := jsonParser.Decode(&vf); err != nil {
			return nil, err
		}

		return &vf, nil
	case DirEntryTypeImage:
		var imf ImageFile
		if err := jsonParser.Decode(&imf); err != nil {
			return nil, err
		}

		return &imf, nil
	case DirEntryTypeFile:
		var bff BaseFile
		if err := jsonParser.Decode(&bff); err != nil {
			return nil, err
		}

		return &bff, nil
	default:
		return nil, errors.New("unknown file type")
	}
}

func SaveFileFile(filePath string, file DirEntry) error {
	if file == nil {
		return fmt.Errorf("file must not be nil")
	}
	return marshalToFile(filePath, file)
}
