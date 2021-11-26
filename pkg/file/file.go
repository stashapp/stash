package file

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"

	"github.com/stashapp/stash/pkg/models"
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

func (f *fsFile) ZipFileID() int {
	return 0
}

func FSFile(path string, info fs.FileInfo) SourceFile {
	return &fsFile{
		path: path,
		info: info,
	}
}

type FSStatter struct {
}

func (s *FSStatter) Stat(reader models.FileReader, f models.File) (fs.FileInfo, error) {
	if !f.ZipFileID.Valid {
		// direct file system
		return os.Stat(f.Path)
	}

	zipFiles, err := reader.Find([]int{int(f.ZipFileID.Int64)})
	if err != nil {
		return nil, err
	}

	readCloser, err := zip.OpenReader(zipFiles[0].Path)
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	zf, err := readCloser.Open(f.Path)
	if err != nil {
		return nil, err
	}

	return zf.Stat()
}
