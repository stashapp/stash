package models

type GroupFilterType struct {
	OperatorFilter[GroupFilterType]
	Name     *StringCriterionInput `json:"name"`
	Director *StringCriterionInput `json:"director"`
	Synopsis *StringCriterionInput `json:"synopsis"`
	// Filter by duration (in seconds)
	Duration *IntCriterionInput `json:"duration"`
	// Filter by rating expressed as 1-100
	Rating100 *IntCriterionInput `json:"rating100"`
	// Filter to only include movies with this studio
	Studios *HierarchicalMultiCriterionInput `json:"studios"`
	// Filter to only include movies missing this property
	IsMissing *string `json:"is_missing"`
	// Filter by url
	URL *StringCriterionInput `json:"url"`
	// Filter to only include movies where performer appears in a scene
	Performers *MultiCriterionInput `json:"performers"`
	// Filter to only include performers with these tags
	Tags *HierarchicalMultiCriterionInput `json:"tags"`
	// Filter by tag count
	TagCount *IntCriterionInput `json:"tag_count"`
	// Filter by date
	Date *DateCriterionInput `json:"date"`
	// Filter by containing groups
	ContainingGroups *HierarchicalMultiCriterionInput `json:"containing_groups"`
	// Filter by sub groups
	SubGroups *HierarchicalMultiCriterionInput `json:"sub_groups"`
	// Filter by number of containing groups the group has
	ContainingGroupCount *IntCriterionInput `json:"containing_group_count"`
	// Filter by number of sub-groups the group has
	SubGroupCount *IntCriterionInput `json:"sub_group_count"`
	// Filter by related scenes that meet this criteria
	ScenesFilter *SceneFilterType `json:"scenes_filter"`
	// Filter by related studios that meet this criteria
	StudiosFilter *StudioFilterType `json:"studios_filter"`
	// Filter by created at
	CreatedAt *TimestampCriterionInput `json:"created_at"`
	// Filter by updated at
	UpdatedAt *TimestampCriterionInput `json:"updated_at"`
}
