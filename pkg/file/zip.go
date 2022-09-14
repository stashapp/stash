package file

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
)

var (
	errNotReaderAt  = errors.New("not a ReaderAt")
	errZipFSOpenZip = errors.New("cannot open zip file inside zip file")
)

// ZipFS is a file system backed by a zip file.
type ZipFS struct {
	*zip.Reader
	zipFileCloser io.Closer
	zipInfo       fs.FileInfo
	zipPath       string
}

func newZipFS(fs FS, path string, info fs.FileInfo) (*ZipFS, error) {
	reader, err := fs.Open(path)
	if err != nil {
		return nil, err
	}

	asReaderAt, _ := reader.(io.ReaderAt)
	if asReaderAt == nil {
		reader.Close()
		return nil, errNotReaderAt
	}

	zipReader, err := zip.NewReader(asReaderAt, info.Size())
	if err != nil {
		reader.Close()
		return nil, err
	}

	return &ZipFS{
		Reader:        zipReader,
		zipFileCloser: reader,
		zipInfo:       info,
		zipPath:       path,
	}, nil
}

func (f *ZipFS) rel(name string) (string, error) {
	if f.zipPath == name {
		return ".", nil
	}

	relName, err := filepath.Rel(f.zipPath, name)
	if err != nil {
		return "", fmt.Errorf("internal error getting relative path: %w", err)
	}

	// convert relName to use slash, since zip files do so regardless
	// of os
	relName = filepath.ToSlash(relName)

	return relName, nil
}

func (f *ZipFS) Lstat(name string) (fs.FileInfo, error) {
	reader, err := f.Open(name)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return reader.Stat()
}

func (f *ZipFS) OpenZip(name string) (*ZipFS, error) {
	return nil, errZipFSOpenZip
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

func (f *ZipFS) Close() error {
	return f.zipFileCloser.Close()
}

// openOnly returns a ReadCloser where calling Close will close the zip fs as well.
func (f *ZipFS) OpenOnly(name string) (io.ReadCloser, error) {
	r, err := f.Open(name)
	if err != nil {
		return nil, err
	}

	return &wrappedReadCloser{
		ReadCloser: r,
		outer:      f,
	}, nil
}

type wrappedReadCloser struct {
	io.ReadCloser
	outer io.Closer
}

func (f *wrappedReadCloser) Close() error {
	_ = f.ReadCloser.Close()
	return f.outer.Close()
}
