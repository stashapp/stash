package models

type JoinReader interface {
	// GetScenePerformers(sceneID int) ([]PerformersScenes, error)
	GetSceneMovies(sceneID int) ([]MoviesScenes, error)
	// GetSceneTags(sceneID int) ([]ScenesTags, error)
}

type JoinWriter interface {
	CreatePerformersScenes(newJoins []PerformersScenes) error
	// AddPerformerScene(sceneID int, performerID int) (bool, error)
	UpdatePerformersScenes(sceneID int, updatedJoins []PerformersScenes) error
	// DestroyPerformersScenes(sceneID int) error
	CreateMoviesScenes(newJoins []MoviesScenes) error
	// AddMoviesScene(sceneID int, movieID int, sceneIdx *int) (bool, error)
	UpdateMoviesScenes(sceneID int, updatedJoins []MoviesScenes) error
	// DestroyMoviesScenes(sceneID int) error
	// CreateScenesTags(newJoins []ScenesTags) error
	UpdateScenesTags(sceneID int, updatedJoins []ScenesTags) error
	// AddSceneTag(sceneID int, tagID int) (bool, error)
	// DestroyScenesTags(sceneID int) error
	// CreateSceneMarkersTags(newJoins []SceneMarkersTags) error
	UpdateSceneMarkersTags(sceneMarkerID int, updatedJoins []SceneMarkersTags) error
	// DestroySceneMarkersTags(sceneMarkerID int, updatedJoins []SceneMarkersTags) error
	// DestroyScenesGalleries(sceneID int) error
	// DestroyScenesMarkers(sceneID int) error
	UpdatePerformersGalleries(galleryID int, updatedJoins []PerformersGalleries) error
	UpdateGalleriesTags(galleryID int, updatedJoins []GalleriesTags) error
	UpdateGalleriesImages(imageID int, updatedJoins []GalleriesImages) error
	UpdatePerformersImages(imageID int, updatedJoins []PerformersImages) error
	UpdateImagesTags(imageID int, updatedJoins []ImagesTags) error
}

type JoinReaderWriter interface {
	JoinReader
	JoinWriter
}
