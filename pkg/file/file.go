package file

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
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

func (f *fsFile) ZipFile() *models.File {
	return nil
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

	zipPath := zipFiles[0].Path
	readCloser, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	// strip zip path from file path
	path := strings.TrimPrefix(f.Path, zipPath+string(filepath.Separator))

	zf, err := readCloser.Open(path)
	if err != nil {
		return nil, err
	}
	defer zf.Close()

	return zf.Stat()
}

// Open opens the provided file for reading. Assumes that the zip path is
// populated if applicable.
func Open(f *models.File) (io.ReadCloser, error) {
	if !f.ZipFileID.Valid || !f.ZipPath.Valid {
		return os.Open(f.Path)
	}

	zipReader, err := zip.OpenReader(f.ZipPath.String)
	if err != nil {
		return nil, err
	}

	// strip zip path from file path
	path := strings.TrimPrefix(f.Path, f.ZipPath.String+string(filepath.Separator))

	src, err := zipReader.Open(path)
	if err != nil {
		zipReader.Close()
		return nil, err
	}

	return &zipFileReadCloser{
		src: src,
		zrc: zipReader,
	}, nil
}

func Info(f *models.File) (fs.FileInfo, error) {
	if !f.ZipFileID.Valid || !f.ZipPath.Valid {
		return os.Stat(f.Path)
	}

	zipReader, err := zip.OpenReader(f.ZipPath.String)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	// strip zip path from file path
	path := strings.TrimPrefix(f.Path, f.ZipPath.String+string(filepath.Separator))
	for _, zf := range zipReader.File {
		if zf.Name == path {
			return zf.FileInfo(), nil
		}
	}

	return nil, fs.ErrNotExist
}

// Serve serves the provided file.
func Serve(w http.ResponseWriter, r *http.Request, f *models.File) {
	w.Header().Add("Cache-Control", "max-age=604800000") // 1 Week
	if !f.ZipFileID.Valid || !f.ZipPath.Valid {
		http.ServeFile(w, r, f.Path)
	} else {
		rc, err := Open(f)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			logger.Warnf("error reading image in zip %q: %v", f.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer rc.Close()

		data, err := io.ReadAll(rc)
		if err != nil {
			logger.Warnf("error reading image in zip %q: %v", f.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if k, err := w.Write(data); err != nil {
			logger.Warnf("failure while serving image (wrote %v bytes out of %v): %v", k, len(data), err)
		}
	}
}
