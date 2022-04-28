package image

import (
	"io"
	"net/http"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	_ "golang.org/x/image/webp"
)

func FileExists(path string) bool {
	f, err := openSourceImage(path)
	if err != nil {
		return false
	}
	defer f.Close()

	return true
}

func openSourceImage(path string) (io.ReadCloser, error) {
	// // may need to read from a zip file
	// zipFilename, filename := file.ZipFilePath(path)
	// if zipFilename != "" {
	// 	r, err := zip.OpenReader(zipFilename)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	// defer closing of zip to the calling function, unless an error
	// 	// is returned, in which case it should be closed immediately

	// 	// find the file matching the filename
	// 	for _, f := range r.File {
	// 		if f.Name == filename {
	// 			src, err := f.Open()
	// 			if err != nil {
	// 				r.Close()
	// 				return nil, err
	// 			}
	// 			return &imageReadCloser{
	// 				src: src,
	// 				zrc: r,
	// 			}, nil
	// 		}
	// 	}

	// 	r.Close()
	// 	return nil, fmt.Errorf("file with name '%s' not found in zip file '%s'", filename, zipFilename)
	// }

	// return os.Open(filename)
	panic("TODO")
}

func Serve(w http.ResponseWriter, r *http.Request, path string) {
	// zipFilename, _ := file.ZipFilePath(path)
	// w.Header().Add("Cache-Control", "max-age=604800000") // 1 Week
	// if zipFilename == "" {
	// 	http.ServeFile(w, r, path)
	// } else {
	// 	rc, err := openSourceImage(path)
	// 	if err != nil {
	// 		// assume not found
	// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	// 		return
	// 	}
	// 	defer rc.Close()

	// 	data, err := io.ReadAll(rc)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	if k, err := w.Write(data); err != nil {
	// 		logger.Warnf("failure while serving image (wrote %v bytes out of %v): %v", k, len(data), err)
	// 	}
	// }
}

func IsCover(img *models.Image) bool {
	return strings.HasSuffix(img.Path(), "cover.jpg")
}
