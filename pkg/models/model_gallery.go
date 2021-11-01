package models

import (
	"database/sql"
	"path/filepath"
	"time"
)

type Gallery struct {
	ID          int                 `db:"id" json:"id"`
	Path        sql.NullString      `db:"path" json:"path"`
	Checksum    string              `db:"checksum" json:"checksum"`
	Zip         bool                `db:"zip" json:"zip"`
	Title       sql.NullString      `db:"title" json:"title"`
	URL         sql.NullString      `db:"url" json:"url"`
	Date        SQLiteDate          `db:"date" json:"date"`
	Details     sql.NullString      `db:"details" json:"details"`
	Rating      sql.NullInt64       `db:"rating" json:"rating"`
	Organized   bool                `db:"organized" json:"organized"`
	StudioID    sql.NullInt64       `db:"studio_id,omitempty" json:"studio_id"`
	FileModTime NullSQLiteTimestamp `db:"file_mod_time" json:"file_mod_time"`
	CreatedAt   SQLiteTimestamp     `db:"created_at" json:"created_at"`
	UpdatedAt   SQLiteTimestamp     `db:"updated_at" json:"updated_at"`
}

// GalleryPartial represents part of a Gallery object. It is used to update
// the database entry. Only non-nil fields will be updated.
type GalleryPartial struct {
	ID          int                  `db:"id" json:"id"`
	Path        *sql.NullString      `db:"path" json:"path"`
	Checksum    *string              `db:"checksum" json:"checksum"`
	Title       *sql.NullString      `db:"title" json:"title"`
	URL         *sql.NullString      `db:"url" json:"url"`
	Date        *SQLiteDate          `db:"date" json:"date"`
	Details     *sql.NullString      `db:"details" json:"details"`
	Rating      *sql.NullInt64       `db:"rating" json:"rating"`
	Organized   *bool                `db:"organized" json:"organized"`
	StudioID    *sql.NullInt64       `db:"studio_id,omitempty" json:"studio_id"`
	FileModTime *NullSQLiteTimestamp `db:"file_mod_time" json:"file_mod_time"`
	CreatedAt   *SQLiteTimestamp     `db:"created_at" json:"created_at"`
	UpdatedAt   *SQLiteTimestamp     `db:"updated_at" json:"updated_at"`
}

func (s *Gallery) File() File {
	ret := File{
		Path: s.Path.String,
	}

	ret.Checksum = s.Checksum

	if s.FileModTime.Valid {
		ret.FileModTime = s.FileModTime.Timestamp
	}

	return ret
}

func (s *Gallery) SetFile(f File) {
	path := f.Path
	s.Path = sql.NullString{
		String: path,
		Valid:  true,
	}

	if f.Checksum != "" {
		s.Checksum = f.Checksum
	}

	zeroTime := time.Time{}
	if f.FileModTime != zeroTime {
		s.FileModTime = NullSQLiteTimestamp{
			Timestamp: f.FileModTime,
			Valid:     true,
		}
	}
}

// GetTitle returns the title of the scene. If the Title field is empty,
// then the base filename is returned.
func (s Gallery) GetTitle() string {
	if s.Title.String != "" {
		return s.Title.String
	}

	if s.Path.Valid {
		return filepath.Base(s.Path.String)
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
