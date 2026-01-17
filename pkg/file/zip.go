package file

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/xWTF/chardet"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

var (
	ErrNotReaderAt  = errors.New("invalid reader: does not implement io.ReaderAt")
	errZipFSOpenZip = errors.New("cannot open zip file inside zip file")
)

// ZipFS is a file system backed by a zip file.
type zipFS struct {
	*zip.Reader
	zipFileCloser io.Closer
	zipPath       string
}

func newZipFS(fs models.FS, path string, size int64) (*zipFS, error) {
	reader, err := fs.Open(path)
	if err != nil {
		return nil, err
	}

	asReaderAt, _ := reader.(io.ReaderAt)
	if asReaderAt == nil {
		reader.Close()
		return nil, ErrNotReaderAt
	}

	zipReader, err := zip.NewReader(asReaderAt, size)
	if err != nil {
		reader.Close()
		return nil, err
	}

	// Concat all Name and Comment for better detection result
	var buffer bytes.Buffer
	for _, f := range zipReader.File {
		buffer.WriteString(f.Name)
		buffer.WriteString(f.Comment)
	}
	buffer.WriteString(zipReader.Comment)

	// Detect encoding
	d, err := chardet.NewTextDetector().DetectBest(buffer.Bytes())
	if err != nil {
		// If we can't detect the encoding, just assume it's UTF8
		logger.Warnf("Unable to detect decoding for %s: %w", path, err)
	}

	// If the charset is not UTF8, decode'em
	if d != nil && d.Charset != "UTF-8" {
		logger.Debugf("Detected non-utf8 zip charset %s (%s): %s", d.Charset, d.Language, path)

		e, _ := charset.Lookup(d.Charset)
		if e == nil {
			// if we can't find the encoding, just assume it's UTF8
			logger.Warnf("Failed to lookup charset %s, language %s", d.Charset, d.Language)
		} else {
			decoder := e.NewDecoder()
			for _, f := range zipReader.File {
				newName, _, err := transform.String(decoder, f.Name)
				if err != nil {
					reader.Close()
					logger.Warnf("Failed to decode %v: %v", []byte(f.Name), err)
				} else {
					f.Name = newName
				}
				// Comments are not decoded cuz stash doesn't use that
			}
		}
	}

	return &zipFS{
		Reader:        zipReader,
		zipFileCloser: reader,
		zipPath:       path,
	}, nil
}

func (f *zipFS) rel(name string) (string, error) {
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

func (f *zipFS) Stat(name string) (fs.FileInfo, error) {
	reader, err := f.Open(name)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return reader.Stat()
}

func (f *zipFS) Lstat(name string) (fs.FileInfo, error) {
	return f.Stat(name)
}

func (f *zipFS) OpenZip(name string, size int64) (models.ZipFS, error) {
	return nil, errZipFSOpenZip
}

func (f *zipFS) IsPathCaseSensitive(path string) (bool, error) {
	return true, nil
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

func (f *zipFS) Open(name string) (fs.ReadDirFile, error) {
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

func (f *zipFS) Close() error {
	return f.zipFileCloser.Close()
}

// openOnly returns a ReadCloser where calling Close will close the zip fs as well.
func (f *zipFS) OpenOnly(name string) (io.ReadCloser, error) {
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
