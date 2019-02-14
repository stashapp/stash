package paths

import (
	"path/filepath"
)

type jsonPaths struct {
	MappingsFile string
	ScrapedFile  string

	Performers string
	Scenes     string
	Galleries  string
	Studios    string
}

func newJSONPaths(p Paths) *jsonPaths {
	jp := jsonPaths{}
	jp.MappingsFile = filepath.Join(p.Config.Metadata, "mappings.json")
	jp.ScrapedFile = filepath.Join(p.Config.Metadata, "scraped.json")
	jp.Performers = filepath.Join(p.Config.Metadata, "performers")
	jp.Scenes = filepath.Join(p.Config.Metadata, "scenes")
	jp.Galleries = filepath.Join(p.Config.Metadata, "galleries")
	jp.Studios = filepath.Join(p.Config.Metadata, "studios")
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
