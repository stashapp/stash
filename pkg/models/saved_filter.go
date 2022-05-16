package models

type SavedFilterReader interface {
	Find(id int) (*SavedFilter, error)
	FindByMode(mode FilterMode) ([]*SavedFilter, error)
	FindDefault(mode FilterMode) (*SavedFilter, error)
	FindRecommended() ([]*SavedFilter, error)
}

type SavedFilterWriter interface {
	Create(obj SavedFilter) (*SavedFilter, error)
	Update(obj SavedFilter) (*SavedFilter, error)
	SetDefault(obj SavedFilter) (*SavedFilter, error)
	Destroy(id int) error
}

type SavedFilterReaderWriter interface {
	SavedFilterReader
	SavedFilterWriter
}
