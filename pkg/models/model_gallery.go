package models

import (
	"path/filepath"
	"time"
)

type Gallery struct {
	ID          int        `json:"id"`
	Path        *string    `json:"path"`
	Checksum    string     `json:"checksum"`
	Zip         bool       `json:"zip"`
	Title       string     `json:"title"`
	URL         string     `json:"url"`
	Date        *Date      `json:"date"`
	Details     string     `json:"details"`
	Rating      *int       `json:"rating"`
	Organized   bool       `json:"organized"`
	StudioID    *int       `json:"studio_id"`
	FileModTime *time.Time `json:"file_mod_time"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	SceneIDs     []int `json:"scene_ids"`
	TagIDs       []int `json:"tag_ids"`
	PerformerIDs []int `json:"performer_ids"`
}

// GalleryPartial represents part of a Gallery object. It is used to update
// the database entry. Only non-nil fields will be updated.
type GalleryPartial struct {
	Path        OptionalString
	Checksum    OptionalString
	Zip         OptionalBool
	Title       OptionalString
	URL         OptionalString
	Date        OptionalDate
	Details     OptionalString
	Rating      OptionalInt
	Organized   OptionalBool
	StudioID    OptionalInt
	FileModTime OptionalTime
	CreatedAt   OptionalTime
	UpdatedAt   OptionalTime

	SceneIDs     *UpdateIDs
	TagIDs       *UpdateIDs
	PerformerIDs *UpdateIDs
}

func NewGalleryPartial() GalleryPartial {
	updatedTime := time.Now()
	return GalleryPartial{
		UpdatedAt: NewOptionalTime(updatedTime),
	}
}

func (s *Gallery) File() File {
	var path string
	if s.Path != nil {
		path = *s.Path
	}
	ret := File{
		Path: path,
	}

	ret.Checksum = s.Checksum

	if s.FileModTime != nil {
		ret.FileModTime = *s.FileModTime
	}

	return ret
}

func (s *Gallery) SetFile(f File) {
	path := f.Path
	s.Path = &path

	if f.Checksum != "" {
		s.Checksum = f.Checksum
	}

	zeroTime := time.Time{}
	if f.FileModTime != zeroTime {
		s.FileModTime = &f.FileModTime
	}
}

// GetTitle returns the title of the scene. If the Title field is empty,
// then the base filename is returned.
func (s Gallery) GetTitle() string {
	if s.Title != "" {
		return s.Title
	}

	if s.Path != nil {
		return filepath.Base(*s.Path)
	}

	return ""
}

const DefaultGthumbWidth int = 640

type Galleries []*Gallery

func (g *Galleries) Append(o interface{}) {
	*g = append(*g, o.(*Gallery))
}

func (g *Galleries) New() interface{} {
	return &Gallery{}
}
