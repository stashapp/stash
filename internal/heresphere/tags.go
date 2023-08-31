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
	"github.com/stashapp/stash/pkg/txn"
)

/*
 * This auxiliary function gathers various tags from the scene to feed the api.
 */
func getVideoTags(ctx context.Context, rs Routes, scene *models.Scene) []HeresphereVideoTag {
	processedTags := []HeresphereVideoTag{}

	// Load all relationships
	if err := txn.WithReadTxn(ctx, rs.TxnManager, func(ctx context.Context) error {
		return scene.LoadRelationships(ctx, rs.Repository.Scene)
	}); err != nil {
		logger.Errorf("Heresphere getVideoTags LoadRelationships error: %s\n", err.Error())
		return processedTags
	}

	if err := txn.WithReadTxn(ctx, rs.TxnManager, func(ctx context.Context) error {
		mark_ids, err := rs.Repository.SceneMarker.FindBySceneID(ctx, scene.ID)
		if err == nil {
			for _, mark := range mark_ids {
				tagName := mark.Title

				// Add tag name
				if ret, err := rs.Repository.Tag.Find(ctx, mark.PrimaryTagID); err == nil {
					if len(tagName) == 0 {
						tagName = ret.Name
					} else {
						tagName = fmt.Sprintf("%s - %s", tagName, ret.Name)
					}
				}

				genTag := HeresphereVideoTag{
					Name:  fmt.Sprintf("Marker:%s", tagName),
					Start: mark.Seconds * 1000,
					End:   (mark.Seconds + 60) * 1000,
				}
				processedTags = append(processedTags, genTag)
			}
		} else {
			logger.Errorf("Heresphere getVideoTags SceneMarker.FindBySceneID error: %s\n", err.Error())
		}

		tag_ids, err := rs.Repository.Tag.FindBySceneID(ctx, scene.ID)
		if err == nil {
			for _, tag := range tag_ids {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Tag:%s", tag.Name),
				}
				processedTags = append(processedTags, genTag)
			}
		} else {
			logger.Errorf("Heresphere getVideoTags Tag.FindBySceneID error: %s\n", err.Error())
		}

		perf_ids, err := rs.Repository.Performer.FindBySceneID(ctx, scene.ID)
		if err == nil {
			for _, perf := range perf_ids {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Performer:%s", perf.Name),
				}
				processedTags = append(processedTags, genTag)
			}
		} else {
			logger.Errorf("Heresphere getVideoTags Performer.FindBySceneID error: %s\n", err.Error())
		}

		if scene.GalleryIDs.Loaded() {
			galleries, err := rs.Repository.Gallery.FindMany(ctx, scene.GalleryIDs.List())
			if err == nil {
				for _, gallery := range galleries {
					genTag := HeresphereVideoTag{
						Name: fmt.Sprintf("Gallery:%s", gallery.Title),
					}
					processedTags = append(processedTags, genTag)
				}
			} else {
				logger.Errorf("Heresphere getVideoTags Gallery.FindMany error: %s\n", err.Error())
			}
		}

		if scene.Movies.Loaded() {
			lst := scene.Movies.List()
			idx := make([]int, 0, len(lst))
			for _, movie := range lst {
				idx = append(idx, movie.MovieID)
			}

			movies, err := rs.Repository.Movie.FindMany(ctx, idx)
			if err == nil {
				for _, movie := range movies {
					genTag := HeresphereVideoTag{
						Name: fmt.Sprintf("Movie:%s", movie.Name),
					}
					processedTags = append(processedTags, genTag)
				}
			} else {
				logger.Errorf("Heresphere getVideoTags Movie.FindMany error: %s\n", err.Error())
			}
		}

		if scene.StudioID != nil {
			studio, err := rs.Repository.Studio.Find(ctx, *scene.StudioID)
			if studio != nil {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Studio:%s", studio.Name),
				}
				processedTags = append(processedTags, genTag)
			}
			if err != nil {
				logger.Errorf("Heresphere getVideoTags Studio.Find error: %s\n", err.Error())
			}
		}

		primaryFile := scene.Files.Primary()
		if primaryFile != nil {
			genTag := HeresphereVideoTag{
				Name: fmt.Sprintf("%s:%s",
					string(HeresphereCustomTagInteractive),
					strconv.FormatBool(primaryFile.Interactive),
				),
			}
			processedTags = append(processedTags, genTag)

			if primaryFile.Interactive {
				funSpeed := 0
				if primaryFile.InteractiveSpeed != nil {
					funSpeed = *primaryFile.InteractiveSpeed
				}
				genTag = HeresphereVideoTag{
					Name: fmt.Sprintf("Funspeed:%d",
						funSpeed,
					),
				}
				processedTags = append(processedTags, genTag)
			}
		}

		return err
	}); err != nil {
		fmt.Printf("Heresphere getVideoTags scene reader error: %s\n", err.Error())
	}

	if len(scene.Director) > 0 {
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("Director:%s", scene.Director),
		}
		processedTags = append(processedTags, genTag)
	}

	if scene.Rating != nil {
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("Rating:%d",
				models.Rating100To5(*scene.Rating),
			),
		}
		processedTags = append(processedTags, genTag)
	}

	{
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("%s:%s",
				string(HeresphereCustomTagWatched),
				strconv.FormatBool(scene.PlayCount > 0),
			),
		}
		processedTags = append(processedTags, genTag)
	}

	{
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("%s:%s",
				string(HeresphereCustomTagOrganized),
				strconv.FormatBool(scene.Organized),
			),
		}
		processedTags = append(processedTags, genTag)
	}

	{
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("%s:%s",
				string(HeresphereCustomTagRated),
				strconv.FormatBool(scene.Rating != nil),
			),
		}
		processedTags = append(processedTags, genTag)
	}

	{
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("%s:%s",
				string(HeresphereCustomTagOrgasmed),
				strconv.FormatBool(scene.OCounter > 0),
			),
		}
		processedTags = append(processedTags, genTag)
	}

	{
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("%s:%d", string(HeresphereCustomTagPlayCount), scene.PlayCount),
		}
		processedTags = append(processedTags, genTag)
	}
	{
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("%s:%d", string(HeresphereCustomTagOCounter), scene.OCounter),
		}
		processedTags = append(processedTags, genTag)
	}

	return processedTags
}

/*
 * Processes tags and updates scene tags if applicable
 */
func handleTags(ctx context.Context, scn *models.Scene, user *HeresphereAuthReq, rs Routes, ret *scene.UpdateSet) (bool, error) {
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
		if handleAddTag(ctx, rs, tagI, &tagIDs) {
			continue
		}
		if handleAddPerformer(ctx, rs, tagI, &perfIDs) {
			continue
		}
		if handleAddMarker(ctx, rs, tagI, scn) {
			continue
		}
		if handleAddMovie(ctx, rs, tagI, scn, ret) {
			continue
		}
		if handleAddStudio(ctx, rs, tagI, scn, ret) {
			continue
		}
		if handleAddDirector(ctx, rs, tagI, scn, ret) {
			continue
		}

		// Custom
		if handleSetWatched(ctx, rs, tagI, scn, ret) {
			continue
		}
		if handleSetOrganized(ctx, rs, tagI, scn, ret) {
			continue
		}
		if handleSetRated(ctx, rs, tagI, scn, ret) {
			continue
		}
		if handleSetPlayCount(ctx, rs, tagI, scn, ret) {
			continue
		}
		if handleSetOCount(ctx, rs, tagI, scn, ret) {
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

func handleAddTag(ctx context.Context, rs Routes, tag HeresphereVideoTag, tagIDs *[]int) bool {
	if !strings.HasPrefix(tag.Name, "Tag:") {
		return false
	}

	after := strings.TrimPrefix(tag.Name, "Tag:")
	var err error
	var tagMod *models.Tag
	if err := txn.WithReadTxn(ctx, rs.TxnManager, func(ctx context.Context) error {
		// Search for tag
		tagMod, err = rs.Repository.Tag.FindByName(ctx, after, true)
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
func handleAddPerformer(ctx context.Context, rs Routes, tag HeresphereVideoTag, perfIDs *[]int) bool {
	if !strings.HasPrefix(tag.Name, "Performer:") {
		return false
	}

	after := strings.TrimPrefix(tag.Name, "Performer:")
	var err error
	var tagMod *models.Performer
	if err := txn.WithReadTxn(ctx, rs.TxnManager, func(ctx context.Context) error {
		var tagMods []*models.Performer

		// Search for performer
		if tagMods, err = rs.Repository.Performer.FindByNames(ctx, []string{after}, true); err == nil && len(tagMods) > 0 {
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
func handleAddMarker(ctx context.Context, rs Routes, tag HeresphereVideoTag, scene *models.Scene) bool {
	if !strings.HasPrefix(tag.Name, "Marker:") {
		return false
	}

	after := strings.TrimPrefix(tag.Name, "Marker:")
	var tagId *string

	if err := txn.WithReadTxn(ctx, rs.TxnManager, func(ctx context.Context) error {
		var err error
		var markerResult []*models.MarkerStringsResultType
		searchType := "count"

		// Search for marker
		if markerResult, err = rs.Repository.SceneMarker.GetMarkerStrings(ctx, &after, &searchType); len(markerResult) > 0 {
			tagId = &markerResult[0].ID

			// Search for tag
			if markers, err := rs.Repository.SceneMarker.FindBySceneID(ctx, scene.ID); err == nil {
				i, e := strconv.Atoi(*tagId)
				if e == nil {
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
	}); tagId != nil {
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

			if rs.Repository.SceneMarker.Create(ctx, &newMarker) != nil {
				logger.Errorf("Heresphere handleTags SceneMarker.Create error: %s\n", err.Error())
			}
		}
	}

	return true
}
func handleAddMovie(ctx context.Context, rs Routes, tag HeresphereVideoTag, scene *models.Scene, ret *scene.UpdateSet) bool {
	if !strings.HasPrefix(tag.Name, "Movie:") {
		return false
	}

	after := strings.TrimPrefix(tag.Name, "Movie:")

	var err error
	var tagMod *models.Movie
	if err := txn.WithReadTxn(ctx, rs.TxnManager, func(ctx context.Context) error {
		// Search for performer
		tagMod, err = rs.Repository.Movie.FindByName(ctx, after, true)
		return err
	}); err == nil {
		ret.Partial.MovieIDs.Mode = models.RelationshipUpdateModeSet
		ret.Partial.MovieIDs.AddUnique(models.MoviesScenes{
			MovieID:    tagMod.ID,
			SceneIndex: &scene.ID,
		})
	}

	return true
}
func handleAddStudio(ctx context.Context, rs Routes, tag HeresphereVideoTag, scene *models.Scene, ret *scene.UpdateSet) bool {
	if !strings.HasPrefix(tag.Name, "Studio:") {
		return false
	}

	after := strings.TrimPrefix(tag.Name, "Studio:")

	var err error
	var tagMod *models.Studio
	if err := txn.WithReadTxn(ctx, rs.TxnManager, func(ctx context.Context) error {
		// Search for performer
		tagMod, err = rs.Repository.Studio.FindByName(ctx, after, true)
		return err
	}); err == nil {
		ret.Partial.StudioID.Set = true
		ret.Partial.StudioID.Value = tagMod.ID
	}

	return true
}
func handleAddDirector(ctx context.Context, rs Routes, tag HeresphereVideoTag, scene *models.Scene, ret *scene.UpdateSet) bool {
	if !strings.HasPrefix(tag.Name, "Director:") {
		return false
	}

	after := strings.TrimPrefix(tag.Name, "Director:")
	ret.Partial.Director.Set = true
	ret.Partial.Director.Value = after

	return true
}

// Will be overwritten if PlayCount tag is updated
func handleSetWatched(ctx context.Context, rs Routes, tag HeresphereVideoTag, scene *models.Scene, ret *scene.UpdateSet) bool {
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
func handleSetOrganized(ctx context.Context, rs Routes, tag HeresphereVideoTag, scene *models.Scene, ret *scene.UpdateSet) bool {
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
func handleSetRated(ctx context.Context, rs Routes, tag HeresphereVideoTag, scene *models.Scene, ret *scene.UpdateSet) bool {
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
func handleSetPlayCount(ctx context.Context, rs Routes, tag HeresphereVideoTag, scene *models.Scene, ret *scene.UpdateSet) bool {
	prefix := string(HeresphereCustomTagPlayCount) + ":"
	if !strings.HasPrefix(tag.Name, prefix) {
		return false
	}

	after := strings.TrimPrefix(tag.Name, prefix)
	if numRes, err := strconv.Atoi(after); err != nil {
		ret.Partial.PlayCount.Set = true
		ret.Partial.PlayCount.Value = numRes
	}

	return true
}
func handleSetOCount(ctx context.Context, rs Routes, tag HeresphereVideoTag, scene *models.Scene, ret *scene.UpdateSet) bool {
	prefix := string(HeresphereCustomTagOCounter) + ":"
	if !strings.HasPrefix(tag.Name, prefix) {
		return false
	}

	after := strings.TrimPrefix(tag.Name, prefix)
	if numRes, err := strconv.Atoi(after); err != nil {
		ret.Partial.OCounter.Set = true
		ret.Partial.OCounter.Value = numRes
	}

	return true
}
