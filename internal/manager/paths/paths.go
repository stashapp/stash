package paths

import (
	"github.com/stashapp/stash/internal/manager/jsonschema"
	"github.com/stashapp/stash/internal/utils"
	"os"
	"os/user"
	"path/filepath"
)

type Paths struct {
	FixedPaths *fixedPaths
	Config     *jsonschema.Config
	Generated  *generatedPaths
	JSON       *jsonPaths

	Gallery *galleryPaths
	Scene      *scenePaths
	SceneMarkers *sceneMarkerPaths
}

func RefreshPaths() *Paths {
	fp := newFixedPaths()
	ensureConfigFile(fp)
	return newPaths(fp)
}

func newPaths(fp *fixedPaths) *Paths {
	p := Paths{}
	p.FixedPaths = fp
	p.Config = jsonschema.LoadConfigFile(p.FixedPaths.ConfigFile)
	p.Generated = newGeneratedPaths(p)
	p.JSON = newJSONPaths(p)

	p.Gallery = newGalleryPaths(p.Config)
	p.Scene = newScenePaths(p)
	p.SceneMarkers = newSceneMarkerPaths(p)
	return &p
}

func getExecutionDirectory() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func getHomeDirectory() string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.HomeDir
}

func ensureConfigFile(fp *fixedPaths) {
	configFileExists, _ := utils.FileExists(fp.ConfigFile) // TODO: Verify JSON is correct.  Pass verified
	if configFileExists {
		return
	}

	panic("No config file found")
}