package file

import (
	"context"
	"io/fs"

	"github.com/stashapp/stash/pkg/models"
)

// PathFilter provides a filter function for paths.
type PathFilter interface {
	Accept(ctx context.Context, path string, info fs.FileInfo) bool
}

type PathFilterFunc func(path string) bool

func (pff PathFilterFunc) Accept(path string) bool {
	return pff(path)
}

// Filter provides a filter function for Files.
type Filter interface {
	Accept(ctx context.Context, f models.File) bool
}

type FilterFunc func(ctx context.Context, f models.File) bool

func (ff FilterFunc) Accept(ctx context.Context, f models.File) bool {
	return ff(ctx, f)
}

// Handler provides a handler for Files.
type Handler interface {
	Handle(ctx context.Context, f models.File, oldFile models.File) error
}

// FilteredHandler is a Handler runs only if the filter accepts the file.
type FilteredHandler struct {
	Handler
	Filter
}

// Handle runs the handler if the filter accepts the file.
func (h *FilteredHandler) Handle(ctx context.Context, f models.File, oldFile models.File) error {
	if h.Accept(ctx, f) {
		return h.Handler.Handle(ctx, f, oldFile)
	}
	return nil
}

// CleanHandler provides a handler for cleaning Files and Folders.
type CleanHandler interface {
	HandleFile(ctx context.Context, fileDeleter *Deleter, fileID models.FileID) error
	HandleFolder(ctx context.Context, fileDeleter *Deleter, folderID models.FolderID) error
}
