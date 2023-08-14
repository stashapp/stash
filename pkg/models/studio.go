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
	// Filter by StashID Endpoint
	StashIDEndpoint *StashIDCriterionInput `json:"stash_id_endpoint"`
	// Filter to only include studios missing this property
	IsMissing *string `json:"is_missing"`
	// Filter by rating expressed as 1-5
	Rating *IntCriterionInput `json:"rating"`
	// Filter by rating expressed as 1-100
	Rating100 *IntCriterionInput `json:"rating100"`
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
	// Filter by created at
	CreatedAt *TimestampCriterionInput `json:"created_at"`
	// Filter by updated at
	UpdatedAt *TimestampCriterionInput `json:"updated_at"`
}

type StudioFinder interface {
	FindMany(ctx context.Context, ids []int) ([]*Studio, error)
}

type StudioReader interface {
	Find(ctx context.Context, id int) (*Studio, error)
	StudioFinder
	FindChildren(ctx context.Context, id int) ([]*Studio, error)
	FindByName(ctx context.Context, name string, nocase bool) (*Studio, error)
	FindByStashID(ctx context.Context, stashID StashID) ([]*Studio, error)
	FindByStashIDStatus(ctx context.Context, hasStashID bool, stashboxEndpoint string) ([]*Studio, error)
	Count(ctx context.Context) (int, error)
	All(ctx context.Context) ([]*Studio, error)
	// TODO - this interface is temporary until the filter schema can fully
	// support the query needed
	QueryForAutoTag(ctx context.Context, words []string) ([]*Studio, error)
	Query(ctx context.Context, studioFilter *StudioFilterType, findFilter *FindFilterType) ([]*Studio, int, error)
	GetImage(ctx context.Context, studioID int) ([]byte, error)
	HasImage(ctx context.Context, studioID int) (bool, error)
	AliasLoader
	StashIDLoader
}

type StudioWriter interface {
	Create(ctx context.Context, newStudio *Studio) error
	UpdatePartial(ctx context.Context, input StudioPartial) (*Studio, error)
	Update(ctx context.Context, updatedStudio *Studio) error
	Destroy(ctx context.Context, id int) error
	UpdateImage(ctx context.Context, studioID int, image []byte) error
}

type StudioReaderWriter interface {
	StudioReader
	StudioWriter
}
