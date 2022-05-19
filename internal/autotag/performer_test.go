package autotag

import (
	"testing"

	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stretchr/testify/assert"
)

func TestPerformerScenes(t *testing.T) {
	t.Parallel()

	type test struct {
		performerName string
		expectedRegex string
	}

	performerNames := []test{
		{
			"performer name",
			`(?i)(?:^|_|[^\p{L}\d])performer[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
		},
		{
			"performer + name",
			`(?i)(?:^|_|[^\p{L}\d])performer[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
		},
		{
			`performer + name\`,
			`(?i)(?:^|_|[^\p{L}\d])performer[.\-_ ]*\+[.\-_ ]*name\\(?:$|_|[^\p{L}\d])`,
		},
	}

	for _, p := range performerNames {
		testPerformerScenes(t, p.performerName, p.expectedRegex)
	}
}

func testPerformerScenes(t *testing.T, performerName, expectedRegex string) {
	mockSceneReader := &mocks.SceneReaderWriter{}

	const performerID = 2

	var scenes []*models.Scene
	matchingPaths, falsePaths := generateTestPaths(performerName, "mp4")
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

	mockSceneReader.On("Query", testCtx, scene.QueryOptions(expectedSceneFilter, expectedFindFilter, false)).
		Return(mocks.SceneQueryResult(scenes, len(scenes)), nil).Once()

	for i := range matchingPaths {
		sceneID := i + 1
		mockSceneReader.On("GetPerformerIDs", testCtx, sceneID).Return(nil, nil).Once()
		mockSceneReader.On("UpdatePerformers", testCtx, sceneID, []int{performerID}).Return(nil).Once()
	}

	err := PerformerScenes(testCtx, &performer, nil, mockSceneReader, nil)

	assert := assert.New(t)

	assert.Nil(err)
	mockSceneReader.AssertExpectations(t)
}

func TestPerformerImages(t *testing.T) {
	t.Parallel()

	type test struct {
		performerName string
		expectedRegex string
	}

	performerNames := []test{
		{
			"performer name",
			`(?i)(?:^|_|[^\p{L}\d])performer[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
		},
		{
			"performer + name",
			`(?i)(?:^|_|[^\p{L}\d])performer[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
		},
	}

	for _, p := range performerNames {
		testPerformerImages(t, p.performerName, p.expectedRegex)
	}
}

func testPerformerImages(t *testing.T, performerName, expectedRegex string) {
	mockImageReader := &mocks.ImageReaderWriter{}

	const performerID = 2

	var images []*models.Image
	matchingPaths, falsePaths := generateTestPaths(performerName, imageExt)
	for i, p := range append(matchingPaths, falsePaths...) {
		images = append(images, &models.Image{
			ID:   i + 1,
			Path: p,
		})
	}

	performer := models.Performer{
		ID:   performerID,
		Name: models.NullString(performerName),
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

	mockImageReader.On("Query", testCtx, image.QueryOptions(expectedImageFilter, expectedFindFilter, false)).
		Return(mocks.ImageQueryResult(images, len(images)), nil).Once()

	for i := range matchingPaths {
		imageID := i + 1
		mockImageReader.On("GetPerformerIDs", testCtx, imageID).Return(nil, nil).Once()
		mockImageReader.On("UpdatePerformers", testCtx, imageID, []int{performerID}).Return(nil).Once()
	}

	err := PerformerImages(testCtx, &performer, nil, mockImageReader, nil)

	assert := assert.New(t)

	assert.Nil(err)
	mockImageReader.AssertExpectations(t)
}

func TestPerformerGalleries(t *testing.T) {
	t.Parallel()

	type test struct {
		performerName string
		expectedRegex string
	}

	performerNames := []test{
		{
			"performer name",
			`(?i)(?:^|_|[^\p{L}\d])performer[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
		},
		{
			"performer + name",
			`(?i)(?:^|_|[^\p{L}\d])performer[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\p{L}\d])`,
		},
	}

	for _, p := range performerNames {
		testPerformerGalleries(t, p.performerName, p.expectedRegex)
	}
}

func testPerformerGalleries(t *testing.T, performerName, expectedRegex string) {
	mockGalleryReader := &mocks.GalleryReaderWriter{}

	const performerID = 2

	var galleries []*models.Gallery
	matchingPaths, falsePaths := generateTestPaths(performerName, galleryExt)
	for i, p := range append(matchingPaths, falsePaths...) {
		galleries = append(galleries, &models.Gallery{
			ID:   i + 1,
			Path: models.NullString(p),
		})
	}

	performer := models.Performer{
		ID:   performerID,
		Name: models.NullString(performerName),
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

	mockGalleryReader.On("Query", testCtx, expectedGalleryFilter, expectedFindFilter).Return(galleries, len(galleries), nil).Once()

	for i := range matchingPaths {
		galleryID := i + 1
		mockGalleryReader.On("GetPerformerIDs", testCtx, galleryID).Return(nil, nil).Once()
		mockGalleryReader.On("UpdatePerformers", testCtx, galleryID, []int{performerID}).Return(nil).Once()
	}

	err := PerformerGalleries(testCtx, &performer, nil, mockGalleryReader, nil)

	assert := assert.New(t)

	assert.Nil(err)
	mockGalleryReader.AssertExpectations(t)
}
