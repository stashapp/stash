package heresphere

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/txn"
)

/*
 * This auxiliary function gathers various tags from the scene to feed the api.
 */
func getVideoTags(rs Routes, ctx context.Context, scene *models.Scene) []HeresphereVideoTag {
	processedTags := []HeresphereVideoTag{}

	// Load all relationships
	if err := txn.WithReadTxn(ctx, rs.TxnManager, func(ctx context.Context) error {
		return scene.LoadRelationships(ctx, rs.Repository.Scene)
	}); err != nil {
		return processedTags
	}

	txn.WithReadTxn(ctx, rs.TxnManager, func(ctx context.Context) error {
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

		}

		tag_ids, err := rs.Repository.Tag.FindBySceneID(ctx, scene.ID)
		if err == nil {
			for _, tag := range tag_ids {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Tag:%s", tag.Name),
				}
				processedTags = append(processedTags, genTag)
			}
		}

		perf_ids, err := rs.Repository.Performer.FindBySceneID(ctx, scene.ID)
		if err == nil {
			for _, perf := range perf_ids {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Performer:%s", perf.Name),
				}
				processedTags = append(processedTags, genTag)
			}
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
			}
		}

		if scene.StudioID != nil {
			studio, _ := rs.Repository.Studio.Find(ctx, *scene.StudioID)
			if studio != nil {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Studio:%s", studio.Name),
				}
				processedTags = append(processedTags, genTag)
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

		return err
	})

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
func handleTags(tags *[]HeresphereVideoTag, ctx context.Context, scn *models.Scene, user HeresphereAuthReq, rs Routes, ret *scene.UpdateSet) (bool, error) {
	// Search input tags and add/create any new ones
	var tagIDs []int
	var perfIDs []int

	for _, tagI := range *user.Tags {
		// If missing
		if len(tagI.Name) == 0 {
			continue
		}

		// If add tag
		// FUTURE IMPROVEMENT: Switch to CutPrefix as it's nicer (1.20+)
		// FUTURE IMPROVEMENT: Consider batching searches
		if strings.HasPrefix(tagI.Name, "Tag:") {
			after := strings.TrimPrefix(tagI.Name, "Tag:")
			var err error
			var tagMod *models.Tag
			if err := txn.WithReadTxn(ctx, rs.TxnManager, func(ctx context.Context) error {
				// Search for tag
				tagMod, err = rs.Repository.Tag.FindByName(ctx, after, true)
				return err
			}); err != nil {
				tagMod = nil
			}

			if tagMod != nil {
				tagIDs = append(tagIDs, tagMod.ID)
			}
			continue
		}

		// If add performer
		if strings.HasPrefix(tagI.Name, "Performer:") {
			after := strings.TrimPrefix(tagI.Name, "Performer:")
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
				tagMod = nil
			}

			if tagMod != nil {
				perfIDs = append(perfIDs, tagMod.ID)
			}
			continue
		}

		// If add marker
		if strings.HasPrefix(tagI.Name, "Marker:") {
			after := strings.TrimPrefix(tagI.Name, "Marker:")
			var tagId *string
			if err := txn.WithReadTxn(ctx, rs.TxnManager, func(ctx context.Context) error {
				var err error
				var markerResult []*models.MarkerStringsResultType
				searchType := "count"

				// Search for marker
				if markerResult, err = rs.Repository.SceneMarker.GetMarkerStrings(ctx, &after, &searchType); err == nil && len(markerResult) > 0 {
					tagId = &markerResult[0].ID

					// Search for tag
					if markers, err := rs.Repository.SceneMarker.FindBySceneID(ctx, scn.ID); err == nil {
						i, e := strconv.Atoi(*tagId)
						if e == nil {
							// Note: Currently we search if a marker exists.
							// If it doesn't, create it.
							// This also means that markers CANNOT be deleted using the api.
							for _, marker := range markers {
								if marker.Seconds == tagI.Start &&
									marker.SceneID == scn.ID &&
									marker.PrimaryTagID == i {
									tagId = nil
								}
							}
						}
					}
				}

				return err
			}); err != nil {
				// Create marker
				/*if tagId != nil {
					newTag := SceneMarkerCreateInput{
						Seconds:      tagI.Start,
						SceneID:      strconv.Itoa(scn.ID),
						PrimaryTagID: *tagId,
					}
					if _, err := rs.resolver.Mutation().SceneMarkerCreate(context.Background(), newTag); err != nil {
						return false, err
					}
				}*/
			}
			continue
		}

		if strings.HasPrefix(tagI.Name, "Movie:") {
			after := strings.TrimPrefix(tagI.Name, "Movie:")

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
					SceneIndex: &scn.ID,
				})
			}
			continue
		}
		if strings.HasPrefix(tagI.Name, "Studio:") {
			after := strings.TrimPrefix(tagI.Name, "Studio:")

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
			continue
		}
		if strings.HasPrefix(tagI.Name, "Director:") {
			after := strings.TrimPrefix(tagI.Name, "Director:")
			ret.Partial.Director.Set = true
			ret.Partial.Director.Value = after
			continue
		}

		// Custom
		{
			tagName := tagI.Name

			// Will be overwritten if PlayCount tag is updated
			prefix := string(HeresphereCustomTagWatched) + ":"
			if strings.HasPrefix(tagName, prefix) {
				after := strings.TrimPrefix(tagName, prefix)
				if b, err := strconv.ParseBool(after); err == nil {
					// Plays chicken
					if b && scn.PlayCount == 0 {
						ret.Partial.PlayCount.Set = true
						ret.Partial.PlayCount.Value = 1
					} else if !b {
						ret.Partial.PlayCount.Set = true
						ret.Partial.PlayCount.Value = 0
					}
				}
				continue
			}
			prefix = string(HeresphereCustomTagOrganized) + ":"
			if strings.HasPrefix(tagName, prefix) {
				after := strings.TrimPrefix(tagName, prefix)
				if b, err := strconv.ParseBool(after); err == nil {
					ret.Partial.Organized.Set = true
					ret.Partial.Organized.Value = b
				}
				continue
			}
			prefix = string(HeresphereCustomTagRated) + ":"
			if strings.HasPrefix(tagName, prefix) {
				after := strings.TrimPrefix(tagName, prefix)
				if b, err := strconv.ParseBool(after); err == nil && !b {
					ret.Partial.Rating.Set = true
					ret.Partial.Rating.Null = true
				}
				continue
			}

			// Set numbers
			prefix = string(HeresphereCustomTagPlayCount) + ":"
			if strings.HasPrefix(tagName, prefix) {
				after := strings.TrimPrefix(tagName, prefix)
				if numRes, err := strconv.Atoi(after); err != nil {
					ret.Partial.PlayCount.Set = true
					ret.Partial.PlayCount.Value = numRes
				}
				continue
			}
			prefix = string(HeresphereCustomTagOCounter) + ":"
			if strings.HasPrefix(tagName, prefix) {
				after := strings.TrimPrefix(tagName, prefix)
				if numRes, err := strconv.Atoi(after); err != nil {
					ret.Partial.OCounter.Set = true
					ret.Partial.OCounter.Value = numRes
				}
				continue
			}
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
