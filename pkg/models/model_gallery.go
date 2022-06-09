package models

import (
	"path/filepath"
	"time"

	"github.com/stashapp/stash/pkg/file"
)

type Gallery struct {
	ID int `json:"id"`

	// Path        *string    `json:"path"`
	// Checksum    string     `json:"checksum"`
	// Zip         bool       `json:"zip"`

	Title     string `json:"title"`
	URL       string `json:"url"`
	Date      *Date  `json:"date"`
	Details   string `json:"details"`
	Rating    *int   `json:"rating"`
	Organized bool   `json:"organized"`
	StudioID  *int   `json:"studio_id"`

	// FileModTime *time.Time `json:"file_mod_time"`

	// transient - not persisted
	Files []file.File

	FolderID *file.FolderID `json:"folder_id"`

	// transient - not persisted
	FolderPath string `json:"folder_path"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	SceneIDs     []int `json:"scene_ids"`
	TagIDs       []int `json:"tag_ids"`
	PerformerIDs []int `json:"performer_ids"`
}

func (g Gallery) PrimaryFile() file.File {
	if len(g.Files) == 0 {
		return nil
	}

	return g.Files[0]
}

func (g Gallery) Path() string {
	if p := g.PrimaryFile(); p != nil {
		return p.Base().Path
	}

	return g.FolderPath
}

func (g Gallery) Checksum() string {
	if p := g.PrimaryFile(); p != nil {
		v := p.Base().Fingerprints.Get(file.FingerprintTypeMD5)
		if v == nil {
			return ""
		}

		return v.(string)
	}
	return ""
}

// GalleryPartial represents part of a Gallery object. It is used to update
// the database entry. Only non-nil fields will be updated.
type GalleryPartial struct {
	// Path        OptionalString
	// Checksum    OptionalString
	// Zip         OptionalBool
	Title     OptionalString
	URL       OptionalString
	Date      OptionalDate
	Details   OptionalString
	Rating    OptionalInt
	Organized OptionalBool
	StudioID  OptionalInt
	// FileModTime OptionalTime
	CreatedAt OptionalTime
	UpdatedAt OptionalTime

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

// GetTitle returns the title of the scene. If the Title field is empty,
// then the base filename is returned.
func (g Gallery) GetTitle() string {
	if g.Title != "" {
		return g.Title
	}

	if len(g.Files) > 0 {
		return filepath.Base(g.Path())
	}

	if g.FolderPath != "" {
		return g.FolderPath
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
