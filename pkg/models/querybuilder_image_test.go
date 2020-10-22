// +build integration

package models_test

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/models"
)

func TestImageFind(t *testing.T) {
	// assume that the first image is imageWithGalleryPath
	sqb := models.NewImageQueryBuilder()

	const imageIdx = 0
	imageID := imageIDs[imageIdx]
	image, err := sqb.Find(imageID)

	if err != nil {
		t.Fatalf("Error finding image: %s", err.Error())
	}

	assert.Equal(t, getImageStringValue(imageIdx, "Path"), image.Path)

	imageID = 0
	image, err = sqb.Find(imageID)

	if err != nil {
		t.Fatalf("Error finding image: %s", err.Error())
	}

	assert.Nil(t, image)
}

func TestImageFindByPath(t *testing.T) {
	sqb := models.NewImageQueryBuilder()

	const imageIdx = 1
	imagePath := getImageStringValue(imageIdx, "Path")
	image, err := sqb.FindByPath(imagePath)

	if err != nil {
		t.Fatalf("Error finding image: %s", err.Error())
	}

	assert.Equal(t, imageIDs[imageIdx], image.ID)
	assert.Equal(t, imagePath, image.Path)

	imagePath = "not exist"
	image, err = sqb.FindByPath(imagePath)

	if err != nil {
		t.Fatalf("Error finding image: %s", err.Error())
	}

	assert.Nil(t, image)
}

func TestImageCountByPerformerID(t *testing.T) {
	sqb := models.NewImageQueryBuilder()
	count, err := sqb.CountByPerformerID(performerIDs[performerIdxWithImage])

	if err != nil {
		t.Fatalf("Error counting images: %s", err.Error())
	}

	assert.Equal(t, 1, count)

	count, err = sqb.CountByPerformerID(0)

	if err != nil {
		t.Fatalf("Error counting images: %s", err.Error())
	}

	assert.Equal(t, 0, count)
}

func TestImageQueryQ(t *testing.T) {
	const imageIdx = 2

	q := getImageStringValue(imageIdx, titleField)

	sqb := models.NewImageQueryBuilder()

	imageQueryQ(t, sqb, q, imageIdx)
}

func imageQueryQ(t *testing.T, sqb models.ImageQueryBuilder, q string, expectedImageIdx int) {
	filter := models.FindFilterType{
		Q: &q,
	}
	images, _ := sqb.Query(nil, &filter)

	assert.Len(t, images, 1)
	image := images[0]
	assert.Equal(t, imageIDs[expectedImageIdx], image.ID)

	// no Q should return all results
	filter.Q = nil
	images, _ = sqb.Query(nil, &filter)

	assert.Len(t, images, totalImages)
}

func TestImageQueryPath(t *testing.T) {
	const imageIdx = 1
	imagePath := getImageStringValue(imageIdx, "Path")

	pathCriterion := models.StringCriterionInput{
		Value:    imagePath,
		Modifier: models.CriterionModifierEquals,
	}

	verifyImagePath(t, pathCriterion)

	pathCriterion.Modifier = models.CriterionModifierNotEquals
	verifyImagePath(t, pathCriterion)
}

func verifyImagePath(t *testing.T, pathCriterion models.StringCriterionInput) {
	sqb := models.NewImageQueryBuilder()
	imageFilter := models.ImageFilterType{
		Path: &pathCriterion,
	}

	images, _ := sqb.Query(&imageFilter, nil)

	for _, image := range images {
		verifyString(t, image.Path, pathCriterion)
	}
}

func TestImageQueryRating(t *testing.T) {
	const rating = 3
	ratingCriterion := models.IntCriterionInput{
		Value:    rating,
		Modifier: models.CriterionModifierEquals,
	}

	verifyImagesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotEquals
	verifyImagesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyImagesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierLessThan
	verifyImagesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierIsNull
	verifyImagesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotNull
	verifyImagesRating(t, ratingCriterion)
}

func verifyImagesRating(t *testing.T, ratingCriterion models.IntCriterionInput) {
	sqb := models.NewImageQueryBuilder()
	imageFilter := models.ImageFilterType{
		Rating: &ratingCriterion,
	}

	images, _ := sqb.Query(&imageFilter, nil)

	for _, image := range images {
		verifyInt64(t, image.Rating, ratingCriterion)
	}
}

func TestImageQueryOCounter(t *testing.T) {
	const oCounter = 1
	oCounterCriterion := models.IntCriterionInput{
		Value:    oCounter,
		Modifier: models.CriterionModifierEquals,
	}

	verifyImagesOCounter(t, oCounterCriterion)

	oCounterCriterion.Modifier = models.CriterionModifierNotEquals
	verifyImagesOCounter(t, oCounterCriterion)

	oCounterCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyImagesOCounter(t, oCounterCriterion)

	oCounterCriterion.Modifier = models.CriterionModifierLessThan
	verifyImagesOCounter(t, oCounterCriterion)
}

func verifyImagesOCounter(t *testing.T, oCounterCriterion models.IntCriterionInput) {
	sqb := models.NewImageQueryBuilder()
	imageFilter := models.ImageFilterType{
		OCounter: &oCounterCriterion,
	}

	images, _ := sqb.Query(&imageFilter, nil)

	for _, image := range images {
		verifyInt(t, image.OCounter, oCounterCriterion)
	}
}

func TestImageQueryResolution(t *testing.T) {
	verifyImagesResolution(t, models.ResolutionEnumLow)
	verifyImagesResolution(t, models.ResolutionEnumStandard)
	verifyImagesResolution(t, models.ResolutionEnumStandardHd)
	verifyImagesResolution(t, models.ResolutionEnumFullHd)
	verifyImagesResolution(t, models.ResolutionEnumFourK)
	verifyImagesResolution(t, models.ResolutionEnum("unknown"))
}

func verifyImagesResolution(t *testing.T, resolution models.ResolutionEnum) {
	sqb := models.NewImageQueryBuilder()
	imageFilter := models.ImageFilterType{
		Resolution: &resolution,
	}

	images, _ := sqb.Query(&imageFilter, nil)

	for _, image := range images {
		verifyImageResolution(t, image.Height, resolution)
	}
}

func verifyImageResolution(t *testing.T, height sql.NullInt64, resolution models.ResolutionEnum) {
	assert := assert.New(t)
	h := height.Int64

	switch resolution {
	case models.ResolutionEnumLow:
		assert.True(h < 480)
	case models.ResolutionEnumStandard:
		assert.True(h >= 480 && h < 720)
	case models.ResolutionEnumStandardHd:
		assert.True(h >= 720 && h < 1080)
	case models.ResolutionEnumFullHd:
		assert.True(h >= 1080 && h < 2160)
	case models.ResolutionEnumFourK:
		assert.True(h >= 2160)
	}
}

func TestImageQueryIsMissingGalleries(t *testing.T) {
	sqb := models.NewImageQueryBuilder()
	isMissing := "galleries"
	imageFilter := models.ImageFilterType{
		IsMissing: &isMissing,
	}

	q := getImageStringValue(imageIdxWithGallery, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	images, _ := sqb.Query(&imageFilter, &findFilter)

	assert.Len(t, images, 0)

	findFilter.Q = nil
	images, _ = sqb.Query(&imageFilter, &findFilter)

	// ensure non of the ids equal the one with gallery
	for _, image := range images {
		assert.NotEqual(t, imageIDs[imageIdxWithGallery], image.ID)
	}
}

func TestImageQueryIsMissingStudio(t *testing.T) {
	sqb := models.NewImageQueryBuilder()
	isMissing := "studio"
	imageFilter := models.ImageFilterType{
		IsMissing: &isMissing,
	}

	q := getImageStringValue(imageIdxWithStudio, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	images, _ := sqb.Query(&imageFilter, &findFilter)

	assert.Len(t, images, 0)

	findFilter.Q = nil
	images, _ = sqb.Query(&imageFilter, &findFilter)

	// ensure non of the ids equal the one with studio
	for _, image := range images {
		assert.NotEqual(t, imageIDs[imageIdxWithStudio], image.ID)
	}
}

func TestImageQueryIsMissingPerformers(t *testing.T) {
	sqb := models.NewImageQueryBuilder()
	isMissing := "performers"
	imageFilter := models.ImageFilterType{
		IsMissing: &isMissing,
	}

	q := getImageStringValue(imageIdxWithPerformer, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	images, _ := sqb.Query(&imageFilter, &findFilter)

	assert.Len(t, images, 0)

	findFilter.Q = nil
	images, _ = sqb.Query(&imageFilter, &findFilter)

	assert.True(t, len(images) > 0)

	// ensure non of the ids equal the one with movies
	for _, image := range images {
		assert.NotEqual(t, imageIDs[imageIdxWithPerformer], image.ID)
	}
}

func TestImageQueryIsMissingTags(t *testing.T) {
	sqb := models.NewImageQueryBuilder()
	isMissing := "tags"
	imageFilter := models.ImageFilterType{
		IsMissing: &isMissing,
	}

	q := getImageStringValue(imageIdxWithTwoTags, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	images, _ := sqb.Query(&imageFilter, &findFilter)

	assert.Len(t, images, 0)

	findFilter.Q = nil
	images, _ = sqb.Query(&imageFilter, &findFilter)

	assert.True(t, len(images) > 0)
}

func TestImageQueryIsMissingRating(t *testing.T) {
	sqb := models.NewImageQueryBuilder()
	isMissing := "rating"
	imageFilter := models.ImageFilterType{
		IsMissing: &isMissing,
	}

	images, _ := sqb.Query(&imageFilter, nil)

	assert.True(t, len(images) > 0)

	// ensure date is null, empty or "0001-01-01"
	for _, image := range images {
		assert.True(t, !image.Rating.Valid)
	}
}

func TestImageQueryPerformers(t *testing.T) {
	sqb := models.NewImageQueryBuilder()
	performerCriterion := models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(performerIDs[performerIdxWithImage]),
			strconv.Itoa(performerIDs[performerIdx1WithImage]),
		},
		Modifier: models.CriterionModifierIncludes,
	}

	imageFilter := models.ImageFilterType{
		Performers: &performerCriterion,
	}

	images, _ := sqb.Query(&imageFilter, nil)

	assert.Len(t, images, 2)

	// ensure ids are correct
	for _, image := range images {
		assert.True(t, image.ID == imageIDs[imageIdxWithPerformer] || image.ID == imageIDs[imageIdxWithTwoPerformers])
	}

	performerCriterion = models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(performerIDs[performerIdx1WithImage]),
			strconv.Itoa(performerIDs[performerIdx2WithImage]),
		},
		Modifier: models.CriterionModifierIncludesAll,
	}

	images, _ = sqb.Query(&imageFilter, nil)

	assert.Len(t, images, 1)
	assert.Equal(t, imageIDs[imageIdxWithTwoPerformers], images[0].ID)

	performerCriterion = models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(performerIDs[performerIdx1WithImage]),
		},
		Modifier: models.CriterionModifierExcludes,
	}

	q := getImageStringValue(imageIdxWithTwoPerformers, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	images, _ = sqb.Query(&imageFilter, &findFilter)
	assert.Len(t, images, 0)
}

func TestImageQueryTags(t *testing.T) {
	sqb := models.NewImageQueryBuilder()
	tagCriterion := models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(tagIDs[tagIdxWithImage]),
			strconv.Itoa(tagIDs[tagIdx1WithImage]),
		},
		Modifier: models.CriterionModifierIncludes,
	}

	imageFilter := models.ImageFilterType{
		Tags: &tagCriterion,
	}

	images, _ := sqb.Query(&imageFilter, nil)

	assert.Len(t, images, 2)

	// ensure ids are correct
	for _, image := range images {
		assert.True(t, image.ID == imageIDs[imageIdxWithTag] || image.ID == imageIDs[imageIdxWithTwoTags])
	}

	tagCriterion = models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(tagIDs[tagIdx1WithImage]),
			strconv.Itoa(tagIDs[tagIdx2WithImage]),
		},
		Modifier: models.CriterionModifierIncludesAll,
	}

	images, _ = sqb.Query(&imageFilter, nil)

	assert.Len(t, images, 1)
	assert.Equal(t, imageIDs[imageIdxWithTwoTags], images[0].ID)

	tagCriterion = models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(tagIDs[tagIdx1WithImage]),
		},
		Modifier: models.CriterionModifierExcludes,
	}

	q := getImageStringValue(imageIdxWithTwoTags, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	images, _ = sqb.Query(&imageFilter, &findFilter)
	assert.Len(t, images, 0)
}

func TestImageQueryStudio(t *testing.T) {
	sqb := models.NewImageQueryBuilder()
	studioCriterion := models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(studioIDs[studioIdxWithImage]),
		},
		Modifier: models.CriterionModifierIncludes,
	}

	imageFilter := models.ImageFilterType{
		Studios: &studioCriterion,
	}

	images, _ := sqb.Query(&imageFilter, nil)

	assert.Len(t, images, 1)

	// ensure id is correct
	assert.Equal(t, imageIDs[imageIdxWithStudio], images[0].ID)

	studioCriterion = models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(studioIDs[studioIdxWithImage]),
		},
		Modifier: models.CriterionModifierExcludes,
	}

	q := getImageStringValue(imageIdxWithStudio, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	images, _ = sqb.Query(&imageFilter, &findFilter)
	assert.Len(t, images, 0)
}

func TestImageQuerySorting(t *testing.T) {
	sort := titleField
	direction := models.SortDirectionEnumAsc
	findFilter := models.FindFilterType{
		Sort:      &sort,
		Direction: &direction,
	}

	sqb := models.NewImageQueryBuilder()
	images, _ := sqb.Query(nil, &findFilter)

	// images should be in same order as indexes
	firstImage := images[0]
	lastImage := images[len(images)-1]

	assert.Equal(t, imageIDs[0], firstImage.ID)
	assert.Equal(t, imageIDs[len(imageIDs)-1], lastImage.ID)

	// sort in descending order
	direction = models.SortDirectionEnumDesc

	images, _ = sqb.Query(nil, &findFilter)
	firstImage = images[0]
	lastImage = images[len(images)-1]

	assert.Equal(t, imageIDs[len(imageIDs)-1], firstImage.ID)
	assert.Equal(t, imageIDs[0], lastImage.ID)
}

func TestImageQueryPagination(t *testing.T) {
	perPage := 1
	findFilter := models.FindFilterType{
		PerPage: &perPage,
	}

	sqb := models.NewImageQueryBuilder()
	images, _ := sqb.Query(nil, &findFilter)

	assert.Len(t, images, 1)

	firstID := images[0].ID

	page := 2
	findFilter.Page = &page
	images, _ = sqb.Query(nil, &findFilter)

	assert.Len(t, images, 1)
	secondID := images[0].ID
	assert.NotEqual(t, firstID, secondID)

	perPage = 2
	page = 1

	images, _ = sqb.Query(nil, &findFilter)
	assert.Len(t, images, 2)
	assert.Equal(t, firstID, images[0].ID)
	assert.Equal(t, secondID, images[1].ID)
}

func TestImageCountByTagID(t *testing.T) {
	sqb := models.NewImageQueryBuilder()

	imageCount, err := sqb.CountByTagID(tagIDs[tagIdxWithImage])

	if err != nil {
		t.Fatalf("error calling CountByTagID: %s", err.Error())
	}

	assert.Equal(t, 1, imageCount)

	imageCount, err = sqb.CountByTagID(0)

	if err != nil {
		t.Fatalf("error calling CountByTagID: %s", err.Error())
	}

	assert.Equal(t, 0, imageCount)
}

func TestImageCountByStudioID(t *testing.T) {
	sqb := models.NewImageQueryBuilder()

	imageCount, err := sqb.CountByStudioID(studioIDs[studioIdxWithImage])

	if err != nil {
		t.Fatalf("error calling CountByStudioID: %s", err.Error())
	}

	assert.Equal(t, 1, imageCount)

	imageCount, err = sqb.CountByStudioID(0)

	if err != nil {
		t.Fatalf("error calling CountByStudioID: %s", err.Error())
	}

	assert.Equal(t, 0, imageCount)
}

func TestImageFindByPerformerID(t *testing.T) {
	sqb := models.NewImageQueryBuilder()

	images, err := sqb.FindByPerformerID(performerIDs[performerIdxWithImage])

	if err != nil {
		t.Fatalf("error calling FindByPerformerID: %s", err.Error())
	}

	assert.Len(t, images, 1)
	assert.Equal(t, imageIDs[imageIdxWithPerformer], images[0].ID)

	images, err = sqb.FindByPerformerID(0)

	if err != nil {
		t.Fatalf("error calling FindByPerformerID: %s", err.Error())
	}

	assert.Len(t, images, 0)
}

func TestImageFindByStudioID(t *testing.T) {
	sqb := models.NewImageQueryBuilder()

	images, err := sqb.FindByStudioID(performerIDs[studioIdxWithImage])

	if err != nil {
		t.Fatalf("error calling FindByStudioID: %s", err.Error())
	}

	assert.Len(t, images, 1)
	assert.Equal(t, imageIDs[imageIdxWithStudio], images[0].ID)

	images, err = sqb.FindByStudioID(0)

	if err != nil {
		t.Fatalf("error calling FindByStudioID: %s", err.Error())
	}

	assert.Len(t, images, 0)
}

// TODO Update
// TODO IncrementOCounter
// TODO DecrementOCounter
// TODO ResetOCounter
// TODO Destroy
// TODO FindByChecksum
// TODO Count
// TODO SizeCount
// TODO All
