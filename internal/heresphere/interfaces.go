package heresphere

import (
	"github.com/stashapp/stash/pkg/models"
)

type sceneFinder interface {
	models.SceneQueryer
	models.SceneGetter
	models.SceneReader
	models.SceneWriter
}

type sceneMarkerFinder interface {
	models.SceneMarkerFinder
	models.SceneMarkerCreator
	models.SceneMarkerReader
}

type tagFinder interface {
	models.TagFinder
	models.TagCreator
}

type fileFinder interface {
	models.FileFinder
	models.FileReader
	models.FileDestroyer
}

type savedfilterFinder interface {
	models.SavedFilterReader
}

type performerFinder interface {
	models.PerformerFinder
}

type galleryFinder interface {
	models.GalleryFinder
}

type movieFinder interface {
	models.MovieFinder
}

type studioFinder interface {
	models.StudioFinder
}
