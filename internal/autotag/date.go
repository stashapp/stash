package autotag

import (
	"context"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/match"
)

// extracts and sets the date from a scene's file path.
func SceneDate(ctx context.Context, s *models.Scene, rw models.SceneUpdater) error {
	// Skip if the date is already set
	if s.Date != nil {
		return nil
	}

	// Extract date from file path
	date := match.PathToDate(s.Path)
	if date == nil {
		return nil // Date not found
	}

	// Update scene object
	partial := models.NewScenePartial()
	dateModel := models.Date{Time: *date}
	partial.Date = models.NewOptionalDate(dateModel)

	// Update the database
	_, err := rw.UpdatePartial(ctx, s.ID, partial)
	return err
}
