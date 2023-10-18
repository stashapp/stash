package heresphere

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

/*
 * This auxiliary function gathers various tags from the scene to feed the api.
 */
func (rs routes) getVideoTags(ctx context.Context, scene *models.Scene) []HeresphereVideoTag {
	processedTags := []HeresphereVideoTag{}

	if err := rs.withReadTxn(ctx, func(ctx context.Context) error {
		err := scene.LoadRelationships(ctx, rs.SceneFinder)

		processedTags = append(processedTags, rs.generateMarkerTags(ctx, scene)...)
		processedTags = append(processedTags, rs.generateTagTags(ctx, scene)...)
		processedTags = append(processedTags, rs.generatePerformerTags(ctx, scene)...)
		processedTags = append(processedTags, rs.generateGalleryTags(ctx, scene)...)
		processedTags = append(processedTags, rs.generateMovieTags(ctx, scene)...)
		processedTags = append(processedTags, rs.generateStudioTag(ctx, scene)...)
		processedTags = append(processedTags, rs.generateInteractiveTag(scene)...)
		processedTags = append(processedTags, rs.generateDirectorTag(scene)...)
		processedTags = append(processedTags, rs.generateRatingTag(scene)...)
		processedTags = append(processedTags, rs.generateWatchedTag(scene)...)
		processedTags = append(processedTags, rs.generateOrganizedTag(scene)...)
		processedTags = append(processedTags, rs.generateRatedTag(scene)...)
		processedTags = append(processedTags, rs.generateOrgasmedTag(scene)...)
		processedTags = append(processedTags, rs.generatePlayCountTag(scene)...)
		processedTags = append(processedTags, rs.generateOCounterTag(scene)...)

		return err
	}); err != nil {
		logger.Errorf("Heresphere getVideoTags generate tags error: %s\n", err.Error())
	}

	return processedTags
}
func (rs routes) generateMarkerTags(ctx context.Context, scene *models.Scene) []HeresphereVideoTag {
	// Generate marker tags
	tags := []HeresphereVideoTag{}

	markIDs, err := rs.SceneMarkerFinder.FindBySceneID(ctx, scene.ID)
	if err != nil {
		logger.Errorf("Heresphere generateMarkerTags SceneMarker.FindBySceneID error: %s\n", err.Error())
		return tags
	}

	for _, mark := range markIDs {
		tagName := mark.Title

		if ret, err := rs.TagFinder.Find(ctx, mark.PrimaryTagID); err == nil {
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
func (rs routes) generateTagTags(ctx context.Context, scene *models.Scene) []HeresphereVideoTag {
	// Generate tag tags
	tags := []HeresphereVideoTag{}

	tagIDs, err := rs.TagFinder.FindBySceneID(ctx, scene.ID)
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

func (rs routes) generatePerformerTags(ctx context.Context, scene *models.Scene) []HeresphereVideoTag {
	// Generate performer tags
	tags := []HeresphereVideoTag{}

	perfIDs, err := rs.PerformerFinder.FindBySceneID(ctx, scene.ID)
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

func (rs routes) generateGalleryTags(ctx context.Context, scene *models.Scene) []HeresphereVideoTag {
	// Generate gallery tags
	tags := []HeresphereVideoTag{}

	if scene.GalleryIDs.Loaded() {
		galleries, err := rs.GalleryFinder.FindMany(ctx, scene.GalleryIDs.List())
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

func (rs routes) generateMovieTags(ctx context.Context, scene *models.Scene) []HeresphereVideoTag {
	// Generate movie tags
	tags := []HeresphereVideoTag{}

	if scene.Movies.Loaded() {
		lst := scene.Movies.List()
		idx := make([]int, 0, len(lst))
		for _, movie := range lst {
			idx = append(idx, movie.MovieID)
		}

		movies, err := rs.MovieFinder.FindMany(ctx, idx)
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

func (rs routes) generateStudioTag(ctx context.Context, scene *models.Scene) []HeresphereVideoTag {
	// Generate studio tag
	tags := []HeresphereVideoTag{}

	if scene.StudioID != nil {
		studio, err := rs.StudioFinder.Find(ctx, *scene.StudioID)
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

func (rs routes) generateInteractiveTag(scene *models.Scene) []HeresphereVideoTag {
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

func (rs routes) generateDirectorTag(scene *models.Scene) []HeresphereVideoTag {
	// Generate director tag
	tags := []HeresphereVideoTag{}

	if len(scene.Director) > 0 {
		tags = append(tags, HeresphereVideoTag{
			Name: fmt.Sprintf("Director:%s", scene.Director),
		})
	}

	return tags
}

func (rs routes) generateRatingTag(scene *models.Scene) []HeresphereVideoTag {
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

func (rs routes) generateWatchedTag(scene *models.Scene) []HeresphereVideoTag {
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

func (rs routes) generateOrganizedTag(scene *models.Scene) []HeresphereVideoTag {
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

func (rs routes) generateRatedTag(scene *models.Scene) []HeresphereVideoTag {
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

func (rs routes) generateOrgasmedTag(scene *models.Scene) []HeresphereVideoTag {
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

func (rs routes) generatePlayCountTag(scene *models.Scene) []HeresphereVideoTag {
	return []HeresphereVideoTag{
		{
			Name: fmt.Sprintf("%s:%d", string(HeresphereCustomTagPlayCount), scene.PlayCount),
		},
	}
}

func (rs routes) generateOCounterTag(scene *models.Scene) []HeresphereVideoTag {
	return []HeresphereVideoTag{
		{
			Name: fmt.Sprintf("%s:%d", string(HeresphereCustomTagOCounter), scene.OCounter),
		},
	}
}
