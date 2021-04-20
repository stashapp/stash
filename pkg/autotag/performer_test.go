package autotag

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPerformerScenes(t *testing.T) {
	mockSceneReader := &mocks.SceneReaderWriter{}

	const performerID = 2
	const performerName = "performer name"

	var scenes []*models.Scene
	matchingPaths, falsePaths := generateScenePaths(performerName)
	for i, p := range append(matchingPaths, falsePaths...) {
		scenes = append(scenes, &models.Scene{
			ID:   i + 1,
			Path: p,
		})
	}

	performer := models.Performer{
		ID:   performerID,
		Name: models.NullString(performerName),
	}

	organized := false
	perPage := 0

	expectedSceneFilter := &models.SceneFilterType{
		Organized: &organized,
		Path: &models.StringCriterionInput{
			Value:    `(?i)(?:^|_|[^\w\d])performer[.\-_ ]*name(?:$|_|[^\w\d])`,
			Modifier: models.CriterionModifierMatchesRegex,
		},
	}

	expectedFindFilter := &models.FindFilterType{
		PerPage: &perPage,
	}

	mockSceneReader.On("Query", expectedSceneFilter, expectedFindFilter).Return(scenes, len(scenes), nil).Once()

	for i := range matchingPaths {
		sceneID := i + 1
		mockSceneReader.On("GetPerformerIDs", sceneID).Return(nil, nil).Once()
		mockSceneReader.On("UpdatePerformers", sceneID, []int{performerID}).Return(nil).Once()
	}

	err := PerformerScenes(&performer, nil, mockSceneReader)

	assert := assert.New(t)

	assert.Nil(err)
	mockSceneReader.AssertExpectations(t)
}
