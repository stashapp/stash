package paths

import (
	"github.com/stashapp/stash/pkg/manager/jsonschema"
)

type Paths struct {
	Config    *jsonschema.Config
	Generated *generatedPaths
	JSON      *jsonPaths

	Gallery      *galleryPaths
	Scene        *scenePaths
	SceneMarkers *sceneMarkerPaths
}

func NewPaths(config *jsonschema.Config) *Paths {
	p := Paths{}
	p.Config = config
	p.Generated = newGeneratedPaths(p)
	p.JSON = newJSONPaths(p)

	p.Gallery = newGalleryPaths(p.Config)
	p.Scene = newScenePaths(p)
	p.SceneMarkers = newSceneMarkerPaths(p)
	return &p
}
