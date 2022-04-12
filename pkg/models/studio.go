package models

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
	Find(id int) (*Studio, error)
	FindMany(ids []int) ([]*Studio, error)
	FindChildren(id int) ([]*Studio, error)
	FindByName(name string, nocase bool) (*Studio, error)
	FindByStashID(stashID StashID) ([]*Studio, error)
	Count() (int, error)
	All() ([]*Studio, error)
	// TODO - this interface is temporary until the filter schema can fully
	// support the query needed
	QueryForAutoTag(words []string) ([]*Studio, error)
	Query(studioFilter *StudioFilterType, findFilter *FindFilterType) ([]*Studio, int, error)
	GetImage(studioID int) ([]byte, error)
	HasImage(studioID int) (bool, error)
	GetStashIDs(studioID int) ([]*StashID, error)
	GetAliases(studioID int) ([]string, error)
}

type StudioWriter interface {
	Create(newStudio Studio) (*Studio, error)
	Update(updatedStudio StudioPartial) (*Studio, error)
	UpdateFull(updatedStudio Studio) (*Studio, error)
	Destroy(id int) error
	UpdateImage(studioID int, image []byte) error
	DestroyImage(studioID int) error
	UpdateStashIDs(studioID int, stashIDs []StashID) error
	UpdateAliases(studioID int, aliases []string) error
}

type StudioReaderWriter interface {
	StudioReader
	StudioWriter
}
