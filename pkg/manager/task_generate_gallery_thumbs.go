package manager

import (
	"github.com/remeh/sizedwaitgroup"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type GenerateGthumbsTask struct {
	Gallery   models.Gallery
	Overwrite bool
}

func (t *GenerateGthumbsTask) Start(wg *sizedwaitgroup.SizedWaitGroup) {
	defer wg.Done()
	generated := 0
	count := t.Gallery.ImageCount()
	for i := 0; i < count; i++ {
		thumbPath := paths.GetGthumbPath(t.Gallery.Checksum, i, models.DefaultGthumbWidth)
		exists, _ := utils.FileExists(thumbPath)
		if !t.Overwrite && exists {
			continue
		}
		data := t.Gallery.GetThumbnail(i, models.DefaultGthumbWidth)
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
