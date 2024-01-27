package models

import "context"

// SceneGetter provides methods to get scenes by ID.
type SceneGetter interface {
	// TODO - rename this to Find and remove existing method
	FindMany(ctx context.Context, ids []int) ([]*Scene, error)
	Find(ctx context.Context, id int) (*Scene, error)
}

// SceneFinder provides methods to find scenes.
type SceneFinder interface {
	SceneGetter
	FindByFingerprints(ctx context.Context, fp []Fingerprint) ([]*Scene, error)
	FindByChecksum(ctx context.Context, checksum string) ([]*Scene, error)
	FindByOSHash(ctx context.Context, oshash string) ([]*Scene, error)
	FindByPath(ctx context.Context, path string) ([]*Scene, error)
	FindByFileID(ctx context.Context, fileID FileID) ([]*Scene, error)
	FindByPrimaryFileID(ctx context.Context, fileID FileID) ([]*Scene, error)
	FindByPerformerID(ctx context.Context, performerID int) ([]*Scene, error)
	FindByGalleryID(ctx context.Context, performerID int) ([]*Scene, error)
	FindByMovieID(ctx context.Context, movieID int) ([]*Scene, error)
	FindDuplicates(ctx context.Context, distance int, durationDiff float64) ([][]*Scene, error)
}

// SceneQueryer provides methods to query scenes.
type SceneQueryer interface {
	Query(ctx context.Context, options SceneQueryOptions) (*SceneQueryResult, error)
	QueryCount(ctx context.Context, sceneFilter *SceneFilterType, findFilter *FindFilterType) (int, error)
}

// SceneCounter provides methods to count scenes.
type SceneCounter interface {
	Count(ctx context.Context) (int, error)
	CountByPerformerID(ctx context.Context, performerID int) (int, error)
	CountByMovieID(ctx context.Context, movieID int) (int, error)
	CountByFileID(ctx context.Context, fileID FileID) (int, error)
	CountByStudioID(ctx context.Context, studioID int) (int, error)
	CountByTagID(ctx context.Context, tagID int) (int, error)
	CountMissingChecksum(ctx context.Context) (int, error)
	CountMissingOSHash(ctx context.Context) (int, error)
	OCount(ctx context.Context) (int, error)
	OCountByPerformerID(ctx context.Context, performerID int) (int, error)
	PlayCount(ctx context.Context) (int, error)
	UniqueScenePlayCount(ctx context.Context) (int, error)
}

// SceneCreator provides methods to create scenes.
type SceneCreator interface {
	Create(ctx context.Context, newScene *Scene, fileIDs []FileID) error
}

// SceneUpdater provides methods to update scenes.
type SceneUpdater interface {
	Update(ctx context.Context, updatedScene *Scene) error
	UpdatePartial(ctx context.Context, id int, updatedScene ScenePartial) (*Scene, error)
	UpdateCover(ctx context.Context, sceneID int, cover []byte) error
}

// SceneDestroyer provides methods to destroy scenes.
type SceneDestroyer interface {
	Destroy(ctx context.Context, id int) error
}

type SceneCreatorUpdater interface {
	SceneCreator
	SceneUpdater
}

// SceneReader provides all methods to read scenes.
type SceneReader interface {
	SceneFinder
	SceneQueryer
	SceneCounter

	URLLoader
	FileIDLoader
	GalleryIDLoader
	PerformerIDLoader
	TagIDLoader
	SceneMovieLoader
	StashIDLoader
	VideoFileLoader

	All(ctx context.Context) ([]*Scene, error)
	Wall(ctx context.Context, q *string) ([]*Scene, error)
	Size(ctx context.Context) (float64, error)
	Duration(ctx context.Context) (float64, error)
	PlayDuration(ctx context.Context) (float64, error)
	GetCover(ctx context.Context, sceneID int) ([]byte, error)
	HasCover(ctx context.Context, sceneID int) (bool, error)
}

// SceneWriter provides all methods to modify scenes.
type SceneWriter interface {
	SceneCreator
	SceneUpdater
	SceneDestroyer

	AddFileID(ctx context.Context, id int, fileID FileID) error
	AddGalleryIDs(ctx context.Context, sceneID int, galleryIDs []int) error
	AssignFiles(ctx context.Context, sceneID int, fileID []FileID) error
	IncrementOCounter(ctx context.Context, id int) (int, error)
	DecrementOCounter(ctx context.Context, id int) (int, error)
	ResetOCounter(ctx context.Context, id int) (int, error)
	SaveActivity(ctx context.Context, sceneID int, resumeTime *float64, playDuration *float64) (bool, error)
	IncrementWatchCount(ctx context.Context, sceneID int) (int, error)
}

// SceneReaderWriter provides all scene methods.
type SceneReaderWriter interface {
	SceneReader
	SceneWriter
}
