package file

// VisualFile is an interface for files that have a width and height.
type VisualFile interface {
	File
	GetWidth() int
	GetHeight() int
	GetFormat() string
}

func GetMinResolution(f VisualFile) int {
	w := f.GetWidth()
	h := f.GetHeight()

	if w < h {
		return w
	}

	return h
}
