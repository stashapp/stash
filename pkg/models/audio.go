// TODO(audio): update this file

package models

import "context"

type AudioFilterType struct {
	OperatorFilter[AudioFilterType]
	ID      *IntCriterionInput    `json:"id"`
	Title   *StringCriterionInput `json:"title"`
	Code    *StringCriterionInput `json:"code"`
	Details *StringCriterionInput `json:"details"`
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
	// Filter Audios by duplication criteria
	Duplicated *DuplicationCriterionInput `json:"duplicated"`
	// Filter by resolution
	Resolution *ResolutionCriterionInput `json:"resolution"`
	// Filter by orientation
	Orientation *OrientationCriterionInput `json:"orientation"`
	// Filter by samplerate
	Samplerate *IntCriterionInput `json:"samplerate"`
	// Filter by bitrate
	Bitrate *IntCriterionInput `json:"bitrate"`
	// Filter by audio codec
	AudioCodec *StringCriterionInput `json:"audio_codec"`
	// Filter by duration (in seconds)
	Duration *IntCriterionInput `json:"duration"`
	// Filter to only include audios which have markers. `true` or `false`
	HasMarkers *string `json:"has_markers"`
	// Filter to only include audios missing this property
	IsMissing *string `json:"is_missing"`
	// Filter to only include audios with this studio
	Studios *HierarchicalMultiCriterionInput `json:"studios"`
	// Filter to only include audios with this group
	Groups *HierarchicalMultiCriterionInput `json:"groups"`
	// Filter to only include audios with this movie
	Movies *MultiCriterionInput `json:"movies"`
	// Filter to only include audios with this gallery
	Galleries *MultiCriterionInput `json:"galleries"`
	// Filter to only include audios with these tags
	Tags *HierarchicalMultiCriterionInput `json:"tags"`
	// Filter by tag count
	TagCount *IntCriterionInput `json:"tag_count"`
	// Filter to only include audios with performers with these tags
	PerformerTags *HierarchicalMultiCriterionInput `json:"performer_tags"`
	// Filter audios that have performers that have been favorited
	PerformerFavorite *bool `json:"performer_favorite"`
	// Filter audios by performer age at time of audio
	PerformerAge *IntCriterionInput `json:"performer_age"`
	// Filter to only include audios with these performers
	Performers *MultiCriterionInput `json:"performers"`
	// Filter by performer count
	PerformerCount *IntCriterionInput `json:"performer_count"`
	// Filter by url
	URL *StringCriterionInput `json:"url"`
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
	GroupsFilter *GroupFilterType `json:"groups_filter"`
	// Filter by related files that meet this criteria
	FilesFilter *FileFilterType `json:"files_filter"`
	// Filter by created at
	CreatedAt *TimestampCriterionInput `json:"created_at"`
	// Filter by updated at
	UpdatedAt *TimestampCriterionInput `json:"updated_at"`

	// Filter by custom fields
	CustomFields []CustomFieldCriterionInput `json:"custom_fields"`
}

type AudioQueryOptions struct {
	QueryOptions
	AudioFilter *AudioFilterType

	TotalDuration bool
	TotalSize     bool
}

type AudioQueryResult struct {
	QueryResult[int]
	TotalDuration float64
	TotalSize     float64

	getter     AudioGetter
	audios     []*Audio
	resolveErr error
}

type AudioGroupInput struct {
	GroupID    string `json:"group_id"`
	AudioIndex *int   `json:"audio_index"`
}

type AudioCreateInput struct {
	Title        *string           `json:"title"`
	Code         *string           `json:"code"`
	Details      *string           `json:"details"`
	URL          *string           `json:"url"`
	Urls         []string          `json:"urls"`
	Date         *string           `json:"date"`
	Rating100    *int              `json:"rating100"`
	Organized    *bool             `json:"organized"`
	StudioID     *string           `json:"studio_id"`
	GalleryIds   []string          `json:"gallery_ids"`
	PerformerIds []string          `json:"performer_ids"`
	Groups       []AudioGroupInput `json:"groups"`
	TagIds       []string          `json:"tag_ids"`
	// This should be a URL or a base64 encoded data URL
	CoverImage *string `json:"cover_image"`
	// The first id will be assigned as primary.
	// Files will be reassigned from existing audios if applicable.
	// Files must not already be primary for another audio.
	FileIds      []string       `json:"file_ids"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

type AudioUpdateInput struct {
	ClientMutationID *string           `json:"clientMutationId"`
	ID               string            `json:"id"`
	Title            *string           `json:"title"`
	Code             *string           `json:"code"`
	Details          *string           `json:"details"`
	URL              *string           `json:"url"`
	Urls             []string          `json:"urls"`
	Date             *string           `json:"date"`
	Rating100        *int              `json:"rating100"`
	OCounter         *int              `json:"o_counter"`
	Organized        *bool             `json:"organized"`
	StudioID         *string           `json:"studio_id"`
	GalleryIds       []string          `json:"gallery_ids"`
	PerformerIds     []string          `json:"performer_ids"`
	Groups           []AudioGroupInput `json:"groups"`
	TagIds           []string          `json:"tag_ids"`
	// This should be a URL or a base64 encoded data URL
	CoverImage    *string  `json:"cover_image"`
	ResumeTime    *float64 `json:"resume_time"`
	PlayDuration  *float64 `json:"play_duration"`
	PlayCount     *int     `json:"play_count"`
	PrimaryFileID *string  `json:"primary_file_id"`
	CustomFields  *CustomFieldsInput
}

type AudioDestroyInput struct {
	ID               string `json:"id"`
	DeleteFile       *bool  `json:"delete_file"`
	DeleteGenerated  *bool  `json:"delete_generated"`
	DestroyFileEntry *bool  `json:"destroy_file_entry"`
}

type AudiosDestroyInput struct {
	Ids              []string `json:"ids"`
	DeleteFile       *bool    `json:"delete_file"`
	DeleteGenerated  *bool    `json:"delete_generated"`
	DestroyFileEntry *bool    `json:"destroy_file_entry"`
}

func NewAudioQueryResult(getter AudioGetter) *AudioQueryResult {
	return &AudioQueryResult{
		getter: getter,
	}
}

func (r *AudioQueryResult) Resolve(ctx context.Context) ([]*Audio, error) {
	// cache results
	if r.audios == nil && r.resolveErr == nil {
		r.audios, r.resolveErr = r.getter.FindMany(ctx, r.IDs)
	}
	return r.audios, r.resolveErr
}
