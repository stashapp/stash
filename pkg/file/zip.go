package file

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

const zipSeparator = "\x00"

type zipFile struct {
	zipFile *models.File
	file    *zip.File
}

func (f *zipFile) Open() (io.ReadCloser, error) {
	return f.file.Open()
}

func (f *zipFile) Path() string {
	return f.file.Name
}

func (f *zipFile) FileInfo() fs.FileInfo {
	return f.file.FileInfo()
}

func (f *zipFile) ZipFile() *models.File {
	return f.zipFile
}

func ZipFile(zf *models.File, file *zip.File) SourceFile {
	return &zipFile{
		zipFile: zf,
		file:    file,
	}
}

// IsZipPath returns true if the path includes the zip separator byte,
// indicating it is within a zip file.
func IsZipPath(p string) bool {
	return strings.Contains(p, zipSeparator)
}

// ZipPathDisplayName converts an zip path for display. It translates the zip
// file separator character into '/', since this character is also used for
// path separators within zip files. It returns the original provided path
// if it does not contain the zip file separator character.
func ZipPathDisplayName(path string) string {
	return strings.ReplaceAll(path, zipSeparator, "/")
}

func ZipFilePath(path string) (zipFilename, filename string) {
	nullIndex := strings.Index(path, zipSeparator)
	if nullIndex != -1 {
		zipFilename = path[0:nullIndex]
		filename = path[nullIndex+1:]
	} else {
		filename = path
	}
	return
}

type zipFileReadCloser struct {
	src io.ReadCloser
	zrc *zip.ReadCloser
}

func (i *zipFileReadCloser) Read(p []byte) (n int, err error) {
	return i.src.Read(p)
}

func (i *zipFileReadCloser) Close() error {
	err := i.src.Close()
	var err2 error
	if i.zrc != nil {
		err2 = i.zrc.Close()
	}

	if err != nil {
		return err
	}
	return err2
}

func WalkZip(path string, walkFunc func(file *zip.File) error) error {
	readCloser, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer readCloser.Close()

	for _, file := range readCloser.File {
		if file.FileInfo().IsDir() {
			continue
		}

		err := walkFunc(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func CountImagesInZip(path string, validExtensions []string) (int, error) {
	ret := 0
	err := WalkZip(path, func(file *zip.File) error {
		if strings.Contains(file.Name, "__MACOSX") {
			return nil
		}

		if !utils.MatchExtension(file.Name, validExtensions) {
			return nil
		}

		ret++
		return nil
	})

	if err != nil {
		return ret, fmt.Errorf("walking zip: %w", err)
	}

	return ret, nil
}
