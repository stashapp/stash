package models

type SceneReader interface {
	Find(id int) (*Scene, error)
	FindMany(ids []int) ([]*Scene, error)
	FindByChecksum(checksum string) (*Scene, error)
	FindByOSHash(oshash string) (*Scene, error)
	FindByPath(path string) (*Scene, error)
	FindByPerformerID(performerID int) ([]*Scene, error)
	CountByPerformerID(performerID int) (int, error)
	// FindByStudioID(studioID int) ([]*Scene, error)
	FindByMovieID(movieID int) ([]*Scene, error)
	CountByMovieID(movieID int) (int, error)
	// Count() (int, error)
	// SizeCount() (string, error)
	CountByStudioID(studioID int) (int, error)
	CountByTagID(tagID int) (int, error)
	// CountMissingChecksum() (int, error)
	// CountMissingOSHash() (int, error)
	// Wall(q *string) ([]*Scene, error)
	All() ([]*Scene, error)
	// Query(sceneFilter *SceneFilterType, findFilter *FindFilterType) ([]*Scene, int)
	// QueryAllByPathRegex(regex string) ([]*Scene, error)
	QueryByPathRegex(findFilter *FindFilterType) ([]*Scene, int)
	GetSceneCover(sceneID int) ([]byte, error)
	GetMovies(sceneID int) ([]MoviesScenes, error)
	GetTagIDs(imageID int) ([]int, error)
	GetPerformerIDs(imageID int) ([]int, error)
	GetStashIDs(performerID int) ([]*StashID, error)
}

type SceneWriter interface {
	Create(newScene Scene) (*Scene, error)
	Update(updatedScene ScenePartial) (*Scene, error)
	UpdateFull(updatedScene Scene) (*Scene, error)
	IncrementOCounter(id int) (int, error)
	DecrementOCounter(id int) (int, error)
	ResetOCounter(id int) (int, error)
	UpdateFileModTime(id int, modTime NullSQLiteTimestamp) error
	Destroy(id int) error
	UpdateSceneCover(sceneID int, cover []byte) error
	// DestroySceneCover(sceneID int) error
	UpdatePerformers(sceneID int, performerIDs []int) error
	UpdateTags(sceneID int, tagIDs []int) error
	UpdateMovies(sceneID int, movies []MoviesScenes) error
	UpdateStashIDs(sceneID int, stashIDs []StashID) error
}

type SceneReaderWriter interface {
	SceneReader
	SceneWriter
}
