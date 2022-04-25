package models

type TagFilterType struct {
	And *TagFilterType `json:"AND"`
	Or  *TagFilterType `json:"OR"`
	Not *TagFilterType `json:"NOT"`
	// Filter by tag name
	Name *StringCriterionInput `json:"name"`
	// Filter by tag aliases
	Aliases *StringCriterionInput `json:"aliases"`
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
}

type TagReader interface {
	Find(id int) (*Tag, error)
	FindMany(ids []int) ([]*Tag, error)
	FindBySceneID(sceneID int) ([]*Tag, error)
	FindByPerformerID(performerID int) ([]*Tag, error)
	FindBySceneMarkerID(sceneMarkerID int) ([]*Tag, error)
	FindByImageID(imageID int) ([]*Tag, error)
	FindByGalleryID(galleryID int) ([]*Tag, error)
	FindByName(name string, nocase bool) (*Tag, error)
	FindByNames(names []string, nocase bool) ([]*Tag, error)
	FindByParentTagID(parentID int) ([]*Tag, error)
	FindByChildTagID(childID int) ([]*Tag, error)
	Count() (int, error)
	All() ([]*Tag, error)
	// TODO - this interface is temporary until the filter schema can fully
	// support the query needed
	QueryForAutoTag(words []string) ([]*Tag, error)
	Query(tagFilter *TagFilterType, findFilter *FindFilterType) ([]*Tag, int, error)
	GetImage(tagID int) ([]byte, error)
	GetAliases(tagID int) ([]string, error)
	FindAllAncestors(tagID int, excludeIDs []int) ([]*TagPath, error)
	FindAllDescendants(tagID int, excludeIDs []int) ([]*TagPath, error)
}

type TagWriter interface {
	Create(newTag Tag) (*Tag, error)
	Update(updateTag TagPartial) (*Tag, error)
	UpdateFull(updatedTag Tag) (*Tag, error)
	Destroy(id int) error
	UpdateImage(tagID int, image []byte) error
	DestroyImage(tagID int) error
	UpdateAliases(tagID int, aliases []string) error
	Merge(source []int, destination int) error
	UpdateParentTags(tagID int, parentIDs []int) error
	UpdateChildTags(tagID int, parentIDs []int) error
}

type TagReaderWriter interface {
	TagReader
	TagWriter
}
