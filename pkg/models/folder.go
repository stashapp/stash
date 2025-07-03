package models

import (
	"context"
	"path/filepath"
	"strings"
)

type FolderQueryOptions struct {
	QueryOptions
	FolderFilter *FolderFilterType

	TotalDuration bool
	Megapixels    bool
	TotalSize     bool
}

type FolderFilterType struct {
	OperatorFilter[FolderFilterType]

	Path     *StringCriterionInput `json:"path,omitempty"`
	Basename *StringCriterionInput `json:"basename,omitempty"`
	// Filter by parent directory path
	Dir          *StringCriterionInput            `json:"dir,omitempty"`
	ParentFolder *HierarchicalMultiCriterionInput `json:"parent_folder,omitempty"`
	ZipFile      *MultiCriterionInput             `json:"zip_file,omitempty"`
	// Filter by modification time
	ModTime      *TimestampCriterionInput `json:"mod_time,omitempty"`
	GalleryCount *IntCriterionInput       `json:"gallery_count,omitempty"`
	// Filter by files that meet this criteria
	FilesFilter *FileFilterType `json:"files_filter,omitempty"`
	// Filter by related galleries that meet this criteria
	GalleriesFilter *GalleryFilterType `json:"galleries_filter,omitempty"`
	// Filter by creation time
	CreatedAt *TimestampCriterionInput `json:"created_at,omitempty"`
	// Filter by last update time
	UpdatedAt *TimestampCriterionInput `json:"updated_at,omitempty"`
}

func PathsFolderFilter(paths []string) *FileFilterType {
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

type FolderQueryResult struct {
	QueryResult[FolderID]

	getter     FolderGetter
	folders    []*Folder
	resolveErr error
}

func NewFolderQueryResult(folderGetter FolderGetter) *FolderQueryResult {
	return &FolderQueryResult{
		getter: folderGetter,
	}
}

func (r *FolderQueryResult) Resolve(ctx context.Context) ([]*Folder, error) {
	// cache results
	if r.folders == nil && r.resolveErr == nil {
		r.folders, r.resolveErr = r.getter.FindMany(ctx, r.IDs)
	}
	return r.folders, r.resolveErr
}
