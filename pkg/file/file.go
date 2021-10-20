package file

import (
	"io"
	"io/fs"
	"os"
)

type fsFile struct {
	path string
	info fs.FileInfo
}

func (f *fsFile) Open() (io.ReadCloser, error) {
	return os.Open(f.path)
}

func (f *fsFile) Path() string {
	return f.path
}

func (f *fsFile) FileInfo() fs.FileInfo {
	return f.info
}

func FSFile(path string, info fs.FileInfo) SourceFile {
	return &fsFile{
		path: path,
		info: info,
	}
}
