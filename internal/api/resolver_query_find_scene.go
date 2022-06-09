package api

import (
	"context"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

func (r *queryResolver) FindScene(ctx context.Context, id *string, checksum *string) (*models.Scene, error) {
	var scene *models.Scene
	if err := r.withTxn(ctx, func(ctx context.Context) error {
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
			scene, err = qb.FindByChecksum(ctx, *checksum)
		}

		return err
	}); err != nil {
		return nil, err
	}

	return scene, nil
}

func (r *queryResolver) FindSceneByHash(ctx context.Context, input SceneHashInput) (*models.Scene, error) {
	var scene *models.Scene

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Scene
		var err error
		if input.Checksum != nil {
			scene, err = qb.FindByChecksum(ctx, *input.Checksum)
			if err != nil {
				return err
			}
		}

		if scene == nil && input.Oshash != nil {
			scene, err = qb.FindByOSHash(ctx, *input.Oshash)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return scene, nil
}

func (r *queryResolver) FindScenes(ctx context.Context, sceneFilter *models.SceneFilterType, sceneIDs []int, filter *models.FindFilterType) (ret *FindScenesResultType, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		var scenes []*models.Scene
		var err error

		fields := graphql.CollectAllFields(ctx)
		result := &models.SceneQueryResult{}

		if len(sceneIDs) > 0 {
			scenes, err = r.repository.Scene.FindMany(ctx, sceneIDs)
			if err == nil {
				result.Count = len(scenes)
				for _, s := range scenes {
					result.TotalDuration += float64PtrToFloat64(s.Duration)

					size, _ := strconv.ParseFloat(stringPtrToString(s.Size), 64)
					result.TotalSize += size
				}
			}
		} else {
			result, err = r.repository.Scene.Query(ctx, models.SceneQueryOptions{
				QueryOptions: models.QueryOptions{
					FindFilter: filter,
					Count:      stringslice.StrInclude(fields, "count"),
				},
				SceneFilter:   sceneFilter,
				TotalDuration: stringslice.StrInclude(fields, "duration"),
				TotalSize:     stringslice.StrInclude(fields, "filesize"),
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
	if err := r.withTxn(ctx, func(ctx context.Context) error {

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
				Count:      stringslice.StrInclude(fields, "count"),
			},
			SceneFilter:   sceneFilter,
			TotalDuration: stringslice.StrInclude(fields, "duration"),
			TotalSize:     stringslice.StrInclude(fields, "filesize"),
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

func (r *queryResolver) ParseSceneFilenames(ctx context.Context, filter *models.FindFilterType, config manager.SceneParserInput) (ret *SceneParserResultType, err error) {
	parser := manager.NewSceneFilenameParser(filter, config)

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		result, count, err := parser.Parse(ctx, manager.SceneFilenameParserRepository{
			Scene:     r.repository.Scene,
			Performer: r.repository.Performer,
			Studio:    r.repository.Studio,
			Movie:     r.repository.Movie,
			Tag:       r.repository.Tag,
		})

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

func (r *queryResolver) FindDuplicateScenes(ctx context.Context, distance *int) (ret [][]*models.Scene, err error) {
	dist := 0
	if distance != nil {
		dist = *distance
	}
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Scene.FindDuplicates(ctx, dist)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
