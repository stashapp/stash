package api

import (
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type BaseFile interface {
	IsBaseFile()
}

type VisualFile interface {
	IsVisualFile()
}

func convertVisualFile(f models.File) (VisualFile, error) {
	switch f := f.(type) {
	case VisualFile:
		return f, nil
	case *models.VideoFile:
		return &VideoFile{VideoFile: f}, nil
	case *models.ImageFile:
		return &ImageFile{ImageFile: f}, nil
	default:
		return nil, fmt.Errorf("file %s is not a visual file", f.Base().Path)
	}
}

type GalleryFile struct {
	*models.BaseFile
}

func (GalleryFile) IsBaseFile() {}

func (GalleryFile) IsVisualFile() {}

func (f *GalleryFile) Fingerprints() []models.Fingerprint {
	return f.BaseFile.Fingerprints
}

type VideoFile struct {
	*models.VideoFile
}

func (VideoFile) IsBaseFile() {}

func (VideoFile) IsVisualFile() {}

func (f *VideoFile) Fingerprints() []models.Fingerprint {
	return f.VideoFile.Fingerprints
}

type ImageFile struct {
	*models.ImageFile
}

func (ImageFile) IsBaseFile() {}

func (ImageFile) IsVisualFile() {}

func (f *ImageFile) Fingerprints() []models.Fingerprint {
	return f.ImageFile.Fingerprints
}
