package paths

import (
	"github.com/stashapp/stash/pkg/utils"
	"path/filepath"
)

type Paths struct {
	Generated *generatedPaths
	JSON      *jsonPaths

	Gallery      *galleryPaths
	Scene        *scenePaths
	SceneMarkers *sceneMarkerPaths
}

func NewPaths() *Paths {
	p := Paths{}
	p.Generated = newGeneratedPaths()
	p.JSON = newJSONPaths()

	p.Gallery = newGalleryPaths()
	p.Scene = newScenePaths(p)
	p.SceneMarkers = newSceneMarkerPaths(p)
	return &p
}

func GetConfigDirectory() string {
	return filepath.Join(utils.GetHomeDirectory(), ".stash")
}

func GetDefaultDatabaseFilePath() string {
	return filepath.Join(GetConfigDirectory(), "stash-go-tagger.sqlite")
}

func GetDefaultConfigFilePath() string {
	return filepath.Join(GetConfigDirectory(), "config-tagger.yml")
}

func GetSSLKey() string {
	return filepath.Join(GetConfigDirectory(), "stash.key")
}

func GetSSLCert() string {
	return filepath.Join(GetConfigDirectory(), "stash.crt")
}
