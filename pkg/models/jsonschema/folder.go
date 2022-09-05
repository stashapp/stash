package jsonschema

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models/json"
)

type Folder struct {
	BaseDirEntry

	Path string `json:"path,omitempty"`

	CreatedAt json.JSONTime `json:"created_at,omitempty"`
	UpdatedAt json.JSONTime `json:"updated_at,omitempty"`
}

func (f *Folder) Filename() string {
	// prefix with the path depth so that we can import lower-level folders first
	depth := strings.Count(f.Path, string(filepath.Separator))

	// hash the full path for a unique filename
	hash := md5.FromString(f.Path)

	basename := filepath.Base(f.Path)

	return fmt.Sprintf("%2x.%s.%s.json", depth, basename, hash)
}

func LoadFolderFile(filePath string) (*Folder, error) {
	var folder Folder
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&folder)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func SaveFolderFile(filePath string, folder *Folder) error {
	if folder == nil {
		return fmt.Errorf("folder must not be nil")
	}
	return marshalToFile(filePath, folder)
}
