package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

func (r *queryResolver) FindScene(ctx context.Context, id *string, checksum *string) (*models.Scene, error) {
	qb := sqlite.NewSceneQueryBuilder()
	var scene *models.Scene
	var err error
	if id != nil {
		idInt, _ := strconv.Atoi(*id)
		scene, err = qb.Find(idInt)
	} else if checksum != nil {
		scene, err = qb.FindByChecksum(*checksum)
	}
	return scene, err
}

func (r *queryResolver) FindSceneByHash(ctx context.Context, input models.SceneHashInput) (*models.Scene, error) {
	qb := sqlite.NewSceneQueryBuilder()
	var scene *models.Scene
	var err error

	if input.Checksum != nil {
		scene, err = qb.FindByChecksum(*input.Checksum)
		if err != nil {
			return nil, err
		}
	}

	if scene == nil && input.Oshash != nil {
		scene, err = qb.FindByOSHash(*input.Oshash)
		if err != nil {
			return nil, err
		}
	}

	return scene, err
}

func (r *queryResolver) FindScenes(ctx context.Context, sceneFilter *models.SceneFilterType, sceneIds []int, filter *models.FindFilterType) (*models.FindScenesResultType, error) {
	qb := sqlite.NewSceneQueryBuilder()
	scenes, total := qb.Query(sceneFilter, filter)
	return &models.FindScenesResultType{
		Count:  total,
		Scenes: scenes,
	}, nil
}

func (r *queryResolver) FindScenesByPathRegex(ctx context.Context, filter *models.FindFilterType) (*models.FindScenesResultType, error) {
	qb := sqlite.NewSceneQueryBuilder()

	scenes, total := qb.QueryByPathRegex(filter)
	return &models.FindScenesResultType{
		Count:  total,
		Scenes: scenes,
	}, nil
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
