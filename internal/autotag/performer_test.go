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

func TestPerformerScenes(t *testing.T) {
	t.Parallel()

	type test struct {
		performerName        string
		expectedRegex        string
		aliases              []string
		expectedAliasRegexes []string
	}

	performerNames := []test{
		{
			"performer name",
			`(?i)(?:^|_|[^\p{L}\d])performer[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
			nil,
			nil,
		},
		{
			"performer + name",
			`(?i)(?:^|_|[^\p{L}\d])performer[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
			[]string{"alias name"},
			[]string{`(?i)(?:^|_|[^\p{L}\d])alias[.\-_ ]*name(?:$|_|[^\p{L}\d])`},
		},
	}

	// trailing backslash tests only work where filepath separator is not backslash
	if filepath.Separator != '\\' {
		performerNames = append(performerNames, test{
			`performer + name\`,
			`(?i)(?:^|_|[^\p{L}\d])performer[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\p{L}\d])`,
			nil,
			nil,
		})
	}

	for _, p := range performerNames {
		testPerformerScenes(t, p.performerName, p.expectedRegex, p.aliases, p.expectedAliasRegexes, true)
		testPerformerScenes(t, p.performerName, p.expectedRegex, p.aliases, p.expectedAliasRegexes, false)
	}
}

func testPerformerScenes(t *testing.T, performerName, expectedRegex string, aliases []string, expectedAliasRegexes []string, matchAlias bool) {
	db := mocks.NewDatabase()

	const performerID = 2

	var scenes []*models.Scene
	matchingPaths, falsePaths := generateTestPaths(performerName, "mp4")
	for i, p := range append(matchingPaths, falsePaths...) {
		scenes = append(scenes, &models.Scene{
			ID:           i + 1,
			Path:         p,
			PerformerIDs: models.NewRelatedIDs([]int{}),
		})
	}

	if aliases == nil {
		aliases = []string{}
	}

	performer := models.Performer{
		ID:      performerID,
		Name:    performerName,
		Aliases: models.NewRelatedStrings(aliases),
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

	db.Scene.On("Query", mock.Anything, scene.QueryOptions(expectedSceneFilter, expectedFindFilter, false)).
		Return(mocks.SceneQueryResult(scenes, len(scenes)), nil).Once()

	if matchAlias {
		for _, aliasRegex := range expectedAliasRegexes {
			expectedAliasFilter := &models.SceneFilterType{
				Organized: &organized,
				Path: &models.StringCriterionInput{
					Value:    aliasRegex,
					Modifier: models.CriterionModifierMatchesRegex,
				},
			}
			db.Scene.On("Query", mock.Anything, scene.QueryOptions(expectedAliasFilter, expectedFindFilter, false)).
				Return(mocks.SceneQueryResult([]*models.Scene{}, 0), nil).Once()
		}
	}

	for i := range matchingPaths {
		sceneID := i + 1

		matchPartial := mock.MatchedBy(func(got models.ScenePartial) bool {
			expected := models.ScenePartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			}

			return scenePartialsEqual(got, expected)
		})
		db.Scene.On("UpdatePartial", mock.Anything, sceneID, matchPartial).Return(nil, nil).Once()
	}

	tagger := Tagger{
		TxnManager: db,
	}

	err := tagger.PerformerScenes(testCtx, &performer, nil, db.Scene, matchAlias)

	assert := assert.New(t)

	assert.Nil(err)
	db.AssertExpectations(t)
}

func TestPerformerImages(t *testing.T) {
	t.Parallel()

	type test struct {
		performerName        string
		expectedRegex        string
		aliases              []string
		expectedAliasRegexes []string
	}

	performerNames := []test{
		{
			"performer name",
			`(?i)(?:^|_|[^\p{L}\d])performer[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
			nil,
			nil,
		},
		{
			"performer + name",
			`(?i)(?:^|_|[^\p{L}\d])performer[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
			[]string{"alias name"},
			[]string{`(?i)(?:^|_|[^\p{L}\d])alias[.\-_ ]*name(?:$|_|[^\p{L}\d])`},
		},
	}

	for _, p := range performerNames {
		testPerformerImages(t, p.performerName, p.expectedRegex, p.aliases, p.expectedAliasRegexes, true)
		testPerformerImages(t, p.performerName, p.expectedRegex, p.aliases, p.expectedAliasRegexes, false)
	}
}

func testPerformerImages(t *testing.T, performerName, expectedRegex string, aliases []string, expectedAliasRegexes []string, matchAlias bool) {
	db := mocks.NewDatabase()

	const performerID = 2

	var images []*models.Image
	matchingPaths, falsePaths := generateTestPaths(performerName, imageExt)
	for i, p := range append(matchingPaths, falsePaths...) {
		images = append(images, &models.Image{
			ID:           i + 1,
			Path:         p,
			PerformerIDs: models.NewRelatedIDs([]int{}),
		})
	}

	if aliases == nil {
		aliases = []string{}
	}

	performer := models.Performer{
		ID:      performerID,
		Name:    performerName,
		Aliases: models.NewRelatedStrings(aliases),
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

	db.Image.On("Query", mock.Anything, image.QueryOptions(expectedImageFilter, expectedFindFilter, false)).
		Return(mocks.ImageQueryResult(images, len(images)), nil).Once()

	if matchAlias {
		for _, aliasRegex := range expectedAliasRegexes {
			expectedAliasFilter := &models.ImageFilterType{
				Organized: &organized,
				Path: &models.StringCriterionInput{
					Value:    aliasRegex,
					Modifier: models.CriterionModifierMatchesRegex,
				},
			}
			db.Image.On("Query", mock.Anything, image.QueryOptions(expectedAliasFilter, expectedFindFilter, false)).
				Return(mocks.ImageQueryResult([]*models.Image{}, 0), nil).Once()
		}
	}

	for i := range matchingPaths {
		imageID := i + 1

		matchPartial := mock.MatchedBy(func(got models.ImagePartial) bool {
			expected := models.ImagePartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			}

			return imagePartialsEqual(got, expected)
		})
		db.Image.On("UpdatePartial", mock.Anything, imageID, matchPartial).Return(nil, nil).Once()
	}

	tagger := Tagger{
		TxnManager: db,
	}

	err := tagger.PerformerImages(testCtx, &performer, nil, db.Image, matchAlias)

	assert := assert.New(t)

	assert.Nil(err)
	db.AssertExpectations(t)
}

func TestPerformerGalleries(t *testing.T) {
	t.Parallel()

	type test struct {
		performerName        string
		expectedRegex        string
		aliases              []string
		expectedAliasRegexes []string
	}

	performerNames := []test{
		{
			"performer name",
			`(?i)(?:^|_|[^\p{L}\d])performer[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
			nil,
			nil,
		},
		{
			"performer + name",
			`(?i)(?:^|_|[^\p{L}\d])performer[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
			[]string{"alias name"},
			[]string{`(?i)(?:^|_|[^\p{L}\d])alias[.\-_ ]*name(?:$|_|[^\p{L}\d])`},
		},
	}

	for _, p := range performerNames {
		testPerformerGalleries(t, p.performerName, p.expectedRegex, p.aliases, p.expectedAliasRegexes, true)
		testPerformerGalleries(t, p.performerName, p.expectedRegex, p.aliases, p.expectedAliasRegexes, false)
	}
}

func testPerformerGalleries(t *testing.T, performerName, expectedRegex string, aliases []string, expectedAliasRegexes []string, matchAlias bool) {
	db := mocks.NewDatabase()

	const performerID = 2

	var galleries []*models.Gallery
	matchingPaths, falsePaths := generateTestPaths(performerName, galleryExt)
	for i, p := range append(matchingPaths, falsePaths...) {
		v := p
		galleries = append(galleries, &models.Gallery{
			ID:           i + 1,
			Path:         v,
			PerformerIDs: models.NewRelatedIDs([]int{}),
		})
	}

	if aliases == nil {
		aliases = []string{}
	}

	performer := models.Performer{
		ID:      performerID,
		Name:    performerName,
		Aliases: models.NewRelatedStrings(aliases),
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

	db.Gallery.On("Query", mock.Anything, expectedGalleryFilter, expectedFindFilter).Return(galleries, len(galleries), nil).Once()

	if matchAlias {
		for _, aliasRegex := range expectedAliasRegexes {
			expectedAliasFilter := &models.GalleryFilterType{
				Organized: &organized,
				Path: &models.StringCriterionInput{
					Value:    aliasRegex,
					Modifier: models.CriterionModifierMatchesRegex,
				},
			}
			db.Gallery.On("Query", mock.Anything, expectedAliasFilter, expectedFindFilter).
				Return([]*models.Gallery{}, 0, nil).Once()
		}
	}

	for i := range matchingPaths {
		galleryID := i + 1

		matchPartial := mock.MatchedBy(func(got models.GalleryPartial) bool {
			expected := models.GalleryPartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			}

			return galleryPartialsEqual(got, expected)
		})
		db.Gallery.On("UpdatePartial", mock.Anything, galleryID, matchPartial).Return(nil, nil).Once()
	}

	tagger := Tagger{
		TxnManager: db,
	}

	err := tagger.PerformerGalleries(testCtx, &performer, nil, db.Gallery, matchAlias)

	assert := assert.New(t)

	assert.Nil(err)
	db.AssertExpectations(t)
}
