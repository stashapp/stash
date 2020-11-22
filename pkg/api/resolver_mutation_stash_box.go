package api

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper/stashbox"
)

func (r *mutationResolver) SubmitStashBoxFingerprints(ctx context.Context, input models.StashBoxFingerprintSubmissionInput) (bool, error) {
	boxes := config.GetStashBoxes()

	if input.StashBoxIndex < 0 || input.StashBoxIndex >= len(boxes) {
		return false, fmt.Errorf("invalid stash_box_index %d", input.StashBoxIndex)
	}

	client := stashbox.NewClient(*boxes[input.StashBoxIndex])

	return client.SubmitStashBoxFingerprints(input.SceneIds, boxes[input.StashBoxIndex].Endpoint)
}
