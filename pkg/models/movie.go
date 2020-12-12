package models

type MovieReader interface {
	Find(id int) (*Movie, error)
	FindMany(ids []int) ([]*Movie, error)
	// FindBySceneID(sceneID int) ([]*Movie, error)
	FindByName(name string, nocase bool) (*Movie, error)
	FindByNames(names []string, nocase bool) ([]*Movie, error)
	All() ([]*Movie, error)
	// AllSlim() ([]*Movie, error)
	// Query(movieFilter *MovieFilterType, findFilter *FindFilterType) ([]*Movie, int)
	GetFrontImage(movieID int) ([]byte, error)
	GetBackImage(movieID int) ([]byte, error)
}

type MovieWriter interface {
	Create(newMovie Movie) (*Movie, error)
	Update(updatedMovie MoviePartial) (*Movie, error)
	UpdateFull(updatedMovie Movie) (*Movie, error)
	// Destroy(id string) error
	UpdateMovieImages(movieID int, frontImage []byte, backImage []byte) error
	// DestroyMovieImages(movieID int) error
}

type MovieReaderWriter interface {
	MovieReader
	MovieWriter
}
