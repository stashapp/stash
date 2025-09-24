package models

import (
	"context"
	"path/filepath"
	"strings"
)

type FileQueryOptions struct {
	QueryOptions
	FileFilter *FileFilterType

	TotalDuration bool
	Megapixels    bool
	TotalSize     bool
}

type FileFilterType struct {
	OperatorFilter[FileFilterType]

	// Filter by path
	Path *StringCriterionInput `json:"path"`

	Basename        *StringCriterionInput            `json:"basename"`
	Dir             *StringCriterionInput            `json:"dir"`
	ParentFolder    *HierarchicalMultiCriterionInput `json:"parent_folder"`
	ZipFile         *MultiCriterionInput             `json:"zip_file"`
	ModTime         *TimestampCriterionInput         `json:"mod_time"`
	Duplicated      *PHashDuplicationCriterionInput  `json:"duplicated"`
	Hashes          []*FingerprintFilterInput        `json:"hashes"`
	VideoFileFilter *VideoFileFilterInput            `json:"video_file_filter"`
	ImageFileFilter *ImageFileFilterInput            `json:"image_file_filter"`
	SceneCount      *IntCriterionInput               `json:"scene_count"`
	ImageCount      *IntCriterionInput               `json:"image_count"`
	GalleryCount    *IntCriterionInput               `json:"gallery_count"`
	ScenesFilter    *SceneFilterType                 `json:"scenes_filter"`
	ImagesFilter    *ImageFilterType                 `json:"images_filter"`
	GalleriesFilter *GalleryFilterType               `json:"galleries_filter"`
	CreatedAt       *TimestampCriterionInput         `json:"created_at"`
	UpdatedAt       *TimestampCriterionInput         `json:"updated_at"`
}

func PathsFileFilter(paths []string) *FileFilterType {
	if paths == nil {
		return nil
	}

	sep := string(filepath.Separator)

	var ret *FileFilterType
	var or *FileFilterType
	for _, p := range paths {
		newOr := &FileFilterType{}
		if or != nil {
			or.Or = newOr
		} else {
			ret = newOr
		}

		or = newOr

		if !strings.HasSuffix(p, sep) {
			p += sep
		}

		or.Path = &StringCriterionInput{
			Modifier: CriterionModifierEquals,
			Value:    p + "%",
		}
	}

	return ret
}

type FileQueryResult struct {
	QueryResult[FileID]
	TotalDuration float64
	Megapixels    float64
	TotalSize     int64

	getter     FileGetter
	files      []File
	resolveErr error
}

func NewFileQueryResult(fileGetter FileGetter) *FileQueryResult {
	return &FileQueryResult{
		getter: fileGetter,
	}
}

func (r *FileQueryResult) Resolve(ctx context.Context) ([]File, error) {
	// cache results
	if r.files == nil && r.resolveErr == nil {
		r.files, r.resolveErr = r.getter.Find(ctx, r.IDs...)
	}
	return r.files, r.resolveErr
}
