package models

import (
	"path/filepath"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/file"
)

// Scene stores the metadata for a single video scene.
type Scene struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Details   string `json:"details"`
	URL       string `json:"url"`
	Date      *Date  `json:"date"`
	Rating    *int   `json:"rating"`
	Organized bool   `json:"organized"`
	OCounter  int    `json:"o_counter"`
	StudioID  *int   `json:"studio_id"`

	// transient - not persisted
	Files []*file.VideoFile

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	GalleryIDs   []int          `json:"gallery_ids"`
	TagIDs       []int          `json:"tag_ids"`
	PerformerIDs []int          `json:"performer_ids"`
	Movies       []MoviesScenes `json:"movies"`
	StashIDs     []StashID      `json:"stash_ids"`
}

func (s Scene) PrimaryFile() *file.VideoFile {
	if len(s.Files) == 0 {
		return nil
	}

	return s.Files[0]
}

func (s Scene) Path() string {
	if p := s.PrimaryFile(); p != nil {
		return p.Base().Path
	}

	return ""
}

func (s Scene) getHash(type_ string) string {
	if p := s.PrimaryFile(); p != nil {
		v := p.Base().Fingerprints.Get(type_)
		if v == nil {
			return ""
		}

		return v.(string)
	}
	return ""
}

func (s Scene) Checksum() string {
	return s.getHash(file.FingerprintTypeMD5)
}

func (s Scene) OSHash() string {
	return s.getHash(file.FingerprintTypeOshash)
}

func (s Scene) Phash() int64 {
	if p := s.PrimaryFile(); p != nil {
		v := p.Base().Fingerprints.Get(file.FingerprintTypePhash)
		if v == nil {
			return 0
		}

		return v.(int64)
	}
	return 0
}

func (s Scene) Duration() float64 {
	if p := s.PrimaryFile(); p != nil {
		return p.Duration
	}

	return 0
}

func (s Scene) Format() string {
	if p := s.PrimaryFile(); p != nil {
		return p.Format
	}

	return ""
}

func (s Scene) VideoCodec() string {
	if p := s.PrimaryFile(); p != nil {
		return p.VideoCodec
	}

	return ""
}

func (s Scene) AudioCodec() string {
	if p := s.PrimaryFile(); p != nil {
		return p.AudioCodec
	}

	return ""
}

// ScenePartial represents part of a Scene object. It is used to update
// the database entry.
type ScenePartial struct {
	Title     OptionalString
	Details   OptionalString
	URL       OptionalString
	Date      OptionalDate
	Rating    OptionalInt
	Organized OptionalBool
	OCounter  OptionalInt
	StudioID  OptionalInt
	CreatedAt OptionalTime
	UpdatedAt OptionalTime

	GalleryIDs   *UpdateIDs
	TagIDs       *UpdateIDs
	PerformerIDs *UpdateIDs
	MovieIDs     *UpdateMovieIDs
	StashIDs     *UpdateStashIDs
}

func NewScenePartial() ScenePartial {
	updatedTime := time.Now()
	return ScenePartial{
		UpdatedAt: NewOptionalTime(updatedTime),
	}
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
func (s ScenePartial) UpdateInput(id int) SceneUpdateInput {
	var dateStr *string
	if s.Date.Set {
		d := s.Date.Value
		v := d.String()
		dateStr = &v
	}

	var stashIDs []StashID
	if s.StashIDs != nil {
		stashIDs = s.StashIDs.StashIDs
	}

	return SceneUpdateInput{
		ID:           strconv.Itoa(id),
		Title:        s.Title.Ptr(),
		Details:      s.Details.Ptr(),
		URL:          s.URL.Ptr(),
		Date:         dateStr,
		Rating:       s.Rating.Ptr(),
		Organized:    s.Organized.Ptr(),
		StudioID:     s.StudioID.StringPtr(),
		GalleryIds:   s.GalleryIDs.IDStrings(),
		PerformerIds: s.PerformerIDs.IDStrings(),
		Movies:       s.MovieIDs.SceneMovieInputs(),
		TagIds:       s.TagIDs.IDStrings(),
		StashIds:     stashIDs,
	}
}

// GetTitle returns the title of the scene. If the Title field is empty,
// then the base filename is returned.
func (s Scene) GetTitle() string {
	if s.Title != "" {
		return s.Title
	}

	return filepath.Base(s.Path())
}

// GetHash returns the hash of the scene, based on the hash algorithm provided. If
// hash algorithm is MD5, then Checksum is returned. Otherwise, OSHash is returned.
func (s Scene) GetHash(hashAlgorithm HashAlgorithm) string {
	f := s.PrimaryFile()
	if f == nil {
		return ""
	}

	switch hashAlgorithm {
	case HashAlgorithmMd5:
		return f.Base().Fingerprints.Get(file.FingerprintTypeMD5).(string)
	case HashAlgorithmOshash:
		return f.Base().Fingerprints.Get(file.FingerprintTypeOshash).(string)
	}

	return ""
}

func (s Scene) GetMinResolution() int {
	f := s.PrimaryFile()
	if f == nil {
		return 0
	}

	w := f.Width
	h := f.Height

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

type VideoCaption struct {
	LanguageCode string `json:"language_code"`
	Filename     string `json:"filename"`
	CaptionType  string `json:"caption_type"`
}

func (c VideoCaption) Path(filePath string) string {
	return filepath.Join(filepath.Dir(filePath), c.Filename)
}
