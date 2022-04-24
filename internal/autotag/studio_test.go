package autotag

import (
	"testing"

	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stretchr/testify/assert"
)

type testStudioCase struct {
	studioName    string
	expectedRegex string
	aliasName     string
	aliasRegex    string
}

var testStudioCases = []testStudioCase{
	{
		"studio name",
		`(?i)(?:^|_|[^\p{L}\d])studio[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
		"",
		"",
	},
	{
		"studio + name",
		`(?i)(?:^|_|[^\p{L}\d])studio[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
		"",
		"",
	},
	{
		`studio + name\`,
		`(?i)(?:^|_|[^\p{L}\d])studio[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\p{L}\d])`,
		"",
		"",
	},
	{
		"studio name",
		`(?i)(?:^|_|[^\p{L}\d])studio[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
		"alias name",
		`(?i)(?:^|_|[^\p{L}\d])alias[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
	},
	{
		"studio + name",
		`(?i)(?:^|_|[^\p{L}\d])studio[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
		"alias + name",
		`(?i)(?:^|_|[^\p{L}\d])alias[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
	},
	{
		`studio + name\`,
		`(?i)(?:^|_|[^\p{L}\d])studio[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\p{L}\d])`,
		`alias + name\`,
		`(?i)(?:^|_|[^\p{L}\d])alias[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\p{L}\d])`,
	},
}

func TestStudioScenes(t *testing.T) {
	t.Parallel()

	for _, p := range testStudioCases {
		testStudioScenes(t, p)
	}
}

func testStudioScenes(t *testing.T, tc testStudioCase) {
	studioName := tc.studioName
	expectedRegex := tc.expectedRegex
	aliasName := tc.aliasName
	aliasRegex := tc.aliasRegex

	mockSceneReader := &mocks.SceneReaderWriter{}

	var studioID = 2

	var aliases []string

	testPathName := studioName
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

	studio := models.Studio{
		ID:   studioID,
		Name: models.NullString(studioName),
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
		expectedStudioID := studioID
		mockSceneReader.On("UpdatePartial", testCtx, sceneID, models.ScenePartial{
			StudioID: models.NewOptionalInt(expectedStudioID),
		}).Return(nil, nil).Once()
	}

	err := StudioScenes(testCtx, &studio, nil, aliases, mockSceneReader, nil)

	assert := assert.New(t)

	assert.Nil(err)
	mockSceneReader.AssertExpectations(t)
}

func TestStudioImages(t *testing.T) {
	t.Parallel()

	for _, p := range testStudioCases {
		testStudioImages(t, p)
	}
}

func testStudioImages(t *testing.T, tc testStudioCase) {
	studioName := tc.studioName
	expectedRegex := tc.expectedRegex
	aliasName := tc.aliasName
	aliasRegex := tc.aliasRegex

	mockImageReader := &mocks.ImageReaderWriter{}

	var studioID = 2

	var aliases []string

	testPathName := studioName
	if aliasName != "" {
		aliases = []string{aliasName}
		testPathName = aliasName
	}

	var images []*models.Image
	matchingPaths, falsePaths := generateTestPaths(testPathName, imageExt)
	for i, p := range append(matchingPaths, falsePaths...) {
		images = append(images, &models.Image{
			ID:   i + 1,
			Path: p,
		})
	}

	studio := models.Studio{
		ID:   studioID,
		Name: models.NullString(studioName),
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
		expectedStudioID := studioID
		mockImageReader.On("UpdatePartial", testCtx, imageID, models.ImagePartial{
			StudioID: models.NewOptionalInt(expectedStudioID),
		}).Return(nil, nil).Once()
	}

	err := StudioImages(testCtx, &studio, nil, aliases, mockImageReader, nil)

	assert := assert.New(t)

	assert.Nil(err)
	mockImageReader.AssertExpectations(t)
}

func TestStudioGalleries(t *testing.T) {
	t.Parallel()

	for _, p := range testStudioCases {
		testStudioGalleries(t, p)
	}
}

func testStudioGalleries(t *testing.T, tc testStudioCase) {
	studioName := tc.studioName
	expectedRegex := tc.expectedRegex
	aliasName := tc.aliasName
	aliasRegex := tc.aliasRegex
	mockGalleryReader := &mocks.GalleryReaderWriter{}

	var studioID = 2

	var aliases []string

	testPathName := studioName
	if aliasName != "" {
		aliases = []string{aliasName}
		testPathName = aliasName
	}

	var galleries []*models.Gallery
	matchingPaths, falsePaths := generateTestPaths(testPathName, galleryExt)
	for i, p := range append(matchingPaths, falsePaths...) {
		v := p
		galleries = append(galleries, &models.Gallery{
			ID:   i + 1,
			Path: &v,
		})
	}

	studio := models.Studio{
		ID:   studioID,
		Name: models.NullString(studioName),
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
		expectedStudioID := studioID
		mockGalleryReader.On("UpdatePartial", testCtx, galleryID, models.GalleryPartial{
			StudioID: models.NewOptionalInt(expectedStudioID),
		}).Return(nil, nil).Once()
	}

	err := StudioGalleries(testCtx, &studio, nil, aliases, mockGalleryReader, nil)

	assert := assert.New(t)

	assert.Nil(err)
	mockGalleryReader.AssertExpectations(t)
}
