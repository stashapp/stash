package file

import (
	"io/fs"
	"os"
)

type FS interface {
	Lstat(name string) (fs.FileInfo, error)
	Open(name string) (fs.ReadDirFile, error)
}

type OsFS struct{}

func (f *OsFS) Lstat(name string) (fs.FileInfo, error) {
	return os.Lstat(name)
}

func (f *OsFS) Open(name string) (fs.ReadDirFile, error) {
	return os.Open(name)
}
