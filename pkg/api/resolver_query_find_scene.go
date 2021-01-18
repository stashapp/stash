package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
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

func (r *queryResolver) FindScenes(ctx context.Context, sceneFilter *models.SceneFilterType, sceneIds []int, filter *models.FindFilterType) (ret *models.FindScenesResultType, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		scenes, total, err := repo.Scene().Query(sceneFilter, filter)
		if err != nil {
			return err
		}
		ret = &models.FindScenesResultType{
			Count:  total,
			Scenes: scenes,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindScenesByPathRegex(ctx context.Context, filter *models.FindFilterType) (ret *models.FindScenesResultType, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		scenes, total, err := repo.Scene().QueryByPathRegex(filter)
		if err != nil {
			return err
		}

		ret = &models.FindScenesResultType{
			Count:  total,
			Scenes: scenes,
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
