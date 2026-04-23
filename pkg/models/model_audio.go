// TODO(audio): update this file

package models

import (
	"context"
	"errors"
	"path/filepath"
	"strconv"
	"time"
)

// Audio stores the metadata for a single video audio.
type Audio struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Code    string `json:"code"`
	Details string `json:"details"`
	Artists string `json:"artists"`
	Date    *Date  `json:"date"`
	// Rating expressed in 1-100 scale
	Rating    *int `json:"rating"`
	Organized bool `json:"organized"`
	StudioID  *int `json:"studio_id"`

	// transient - not persisted
	Files         RelatedAudioFiles
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

	URLs         RelatedStrings     `json:"urls"`
	GalleryIDs   RelatedIDs         `json:"gallery_ids"`
	TagIDs       RelatedIDs         `json:"tag_ids"`
	PerformerIDs RelatedIDs         `json:"performer_ids"`
	Groups       RelatedGroupsAudio `json:"groups"`
}

func NewAudio() Audio {
	currentTime := time.Now()
	return Audio{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}

type CreateAudioInput struct {
	*Audio

	FileIDs      []FileID
	CoverImage   []byte
	CustomFields CustomFieldMap `json:"custom_fields"`
}

type UpdateAudioInput struct {
	*Audio

	CustomFields CustomFieldsInput `json:"custom_fields"`
}

// AudioPartial represents part of a Audio object. It is used to update
// the database entry.
type AudioPartial struct {
	Title   OptionalString
	Code    OptionalString
	Details OptionalString
	Date    OptionalDate
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
	GroupIDs      *UpdateGroupIDsAudio
	PrimaryFileID *FileID
}

func NewAudioPartial() AudioPartial {
	currentTime := time.Now()
	return AudioPartial{
		UpdatedAt: NewOptionalTime(currentTime),
	}
}

func (s *Audio) LoadURLs(ctx context.Context, l URLLoader) error {
	return s.URLs.load(func() ([]string, error) {
		return l.GetURLs(ctx, s.ID)
	})
}

func (s *Audio) LoadFiles(ctx context.Context, l AudioFileLoader) error {
	return s.Files.load(func() ([]*AudioFile, error) {
		return l.GetFiles(ctx, s.ID)
	})
}

func (s *Audio) LoadPrimaryFile(ctx context.Context, l FileGetter) error {
	return s.Files.loadPrimary(func() (*AudioFile, error) {
		if s.PrimaryFileID == nil {
			return nil, nil
		}

		f, err := l.Find(ctx, *s.PrimaryFileID)
		if err != nil {
			return nil, err
		}

		var vf *AudioFile
		if len(f) > 0 {
			var ok bool
			vf, ok = f[0].(*AudioFile)
			if !ok {
				return nil, errors.New("not a video file")
			}
		}
		return vf, nil
	})
}

func (s *Audio) LoadGalleryIDs(ctx context.Context, l GalleryIDLoader) error {
	return s.GalleryIDs.load(func() ([]int, error) {
		return l.GetGalleryIDs(ctx, s.ID)
	})
}

func (s *Audio) LoadPerformerIDs(ctx context.Context, l PerformerIDLoader) error {
	return s.PerformerIDs.load(func() ([]int, error) {
		return l.GetPerformerIDs(ctx, s.ID)
	})
}

func (s *Audio) LoadTagIDs(ctx context.Context, l TagIDLoader) error {
	return s.TagIDs.load(func() ([]int, error) {
		return l.GetTagIDs(ctx, s.ID)
	})
}

func (s *Audio) LoadGroups(ctx context.Context, l AudioGroupLoader) error {
	return s.Groups.load(func() ([]GroupsAudios, error) {
		return l.GetGroups(ctx, s.ID)
	})
}

func (s *Audio) LoadRelationships(ctx context.Context, l AudioReader) error {
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

	if err := s.LoadFiles(ctx, l); err != nil {
		return err
	}

	return nil
}

// UpdateInput constructs a AudioUpdateInput using the populated fields in the AudioPartial object.
func (s AudioPartial) UpdateInput(id int) AudioUpdateInput {
	var dateStr *string
	if s.Date.Set {
		d := s.Date.Value
		v := d.String()
		dateStr = &v
	}

	ret := AudioUpdateInput{
		ID:           strconv.Itoa(id),
		Title:        s.Title.Ptr(),
		Code:         s.Code.Ptr(),
		Details:      s.Details.Ptr(),
		Urls:         s.URLs.Strings(),
		Date:         dateStr,
		Rating100:    s.Rating.Ptr(),
		Organized:    s.Organized.Ptr(),
		StudioID:     s.StudioID.StringPtr(),
		PerformerIds: s.PerformerIDs.IDStrings(),
		Groups:       s.GroupIDs.GroupInputs(),
		TagIds:       s.TagIDs.IDStrings(),
	}

	return ret
}

// GetTitle returns the title of the audio. If the Title field is empty,
// then the base filename is returned.
func (s Audio) GetTitle() string {
	if s.Title != "" {
		return s.Title
	}

	return filepath.Base(s.Path)
}

// DisplayName returns a display name for the audio for logging purposes.
// It returns Path if not empty, otherwise it returns the ID.
func (s Audio) DisplayName() string {
	if s.Path != "" {
		return s.Path
	}

	return strconv.Itoa(s.ID)
}

// GetHash returns the hash of the audio, based on the hash algorithm provided. If
// hash algorithm is MD5, then Checksum is returned. Otherwise, OSHash is returned.
func (s Audio) GetHash(hashAlgorithm HashAlgorithm) string {
	switch hashAlgorithm {
	case HashAlgorithmMd5:
		return s.Checksum
	case HashAlgorithmOshash:
		return s.OSHash
	}

	return ""
}

// AudioFileType represents the file metadata for a audio.
type AudioFileType struct {
	Size       *string  `graphql:"size" json:"size"`
	Duration   *float64 `graphql:"duration" json:"duration"`
	AudioCodec *string  `graphql:"audio_codec" json:"audio_codec"`
	Samplerate *float64 `graphql:"samplerate" json:"samplerate"`
	Bitrate    *int     `graphql:"bitrate" json:"bitrate"`
}

// TODO(audio): don't know if we need this, using VideoCaption for now due to `pkg/models/repository_file.go` and `FileReader` using
// type AudioCaption struct {
// 	LanguageCode string `json:"language_code"`
// 	Filename     string `json:"filename"`
// 	CaptionType  string `json:"caption_type"`
// }

// func (c AudioCaption) Path(filePath string) string {
// 	return filepath.Join(filepath.Dir(filePath), c.Filename)
// }
