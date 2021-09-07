package paths

import (
	"path/filepath"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
)

type JSONPaths struct {
	Metadata string

	MappingsFile string
	ScrapedFile  string

	Performers string
	Scenes     string
	Images     string
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
	jp.Images = filepath.Join(baseDir, "images")
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
	err := utils.EnsureDir(jsonPaths.Metadata)
	if err != nil {
		logger.Warnf("couldn't create directories for Metadata: %v", err)
	}
	err = utils.EnsureDir(jsonPaths.Scenes)
	if err != nil {
		logger.Warnf("couldn't create directories for Scenes: %v", err)
	}
	err = utils.EnsureDir(jsonPaths.Images)
	if err != nil {
		logger.Warnf("couldn't create directories for Images: %v", err)
	}
	err = utils.EnsureDir(jsonPaths.Galleries)
	if err != nil {
		logger.Warnf("couldn't create directories for Galleries: %v", err)
	}
	err = utils.EnsureDir(jsonPaths.Performers)
	if err != nil {
		logger.Warnf("couldn't create directories for Performers: %v", err)
	}
	err = utils.EnsureDir(jsonPaths.Studios)
	if err != nil {
		logger.Warnf("couldn't create directories for Studios: %v", err)
	}
	err = utils.EnsureDir(jsonPaths.Movies)
	if err != nil {
		logger.Warnf("couldn't create directories for Movies: %v", err)
	}
	err = utils.EnsureDir(jsonPaths.Tags)
	if err != nil {
		logger.Warnf("couldn't create directories for Tags: %v", err)
	}
}

func (jp *JSONPaths) PerformerJSONPath(checksum string) string {
	return filepath.Join(jp.Performers, checksum+".json")
}

func (jp *JSONPaths) SceneJSONPath(checksum string) string {
	return filepath.Join(jp.Scenes, checksum+".json")
}

func (jp *JSONPaths) ImageJSONPath(checksum string) string {
	return filepath.Join(jp.Images, checksum+".json")
}

func (jp *JSONPaths) GalleryJSONPath(checksum string) string {
	return filepath.Join(jp.Galleries, checksum+".json")
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
