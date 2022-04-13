package models

import "context"

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
	Find(ctx context.Context, id int) (*Tag, error)
	FindMany(ctx context.Context, ids []int) ([]*Tag, error)
	FindBySceneID(ctx context.Context, sceneID int) ([]*Tag, error)
	FindByPerformerID(ctx context.Context, performerID int) ([]*Tag, error)
	FindBySceneMarkerID(ctx context.Context, sceneMarkerID int) ([]*Tag, error)
	FindByImageID(ctx context.Context, imageID int) ([]*Tag, error)
	FindByGalleryID(ctx context.Context, galleryID int) ([]*Tag, error)
	FindByName(ctx context.Context, name string, nocase bool) (*Tag, error)
	FindByNames(ctx context.Context, names []string, nocase bool) ([]*Tag, error)
	FindByParentTagID(ctx context.Context, parentID int) ([]*Tag, error)
	FindByChildTagID(ctx context.Context, childID int) ([]*Tag, error)
	Count(ctx context.Context) (int, error)
	All(ctx context.Context) ([]*Tag, error)
	// TODO - this interface is temporary until the filter schema can fully
	// support the query needed
	QueryForAutoTag(ctx context.Context, words []string) ([]*Tag, error)
	Query(ctx context.Context, tagFilter *TagFilterType, findFilter *FindFilterType) ([]*Tag, int, error)
	GetImage(ctx context.Context, tagID int) ([]byte, error)
	GetAliases(ctx context.Context, tagID int) ([]string, error)
	FindAllAncestors(ctx context.Context, tagID int, excludeIDs []int) ([]*TagPath, error)
	FindAllDescendants(ctx context.Context, tagID int, excludeIDs []int) ([]*TagPath, error)
}

type TagWriter interface {
	Create(ctx context.Context, newTag Tag) (*Tag, error)
	Update(ctx context.Context, updateTag TagPartial) (*Tag, error)
	UpdateFull(ctx context.Context, updatedTag Tag) (*Tag, error)
	Destroy(ctx context.Context, id int) error
	UpdateImage(ctx context.Context, tagID int, image []byte) error
	DestroyImage(ctx context.Context, tagID int) error
	UpdateAliases(ctx context.Context, tagID int, aliases []string) error
	Merge(ctx context.Context, source []int, destination int) error
	UpdateParentTags(ctx context.Context, tagID int, parentIDs []int) error
	UpdateChildTags(ctx context.Context, tagID int, parentIDs []int) error
}

type TagReaderWriter interface {
	TagReader
	TagWriter
}
