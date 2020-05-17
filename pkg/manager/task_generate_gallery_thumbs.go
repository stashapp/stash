package manager

import (
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
	"sync"
)

type GenerateGthumbsTask struct {
	Gallery models.Gallery
}

func (t *GenerateGthumbsTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	generated := 0
	count := utils.ImageCount(t.Gallery.Path)
	for i := 0; i < count; i++ {
		thumbPath := paths.GetGthumbPath(t.Gallery.Checksum, i, models.DefaultGthumbWidth)
		exists, _ := utils.FileExists(thumbPath)
		if exists {
			continue
		}
		data := utils.GetThumbnail(i, models.DefaultGthumbWidth, t.Gallery.Path)
		err := utils.WriteFile(thumbPath, data)
		if err != nil {
			logger.Errorf("error writing gallery thumbnail: %s", err)
		} else {
			generated++
		}

	}
	if generated > 0 {
		logger.Infof("Generated %d thumbnails for %s", generated, t.Gallery.Path)
	}
}
