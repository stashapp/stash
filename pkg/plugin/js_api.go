package plugin

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type jsAPI struct {
	ctx context.Context
	r   models.ResolverRoot
}

func (api *jsAPI) FindScenesByPathRegex(filter *models.FindFilterType) (*models.FindScenesResultType, error) {
	return api.r.Query().FindScenesByPathRegex(api.ctx, filter)
}

func (api *jsAPI) FindTags(tagFilter *models.TagFilterType, filter *models.FindFilterType) (*models.FindTagsResultType, error) {
	return api.r.Query().FindTags(api.ctx, tagFilter, filter)
}

func (api *jsAPI) AllTags() ([]*models.Tag, error) {
	return api.r.Query().AllTags(api.ctx)
}

func (api *jsAPI) TagCreate(input models.TagCreateInput) (*models.Tag, error) {
	return api.r.Mutation().TagCreate(api.ctx, input)
}

func (api *jsAPI) FindScenes(sceneFilter *models.SceneFilterType, sceneIds []int, filter *models.FindFilterType) (*models.FindScenesResultType, error) {
	return api.r.Query().FindScenes(api.ctx, sceneFilter, sceneIds, filter)
}
