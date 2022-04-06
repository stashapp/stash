package paths

import (
	"path/filepath"

	"github.com/stashapp/stash/pkg/fsutil"
)

type Config interface {
	GetGeneratedPath() string
	GetDataPath() string
}

type Paths struct {
	Generated *generatedPaths
	Data      *dataPaths

	Scene        *scenePaths
	SceneMarkers *sceneMarkerPaths
}

func NewPaths(config Config) *Paths {
	p := Paths{}
	p.Generated = newGeneratedPaths(config.GetGeneratedPath())
	p.Data = newDataPaths(config.GetDataPath())

	p.Scene = newScenePaths(p)
	p.SceneMarkers = newSceneMarkerPaths(p)
	return &p
}

func GetStashHomeDirectory() string {
	return filepath.Join(fsutil.GetHomeDirectory(), ".stash")
}

func GetDefaultDatabaseFilePath() string {
	return filepath.Join(GetStashHomeDirectory(), "stash-go.sqlite")
}
