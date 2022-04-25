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
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		qb := repo.Scene()
		var err error
		if id != nil {
			idInt, err := strconv.Atoi(*id)
			if err != nil {
				return err
			}
			scene, err = qb.Find(idInt)
			if err != nil {
				return err
			}
		} else if checksum != nil {
			scene, err = qb.FindByChecksum(*checksum)
		}

		return err
	}); err != nil {
		return nil, err
	}

	return scene, nil
}

func (r *queryResolver) FindSceneByHash(ctx context.Context, input SceneHashInput) (*models.Scene, error) {
	var scene *models.Scene

	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		qb := repo.Scene()
		var err error
		if input.Checksum != nil {
			scene, err = qb.FindByChecksum(*input.Checksum)
			if err != nil {
				return err
			}
		}

		if scene == nil && input.Oshash != nil {
			scene, err = qb.FindByOSHash(*input.Oshash)
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
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		var scenes []*models.Scene
		var err error

		fields := graphql.CollectAllFields(ctx)
		result := &models.SceneQueryResult{}

		if len(sceneIDs) > 0 {
			scenes, err = repo.Scene().FindMany(sceneIDs)
			if err == nil {
				result.Count = len(scenes)
				for _, s := range scenes {
					result.TotalDuration += s.Duration.Float64
					size, _ := strconv.ParseFloat(s.Size.String, 64)
					result.TotalSize += size
				}
			}
		} else {
			result, err = repo.Scene().Query(models.SceneQueryOptions{
				QueryOptions: models.QueryOptions{
					FindFilter: filter,
					Count:      stringslice.StrInclude(fields, "count"),
				},
				SceneFilter:   sceneFilter,
				TotalDuration: stringslice.StrInclude(fields, "duration"),
				TotalSize:     stringslice.StrInclude(fields, "filesize"),
			})
			if err == nil {
				scenes, err = result.Resolve()
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
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {

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

		result, err := repo.Scene().Query(models.SceneQueryOptions{
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

		scenes, err := result.Resolve()
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

	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		result, count, err := parser.Parse(repo)

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
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.Scene().FindDuplicates(dist)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
