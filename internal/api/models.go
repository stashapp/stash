package api

import (
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
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

func convertBaseFile(f models.File) BaseFile {
	if f == nil {
		return nil
	}

	switch f := f.(type) {
	case BaseFile:
		return f
	case *models.VideoFile:
		return &VideoFile{VideoFile: f}
	case *models.ImageFile:
		return &ImageFile{ImageFile: f}
	case *models.BaseFile:
		// assume gallery file if it's not a video or image file
		return &GalleryFile{BaseFile: f}
	default:
		panic(fmt.Errorf("unknown file type %T", f))
	}
}

func convertBaseFiles(files []models.File) []BaseFile {
	return sliceutil.Map(files, convertBaseFile)
}

type GalleryFile struct {
	*models.BaseFile
}

func (GalleryFile) IsBaseFile() {}

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

type AudioFile struct {
	*models.AudioFile
}

func (AudioFile) IsBaseFile() {}

func (f *AudioFile) Fingerprints() []models.Fingerprint {
	return f.AudioFile.Fingerprints
}

type ImageFile struct {
	*models.ImageFile
}

func (ImageFile) IsBaseFile() {}

func (ImageFile) IsVisualFile() {}

func (f *ImageFile) Fingerprints() []models.Fingerprint {
	return f.ImageFile.Fingerprints
}

type BasicFile struct {
	*models.BaseFile
}

func (BasicFile) IsBaseFile() {}

func (BasicFile) IsVisualFile() {}

func (f *BasicFile) Fingerprints() []models.Fingerprint {
	return f.BaseFile.Fingerprints
}
