package api

import (
	"fmt"
	"strings"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/stashbox"
)

func (r *Resolver) newStashBoxClient(box models.StashBox) *stashbox.Client {
	return stashbox.NewClient(box, stashbox.ExcludeTagPatterns(manager.GetInstance().Config.GetScraperExcludeTagPatterns()))
}

func resolveStashBoxFn(indexField, endpointField string) func(index *int, endpoint *string) (*models.StashBox, error) {
	return func(index *int, endpoint *string) (*models.StashBox, error) {
		boxes := config.GetInstance().GetStashBoxes()

		// prefer endpoint over index
		if endpoint != nil {
			for _, box := range boxes {
				if strings.EqualFold(*endpoint, box.Endpoint) {
					return box, nil
				}
			}
			return nil, fmt.Errorf("stash box not found")
		}

		if index != nil {
			if *index < 0 || *index >= len(boxes) {
				return nil, fmt.Errorf("invalid %s %d", indexField, index)
			}

			return boxes[*index], nil
		}

		return nil, fmt.Errorf("%s not provided", endpointField)
	}
}

var (
	resolveStashBox              = resolveStashBoxFn("stash_box_index", "stash_box_endpoint")
	resolveStashBoxBatchTagInput = resolveStashBoxFn("endpoint", "stash_box_endpoint")
)
