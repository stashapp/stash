package models

import (
	"io"
	"io/fs"
)

// FileOpener provides an interface to open a file.
type FileOpener interface {
	Open() (io.ReadCloser, error)
}

// FS represents a file system.
type FS interface {
	Stat(name string) (fs.FileInfo, error)
	Lstat(name string) (fs.FileInfo, error)
	Open(name string) (fs.ReadDirFile, error)
	OpenZip(name string, size int64) (ZipFS, error)
	IsPathCaseSensitive(path string) (bool, error)
}

// ZipFS represents a zip file system.
type ZipFS interface {
	FS
	io.Closer
	OpenOnly(name string) (io.ReadCloser, error)
}
