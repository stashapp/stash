package models

type StudioReader interface {
	Find(id int) (*Studio, error)
	FindMany(ids []int) ([]*Studio, error)
	FindChildren(id int) ([]*Studio, error)
	FindByName(name string, nocase bool) (*Studio, error)
	// Count() (int, error)
	All() ([]*Studio, error)
	// AllSlim() ([]*Studio, error)
	// Query(studioFilter *StudioFilterType, findFilter *FindFilterType) ([]*Studio, int)
	GetStudioImage(studioID int) ([]byte, error)
	HasStudioImage(studioID int) (bool, error)
}

type StudioWriter interface {
	Create(newStudio Studio) (*Studio, error)
	Update(updatedStudio StudioPartial) (*Studio, error)
	UpdateFull(updatedStudio Studio) (*Studio, error)
	// Destroy(id string) error
	UpdateStudioImage(studioID int, image []byte) error
	// DestroyStudioImage(studioID int) error
}

type StudioReaderWriter interface {
	StudioReader
	StudioWriter
}
