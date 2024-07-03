//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/models"
)

func loadGroupRelationships(ctx context.Context, expected models.Group, actual *models.Group) error {
	if expected.URLs.Loaded() {
		if err := actual.LoadURLs(ctx, db.Group); err != nil {
			return err
		}
	}
	if expected.TagIDs.Loaded() {
		if err := actual.LoadTagIDs(ctx, db.Group); err != nil {
			return err
		}
	}

	return nil
}

func Test_GroupStore_Create(t *testing.T) {
	var (
		name      = "name"
		url       = "url"
		aliases   = "alias1, alias2"
		director  = "director"
		rating    = 60
		duration  = 34
		synopsis  = "synopsis"
		date, _   = models.ParseDate("2003-02-01")
		createdAt = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name      string
		newObject models.Group
		wantErr   bool
	}{
		{
			"full",
			models.Group{
				Name:      name,
				Duration:  &duration,
				Date:      &date,
				Rating:    &rating,
				StudioID:  &studioIDs[studioIdxWithGroup],
				Director:  director,
				Synopsis:  synopsis,
				URLs:      models.NewRelatedStrings([]string{url}),
				TagIDs:    models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithGroup]}),
				Aliases:   aliases,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"invalid tag id",
			models.Group{
				Name:   name,
				TagIDs: models.NewRelatedIDs([]int{invalidID}),
			},
			true,
		},
	}

	qb := db.Group

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			p := tt.newObject
			if err := qb.Create(ctx, &p); (err != nil) != tt.wantErr {
				t.Errorf("GroupStore.Create() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				assert.Zero(p.ID)
				return
			}

			assert.NotZero(p.ID)

			copy := tt.newObject
			copy.ID = p.ID

			// load relationships
			if err := loadGroupRelationships(ctx, copy, &p); err != nil {
				t.Errorf("loadGroupRelationships() error = %v", err)
				return
			}

			assert.Equal(copy, p)

			// ensure can find the group
			found, err := qb.Find(ctx, p.ID)
			if err != nil {
				t.Errorf("GroupStore.Find() error = %v", err)
			}

			if !assert.NotNil(found) {
				return
			}

			// load relationships
			if err := loadGroupRelationships(ctx, copy, found); err != nil {
				t.Errorf("loadGroupRelationships() error = %v", err)
				return
			}
			assert.Equal(copy, *found)

			return
		})
	}
}

func Test_groupQueryBuilder_Update(t *testing.T) {
	var (
		name      = "name"
		url       = "url"
		aliases   = "alias1, alias2"
		director  = "director"
		rating    = 60
		duration  = 34
		synopsis  = "synopsis"
		date, _   = models.ParseDate("2003-02-01")
		createdAt = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name          string
		updatedObject *models.Group
		wantErr       bool
	}{
		{
			"full",
			&models.Group{
				ID:        groupIDs[groupIdxWithTag],
				Name:      name,
				Duration:  &duration,
				Date:      &date,
				Rating:    &rating,
				StudioID:  &studioIDs[studioIdxWithGroup],
				Director:  director,
				Synopsis:  synopsis,
				URLs:      models.NewRelatedStrings([]string{url}),
				TagIDs:    models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithGroup]}),
				Aliases:   aliases,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"clear tag ids",
			&models.Group{
				ID:     groupIDs[groupIdxWithTag],
				Name:   name,
				TagIDs: models.NewRelatedIDs([]int{}),
			},
			false,
		},
		{
			"invalid studio id",
			&models.Group{
				ID:       groupIDs[groupIdxWithScene],
				Name:     name,
				StudioID: &invalidID,
			},
			true,
		},
		{
			"invalid tag id",
			&models.Group{
				ID:     groupIDs[groupIdxWithScene],
				Name:   name,
				TagIDs: models.NewRelatedIDs([]int{invalidID}),
			},
			true,
		},
	}

	qb := db.Group
	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			copy := *tt.updatedObject

			if err := qb.Update(ctx, tt.updatedObject); (err != nil) != tt.wantErr {
				t.Errorf("groupQueryBuilder.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			s, err := qb.Find(ctx, tt.updatedObject.ID)
			if err != nil {
				t.Errorf("groupQueryBuilder.Find() error = %v", err)
			}

			// load relationships
			if err := loadGroupRelationships(ctx, copy, s); err != nil {
				t.Errorf("loadGroupRelationships() error = %v", err)
				return
			}

			assert.Equal(copy, *s)
		})
	}
}

func clearGroupPartial() models.GroupPartial {
	// leave mandatory fields
	return models.GroupPartial{
		Aliases:  models.OptionalString{Set: true, Null: true},
		Synopsis: models.OptionalString{Set: true, Null: true},
		Director: models.OptionalString{Set: true, Null: true},
		Duration: models.OptionalInt{Set: true, Null: true},
		URLs:     &models.UpdateStrings{Mode: models.RelationshipUpdateModeSet},
		Date:     models.OptionalDate{Set: true, Null: true},
		Rating:   models.OptionalInt{Set: true, Null: true},
		StudioID: models.OptionalInt{Set: true, Null: true},
		TagIDs:   &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
	}
}

func Test_groupQueryBuilder_UpdatePartial(t *testing.T) {
	var (
		name      = "name"
		url       = "url"
		aliases   = "alias1, alias2"
		director  = "director"
		rating    = 60
		duration  = 34
		synopsis  = "synopsis"
		date, _   = models.ParseDate("2003-02-01")
		createdAt = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name    string
		id      int
		partial models.GroupPartial
		want    models.Group
		wantErr bool
	}{
		{
			"full",
			groupIDs[groupIdxWithScene],
			models.GroupPartial{
				Name:     models.NewOptionalString(name),
				Director: models.NewOptionalString(director),
				Synopsis: models.NewOptionalString(synopsis),
				Aliases:  models.NewOptionalString(aliases),
				URLs: &models.UpdateStrings{
					Values: []string{url},
					Mode:   models.RelationshipUpdateModeSet,
				},
				Date:      models.NewOptionalDate(date),
				Duration:  models.NewOptionalInt(duration),
				Rating:    models.NewOptionalInt(rating),
				StudioID:  models.NewOptionalInt(studioIDs[studioIdxWithGroup]),
				CreatedAt: models.NewOptionalTime(createdAt),
				UpdatedAt: models.NewOptionalTime(updatedAt),
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithGroup], tagIDs[tagIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeSet,
				},
			},
			models.Group{
				ID:        groupIDs[groupIdxWithScene],
				Name:      name,
				Director:  director,
				Synopsis:  synopsis,
				Aliases:   aliases,
				URLs:      models.NewRelatedStrings([]string{url}),
				Date:      &date,
				Duration:  &duration,
				Rating:    &rating,
				StudioID:  &studioIDs[studioIdxWithGroup],
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				TagIDs:    models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithGroup]}),
			},
			false,
		},
		{
			"clear all",
			groupIDs[groupIdxWithScene],
			clearGroupPartial(),
			models.Group{
				ID:     groupIDs[groupIdxWithScene],
				Name:   groupNames[groupIdxWithScene],
				TagIDs: models.NewRelatedIDs([]int{}),
			},
			false,
		},
		{
			"invalid id",
			invalidID,
			models.GroupPartial{},
			models.Group{},
			true,
		},
	}
	for _, tt := range tests {
		qb := db.Group

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			got, err := qb.UpdatePartial(ctx, tt.id, tt.partial)
			if (err != nil) != tt.wantErr {
				t.Errorf("groupQueryBuilder.UpdatePartial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// load relationships
			if err := loadGroupRelationships(ctx, tt.want, got); err != nil {
				t.Errorf("loadGroupRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, *got)

			s, err := qb.Find(ctx, tt.id)
			if err != nil {
				t.Errorf("groupQueryBuilder.Find() error = %v", err)
			}

			// load relationships
			if err := loadGroupRelationships(ctx, tt.want, s); err != nil {
				t.Errorf("loadGroupRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, *s)
		})
	}
}

func TestGroupFindByName(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		mqb := db.Group

		name := groupNames[groupIdxWithScene] // find a group by name

		group, err := mqb.FindByName(ctx, name, false)

		if err != nil {
			t.Errorf("Error finding groups: %s", err.Error())
		}

		assert.Equal(t, groupNames[groupIdxWithScene], group.Name)

		name = groupNames[groupIdxWithDupName] // find a group by name nocase

		group, err = mqb.FindByName(ctx, name, true)

		if err != nil {
			t.Errorf("Error finding groups: %s", err.Error())
		}
		// groupIdxWithDupName and groupIdxWithScene should have similar names ( only diff should be Name vs NaMe)
		//group.Name should match with groupIdxWithScene since its ID is before moveIdxWithDupName
		assert.Equal(t, groupNames[groupIdxWithScene], group.Name)
		//group.Name should match with groupIdxWithDupName if the check is not case sensitive
		assert.Equal(t, strings.ToLower(groupNames[groupIdxWithDupName]), strings.ToLower(group.Name))

		return nil
	})
}

func TestGroupFindByNames(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		var names []string

		mqb := db.Group

		names = append(names, groupNames[groupIdxWithScene]) // find groups by names

		groups, err := mqb.FindByNames(ctx, names, false)
		if err != nil {
			t.Errorf("Error finding groups: %s", err.Error())
		}
		assert.Len(t, groups, 1)
		assert.Equal(t, groupNames[groupIdxWithScene], groups[0].Name)

		groups, err = mqb.FindByNames(ctx, names, true) // find groups by names nocase
		if err != nil {
			t.Errorf("Error finding groups: %s", err.Error())
		}
		assert.Len(t, groups, 2) // groupIdxWithScene and groupIdxWithDupName
		assert.Equal(t, strings.ToLower(groupNames[groupIdxWithScene]), strings.ToLower(groups[0].Name))
		assert.Equal(t, strings.ToLower(groupNames[groupIdxWithScene]), strings.ToLower(groups[1].Name))

		return nil
	})
}

func groupsToIDs(i []*models.Group) []int {
	ret := make([]int, len(i))
	for i, v := range i {
		ret[i] = v.ID
	}

	return ret
}

func TestGroupQuery(t *testing.T) {
	var (
		frontImage = "front_image"
		backImage  = "back_image"
	)

	tests := []struct {
		name        string
		findFilter  *models.FindFilterType
		filter      *models.GroupFilterType
		includeIdxs []int
		excludeIdxs []int
		wantErr     bool
	}{
		{
			"is missing front image",
			nil,
			&models.GroupFilterType{
				IsMissing: &frontImage,
			},
			// just ensure that it doesn't error
			nil,
			nil,
			false,
		},
		{
			"is missing back image",
			nil,
			&models.GroupFilterType{
				IsMissing: &backImage,
			},
			// just ensure that it doesn't error
			nil,
			nil,
			false,
		},
	}

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			results, _, err := db.Group.Query(ctx, tt.filter, tt.findFilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupQueryBuilder.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ids := groupsToIDs(results)
			include := indexesToIDs(performerIDs, tt.includeIdxs)
			exclude := indexesToIDs(performerIDs, tt.excludeIdxs)

			for _, i := range include {
				assert.Contains(ids, i)
			}
			for _, e := range exclude {
				assert.NotContains(ids, e)
			}
		})
	}
}

func TestGroupQueryStudio(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		mqb := db.Group
		studioCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithGroup]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		groupFilter := models.GroupFilterType{
			Studios: &studioCriterion,
		}

		groups, _, err := mqb.Query(ctx, &groupFilter, nil)
		if err != nil {
			t.Errorf("Error querying group: %s", err.Error())
		}

		assert.Len(t, groups, 1)

		// ensure id is correct
		assert.Equal(t, groupIDs[groupIdxWithStudio], groups[0].ID)

		studioCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithGroup]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getGroupStringValue(groupIdxWithStudio, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		groups, _, err = mqb.Query(ctx, &groupFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying group: %s", err.Error())
		}
		assert.Len(t, groups, 0)

		return nil
	})
}

func TestGroupQueryURL(t *testing.T) {
	const sceneIdx = 1
	groupURL := getGroupStringValue(sceneIdx, urlField)

	urlCriterion := models.StringCriterionInput{
		Value:    groupURL,
		Modifier: models.CriterionModifierEquals,
	}

	filter := models.GroupFilterType{
		URL: &urlCriterion,
	}

	verifyFn := func(n *models.Group) {
		t.Helper()

		urls := n.URLs.List()
		var url string
		if len(urls) > 0 {
			url = urls[0]
		}

		verifyString(t, url, urlCriterion)
	}

	verifyGroupQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotEquals
	verifyGroupQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierMatchesRegex
	urlCriterion.Value = "group_.*1_URL"
	verifyGroupQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotMatchesRegex
	verifyGroupQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierIsNull
	urlCriterion.Value = ""
	verifyGroupQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotNull
	verifyGroupQuery(t, filter, verifyFn)
}

func TestGroupQueryURLExcludes(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		mqb := db.Group

		// create group with two URLs
		group := models.Group{
			Name: "TestGroupQueryURLExcludes",
			URLs: models.NewRelatedStrings([]string{
				"aaa",
				"bbb",
			}),
		}

		err := mqb.Create(ctx, &group)

		if err != nil {
			return fmt.Errorf("Error creating group: %w", err)
		}

		// query for groups that exclude the URL "aaa"
		urlCriterion := models.StringCriterionInput{
			Value:    "aaa",
			Modifier: models.CriterionModifierExcludes,
		}

		nameCriterion := models.StringCriterionInput{
			Value:    group.Name,
			Modifier: models.CriterionModifierEquals,
		}

		filter := models.GroupFilterType{
			URL:  &urlCriterion,
			Name: &nameCriterion,
		}

		groups := queryGroups(ctx, t, &filter, nil)
		assert.Len(t, groups, 0, "Expected no groups to be found")

		// query for groups that exclude the URL "ccc"
		urlCriterion.Value = "ccc"
		groups = queryGroups(ctx, t, &filter, nil)

		if assert.Len(t, groups, 1, "Expected one group to be found") {
			assert.Equal(t, group.Name, groups[0].Name)
		}

		return nil
	})
}

func verifyGroupQuery(t *testing.T, filter models.GroupFilterType, verifyFn func(s *models.Group)) {
	withTxn(func(ctx context.Context) error {
		t.Helper()
		sqb := db.Group

		groups := queryGroups(ctx, t, &filter, nil)

		for _, group := range groups {
			if err := group.LoadURLs(ctx, sqb); err != nil {
				t.Errorf("Error loading group relationships: %v", err)
			}
		}

		// assume it should find at least one
		assert.Greater(t, len(groups), 0)

		for _, m := range groups {
			verifyFn(m)
		}

		return nil
	})
}

func queryGroups(ctx context.Context, t *testing.T, groupFilter *models.GroupFilterType, findFilter *models.FindFilterType) []*models.Group {
	sqb := db.Group
	groups, _, err := sqb.Query(ctx, groupFilter, findFilter)
	if err != nil {
		t.Errorf("Error querying group: %s", err.Error())
	}

	return groups
}

func TestGroupQueryTags(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		tagCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdxWithGroup]),
				strconv.Itoa(tagIDs[tagIdx1WithGroup]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		groupFilter := models.GroupFilterType{
			Tags: &tagCriterion,
		}

		// ensure ids are correct
		groups := queryGroups(ctx, t, &groupFilter, nil)
		assert.Len(t, groups, 3)
		for _, group := range groups {
			assert.True(t, group.ID == groupIDs[groupIdxWithTag] || group.ID == groupIDs[groupIdxWithTwoTags] || group.ID == groupIDs[groupIdxWithThreeTags])
		}

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithGroup]),
				strconv.Itoa(tagIDs[tagIdx2WithGroup]),
			},
			Modifier: models.CriterionModifierIncludesAll,
		}

		groups = queryGroups(ctx, t, &groupFilter, nil)

		if assert.Len(t, groups, 2) {
			assert.Equal(t, sceneIDs[groupIdxWithTwoTags], groups[0].ID)
			assert.Equal(t, sceneIDs[groupIdxWithThreeTags], groups[1].ID)
		}

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithGroup]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getSceneStringValue(groupIdxWithTwoTags, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		groups = queryGroups(ctx, t, &groupFilter, &findFilter)
		assert.Len(t, groups, 0)

		return nil
	})
}

func TestGroupQueryTagCount(t *testing.T) {
	const tagCount = 1
	tagCountCriterion := models.IntCriterionInput{
		Value:    tagCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyGroupsTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyGroupsTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyGroupsTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyGroupsTagCount(t, tagCountCriterion)
}

func verifyGroupsTagCount(t *testing.T, tagCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Group
		groupFilter := models.GroupFilterType{
			TagCount: &tagCountCriterion,
		}

		groups := queryGroups(ctx, t, &groupFilter, nil)
		assert.Greater(t, len(groups), 0)

		for _, group := range groups {
			ids, err := sqb.GetTagIDs(ctx, group.ID)
			if err != nil {
				return err
			}
			verifyInt(t, len(ids), tagCountCriterion)
		}

		return nil
	})
}

func TestGroupQuerySorting(t *testing.T) {
	sort := "scenes_count"
	direction := models.SortDirectionEnumDesc
	findFilter := models.FindFilterType{
		Sort:      &sort,
		Direction: &direction,
	}

	withTxn(func(ctx context.Context) error {
		groups := queryGroups(ctx, t, nil, &findFilter)

		// scenes should be in same order as indexes
		firstGroup := groups[0]

		assert.Equal(t, groupIDs[groupIdxWithScene], firstGroup.ID)

		// sort in descending order
		direction = models.SortDirectionEnumAsc

		groups = queryGroups(ctx, t, nil, &findFilter)
		lastGroup := groups[len(groups)-1]

		assert.Equal(t, groupIDs[groupIdxWithScene], lastGroup.ID)

		return nil
	})
}

func TestGroupUpdateFrontImage(t *testing.T) {
	if err := withRollbackTxn(func(ctx context.Context) error {
		qb := db.Group

		// create group to test against
		const name = "TestGroupUpdateGroupImages"
		group := models.Group{
			Name: name,
		}
		err := qb.Create(ctx, &group)
		if err != nil {
			return fmt.Errorf("Error creating group: %s", err.Error())
		}

		return testUpdateImage(t, ctx, group.ID, qb.UpdateFrontImage, qb.GetFrontImage)
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestGroupUpdateBackImage(t *testing.T) {
	if err := withRollbackTxn(func(ctx context.Context) error {
		qb := db.Group

		// create group to test against
		const name = "TestGroupUpdateGroupImages"
		group := models.Group{
			Name: name,
		}
		err := qb.Create(ctx, &group)
		if err != nil {
			return fmt.Errorf("Error creating group: %s", err.Error())
		}

		return testUpdateImage(t, ctx, group.ID, qb.UpdateBackImage, qb.GetBackImage)
	}); err != nil {
		t.Error(err.Error())
	}
}

// TODO Update
// TODO Destroy - ensure image is destroyed
// TODO Find
// TODO Count
// TODO All
// TODO Query
