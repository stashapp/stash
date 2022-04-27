package file

import "context"

// Filter provides a filter function for Files.
type Filter interface {
	Accept(f File) bool
}

type FilterFunc func(f File) bool

func (ff FilterFunc) Accept(f File) bool {
	return ff(f)
}

// Handler provides a handler for Files.
type Handler interface {
	Handle(ctx context.Context, fs FS, f File) error
}

// FilteredHandler is a Handler runs only if the filter accepts the file.
type FilteredHandler struct {
	Handler
	Filter
}

// Handle runs the handler if the filter accepts the file.
func (h *FilteredHandler) Handle(ctx context.Context, fs FS, f File) error {
	if h.Accept(f) {
		return h.Handler.Handle(ctx, fs, f)
	}
	return nil
}
