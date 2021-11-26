package file

import (
	"archive/zip"
	"io"
	"io/fs"
	"strings"

	"github.com/stashapp/stash/pkg/models"
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

func (f *zipFile) ZipFileID() int {
	return f.zipFile.ID
}

func ZipFile(zf *models.File, file *zip.File) SourceFile {
	return &zipFile{
		zipFile: zf,
		file:    file,
	}
}

func ZipFilename(zipFilename, filenameInZip string) string {
	return zipFilename + zipSeparator + filenameInZip
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
