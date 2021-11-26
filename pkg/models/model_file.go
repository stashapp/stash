package models

import (
	"database/sql"
	"time"
)

type File struct {
	ID          int             `db:"id" json:"id"`
	Checksum    string          `db:"checksum" json:"checksum"`
	OSHash      string          `db:"oshash" json:"oshash"`
	Path        string          `db:"path" json:"path"`
	ZipFileID   sql.NullInt64   `db:"zip_file_id" json:"zip_file_id"`
	Size        int64           `db:"size" json:"size"`
	Duration    sql.NullFloat64 `db:"duration" json:"duration"`
	VideoCodec  sql.NullString  `db:"video_codec" json:"video_codec"`
	Format      sql.NullString  `db:"format" json:"format_name"`
	AudioCodec  sql.NullString  `db:"audio_codec" json:"audio_codec"`
	Width       sql.NullInt64   `db:"width" json:"width"`
	Height      sql.NullInt64   `db:"height" json:"height"`
	Framerate   sql.NullFloat64 `db:"framerate" json:"framerate"`
	Bitrate     sql.NullInt64   `db:"bitrate" json:"bitrate"`
	FileModTime time.Time       `db:"mod_time" json:"mod_time"`
	CreatedAt   SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt   SQLiteTimestamp `db:"updated_at" json:"updated_at"`

	// resolved, not stored
	ZipPath sql.NullString `db:"zip_path" json:"zip_path" resolved:"true"`
}

// GetHash returns the hash of the scene, based on the hash algorithm provided. If
// hash algorithm is MD5, then Checksum is returned. Otherwise, OSHash is returned.
func (s File) GetHash(hashAlgorithm HashAlgorithm) string {
	switch hashAlgorithm {
	case HashAlgorithmMd5:
		return s.Checksum
	case HashAlgorithmOshash:
		return s.OSHash
	default:
		panic("unknown hash algorithm")
	}
}

// Equal returns true if path, checksum, size and mod time are equal to the
// values in the provided file.
func (s File) Equal(o File) bool {
	return s.Path == o.Path && s.Checksum == o.Checksum && s.OSHash == o.OSHash && s.Size == o.Size && s.FileModTime.Equal(o.FileModTime)
}

type Files []*File

func (s *Files) Append(o interface{}) {
	*s = append(*s, o.(*File))
}

func (s *Files) New() interface{} {
	return &File{}
}
