package models

import "context"

type StudioFilterType struct {
	And     *StudioFilterType     `json:"AND"`
	Or      *StudioFilterType     `json:"OR"`
	Not     *StudioFilterType     `json:"NOT"`
	Name    *StringCriterionInput `json:"name"`
	Details *StringCriterionInput `json:"details"`
	// Filter to only include studios with this parent studio
	Parents *MultiCriterionInput `json:"parents"`
	// Filter by StashID
	StashID *StringCriterionInput `json:"stash_id"`
	// Filter to only include studios missing this property
	IsMissing *string `json:"is_missing"`
	// Filter by rating
	Rating *IntCriterionInput `json:"rating"`
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
	// Filter by autotag ignore value
	IgnoreAutoTag *bool `json:"ignore_auto_tag"`
}

type StudioReader interface {
	Find(ctx context.Context, id int) (*Studio, error)
	FindMany(ctx context.Context, ids []int) ([]*Studio, error)
	FindChildren(ctx context.Context, id int) ([]*Studio, error)
	FindByName(ctx context.Context, name string, nocase bool) (*Studio, error)
	FindByStashID(ctx context.Context, stashID StashID) ([]*Studio, error)
	Count(ctx context.Context) (int, error)
	All(ctx context.Context) ([]*Studio, error)
	// TODO - this interface is temporary until the filter schema can fully
	// support the query needed
	QueryForAutoTag(ctx context.Context, words []string) ([]*Studio, error)
	Query(ctx context.Context, studioFilter *StudioFilterType, findFilter *FindFilterType) ([]*Studio, int, error)
	GetImage(ctx context.Context, studioID int) ([]byte, error)
	HasImage(ctx context.Context, studioID int) (bool, error)
	GetStashIDs(ctx context.Context, studioID int) ([]*StashID, error)
	GetAliases(ctx context.Context, studioID int) ([]string, error)
}

type StudioWriter interface {
	Create(ctx context.Context, newStudio Studio) (*Studio, error)
	Update(ctx context.Context, updatedStudio StudioPartial) (*Studio, error)
	UpdateFull(ctx context.Context, updatedStudio Studio) (*Studio, error)
	Destroy(ctx context.Context, id int) error
	UpdateImage(ctx context.Context, studioID int, image []byte) error
	DestroyImage(ctx context.Context, studioID int) error
	UpdateStashIDs(ctx context.Context, studioID int, stashIDs []*StashID) error
	UpdateAliases(ctx context.Context, studioID int, aliases []string) error
}

type StudioReaderWriter interface {
	StudioReader
	StudioWriter
}
