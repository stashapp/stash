package image

import (
	"context"
	"fmt"
	"image"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/stashapp/stash/pkg/file"
	_ "golang.org/x/image/webp"
)

// Decorator adds image specific fields to a File.
type Decorator struct {
}

func (d *Decorator) Decorate(ctx context.Context, fs file.FS, f file.File) (file.File, error) {
	base := f.Base()
	r, err := fs.Open(base.Path)
	if err != nil {
		return f, fmt.Errorf("reading image file %q: %w", base.Path, err)
	}
	defer r.Close()

	c, format, err := image.DecodeConfig(r)
	if err != nil {
		return f, fmt.Errorf("decoding image file %q: %w", base.Path, err)
	}

	return &file.ImageFile{
		BaseFile: base,
		Format:   format,
		Width:    c.Width,
		Height:   c.Height,
	}, nil
}
