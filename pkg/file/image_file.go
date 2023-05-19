package file

// ImageFile is an extension of BaseFile to represent image files.
type ImageFile struct {
	*BaseFile
	Format string `json:"format"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func (f ImageFile) GetWidth() int {
	return f.Width
}

func (f ImageFile) GetHeight() int {
	return f.Height
}

func (f ImageFile) GetFormat() string {
	return f.Format
}
