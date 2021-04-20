package autotag

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
)

func TestTagScenes(t *testing.T) {
	mockSceneReader := &mocks.SceneReaderWriter{}

	const tagID = 2
	const tagName = "tag name"

	var scenes []*models.Scene
	matchingPaths, falsePaths := generateScenePaths(tagName)
	for i, p := range append(matchingPaths, falsePaths...) {
		scenes = append(scenes, &models.Scene{
			ID:   i + 1,
			Path: p,
		})
	}

	tag := models.Tag{
		ID:   tagID,
		Name: tagName,
	}

	organized := false
	perPage := 0

	expectedSceneFilter := &models.SceneFilterType{
		Organized: &organized,
		Path: &models.StringCriterionInput{
			Value:    `(?i)(?:^|_|[^\w\d])tag[.\-_ ]*name(?:$|_|[^\w\d])`,
			Modifier: models.CriterionModifierMatchesRegex,
		},
	}

	expectedFindFilter := &models.FindFilterType{
		PerPage: &perPage,
	}

	mockSceneReader.On("Query", expectedSceneFilter, expectedFindFilter).Return(scenes, len(scenes), nil).Once()

	for i := range matchingPaths {
		sceneID := i + 1
		mockSceneReader.On("GetTagIDs", sceneID).Return(nil, nil).Once()
		mockSceneReader.On("UpdateTags", sceneID, []int{tagID}).Return(nil).Once()
	}

	err := TagScenes(&tag, nil, mockSceneReader)

	assert := assert.New(t)

	assert.Nil(err)
	mockSceneReader.AssertExpectations(t)
}
