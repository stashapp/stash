package autotag

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
)

func TestStudioScenes(t *testing.T) {
	type test struct {
		studioName    string
		expectedRegex string
	}

	studioNames := []test{
		{
			"studio name",
			`(?i)(?:^|_|[^\w\d])studio[.\-_ ]*name(?:$|_|[^\w\d])`,
		},
		{
			"studio + name",
			`(?i)(?:^|_|[^\w\d])studio[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\w\d])`,
		},
	}

	for _, p := range studioNames {
		testStudioScenes(t, p.studioName, p.expectedRegex)
	}
}

func testStudioScenes(t *testing.T, studioName, expectedRegex string) {
	mockSceneReader := &mocks.SceneReaderWriter{}

	const studioID = 2

	var scenes []*models.Scene
	matchingPaths, falsePaths := generateTestPaths(studioName, sceneExt)
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

	mockSceneReader.On("Query", expectedSceneFilter, expectedFindFilter).Return(scenes, len(scenes), nil).Once()

	for i := range matchingPaths {
		sceneID := i + 1
		mockSceneReader.On("Find", sceneID).Return(&models.Scene{}, nil).Once()
		expectedStudioID := models.NullInt64(studioID)
		mockSceneReader.On("Update", models.ScenePartial{
			ID:       sceneID,
			StudioID: &expectedStudioID,
		}).Return(nil, nil).Once()
	}

	err := StudioScenes(&studio, nil, mockSceneReader)

	assert := assert.New(t)

	assert.Nil(err)
	mockSceneReader.AssertExpectations(t)
}

func TestStudioImages(t *testing.T) {
	type test struct {
		studioName    string
		expectedRegex string
	}

	studioNames := []test{
		{
			"studio name",
			`(?i)(?:^|_|[^\w\d])studio[.\-_ ]*name(?:$|_|[^\w\d])`,
		},
		{
			"studio + name",
			`(?i)(?:^|_|[^\w\d])studio[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\w\d])`,
		},
	}

	for _, p := range studioNames {
		testStudioImages(t, p.studioName, p.expectedRegex)
	}
}

func testStudioImages(t *testing.T, studioName, expectedRegex string) {
	mockImageReader := &mocks.ImageReaderWriter{}

	const studioID = 2

	var images []*models.Image
	matchingPaths, falsePaths := generateTestPaths(studioName, imageExt)
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

	mockImageReader.On("Query", expectedImageFilter, expectedFindFilter).Return(images, len(images), nil).Once()

	for i := range matchingPaths {
		imageID := i + 1
		mockImageReader.On("Find", imageID).Return(&models.Image{}, nil).Once()
		expectedStudioID := models.NullInt64(studioID)
		mockImageReader.On("Update", models.ImagePartial{
			ID:       imageID,
			StudioID: &expectedStudioID,
		}).Return(nil, nil).Once()
	}

	err := StudioImages(&studio, nil, mockImageReader)

	assert := assert.New(t)

	assert.Nil(err)
	mockImageReader.AssertExpectations(t)
}

func TestStudioGalleries(t *testing.T) {
	type test struct {
		studioName    string
		expectedRegex string
	}

	studioNames := []test{
		{
			"studio name",
			`(?i)(?:^|_|[^\w\d])studio[.\-_ ]*name(?:$|_|[^\w\d])`,
		},
		{
			"studio + name",
			`(?i)(?:^|_|[^\w\d])studio[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\w\d])`,
		},
	}

	for _, p := range studioNames {
		testStudioGalleries(t, p.studioName, p.expectedRegex)
	}
}

func testStudioGalleries(t *testing.T, studioName, expectedRegex string) {
	mockGalleryReader := &mocks.GalleryReaderWriter{}

	const studioID = 2

	var galleries []*models.Gallery
	matchingPaths, falsePaths := generateTestPaths(studioName, galleryExt)
	for i, p := range append(matchingPaths, falsePaths...) {
		galleries = append(galleries, &models.Gallery{
			ID:   i + 1,
			Path: models.NullString(p),
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

	mockGalleryReader.On("Query", expectedGalleryFilter, expectedFindFilter).Return(galleries, len(galleries), nil).Once()

	for i := range matchingPaths {
		galleryID := i + 1
		mockGalleryReader.On("Find", galleryID).Return(&models.Gallery{}, nil).Once()
		expectedStudioID := models.NullInt64(studioID)
		mockGalleryReader.On("UpdatePartial", models.GalleryPartial{
			ID:       galleryID,
			StudioID: &expectedStudioID,
		}).Return(nil, nil).Once()
	}

	err := StudioGalleries(&studio, nil, mockGalleryReader)

	assert := assert.New(t)

	assert.Nil(err)
	mockGalleryReader.AssertExpectations(t)
}
