package autotag

import (
	"path/filepath"
	"testing"

	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testTagCase struct {
	tagName       string
	expectedRegex string
	aliasName     string
	aliasRegex    string
}

var (
	testTagCases = []testTagCase{
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
	}

	trailingBackslashCases = []testTagCase{
		{
			`tag + name\`,
			`(?i)(?:^|_|[^\p{L}\d])tag[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\p{L}\d])`,
			"",
			"",
		},
		{
			`tag + name\`,
			`(?i)(?:^|_|[^\p{L}\d])tag[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\p{L}\d])`,
			`alias + name\`,
			`(?i)(?:^|_|[^\p{L}\d])alias[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\p{L}\d])`,
		},
	}
)

func TestTagScenes(t *testing.T) {
	t.Parallel()

	tc := testTagCases
	// trailing backslash tests only work where filepath separator is not backslash
	if filepath.Separator != '\\' {
		tc = append(tc, trailingBackslashCases...)
	}

	for _, p := range tc {
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
			ID:     i + 1,
			Path:   p,
			TagIDs: models.NewRelatedIDs([]int{}),
		})
	}

	tag := models.Tag{
		ID:   tagID,
		Name: tagName,
	}

	organized := false
	perPage := 1000
	sort := "id"
	direction := models.SortDirectionEnumAsc

	expectedSceneFilter := &models.SceneFilterType{
		Organized: &organized,
		Path: &models.StringCriterionInput{
			Value:    expectedRegex,
			Modifier: models.CriterionModifierMatchesRegex,
		},
	}

	expectedFindFilter := &models.FindFilterType{
		PerPage:   &perPage,
		Sort:      &sort,
		Direction: &direction,
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

		mockSceneReader.On("Query", mock.Anything, scene.QueryOptions(expectedAliasFilter, expectedFindFilter, false)).
			Return(mocks.SceneQueryResult(scenes, len(scenes)), nil).Once()
	}

	for i := range matchingPaths {
		sceneID := i + 1
		mockSceneReader.On("UpdatePartial", mock.Anything, sceneID, models.ScenePartial{
			TagIDs: &models.UpdateIDs{
				IDs:  []int{tagID},
				Mode: models.RelationshipUpdateModeAdd,
			},
		}).Return(nil, nil).Once()
	}

	tagger := Tagger{
		TxnManager: &mocks.TxnManager{},
	}

	err := tagger.TagScenes(testCtx, &tag, nil, aliases, mockSceneReader)

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
			ID:     i + 1,
			Path:   p,
			TagIDs: models.NewRelatedIDs([]int{}),
		})
	}

	tag := models.Tag{
		ID:   tagID,
		Name: tagName,
	}

	organized := false
	perPage := 1000
	sort := "id"
	direction := models.SortDirectionEnumAsc

	expectedImageFilter := &models.ImageFilterType{
		Organized: &organized,
		Path: &models.StringCriterionInput{
			Value:    expectedRegex,
			Modifier: models.CriterionModifierMatchesRegex,
		},
	}

	expectedFindFilter := &models.FindFilterType{
		PerPage:   &perPage,
		Sort:      &sort,
		Direction: &direction,
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

		mockImageReader.On("Query", mock.Anything, image.QueryOptions(expectedAliasFilter, expectedFindFilter, false)).
			Return(mocks.ImageQueryResult(images, len(images)), nil).Once()
	}

	for i := range matchingPaths {
		imageID := i + 1

		mockImageReader.On("UpdatePartial", mock.Anything, imageID, models.ImagePartial{
			TagIDs: &models.UpdateIDs{
				IDs:  []int{tagID},
				Mode: models.RelationshipUpdateModeAdd,
			},
		}).Return(nil, nil).Once()
	}

	tagger := Tagger{
		TxnManager: &mocks.TxnManager{},
	}

	err := tagger.TagImages(testCtx, &tag, nil, aliases, mockImageReader)

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
		v := p
		galleries = append(galleries, &models.Gallery{
			ID:     i + 1,
			Path:   v,
			TagIDs: models.NewRelatedIDs([]int{}),
		})
	}

	tag := models.Tag{
		ID:   tagID,
		Name: tagName,
	}

	organized := false
	perPage := 1000
	sort := "id"
	direction := models.SortDirectionEnumAsc

	expectedGalleryFilter := &models.GalleryFilterType{
		Organized: &organized,
		Path: &models.StringCriterionInput{
			Value:    expectedRegex,
			Modifier: models.CriterionModifierMatchesRegex,
		},
	}

	expectedFindFilter := &models.FindFilterType{
		PerPage:   &perPage,
		Sort:      &sort,
		Direction: &direction,
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

		mockGalleryReader.On("Query", mock.Anything, expectedAliasFilter, expectedFindFilter).Return(galleries, len(galleries), nil).Once()
	}

	for i := range matchingPaths {
		galleryID := i + 1

		mockGalleryReader.On("UpdatePartial", mock.Anything, galleryID, models.GalleryPartial{
			TagIDs: &models.UpdateIDs{
				IDs:  []int{tagID},
				Mode: models.RelationshipUpdateModeAdd,
			},
		}).Return(nil, nil).Once()

	}

	tagger := Tagger{
		TxnManager: &mocks.TxnManager{},
	}

	err := tagger.TagGalleries(testCtx, &tag, nil, aliases, mockGalleryReader)

	assert := assert.New(t)

	assert.Nil(err)
	mockGalleryReader.AssertExpectations(t)
}
