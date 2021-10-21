package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
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

func (r *queryResolver) FindSceneByHash(ctx context.Context, input models.SceneHashInput) (*models.Scene, error) {
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

func (r *queryResolver) FindScenes(ctx context.Context, sceneFilter *models.SceneFilterType, sceneIDs []int, filter *models.FindFilterType) (ret *models.FindScenesResultType, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		var scenes []*models.Scene
		var count int
		var err error

		if len(sceneIDs) > 0 {
			scenes, err = repo.Scene().FindMany(sceneIDs)
			if err == nil {
				count = len(scenes)
				for i := range scenes {
					fmt.Printf("i: %v\n", i)
				}
			}
		} else {
			scenes, count, err = scene.QueryWithCount(repo.Scene(), sceneFilter, filter)
		}

		if err != nil {
			return err
		}

		ret = &models.FindScenesResultType{
			Count:  count,
			Scenes: scenes,
			// Duration: duration,
			// Filesize: filesize,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindScenesByPathRegex(ctx context.Context, filter *models.FindFilterType) (ret *models.FindScenesResultType, err error) {
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

		scenes, total, err := scene.QueryWithCount(repo.Scene(), sceneFilter, queryFilter)
		if err != nil {
			return err
		}

		ret = &models.FindScenesResultType{
			Count:  total,
			Scenes: scenes,
			// Duration: duration,
			// Filesize: filesize,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) ParseSceneFilenames(ctx context.Context, filter *models.FindFilterType, config models.SceneParserInput) (ret *models.SceneParserResultType, err error) {
	parser := manager.NewSceneFilenameParser(filter, config)

	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		result, count, err := parser.Parse(repo)

		if err != nil {
			return err
		}

		ret = &models.SceneParserResultType{
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
