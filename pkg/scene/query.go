package scene

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/models"
)

func BatchProcess(ctx context.Context, reader models.SceneReader, sceneFilter *models.SceneFilterType, findFilter *models.FindFilterType, fn func(scene *models.Scene) error) error {
	const batchSize = 1000

	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	page := 1
	perPage := batchSize
	findFilter.Page = &page
	findFilter.PerPage = &perPage

	more := true
	for more {
		if job.IsCancelled(ctx) {
			return nil
		}

		scenes, _, err := reader.Query(sceneFilter, findFilter)
		if err != nil {
			return fmt.Errorf("error querying for scenes: %w", err)
		}

		for _, scene := range scenes {
			if err := fn(scene); err != nil {
				return err
			}
		}

		if len(scenes) != batchSize {
			more = false
		} else {
			*findFilter.Page++
		}
	}

	return nil
}
