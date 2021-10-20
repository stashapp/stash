package autotag

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
)

type testTagCase struct {
	tagName       string
	expectedRegex string
	aliasName     string
	aliasRegex    string
}

var testTagCases = []testTagCase{
	{
		"tag name",
		`(?i)(?:^|_|[^\w\d])tag[.\-_ ]*name(?:$|_|[^\w\d])`,
		"",
		"",
	},
	{
		"tag + name",
		`(?i)(?:^|_|[^\w\d])tag[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\w\d])`,
		"",
		"",
	},
	{
		`tag + name\`,
		`(?i)(?:^|_|[^\w\d])tag[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\w\d])`,
		"",
		"",
	},
	{
		"tag name",
		`(?i)(?:^|_|[^\w\d])tag[.\-_ ]*name(?:$|_|[^\w\d])`,
		"alias name",
		`(?i)(?:^|_|[^\w\d])alias[.\-_ ]*name(?:$|_|[^\w\d])`,
	},
	{
		"tag + name",
		`(?i)(?:^|_|[^\w\d])tag[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\w\d])`,
		"alias + name",
		`(?i)(?:^|_|[^\w\d])alias[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\w\d])`,
	},
	{
		`tag + name\`,
		`(?i)(?:^|_|[^\w\d])tag[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\w\d])`,
		`alias + name\`,
		`(?i)(?:^|_|[^\w\d])alias[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\w\d])`,
	},
}

func TestTagScenes(t *testing.T) {
	t.Parallel()

	for _, p := range testTagCases {
		testTagScenes(t, p)
	}
}

func testTagScenes(t *testing.T, tc testTagCase) {
	tagName := tc.tagName
	expectedRegex := tc.expectedRegex
	aliasName := tc.aliasName
	aliasRegex := tc.aliasRegex

	mockSceneReader := &mocks.SceneReaderWriter{}

	const tagID = 2

	var aliases []string

	testPathName := tagName
	if aliasName != "" {
		aliases = []string{aliasName}
		testPathName = aliasName
	}

	matchingPaths, falsePaths := generateTestPaths(testPathName, "mp4")

	var scenes []*models.Scene
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
	perPage := models.PerPageAll

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

	// if alias provided, then don't find by name
	onNameQuery := mockSceneReader.On("Query", expectedSceneFilter, expectedFindFilter)
	if aliasName == "" {
		onNameQuery.Return(scenes, len(scenes), nil).Once()
	} else {
		onNameQuery.Return(nil, 0, nil).Once()

		expectedAliasFilter := &models.SceneFilterType{
			Organized: &organized,
			Path: &models.StringCriterionInput{
				Value:    aliasRegex,
				Modifier: models.CriterionModifierMatchesRegex,
			},
		}

		mockSceneReader.On("Query", expectedAliasFilter, expectedFindFilter).Return(scenes, len(scenes), nil).Once()
	}

	for i := range matchingPaths {
		sceneID := i + 1
		mockSceneReader.On("GetTagIDs", sceneID).Return(nil, nil).Once()
		mockSceneReader.On("UpdateTags", sceneID, []int{tagID}).Return(nil).Once()
	}

	err := TagScenes(&tag, nil, aliases, mockSceneReader)

	assert := assert.New(t)

	assert.Nil(err)
	mockSceneReader.AssertExpectations(t)
}

func TestTagImages(t *testing.T) {
	t.Parallel()

	for _, p := range testTagCases {
		testTagImages(t, p)
	}
}

func testTagImages(t *testing.T, tc testTagCase) {
	tagName := tc.tagName
	expectedRegex := tc.expectedRegex
	aliasName := tc.aliasName
	aliasRegex := tc.aliasRegex

	mockImageReader := &mocks.ImageReaderWriter{}

	const tagID = 2

	var aliases []string

	testPathName := tagName
	if aliasName != "" {
		aliases = []string{aliasName}
		testPathName = aliasName
	}

	var images []*models.Image
	matchingPaths, falsePaths := generateTestPaths(testPathName, "mp4")
	for i, p := range append(matchingPaths, falsePaths...) {
		images = append(images, &models.Image{
			ID:   i + 1,
			Path: p,
		})
	}

	tag := models.Tag{
		ID:   tagID,
		Name: tagName,
	}

	organized := false
	perPage := models.PerPageAll

	expectedImageFilter := &models.ImageFilterType{
		Organized: &organized,
		Path: &models.StringCriterionInput{
			Value:    expectedRegex,
			Modifier: models.CriterionModifierMatchesRegex,
		},
	}

	expectedFindFilter := &models.FindFilterType{
		PerPage: &perPage,
	}

	// if alias provided, then don't find by name
	onNameQuery := mockImageReader.On("Query", expectedImageFilter, expectedFindFilter)
	if aliasName == "" {
		onNameQuery.Return(images, len(images), nil).Once()
	} else {
		onNameQuery.Return(nil, 0, nil).Once()

		expectedAliasFilter := &models.ImageFilterType{
			Organized: &organized,
			Path: &models.StringCriterionInput{
				Value:    aliasRegex,
				Modifier: models.CriterionModifierMatchesRegex,
			},
		}

		mockImageReader.On("Query", expectedAliasFilter, expectedFindFilter).Return(images, len(images), nil).Once()
	}

	for i := range matchingPaths {
		imageID := i + 1
		mockImageReader.On("GetTagIDs", imageID).Return(nil, nil).Once()
		mockImageReader.On("UpdateTags", imageID, []int{tagID}).Return(nil).Once()
	}

	err := TagImages(&tag, nil, aliases, mockImageReader)

	assert := assert.New(t)

	assert.Nil(err)
	mockImageReader.AssertExpectations(t)
}

func TestTagGalleries(t *testing.T) {
	t.Parallel()

	for _, p := range testTagCases {
		testTagGalleries(t, p)
	}
}

func testTagGalleries(t *testing.T, tc testTagCase) {
	tagName := tc.tagName
	expectedRegex := tc.expectedRegex
	aliasName := tc.aliasName
	aliasRegex := tc.aliasRegex

	mockGalleryReader := &mocks.GalleryReaderWriter{}

	const tagID = 2

	var aliases []string

	testPathName := tagName
	if aliasName != "" {
		aliases = []string{aliasName}
		testPathName = aliasName
	}

	var galleries []*models.Gallery
	matchingPaths, falsePaths := generateTestPaths(testPathName, "mp4")
	for i, p := range append(matchingPaths, falsePaths...) {
		galleries = append(galleries, &models.Gallery{
			ID:   i + 1,
			Path: models.NullString(p),
		})
	}

	tag := models.Tag{
		ID:   tagID,
		Name: tagName,
	}

	organized := false
	perPage := models.PerPageAll

	expectedGalleryFilter := &models.GalleryFilterType{
		Organized: &organized,
		Path: &models.StringCriterionInput{
			Value:    expectedRegex,
			Modifier: models.CriterionModifierMatchesRegex,
		},
	}

	expectedFindFilter := &models.FindFilterType{
		PerPage: &perPage,
	}

	// if alias provided, then don't find by name
	onNameQuery := mockGalleryReader.On("Query", expectedGalleryFilter, expectedFindFilter)
	if aliasName == "" {
		onNameQuery.Return(galleries, len(galleries), nil).Once()
	} else {
		onNameQuery.Return(nil, 0, nil).Once()

		expectedAliasFilter := &models.GalleryFilterType{
			Organized: &organized,
			Path: &models.StringCriterionInput{
				Value:    aliasRegex,
				Modifier: models.CriterionModifierMatchesRegex,
			},
		}

		mockGalleryReader.On("Query", expectedAliasFilter, expectedFindFilter).Return(galleries, len(galleries), nil).Once()
	}

	for i := range matchingPaths {
		galleryID := i + 1
		mockGalleryReader.On("GetTagIDs", galleryID).Return(nil, nil).Once()
		mockGalleryReader.On("UpdateTags", galleryID, []int{tagID}).Return(nil).Once()
	}

	err := TagGalleries(&tag, nil, aliases, mockGalleryReader)

	assert := assert.New(t)

	assert.Nil(err)
	mockGalleryReader.AssertExpectations(t)
}
