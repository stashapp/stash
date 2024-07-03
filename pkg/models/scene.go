package models

import "context"

type PHashDuplicationCriterionInput struct {
	Duplicated *bool `json:"duplicated"`
	// Currently unimplemented
	Distance *int `json:"distance"`
}

type SceneFilterType struct {
	OperatorFilter[SceneFilterType]
	ID       *IntCriterionInput    `json:"id"`
	Title    *StringCriterionInput `json:"title"`
	Code     *StringCriterionInput `json:"code"`
	Details  *StringCriterionInput `json:"details"`
	Director *StringCriterionInput `json:"director"`
	// Filter by file oshash
	Oshash *StringCriterionInput `json:"oshash"`
	// Filter by file checksum
	Checksum *StringCriterionInput `json:"checksum"`
	// Filter by file phash
	Phash *StringCriterionInput `json:"phash"`
	// Filter by phash distance
	PhashDistance *PhashDistanceCriterionInput `json:"phash_distance"`
	// Filter by path
	Path *StringCriterionInput `json:"path"`
	// Filter by file count
	FileCount *IntCriterionInput `json:"file_count"`
	// Filter by rating expressed as 1-100
	Rating100 *IntCriterionInput `json:"rating100"`
	// Filter by organized
	Organized *bool `json:"organized"`
	// Filter by o-counter
	OCounter *IntCriterionInput `json:"o_counter"`
	// Filter Scenes that have an exact phash match available
	Duplicated *PHashDuplicationCriterionInput `json:"duplicated"`
	// Filter by resolution
	Resolution *ResolutionCriterionInput `json:"resolution"`
	// Filter by orientation
	Orientation *OrientationCriterionInput `json:"orientation"`
	// Filter by framerate
	Framerate *IntCriterionInput `json:"framerate"`
	// Filter by bitrate
	Bitrate *IntCriterionInput `json:"bitrate"`
	// Filter by video codec
	VideoCodec *StringCriterionInput `json:"video_codec"`
	// Filter by audio codec
	AudioCodec *StringCriterionInput `json:"audio_codec"`
	// Filter by duration (in seconds)
	Duration *IntCriterionInput `json:"duration"`
	// Filter to only include scenes which have markers. `true` or `false`
	HasMarkers *string `json:"has_markers"`
	// Filter to only include scenes missing this property
	IsMissing *string `json:"is_missing"`
	// Filter to only include scenes with this studio
	Studios *HierarchicalMultiCriterionInput `json:"studios"`
	// Filter to only include scenes with this group
	Groups *MultiCriterionInput `json:"groups"`
	// Filter to only include scenes with this movie
	Movies *MultiCriterionInput `json:"movies"`
	// Filter to only include scenes with this gallery
	Galleries *MultiCriterionInput `json:"galleries"`
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
	// Filter by StashID Endpoint
	StashIDEndpoint *StashIDCriterionInput `json:"stash_id_endpoint"`
	// Filter by url
	URL *StringCriterionInput `json:"url"`
	// Filter by interactive
	Interactive *bool `json:"interactive"`
	// Filter by InteractiveSpeed
	InteractiveSpeed *IntCriterionInput `json:"interactive_speed"`
	// Filter by captions
	Captions *StringCriterionInput `json:"captions"`
	// Filter by resume time
	ResumeTime *IntCriterionInput `json:"resume_time"`
	// Filter by play count
	PlayCount *IntCriterionInput `json:"play_count"`
	// Filter by play duration (in seconds)
	PlayDuration *IntCriterionInput `json:"play_duration"`
	// Filter by last played at
	LastPlayedAt *TimestampCriterionInput `json:"last_played_at"`
	// Filter by date
	Date *DateCriterionInput `json:"date"`
	// Filter by related galleries that meet this criteria
	GalleriesFilter *GalleryFilterType `json:"galleries_filter"`
	// Filter by related performers that meet this criteria
	PerformersFilter *PerformerFilterType `json:"performers_filter"`
	// Filter by related studios that meet this criteria
	StudiosFilter *StudioFilterType `json:"studios_filter"`
	// Filter by related tags that meet this criteria
	TagsFilter *TagFilterType `json:"tags_filter"`
	// Filter by related groups that meet this criteria
	GroupsFilter *MovieFilterType `json:"groups_filter"`
	// Filter by related movies that meet this criteria
	MoviesFilter *MovieFilterType `json:"movies_filter"`
	// Filter by related markers that meet this criteria
	MarkersFilter *SceneMarkerFilterType `json:"markers_filter"`
	// Filter by created at
	CreatedAt *TimestampCriterionInput `json:"created_at"`
	// Filter by updated at
	UpdatedAt *TimestampCriterionInput `json:"updated_at"`
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

	getter     SceneGetter
	scenes     []*Scene
	resolveErr error
}

// SceneMovieInput is used for groups and movies
type SceneMovieInput struct {
	MovieID    string `json:"movie_id"`
	SceneIndex *int   `json:"scene_index"`
}

type SceneGroupInput struct {
	GroupID    string `json:"group_id"`
	SceneIndex *int   `json:"scene_index"`
}

type SceneCreateInput struct {
	Title        *string           `json:"title"`
	Code         *string           `json:"code"`
	Details      *string           `json:"details"`
	Director     *string           `json:"director"`
	URL          *string           `json:"url"`
	Urls         []string          `json:"urls"`
	Date         *string           `json:"date"`
	Rating100    *int              `json:"rating100"`
	Organized    *bool             `json:"organized"`
	StudioID     *string           `json:"studio_id"`
	GalleryIds   []string          `json:"gallery_ids"`
	PerformerIds []string          `json:"performer_ids"`
	Movies       []SceneMovieInput `json:"movies"`
	Groups       []SceneGroupInput `json:"groups"`
	TagIds       []string          `json:"tag_ids"`
	// This should be a URL or a base64 encoded data URL
	CoverImage *string   `json:"cover_image"`
	StashIds   []StashID `json:"stash_ids"`
	// The first id will be assigned as primary.
	// Files will be reassigned from existing scenes if applicable.
	// Files must not already be primary for another scene.
	FileIds []string `json:"file_ids"`
}

type SceneUpdateInput struct {
	ClientMutationID *string           `json:"clientMutationId"`
	ID               string            `json:"id"`
	Title            *string           `json:"title"`
	Code             *string           `json:"code"`
	Details          *string           `json:"details"`
	Director         *string           `json:"director"`
	URL              *string           `json:"url"`
	Urls             []string          `json:"urls"`
	Date             *string           `json:"date"`
	Rating100        *int              `json:"rating100"`
	OCounter         *int              `json:"o_counter"`
	Organized        *bool             `json:"organized"`
	StudioID         *string           `json:"studio_id"`
	GalleryIds       []string          `json:"gallery_ids"`
	PerformerIds     []string          `json:"performer_ids"`
	Movies           []SceneMovieInput `json:"movies"`
	Groups           []SceneGroupInput `json:"groups"`
	TagIds           []string          `json:"tag_ids"`
	// This should be a URL or a base64 encoded data URL
	CoverImage    *string   `json:"cover_image"`
	StashIds      []StashID `json:"stash_ids"`
	ResumeTime    *float64  `json:"resume_time"`
	PlayDuration  *float64  `json:"play_duration"`
	PlayCount     *int      `json:"play_count"`
	PrimaryFileID *string   `json:"primary_file_id"`
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

func NewSceneQueryResult(getter SceneGetter) *SceneQueryResult {
	return &SceneQueryResult{
		getter: getter,
	}
}

func (r *SceneQueryResult) Resolve(ctx context.Context) ([]*Scene, error) {
	// cache results
	if r.scenes == nil && r.resolveErr == nil {
		r.scenes, r.resolveErr = r.getter.FindMany(ctx, r.IDs)
	}
	return r.scenes, r.resolveErr
}
