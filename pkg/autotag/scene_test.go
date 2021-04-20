package autotag

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testSeparators = []string{
	".",
	"-",
	"_",
	" ",
}

var testEndSeparators = []string{
	"{",
	"}",
	"(",
	")",
	",",
}

func generateNamePatterns(name, separator string) []string {
	var ret []string
	ret = append(ret, fmt.Sprintf("%s%saaa.mp4", name, separator))
	ret = append(ret, fmt.Sprintf("aaa%s%s.mp4", separator, name))
	ret = append(ret, fmt.Sprintf("aaa%s%s%sbbb.mp4", separator, name, separator))
	ret = append(ret, fmt.Sprintf("dir/%s%saaa.mp4", name, separator))
	ret = append(ret, fmt.Sprintf("dir\\%s%saaa.mp4", name, separator))
	ret = append(ret, fmt.Sprintf("%s%saaa/dir/bbb.mp4", name, separator))
	ret = append(ret, fmt.Sprintf("%s%saaa\\dir\\bbb.mp4", name, separator))
	ret = append(ret, fmt.Sprintf("dir/%s%s/aaa.mp4", name, separator))
	ret = append(ret, fmt.Sprintf("dir\\%s%s\\aaa.mp4", name, separator))

	return ret
}

func generateFalseNamePattern(name string, separator string) string {
	splitted := strings.Split(name, " ")

	return fmt.Sprintf("%s%saaa%s%s.mp4", splitted[0], separator, separator, splitted[1])
}

func generateScenePaths(testName string) (scenePatterns []string, falseScenePatterns []string) {
	separators := append(testSeparators, testEndSeparators...)

	for _, separator := range separators {
		scenePatterns = append(scenePatterns, generateNamePatterns(testName, separator)...)
		scenePatterns = append(scenePatterns, generateNamePatterns(strings.ToLower(testName), separator)...)
		falseScenePatterns = append(falseScenePatterns, generateFalseNamePattern(testName, separator))
	}

	return
}

type pathTestTable struct {
	ScenePath string
	Matches   bool
}

func generateTestTable(testName string) []pathTestTable {
	var ret []pathTestTable

	var scenePatterns []string
	var falseScenePatterns []string

	separators := append(testSeparators, testEndSeparators...)

	for _, separator := range separators {
		scenePatterns = append(scenePatterns, generateNamePatterns(testName, separator)...)
		scenePatterns = append(scenePatterns, generateNamePatterns(strings.ToLower(testName), separator)...)
		falseScenePatterns = append(falseScenePatterns, generateFalseNamePattern(testName, separator))
	}

	for _, p := range scenePatterns {
		t := pathTestTable{
			ScenePath: p,
			Matches:   true,
		}

		ret = append(ret, t)
	}

	for _, p := range falseScenePatterns {
		t := pathTestTable{
			ScenePath: p,
			Matches:   false,
		}

		ret = append(ret, t)
	}

	return ret
}

func TestScenePerformers(t *testing.T) {
	const sceneID = 1
	const performerName = "performer name"
	const performerID = 2
	performer := models.Performer{
		ID:   performerID,
		Name: models.NullString(performerName),
	}

	const reversedPerformerName = "name performer"
	const reversedPerformerID = 3
	reversedPerformer := models.Performer{
		ID:   reversedPerformerID,
		Name: models.NullString(reversedPerformerName),
	}

	testTables := generateTestTable(performerName)

	assert := assert.New(t)

	for _, test := range testTables {
		mockPerformerReader := &mocks.PerformerReaderWriter{}
		mockSceneReader := &mocks.SceneReaderWriter{}

		mockPerformerReader.On("QueryForAutoTag", mock.Anything).Return([]*models.Performer{&performer, &reversedPerformer}, nil).Once()

		if test.Matches {
			mockSceneReader.On("GetPerformerIDs", sceneID).Return(nil, nil).Once()
			mockSceneReader.On("UpdatePerformers", sceneID, []int{performerID}).Return(nil).Once()
		}

		scene := models.Scene{
			ID:   sceneID,
			Path: test.ScenePath,
		}
		err := ScenePerformers(&scene, mockSceneReader, mockPerformerReader)

		assert.Nil(err)
		mockPerformerReader.AssertExpectations(t)
		mockSceneReader.AssertExpectations(t)
	}
}

func TestSceneStudios(t *testing.T) {
	const sceneID = 1
	const studioName = "studio name"
	const studioID = 2
	studio := models.Studio{
		ID:   studioID,
		Name: models.NullString(studioName),
	}

	const reversedStudioName = "name studio"
	const reversedStudioID = 3
	reversedStudio := models.Studio{
		ID:   reversedStudioID,
		Name: models.NullString(reversedStudioName),
	}

	testTables := generateTestTable(studioName)

	assert := assert.New(t)

	for _, test := range testTables {
		mockStudioReader := &mocks.StudioReaderWriter{}
		mockSceneReader := &mocks.SceneReaderWriter{}

		mockStudioReader.On("QueryForAutoTag", mock.Anything).Return([]*models.Studio{&studio, &reversedStudio}, nil).Once()

		if test.Matches {
			mockSceneReader.On("Find", sceneID).Return(&models.Scene{}, nil).Once()
			expectedStudioID := models.NullInt64(studioID)
			mockSceneReader.On("Update", models.ScenePartial{
				ID:       sceneID,
				StudioID: &expectedStudioID,
			}).Return(nil, nil).Once()
		}

		scene := models.Scene{
			ID:   sceneID,
			Path: test.ScenePath,
		}
		err := SceneStudios(&scene, mockSceneReader, mockStudioReader)

		assert.Nil(err)
		mockStudioReader.AssertExpectations(t)
		mockSceneReader.AssertExpectations(t)
	}
}

func TestSceneTags(t *testing.T) {
	const sceneID = 1
	const tagName = "tag name"
	const tagID = 2
	tag := models.Tag{
		ID:   tagID,
		Name: tagName,
	}

	const reversedTagName = "name tag"
	const reversedTagID = 3
	reversedTag := models.Tag{
		ID:   reversedTagID,
		Name: reversedTagName,
	}

	testTables := generateTestTable(tagName)

	assert := assert.New(t)

	for _, test := range testTables {
		mockTagReader := &mocks.TagReaderWriter{}
		mockSceneReader := &mocks.SceneReaderWriter{}

		mockTagReader.On("QueryForAutoTag", mock.Anything).Return([]*models.Tag{&tag, &reversedTag}, nil).Once()

		if test.Matches {
			mockSceneReader.On("GetTagIDs", sceneID).Return(nil, nil).Once()
			mockSceneReader.On("UpdateTags", sceneID, []int{tagID}).Return(nil).Once()
		}

		scene := models.Scene{
			ID:   sceneID,
			Path: test.ScenePath,
		}
		err := SceneTags(&scene, mockSceneReader, mockTagReader)

		assert.Nil(err)
		mockTagReader.AssertExpectations(t)
		mockSceneReader.AssertExpectations(t)
	}
}
