package models

type FileReader interface {
	Find(ids []int) ([]*File, error)
	FindByChecksum(checksum string) ([]*File, error)
	FindByOSHash(oshash string) ([]*File, error)
	FindByPath(path string, zipFileID int) (*File, error)
	// AllOfType(fileType FileType) ([]*File, error)
}

type FileWriter interface {
	Create(newFile File) (*File, error)
	UpdateFull(updatedFile File) (*File, error)
	Destroy(id int) error
}

type FileReaderWriter interface {
	FileReader
	FileWriter
}

type FileJoinReader interface {
	GetFileIDs(id int) ([]int, error)
}

type FileJoinWriter interface {
	UpdateFiles(id int, fileIDs []int) error
}

type FileJoinReaderWriter interface {
	FileJoinReader
	FileJoinWriter
}
