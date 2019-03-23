package paths

import (
	"os/user"
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

func GetHomeDirectory() string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.HomeDir
}

func GetConfigDirectory() string {
	return filepath.Join(GetHomeDirectory(), ".stash")
}

func GetDefaultDatabaseFilePath() string {
	return filepath.Join(GetConfigDirectory(), "stash-go.sqlite")
}

func GetDefaultConfigFilePath() string {
	return filepath.Join(GetConfigDirectory(), "config.yml")
}
