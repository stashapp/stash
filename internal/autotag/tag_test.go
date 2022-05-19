package autotag

import (
	"testing"

	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/scene"
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
		`(?i)(?:^|_|[^\p{L}\d])tag[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
		"",
		"",
	},
	{
		"tag + name",
		`(?i)(?:^|_|[^\p{L}\d])tag[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
		"",
		"",
	},
	{
		`tag + name\`,
		`(?i)(?:^|_|[^\p{L}\d])tag[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\p{L}\d])`,
		"",
		"",
	},
	{
		"tag name",
		`(?i)(?:^|_|[^\p{L}\d])tag[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
		"alias name",
		`(?i)(?:^|_|[^\p{L}\d])alias[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
	},
	{
		"tag + name",
		`(?i)(?:^|_|[^\p{L}\d])tag[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
		"alias + name",
		`(?i)(?:^|_|[^\p{L}\d])alias[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
	},
	{
		`tag + name\`,
		`(?i)(?:^|_|[^\p{L}\d])tag[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\p{L}\d])`,
		`alias + name\`,
		`(?i)(?:^|_|[^\p{L}\d])alias[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\p{L}\d])`,
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
	onNameQuery := mockSceneReader.On("Query", testCtx, scene.QueryOptions(expectedSceneFilter, expectedFindFilter, false))
	if aliasName == "" {
		onNameQuery.Return(mocks.SceneQueryResult(scenes, len(scenes)), nil).Once()
	} else {
		onNameQuery.Return(mocks.SceneQueryResult(nil, 0), nil).Once()

		expectedAliasFilter := &models.SceneFilterType{
			Organized: &organized,
			Path: &models.StringCriterionInput{
				Value:    aliasRegex,
				Modifier: models.CriterionModifierMatchesRegex,
			},
		}

		mockSceneReader.On("Query", testCtx, scene.QueryOptions(expectedAliasFilter, expectedFindFilter, false)).
			Return(mocks.SceneQueryResult(scenes, len(scenes)), nil).Once()
	}

	for i := range matchingPaths {
		sceneID := i + 1
		mockSceneReader.On("GetTagIDs", testCtx, sceneID).Return(nil, nil).Once()
		mockSceneReader.On("UpdateTags", testCtx, sceneID, []int{tagID}).Return(nil).Once()
	}

	err := TagScenes(testCtx, &tag, nil, aliases, mockSceneReader, nil)

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
	onNameQuery := mockImageReader.On("Query", testCtx, image.QueryOptions(expectedImageFilter, expectedFindFilter, false))
	if aliasName == "" {
		onNameQuery.Return(mocks.ImageQueryResult(images, len(images)), nil).Once()
	} else {
		onNameQuery.Return(mocks.ImageQueryResult(nil, 0), nil).Once()

		expectedAliasFilter := &models.ImageFilterType{
			Organized: &organized,
			Path: &models.StringCriterionInput{
				Value:    aliasRegex,
				Modifier: models.CriterionModifierMatchesRegex,
			},
		}

		mockImageReader.On("Query", testCtx, image.QueryOptions(expectedAliasFilter, expectedFindFilter, false)).
			Return(mocks.ImageQueryResult(images, len(images)), nil).Once()
	}

	for i := range matchingPaths {
		imageID := i + 1
		mockImageReader.On("GetTagIDs", testCtx, imageID).Return(nil, nil).Once()
		mockImageReader.On("UpdateTags", testCtx, imageID, []int{tagID}).Return(nil).Once()
	}

	err := TagImages(testCtx, &tag, nil, aliases, mockImageReader, nil)

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
	onNameQuery := mockGalleryReader.On("Query", testCtx, expectedGalleryFilter, expectedFindFilter)
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

		mockGalleryReader.On("Query", testCtx, expectedAliasFilter, expectedFindFilter).Return(galleries, len(galleries), nil).Once()
	}

	for i := range matchingPaths {
		galleryID := i + 1
		mockGalleryReader.On("GetTagIDs", testCtx, galleryID).Return(nil, nil).Once()
		mockGalleryReader.On("UpdateTags", testCtx, galleryID, []int{tagID}).Return(nil).Once()
	}

	err := TagGalleries(testCtx, &tag, nil, aliases, mockGalleryReader, nil)

	assert := assert.New(t)

	assert.Nil(err)
	mockGalleryReader.AssertExpectations(t)
}
