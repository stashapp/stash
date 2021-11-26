package image

import (
	"archive/zip"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
	_ "golang.org/x/image/webp"
)

func DecodeSourceImage(r io.Reader) (*image.Config, *string, error) {
	config, format, err := image.DecodeConfig(r)

	return &config, &format, err
}

func FileExists(path string) bool {
	f, err := openSourceImage(path)
	if err != nil {
		return false
	}
	defer f.Close()

	return true
}

type imageReadCloser struct {
	src io.ReadCloser
	zrc *zip.ReadCloser
}

func (i *imageReadCloser) Read(p []byte) (n int, err error) {
	return i.src.Read(p)
}

func (i *imageReadCloser) Close() error {
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

func openSourceImage(path string) (io.ReadCloser, error) {
	// may need to read from a zip file
	zipFilename, filename := file.ZipFilePath(path)
	if zipFilename != "" {
		r, err := zip.OpenReader(zipFilename)
		if err != nil {
			return nil, err
		}

		// defer closing of zip to the calling function, unless an error
		// is returned, in which case it should be closed immediately

		// find the file matching the filename
		for _, f := range r.File {
			if f.Name == filename {
				src, err := f.Open()
				if err != nil {
					r.Close()
					return nil, err
				}
				return &imageReadCloser{
					src: src,
					zrc: r,
				}, nil
			}
		}

		r.Close()
		return nil, fmt.Errorf("file with name '%s' not found in zip file '%s'", filename, zipFilename)
	}

	return os.Open(filename)
}

// GetFileModTime gets the file modification time, handling files in zip files.
func GetFileModTime(path string) (time.Time, error) {
	fi, err := stat(path)
	if err != nil {
		return time.Time{}, fmt.Errorf("error performing stat on %s: %s", path, err.Error())
	}

	ret := fi.ModTime()
	// truncate to seconds, since we don't store beyond that in the database
	ret = ret.Truncate(time.Second)

	return ret, nil
}

func stat(path string) (os.FileInfo, error) {
	// may need to read from a zip file
	zipFilename, filename := file.ZipFilePath(path)
	if zipFilename != "" {
		r, err := zip.OpenReader(zipFilename)
		if err != nil {
			return nil, err
		}
		defer r.Close()

		// find the file matching the filename
		for _, f := range r.File {
			if f.Name == filename {
				return f.FileInfo(), nil
			}
		}

		return nil, fmt.Errorf("file with name '%s' not found in zip file '%s'", filename, zipFilename)
	}

	return os.Stat(filename)
}

func IsCover(img *models.Image) bool {
	return strings.HasSuffix(strings.ToLower(img.Path), "cover.jpg")
}

func GetTitle(s *models.Image) string {
	if s.Title.String != "" {
		return s.Title.String
	}

	return filepath.Base(s.Path)
}

// GetFilename gets the base name of the image file
// If stripExt is set the file extension is omitted from the name
func GetFilename(s *models.Image, stripExt bool) string {
	return utils.GetNameFromPath(s.Path, stripExt)
}
