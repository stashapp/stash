package models

import (
	"context"
	"errors"
	"path/filepath"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/file"
)

// Image stores the metadata for a single image.
type Image struct {
	ID int `json:"id"`

	Title string `json:"title"`
	// Rating expressed in 1-100 scale
	Rating    *int `json:"rating"`
	Organized bool `json:"organized"`
	OCounter  int  `json:"o_counter"`
	StudioID  *int `json:"studio_id"`

	// transient - not persisted
	Files         RelatedImageFiles
	PrimaryFileID *file.ID
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

func (i *Image) LoadFiles(ctx context.Context, l ImageFileLoader) error {
	return i.Files.load(func() ([]*file.ImageFile, error) {
		return l.GetFiles(ctx, i.ID)
	})
}

func (i *Image) LoadPrimaryFile(ctx context.Context, l file.Finder) error {
	return i.Files.loadPrimary(func() (*file.ImageFile, error) {
		if i.PrimaryFileID == nil {
			return nil, nil
		}

		f, err := l.Find(ctx, *i.PrimaryFileID)
		if err != nil {
			return nil, err
		}

		var vf *file.ImageFile
		if len(f) > 0 {
			var ok bool
			vf, ok = f[0].(*file.ImageFile)
			if !ok {
				return nil, errors.New("not an image file")
			}
		}
		return vf, nil
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
	FileIDs []file.ID
}

type ImagePartial struct {
	Title OptionalString
	// Rating expressed in 1-100 scale
	Rating    OptionalInt
	Organized OptionalBool
	OCounter  OptionalInt
	StudioID  OptionalInt
	CreatedAt OptionalTime
	UpdatedAt OptionalTime

	GalleryIDs    *UpdateIDs
	TagIDs        *UpdateIDs
	PerformerIDs  *UpdateIDs
	PrimaryFileID *file.ID
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
