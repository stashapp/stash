package models

type MovieFilterType struct {
	Name     *StringCriterionInput `json:"name"`
	Director *StringCriterionInput `json:"director"`
	Synopsis *StringCriterionInput `json:"synopsis"`
	// Filter by duration (in seconds)
	Duration *IntCriterionInput `json:"duration"`
	// Filter by rating
	Rating *IntCriterionInput `json:"rating"`
	// Filter to only include movies with this studio
	Studios *HierarchicalMultiCriterionInput `json:"studios"`
	// Filter to only include movies missing this property
	IsMissing *string `json:"is_missing"`
	// Filter by url
	URL *StringCriterionInput `json:"url"`
	// Filter to only include movies where performer appears in a scene
	Performers *MultiCriterionInput `json:"performers"`
}

type MovieReader interface {
	Find(id int) (*Movie, error)
	FindMany(ids []int) ([]*Movie, error)
	// FindBySceneID(sceneID int) ([]*Movie, error)
	FindByName(name string, nocase bool) (*Movie, error)
	FindByNames(names []string, nocase bool) ([]*Movie, error)
	All() ([]*Movie, error)
	Count() (int, error)
	Query(movieFilter *MovieFilterType, findFilter *FindFilterType) ([]*Movie, int, error)
	GetFrontImage(movieID int) ([]byte, error)
	GetBackImage(movieID int) ([]byte, error)
	FindByPerformerID(performerID int) ([]*Movie, error)
	CountByPerformerID(performerID int) (int, error)
	FindByStudioID(studioID int) ([]*Movie, error)
	CountByStudioID(studioID int) (int, error)
}

type MovieWriter interface {
	Create(newMovie Movie) (*Movie, error)
	Update(updatedMovie MoviePartial) (*Movie, error)
	UpdateFull(updatedMovie Movie) (*Movie, error)
	Destroy(id int) error
	UpdateImages(movieID int, frontImage []byte, backImage []byte) error
	DestroyImages(movieID int) error
}

type MovieReaderWriter interface {
	MovieReader
	MovieWriter
}
