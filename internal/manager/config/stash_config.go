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

// GetStashFromPath returns the most specific stash configuration containing path.
func (s StashConfigs) GetStashFromPath(path string) *StashConfig {
	return s.GetStashFromDirPath(filepath.Dir(path))
}

// GetStashFromDirPath returns the most specific stash configuration containing dirPath.
func (s StashConfigs) GetStashFromDirPath(dirPath string) *StashConfig {
	var ret *StashConfig
	longestPath := -1

	for _, f := range s {
		if f == nil {
			continue
		}

		path := filepath.Clean(f.Path)
		if fsutil.IsPathInDir(path, dirPath) && len(path) > longestPath {
			ret = f
			longestPath = len(path)
		}
	}

	return ret
}

// GetStashRootFromDirPath returns the topmost configured stash path containing dirPath.
func (s StashConfigs) GetStashRootFromDirPath(dirPath string) string {
	var ret string
	shortestPath := -1

	for _, f := range s {
		if f == nil {
			continue
		}

		path := filepath.Clean(f.Path)
		if fsutil.IsPathInDir(path, dirPath) && (shortestPath == -1 || len(path) < shortestPath) {
			ret = path
			shortestPath = len(path)
		}
	}

	return ret
}

func (s StashConfigs) Paths() []string {
	paths := make([]string, len(s))
	for i, c := range s {
		// #6618 - clean the path to ensure comparison works correctly
		paths[i] = filepath.Clean(c.Path)
	}
	return paths
}
