package models

type TagFilterType struct {
	OperatorFilter[TagFilterType]
	// Filter by tag name
	Name *StringCriterionInput `json:"name"`
	// Filter by tag aliases
	Aliases *StringCriterionInput `json:"aliases"`
	// Filter by tag favorites
	Favorite *bool `json:"favorite"`
	// Filter by tag description
	Description *StringCriterionInput `json:"description"`
	// Filter to only include tags missing this property
	IsMissing *string `json:"is_missing"`
	// Filter by number of scenes with this tag
	SceneCount *IntCriterionInput `json:"scene_count"`
	// Filter by number of images with this tag
	ImageCount *IntCriterionInput `json:"image_count"`
	// Filter by number of galleries with this tag
	GalleryCount *IntCriterionInput `json:"gallery_count"`
	// Filter by number of performers with this tag
	PerformerCount *IntCriterionInput `json:"performer_count"`
	// Filter by number of studios with this tag
	StudioCount *IntCriterionInput `json:"studio_count"`
	// Filter by number of movies with this tag
	MovieCount *IntCriterionInput `json:"movie_count"`
	// Filter by number of markers with this tag
	MarkerCount *IntCriterionInput `json:"marker_count"`
	// Filter by parent tags
	Parents *HierarchicalMultiCriterionInput `json:"parents"`
	// Filter by child tags
	Children *HierarchicalMultiCriterionInput `json:"children"`
	// Filter by number of parent tags the tag has
	ParentCount *IntCriterionInput `json:"parent_count"`
	// Filter by number f child tags the tag has
	ChildCount *IntCriterionInput `json:"child_count"`
	// Filter by autotag ignore value
	IgnoreAutoTag *bool `json:"ignore_auto_tag"`
	// Filter by related scenes that meet this criteria
	ScenesFilter *SceneFilterType `json:"scenes_filter"`
	// Filter by related images that meet this criteria
	ImagesFilter *ImageFilterType `json:"images_filter"`
	// Filter by related galleries that meet this criteria
	GalleriesFilter *GalleryFilterType `json:"galleries_filter"`
	// Filter by created at
	CreatedAt *TimestampCriterionInput `json:"created_at"`
	// Filter by updated at
	UpdatedAt *TimestampCriterionInput `json:"updated_at"`
}
