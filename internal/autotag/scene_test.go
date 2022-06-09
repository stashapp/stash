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

const sceneExt = "mp4"

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

func generateNamePatterns(name, separator, ext string) []string {
	var ret []string
	ret = append(ret, fmt.Sprintf("%s%saaa.%s", name, separator, ext))
	ret = append(ret, fmt.Sprintf("aaa%s%s.%s", separator, name, ext))
	ret = append(ret, fmt.Sprintf("aaa%s%s%sbbb.%s", separator, name, separator, ext))
	ret = append(ret, fmt.Sprintf("dir/%s%saaa.%s", name, separator, ext))
	ret = append(ret, fmt.Sprintf("dir%sdir/%s%saaa.%s", separator, name, separator, ext))
	ret = append(ret, fmt.Sprintf("dir\\%s%saaa.%s", name, separator, ext))
	ret = append(ret, fmt.Sprintf("%s%saaa/dir/bbb.%s", name, separator, ext))
	ret = append(ret, fmt.Sprintf("%s%saaa\\dir\\bbb.%s", name, separator, ext))
	ret = append(ret, fmt.Sprintf("dir/%s%s/aaa.%s", name, separator, ext))
	ret = append(ret, fmt.Sprintf("dir\\%s%s\\aaa.%s", name, separator, ext))

	return ret
}

func generateSplitNamePatterns(name, separator, ext string) []string {
	var ret []string
	splitted := strings.Split(name, " ")
	// only do this for names that are split into two
	if len(splitted) == 2 {
		ret = append(ret, fmt.Sprintf("%s%s%s.%s", splitted[0], separator, splitted[1], ext))
	}

	return ret
}

func generateFalseNamePatterns(name string, separator, ext string) []string {
	splitted := strings.Split(name, " ")

	var ret []string
	// only do this for names that are split into two
	if len(splitted) == 2 {
		ret = append(ret, fmt.Sprintf("%s%saaa%s%s.%s", splitted[0], separator, separator, splitted[1], ext))
	}

	return ret
}

func generateTestPaths(testName, ext string) (scenePatterns []string, falseScenePatterns []string) {
	separators := testSeparators
	separators = append(separators, testEndSeparators...)

	for _, separator := range separators {
		scenePatterns = append(scenePatterns, generateNamePatterns(testName, separator, ext)...)
		scenePatterns = append(scenePatterns, generateNamePatterns(strings.ToLower(testName), separator, ext)...)
		scenePatterns = append(scenePatterns, generateNamePatterns(strings.ReplaceAll(testName, " ", ""), separator, ext)...)
		falseScenePatterns = append(falseScenePatterns, generateFalseNamePatterns(testName, separator, ext)...)
	}

	// add test cases for intra-name separators
	for _, separator := range testSeparators {
		if separator != " " {
			scenePatterns = append(scenePatterns, generateNamePatterns(strings.ReplaceAll(testName, " ", separator), separator, ext)...)
		}
	}

	// add basic false scenarios
	falseScenePatterns = append(falseScenePatterns, fmt.Sprintf("aaa%s.%s", testName, ext))
	falseScenePatterns = append(falseScenePatterns, fmt.Sprintf("%saaa.%s", testName, ext))

	// add path separator false scenarios
	falseScenePatterns = append(falseScenePatterns, generateFalseNamePatterns(testName, "/", ext)...)
	falseScenePatterns = append(falseScenePatterns, generateFalseNamePatterns(testName, "\\", ext)...)

	// split patterns only valid for ._- and whitespace
	for _, separator := range testSeparators {
		scenePatterns = append(scenePatterns, generateSplitNamePatterns(testName, separator, ext)...)
	}

	// false patterns for other separators
	for _, separator := range testEndSeparators {
		falseScenePatterns = append(falseScenePatterns, generateSplitNamePatterns(testName, separator, ext)...)
	}

	return
}

type pathTestTable struct {
	Path    string
	Matches bool
}

func generateTestTable(testName, ext string) []pathTestTable {
	var ret []pathTestTable

	var scenePatterns []string
	var falseScenePatterns []string

	separators := testSeparators
	separators = append(separators, testEndSeparators...)

	for _, separator := range separators {
		scenePatterns = append(scenePatterns, generateNamePatterns(testName, separator, ext)...)
		scenePatterns = append(scenePatterns, generateNamePatterns(strings.ToLower(testName), separator, ext)...)
		falseScenePatterns = append(falseScenePatterns, generateFalseNamePatterns(testName, separator, ext)...)
	}

	for _, p := range scenePatterns {
		t := pathTestTable{
			Path:    p,
			Matches: true,
		}

		ret = append(ret, t)
	}

	for _, p := range falseScenePatterns {
		t := pathTestTable{
			Path:    p,
			Matches: false,
		}

		ret = append(ret, t)
	}

	return ret
}

func TestScenePerformers(t *testing.T) {
	t.Parallel()

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

	testTables := generateTestTable(performerName, sceneExt)

	assert := assert.New(t)

	for _, test := range testTables {
		mockPerformerReader := &mocks.PerformerReaderWriter{}
		mockSceneReader := &mocks.SceneReaderWriter{}

		mockPerformerReader.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		mockPerformerReader.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Performer{&performer, &reversedPerformer}, nil).Once()

		scene := models.Scene{
			ID:   sceneID,
			Path: test.Path,
		}

		if test.Matches {
			mockSceneReader.On("UpdatePartial", testCtx, sceneID, models.ScenePartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			}).Return(nil, nil).Once()
		}

		err := ScenePerformers(testCtx, &scene, mockSceneReader, mockPerformerReader, nil)

		assert.Nil(err)
		mockPerformerReader.AssertExpectations(t)
		mockSceneReader.AssertExpectations(t)
	}
}

func TestSceneStudios(t *testing.T) {
	t.Parallel()

	var (
		sceneID    = 1
		studioName = "studio name"
		studioID   = 2
	)
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

	testTables := generateTestTable(studioName, sceneExt)

	assert := assert.New(t)

	doTest := func(mockStudioReader *mocks.StudioReaderWriter, mockSceneReader *mocks.SceneReaderWriter, test pathTestTable) {
		if test.Matches {
			expectedStudioID := studioID
			mockSceneReader.On("UpdatePartial", testCtx, sceneID, models.ScenePartial{
				StudioID: models.NewOptionalInt(expectedStudioID),
			}).Return(nil, nil).Once()
		}

		scene := models.Scene{
			ID:   sceneID,
			Path: test.Path,
		}
		err := SceneStudios(testCtx, &scene, mockSceneReader, mockStudioReader, nil)

		assert.Nil(err)
		mockStudioReader.AssertExpectations(t)
		mockSceneReader.AssertExpectations(t)
	}

	for _, test := range testTables {
		mockStudioReader := &mocks.StudioReaderWriter{}
		mockSceneReader := &mocks.SceneReaderWriter{}

		mockStudioReader.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		mockStudioReader.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Studio{&studio, &reversedStudio}, nil).Once()
		mockStudioReader.On("GetAliases", testCtx, mock.Anything).Return([]string{}, nil).Maybe()

		doTest(mockStudioReader, mockSceneReader, test)
	}

	const unmatchedName = "unmatched"
	studio.Name.String = unmatchedName

	// test against aliases
	for _, test := range testTables {
		mockStudioReader := &mocks.StudioReaderWriter{}
		mockSceneReader := &mocks.SceneReaderWriter{}

		mockStudioReader.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		mockStudioReader.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Studio{&studio, &reversedStudio}, nil).Once()
		mockStudioReader.On("GetAliases", testCtx, studioID).Return([]string{
			studioName,
		}, nil).Once()
		mockStudioReader.On("GetAliases", testCtx, reversedStudioID).Return([]string{}, nil).Once()

		doTest(mockStudioReader, mockSceneReader, test)
	}
}

func TestSceneTags(t *testing.T) {
	t.Parallel()

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

	testTables := generateTestTable(tagName, sceneExt)

	assert := assert.New(t)

	doTest := func(mockTagReader *mocks.TagReaderWriter, mockSceneReader *mocks.SceneReaderWriter, test pathTestTable) {
		if test.Matches {
			mockSceneReader.On("UpdatePartial", testCtx, sceneID, models.ScenePartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			}).Return(nil, nil).Once()
		}

		scene := models.Scene{
			ID:   sceneID,
			Path: test.Path,
		}
		err := SceneTags(testCtx, &scene, mockSceneReader, mockTagReader, nil)

		assert.Nil(err)
		mockTagReader.AssertExpectations(t)
		mockSceneReader.AssertExpectations(t)
	}

	for _, test := range testTables {
		mockTagReader := &mocks.TagReaderWriter{}
		mockSceneReader := &mocks.SceneReaderWriter{}

		mockTagReader.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		mockTagReader.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Tag{&tag, &reversedTag}, nil).Once()
		mockTagReader.On("GetAliases", testCtx, mock.Anything).Return([]string{}, nil).Maybe()

		doTest(mockTagReader, mockSceneReader, test)
	}

	const unmatchedName = "unmatched"
	tag.Name = unmatchedName

	// test against aliases
	for _, test := range testTables {
		mockTagReader := &mocks.TagReaderWriter{}
		mockSceneReader := &mocks.SceneReaderWriter{}

		mockTagReader.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		mockTagReader.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Tag{&tag, &reversedTag}, nil).Once()
		mockTagReader.On("GetAliases", testCtx, tagID).Return([]string{
			tagName,
		}, nil).Once()
		mockTagReader.On("GetAliases", testCtx, reversedTagID).Return([]string{}, nil).Once()

		doTest(mockTagReader, mockSceneReader, test)
	}
}
