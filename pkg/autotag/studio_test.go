package autotag

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
)

func TestStudioScenes(t *testing.T) {
	mockSceneReader := &mocks.SceneReaderWriter{}

	const studioID = 2
	const studioName = "studio name"

	var scenes []*models.Scene
	matchingPaths, falsePaths := generateScenePaths(studioName)
	for i, p := range append(matchingPaths, falsePaths...) {
		scenes = append(scenes, &models.Scene{
			ID:   i + 1,
			Path: p,
		})
	}

	studio := models.Studio{
		ID:   studioID,
		Name: models.NullString(studioName),
	}

	organized := false
	perPage := 0

	expectedSceneFilter := &models.SceneFilterType{
		Organized: &organized,
		Path: &models.StringCriterionInput{
			Value:    `(?i)(?:^|_|[^\w\d])studio[.\-_ ]*name(?:$|_|[^\w\d])`,
			Modifier: models.CriterionModifierMatchesRegex,
		},
	}

	expectedFindFilter := &models.FindFilterType{
		PerPage: &perPage,
	}

	mockSceneReader.On("Query", expectedSceneFilter, expectedFindFilter).Return(scenes, len(scenes), nil).Once()

	for i := range matchingPaths {
		sceneID := i + 1
		expectedStudioID := models.NullInt64(studioID)
		mockSceneReader.On("Update", models.ScenePartial{
			ID:       sceneID,
			StudioID: &expectedStudioID,
		}).Return(nil, nil).Once()
	}

	err := StudioScenes(&studio, mockSceneReader)

	assert := assert.New(t)

	assert.Nil(err)
	mockSceneReader.AssertExpectations(t)
}
