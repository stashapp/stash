package jsonschema

import (
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models"
)

type SavedFilter struct {
	Mode         models.FilterMode      `db:"mode" json:"mode"`
	Name         string                 `db:"name" json:"name"`
	FindFilter   *models.FindFilterType `json:"find_filter"`
	ObjectFilter map[string]interface{} `json:"object_filter"`
	UIOptions    map[string]interface{} `json:"ui_options"`
}

func (s SavedFilter) Filename() string {
	ret := fsutil.SanitiseBasename(s.Name + "_" + s.Mode.String())
	return ret + ".json"
}

func LoadSavedFilterFile(filePath string) (*SavedFilter, error) {
	return loadFile[SavedFilter](filePath)
}

func SaveSavedFilterFile(filePath string, image *SavedFilter) error {
	return saveFile[SavedFilter](filePath, image)
}
