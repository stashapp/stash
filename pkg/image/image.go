package image

import (
	"net/http"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	_ "golang.org/x/image/webp"
)

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
