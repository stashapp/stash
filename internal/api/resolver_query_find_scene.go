package api

import (
	"context"
	"slices"
	"strconv"

	"github.com/99designs/gqlgen/graphql"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
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
	sceneIDs []int,
	ids []string,
	filter *models.FindFilterType,
) (ret *FindScenesResultType, err error) {
	if len(ids) > 0 {
		sceneIDs, err = stringslice.StringSliceToIntSlice(ids)
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
				SceneFilter:   sceneFilter,
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

		sceneFilter := &models.SceneFilterType{}

		if filter != nil && filter.Q != nil {
			sceneFilter.Path = &models.StringCriterionInput{
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
			SceneFilter:   sceneFilter,
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
