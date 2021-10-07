package models

import (
	"database/sql"
	"path/filepath"
	"strconv"
	"time"
)

// Image stores the metadata for a single image.
type Image struct {
	ID          int                 `db:"id" json:"id"`
	Checksum    string              `db:"checksum" json:"checksum"`
	Path        string              `db:"path" json:"path"`
	Title       sql.NullString      `db:"title" json:"title"`
	Rating      sql.NullInt64       `db:"rating" json:"rating"`
	Organized   bool                `db:"organized" json:"organized"`
	OCounter    int                 `db:"o_counter" json:"o_counter"`
	Size        sql.NullInt64       `db:"size" json:"size"`
	Width       sql.NullInt64       `db:"width" json:"width"`
	Height      sql.NullInt64       `db:"height" json:"height"`
	StudioID    sql.NullInt64       `db:"studio_id,omitempty" json:"studio_id"`
	FileModTime NullSQLiteTimestamp `db:"file_mod_time" json:"file_mod_time"`
	CreatedAt   SQLiteTimestamp     `db:"created_at" json:"created_at"`
	UpdatedAt   SQLiteTimestamp     `db:"updated_at" json:"updated_at"`
}

// ImagePartial represents part of a Image object. It is used to update
// the database entry. Only non-nil fields will be updated.
type ImagePartial struct {
	ID          int                  `db:"id" json:"id"`
	Checksum    *string              `db:"checksum" json:"checksum"`
	Path        *string              `db:"path" json:"path"`
	Title       *sql.NullString      `db:"title" json:"title"`
	Rating      *sql.NullInt64       `db:"rating" json:"rating"`
	Organized   *bool                `db:"organized" json:"organized"`
	Size        *sql.NullInt64       `db:"size" json:"size"`
	Width       *sql.NullInt64       `db:"width" json:"width"`
	Height      *sql.NullInt64       `db:"height" json:"height"`
	StudioID    *sql.NullInt64       `db:"studio_id,omitempty" json:"studio_id"`
	FileModTime *NullSQLiteTimestamp `db:"file_mod_time" json:"file_mod_time"`
	CreatedAt   *SQLiteTimestamp     `db:"created_at" json:"created_at"`
	UpdatedAt   *SQLiteTimestamp     `db:"updated_at" json:"updated_at"`
}

func (i *Image) File() File {
	ret := File{
		Path: i.Path,
	}

	ret.Checksum = i.Checksum
	if i.FileModTime.Valid {
		ret.FileModTime = i.FileModTime.Timestamp
	}
	if i.Size.Valid {
		ret.Size = strconv.FormatInt(i.Size.Int64, 10)
	}

	return ret
}

func (i *Image) SetFile(f File) {
	path := f.Path
	i.Path = path

	if f.Checksum != "" {
		i.Checksum = f.Checksum
	}
	zeroTime := time.Time{}
	if f.FileModTime != zeroTime {
		i.FileModTime = NullSQLiteTimestamp{
			Timestamp: f.FileModTime,
			Valid:     true,
		}
	}
	if f.Size != "" {
		size, err := strconv.ParseInt(f.Size, 10, 64)
		if err == nil {
			i.Size = sql.NullInt64{
				Int64: size,
				Valid: true,
			}
		}
	}
}

// GetTitle returns the title of the image. If the Title field is empty,
// then the base filename is returned.
func (i *Image) GetTitle() string {
	if i.Title.String != "" {
		return i.Title.String
	}

	return filepath.Base(i.Path)
}

// ImageFileType represents the file metadata for an image.
type ImageFileType struct {
	Size   *int `graphql:"size" json:"size"`
	Width  *int `graphql:"width" json:"width"`
	Height *int `graphql:"height" json:"height"`
}

type Images []*Image

func (i *Images) Append(o interface{}) {
	*i = append(*i, o.(*Image))
}

func (i *Images) New() interface{} {
	return &Image{}
}
