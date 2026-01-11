package models

type StudioFilterType struct {
	OperatorFilter[StudioFilterType]
	Name    *StringCriterionInput `json:"name"`
	Details *StringCriterionInput `json:"details"`
	// Filter to only include studios with this parent studio
	Parents *MultiCriterionInput `json:"parents"`
	// Filter by StashID
	StashID *StringCriterionInput `json:"stash_id"`
	// Filter by StashID Endpoint
	StashIDEndpoint *StashIDCriterionInput `json:"stash_id_endpoint"`
	// Filter by StashIDs Endpoint
	StashIDsEndpoint *StashIDsCriterionInput `json:"stash_ids_endpoint"`
	// Filter to only include studios missing this property
	IsMissing *string `json:"is_missing"`
	// Filter by rating expressed as 1-100
	Rating100 *IntCriterionInput `json:"rating100"`
	// Filter to only include studios with these tags
	Tags *HierarchicalMultiCriterionInput `json:"tags"`
	// Filter by tag count
	TagCount *IntCriterionInput `json:"tag_count"`
	// Filter by favorite
	Favorite *bool `json:"favorite"`
	// Filter by scene count
	SceneCount *IntCriterionInput `json:"scene_count"`
	// Filter by image count
	ImageCount *IntCriterionInput `json:"image_count"`
	// Filter by gallery count
	GalleryCount *IntCriterionInput `json:"gallery_count"`
	// Filter by url
	URL *StringCriterionInput `json:"url"`
	// Filter by studio aliases
	Aliases *StringCriterionInput `json:"aliases"`
	// Filter by subsidiary studio count
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

type StudioCreateInput struct {
	Name     string   `json:"name"`
	URL      *string  `json:"url"` // deprecated
	Urls     []string `json:"urls"`
	ParentID *string  `json:"parent_id"`
	// This should be a URL or a base64 encoded data URL
	Image         *string        `json:"image"`
	StashIds      []StashIDInput `json:"stash_ids"`
	Rating100     *int           `json:"rating100"`
	Favorite      *bool          `json:"favorite"`
	Details       *string        `json:"details"`
	Aliases       []string       `json:"aliases"`
	TagIds        []string       `json:"tag_ids"`
	IgnoreAutoTag *bool          `json:"ignore_auto_tag"`
}

type StudioUpdateInput struct {
	ID       string   `json:"id"`
	Name     *string  `json:"name"`
	URL      *string  `json:"url"` // deprecated
	Urls     []string `json:"urls"`
	ParentID *string  `json:"parent_id"`
	// This should be a URL or a base64 encoded data URL
	Image         *string        `json:"image"`
	StashIds      []StashIDInput `json:"stash_ids"`
	Rating100     *int           `json:"rating100"`
	Favorite      *bool          `json:"favorite"`
	Details       *string        `json:"details"`
	Aliases       []string       `json:"aliases"`
	TagIds        []string       `json:"tag_ids"`
	IgnoreAutoTag *bool          `json:"ignore_auto_tag"`
}
