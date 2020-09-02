package paths

import (
	"path/filepath"

	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/utils"
)

type jsonPaths struct {
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

func newJSONPaths() *jsonPaths {
	jp := jsonPaths{}
	jp.Metadata = config.GetMetadataPath()
	jp.MappingsFile = filepath.Join(config.GetMetadataPath(), "mappings.json")
	jp.ScrapedFile = filepath.Join(config.GetMetadataPath(), "scraped.json")
	jp.Performers = filepath.Join(config.GetMetadataPath(), "performers")
	jp.Scenes = filepath.Join(config.GetMetadataPath(), "scenes")
	jp.Galleries = filepath.Join(config.GetMetadataPath(), "galleries")
	jp.Studios = filepath.Join(config.GetMetadataPath(), "studios")
	jp.Movies = filepath.Join(config.GetMetadataPath(), "movies")
	jp.Tags = filepath.Join(config.GetMetadataPath(), "tags")
	return &jp
}

func GetJSONPaths() *jsonPaths {
	jp := newJSONPaths()
	return jp
}

func EnsureJSONDirs() {
	jsonPaths := GetJSONPaths()
	utils.EnsureDir(jsonPaths.Metadata)
	utils.EnsureDir(jsonPaths.Scenes)
	utils.EnsureDir(jsonPaths.Galleries)
	utils.EnsureDir(jsonPaths.Performers)
	utils.EnsureDir(jsonPaths.Studios)
	utils.EnsureDir(jsonPaths.Movies)
	utils.EnsureDir(jsonPaths.Tags)
}

func (jp *jsonPaths) PerformerJSONPath(checksum string) string {
	return filepath.Join(jp.Performers, checksum+".json")
}

func (jp *jsonPaths) SceneJSONPath(checksum string) string {
	return filepath.Join(jp.Scenes, checksum+".json")
}

func (jp *jsonPaths) StudioJSONPath(checksum string) string {
	return filepath.Join(jp.Studios, checksum+".json")
}

func (jp *jsonPaths) TagJSONPath(checksum string) string {
	return filepath.Join(jp.Tags, checksum+".json")
}

func (jp *jsonPaths) MovieJSONPath(checksum string) string {
	return filepath.Join(jp.Movies, checksum+".json")
}
