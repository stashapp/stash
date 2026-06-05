package models

import "context"

// MangaGetter provides methods to get mangas by ID.
type MangaGetter interface {
	FindMany(ctx context.Context, ids []int) ([]*Manga, error)
	Find(ctx context.Context, id int) (*Manga, error)
}

// MangaFinder provides methods to find mangas.
type MangaFinder interface {
	MangaGetter
}

// MangaQueryer provides methods to query mangas.
type MangaQueryer interface {
	Query(ctx context.Context, mangaFilter *MangaFilterType, findFilter *FindFilterType) ([]*Manga, int, error)
	QueryCount(ctx context.Context, mangaFilter *MangaFilterType, findFilter *FindFilterType) (int, error)
}

// MangaCounter provides methods to count mangas.
type MangaCounter interface {
	Count(ctx context.Context) (int, error)
}

// MangaCreator provides methods to create mangas.
type MangaCreator interface {
	Create(ctx context.Context, newManga *CreateMangaInput) error
}

// MangaUpdater provides methods to update mangas.
type MangaUpdater interface {
	Update(ctx context.Context, updatedManga *UpdateMangaInput) error
	UpdatePartial(ctx context.Context, id int, updatedManga MangaPartial) (*Manga, error)
}

// MangaDestroyer provides methods to destroy mangas.
type MangaDestroyer interface {
	Destroy(ctx context.Context, id int) error
}

type MangaCreatorUpdater interface {
	MangaCreator
	MangaUpdater
}

// MangaReader provides all methods to read mangas.
type MangaReader interface {
	MangaFinder
	MangaQueryer
	MangaCounter

	URLLoader
	PerformerIDLoader
	TagIDLoader

	All(ctx context.Context) ([]*Manga, error)
}

// MangaWriter provides all methods to modify mangas.
type MangaWriter interface {
	MangaCreator
	MangaUpdater
	MangaDestroyer
}

// MangaReaderWriter provides all manga methods.
type MangaReaderWriter interface {
	MangaReader
	MangaWriter
}
