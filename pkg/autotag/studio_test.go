package autotag

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
)

func TestStudioScenes(t *testing.T) {
	type test struct {
		studioName    string
		expectedRegex string
	}

	studioNames := []test{
		{
			"studio name",
			`(?i)(?:^|_|[^\w\d])studio[.\-_ ]*name(?:$|_|[^\w\d])`,
		},
		{
			"studio + name",
			`(?i)(?:^|_|[^\w\d])studio[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\w\d])`,
		},
	}

	for _, p := range studioNames {
		testStudioScenes(t, p.studioName, p.expectedRegex)
	}
}

func testStudioScenes(t *testing.T, studioName, expectedRegex string) {
	mockSceneReader := &mocks.SceneReaderWriter{}

	const studioID = 2

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
		mockSceneReader.On("Find", sceneID).Return(&models.Scene{}, nil).Once()
		expectedStudioID := models.NullInt64(studioID)
		mockSceneReader.On("Update", models.ScenePartial{
			ID:       sceneID,
			StudioID: &expectedStudioID,
		}).Return(nil, nil).Once()
	}

	err := StudioScenes(&studio, nil, mockSceneReader)

	assert := assert.New(t)

	assert.Nil(err)
	mockSceneReader.AssertExpectations(t)
}
