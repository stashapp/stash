package paths

import (
	"github.com/stashapp/stash/manager/jsonschema"
	"github.com/stashapp/stash/utils"
)

type Paths struct {
	Config     *jsonschema.Config
	Generated  *generatedPaths
	JSON       *jsonPaths

	Gallery *galleryPaths
	Scene      *scenePaths
	SceneMarkers *sceneMarkerPaths
}

func RefreshPaths() *Paths {
	ensureConfigFile()
	return newPaths()
}

func newPaths() *Paths {
	p := Paths{}
	p.Config = jsonschema.LoadConfigFile(StaticPaths.ConfigFile)
	p.Generated = newGeneratedPaths(p)
	p.JSON = newJSONPaths(p)

	p.Gallery = newGalleryPaths(p.Config)
	p.Scene = newScenePaths(p)
	p.SceneMarkers = newSceneMarkerPaths(p)
	return &p
}

func ensureConfigFile() {
	configFileExists, _ := utils.FileExists(StaticPaths.ConfigFile) // TODO: Verify JSON is correct.  Pass verified
	if configFileExists {
		return
	}

	panic("No config file found")
}