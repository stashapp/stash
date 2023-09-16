package heresphere

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

/*
 * This auxiliary function gathers various tags from the scene to feed the api.
 */
func getVideoTags(ctx context.Context, rs Routes, scene *models.Scene) []HeresphereVideoTag {
	processedTags := []HeresphereVideoTag{}

	if err := txn.WithReadTxn(ctx, rs.TxnManager, func(ctx context.Context) error {
		err := scene.LoadRelationships(ctx, rs.Repository.Scene)

		processedTags = append(processedTags, generateMarkerTags(ctx, rs, scene)...)
		processedTags = append(processedTags, generateTagTags(ctx, rs, scene)...)
		processedTags = append(processedTags, generatePerformerTags(ctx, rs, scene)...)
		processedTags = append(processedTags, generateGalleryTags(ctx, rs, scene)...)
		processedTags = append(processedTags, generateMovieTags(ctx, rs, scene)...)
		processedTags = append(processedTags, generateStudioTag(ctx, rs, scene)...)
		processedTags = append(processedTags, generateInteractiveTag(scene)...)
		processedTags = append(processedTags, generateDirectorTag(scene)...)
		processedTags = append(processedTags, generateRatingTag(scene)...)
		processedTags = append(processedTags, generateWatchedTag(scene)...)
		processedTags = append(processedTags, generateOrganizedTag(scene)...)
		processedTags = append(processedTags, generateRatedTag(scene)...)
		processedTags = append(processedTags, generateOrgasmedTag(scene)...)
		processedTags = append(processedTags, generatePlayCountTag(scene)...)
		processedTags = append(processedTags, generateOCounterTag(scene)...)

		return err
	}); err != nil {
		logger.Errorf("Heresphere getVideoTags generate tags error: %s\n", err.Error())
	}

	return processedTags
}
func generateMarkerTags(ctx context.Context, rs Routes, scene *models.Scene) []HeresphereVideoTag {
	// Generate marker tags
	tags := []HeresphereVideoTag{}

	markIDs, err := rs.Repository.SceneMarker.FindBySceneID(ctx, scene.ID)
	if err != nil {
		logger.Errorf("Heresphere generateMarkerTags SceneMarker.FindBySceneID error: %s\n", err.Error())
		return tags
	}

	for _, mark := range markIDs {
		tagName := mark.Title

		if ret, err := rs.Repository.Tag.Find(ctx, mark.PrimaryTagID); err == nil {
			if len(tagName) == 0 {
				tagName = ret.Name
			} else {
				tagName = fmt.Sprintf("%s - %s", tagName, ret.Name)
			}
		}

		tags = append(tags, HeresphereVideoTag{
			Name:  fmt.Sprintf("Marker:%s", tagName),
			Start: mark.Seconds * 1000,
			End:   (mark.Seconds + 60) * 1000,
		})
	}

	return tags
}
func generateTagTags(ctx context.Context, rs Routes, scene *models.Scene) []HeresphereVideoTag {
	// Generate tag tags
	tags := []HeresphereVideoTag{}

	tagIDs, err := rs.Repository.Tag.FindBySceneID(ctx, scene.ID)
	if err != nil {
		logger.Errorf("Heresphere generateTagTags Tag.FindBySceneID error: %s\n", err.Error())
		return tags
	}

	for _, tag := range tagIDs {
		tags = append(tags, HeresphereVideoTag{
			Name: fmt.Sprintf("Tag:%s", tag.Name),
		})
	}

	return tags
}

func generatePerformerTags(ctx context.Context, rs Routes, scene *models.Scene) []HeresphereVideoTag {
	// Generate performer tags
	tags := []HeresphereVideoTag{}

	perfIDs, err := rs.Repository.Performer.FindBySceneID(ctx, scene.ID)
	if err != nil {
		logger.Errorf("Heresphere generatePerformerTags Performer.FindBySceneID error: %s\n", err.Error())
		return tags
	}

	hasFavPerformer := false
	for _, perf := range perfIDs {
		tags = append(tags, HeresphereVideoTag{
			Name: fmt.Sprintf("Performer:%s", perf.Name),
		})
		hasFavPerformer = hasFavPerformer || perf.Favorite
	}

	tags = append(tags, HeresphereVideoTag{
		Name: fmt.Sprintf("HasFavoritedPerformer:%s", strconv.FormatBool(hasFavPerformer)),
	})

	return tags
}

func generateGalleryTags(ctx context.Context, rs Routes, scene *models.Scene) []HeresphereVideoTag {
	// Generate gallery tags
	tags := []HeresphereVideoTag{}

	if scene.GalleryIDs.Loaded() {
		galleries, err := rs.Repository.Gallery.FindMany(ctx, scene.GalleryIDs.List())
		if err != nil {
			logger.Errorf("Heresphere generateGalleryTags Gallery.FindMany error: %s\n", err.Error())
			return tags
		}

		for _, gallery := range galleries {
			tags = append(tags, HeresphereVideoTag{
				Name: fmt.Sprintf("Gallery:%s", gallery.Title),
			})
		}
	}

	return tags
}

func generateMovieTags(ctx context.Context, rs Routes, scene *models.Scene) []HeresphereVideoTag {
	// Generate movie tags
	tags := []HeresphereVideoTag{}

	if scene.Movies.Loaded() {
		lst := scene.Movies.List()
		idx := make([]int, 0, len(lst))
		for _, movie := range lst {
			idx = append(idx, movie.MovieID)
		}

		movies, err := rs.Repository.Movie.FindMany(ctx, idx)
		if err != nil {
			logger.Errorf("Heresphere generateMovieTags Movie.FindMany error: %s\n", err.Error())
			return tags
		}

		for _, movie := range movies {
			tags = append(tags, HeresphereVideoTag{
				Name: fmt.Sprintf("Movie:%s", movie.Name),
			})
		}
	}

	return tags
}

func generateStudioTag(ctx context.Context, rs Routes, scene *models.Scene) []HeresphereVideoTag {
	// Generate studio tag
	tags := []HeresphereVideoTag{}

	if scene.StudioID != nil {
		studio, err := rs.Repository.Studio.Find(ctx, *scene.StudioID)
		if err != nil {
			logger.Errorf("Heresphere generateStudioTag Studio.Find error: %s\n", err.Error())
			return tags
		}

		tags = append(tags, HeresphereVideoTag{
			Name: fmt.Sprintf("Studio:%s", studio.Name),
		})
	}

	return tags
}

func generateInteractiveTag(scene *models.Scene) []HeresphereVideoTag {
	// Generate interactive tag
	tags := []HeresphereVideoTag{}

	primaryFile := scene.Files.Primary()
	if primaryFile != nil {
		tags = append(tags, HeresphereVideoTag{
			Name: fmt.Sprintf("%s:%s",
				string(HeresphereCustomTagInteractive),
				strconv.FormatBool(primaryFile.Interactive),
			),
		})

		if primaryFile.Interactive {
			funSpeed := 0
			if primaryFile.InteractiveSpeed != nil {
				funSpeed = *primaryFile.InteractiveSpeed
			}
			tags = append(tags, HeresphereVideoTag{
				Name: fmt.Sprintf("Funspeed:%d",
					funSpeed,
				),
			})
		}
	}

	return tags
}

func generateDirectorTag(scene *models.Scene) []HeresphereVideoTag {
	// Generate director tag
	tags := []HeresphereVideoTag{}

	if len(scene.Director) > 0 {
		tags = append(tags, HeresphereVideoTag{
			Name: fmt.Sprintf("Director:%s", scene.Director),
		})
	}

	return tags
}

func generateRatingTag(scene *models.Scene) []HeresphereVideoTag {
	// Generate rating tag
	tags := []HeresphereVideoTag{}

	if scene.Rating != nil {
		tags = append(tags, HeresphereVideoTag{
			Name: fmt.Sprintf("Rating:%d",
				models.Rating100To5(*scene.Rating),
			),
		})
	}

	return tags
}

func generateWatchedTag(scene *models.Scene) []HeresphereVideoTag {
	// Generate watched tag
	return []HeresphereVideoTag{
		{
			Name: fmt.Sprintf("%s:%s",
				string(HeresphereCustomTagWatched),
				strconv.FormatBool(scene.PlayCount > 0),
			),
		},
	}
}

func generateOrganizedTag(scene *models.Scene) []HeresphereVideoTag {
	// Generate organized tag
	return []HeresphereVideoTag{
		{
			Name: fmt.Sprintf("%s:%s",
				string(HeresphereCustomTagOrganized),
				strconv.FormatBool(scene.Organized),
			),
		},
	}
}

func generateRatedTag(scene *models.Scene) []HeresphereVideoTag {
	// Generate rated tag
	return []HeresphereVideoTag{
		{
			Name: fmt.Sprintf("%s:%s",
				string(HeresphereCustomTagRated),
				strconv.FormatBool(scene.Rating != nil),
			),
		},
	}
}

func generateOrgasmedTag(scene *models.Scene) []HeresphereVideoTag {
	// Generate orgasmed tag
	return []HeresphereVideoTag{
		{
			Name: fmt.Sprintf("%s:%s",
				string(HeresphereCustomTagOrgasmed),
				strconv.FormatBool(scene.OCounter > 0),
			),
		},
	}
}

func generatePlayCountTag(scene *models.Scene) []HeresphereVideoTag {
	return []HeresphereVideoTag{
		{
			Name: fmt.Sprintf("%s:%d", string(HeresphereCustomTagPlayCount), scene.PlayCount),
		},
	}
}

func generateOCounterTag(scene *models.Scene) []HeresphereVideoTag {
	return []HeresphereVideoTag{
		{
			Name: fmt.Sprintf("%s:%d", string(HeresphereCustomTagOCounter), scene.OCounter),
		},
	}
}
