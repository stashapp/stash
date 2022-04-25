package file

import "context"

type Filter interface {
	Accept(f File) bool
}

type Handler interface {
	Handle(ctx context.Context, fs FS, f File) error
}

type FilteredHandler struct {
	Handler
	Filter
}

func (h *FilteredHandler) Handle(ctx context.Context, fs FS, f File) error {
	if h.Accept(f) {
		return h.Handler.Handle(ctx, fs, f)
	}
	return nil
}
