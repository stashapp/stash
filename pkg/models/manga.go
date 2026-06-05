package models

type MangaFilterType struct {
	OperatorFilter[MangaFilterType]
	ID      *IntCriterionInput    `json:"id"`
	Title   *StringCriterionInput `json:"title"`
	Details *StringCriterionInput `json:"details"`
	// Filter by rating expressed as 1-100
	Rating100 *IntCriterionInput `json:"rating100"`
	// Filter by organized
	Organized *bool `json:"organized"`
	// Filter to only include mangas with this studio
	Studios *HierarchicalMultiCriterionInput `json:"studios"`
	// Filter to only include mangas with these tags
	Tags *HierarchicalMultiCriterionInput `json:"tags"`
	// Filter by tag count
	TagCount *IntCriterionInput `json:"tag_count"`
	// Filter to only include mangas with performers with these tags
	PerformerTags *HierarchicalMultiCriterionInput `json:"performer_tags"`
	// Filter to only include mangas with these performers
	Performers *MultiCriterionInput `json:"performers"`
	// Filter by performer count
	PerformerCount *IntCriterionInput `json:"performer_count"`
	// Filter mangas that have performers that have been favorited
	PerformerFavorite *bool `json:"performer_favorite"`
	// Filter mangas by performer age at time of manga
	PerformerAge *IntCriterionInput `json:"performer_age"`
	// Filter by url
	URL *StringCriterionInput `json:"url"`
	// Filter by date
	Date *DateCriterionInput `json:"date"`
	// Filter by related performers that meet this criteria
	PerformersFilter *PerformerFilterType `json:"performers_filter"`
	// Filter by related studios that meet this criteria
	StudiosFilter *StudioFilterType `json:"studios_filter"`
	// Filter by related tags that meet this criteria
	TagsFilter *TagFilterType `json:"tags_filter"`
	// Filter by created at
	CreatedAt *TimestampCriterionInput `json:"created_at"`
	// Filter by updated at
	UpdatedAt *TimestampCriterionInput `json:"updated_at"`

	// Filter by custom fields
	CustomFields []CustomFieldCriterionInput `json:"custom_fields"`
}

type MangaUpdateInput struct {
	ClientMutationID *string  `json:"clientMutationId"`
	ID               string   `json:"id"`
	Title            *string  `json:"title"`
	URL              *string  `json:"url"`
	Date             *string  `json:"date"`
	Details          *string  `json:"details"`
	Rating100        *int     `json:"rating100"`
	Organized        *bool    `json:"organized"`
	StudioID         *string  `json:"studio_id"`
	TagIds           []string `json:"tag_ids"`
	PerformerIds     []string `json:"performer_ids"`
}

type MangaDestroyInput struct {
	Ids []string `json:"ids"`
}
