package mocks

import "github.com/stashapp/stash/pkg/models"

type sceneResolver struct {
	scenes []*models.Scene
}

func (s *sceneResolver) Find(id int) (*models.Scene, error) {
	panic("not implemented")
}

func (s *sceneResolver) FindMany(ids []int) ([]*models.Scene, error) {
	return s.scenes, nil
}

func SceneQueryResult(scenes []*models.Scene, count int) *models.SceneQueryResult {
	ret := models.NewSceneQueryResult(&sceneResolver{
		scenes: scenes,
	})

	ret.Count = count
	return ret
}
