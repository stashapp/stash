package autotag

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
)

func TestTagScenes(t *testing.T) {
	type test struct {
		tagName       string
		expectedRegex string
	}

	tagNames := []test{
		{
			"tag name",
			`(?i)(?:^|_|[^\w\d])tag[.\-_ ]*name(?:$|_|[^\w\d])`,
		},
		{
			"tag + name",
			`(?i)(?:^|_|[^\w\d])tag[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\w\d])`,
		},
	}

	for _, p := range tagNames {
		testTagScenes(t, p.tagName, p.expectedRegex)
	}
}

func testTagScenes(t *testing.T, tagName, expectedRegex string) {
	mockSceneReader := &mocks.SceneReaderWriter{}

	const tagID = 2

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
			Value:    expectedRegex,
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
