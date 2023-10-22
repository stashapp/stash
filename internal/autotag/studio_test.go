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

type testStudioCase struct {
	studioName    string
	expectedRegex string
	aliasName     string
	aliasRegex    string
}

var (
	testStudioCases = []testStudioCase{
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
	}

	trailingBackslashStudioCases = []testStudioCase{
		{
			`studio + name\`,
			`(?i)(?:^|_|[^\p{L}\d])studio[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\p{L}\d])`,
			"",
			"",
		},
		{
			`studio + name\`,
			`(?i)(?:^|_|[^\p{L}\d])studio[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\p{L}\d])`,
			`alias + name\`,
			`(?i)(?:^|_|[^\p{L}\d])alias[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\p{L}\d])`,
		},
	}
)

func TestStudioScenes(t *testing.T) {
	t.Parallel()

	tc := testStudioCases
	// trailing backslash tests only work where filepath separator is not backslash
	if filepath.Separator != '\\' {
		tc = append(tc, trailingBackslashStudioCases...)
	}

	for _, p := range tc {
		testStudioScenes(t, p)
	}
}

func testStudioScenes(t *testing.T, tc testStudioCase) {
	studioName := tc.studioName
	expectedRegex := tc.expectedRegex
	aliasName := tc.aliasName
	aliasRegex := tc.aliasRegex

	db := mocks.NewDatabase()

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
		Name: studioName,
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
	onNameQuery := db.Scene.On("Query", testCtx, scene.QueryOptions(expectedSceneFilter, expectedFindFilter, false))

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

		db.Scene.On("Query", mock.Anything, scene.QueryOptions(expectedAliasFilter, expectedFindFilter, false)).
			Return(mocks.SceneQueryResult(scenes, len(scenes)), nil).Once()
	}

	for i := range matchingPaths {
		sceneID := i + 1

		matchPartial := mock.MatchedBy(func(got models.ScenePartial) bool {
			expected := models.ScenePartial{
				StudioID: models.NewOptionalInt(studioID),
			}

			return scenePartialsEqual(got, expected)
		})
		db.Scene.On("UpdatePartial", mock.Anything, sceneID, matchPartial).Return(nil, nil).Once()
	}

	tagger := Tagger{
		TxnManager: db,
	}

	err := tagger.StudioScenes(testCtx, &studio, nil, aliases, db.Scene)

	assert := assert.New(t)

	assert.Nil(err)
	db.AssertExpectations(t)
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

	db := mocks.NewDatabase()

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
		Name: studioName,
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
	onNameQuery := db.Image.On("Query", mock.Anything, image.QueryOptions(expectedImageFilter, expectedFindFilter, false))
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

		db.Image.On("Query", mock.Anything, image.QueryOptions(expectedAliasFilter, expectedFindFilter, false)).
			Return(mocks.ImageQueryResult(images, len(images)), nil).Once()
	}

	for i := range matchingPaths {
		imageID := i + 1

		matchPartial := mock.MatchedBy(func(got models.ImagePartial) bool {
			expected := models.ImagePartial{
				StudioID: models.NewOptionalInt(studioID),
			}

			return imagePartialsEqual(got, expected)
		})
		db.Image.On("UpdatePartial", mock.Anything, imageID, matchPartial).Return(nil, nil).Once()
	}

	tagger := Tagger{
		TxnManager: db,
	}

	err := tagger.StudioImages(testCtx, &studio, nil, aliases, db.Image)

	assert := assert.New(t)

	assert.Nil(err)
	db.AssertExpectations(t)
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

	db := mocks.NewDatabase()

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
			Path: v,
		})
	}

	studio := models.Studio{
		ID:   studioID,
		Name: studioName,
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
	onNameQuery := db.Gallery.On("Query", mock.Anything, expectedGalleryFilter, expectedFindFilter)
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

		db.Gallery.On("Query", mock.Anything, expectedAliasFilter, expectedFindFilter).Return(galleries, len(galleries), nil).Once()
	}

	for i := range matchingPaths {
		galleryID := i + 1

		matchPartial := mock.MatchedBy(func(got models.GalleryPartial) bool {
			expected := models.GalleryPartial{
				StudioID: models.NewOptionalInt(studioID),
			}

			return galleryPartialsEqual(got, expected)
		})
		db.Gallery.On("UpdatePartial", mock.Anything, galleryID, matchPartial).Return(nil, nil).Once()
	}

	tagger := Tagger{
		TxnManager: db,
	}

	err := tagger.StudioGalleries(testCtx, &studio, nil, aliases, db.Gallery)

	assert := assert.New(t)

	assert.Nil(err)
	db.AssertExpectations(t)
}
