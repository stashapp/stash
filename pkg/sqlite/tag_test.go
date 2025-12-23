//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestMarkerFindBySceneMarkerID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		tqb := db.Tag

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

func TestTagFindByGroupID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		tqb := db.Tag

		groupID := groupIDs[groupIdxWithTag]

		tags, err := tqb.FindByGroupID(ctx, groupID)

		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}

		assert.Len(t, tags, 1)
		assert.Equal(t, tagIDs[tagIdxWithGroup], tags[0].ID)

		tags, err = tqb.FindByGroupID(ctx, 0)

		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}

		assert.Len(t, tags, 0)

		return nil
	})
}

func TestTagFindByName(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		tqb := db.Tag

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

		sqb := db.Tag

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
		tqb := db.Tag

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
		tqb := db.Tag

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
		sqb := db.Tag

		sortBy := "scenes_count"
		dir := models.SortDirectionEnumDesc
		findFilter := &models.FindFilterType{
			Sort:      &sortBy,
			Direction: &dir,
		}

		tags := queryTags(ctx, t, sqb, nil, findFilter)
		assert := assert.New(t)
		assert.Equal(tagIDs[tagIdx2WithScene], tags[0].ID)

		sortBy = "scene_markers_count"
		tags = queryTags(ctx, t, sqb, nil, findFilter)
		assert.Equal(tagIDs[tagIdxWithPrimaryMarkers], tags[0].ID)

		sortBy = "images_count"
		tags = queryTags(ctx, t, sqb, nil, findFilter)
		assert.Equal(tagIDs[tagIdx1WithImage], tags[0].ID)

		sortBy = "galleries_count"
		tags = queryTags(ctx, t, sqb, nil, findFilter)
		assert.Equal(tagIDs[tagIdx1WithGallery], tags[0].ID)

		sortBy = "performers_count"
		tags = queryTags(ctx, t, sqb, nil, findFilter)
		assert.Equal(tagIDs[tagIdx2WithPerformer], tags[0].ID)

		sortBy = "studios_count"
		tags = queryTags(ctx, t, sqb, nil, findFilter)
		assert.Equal(tagIDs[tagIdx2WithStudio], tags[0].ID)

		sortBy = "movies_count"
		tags = queryTags(ctx, t, sqb, nil, findFilter)
		assert.Equal(tagIDs[tagIdx1WithGroup], tags[0].ID)

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
		aliases, err := db.Tag.GetAliases(ctx, tag.ID)
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

	aliasCriterion.Modifier = models.CriterionModifierIsNull
	aliasCriterion.Value = ""
	verifyTagQuery(t, tagFilter, nil, verifyFn)

	aliasCriterion.Modifier = models.CriterionModifierNotNull
	verifyTagQuery(t, tagFilter, nil, verifyFn)
}

func verifyTagQuery(t *testing.T, tagFilter *models.TagFilterType, findFilter *models.FindFilterType, verifyFn func(ctx context.Context, t *models.Tag)) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Tag

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

func tagsToIDs(i []*models.Tag) []int {
	ret := make([]int, len(i))
	for i, v := range i {
		ret[i] = v.ID
	}

	return ret
}

func TestTagQuery(t *testing.T) {
	var (
		endpoint = tagStashID(tagIdxWithPerformer).Endpoint
		stashID  = tagStashID(tagIdxWithPerformer).StashID
	)

	tests := []struct {
		name        string
		findFilter  *models.FindFilterType
		filter      *models.TagFilterType
		includeIdxs []int
		excludeIdxs []int
		wantErr     bool
	}{
		{
			"stash id with endpoint",
			nil,
			&models.TagFilterType{
				StashIDEndpoint: &models.StashIDCriterionInput{
					Endpoint: &endpoint,
					StashID:  &stashID,
					Modifier: models.CriterionModifierEquals,
				},
			},
			[]int{tagIdxWithPerformer},
			nil,
			false,
		},
		{
			"exclude stash id with endpoint",
			nil,
			&models.TagFilterType{
				StashIDEndpoint: &models.StashIDCriterionInput{
					Endpoint: &endpoint,
					StashID:  &stashID,
					Modifier: models.CriterionModifierNotEquals,
				},
			},
			nil,
			[]int{tagIdxWithPerformer},
			false,
		},
		{
			"null stash id with endpoint",
			nil,
			&models.TagFilterType{
				StashIDEndpoint: &models.StashIDCriterionInput{
					Endpoint: &endpoint,
					Modifier: models.CriterionModifierIsNull,
				},
			},
			nil,
			[]int{tagIdxWithPerformer},
			false,
		},
		{
			"not null stash id with endpoint",
			nil,
			&models.TagFilterType{
				StashIDEndpoint: &models.StashIDCriterionInput{
					Endpoint: &endpoint,
					Modifier: models.CriterionModifierNotNull,
				},
			},
			[]int{tagIdxWithPerformer},
			nil,
			false,
		},
	}

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			tags, _, err := db.Tag.Query(ctx, tt.filter, tt.findFilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("PerformerStore.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ids := tagsToIDs(tags)
			include := indexesToIDs(tagIDs, tt.includeIdxs)
			exclude := indexesToIDs(tagIDs, tt.excludeIdxs)

			for _, i := range include {
				assert.Contains(ids, i)
			}
			for _, e := range exclude {
				assert.NotContains(ids, e)
			}
		})
	}
}

func TestTagQueryIsMissingImage(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		qb := db.Tag
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
		qb := db.Tag
		tagFilter := models.TagFilterType{
			SceneCount: &sceneCountCriterion,
		}

		tags, _, err := qb.Query(ctx, &tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt(t, getTagSceneCount(tag.ID), sceneCountCriterion)
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
		qb := db.Tag
		tagFilter := models.TagFilterType{
			MarkerCount: &markerCountCriterion,
		}

		tags, _, err := qb.Query(ctx, &tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt(t, getTagMarkerCount(tag.ID), markerCountCriterion)
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
		qb := db.Tag
		tagFilter := models.TagFilterType{
			ImageCount: &imageCountCriterion,
		}

		tags, _, err := qb.Query(ctx, &tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt(t, getTagImageCount(tag.ID), imageCountCriterion)
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
		qb := db.Tag
		tagFilter := models.TagFilterType{
			GalleryCount: &imageCountCriterion,
		}

		tags, _, err := qb.Query(ctx, &tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt(t, getTagGalleryCount(tag.ID), imageCountCriterion)
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
		qb := db.Tag
		tagFilter := models.TagFilterType{
			PerformerCount: &imageCountCriterion,
		}

		tags, _, err := qb.Query(ctx, &tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt(t, getTagPerformerCount(tag.ID), imageCountCriterion)
		}

		return nil
	})
}

func TestTagQueryStudioCount(t *testing.T) {
	countCriterion := models.IntCriterionInput{
		Value:    1,
		Modifier: models.CriterionModifierEquals,
	}

	verifyTagStudioCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierNotEquals
	verifyTagStudioCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierLessThan
	verifyTagStudioCount(t, countCriterion)

	countCriterion.Value = 0
	countCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyTagStudioCount(t, countCriterion)
}

func verifyTagStudioCount(t *testing.T, imageCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		qb := db.Tag
		tagFilter := models.TagFilterType{
			StudioCount: &imageCountCriterion,
		}

		tags, _, err := qb.Query(ctx, &tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt(t, getTagStudioCount(tag.ID), imageCountCriterion)
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
		qb := db.Tag
		tagFilter := models.TagFilterType{
			ParentCount: &sceneCountCriterion,
		}

		tags := queryTags(ctx, t, qb, &tagFilter, nil)

		if len(tags) == 0 {
			t.Error("Expected at least one tag")
		}

		for _, tag := range tags {
			verifyInt(t, getTagParentCount(tag.ID), sceneCountCriterion)
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
		qb := db.Tag
		tagFilter := models.TagFilterType{
			ChildCount: &sceneCountCriterion,
		}

		tags := queryTags(ctx, t, qb, &tagFilter, nil)

		if len(tags) == 0 {
			t.Error("Expected at least one tag")
		}

		for _, tag := range tags {
			verifyInt(t, getTagChildCount(tag.ID), sceneCountCriterion)
		}

		return nil
	})
}

func TestTagQueryParent(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		const nameField = "Name"
		sqb := db.Tag
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
		assert.Equal(t, tagIDs[tagIdxWithParentTag], tags[0].ID)

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

		sqb := db.Tag
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
		qb := db.Tag

		// create tag to test against
		const name = "TestTagUpdateTagImage"
		tag := models.Tag{
			Name: name,
		}
		err := qb.Create(ctx, &tag)
		if err != nil {
			return fmt.Errorf("Error creating tag: %s", err.Error())
		}

		return testUpdateImage(t, ctx, tag.ID, qb.UpdateImage, qb.GetImage)
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestTagUpdateAlias(t *testing.T) {
	if err := withTxn(func(ctx context.Context) error {
		qb := db.Tag

		// create tag to test against
		const name = "TestTagUpdateAlias"
		tag := models.Tag{
			Name: name,
		}
		err := qb.Create(ctx, &tag)
		if err != nil {
			return fmt.Errorf("Error creating tag: %s", err.Error())
		}

		aliases := []string{"alias1", "alias2"}
		err = qb.UpdateAliases(ctx, tag.ID, aliases)
		if err != nil {
			return fmt.Errorf("Error updating tag aliases: %s", err.Error())
		}

		// ensure aliases set
		storedAliases, err := qb.GetAliases(ctx, tag.ID)
		if err != nil {
			return fmt.Errorf("Error getting aliases: %s", err.Error())
		}
		assert.Equal(t, aliases, storedAliases)

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestTagStashIDs(t *testing.T) {
	if err := withTxn(func(ctx context.Context) error {
		qb := db.Tag

		// create tag to test against
		const name = "TestTagStashIDs"
		tag := models.Tag{
			Name: name,
		}
		err := qb.Create(ctx, &tag)
		if err != nil {
			return fmt.Errorf("Error creating tag: %s", err.Error())
		}

		testStashIDReaderWriter(ctx, t, qb, tag.ID)

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestTagFindByStashID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		qb := db.Tag

		// create tag to test against
		const name = "TestTagFindByStashID"
		const stashID = "stashid"
		const endpoint = "endpoint"
		tag := models.Tag{
			Name:     name,
			StashIDs: models.NewRelatedStashIDs([]models.StashID{{StashID: stashID, Endpoint: endpoint}}),
		}
		err := qb.Create(ctx, &tag)
		if err != nil {
			return fmt.Errorf("Error creating tag: %s", err.Error())
		}

		// find by stash ID
		tags, err := qb.FindByStashID(ctx, models.StashID{StashID: stashID, Endpoint: endpoint})
		if err != nil {
			return fmt.Errorf("Error finding by stash ID: %s", err.Error())
		}

		assert.Len(t, tags, 1)
		assert.Equal(t, tag.ID, tags[0].ID)

		// find by non-existent stash ID
		tags, err = qb.FindByStashID(ctx, models.StashID{StashID: "nonexistent", Endpoint: endpoint})
		if err != nil {
			return fmt.Errorf("Error finding by stash ID: %s", err.Error())
		}

		assert.Len(t, tags, 0)

		return nil
	})
}

func TestTagMerge(t *testing.T) {
	assert := assert.New(t)

	// merge tests - perform these in a transaction that we'll rollback
	if err := withRollbackTxn(func(ctx context.Context) error {
		qb := db.Tag
		mqb := db.SceneMarker

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
			tagIdxWithStudio,
			tagIdx1WithStudio,
			tagIdx2WithStudio,
			tagIdxWithGallery,
			tagIdx1WithGallery,
			tagIdx2WithGallery,
			tagIdx1WithGroup,
			tagIdx2WithGroup,
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
		s, err := db.Scene.Find(ctx, sceneIDs[sceneIdxWithTwoTags])
		if err != nil {
			return err
		}
		if err := s.LoadTagIDs(ctx, db.Scene); err != nil {
			return err
		}
		sceneTagIDs := s.TagIDs.List()

		assert.Contains(sceneTagIDs, destID)

		// ensure marker points to new tag
		marker, err := mqb.Find(ctx, markerIDs[markerIdxWithTag])
		if err != nil {
			return err
		}

		assert.Equal(destID, marker.PrimaryTagID)

		markerTagIDs, err := mqb.GetTagIDs(ctx, marker.ID)
		if err != nil {
			return err
		}

		assert.Contains(markerTagIDs, destID)

		// ensure image points to new tag
		imageTagIDs, err := db.Image.GetTagIDs(ctx, imageIDs[imageIdxWithTwoTags])
		if err != nil {
			return err
		}

		assert.Contains(imageTagIDs, destID)

		g, err := db.Gallery.Find(ctx, galleryIDs[galleryIdxWithTwoTags])
		if err != nil {
			return err
		}

		if err := g.LoadTagIDs(ctx, db.Gallery); err != nil {
			return err
		}

		// ensure gallery points to new tag
		assert.Contains(g.TagIDs.List(), destID)

		// ensure performer points to new tag
		performerTagIDs, err := db.Performer.GetTagIDs(ctx, performerIDs[performerIdxWithTwoTags])
		if err != nil {
			return err
		}

		assert.Contains(performerTagIDs, destID)

		// ensure studio points to new tag
		studioTagIDs, err := db.Studio.GetTagIDs(ctx, studioIDs[studioIdxWithTwoTags])
		if err != nil {
			return err
		}

		assert.Contains(studioTagIDs, destID)

		// ensure group points to new tag
		group, err := db.Group.Find(ctx, groupIDs[groupIdxWithTwoTags])
		if err != nil {
			return err
		}
		if err := group.LoadTagIDs(ctx, db.Group); err != nil {
			return err
		}
		groupTagIDs := group.TagIDs.List()

		assert.Contains(groupTagIDs, destID)

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
