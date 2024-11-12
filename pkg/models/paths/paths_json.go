package paths

import (
	"path/filepath"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

type JSONPaths struct {
	Metadata string

	ScrapedFile string

	Performers   string
	Scenes       string
	Images       string
	Galleries    string
	Studios      string
	Tags         string
	Groups       string
	Files        string
	SavedFilters string
}

func newJSONPaths(baseDir string) *JSONPaths {
	jp := JSONPaths{}
	jp.Metadata = baseDir
	jp.ScrapedFile = filepath.Join(baseDir, "scraped.json")
	jp.Performers = filepath.Join(baseDir, "performers")
	jp.Scenes = filepath.Join(baseDir, "scenes")
	jp.Images = filepath.Join(baseDir, "images")
	jp.Galleries = filepath.Join(baseDir, "galleries")
	jp.Studios = filepath.Join(baseDir, "studios")
	jp.Groups = filepath.Join(baseDir, "movies")
	jp.Tags = filepath.Join(baseDir, "tags")
	jp.Files = filepath.Join(baseDir, "files")
	jp.SavedFilters = filepath.Join(baseDir, "saved_filters")
	return &jp
}

func GetJSONPaths(baseDir string) *JSONPaths {
	jp := newJSONPaths(baseDir)
	return jp
}

func EmptyJSONDirs(baseDir string) {
	jsonPaths := GetJSONPaths(baseDir)
	_ = fsutil.EmptyDir(jsonPaths.Scenes)
	_ = fsutil.EmptyDir(jsonPaths.Images)
	_ = fsutil.EmptyDir(jsonPaths.Galleries)
	_ = fsutil.EmptyDir(jsonPaths.Performers)
	_ = fsutil.EmptyDir(jsonPaths.Studios)
	_ = fsutil.EmptyDir(jsonPaths.Groups)
	_ = fsutil.EmptyDir(jsonPaths.Tags)
	_ = fsutil.EmptyDir(jsonPaths.Files)
	_ = fsutil.EmptyDir(jsonPaths.SavedFilters)
}

func EnsureJSONDirs(baseDir string) {
	jsonPaths := GetJSONPaths(baseDir)
	if err := fsutil.EnsureDir(jsonPaths.Metadata); err != nil {
		logger.Warnf("couldn't create directories for Metadata: %v", err)
	}
	if err := fsutil.EnsureDir(jsonPaths.Scenes); err != nil {
		logger.Warnf("couldn't create directories for Scenes: %v", err)
	}
	if err := fsutil.EnsureDir(jsonPaths.Images); err != nil {
		logger.Warnf("couldn't create directories for Images: %v", err)
	}
	if err := fsutil.EnsureDir(jsonPaths.Galleries); err != nil {
		logger.Warnf("couldn't create directories for Galleries: %v", err)
	}
	if err := fsutil.EnsureDir(jsonPaths.Performers); err != nil {
		logger.Warnf("couldn't create directories for Performers: %v", err)
	}
	if err := fsutil.EnsureDir(jsonPaths.Studios); err != nil {
		logger.Warnf("couldn't create directories for Studios: %v", err)
	}
	if err := fsutil.EnsureDir(jsonPaths.Groups); err != nil {
		logger.Warnf("couldn't create directories for Groups: %v", err)
	}
	if err := fsutil.EnsureDir(jsonPaths.Tags); err != nil {
		logger.Warnf("couldn't create directories for Tags: %v", err)
	}
	if err := fsutil.EnsureDir(jsonPaths.Files); err != nil {
		logger.Warnf("couldn't create directories for Files: %v", err)
	}
	if err := fsutil.EnsureDir(jsonPaths.SavedFilters); err != nil {
		logger.Warnf("couldn't create directories for Saved Filters: %v", err)
	}
}
