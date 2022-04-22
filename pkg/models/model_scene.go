package models

import (
	"path/filepath"
	"strconv"
	"time"
)

// Scene stores the metadata for a single video scene.
type Scene struct {
	ID               int        `json:"id"`
	Checksum         *string    `json:"checksum"`
	OSHash           *string    `json:"oshash"`
	Path             string     `json:"path"`
	Title            *string    `json:"title"`
	Details          *string    `json:"details"`
	URL              *string    `json:"url"`
	Date             *Date      `json:"date"`
	Rating           *int       `json:"rating"`
	Organized        bool       `json:"organized"`
	OCounter         int        `json:"o_counter"`
	Size             *string    `json:"size"`
	Duration         *float64   `json:"duration"`
	VideoCodec       *string    `json:"video_codec"`
	Format           *string    `json:"format_name"`
	AudioCodec       *string    `json:"audio_codec"`
	Width            *int       `json:"width"`
	Height           *int       `json:"height"`
	Framerate        *float64   `json:"framerate"`
	Bitrate          *int64     `json:"bitrate"`
	StudioID         *int       `json:"studio_id"`
	FileModTime      *time.Time `json:"file_mod_time"`
	Phash            *int64     `json:"phash"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	Interactive      bool       `json:"interactive"`
	InteractiveSpeed *int       `json:"interactive_speed"`

	GalleryIDs   []int          `json:"gallery_ids"`
	TagIDs       []int          `json:"tag_ids"`
	PerformerIDs []int          `json:"performer_ids"`
	Movies       []MoviesScenes `json:"movies"`
	StashIDs     []StashID      `json:"stash_ids"`
}

func (s *Scene) File() File {
	ret := File{
		Path: s.Path,
	}

	if s.Checksum != nil {
		ret.Checksum = *s.Checksum
	}
	if s.OSHash != nil {
		ret.OSHash = *s.OSHash
	}
	if s.FileModTime != nil {
		ret.FileModTime = *s.FileModTime
	}
	if s.Size != nil {
		ret.Size = *s.Size
	}

	return ret
}

func (s *Scene) SetFile(f File) {
	path := f.Path
	s.Path = path

	if f.Checksum != "" {
		s.Checksum = &f.Checksum
	}
	if f.OSHash != "" {
		s.OSHash = &f.OSHash
	}
	zeroTime := time.Time{}
	if f.FileModTime != zeroTime {
		s.FileModTime = &f.FileModTime
	}
	if f.Size != "" {
		s.Size = &f.Size
	}
}

// ScenePartial represents part of a Scene object. It is used to update
// the database entry. Only non-nil fields will be updated.
type ScenePartial struct {
	ID               int
	Checksum         **string
	OSHash           **string
	Path             *string
	Title            **string
	Details          **string
	URL              **string
	Date             **Date
	Rating           **int
	Organized        *bool
	Size             **string
	Duration         **float64
	VideoCodec       **string
	Format           **string
	AudioCodec       **string
	Width            **int
	Height           **int
	Framerate        **float64
	Bitrate          **int
	StudioID         **int
	FileModTime      *time.Time
	Phash            **int64
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
	Interactive      *bool
	InteractiveSpeed **int

	GalleryIDs   *UpdateIDs
	TagIDs       *UpdateIDs
	PerformerIDs *UpdateIDs
	MovieIDs     *UpdateMovieIDs
	StashIDs     *UpdateStashIDs
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
	CoverImage *string   `json:"cover_image"`
	StashIds   []StashID `json:"stash_ids"`
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

	var dateStr *string
	if s.Date != nil {
		d := *s.Date
		v := d.String()
		dateStr = &v
	}

	return SceneUpdateInput{
		ID:           strconv.Itoa(s.ID),
		Title:        stringDblPtrToPtr(s.Title),
		Details:      stringDblPtrToPtr(s.Details),
		URL:          stringDblPtrToPtr(s.URL),
		Date:         dateStr,
		Rating:       intDblPtrToPtr(s.Rating),
		Organized:    boolPtrCopy(s.Organized),
		StudioID:     intDblPtrToStringPtr(s.StudioID),
		GalleryIds:   s.GalleryIDs.IDStrings(),
		PerformerIds: s.PerformerIDs.IDStrings(),
		Movies:       s.MovieIDs.SceneMovieInputs(),
		TagIds:       s.TagIDs.IDStrings(),
		StashIds:     s.StashIDs.StashIDs,
	}
}

func (s *ScenePartial) SetFile(f File) {
	path := f.Path
	s.Path = &path

	if f.Checksum != "" {
		v := &f.Checksum
		s.Checksum = &v
	}
	if f.OSHash != "" {
		v := &f.OSHash
		s.OSHash = &v
	}
	zeroTime := time.Time{}
	if f.FileModTime != zeroTime {
		s.FileModTime = &f.FileModTime
	}
	if f.Size != "" {
		v := &f.Size
		s.Size = &v
	}
}

// GetTitle returns the title of the scene. If the Title field is empty,
// then the base filename is returned.
func (s Scene) GetTitle() string {
	if s.Title != nil && *s.Title != "" {
		return *s.Title
	}

	return filepath.Base(s.Path)
}

// GetHash returns the hash of the scene, based on the hash algorithm provided. If
// hash algorithm is MD5, then Checksum is returned. Otherwise, OSHash is returned.
func (s Scene) GetHash(hashAlgorithm HashAlgorithm) string {
	return s.File().GetHash(hashAlgorithm)
}

func (s Scene) GetMinResolution() int {
	var w, h int
	if s.Width != nil {
		w = *s.Width
	}
	if s.Height != nil {
		h = *s.Height
	}
	if w < h {
		return w
	}

	return h
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
