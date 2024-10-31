package models

import (
	"context"
	"errors"
	"path/filepath"
	"strconv"
	"time"
)

// Scene stores the metadata for a single video scene.
type Scene struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Code     string `json:"code"`
	Details  string `json:"details"`
	Director string `json:"director"`
	Date     *Date  `json:"date"`
	// Rating expressed in 1-100 scale
	Rating    *int `json:"rating"`
	Organized bool `json:"organized"`
	StudioID  *int `json:"studio_id"`

	// transient - not persisted
	Files         RelatedVideoFiles
	PrimaryFileID *FileID
	// transient - path of primary file - empty if no files
	Path string
	// transient - oshash of primary file - empty if no files
	OSHash string
	// transient - checksum of primary file - empty if no files
	Checksum string

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ResumeTime   float64 `json:"resume_time"`
	PlayDuration float64 `json:"play_duration"`

	URLs         RelatedStrings  `json:"urls"`
	GalleryIDs   RelatedIDs      `json:"gallery_ids"`
	TagIDs       RelatedIDs      `json:"tag_ids"`
	PerformerIDs RelatedIDs      `json:"performer_ids"`
	Groups       RelatedGroups   `json:"groups"`
	StashIDs     RelatedStashIDs `json:"stash_ids"`
}

func NewScene() Scene {
	currentTime := time.Now()
	return Scene{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}

// ScenePartial represents part of a Scene object. It is used to update
// the database entry.
type ScenePartial struct {
	Title    OptionalString
	Code     OptionalString
	Details  OptionalString
	Director OptionalString
	Date     OptionalDate
	// Rating expressed in 1-100 scale
	Rating       OptionalInt
	Organized    OptionalBool
	StudioID     OptionalInt
	CreatedAt    OptionalTime
	UpdatedAt    OptionalTime
	ResumeTime   OptionalFloat64
	PlayDuration OptionalFloat64

	URLs          *UpdateStrings
	GalleryIDs    *UpdateIDs
	TagIDs        *UpdateIDs
	PerformerIDs  *UpdateIDs
	GroupIDs      *UpdateGroupIDs
	StashIDs      *UpdateStashIDs
	PrimaryFileID *FileID
}

func NewScenePartial() ScenePartial {
	currentTime := time.Now()
	return ScenePartial{
		UpdatedAt: NewOptionalTime(currentTime),
	}
}

func (s *Scene) LoadURLs(ctx context.Context, l URLLoader) error {
	return s.URLs.load(func() ([]string, error) {
		return l.GetURLs(ctx, s.ID)
	})
}

func (s *Scene) LoadFiles(ctx context.Context, l VideoFileLoader) error {
	return s.Files.load(func() ([]*VideoFile, error) {
		return l.GetFiles(ctx, s.ID)
	})
}

func (s *Scene) LoadPrimaryFile(ctx context.Context, l FileGetter) error {
	return s.Files.loadPrimary(func() (*VideoFile, error) {
		if s.PrimaryFileID == nil {
			return nil, nil
		}

		f, err := l.Find(ctx, *s.PrimaryFileID)
		if err != nil {
			return nil, err
		}

		var vf *VideoFile
		if len(f) > 0 {
			var ok bool
			vf, ok = f[0].(*VideoFile)
			if !ok {
				return nil, errors.New("not a video file")
			}
		}
		return vf, nil
	})
}

func (s *Scene) LoadGalleryIDs(ctx context.Context, l GalleryIDLoader) error {
	return s.GalleryIDs.load(func() ([]int, error) {
		return l.GetGalleryIDs(ctx, s.ID)
	})
}

func (s *Scene) LoadPerformerIDs(ctx context.Context, l PerformerIDLoader) error {
	return s.PerformerIDs.load(func() ([]int, error) {
		return l.GetPerformerIDs(ctx, s.ID)
	})
}

func (s *Scene) LoadTagIDs(ctx context.Context, l TagIDLoader) error {
	return s.TagIDs.load(func() ([]int, error) {
		return l.GetTagIDs(ctx, s.ID)
	})
}

func (s *Scene) LoadGroups(ctx context.Context, l SceneGroupLoader) error {
	return s.Groups.load(func() ([]GroupsScenes, error) {
		return l.GetGroups(ctx, s.ID)
	})
}

func (s *Scene) LoadStashIDs(ctx context.Context, l StashIDLoader) error {
	return s.StashIDs.load(func() ([]StashID, error) {
		return l.GetStashIDs(ctx, s.ID)
	})
}

func (s *Scene) LoadRelationships(ctx context.Context, l SceneReader) error {
	if err := s.LoadURLs(ctx, l); err != nil {
		return err
	}

	if err := s.LoadGalleryIDs(ctx, l); err != nil {
		return err
	}

	if err := s.LoadPerformerIDs(ctx, l); err != nil {
		return err
	}

	if err := s.LoadTagIDs(ctx, l); err != nil {
		return err
	}

	if err := s.LoadGroups(ctx, l); err != nil {
		return err
	}

	if err := s.LoadStashIDs(ctx, l); err != nil {
		return err
	}

	if err := s.LoadFiles(ctx, l); err != nil {
		return err
	}

	return nil
}

// UpdateInput constructs a SceneUpdateInput using the populated fields in the ScenePartial object.
func (s ScenePartial) UpdateInput(id int) SceneUpdateInput {
	var dateStr *string
	if s.Date.Set {
		d := s.Date.Value
		v := d.String()
		dateStr = &v
	}

	var stashIDs StashIDs
	if s.StashIDs != nil {
		stashIDs = StashIDs(s.StashIDs.StashIDs)
	}

	ret := SceneUpdateInput{
		ID:           strconv.Itoa(id),
		Title:        s.Title.Ptr(),
		Code:         s.Code.Ptr(),
		Details:      s.Details.Ptr(),
		Director:     s.Director.Ptr(),
		Urls:         s.URLs.Strings(),
		Date:         dateStr,
		Rating100:    s.Rating.Ptr(),
		Organized:    s.Organized.Ptr(),
		StudioID:     s.StudioID.StringPtr(),
		GalleryIds:   s.GalleryIDs.IDStrings(),
		PerformerIds: s.PerformerIDs.IDStrings(),
		Movies:       s.GroupIDs.SceneMovieInputs(),
		TagIds:       s.TagIDs.IDStrings(),
		StashIds:     stashIDs.ToStashIDInputs(),
	}

	return ret
}

// GetTitle returns the title of the scene. If the Title field is empty,
// then the base filename is returned.
func (s Scene) GetTitle() string {
	if s.Title != "" {
		return s.Title
	}

	return filepath.Base(s.Path)
}

// DisplayName returns a display name for the scene for logging purposes.
// It returns Path if not empty, otherwise it returns the ID.
func (s Scene) DisplayName() string {
	if s.Path != "" {
		return s.Path
	}

	return strconv.Itoa(s.ID)
}

// GetHash returns the hash of the scene, based on the hash algorithm provided. If
// hash algorithm is MD5, then Checksum is returned. Otherwise, OSHash is returned.
func (s Scene) GetHash(hashAlgorithm HashAlgorithm) string {
	switch hashAlgorithm {
	case HashAlgorithmMd5:
		return s.Checksum
	case HashAlgorithmOshash:
		return s.OSHash
	}

	return ""
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

type VideoCaption struct {
	LanguageCode string `json:"language_code"`
	Filename     string `json:"filename"`
	CaptionType  string `json:"caption_type"`
}

func (c VideoCaption) Path(filePath string) string {
	return filepath.Join(filepath.Dir(filePath), c.Filename)
}
