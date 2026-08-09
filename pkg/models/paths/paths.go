// Package paths provides functions to return paths to various resources.
package paths

import (
	"path/filepath"

	"github.com/stashapp/stash/pkg/fsutil"
)

type Paths struct {
	Generated *generatedPaths

	Scene        *scenePaths
	SceneMarkers *sceneMarkerPaths
	Blobs        string
}

func NewPaths(generatedPath string, blobsPath string) Paths {
	p := Paths{}
	p.Generated = newGeneratedPaths(generatedPath)

	p.Scene = newScenePaths(p)
	p.SceneMarkers = newSceneMarkerPaths(p)
	p.Blobs = blobsPath

	return p
}

func GetStashHomeDirectory() string {
	return filepath.Join(fsutil.GetHomeDirectory(), ".stash")
}

func GetDefaultDatabaseFilePath() string {
	return filepath.Join(GetStashHomeDirectory(), "stash-go.sqlite")
}
