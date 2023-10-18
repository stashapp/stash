package heresphere

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

/*
 * Processes tags and updates scene tags if applicable
 */
func (rs routes) handleTags(ctx context.Context, scn *models.Scene, user *HeresphereAuthReq, ret *scene.UpdateSet) (bool, error) {
	// Search input tags and add/create any new ones
	var tagIDs []int
	var perfIDs []int

	for _, tagI := range *user.Tags {
		// If missing
		if len(tagI.Name) == 0 {
			continue
		}

		// FUTURE IMPROVEMENT: Switch to CutPrefix as it's nicer (1.20+)
		// FUTURE IMPROVEMENT: Consider batching searches
		if rs.handleAddTag(ctx, tagI, &tagIDs) {
			continue
		}
		if rs.handleAddPerformer(ctx, tagI, &perfIDs) {
			continue
		}
		if rs.handleAddMarker(ctx, tagI, scn) {
			continue
		}
		if rs.handleAddStudio(ctx, tagI, scn, ret) {
			continue
		}
		if rs.handleAddDirector(ctx, tagI, scn, ret) {
			continue
		}

		// Custom
		if rs.handleSetWatched(ctx, tagI, scn, ret) {
			continue
		}
		if rs.handleSetOrganized(ctx, tagI, scn, ret) {
			continue
		}
		if rs.handleSetRated(ctx, tagI, scn, ret) {
			continue
		}
		if rs.handleSetPlayCount(ctx, tagI, scn, ret) {
			continue
		}
		if rs.handleSetOCount(ctx, tagI, scn, ret) {
			continue
		}
	}

	// Update tags
	ret.Partial.TagIDs = &models.UpdateIDs{
		IDs:  tagIDs,
		Mode: models.RelationshipUpdateModeSet,
	}
	// Update performers
	ret.Partial.PerformerIDs = &models.UpdateIDs{
		IDs:  perfIDs,
		Mode: models.RelationshipUpdateModeSet,
	}

	return true, nil
}

func (rs routes) handleAddTag(ctx context.Context, tag HeresphereVideoTag, tagIDs *[]int) bool {
	if !strings.HasPrefix(tag.Name, "Tag:") {
		return false
	}

	after := strings.TrimPrefix(tag.Name, "Tag:")
	var err error
	var tagMod *models.Tag
	if err := rs.withReadTxn(ctx, func(ctx context.Context) error {
		// Search for tag
		tagMod, err = rs.TagFinder.FindByName(ctx, after, true)
		return err
	}); err != nil {
		fmt.Printf("Heresphere handleTags Tag.FindByName error: %s\n", err.Error())
		tagMod = nil
	}

	if tagMod != nil {
		*tagIDs = append(*tagIDs, tagMod.ID)
	}

	return true
}
func (rs routes) handleAddPerformer(ctx context.Context, tag HeresphereVideoTag, perfIDs *[]int) bool {
	if !strings.HasPrefix(tag.Name, "Performer:") {
		return false
	}

	after := strings.TrimPrefix(tag.Name, "Performer:")
	var err error
	var tagMod *models.Performer
	if err := rs.withReadTxn(ctx, func(ctx context.Context) error {
		var tagMods []*models.Performer

		// Search for performer
		if tagMods, err = rs.PerformerFinder.FindByNames(ctx, []string{after}, true); err == nil && len(tagMods) > 0 {
			tagMod = tagMods[0]
		}

		return err
	}); err != nil {
		fmt.Printf("Heresphere handleTags Performer.FindByNames error: %s\n", err.Error())
		tagMod = nil
	}

	if tagMod != nil {
		*perfIDs = append(*perfIDs, tagMod.ID)
	}

	return true
}
func (rs routes) handleAddMarker(ctx context.Context, tag HeresphereVideoTag, scene *models.Scene) bool {
	if !strings.HasPrefix(tag.Name, "Marker:") {
		return false
	}

	after := strings.TrimPrefix(tag.Name, "Marker:")
	var tagId *string

	if err := rs.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		var markerResult []*models.MarkerStringsResultType
		searchType := "count"

		// Search for marker
		if markerResult, err = rs.SceneMarkerFinder.GetMarkerStrings(ctx, &after, &searchType); len(markerResult) > 0 {
			tagId = &markerResult[0].ID

			// Search for tag
			if markers, err := rs.SceneMarkerFinder.FindBySceneID(ctx, scene.ID); err == nil {
				i, err := strconv.Atoi(*tagId)
				if err == nil {
					// Note: Currently we search if a marker exists.
					// If it doesn't, create it.
					// This also means that markers CANNOT be deleted using the api.
					for _, marker := range markers {
						if marker.Seconds == tag.Start &&
							marker.SceneID == scene.ID &&
							marker.PrimaryTagID == i {
							tagId = nil
						}
					}
				}
			}
		}

		return err
	}); err != nil || tagId != nil {
		// Create marker
		i, e := strconv.Atoi(*tagId)
		if e == nil {
			currentTime := time.Now()
			newMarker := models.SceneMarker{
				Title:        "",
				Seconds:      tag.Start,
				PrimaryTagID: i,
				SceneID:      scene.ID,
				CreatedAt:    currentTime,
				UpdatedAt:    currentTime,
			}

			if err := rs.withTxn(ctx, func(ctx context.Context) error {
				return rs.SceneMarkerFinder.Create(ctx, &newMarker)
			}); err != nil {
				logger.Errorf("Heresphere handleTags SceneMarker.Create error: %s\n", err.Error())
			}
		}
	}

	return true
}
func (rs routes) handleAddStudio(ctx context.Context, tag HeresphereVideoTag, scene *models.Scene, ret *scene.UpdateSet) bool {
	if !strings.HasPrefix(tag.Name, "Studio:") {
		return false
	}

	after := strings.TrimPrefix(tag.Name, "Studio:")

	var err error
	var tagMod *models.Studio
	if err := rs.withReadTxn(ctx, func(ctx context.Context) error {
		// Search for performer
		tagMod, err = rs.StudioFinder.FindByName(ctx, after, true)
		return err
	}); err == nil {
		ret.Partial.StudioID.Set = true
		ret.Partial.StudioID.Value = tagMod.ID
	}

	return true
}
func (rs routes) handleAddDirector(ctx context.Context, tag HeresphereVideoTag, scene *models.Scene, ret *scene.UpdateSet) bool {
	if !strings.HasPrefix(tag.Name, "Director:") {
		return false
	}

	after := strings.TrimPrefix(tag.Name, "Director:")
	ret.Partial.Director.Set = true
	ret.Partial.Director.Value = after

	return true
}

// Will be overwritten if PlayCount tag is updated
func (rs routes) handleSetWatched(ctx context.Context, tag HeresphereVideoTag, scene *models.Scene, ret *scene.UpdateSet) bool {
	prefix := string(HeresphereCustomTagWatched) + ":"
	if !strings.HasPrefix(tag.Name, prefix) {
		return false
	}

	after := strings.TrimPrefix(tag.Name, prefix)
	if b, err := strconv.ParseBool(after); err == nil {
		// Plays chicken
		if b && scene.PlayCount == 0 {
			ret.Partial.PlayCount.Set = true
			ret.Partial.PlayCount.Value = 1
		} else if !b {
			ret.Partial.PlayCount.Set = true
			ret.Partial.PlayCount.Value = 0
		}
	}

	return true
}
func (rs routes) handleSetOrganized(ctx context.Context, tag HeresphereVideoTag, scene *models.Scene, ret *scene.UpdateSet) bool {
	prefix := string(HeresphereCustomTagOrganized) + ":"
	if !strings.HasPrefix(tag.Name, prefix) {
		return false
	}

	after := strings.TrimPrefix(tag.Name, prefix)
	if b, err := strconv.ParseBool(after); err == nil {
		ret.Partial.Organized.Set = true
		ret.Partial.Organized.Value = b
	}

	return true
}
func (rs routes) handleSetRated(ctx context.Context, tag HeresphereVideoTag, scene *models.Scene, ret *scene.UpdateSet) bool {
	prefix := string(HeresphereCustomTagRated) + ":"
	if !strings.HasPrefix(tag.Name, prefix) {
		return false
	}

	after := strings.TrimPrefix(tag.Name, prefix)
	if b, err := strconv.ParseBool(after); err == nil && !b {
		ret.Partial.Rating.Set = true
		ret.Partial.Rating.Null = true
	}

	return true
}
func (rs routes) handleSetPlayCount(ctx context.Context, tag HeresphereVideoTag, scene *models.Scene, ret *scene.UpdateSet) bool {
	prefix := string(HeresphereCustomTagPlayCount) + ":"
	if !strings.HasPrefix(tag.Name, prefix) {
		return false
	}

	after := strings.TrimPrefix(tag.Name, prefix)
	if numRes, err := strconv.Atoi(after); err == nil {
		ret.Partial.PlayCount.Set = true
		ret.Partial.PlayCount.Value = numRes
	}

	return true
}
func (rs routes) handleSetOCount(ctx context.Context, tag HeresphereVideoTag, scene *models.Scene, ret *scene.UpdateSet) bool {
	prefix := string(HeresphereCustomTagOCounter) + ":"
	if !strings.HasPrefix(tag.Name, prefix) {
		return false
	}

	after := strings.TrimPrefix(tag.Name, prefix)
	if numRes, err := strconv.Atoi(after); err == nil {
		ret.Partial.OCounter.Set = true
		ret.Partial.OCounter.Value = numRes
	}

	return true
}
