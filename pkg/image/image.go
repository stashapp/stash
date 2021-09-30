package image

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
	_ "golang.org/x/image/webp"
)

const zipSeparator = "\x00"

func GetSourceImage(i *models.Image) (image.Image, error) {
	f, err := openSourceImage(i.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	srcImage, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return srcImage, nil
}

func DecodeSourceImage(i *models.Image) (*image.Config, *string, error) {
	f, err := openSourceImage(i.Path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	config, format, err := image.DecodeConfig(f)

	return &config, &format, err
}

func CalculateMD5(path string) (string, error) {
	f, err := openSourceImage(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return utils.MD5FromReader(f)
}

func FileExists(path string) bool {
	f, err := openSourceImage(path)
	if err != nil {
		return false
	}
	defer f.Close()

	return true
}

func ZipFilename(zipFilename, filenameInZip string) string {
	return zipFilename + zipSeparator + filenameInZip
}

// IsZipPath returns true if the path includes the zip separator byte,
// indicating it is within a zip file.
// TODO - this should be moved to utils
func IsZipPath(p string) bool {
	return strings.Contains(p, zipSeparator)
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
	zipFilename, filename := getFilePath(path)
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

func getFilePath(path string) (zipFilename, filename string) {
	nullIndex := strings.Index(path, zipSeparator)
	if nullIndex != -1 {
		zipFilename = path[0:nullIndex]
		filename = path[nullIndex+1:]
	} else {
		filename = path
	}
	return
}

// GetFileDetails returns a pointer to an Image object with the
// width, height and size populated.
func GetFileDetails(path string) (*models.Image, error) {
	i := &models.Image{
		Path: path,
	}

	err := SetFileDetails(i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func SetFileDetails(i *models.Image) error {
	f, err := stat(i.Path)
	if err != nil {
		return err
	}

	config, _, err := DecodeSourceImage(i)

	if err == nil {
		i.Width = sql.NullInt64{
			Int64: int64(config.Width),
			Valid: true,
		}
		i.Height = sql.NullInt64{
			Int64: int64(config.Height),
			Valid: true,
		}
	}

	i.Size = sql.NullInt64{
		Int64: int64(f.Size()),
		Valid: true,
	}

	return nil
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
	zipFilename, filename := getFilePath(path)
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

// PathDisplayName converts an image path for display. It translates the zip
// file separator character into '/', since this character is also used for
// path separators within zip files. It returns the original provided path
// if it does not contain the zip file separator character.
func PathDisplayName(path string) string {
	return strings.Replace(path, zipSeparator, "/", -1)
}

func Serve(w http.ResponseWriter, r *http.Request, path string) {
	zipFilename, _ := getFilePath(path)
	w.Header().Add("Cache-Control", "max-age=604800000") // 1 Week
	if zipFilename == "" {
		http.ServeFile(w, r, path)
	} else {
		rc, err := openSourceImage(path)
		if err != nil {
			// assume not found
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		defer rc.Close()

		data, err := io.ReadAll(rc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if k, err := w.Write(data); err != nil {
			logger.Warnf("failure while serving image (wrote %v bytes out of %v): %v", k, len(data), err)
		}
	}
}

func IsCover(img *models.Image) bool {
	_, fn := getFilePath(img.Path)
	return strings.HasSuffix(fn, "cover.jpg")
}

func GetTitle(s *models.Image) string {
	if s.Title.String != "" {
		return s.Title.String
	}

	_, fn := getFilePath(s.Path)
	return filepath.Base(fn)
}

// GetFilename gets the base name of the image file
// If stripExt is set the file extension is omitted from the name
func GetFilename(s *models.Image, stripExt bool) string {
	_, fn := getFilePath(s.Path)
	return utils.GetNameFromPath(fn, stripExt)
}
