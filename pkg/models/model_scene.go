package models

import (
	"database/sql"
	"path/filepath"
	"strconv"
	"time"
)

// Scene stores the metadata for a single video scene.
type Scene struct {
	ID               int                 `db:"id" json:"id"`
	Checksum         sql.NullString      `db:"checksum" json:"checksum"`
	OSHash           sql.NullString      `db:"oshash" json:"oshash"`
	Path             string              `db:"path" json:"path"`
	Title            sql.NullString      `db:"title" json:"title"`
	Details          sql.NullString      `db:"details" json:"details"`
	URL              sql.NullString      `db:"url" json:"url"`
	Date             SQLiteDate          `db:"date" json:"date"`
	Rating           sql.NullInt64       `db:"rating" json:"rating"`
	Organized        bool                `db:"organized" json:"organized"`
	OCounter         int                 `db:"o_counter" json:"o_counter"`
	Size             sql.NullString      `db:"size" json:"size"`
	Duration         sql.NullFloat64     `db:"duration" json:"duration"`
	VideoCodec       sql.NullString      `db:"video_codec" json:"video_codec"`
	Format           sql.NullString      `db:"format" json:"format_name"`
	AudioCodec       sql.NullString      `db:"audio_codec" json:"audio_codec"`
	Width            sql.NullInt64       `db:"width" json:"width"`
	Height           sql.NullInt64       `db:"height" json:"height"`
	Framerate        sql.NullFloat64     `db:"framerate" json:"framerate"`
	Bitrate          sql.NullInt64       `db:"bitrate" json:"bitrate"`
	StudioID         sql.NullInt64       `db:"studio_id,omitempty" json:"studio_id"`
	FileModTime      NullSQLiteTimestamp `db:"file_mod_time" json:"file_mod_time"`
	Phash            sql.NullInt64       `db:"phash,omitempty" json:"phash"`
	CreatedAt        SQLiteTimestamp     `db:"created_at" json:"created_at"`
	UpdatedAt        SQLiteTimestamp     `db:"updated_at" json:"updated_at"`
	Interactive      bool                `db:"interactive" json:"interactive"`
	InteractiveSpeed sql.NullInt64       `db:"interactive_speed" json:"interactive_speed"`
}

func (s *Scene) File() File {
	ret := File{
		Path: s.Path,
	}

	if s.Checksum.Valid {
		ret.Checksum = s.Checksum.String
	}
	if s.OSHash.Valid {
		ret.OSHash = s.OSHash.String
	}
	if s.FileModTime.Valid {
		ret.FileModTime = s.FileModTime.Timestamp
	}
	if s.Size.Valid {
		ret.Size = s.Size.String
	}

	return ret
}

func (s *Scene) SetFile(f File) {
	path := f.Path
	s.Path = path

	if f.Checksum != "" {
		s.Checksum = sql.NullString{
			String: f.Checksum,
			Valid:  true,
		}
	}
	if f.OSHash != "" {
		s.OSHash = sql.NullString{
			String: f.OSHash,
			Valid:  true,
		}
	}
	zeroTime := time.Time{}
	if f.FileModTime != zeroTime {
		s.FileModTime = NullSQLiteTimestamp{
			Timestamp: f.FileModTime,
			Valid:     true,
		}
	}
	if f.Size != "" {
		s.Size = sql.NullString{
			String: f.Size,
			Valid:  true,
		}
	}
}

// ScenePartial represents part of a Scene object. It is used to update
// the database entry. Only non-nil fields will be updated.
type ScenePartial struct {
	ID               int                  `db:"id" json:"id"`
	Checksum         *sql.NullString      `db:"checksum" json:"checksum"`
	OSHash           *sql.NullString      `db:"oshash" json:"oshash"`
	Path             *string              `db:"path" json:"path"`
	Title            *sql.NullString      `db:"title" json:"title"`
	Details          *sql.NullString      `db:"details" json:"details"`
	URL              *sql.NullString      `db:"url" json:"url"`
	Date             *SQLiteDate          `db:"date" json:"date"`
	Rating           *sql.NullInt64       `db:"rating" json:"rating"`
	Organized        *bool                `db:"organized" json:"organized"`
	Size             *sql.NullString      `db:"size" json:"size"`
	Duration         *sql.NullFloat64     `db:"duration" json:"duration"`
	VideoCodec       *sql.NullString      `db:"video_codec" json:"video_codec"`
	Format           *sql.NullString      `db:"format" json:"format_name"`
	AudioCodec       *sql.NullString      `db:"audio_codec" json:"audio_codec"`
	Width            *sql.NullInt64       `db:"width" json:"width"`
	Height           *sql.NullInt64       `db:"height" json:"height"`
	Framerate        *sql.NullFloat64     `db:"framerate" json:"framerate"`
	Bitrate          *sql.NullInt64       `db:"bitrate" json:"bitrate"`
	StudioID         *sql.NullInt64       `db:"studio_id,omitempty" json:"studio_id"`
	MovieID          *sql.NullInt64       `db:"movie_id,omitempty" json:"movie_id"`
	FileModTime      *NullSQLiteTimestamp `db:"file_mod_time" json:"file_mod_time"`
	Phash            *sql.NullInt64       `db:"phash,omitempty" json:"phash"`
	CreatedAt        *SQLiteTimestamp     `db:"created_at" json:"created_at"`
	UpdatedAt        *SQLiteTimestamp     `db:"updated_at" json:"updated_at"`
	Interactive      *bool                `db:"interactive" json:"interactive"`
	InteractiveSpeed *sql.NullInt64       `db:"interactive_speed" json:"interactive_speed"`
}

type SceneMovieInput struct {
	MovieID    string `json:"movie_id"`
	SceneIndex *int   `json:"scene_index"`
}

type SceneUpdateInput struct {
	ClientMutationID *string            `json:"clientMutationId"`
	ID               string             `json:"id"`
	Title            *string            `json:"title"`
	Details          *string            `json:"details"`
	URL              *string            `json:"url"`
	Date             *string            `json:"date"`
	Rating           *int               `json:"rating"`
	Organized        *bool              `json:"organized"`
	StudioID         *string            `json:"studio_id"`
	GalleryIds       []string           `json:"gallery_ids"`
	PerformerIds     []string           `json:"performer_ids"`
	Movies           []*SceneMovieInput `json:"movies"`
	TagIds           []string           `json:"tag_ids"`
	// This should be a URL or a base64 encoded data URL
	CoverImage *string         `json:"cover_image"`
	StashIds   []*StashIDInput `json:"stash_ids"`
}

// UpdateInput constructs a SceneUpdateInput using the populated fields in the ScenePartial object.
func (s ScenePartial) UpdateInput() SceneUpdateInput {
	boolPtrCopy := func(v *bool) *bool {
		if v == nil {
			return nil
		}

		vv := *v
		return &vv
	}

	return SceneUpdateInput{
		ID:        strconv.Itoa(s.ID),
		Title:     nullStringPtrToStringPtr(s.Title),
		Details:   nullStringPtrToStringPtr(s.Details),
		URL:       nullStringPtrToStringPtr(s.URL),
		Date:      s.Date.StringPtr(),
		Rating:    nullInt64PtrToIntPtr(s.Rating),
		Organized: boolPtrCopy(s.Organized),
		StudioID:  nullInt64PtrToStringPtr(s.StudioID),
	}
}

func (s *ScenePartial) SetFile(f File) {
	path := f.Path
	s.Path = &path

	if f.Checksum != "" {
		s.Checksum = &sql.NullString{
			String: f.Checksum,
			Valid:  true,
		}
	}
	if f.OSHash != "" {
		s.OSHash = &sql.NullString{
			String: f.OSHash,
			Valid:  true,
		}
	}
	zeroTime := time.Time{}
	if f.FileModTime != zeroTime {
		s.FileModTime = &NullSQLiteTimestamp{
			Timestamp: f.FileModTime,
			Valid:     true,
		}
	}
	if f.Size != "" {
		s.Size = &sql.NullString{
			String: f.Size,
			Valid:  true,
		}
	}
}

// GetTitle returns the title of the scene. If the Title field is empty,
// then the base filename is returned.
func (s Scene) GetTitle() string {
	if s.Title.String != "" {
		return s.Title.String
	}

	return filepath.Base(s.Path)
}

// GetHash returns the hash of the scene, based on the hash algorithm provided. If
// hash algorithm is MD5, then Checksum is returned. Otherwise, OSHash is returned.
func (s Scene) GetHash(hashAlgorithm HashAlgorithm) string {
	return s.File().GetHash(hashAlgorithm)
}

func (s Scene) GetMinResolution() int64 {
	if s.Width.Int64 < s.Height.Int64 {
		return s.Width.Int64
	}

	return s.Height.Int64
}

// SceneFileType represents the file metadata for a scene.
type SceneFileType struct {
	Size       *string  `graphql:"size" json:"size"`
	Duration   *float64 `graphql:"duration" json:"duration"`
	VideoCodec *string  `graphql:"video_codec" json:"video_codec"`
	AudioCodec *string  `graphql:"audio_codec" json:"audio_codec"`
	Width      *int     `graphql:"width" json:"width"`
	Height     *int     `graphql:"height" json:"height"`
	Framerate  *float64 `graphql:"framerate" json:"framerate"`
	Bitrate    *int     `graphql:"bitrate" json:"bitrate"`
}

type Scenes []*Scene

func (s *Scenes) Append(o interface{}) {
	*s = append(*s, o.(*Scene))
}

func (s *Scenes) New() interface{} {
	return &Scene{}
}

type SceneCaption struct {
	LanguageCode string `json:"language_code"`
	Filename     string `json:"filename"`
	CaptionType  string `json:"caption_type"`
}

func (c SceneCaption) Path(scenePath string) string {
	return filepath.Join(filepath.Dir(scenePath), c.Filename)
}
