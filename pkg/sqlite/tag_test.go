//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestMarkerFindBySceneMarkerID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		tqb := sqlite.TagReaderWriter

		markerID := markerIDs[markerIdxWithTag]

		tags, err := tqb.FindBySceneMarkerID(ctx, markerID)

		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}

		assert.Len(t, tags, 1)
		assert.Equal(t, tagIDs[tagIdxWithMarkers], tags[0].ID)

		tags, err = tqb.FindBySceneMarkerID(ctx, 0)

		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}

		assert.Len(t, tags, 0)

		return nil
	})
}

func TestTagFindByName(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		tqb := sqlite.TagReaderWriter

		name := tagNames[tagIdxWithScene] // find a tag by name

		tag, err := tqb.FindByName(ctx, name, false)

		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}

		assert.Equal(t, tagNames[tagIdxWithScene], tag.Name)

		name = tagNames[tagIdxWithDupName] // find a tag by name nocase

		tag, err = tqb.FindByName(ctx, name, true)

		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}
		// tagIdxWithDupName and tagIdxWithScene should have similar names ( only diff should be Name vs NaMe)
		//tag.Name should match with tagIdxWithScene since its ID is before tagIdxWithDupName
		assert.Equal(t, tagNames[tagIdxWithScene], tag.Name)
		//tag.Name should match with tagIdxWithDupName if the check is not case sensitive
		assert.Equal(t, strings.ToLower(tagNames[tagIdxWithDupName]), strings.ToLower(tag.Name))

		return nil
	})
}

func TestTagQueryIgnoreAutoTag(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		ignoreAutoTag := true
		tagFilter := models.TagFilterType{
			IgnoreAutoTag: &ignoreAutoTag,
		}

		sqb := sqlite.TagReaderWriter

		tags := queryTags(ctx, t, sqb, &tagFilter, nil)

		assert.Len(t, tags, int(math.Ceil(float64(totalTags)/5)))
		for _, s := range tags {
			assert.True(t, s.IgnoreAutoTag)
		}

		return nil
	})
}

func TestTagQueryForAutoTag(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		tqb := sqlite.TagReaderWriter

		name := tagNames[tagIdx1WithScene] // find a tag by name

		tags, err := tqb.QueryForAutoTag(ctx, []string{name})

		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}

		assert.Len(t, tags, 2)
		lcName := tagNames[tagIdx1WithScene]
		assert.Equal(t, strings.ToLower(lcName), strings.ToLower(tags[0].Name))
		assert.Equal(t, strings.ToLower(lcName), strings.ToLower(tags[1].Name))

		// find by alias
		name = getTagStringValue(tagIdx1WithScene, "Alias")
		tags, err = tqb.QueryForAutoTag(ctx, []string{name})

		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}

		assert.Len(t, tags, 1)
		assert.Equal(t, tagIDs[tagIdx1WithScene], tags[0].ID)

		return nil
	})
}

func TestTagFindByNames(t *testing.T) {
	var names []string

	withTxn(func(ctx context.Context) error {
		tqb := sqlite.TagReaderWriter

		names = append(names, tagNames[tagIdxWithScene]) // find tags by names

		tags, err := tqb.FindByNames(ctx, names, false)
		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}
		assert.Len(t, tags, 1)
		assert.Equal(t, tagNames[tagIdxWithScene], tags[0].Name)

		tags, err = tqb.FindByNames(ctx, names, true) // find tags by names nocase
		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}
		assert.Len(t, tags, 2) // tagIdxWithScene and tagIdxWithDupName
		assert.Equal(t, strings.ToLower(tagNames[tagIdxWithScene]), strings.ToLower(tags[0].Name))
		assert.Equal(t, strings.ToLower(tagNames[tagIdxWithScene]), strings.ToLower(tags[1].Name))

		names = append(names, tagNames[tagIdx1WithScene]) // find tags by names ( 2 names )

		tags, err = tqb.FindByNames(ctx, names, false)
		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}
		assert.Len(t, tags, 2) // tagIdxWithScene and tagIdx1WithScene
		assert.Equal(t, tagNames[tagIdxWithScene], tags[0].Name)
		assert.Equal(t, tagNames[tagIdx1WithScene], tags[1].Name)

		tags, err = tqb.FindByNames(ctx, names, true) // find tags by names ( 2 names nocase)
		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}
		assert.Len(t, tags, 4) // tagIdxWithScene and tagIdxWithDupName , tagIdx1WithScene and tagIdx1WithDupName
		assert.Equal(t, tagNames[tagIdxWithScene], tags[0].Name)
		assert.Equal(t, tagNames[tagIdx1WithScene], tags[1].Name)
		assert.Equal(t, tagNames[tagIdx1WithDupName], tags[2].Name)
		assert.Equal(t, tagNames[tagIdxWithDupName], tags[3].Name)

		return nil
	})
}

func TestTagQuerySort(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.TagReaderWriter

		sortBy := "scenes_count"
		dir := models.SortDirectionEnumDesc
		findFilter := &models.FindFilterType{
			Sort:      &sortBy,
			Direction: &dir,
		}

		tags := queryTags(ctx, t, sqb, nil, findFilter)
		assert := assert.New(t)
		assert.Equal(tagIDs[tagIdxWithScene], tags[0].ID)

		sortBy = "scene_markers_count"
		tags = queryTags(ctx, t, sqb, nil, findFilter)
		assert.Equal(tagIDs[tagIdxWithMarkers], tags[0].ID)

		sortBy = "images_count"
		tags = queryTags(ctx, t, sqb, nil, findFilter)
		assert.Equal(tagIDs[tagIdxWithImage], tags[0].ID)

		sortBy = "galleries_count"
		tags = queryTags(ctx, t, sqb, nil, findFilter)
		assert.Equal(tagIDs[tagIdxWithGallery], tags[0].ID)

		sortBy = "performers_count"
		tags = queryTags(ctx, t, sqb, nil, findFilter)
		assert.Equal(tagIDs[tagIdxWithPerformer], tags[0].ID)

		return nil
	})
}

func TestTagQueryName(t *testing.T) {
	const tagIdx = 1
	tagName := getSceneStringValue(tagIdx, "Name")

	nameCriterion := &models.StringCriterionInput{
		Value:    tagName,
		Modifier: models.CriterionModifierEquals,
	}

	tagFilter := &models.TagFilterType{
		Name: nameCriterion,
	}

	verifyFn := func(ctx context.Context, tag *models.Tag) {
		verifyString(t, tag.Name, *nameCriterion)
	}

	verifyTagQuery(t, tagFilter, nil, verifyFn)

	nameCriterion.Modifier = models.CriterionModifierNotEquals
	verifyTagQuery(t, tagFilter, nil, verifyFn)

	nameCriterion.Modifier = models.CriterionModifierMatchesRegex
	nameCriterion.Value = "tag_.*1_Name"
	verifyTagQuery(t, tagFilter, nil, verifyFn)

	nameCriterion.Modifier = models.CriterionModifierNotMatchesRegex
	verifyTagQuery(t, tagFilter, nil, verifyFn)
}

func TestTagQueryAlias(t *testing.T) {
	const tagIdx = 1
	tagName := getSceneStringValue(tagIdx, "Alias")

	aliasCriterion := &models.StringCriterionInput{
		Value:    tagName,
		Modifier: models.CriterionModifierEquals,
	}

	tagFilter := &models.TagFilterType{
		Aliases: aliasCriterion,
	}

	verifyFn := func(ctx context.Context, tag *models.Tag) {
		aliases, err := sqlite.TagReaderWriter.GetAliases(ctx, tag.ID)
		if err != nil {
			t.Errorf("Error querying tags: %s", err.Error())
		}

		var alias string
		if len(aliases) > 0 {
			alias = aliases[0]
		}

		verifyString(t, alias, *aliasCriterion)
	}

	verifyTagQuery(t, tagFilter, nil, verifyFn)

	aliasCriterion.Modifier = models.CriterionModifierNotEquals
	verifyTagQuery(t, tagFilter, nil, verifyFn)

	aliasCriterion.Modifier = models.CriterionModifierMatchesRegex
	aliasCriterion.Value = "tag_.*1_Alias"
	verifyTagQuery(t, tagFilter, nil, verifyFn)

	aliasCriterion.Modifier = models.CriterionModifierNotMatchesRegex
	verifyTagQuery(t, tagFilter, nil, verifyFn)
}

func verifyTagQuery(t *testing.T, tagFilter *models.TagFilterType, findFilter *models.FindFilterType, verifyFn func(ctx context.Context, t *models.Tag)) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.TagReaderWriter

		tags := queryTags(ctx, t, sqb, tagFilter, findFilter)

		for _, tag := range tags {
			verifyFn(ctx, tag)
		}

		return nil
	})
}

func queryTags(ctx context.Context, t *testing.T, qb models.TagReader, tagFilter *models.TagFilterType, findFilter *models.FindFilterType) []*models.Tag {
	t.Helper()
	tags, _, err := qb.Query(ctx, tagFilter, findFilter)
	if err != nil {
		t.Errorf("Error querying tags: %s", err.Error())
	}

	return tags
}

func TestTagQueryIsMissingImage(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		qb := sqlite.TagReaderWriter
		isMissing := "image"
		tagFilter := models.TagFilterType{
			IsMissing: &isMissing,
		}

		q := getTagStringValue(tagIdxWithCoverImage, "name")
		findFilter := models.FindFilterType{
			Q: &q,
		}

		tags, _, err := qb.Query(ctx, &tagFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		assert.Len(t, tags, 0)

		findFilter.Q = nil
		tags, _, err = qb.Query(ctx, &tagFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		// ensure non of the ids equal the one with image
		for _, tag := range tags {
			assert.NotEqual(t, tagIDs[tagIdxWithCoverImage], tag.ID)
		}

		return nil
	})
}

func TestTagQuerySceneCount(t *testing.T) {
	countCriterion := models.IntCriterionInput{
		Value:    1,
		Modifier: models.CriterionModifierEquals,
	}

	verifyTagSceneCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierNotEquals
	verifyTagSceneCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierLessThan
	verifyTagSceneCount(t, countCriterion)

	countCriterion.Value = 0
	countCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyTagSceneCount(t, countCriterion)
}

func verifyTagSceneCount(t *testing.T, sceneCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		qb := sqlite.TagReaderWriter
		tagFilter := models.TagFilterType{
			SceneCount: &sceneCountCriterion,
		}

		tags, _, err := qb.Query(ctx, &tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt64(t, sql.NullInt64{
				Int64: int64(getTagSceneCount(tag.ID)),
				Valid: true,
			}, sceneCountCriterion)
		}

		return nil
	})
}

func TestTagQueryMarkerCount(t *testing.T) {
	countCriterion := models.IntCriterionInput{
		Value:    1,
		Modifier: models.CriterionModifierEquals,
	}

	verifyTagMarkerCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierNotEquals
	verifyTagMarkerCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierLessThan
	verifyTagMarkerCount(t, countCriterion)

	countCriterion.Value = 0
	countCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyTagMarkerCount(t, countCriterion)
}

func verifyTagMarkerCount(t *testing.T, markerCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		qb := sqlite.TagReaderWriter
		tagFilter := models.TagFilterType{
			MarkerCount: &markerCountCriterion,
		}

		tags, _, err := qb.Query(ctx, &tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt64(t, sql.NullInt64{
				Int64: int64(getTagMarkerCount(tag.ID)),
				Valid: true,
			}, markerCountCriterion)
		}

		return nil
	})
}

func TestTagQueryImageCount(t *testing.T) {
	countCriterion := models.IntCriterionInput{
		Value:    1,
		Modifier: models.CriterionModifierEquals,
	}

	verifyTagImageCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierNotEquals
	verifyTagImageCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierLessThan
	verifyTagImageCount(t, countCriterion)

	countCriterion.Value = 0
	countCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyTagImageCount(t, countCriterion)
}

func verifyTagImageCount(t *testing.T, imageCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		qb := sqlite.TagReaderWriter
		tagFilter := models.TagFilterType{
			ImageCount: &imageCountCriterion,
		}

		tags, _, err := qb.Query(ctx, &tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt64(t, sql.NullInt64{
				Int64: int64(getTagImageCount(tag.ID)),
				Valid: true,
			}, imageCountCriterion)
		}

		return nil
	})
}

func TestTagQueryGalleryCount(t *testing.T) {
	countCriterion := models.IntCriterionInput{
		Value:    1,
		Modifier: models.CriterionModifierEquals,
	}

	verifyTagGalleryCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierNotEquals
	verifyTagGalleryCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierLessThan
	verifyTagGalleryCount(t, countCriterion)

	countCriterion.Value = 0
	countCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyTagGalleryCount(t, countCriterion)
}

func verifyTagGalleryCount(t *testing.T, imageCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		qb := sqlite.TagReaderWriter
		tagFilter := models.TagFilterType{
			GalleryCount: &imageCountCriterion,
		}

		tags, _, err := qb.Query(ctx, &tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt64(t, sql.NullInt64{
				Int64: int64(getTagGalleryCount(tag.ID)),
				Valid: true,
			}, imageCountCriterion)
		}

		return nil
	})
}

func TestTagQueryPerformerCount(t *testing.T) {
	countCriterion := models.IntCriterionInput{
		Value:    1,
		Modifier: models.CriterionModifierEquals,
	}

	verifyTagPerformerCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierNotEquals
	verifyTagPerformerCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierLessThan
	verifyTagPerformerCount(t, countCriterion)

	countCriterion.Value = 0
	countCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyTagPerformerCount(t, countCriterion)
}

func verifyTagPerformerCount(t *testing.T, imageCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		qb := sqlite.TagReaderWriter
		tagFilter := models.TagFilterType{
			PerformerCount: &imageCountCriterion,
		}

		tags, _, err := qb.Query(ctx, &tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt64(t, sql.NullInt64{
				Int64: int64(getTagPerformerCount(tag.ID)),
				Valid: true,
			}, imageCountCriterion)
		}

		return nil
	})
}

func TestTagQueryParentCount(t *testing.T) {
	countCriterion := models.IntCriterionInput{
		Value:    1,
		Modifier: models.CriterionModifierEquals,
	}

	verifyTagParentCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierNotEquals
	verifyTagParentCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierLessThan
	verifyTagParentCount(t, countCriterion)

	countCriterion.Value = 0
	countCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyTagParentCount(t, countCriterion)
}

func verifyTagParentCount(t *testing.T, sceneCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		qb := sqlite.TagReaderWriter
		tagFilter := models.TagFilterType{
			ParentCount: &sceneCountCriterion,
		}

		tags := queryTags(ctx, t, qb, &tagFilter, nil)

		if len(tags) == 0 {
			t.Error("Expected at least one tag")
		}

		for _, tag := range tags {
			verifyInt64(t, sql.NullInt64{
				Int64: int64(getTagParentCount(tag.ID)),
				Valid: true,
			}, sceneCountCriterion)
		}

		return nil
	})
}

func TestTagQueryChildCount(t *testing.T) {
	countCriterion := models.IntCriterionInput{
		Value:    1,
		Modifier: models.CriterionModifierEquals,
	}

	verifyTagChildCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierNotEquals
	verifyTagChildCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierLessThan
	verifyTagChildCount(t, countCriterion)

	countCriterion.Value = 0
	countCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyTagChildCount(t, countCriterion)
}

func verifyTagChildCount(t *testing.T, sceneCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		qb := sqlite.TagReaderWriter
		tagFilter := models.TagFilterType{
			ChildCount: &sceneCountCriterion,
		}

		tags := queryTags(ctx, t, qb, &tagFilter, nil)

		if len(tags) == 0 {
			t.Error("Expected at least one tag")
		}

		for _, tag := range tags {
			verifyInt64(t, sql.NullInt64{
				Int64: int64(getTagChildCount(tag.ID)),
				Valid: true,
			}, sceneCountCriterion)
		}

		return nil
	})
}

func TestTagQueryParent(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		const nameField = "Name"
		sqb := sqlite.TagReaderWriter
		tagCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdxWithChildTag]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		tagFilter := models.TagFilterType{
			Parents: &tagCriterion,
		}

		tags := queryTags(ctx, t, sqb, &tagFilter, nil)

		assert.Len(t, tags, 1)

		// ensure id is correct
		assert.Equal(t, sceneIDs[tagIdxWithParentTag], tags[0].ID)

		tagCriterion.Modifier = models.CriterionModifierExcludes

		q := getTagStringValue(tagIdxWithParentTag, nameField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		tags = queryTags(ctx, t, sqb, &tagFilter, &findFilter)
		assert.Len(t, tags, 0)

		depth := -1

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdxWithGrandChild]),
			},
			Modifier: models.CriterionModifierIncludes,
			Depth:    &depth,
		}

		tags = queryTags(ctx, t, sqb, &tagFilter, nil)
		assert.Len(t, tags, 2)

		depth = 1

		tags = queryTags(ctx, t, sqb, &tagFilter, nil)
		assert.Len(t, tags, 2)

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Modifier: models.CriterionModifierIsNull,
		}
		q = getTagStringValue(tagIdxWithGallery, nameField)

		tags = queryTags(ctx, t, sqb, &tagFilter, &findFilter)
		assert.Len(t, tags, 1)
		assert.Equal(t, tagIDs[tagIdxWithGallery], tags[0].ID)

		q = getTagStringValue(tagIdxWithParentTag, nameField)
		tags = queryTags(ctx, t, sqb, &tagFilter, &findFilter)
		assert.Len(t, tags, 0)

		tagCriterion.Modifier = models.CriterionModifierNotNull

		tags = queryTags(ctx, t, sqb, &tagFilter, &findFilter)
		assert.Len(t, tags, 1)
		assert.Equal(t, tagIDs[tagIdxWithParentTag], tags[0].ID)

		q = getTagStringValue(tagIdxWithGallery, nameField)
		tags = queryTags(ctx, t, sqb, &tagFilter, &findFilter)
		assert.Len(t, tags, 0)

		return nil
	})
}

func TestTagQueryChild(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		const nameField = "Name"

		sqb := sqlite.TagReaderWriter
		tagCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdxWithParentTag]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		tagFilter := models.TagFilterType{
			Children: &tagCriterion,
		}

		tags := queryTags(ctx, t, sqb, &tagFilter, nil)

		assert.Len(t, tags, 1)

		// ensure id is correct
		assert.Equal(t, sceneIDs[tagIdxWithChildTag], tags[0].ID)

		tagCriterion.Modifier = models.CriterionModifierExcludes

		q := getTagStringValue(tagIdxWithChildTag, nameField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		tags = queryTags(ctx, t, sqb, &tagFilter, &findFilter)
		assert.Len(t, tags, 0)

		depth := -1

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdxWithGrandParent]),
			},
			Modifier: models.CriterionModifierIncludes,
			Depth:    &depth,
		}

		tags = queryTags(ctx, t, sqb, &tagFilter, nil)
		assert.Len(t, tags, 2)

		depth = 1

		tags = queryTags(ctx, t, sqb, &tagFilter, nil)
		assert.Len(t, tags, 2)

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Modifier: models.CriterionModifierIsNull,
		}
		q = getTagStringValue(tagIdxWithGallery, nameField)

		tags = queryTags(ctx, t, sqb, &tagFilter, &findFilter)
		assert.Len(t, tags, 1)
		assert.Equal(t, tagIDs[tagIdxWithGallery], tags[0].ID)

		q = getTagStringValue(tagIdxWithChildTag, nameField)
		tags = queryTags(ctx, t, sqb, &tagFilter, &findFilter)
		assert.Len(t, tags, 0)

		tagCriterion.Modifier = models.CriterionModifierNotNull

		tags = queryTags(ctx, t, sqb, &tagFilter, &findFilter)
		assert.Len(t, tags, 1)
		assert.Equal(t, tagIDs[tagIdxWithChildTag], tags[0].ID)

		q = getTagStringValue(tagIdxWithGallery, nameField)
		tags = queryTags(ctx, t, sqb, &tagFilter, &findFilter)
		assert.Len(t, tags, 0)

		return nil
	})
}

func TestTagUpdateTagImage(t *testing.T) {
	if err := withTxn(func(ctx context.Context) error {
		qb := sqlite.TagReaderWriter

		// create tag to test against
		const name = "TestTagUpdateTagImage"
		tag := models.Tag{
			Name: name,
		}
		created, err := qb.Create(ctx, tag)
		if err != nil {
			return fmt.Errorf("Error creating tag: %s", err.Error())
		}

		image := []byte("image")
		err = qb.UpdateImage(ctx, created.ID, image)
		if err != nil {
			return fmt.Errorf("Error updating studio image: %s", err.Error())
		}

		// ensure image set
		storedImage, err := qb.GetImage(ctx, created.ID)
		if err != nil {
			return fmt.Errorf("Error getting image: %s", err.Error())
		}
		assert.Equal(t, storedImage, image)

		// set nil image
		err = qb.UpdateImage(ctx, created.ID, nil)
		if err == nil {
			return fmt.Errorf("Expected error setting nil image")
		}

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestTagDestroyTagImage(t *testing.T) {
	if err := withTxn(func(ctx context.Context) error {
		qb := sqlite.TagReaderWriter

		// create performer to test against
		const name = "TestTagDestroyTagImage"
		tag := models.Tag{
			Name: name,
		}
		created, err := qb.Create(ctx, tag)
		if err != nil {
			return fmt.Errorf("Error creating tag: %s", err.Error())
		}

		image := []byte("image")
		err = qb.UpdateImage(ctx, created.ID, image)
		if err != nil {
			return fmt.Errorf("Error updating studio image: %s", err.Error())
		}

		err = qb.DestroyImage(ctx, created.ID)
		if err != nil {
			return fmt.Errorf("Error destroying studio image: %s", err.Error())
		}

		// image should be nil
		storedImage, err := qb.GetImage(ctx, created.ID)
		if err != nil {
			return fmt.Errorf("Error getting image: %s", err.Error())
		}
		assert.Nil(t, storedImage)

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestTagUpdateAlias(t *testing.T) {
	if err := withTxn(func(ctx context.Context) error {
		qb := sqlite.TagReaderWriter

		// create tag to test against
		const name = "TestTagUpdateAlias"
		tag := models.Tag{
			Name: name,
		}
		created, err := qb.Create(ctx, tag)
		if err != nil {
			return fmt.Errorf("Error creating tag: %s", err.Error())
		}

		aliases := []string{"alias1", "alias2"}
		err = qb.UpdateAliases(ctx, created.ID, aliases)
		if err != nil {
			return fmt.Errorf("Error updating tag aliases: %s", err.Error())
		}

		// ensure aliases set
		storedAliases, err := qb.GetAliases(ctx, created.ID)
		if err != nil {
			return fmt.Errorf("Error getting aliases: %s", err.Error())
		}
		assert.Equal(t, aliases, storedAliases)

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestTagMerge(t *testing.T) {
	assert := assert.New(t)

	// merge tests - perform these in a transaction that we'll rollback
	if err := withRollbackTxn(func(ctx context.Context) error {
		qb := sqlite.TagReaderWriter

		// try merging into same tag
		err := qb.Merge(ctx, []int{tagIDs[tagIdx1WithScene]}, tagIDs[tagIdx1WithScene])
		assert.NotNil(err)

		// merge everything into tagIdxWithScene
		srcIdxs := []int{
			tagIdx1WithScene,
			tagIdx2WithScene,
			tagIdxWithPrimaryMarkers,
			tagIdxWithMarkers,
			tagIdxWithCoverImage,
			tagIdxWithImage,
			tagIdx1WithImage,
			tagIdx2WithImage,
			tagIdxWithPerformer,
			tagIdx1WithPerformer,
			tagIdx2WithPerformer,
			tagIdxWithGallery,
			tagIdx1WithGallery,
			tagIdx2WithGallery,
		}
		var srcIDs []int
		for _, idx := range srcIdxs {
			srcIDs = append(srcIDs, tagIDs[idx])
		}

		destID := tagIDs[tagIdxWithScene]
		if err = qb.Merge(ctx, srcIDs, destID); err != nil {
			return err
		}

		// ensure other tags are deleted
		for _, tagId := range srcIDs {
			t, err := qb.Find(ctx, tagId)
			if err != nil {
				return err
			}

			assert.Nil(t)
		}

		// ensure aliases are set on the destination
		destAliases, err := qb.GetAliases(ctx, destID)
		if err != nil {
			return err
		}
		for _, tagIdx := range srcIdxs {
			assert.Contains(destAliases, getTagStringValue(tagIdx, "Name"))
		}

		// ensure scene points to new tag
		sceneTagIDs, err := sqlite.SceneReaderWriter.GetTagIDs(ctx, sceneIDs[sceneIdxWithTwoTags])
		if err != nil {
			return err
		}

		assert.Contains(sceneTagIDs, destID)

		// ensure marker points to new tag
		marker, err := sqlite.SceneMarkerReaderWriter.Find(ctx, markerIDs[markerIdxWithTag])
		if err != nil {
			return err
		}

		assert.Equal(destID, marker.PrimaryTagID)

		markerTagIDs, err := sqlite.SceneMarkerReaderWriter.GetTagIDs(ctx, marker.ID)
		if err != nil {
			return err
		}

		assert.Contains(markerTagIDs, destID)

		// ensure image points to new tag
		imageTagIDs, err := sqlite.ImageReaderWriter.GetTagIDs(ctx, imageIDs[imageIdxWithTwoTags])
		if err != nil {
			return err
		}

		assert.Contains(imageTagIDs, destID)

		// ensure gallery points to new tag
		galleryTagIDs, err := sqlite.GalleryReaderWriter.GetTagIDs(ctx, galleryIDs[galleryIdxWithTwoTags])
		if err != nil {
			return err
		}

		assert.Contains(galleryTagIDs, destID)

		// ensure performer points to new tag
		performerTagIDs, err := sqlite.GalleryReaderWriter.GetTagIDs(ctx, performerIDs[performerIdxWithTwoTags])
		if err != nil {
			return err
		}

		assert.Contains(performerTagIDs, destID)

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

// TODO Create
// TODO Update
// TODO Destroy
// TODO Find
// TODO FindBySceneID
// TODO FindBySceneMarkerID
// TODO Count
// TODO All
// TODO AllSlim
// TODO Query
