package models

type FileQueryOptions struct {
	QueryOptions
}

type FileQueryResult struct {
	QueryResult
	TotalDuration float64
	TotalSize     float64

	finder     FileFinder
	files      []*File
	resolveErr error
}

func NewFileQueryResult(finder FileFinder) *FileQueryResult {
	return &FileQueryResult{
		finder: finder,
	}
}

func (r *FileQueryResult) Resolve() ([]*File, error) {
	// cache results
	if r.files == nil && r.resolveErr == nil {
		r.files, r.resolveErr = r.finder.Find(r.IDs)
	}
	return r.files, r.resolveErr
}

type FileFinder interface {
	Find(ids []int) ([]*File, error)
}

type FileReader interface {
	FileFinder
	FindByChecksum(checksum string) ([]*File, error)
	FindByOSHash(oshash string) ([]*File, error)
	FindByPath(path string) (*File, error)
	Query(options FileQueryOptions) (*FileQueryResult, error)
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
