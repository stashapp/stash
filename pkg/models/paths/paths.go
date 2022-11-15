package paths

import (
	"path/filepath"

	"github.com/stashapp/stash/pkg/fsutil"
)

type Paths struct {
	Generated *generatedPaths

	Scene        *scenePaths
	SceneMarkers *sceneMarkerPaths
}

func NewPaths(generatedPath string) Paths {
	p := Paths{}
	p.Generated = newGeneratedPaths(generatedPath)

	p.Scene = newScenePaths(p)
	p.SceneMarkers = newSceneMarkerPaths(p)
	return p
}

func GetStashHomeDirectory() string {
	return filepath.Join(fsutil.GetHomeDirectory(), ".stash")
}

func GetDefaultDatabaseFilePath() string {
	return filepath.Join(GetStashHomeDirectory(), "stash-go.sqlite")
}
