package models

import (
	"context"
	"path/filepath"
	"strconv"
	"time"
)

// Image stores the metadata for a single image.
type Image struct {
	ID int `json:"id"`

	Title string `json:"title"`
	// Rating expressed in 1-100 scale
	Rating    *int   `json:"rating"`
	Organized bool   `json:"organized"`
	OCounter  int    `json:"o_counter"`
	StudioID  *int   `json:"studio_id"`
	URL       string `json:"url"`
	Date      *Date  `json:"date"`

	// transient - not persisted
	Files         RelatedFiles
	PrimaryFileID *FileID
	// transient - path of primary file - empty if no files
	Path string
	// transient - checksum of primary file - empty if no files
	Checksum string

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	GalleryIDs   RelatedIDs `json:"gallery_ids"`
	TagIDs       RelatedIDs `json:"tag_ids"`
	PerformerIDs RelatedIDs `json:"performer_ids"`
}

func (i *Image) LoadFiles(ctx context.Context, l FileLoader) error {
	return i.Files.load(func() ([]File, error) {
		return l.GetFiles(ctx, i.ID)
	})
}

func (i *Image) LoadPrimaryFile(ctx context.Context, l FileGetter) error {
	return i.Files.loadPrimary(func() (File, error) {
		if i.PrimaryFileID == nil {
			return nil, nil
		}

		f, err := l.Find(ctx, *i.PrimaryFileID)
		if err != nil {
			return nil, err
		}

		if len(f) > 0 {
			return f[0], nil
		}

		return nil, nil
	})
}

func (i *Image) LoadGalleryIDs(ctx context.Context, l GalleryIDLoader) error {
	return i.GalleryIDs.load(func() ([]int, error) {
		return l.GetGalleryIDs(ctx, i.ID)
	})
}

func (i *Image) LoadPerformerIDs(ctx context.Context, l PerformerIDLoader) error {
	return i.PerformerIDs.load(func() ([]int, error) {
		return l.GetPerformerIDs(ctx, i.ID)
	})
}

func (i *Image) LoadTagIDs(ctx context.Context, l TagIDLoader) error {
	return i.TagIDs.load(func() ([]int, error) {
		return l.GetTagIDs(ctx, i.ID)
	})
}

// GetTitle returns the title of the image. If the Title field is empty,
// then the base filename is returned.
func (i Image) GetTitle() string {
	if i.Title != "" {
		return i.Title
	}

	if i.Path != "" {
		return filepath.Base(i.Path)
	}

	return ""
}

// DisplayName returns a display name for the scene for logging purposes.
// It returns Path if not empty, otherwise it returns the ID.
func (i Image) DisplayName() string {
	if i.Path != "" {
		return i.Path
	}

	return strconv.Itoa(i.ID)
}

type ImageCreateInput struct {
	*Image
	FileIDs []FileID
}

type ImagePartial struct {
	Title OptionalString
	// Rating expressed in 1-100 scale
	Rating    OptionalInt
	URL       OptionalString
	Date      OptionalDate
	Organized OptionalBool
	OCounter  OptionalInt
	StudioID  OptionalInt
	CreatedAt OptionalTime
	UpdatedAt OptionalTime

	GalleryIDs    *UpdateIDs
	TagIDs        *UpdateIDs
	PerformerIDs  *UpdateIDs
	PrimaryFileID *FileID
}

func NewImagePartial() ImagePartial {
	updatedTime := time.Now()
	return ImagePartial{
		UpdatedAt: NewOptionalTime(updatedTime),
	}
}

type Images []*Image

func (i *Images) Append(o interface{}) {
	*i = append(*i, o.(*Image))
}

func (i *Images) New() interface{} {
	return &Image{}
}
