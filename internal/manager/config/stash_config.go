package config

import (
	"path/filepath"

	"github.com/stashapp/stash/pkg/fsutil"
)

// Stash configuration details
type StashConfigInput struct {
	Path         string `json:"path"`
	ExcludeVideo bool   `json:"excludeVideo"`
	ExcludeImage bool   `json:"excludeImage"`
}

type StashConfig struct {
	Path         string `json:"path"`
	ExcludeVideo bool   `json:"excludeVideo"`
	ExcludeImage bool   `json:"excludeImage"`
}

type StashConfigs []*StashConfig

func (s StashConfigs) GetStashFromPath(path string) *StashConfig {
	for _, f := range s {
		if fsutil.IsPathInDir(f.Path, filepath.Dir(path)) {
			return f
		}
	}
	return nil
}

func (s StashConfigs) GetStashFromDirPath(dirPath string) *StashConfig {
	for _, f := range s {
		if fsutil.IsPathInDir(f.Path, dirPath) {
			return f
		}
	}
	return nil
}
