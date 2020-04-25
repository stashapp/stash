package api

import (
	"context"
	"github.com/dustin/go-humanize"
	"github.com/go-chi/chi"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type galleryRoutes struct{}

type thumbBuffer struct {
	path string
	dir  string
	data []byte
}

const ThumbCacheLimit int64 = 25 * 1024 * 1024 // cache size limit in bytes // TODO add to options instead

func newCacheThumb(dir string, path string, data []byte) *thumbBuffer {
	t := thumbBuffer{dir: dir, path: path, data: data}
	return &t
}

var writeChan chan *thumbBuffer
var touchChan chan *string

func startThumbCache() { // TODO add extra wait, close chan code if/when stash gets a stop mode

	dir := paths.GetThumbCache()
	info, err := os.Lstat(dir)
	if err == nil {
		var files []utils.DuDetails
		size := utils.DuDir(dir, info, &files) // get cache stats
		logger.Infof("Cache dir used for galleries: %s", dir)
		logger.Infof("Current cache size: %s", humanize.Bytes(uint64(size)))
		if ThumbCacheLimit > 0 { // ThumbCachelimit = 0 means limit disabled
			diff := size - ThumbCacheLimit
			if diff > 0 { // reduce cache by deleting files if needed
				utils.SortDuDetailsByMtime(files) // sort slice so least recently modified files are first
				diff = utils.ReduceDir(files, diff)
				logger.Infof("Reduced cache by %s", humanize.Bytes(uint64(diff)))
			}
		}
	}

	writeChan = make(chan *thumbBuffer, 20)
	touchChan = make(chan *string, 100)
	go thumbnailCacheWriter()
	go thumbnailCacheToucher()
}

//serialize file writes to avoid race conditions
func thumbnailCacheWriter() {

	for thumb := range writeChan {

		exists, _ := utils.FileExists(thumb.path)
		if !exists { // file to write shouldn't exist
			utils.EnsureDirAll(thumb.dir)
			err := ioutil.WriteFile(thumb.path, thumb.data, 0755) // store thumbnail in cache
			if err != nil {
				logger.Errorf("Write error for thumbnail %s: %s ", thumb.path, err)
			}
		}

	}

}

//serialize file touches to avoid race conditions
func thumbnailCacheToucher() {
	for thumbPath := range touchChan {
		time := time.Now()
		os.Chtimes(*thumbPath, time, time)
	}

}

func (rs galleryRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{galleryId}", func(r chi.Router) {
		r.Use(GalleryCtx)
		r.Get("/{fileIndex}", rs.File)
	})

	return r
}

func (rs galleryRoutes) File(w http.ResponseWriter, r *http.Request) {
	gallery := r.Context().Value(galleryKey).(*models.Gallery)
	fileIndex, _ := strconv.Atoi(chi.URLParam(r, "fileIndex"))
	thumb := r.URL.Query().Get("thumb")
	w.Header().Add("Cache-Control", "max-age=604800000") // 1 Week
	if thumb == "true" {
		_, _ = w.Write(cacheThumb(gallery, fileIndex, 200))
	} else if thumb == "" {
		_, _ = w.Write(gallery.GetImage(fileIndex))
	} else {
		width, err := strconv.ParseInt(thumb, 0, 64)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		_, _ = w.Write(cacheThumb(gallery, fileIndex, int(width)))
	}
}

// get thumbnail from cache, otherwise create it and store to cache
func cacheThumb(gallery *models.Gallery, index int, width int) []byte {
	thumbPath := paths.GetThumbPath(gallery.Checksum, index, width)
	exists, _ := utils.FileExists(thumbPath)
	if exists { // if thumbnail exists in cache return that
		content, err := ioutil.ReadFile(thumbPath)
		if err == nil {
			if ThumbCacheLimit > 0 {
				touchChan <- &thumbPath // touch the file so we know which thumbs are rarely accessed
			}
			return content
		} else {
			logger.Errorf("Read Error for file %s : %s", thumbPath, err)
		}

	}
	data := gallery.GetThumbnail(index, width)
	thumbDir := paths.GetThumbDir(gallery.Checksum)
	t := newCacheThumb(thumbDir, thumbPath, data)
	writeChan <- t // write the file to cache
	return data
}

func GalleryCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		galleryID, err := strconv.Atoi(chi.URLParam(r, "galleryId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		qb := models.NewGalleryQueryBuilder()
		gallery, err := qb.Find(galleryID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), galleryKey, gallery)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
