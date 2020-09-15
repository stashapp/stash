package paths

import (
	"path/filepath"

	"github.com/stashapp/stash/pkg/utils"
)

type Paths struct {
	Generated *generatedPaths

	Gallery      *galleryPaths
	Scene        *scenePaths
	SceneMarkers *sceneMarkerPaths
}

func NewPaths() *Paths {
	p := Paths{}
	p.Generated = newGeneratedPaths()

	p.Gallery = newGalleryPaths()
	p.Scene = newScenePaths(p)
	p.SceneMarkers = newSceneMarkerPaths(p)
	return &p
}

func GetConfigDirectory() string {
	return filepath.Join(utils.GetHomeDirectory(), ".stash")
}

func GetDefaultDatabaseFilePath() string {
	return filepath.Join(GetConfigDirectory(), "stash-go.sqlite")
}

func GetDefaultConfigFilePath() string {
	return filepath.Join(GetConfigDirectory(), "config.yml")
}

func GetSSLKey() string {
	return filepath.Join(GetConfigDirectory(), "stash.key")
}

func GetSSLCert() string {
	return filepath.Join(GetConfigDirectory(), "stash.crt")
}
