package file

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// ZipFS is a file system backed by a zip file.
type ZipFS struct {
	*zip.Reader
	zipInfo fs.FileInfo
	zipPath string
}

// zipDirEntry is a special FileInfo that returns the zip file as a directory.
type zipDirEntry struct {
	fs.FileInfo
}

func (d *zipDirEntry) IsDir() bool { return true }

func (f *ZipFS) rel(name string) (string, error) {
	if f.zipPath == name {
		return ".", nil
	}

	relName, err := filepath.Rel(f.zipPath, name)
	if err != nil {
		return "", fmt.Errorf("internal error getting relative path: %w", err)
	}

	return relName, nil
}

func (f *ZipFS) Lstat(name string) (fs.FileInfo, error) {
	if f.zipPath == name {
		return &zipDirEntry{
			FileInfo: f.zipInfo,
		}, nil
	}

	relName, err := f.rel(name)
	if err != nil {
		return nil, err
	}

	for _, ff := range f.File {
		if ff.Name == relName {
			return ff.FileInfo(), nil
		}
	}

	return nil, os.ErrNotExist
}

type zipReadDirFile struct {
	fs.File
}

func (f *zipReadDirFile) ReadDir(n int) ([]fs.DirEntry, error) {
	asReadDirFile, _ := f.File.(fs.ReadDirFile)
	if asReadDirFile == nil {
		return nil, fmt.Errorf("internal error: not a ReadDirFile")
	}

	return asReadDirFile.ReadDir(n)
}

func (f *ZipFS) Open(name string) (fs.ReadDirFile, error) {
	relName, err := f.rel(name)
	if err != nil {
		return nil, err
	}

	r, err := f.Reader.Open(relName)
	if err != nil {
		return nil, err
	}

	return &zipReadDirFile{
		File: r,
	}, nil
}
