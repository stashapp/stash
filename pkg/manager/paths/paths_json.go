package paths

import (
	"github.com/stashapp/stash/pkg/manager/config"
	"path/filepath"
)

type jsonPaths struct {
	MappingsFile string
	ScrapedFile  string

	Performers string
	Scenes     string
	Galleries  string
	Studios    string
	Dvds       string
}

func newJSONPaths() *jsonPaths {
	jp := jsonPaths{}
	jp.MappingsFile = filepath.Join(config.GetMetadataPath(), "mappings.json")
	jp.ScrapedFile = filepath.Join(config.GetMetadataPath(), "scraped.json")
	jp.Performers = filepath.Join(config.GetMetadataPath(), "performers")
	jp.Scenes = filepath.Join(config.GetMetadataPath(), "scenes")
	jp.Galleries = filepath.Join(config.GetMetadataPath(), "galleries")
	jp.Studios = filepath.Join(config.GetMetadataPath(), "studios")
	jp.Dvds = filepath.Join(config.GetMetadataPath(), "dvds")
	return &jp
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

func (jp *jsonPaths) DvdJSONPath(checksum string) string {
	return filepath.Join(jp.Dvds, checksum+".json")
}
