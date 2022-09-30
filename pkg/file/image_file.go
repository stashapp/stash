package file

// ImageFile is an extension of BaseFile to represent image files.
type ImageFile struct {
	*BaseFile
	Format string `json:"format"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
