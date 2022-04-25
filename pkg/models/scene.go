package models

type PHashDuplicationCriterionInput struct {
	Duplicated *bool `json:"duplicated"`
	// Currently unimplemented
	Distance *int `json:"distance"`
}

type SceneFilterType struct {
	And     *SceneFilterType      `json:"AND"`
	Or      *SceneFilterType      `json:"OR"`
	Not     *SceneFilterType      `json:"NOT"`
	Title   *StringCriterionInput `json:"title"`
	Details *StringCriterionInput `json:"details"`
	// Filter by file oshash
	Oshash *StringCriterionInput `json:"oshash"`
	// Filter by file checksum
	Checksum *StringCriterionInput `json:"checksum"`
	// Filter by file phash
	Phash *StringCriterionInput `json:"phash"`
	// Filter by path
	Path *StringCriterionInput `json:"path"`
	// Filter by rating
	Rating *IntCriterionInput `json:"rating"`
	// Filter by organized
	Organized *bool `json:"organized"`
	// Filter by o-counter
	OCounter *IntCriterionInput `json:"o_counter"`
	// Filter Scenes that have an exact phash match available
	Duplicated *PHashDuplicationCriterionInput `json:"duplicated"`
	// Filter by resolution
	Resolution *ResolutionCriterionInput `json:"resolution"`
	// Filter by duration (in seconds)
	Duration *IntCriterionInput `json:"duration"`
	// Filter to only include scenes which have markers. `true` or `false`
	HasMarkers *string `json:"has_markers"`
	// Filter to only include scenes missing this property
	IsMissing *string `json:"is_missing"`
	// Filter to only include scenes with this studio
	Studios *HierarchicalMultiCriterionInput `json:"studios"`
	// Filter to only include scenes with this movie
	Movies *MultiCriterionInput `json:"movies"`
	// Filter to only include scenes with these tags
	Tags *HierarchicalMultiCriterionInput `json:"tags"`
	// Filter by tag count
	TagCount *IntCriterionInput `json:"tag_count"`
	// Filter to only include scenes with performers with these tags
	PerformerTags *HierarchicalMultiCriterionInput `json:"performer_tags"`
	// Filter scenes that have performers that have been favorited
	PerformerFavorite *bool `json:"performer_favorite"`
	// Filter scenes by performer age at time of scene
	PerformerAge *IntCriterionInput `json:"performer_age"`
	// Filter to only include scenes with these performers
	Performers *MultiCriterionInput `json:"performers"`
	// Filter by performer count
	PerformerCount *IntCriterionInput `json:"performer_count"`
	// Filter by StashID
	StashID *StringCriterionInput `json:"stash_id"`
	// Filter by url
	URL *StringCriterionInput `json:"url"`
	// Filter by interactive
	Interactive *bool `json:"interactive"`
	// Filter by InteractiveSpeed
	InteractiveSpeed *IntCriterionInput `json:"interactive_speed"`

	Captions *StringCriterionInput `json:"captions"`
}

type SceneQueryOptions struct {
	QueryOptions
	SceneFilter *SceneFilterType

	TotalDuration bool
	TotalSize     bool
}

type SceneQueryResult struct {
	QueryResult
	TotalDuration float64
	TotalSize     float64

	finder     SceneFinder
	scenes     []*Scene
	resolveErr error
}

type SceneDestroyInput struct {
	ID              string `json:"id"`
	DeleteFile      *bool  `json:"delete_file"`
	DeleteGenerated *bool  `json:"delete_generated"`
}

type ScenesDestroyInput struct {
	Ids             []string `json:"ids"`
	DeleteFile      *bool    `json:"delete_file"`
	DeleteGenerated *bool    `json:"delete_generated"`
}

func NewSceneQueryResult(finder SceneFinder) *SceneQueryResult {
	return &SceneQueryResult{
		finder: finder,
	}
}

func (r *SceneQueryResult) Resolve() ([]*Scene, error) {
	// cache results
	if r.scenes == nil && r.resolveErr == nil {
		r.scenes, r.resolveErr = r.finder.FindMany(r.IDs)
	}
	return r.scenes, r.resolveErr
}

type SceneFinder interface {
	// TODO - rename this to Find and remove existing method
	FindMany(ids []int) ([]*Scene, error)
}

type SceneReader interface {
	SceneFinder
	// TODO - remove this in another PR
	Find(id int) (*Scene, error)
	FindByChecksum(checksum string) (*Scene, error)
	FindByOSHash(oshash string) (*Scene, error)
	FindByPath(path string) (*Scene, error)
	FindByPerformerID(performerID int) ([]*Scene, error)
	FindByGalleryID(performerID int) ([]*Scene, error)
	FindDuplicates(distance int) ([][]*Scene, error)
	CountByPerformerID(performerID int) (int, error)
	// FindByStudioID(studioID int) ([]*Scene, error)
	FindByMovieID(movieID int) ([]*Scene, error)
	CountByMovieID(movieID int) (int, error)
	Count() (int, error)
	Size() (float64, error)
	Duration() (float64, error)
	// SizeCount() (string, error)
	CountByStudioID(studioID int) (int, error)
	CountByTagID(tagID int) (int, error)
	CountMissingChecksum() (int, error)
	CountMissingOSHash() (int, error)
	Wall(q *string) ([]*Scene, error)
	All() ([]*Scene, error)
	Query(options SceneQueryOptions) (*SceneQueryResult, error)
	GetCaptions(sceneID int) ([]*SceneCaption, error)
	GetCover(sceneID int) ([]byte, error)
	GetMovies(sceneID int) ([]MoviesScenes, error)
	GetTagIDs(sceneID int) ([]int, error)
	GetGalleryIDs(sceneID int) ([]int, error)
	GetPerformerIDs(sceneID int) ([]int, error)
	GetStashIDs(sceneID int) ([]*StashID, error)
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
	UpdateCaptions(id int, captions []*SceneCaption) error
	UpdateCover(sceneID int, cover []byte) error
	DestroyCover(sceneID int) error
	UpdatePerformers(sceneID int, performerIDs []int) error
	UpdateTags(sceneID int, tagIDs []int) error
	UpdateGalleries(sceneID int, galleryIDs []int) error
	UpdateMovies(sceneID int, movies []MoviesScenes) error
	UpdateStashIDs(sceneID int, stashIDs []StashID) error
}

type SceneReaderWriter interface {
	SceneReader
	SceneWriter
}
