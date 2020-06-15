// +build integration

package models_test

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/models"
)

func TestSceneFind(t *testing.T) {
	// assume that the first scene is sceneWithGalleryPath
	sqb := models.NewSceneQueryBuilder()

	const sceneIdx = 0
	sceneID := sceneIDs[sceneIdx]
	scene, err := sqb.Find(sceneID)

	if err != nil {
		t.Fatalf("Error finding scene: %s", err.Error())
	}

	assert.Equal(t, getSceneStringValue(sceneIdx, "Path"), scene.Path)

	sceneID = 0
	scene, err = sqb.Find(sceneID)

	if err != nil {
		t.Fatalf("Error finding scene: %s", err.Error())
	}

	assert.Nil(t, scene)
}

func TestSceneFindByPath(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()

	const sceneIdx = 1
	scenePath := getSceneStringValue(sceneIdx, "Path")
	scene, err := sqb.FindByPath(scenePath)

	if err != nil {
		t.Fatalf("Error finding scene: %s", err.Error())
	}

	assert.Equal(t, sceneIDs[sceneIdx], scene.ID)
	assert.Equal(t, scenePath, scene.Path)

	scenePath = "not exist"
	scene, err = sqb.FindByPath(scenePath)

	if err != nil {
		t.Fatalf("Error finding scene: %s", err.Error())
	}

	assert.Nil(t, scene)
}

func TestSceneCountByPerformerID(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()
	count, err := sqb.CountByPerformerID(performerIDs[performerIdxWithScene])

	if err != nil {
		t.Fatalf("Error counting scenes: %s", err.Error())
	}

	assert.Equal(t, 1, count)

	count, err = sqb.CountByPerformerID(0)

	if err != nil {
		t.Fatalf("Error counting scenes: %s", err.Error())
	}

	assert.Equal(t, 0, count)
}

func TestSceneWall(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()

	const sceneIdx = 2
	wallQuery := getSceneStringValue(sceneIdx, "Details")
	scenes, err := sqb.Wall(&wallQuery)

	if err != nil {
		t.Fatalf("Error finding scenes: %s", err.Error())
	}

	assert.Len(t, scenes, 1)
	scene := scenes[0]
	assert.Equal(t, sceneIDs[sceneIdx], scene.ID)
	assert.Equal(t, getSceneStringValue(sceneIdx, "Path"), scene.Path)

	wallQuery = "not exist"
	scenes, err = sqb.Wall(&wallQuery)

	if err != nil {
		t.Fatalf("Error finding scene: %s", err.Error())
	}

	assert.Len(t, scenes, 0)
}

func TestSceneQueryQ(t *testing.T) {
	const sceneIdx = 2

	q := getSceneStringValue(sceneIdx, titleField)

	sqb := models.NewSceneQueryBuilder()

	sceneQueryQ(t, sqb, q, sceneIdx)
}

func sceneQueryQ(t *testing.T, sqb models.SceneQueryBuilder, q string, expectedSceneIdx int) {
	filter := models.FindFilterType{
		Q: &q,
	}
	scenes, _ := sqb.Query(nil, &filter)

	assert.Len(t, scenes, 1)
	scene := scenes[0]
	assert.Equal(t, sceneIDs[expectedSceneIdx], scene.ID)

	// no Q should return all results
	filter.Q = nil
	scenes, _ = sqb.Query(nil, &filter)

	assert.Len(t, scenes, totalScenes)
}

func TestSceneQueryRating(t *testing.T) {
	const rating = 3
	ratingCriterion := models.IntCriterionInput{
		Value:    rating,
		Modifier: models.CriterionModifierEquals,
	}

	verifyScenesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotEquals
	verifyScenesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyScenesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierLessThan
	verifyScenesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierIsNull
	verifyScenesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotNull
	verifyScenesRating(t, ratingCriterion)
}

func verifyScenesRating(t *testing.T, ratingCriterion models.IntCriterionInput) {
	sqb := models.NewSceneQueryBuilder()
	sceneFilter := models.SceneFilterType{
		Rating: &ratingCriterion,
	}

	scenes, _ := sqb.Query(&sceneFilter, nil)

	for _, scene := range scenes {
		verifyInt64(t, scene.Rating, ratingCriterion)
	}
}

func verifyInt64(t *testing.T, value sql.NullInt64, criterion models.IntCriterionInput) {
	assert := assert.New(t)
	if criterion.Modifier == models.CriterionModifierIsNull {
		assert.False(value.Valid, "expect is null values to be null")
	}
	if criterion.Modifier == models.CriterionModifierNotNull {
		assert.True(value.Valid, "expect is null values to be null")
	}
	if criterion.Modifier == models.CriterionModifierEquals {
		assert.Equal(int64(criterion.Value), value.Int64)
	}
	if criterion.Modifier == models.CriterionModifierNotEquals {
		assert.NotEqual(int64(criterion.Value), value.Int64)
	}
	if criterion.Modifier == models.CriterionModifierGreaterThan {
		assert.True(value.Int64 > int64(criterion.Value))
	}
	if criterion.Modifier == models.CriterionModifierLessThan {
		assert.True(value.Int64 < int64(criterion.Value))
	}
}

func TestSceneQueryOCounter(t *testing.T) {
	const oCounter = 1
	oCounterCriterion := models.IntCriterionInput{
		Value:    oCounter,
		Modifier: models.CriterionModifierEquals,
	}

	verifyScenesOCounter(t, oCounterCriterion)

	oCounterCriterion.Modifier = models.CriterionModifierNotEquals
	verifyScenesOCounter(t, oCounterCriterion)

	oCounterCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyScenesOCounter(t, oCounterCriterion)

	oCounterCriterion.Modifier = models.CriterionModifierLessThan
	verifyScenesOCounter(t, oCounterCriterion)
}

func verifyScenesOCounter(t *testing.T, oCounterCriterion models.IntCriterionInput) {
	sqb := models.NewSceneQueryBuilder()
	sceneFilter := models.SceneFilterType{
		OCounter: &oCounterCriterion,
	}

	scenes, _ := sqb.Query(&sceneFilter, nil)

	for _, scene := range scenes {
		verifyInt(t, scene.OCounter, oCounterCriterion)
	}
}

func verifyInt(t *testing.T, value int, criterion models.IntCriterionInput) {
	assert := assert.New(t)
	if criterion.Modifier == models.CriterionModifierEquals {
		assert.Equal(criterion.Value, value)
	}
	if criterion.Modifier == models.CriterionModifierNotEquals {
		assert.NotEqual(criterion.Value, value)
	}
	if criterion.Modifier == models.CriterionModifierGreaterThan {
		assert.True(value > criterion.Value)
	}
	if criterion.Modifier == models.CriterionModifierLessThan {
		assert.True(value < criterion.Value)
	}
}

func TestSceneQueryDuration(t *testing.T) {
	duration := 200.432

	durationCriterion := models.IntCriterionInput{
		Value:    int(duration),
		Modifier: models.CriterionModifierEquals,
	}
	verifyScenesDuration(t, durationCriterion)

	durationCriterion.Modifier = models.CriterionModifierNotEquals
	verifyScenesDuration(t, durationCriterion)

	durationCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyScenesDuration(t, durationCriterion)

	durationCriterion.Modifier = models.CriterionModifierLessThan
	verifyScenesDuration(t, durationCriterion)

	durationCriterion.Modifier = models.CriterionModifierIsNull
	verifyScenesDuration(t, durationCriterion)

	durationCriterion.Modifier = models.CriterionModifierNotNull
	verifyScenesDuration(t, durationCriterion)
}

func verifyScenesDuration(t *testing.T, durationCriterion models.IntCriterionInput) {
	sqb := models.NewSceneQueryBuilder()
	sceneFilter := models.SceneFilterType{
		Duration: &durationCriterion,
	}

	scenes, _ := sqb.Query(&sceneFilter, nil)

	for _, scene := range scenes {
		if durationCriterion.Modifier == models.CriterionModifierEquals {
			assert.True(t, scene.Duration.Float64 >= float64(durationCriterion.Value) && scene.Duration.Float64 < float64(durationCriterion.Value+1))
		} else if durationCriterion.Modifier == models.CriterionModifierNotEquals {
			assert.True(t, scene.Duration.Float64 < float64(durationCriterion.Value) || scene.Duration.Float64 >= float64(durationCriterion.Value+1))
		} else {
			verifyFloat64(t, scene.Duration, durationCriterion)
		}
	}
}

func verifyFloat64(t *testing.T, value sql.NullFloat64, criterion models.IntCriterionInput) {
	assert := assert.New(t)
	if criterion.Modifier == models.CriterionModifierIsNull {
		assert.False(value.Valid, "expect is null values to be null")
	}
	if criterion.Modifier == models.CriterionModifierNotNull {
		assert.True(value.Valid, "expect is null values to be null")
	}
	if criterion.Modifier == models.CriterionModifierEquals {
		assert.Equal(float64(criterion.Value), value.Float64)
	}
	if criterion.Modifier == models.CriterionModifierNotEquals {
		assert.NotEqual(float64(criterion.Value), value.Float64)
	}
	if criterion.Modifier == models.CriterionModifierGreaterThan {
		assert.True(value.Float64 > float64(criterion.Value))
	}
	if criterion.Modifier == models.CriterionModifierLessThan {
		assert.True(value.Float64 < float64(criterion.Value))
	}
}

func TestSceneQueryResolution(t *testing.T) {
	verifyScenesResolution(t, models.ResolutionEnumLow)
	verifyScenesResolution(t, models.ResolutionEnumStandard)
	verifyScenesResolution(t, models.ResolutionEnumStandardHd)
	verifyScenesResolution(t, models.ResolutionEnumFullHd)
	verifyScenesResolution(t, models.ResolutionEnumFourK)
	verifyScenesResolution(t, models.ResolutionEnum("unknown"))
}

func verifyScenesResolution(t *testing.T, resolution models.ResolutionEnum) {
	sqb := models.NewSceneQueryBuilder()
	sceneFilter := models.SceneFilterType{
		Resolution: &resolution,
	}

	scenes, _ := sqb.Query(&sceneFilter, nil)

	for _, scene := range scenes {
		verifySceneResolution(t, scene.Height, resolution)
	}
}

func verifySceneResolution(t *testing.T, height sql.NullInt64, resolution models.ResolutionEnum) {
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

func TestSceneQueryHasMarkers(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()
	hasMarkers := "true"
	sceneFilter := models.SceneFilterType{
		HasMarkers: &hasMarkers,
	}

	q := getSceneStringValue(sceneIdxWithMarker, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	scenes, _ := sqb.Query(&sceneFilter, &findFilter)

	assert.Len(t, scenes, 1)
	assert.Equal(t, sceneIDs[sceneIdxWithMarker], scenes[0].ID)

	hasMarkers = "false"
	scenes, _ = sqb.Query(&sceneFilter, &findFilter)
	assert.Len(t, scenes, 0)

	findFilter.Q = nil
	scenes, _ = sqb.Query(&sceneFilter, &findFilter)

	assert.NotEqual(t, 0, len(scenes))

	// ensure non of the ids equal the one with gallery
	for _, scene := range scenes {
		assert.NotEqual(t, sceneIDs[sceneIdxWithMarker], scene.ID)
	}
}

func TestSceneQueryIsMissingGallery(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()
	isMissing := "gallery"
	sceneFilter := models.SceneFilterType{
		IsMissing: &isMissing,
	}

	q := getSceneStringValue(sceneIdxWithGallery, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	scenes, _ := sqb.Query(&sceneFilter, &findFilter)

	assert.Len(t, scenes, 0)

	findFilter.Q = nil
	scenes, _ = sqb.Query(&sceneFilter, &findFilter)

	// ensure non of the ids equal the one with gallery
	for _, scene := range scenes {
		assert.NotEqual(t, sceneIDs[sceneIdxWithGallery], scene.ID)
	}
}

func TestSceneQueryIsMissingStudio(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()
	isMissing := "studio"
	sceneFilter := models.SceneFilterType{
		IsMissing: &isMissing,
	}

	q := getSceneStringValue(sceneIdxWithStudio, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	scenes, _ := sqb.Query(&sceneFilter, &findFilter)

	assert.Len(t, scenes, 0)

	findFilter.Q = nil
	scenes, _ = sqb.Query(&sceneFilter, &findFilter)

	// ensure non of the ids equal the one with studio
	for _, scene := range scenes {
		assert.NotEqual(t, sceneIDs[sceneIdxWithStudio], scene.ID)
	}
}

func TestSceneQueryIsMissingMovies(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()
	isMissing := "movie"
	sceneFilter := models.SceneFilterType{
		IsMissing: &isMissing,
	}

	q := getSceneStringValue(sceneIdxWithMovie, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	scenes, _ := sqb.Query(&sceneFilter, &findFilter)

	assert.Len(t, scenes, 0)

	findFilter.Q = nil
	scenes, _ = sqb.Query(&sceneFilter, &findFilter)

	// ensure non of the ids equal the one with movies
	for _, scene := range scenes {
		assert.NotEqual(t, sceneIDs[sceneIdxWithMovie], scene.ID)
	}
}

func TestSceneQueryIsMissingPerformers(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()
	isMissing := "performers"
	sceneFilter := models.SceneFilterType{
		IsMissing: &isMissing,
	}

	q := getSceneStringValue(sceneIdxWithPerformer, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	scenes, _ := sqb.Query(&sceneFilter, &findFilter)

	assert.Len(t, scenes, 0)

	findFilter.Q = nil
	scenes, _ = sqb.Query(&sceneFilter, &findFilter)

	assert.True(t, len(scenes) > 0)

	// ensure non of the ids equal the one with movies
	for _, scene := range scenes {
		assert.NotEqual(t, sceneIDs[sceneIdxWithPerformer], scene.ID)
	}
}

func TestSceneQueryIsMissingDate(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()
	isMissing := "date"
	sceneFilter := models.SceneFilterType{
		IsMissing: &isMissing,
	}

	scenes, _ := sqb.Query(&sceneFilter, nil)

	assert.True(t, len(scenes) > 0)

	// ensure date is null, empty or "0001-01-01"
	for _, scene := range scenes {
		assert.True(t, !scene.Date.Valid || scene.Date.String == "" || scene.Date.String == "0001-01-01")
	}
}

func TestSceneQueryIsMissingTags(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()
	isMissing := "tags"
	sceneFilter := models.SceneFilterType{
		IsMissing: &isMissing,
	}

	q := getSceneStringValue(sceneIdxWithTwoTags, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	scenes, _ := sqb.Query(&sceneFilter, &findFilter)

	assert.Len(t, scenes, 0)

	findFilter.Q = nil
	scenes, _ = sqb.Query(&sceneFilter, &findFilter)

	assert.True(t, len(scenes) > 0)
}

func TestSceneQueryIsMissingRating(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()
	isMissing := "rating"
	sceneFilter := models.SceneFilterType{
		IsMissing: &isMissing,
	}

	scenes, _ := sqb.Query(&sceneFilter, nil)

	assert.True(t, len(scenes) > 0)

	// ensure date is null, empty or "0001-01-01"
	for _, scene := range scenes {
		assert.True(t, !scene.Rating.Valid)
	}
}

func TestSceneQueryPerformers(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()
	performerCriterion := models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(performerIDs[performerIdxWithScene]),
			strconv.Itoa(performerIDs[performerIdx1WithScene]),
		},
		Modifier: models.CriterionModifierIncludes,
	}

	sceneFilter := models.SceneFilterType{
		Performers: &performerCriterion,
	}

	scenes, _ := sqb.Query(&sceneFilter, nil)

	assert.Len(t, scenes, 2)

	// ensure ids are correct
	for _, scene := range scenes {
		assert.True(t, scene.ID == sceneIDs[sceneIdxWithPerformer] || scene.ID == sceneIDs[sceneIdxWithTwoPerformers])
	}

	performerCriterion = models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(performerIDs[performerIdx1WithScene]),
			strconv.Itoa(performerIDs[performerIdx2WithScene]),
		},
		Modifier: models.CriterionModifierIncludesAll,
	}

	scenes, _ = sqb.Query(&sceneFilter, nil)

	assert.Len(t, scenes, 1)
	assert.Equal(t, sceneIDs[sceneIdxWithTwoPerformers], scenes[0].ID)

	performerCriterion = models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(performerIDs[performerIdx1WithScene]),
		},
		Modifier: models.CriterionModifierExcludes,
	}

	q := getSceneStringValue(sceneIdxWithTwoPerformers, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	scenes, _ = sqb.Query(&sceneFilter, &findFilter)
	assert.Len(t, scenes, 0)
}

func TestSceneQueryTags(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()
	tagCriterion := models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(tagIDs[tagIdxWithScene]),
			strconv.Itoa(tagIDs[tagIdx1WithScene]),
		},
		Modifier: models.CriterionModifierIncludes,
	}

	sceneFilter := models.SceneFilterType{
		Tags: &tagCriterion,
	}

	scenes, _ := sqb.Query(&sceneFilter, nil)

	assert.Len(t, scenes, 2)

	// ensure ids are correct
	for _, scene := range scenes {
		assert.True(t, scene.ID == sceneIDs[sceneIdxWithTag] || scene.ID == sceneIDs[sceneIdxWithTwoTags])
	}

	tagCriterion = models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(tagIDs[tagIdx1WithScene]),
			strconv.Itoa(tagIDs[tagIdx2WithScene]),
		},
		Modifier: models.CriterionModifierIncludesAll,
	}

	scenes, _ = sqb.Query(&sceneFilter, nil)

	assert.Len(t, scenes, 1)
	assert.Equal(t, sceneIDs[sceneIdxWithTwoTags], scenes[0].ID)

	tagCriterion = models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(tagIDs[tagIdx1WithScene]),
		},
		Modifier: models.CriterionModifierExcludes,
	}

	q := getSceneStringValue(sceneIdxWithTwoTags, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	scenes, _ = sqb.Query(&sceneFilter, &findFilter)
	assert.Len(t, scenes, 0)
}

func TestSceneQueryStudio(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()
	studioCriterion := models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(studioIDs[studioIdxWithScene]),
		},
		Modifier: models.CriterionModifierIncludes,
	}

	sceneFilter := models.SceneFilterType{
		Studios: &studioCriterion,
	}

	scenes, _ := sqb.Query(&sceneFilter, nil)

	assert.Len(t, scenes, 1)

	// ensure id is correct
	assert.Equal(t, sceneIDs[sceneIdxWithStudio], scenes[0].ID)

	studioCriterion = models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(studioIDs[studioIdxWithScene]),
		},
		Modifier: models.CriterionModifierExcludes,
	}

	q := getSceneStringValue(sceneIdxWithStudio, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	scenes, _ = sqb.Query(&sceneFilter, &findFilter)
	assert.Len(t, scenes, 0)
}

func TestSceneQueryMovies(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()
	movieCriterion := models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(movieIDs[movieIdxWithScene]),
		},
		Modifier: models.CriterionModifierIncludes,
	}

	sceneFilter := models.SceneFilterType{
		Movies: &movieCriterion,
	}

	scenes, _ := sqb.Query(&sceneFilter, nil)

	assert.Len(t, scenes, 1)

	// ensure id is correct
	assert.Equal(t, sceneIDs[sceneIdxWithMovie], scenes[0].ID)

	movieCriterion = models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(movieIDs[movieIdxWithScene]),
		},
		Modifier: models.CriterionModifierExcludes,
	}

	q := getSceneStringValue(sceneIdxWithMovie, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	scenes, _ = sqb.Query(&sceneFilter, &findFilter)
	assert.Len(t, scenes, 0)
}

func TestSceneQuerySorting(t *testing.T) {
	sort := titleField
	direction := models.SortDirectionEnumAsc
	findFilter := models.FindFilterType{
		Sort:      &sort,
		Direction: &direction,
	}

	sqb := models.NewSceneQueryBuilder()
	scenes, _ := sqb.Query(nil, &findFilter)

	// scenes should be in same order as indexes
	firstScene := scenes[0]
	lastScene := scenes[len(scenes)-1]

	assert.Equal(t, sceneIDs[0], firstScene.ID)
	assert.Equal(t, sceneIDs[len(sceneIDs)-1], lastScene.ID)

	// sort in descending order
	direction = models.SortDirectionEnumDesc

	scenes, _ = sqb.Query(nil, &findFilter)
	firstScene = scenes[0]
	lastScene = scenes[len(scenes)-1]

	assert.Equal(t, sceneIDs[len(sceneIDs)-1], firstScene.ID)
	assert.Equal(t, sceneIDs[0], lastScene.ID)
}

func TestSceneQueryPagination(t *testing.T) {
	perPage := 1
	findFilter := models.FindFilterType{
		PerPage: &perPage,
	}

	sqb := models.NewSceneQueryBuilder()
	scenes, _ := sqb.Query(nil, &findFilter)

	assert.Len(t, scenes, 1)

	firstID := scenes[0].ID

	page := 2
	findFilter.Page = &page
	scenes, _ = sqb.Query(nil, &findFilter)

	assert.Len(t, scenes, 1)
	secondID := scenes[0].ID
	assert.NotEqual(t, firstID, secondID)

	perPage = 2
	page = 1

	scenes, _ = sqb.Query(nil, &findFilter)
	assert.Len(t, scenes, 2)
	assert.Equal(t, firstID, scenes[0].ID)
	assert.Equal(t, secondID, scenes[1].ID)
}

func TestSceneCountByTagID(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()

	sceneCount, err := sqb.CountByTagID(tagIDs[tagIdxWithScene])

	if err != nil {
		t.Fatalf("error calling CountByTagID: %s", err.Error())
	}

	assert.Equal(t, 1, sceneCount)

	sceneCount, err = sqb.CountByTagID(0)

	if err != nil {
		t.Fatalf("error calling CountByTagID: %s", err.Error())
	}

	assert.Equal(t, 0, sceneCount)
}

func TestSceneCountByMovieID(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()

	sceneCount, err := sqb.CountByMovieID(movieIDs[movieIdxWithScene])

	if err != nil {
		t.Fatalf("error calling CountByMovieID: %s", err.Error())
	}

	assert.Equal(t, 1, sceneCount)

	sceneCount, err = sqb.CountByMovieID(0)

	if err != nil {
		t.Fatalf("error calling CountByMovieID: %s", err.Error())
	}

	assert.Equal(t, 0, sceneCount)
}

func TestSceneCountByStudioID(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()

	sceneCount, err := sqb.CountByStudioID(studioIDs[studioIdxWithScene])

	if err != nil {
		t.Fatalf("error calling CountByStudioID: %s", err.Error())
	}

	assert.Equal(t, 1, sceneCount)

	sceneCount, err = sqb.CountByStudioID(0)

	if err != nil {
		t.Fatalf("error calling CountByStudioID: %s", err.Error())
	}

	assert.Equal(t, 0, sceneCount)
}

func TestFindByMovieID(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()

	scenes, err := sqb.FindByMovieID(movieIDs[movieIdxWithScene])

	if err != nil {
		t.Fatalf("error calling FindByMovieID: %s", err.Error())
	}

	assert.Len(t, scenes, 1)
	assert.Equal(t, sceneIDs[sceneIdxWithMovie], scenes[0].ID)

	scenes, err = sqb.FindByMovieID(0)

	if err != nil {
		t.Fatalf("error calling FindByMovieID: %s", err.Error())
	}

	assert.Len(t, scenes, 0)
}

func TestFindByPerformerID(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()

	scenes, err := sqb.FindByPerformerID(performerIDs[performerIdxWithScene])

	if err != nil {
		t.Fatalf("error calling FindByPerformerID: %s", err.Error())
	}

	assert.Len(t, scenes, 1)
	assert.Equal(t, sceneIDs[sceneIdxWithPerformer], scenes[0].ID)

	scenes, err = sqb.FindByPerformerID(0)

	if err != nil {
		t.Fatalf("error calling FindByPerformerID: %s", err.Error())
	}

	assert.Len(t, scenes, 0)
}

func TestFindByStudioID(t *testing.T) {
	sqb := models.NewSceneQueryBuilder()

	scenes, err := sqb.FindByStudioID(performerIDs[studioIdxWithScene])

	if err != nil {
		t.Fatalf("error calling FindByStudioID: %s", err.Error())
	}

	assert.Len(t, scenes, 1)
	assert.Equal(t, sceneIDs[sceneIdxWithStudio], scenes[0].ID)

	scenes, err = sqb.FindByStudioID(0)

	if err != nil {
		t.Fatalf("error calling FindByStudioID: %s", err.Error())
	}

	assert.Len(t, scenes, 0)
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
