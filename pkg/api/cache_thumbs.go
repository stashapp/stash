package api

import (
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
	"io/ioutil"
)

type thumbBuffer struct {
	path string
	dir  string
	data []byte
}

func newCacheThumb(dir string, path string, data []byte) *thumbBuffer {
	t := thumbBuffer{dir: dir, path: path, data: data}
	return &t
}

var writeChan chan *thumbBuffer
var touchChan chan *string

func startThumbCache() { // TODO add extra wait, close chan code if/when stash gets a stop mode

	writeChan = make(chan *thumbBuffer, 20)
	go thumbnailCacheWriter()
}

//serialize file writes to avoid race conditions
func thumbnailCacheWriter() {

	for thumb := range writeChan {
		exists, _ := utils.FileExists(thumb.path)
		if !exists {
			err := utils.WriteFile(thumb.path, thumb.data)
			if err != nil {
				logger.Errorf("Write error for thumbnail %s: %s ", thumb.path, err)
			}
		}
	}

}

// get thumbnail from cache, otherwise create it and store to cache
func cacheGthumb(gallery *models.Gallery, index int, width int) []byte {
	thumbPath := paths.GetGthumbPath(gallery.Checksum, index, width)
	exists, _ := utils.FileExists(thumbPath)
	if exists { // if thumbnail exists in cache return that
		content, err := ioutil.ReadFile(thumbPath)
		if err == nil {
			return content
		} else {
			logger.Errorf("Read Error for file %s : %s", thumbPath, err)
		}

	}
	data := gallery.GetThumbnail(index, width)
	thumbDir := paths.GetGthumbDir(gallery.Checksum)
	t := newCacheThumb(thumbDir, thumbPath, data)
	writeChan <- t // write the file to cache
	return data
}

// create all thumbs for a given gallery
func CreateGthumbs(gallery *models.Gallery) {
	count := gallery.ImageCount()
	for i := 0; i < count; i++ {
		cacheGthumb(gallery, i, models.DefaultGthumbWidth)
	}
}
