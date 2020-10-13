package image

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

		// find the file matching the filename
		for _, f := range r.File {
			if f.Name == filename {
				src, err := f.Open()
				if err != nil {
					return nil, err
				}
				return &imageReadCloser{
					src: src,
					zrc: r,
				}, nil
			}
		}

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

func SetFileDetails(i *models.Image) error {
	f, err := stat(i.Path)
	if err != nil {
		return err
	}

	src, _ := GetSourceImage(i)

	if src != nil {
		i.Width = sql.NullInt64{
			Int64: int64(src.Bounds().Max.X),
			Valid: true,
		}
		i.Height = sql.NullInt64{
			Int64: int64(src.Bounds().Max.Y),
			Valid: true,
		}
	}

	i.Size = sql.NullInt64{
		Int64: int64(f.Size()),
		Valid: true,
	}

	return nil
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

		data, err := ioutil.ReadAll(rc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(data)
	}
}

func IsCover(img *models.Image) bool {
	_, fn := getFilePath(img.Path)
	return fn == "cover.jpg"
}

func GetTitle(s *models.Image) string {
	if s.Title.String != "" {
		return s.Title.String
	}

	_, fn := getFilePath(s.Path)
	return filepath.Base(fn)
}
