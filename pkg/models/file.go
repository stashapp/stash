package models

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"
)

type FileQueryOptions struct {
	QueryOptions
	FileFilter *FileFilterType
}

type FileFilterType struct {
	And *FileFilterType `json:"AND"`
	Or  *FileFilterType `json:"OR"`
	Not *FileFilterType `json:"NOT"`

	// Filter by path
	Path *StringCriterionInput `json:"path"`
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
	// can't use QueryResult because id type is wrong

	IDs   []FileID
	Count int

	finder     FileFinder
	files      []File
	resolveErr error
}

func NewFileQueryResult(finder FileFinder) *FileQueryResult {
	return &FileQueryResult{
		finder: finder,
	}
}

func (r *FileQueryResult) Resolve(ctx context.Context) ([]File, error) {
	// cache results
	if r.files == nil && r.resolveErr == nil {
		r.files, r.resolveErr = r.finder.Find(ctx, r.IDs...)
	}
	return r.files, r.resolveErr
}

type FileFinder interface {
	Find(ctx context.Context, id ...FileID) ([]File, error)
}

type FileReader interface {
	FileFinder
	FindByPath(ctx context.Context, path string) (File, error)
	FindAllByPath(ctx context.Context, path string) ([]File, error)
	FindAllInPaths(ctx context.Context, p []string, limit, offset int) ([]File, error)
	FindByFingerprint(ctx context.Context, fp Fingerprint) ([]File, error)
	FindByZipFileID(ctx context.Context, zipFileID FileID) ([]File, error)
	FindByFileInfo(ctx context.Context, info fs.FileInfo, size int64) ([]File, error)
	Query(ctx context.Context, options FileQueryOptions) (*FileQueryResult, error)

	CountAllInPaths(ctx context.Context, p []string) (int, error)
	CountByFolderID(ctx context.Context, folderID FolderID) (int, error)

	GetCaptions(ctx context.Context, fileID FileID) ([]*VideoCaption, error)
	IsPrimary(ctx context.Context, fileID FileID) (bool, error)
}

type FileWriter interface {
	Create(ctx context.Context, f File) error
	Update(ctx context.Context, f File) error
	Destroy(ctx context.Context, id FileID) error

	UpdateCaptions(ctx context.Context, fileID FileID, captions []*VideoCaption) error
}

type FileReaderWriter interface {
	FileReader
	FileWriter
}
