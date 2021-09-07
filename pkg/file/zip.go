package file

import (
	"archive/zip"
	"io"
	"io/fs"

	"github.com/stashapp/stash/pkg/image"
)

type zipFile struct {
	zipPath string
	zf      *zip.File
}

func (f *zipFile) Open() (io.ReadCloser, error) {
	return f.zf.Open()
}

func (f *zipFile) Path() string {
	// TODO - fix this
	return image.ZipFilename(f.zipPath, f.zf.Name)
}

func (f *zipFile) FileInfo() fs.FileInfo {
	return f.zf.FileInfo()
}

func ZipFile(zipPath string, zf *zip.File) SourceFile {
	return &zipFile{
		zipPath: zipPath,
		zf:      zf,
	}
}
