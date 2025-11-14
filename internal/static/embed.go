// Package static provides the static files embedded in the application.
package static

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
)

//go:embed performer performer_male performer_sfw scene image gallery tag studio group
var data embed.FS

const (
	Performer                = "performer"
	PerformerMale            = "performer_male"
	DefaultSFWPerformerImage = "performer_sfw/performer.svg"

	Scene             = "scene"
	DefaultSceneImage = "scene/scene.svg"

	Image             = "image"
	DefaultImageImage = "image/image.svg"

	Gallery             = "gallery"
	DefaultGalleryImage = "gallery/gallery.svg"

	Tag             = "tag"
	DefaultTagImage = "tag/tag.svg"

	Studio             = "studio"
	DefaultStudioImage = "studio/studio.svg"

	Group             = "group"
	DefaultGroupImage = "group/group.svg"
)

// Sub returns an FS rooted at path, using fs.Sub.
// It will panic if an error occurs.
func Sub(path string) fs.FS {
	ret, err := fs.Sub(data, path)
	if err != nil {
		panic(fmt.Sprintf("creating static SubFS: %v", err))
	}
	return ret
}

// Open opens the file at path for reading.
// It will panic if an error occurs.
func Open(path string) fs.File {
	f, err := data.Open(path)
	if err != nil {
		panic(fmt.Sprintf("opening static file: %v", err))
	}
	return f
}

// ReadAll returns the contents of the file at path.
// It will panic if an error occurs.
func ReadAll(path string) []byte {
	f := Open(path)
	ret, err := io.ReadAll(f)
	if err != nil {
		panic(fmt.Sprintf("reading static file: %v", err))
	}
	return ret
}
