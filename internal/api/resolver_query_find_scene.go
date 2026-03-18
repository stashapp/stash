package api

import (
	"fmt"

	"errors"

	"context"
	"slices"
	"strconv"

	"github.com/99designs/gqlgen/graphql"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

func (r *queryResolver) FindScene(ctx context.Context, id *string, checksum *string) (*models.Scene, error) {
	var scene *models.Scene
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Scene
		var err error
		if id != nil {
			idInt, err := strconv.Atoi(*id)
			if err != nil {
				return err
			}
			scene, err = qb.Find(ctx, idInt)
			if err != nil {
				return err
			}
		} else if checksum != nil {
			var scenes []*models.Scene
			scenes, err = qb.FindByChecksum(ctx, *checksum)
			if len(scenes) > 0 {
				scene = scenes[0]
			}
		}

		return err
	}); err != nil {
		return nil, err
	}

	return scene, nil
}

func (r *queryResolver) FindSceneByHash(ctx context.Context, input SceneHashInput) (*models.Scene, error) {
	var scene *models.Scene

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Scene
		if input.Checksum != nil {
			scenes, err := qb.FindByChecksum(ctx, *input.Checksum)
			if err != nil {
				return err
			}
			if len(scenes) > 0 {
				scene = scenes[0]
			}
		}

		if scene == nil && input.Oshash != nil {
			scenes, err := qb.FindByOSHash(ctx, *input.Oshash)
			if err != nil {
				return err
			}
			if len(scenes) > 0 {
				scene = scenes[0]
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return scene, nil
}

func (r *queryResolver) FindScenes(
	ctx context.Context,
	sceneFilter *models.SceneFilterType,
	savedFilterID *string,
	sceneIDs []int,
	ids []string,
	filter *models.FindFilterType,
) (ret *FindScenesResultType, err error) {
	if sceneFilter != nil && savedFilterID != nil {
		return nil, errors.New("cannot provide both sceneFilter and saved_filter_id")
	}

	var finalFilter *models.SceneFilterType
	if savedFilterID != nil {
		finalFilter = &models.SceneFilterType{}
		var mode models.FilterMode
		switch "sceneFilter" {
		case "sceneFilter":
			mode = models.FilterModeScenes
		case "performerFilter":
			mode = models.FilterModePerformers
		case "studioFilter":
			mode = models.FilterModeStudios
		case "galleryFilter":
			mode = models.FilterModeGalleries
		case "sceneMarkerFilter":
			mode = models.FilterModeSceneMarkers
		case "movieFilter":
			mode = models.FilterModeMovies
		case "groupFilter":
			mode = models.FilterModeGroups
		case "tagFilter":
			mode = models.FilterModeTags
		case "imageFilter":
			mode = models.FilterModeImages
		default:
			return nil, fmt.Errorf("saved filters are not supported for %s", "sceneFilter")
		}

		mergedFindFilter, err := r.resolveSavedFilter(ctx, *savedFilterID, mode, finalFilter, filter)
		if err != nil {
			return nil, err
		}
		filter = mergedFindFilter
	} else {
		finalFilter = sceneFilter
	}
	if len(ids) > 0 {
		sceneIDs, err = handleIDList(ids, "ids")
		if err != nil {
			return nil, err
		}
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var scenes []*models.Scene
		var err error

		fields := graphql.CollectAllFields(ctx)
		result := &models.SceneQueryResult{}

		if len(sceneIDs) > 0 {
			scenes, err = r.repository.Scene.FindMany(ctx, sceneIDs)
			if err == nil {
				result.Count = len(scenes)
				for _, s := range scenes {
					if err = s.LoadPrimaryFile(ctx, r.repository.File); err != nil {
						break
					}

					f := s.Files.Primary()
					if f == nil {
						continue
					}

					result.TotalDuration += f.Duration

					result.TotalSize += float64(f.Size)
				}
			}
		} else {
			result, err = r.repository.Scene.Query(ctx, models.SceneQueryOptions{
				QueryOptions: models.QueryOptions{
					FindFilter: filter,
					Count:      slices.Contains(fields, "count"),
				},
				SceneFilter:   finalFilter,
				TotalDuration: slices.Contains(fields, "duration"),
				TotalSize:     slices.Contains(fields, "filesize"),
			})
			if err == nil {
				scenes, err = result.Resolve(ctx)
			}
		}

		if err != nil {
			return err
		}

		ret = &FindScenesResultType{
			Count:    result.Count,
			Scenes:   scenes,
			Duration: result.TotalDuration,
			Filesize: result.TotalSize,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindScenesByPathRegex(ctx context.Context, filter *models.FindFilterType) (ret *FindScenesResultType, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {

		finalFilter := &models.SceneFilterType{}

		if filter != nil && filter.Q != nil {
			finalFilter.Path = &models.StringCriterionInput{
				Modifier: models.CriterionModifierMatchesRegex,
				Value:    "(?i)" + *filter.Q,
			}
		}

		// make a copy of the filter if provided, nilling out Q
		var queryFilter *models.FindFilterType
		if filter != nil {
			f := *filter
			queryFilter = &f
			queryFilter.Q = nil
		}

		fields := graphql.CollectAllFields(ctx)

		result, err := r.repository.Scene.Query(ctx, models.SceneQueryOptions{
			QueryOptions: models.QueryOptions{
				FindFilter: queryFilter,
				Count:      slices.Contains(fields, "count"),
			},
			SceneFilter:   finalFilter,
			TotalDuration: slices.Contains(fields, "duration"),
			TotalSize:     slices.Contains(fields, "filesize"),
		})
		if err != nil {
			return err
		}

		scenes, err := result.Resolve(ctx)
		if err != nil {
			return err
		}

		ret = &FindScenesResultType{
			Count:    result.Count,
			Scenes:   scenes,
			Duration: result.TotalDuration,
			Filesize: result.TotalSize,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) ParseSceneFilenames(ctx context.Context, filter *models.FindFilterType, config models.SceneParserInput) (ret *SceneParserResultType, err error) {
	repo := scene.NewFilenameParserRepository(r.repository)
	parser := scene.NewFilenameParser(filter, config, repo)

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		result, count, err := parser.Parse(ctx)

		if err != nil {
			return err
		}

		ret = &SceneParserResultType{
			Count:   count,
			Results: result,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindDuplicateScenes(ctx context.Context, distance *int, durationDiff *float64) (ret [][]*models.Scene, err error) {
	dist := 0
	durDiff := -1.
	if distance != nil {
		dist = *distance
	}
	if durationDiff != nil {
		durDiff = *durationDiff
	}
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Scene.FindDuplicates(ctx, dist, durDiff)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllScenes(ctx context.Context) (ret []*models.Scene, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Scene.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
