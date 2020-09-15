package paths

import (
	"path/filepath"

	"github.com/stashapp/stash/pkg/utils"
)

type JSONPaths struct {
	Metadata string

	MappingsFile string
	ScrapedFile  string

	Performers string
	Scenes     string
	Galleries  string
	Studios    string
	Tags       string
	Movies     string
}

func newJSONPaths(baseDir string) *JSONPaths {
	jp := JSONPaths{}
	jp.Metadata = baseDir
	jp.MappingsFile = filepath.Join(baseDir, "mappings.json")
	jp.ScrapedFile = filepath.Join(baseDir, "scraped.json")
	jp.Performers = filepath.Join(baseDir, "performers")
	jp.Scenes = filepath.Join(baseDir, "scenes")
	jp.Galleries = filepath.Join(baseDir, "galleries")
	jp.Studios = filepath.Join(baseDir, "studios")
	jp.Movies = filepath.Join(baseDir, "movies")
	jp.Tags = filepath.Join(baseDir, "tags")
	return &jp
}

func GetJSONPaths(baseDir string) *JSONPaths {
	jp := newJSONPaths(baseDir)
	return jp
}

func EnsureJSONDirs(baseDir string) {
	jsonPaths := GetJSONPaths(baseDir)
	utils.EnsureDir(jsonPaths.Metadata)
	utils.EnsureDir(jsonPaths.Scenes)
	utils.EnsureDir(jsonPaths.Galleries)
	utils.EnsureDir(jsonPaths.Performers)
	utils.EnsureDir(jsonPaths.Studios)
	utils.EnsureDir(jsonPaths.Movies)
	utils.EnsureDir(jsonPaths.Tags)
}

func (jp *JSONPaths) PerformerJSONPath(checksum string) string {
	return filepath.Join(jp.Performers, checksum+".json")
}

func (jp *JSONPaths) SceneJSONPath(checksum string) string {
	return filepath.Join(jp.Scenes, checksum+".json")
}

func (jp *JSONPaths) StudioJSONPath(checksum string) string {
	return filepath.Join(jp.Studios, checksum+".json")
}

func (jp *JSONPaths) TagJSONPath(checksum string) string {
	return filepath.Join(jp.Tags, checksum+".json")
}

func (jp *JSONPaths) MovieJSONPath(checksum string) string {
	return filepath.Join(jp.Movies, checksum+".json")
}
