//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
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
	if expected.ContainingGroups.Loaded() {
		if err := actual.LoadContainingGroupIDs(ctx, db.Group); err != nil {
			return err
		}
	}
	if expected.SubGroups.Loaded() {
		if err := actual.LoadSubGroupIDs(ctx, db.Group); err != nil {
			return err
		}
	}

	return nil
}

func Test_GroupStore_Create(t *testing.T) {
	var (
		name                       = "name"
		url                        = "url"
		aliases                    = "alias1, alias2"
		director                   = "director"
		rating                     = 60
		duration                   = 34
		synopsis                   = "synopsis"
		date, _                    = models.ParseDate("2003-02-01")
		containingGroupDescription = "containingGroupDescription"
		subGroupDescription        = "subGroupDescription"
		createdAt                  = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt                  = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name      string
		newObject models.Group
		wantErr   bool
	}{
		{
			"full",
			models.Group{
				Name:     name,
				Duration: &duration,
				Date:     &date,
				Rating:   &rating,
				StudioID: &studioIDs[studioIdxWithGroup],
				Director: director,
				Synopsis: synopsis,
				URLs:     models.NewRelatedStrings([]string{url}),
				TagIDs:   models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithGroup]}),
				ContainingGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{
					{GroupID: groupIDs[groupIdxWithScene], Description: containingGroupDescription},
				}),
				SubGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{
					{GroupID: groupIDs[groupIdxWithStudio], Description: subGroupDescription},
				}),
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
		{
			"invalid containing group id",
			models.Group{
				Name:             name,
				ContainingGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{{GroupID: invalidID}}),
			},
			true,
		},
		{
			"invalid sub group id",
			models.Group{
				Name:      name,
				SubGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{{GroupID: invalidID}}),
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
		name                       = "name"
		url                        = "url"
		aliases                    = "alias1, alias2"
		director                   = "director"
		rating                     = 60
		duration                   = 34
		synopsis                   = "synopsis"
		date, _                    = models.ParseDate("2003-02-01")
		containingGroupDescription = "containingGroupDescription"
		subGroupDescription        = "subGroupDescription"
		createdAt                  = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt                  = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name          string
		updatedObject models.Group
		wantErr       bool
	}{
		{
			"full",
			models.Group{
				ID:       groupIDs[groupIdxWithTag],
				Name:     name,
				Duration: &duration,
				Date:     &date,
				Rating:   &rating,
				StudioID: &studioIDs[studioIdxWithGroup],
				Director: director,
				Synopsis: synopsis,
				URLs:     models.NewRelatedStrings([]string{url}),
				TagIDs:   models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithGroup]}),
				ContainingGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{
					{GroupID: groupIDs[groupIdxWithScene], Description: containingGroupDescription},
				}),
				SubGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{
					{GroupID: groupIDs[groupIdxWithStudio], Description: subGroupDescription},
				}),
				Aliases:   aliases,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"clear tag ids",
			models.Group{
				ID:     groupIDs[groupIdxWithTag],
				Name:   name,
				TagIDs: models.NewRelatedIDs([]int{}),
			},
			false,
		},
		{
			"clear containing ids",
			models.Group{
				ID:               groupIDs[groupIdxWithParent],
				Name:             name,
				ContainingGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{}),
			},
			false,
		},
		{
			"clear sub ids",
			models.Group{
				ID:        groupIDs[groupIdxWithChild],
				Name:      name,
				SubGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{}),
			},
			false,
		},
		{
			"invalid studio id",
			models.Group{
				ID:       groupIDs[groupIdxWithScene],
				Name:     name,
				StudioID: &invalidID,
			},
			true,
		},
		{
			"invalid tag id",
			models.Group{
				ID:     groupIDs[groupIdxWithScene],
				Name:   name,
				TagIDs: models.NewRelatedIDs([]int{invalidID}),
			},
			true,
		},
		{
			"invalid containing group id",
			models.Group{
				ID:               groupIDs[groupIdxWithScene],
				Name:             name,
				ContainingGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{{GroupID: invalidID}}),
			},
			true,
		},
		{
			"invalid sub group id",
			models.Group{
				ID:        groupIDs[groupIdxWithScene],
				Name:      name,
				SubGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{{GroupID: invalidID}}),
			},
			true,
		},
	}

	qb := db.Group
	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			actual := tt.updatedObject
			expected := tt.updatedObject

			if err := qb.Update(ctx, &actual); (err != nil) != tt.wantErr {
				t.Errorf("groupQueryBuilder.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			s, err := qb.Find(ctx, actual.ID)
			if err != nil {
				t.Errorf("groupQueryBuilder.Find() error = %v", err)
			}

			// load relationships
			if err := loadGroupRelationships(ctx, expected, s); err != nil {
				t.Errorf("loadGroupRelationships() error = %v", err)
				return
			}

			assert.Equal(expected, *s)
		})
	}
}

var clearGroupPartial = models.GroupPartial{
	// leave mandatory fields
	Aliases:          models.OptionalString{Set: true, Null: true},
	Synopsis:         models.OptionalString{Set: true, Null: true},
	Director:         models.OptionalString{Set: true, Null: true},
	Duration:         models.OptionalInt{Set: true, Null: true},
	URLs:             &models.UpdateStrings{Mode: models.RelationshipUpdateModeSet},
	Date:             models.OptionalDate{Set: true, Null: true},
	Rating:           models.OptionalInt{Set: true, Null: true},
	StudioID:         models.OptionalInt{Set: true, Null: true},
	TagIDs:           &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
	ContainingGroups: &models.UpdateGroupDescriptions{Mode: models.RelationshipUpdateModeSet},
	SubGroups:        &models.UpdateGroupDescriptions{Mode: models.RelationshipUpdateModeSet},
}

func emptyGroup(idx int) models.Group {
	return models.Group{
		ID:               groupIDs[idx],
		Name:             groupNames[idx],
		TagIDs:           models.NewRelatedIDs([]int{}),
		ContainingGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{}),
		SubGroups:        models.NewRelatedGroupDescriptions([]models.GroupIDDescription{}),
	}
}

func Test_groupQueryBuilder_UpdatePartial(t *testing.T) {
	var (
		name                       = "name"
		url                        = "url"
		aliases                    = "alias1, alias2"
		director                   = "director"
		rating                     = 60
		duration                   = 34
		synopsis                   = "synopsis"
		date, _                    = models.ParseDate("2003-02-01")
		containingGroupDescription = "containingGroupDescription"
		subGroupDescription        = "subGroupDescription"
		createdAt                  = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt                  = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
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
				ContainingGroups: &models.UpdateGroupDescriptions{
					Groups: []models.GroupIDDescription{
						{GroupID: groupIDs[groupIdxWithStudio], Description: containingGroupDescription},
						{GroupID: groupIDs[groupIdxWithThreeTags], Description: containingGroupDescription},
					},
					Mode: models.RelationshipUpdateModeSet,
				},
				SubGroups: &models.UpdateGroupDescriptions{
					Groups: []models.GroupIDDescription{
						{GroupID: groupIDs[groupIdxWithTag], Description: subGroupDescription},
						{GroupID: groupIDs[groupIdxWithDupName], Description: subGroupDescription},
					},
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
				ContainingGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{
					{GroupID: groupIDs[groupIdxWithStudio], Description: containingGroupDescription},
					{GroupID: groupIDs[groupIdxWithThreeTags], Description: containingGroupDescription},
				}),
				SubGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{
					{GroupID: groupIDs[groupIdxWithTag], Description: subGroupDescription},
					{GroupID: groupIDs[groupIdxWithDupName], Description: subGroupDescription},
				}),
			},
			false,
		},
		{
			"clear all",
			groupIDs[groupIdxWithScene],
			clearGroupPartial,
			emptyGroup(groupIdxWithScene),
			false,
		},
		{
			"clear tag ids",
			groupIDs[groupIdxWithTag],
			clearGroupPartial,
			emptyGroup(groupIdxWithTag),
			false,
		},
		{
			"clear group relationships",
			groupIDs[groupIdxWithParentAndChild],
			clearGroupPartial,
			emptyGroup(groupIdxWithParentAndChild),
			false,
		},
		{
			"add containing group",
			groupIDs[groupIdxWithParent],
			models.GroupPartial{
				ContainingGroups: &models.UpdateGroupDescriptions{
					Groups: []models.GroupIDDescription{
						{GroupID: groupIDs[groupIdxWithScene], Description: containingGroupDescription},
					},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Group{
				ID:   groupIDs[groupIdxWithParent],
				Name: groupNames[groupIdxWithParent],
				ContainingGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{
					{GroupID: groupIDs[groupIdxWithChild]},
					{GroupID: groupIDs[groupIdxWithScene], Description: containingGroupDescription},
				}),
			},
			false,
		},
		{
			"add sub group",
			groupIDs[groupIdxWithChild],
			models.GroupPartial{
				SubGroups: &models.UpdateGroupDescriptions{
					Groups: []models.GroupIDDescription{
						{GroupID: groupIDs[groupIdxWithScene], Description: subGroupDescription},
					},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Group{
				ID:   groupIDs[groupIdxWithChild],
				Name: groupNames[groupIdxWithChild],
				SubGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{
					{GroupID: groupIDs[groupIdxWithParent]},
					{GroupID: groupIDs[groupIdxWithScene], Description: subGroupDescription},
				}),
			},
			false,
		},
		{
			"remove containing group",
			groupIDs[groupIdxWithParent],
			models.GroupPartial{
				ContainingGroups: &models.UpdateGroupDescriptions{
					Groups: []models.GroupIDDescription{
						{GroupID: groupIDs[groupIdxWithChild]},
					},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Group{
				ID:               groupIDs[groupIdxWithParent],
				Name:             groupNames[groupIdxWithParent],
				ContainingGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{}),
			},
			false,
		},
		{
			"remove sub group",
			groupIDs[groupIdxWithChild],
			models.GroupPartial{
				SubGroups: &models.UpdateGroupDescriptions{
					Groups: []models.GroupIDDescription{
						{GroupID: groupIDs[groupIdxWithParent]},
					},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Group{
				ID:        groupIDs[groupIdxWithChild],
				Name:      groupNames[groupIdxWithChild],
				SubGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{}),
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

		assert.Equal(t, groupIDs[groupIdxWithParentAndScene], lastGroup.ID)

		return nil
	})
}

func TestGroupQuerySortOrderIndex(t *testing.T) {
	sort := "sub_group_order"
	direction := models.SortDirectionEnumDesc
	findFilter := models.FindFilterType{
		Sort:      &sort,
		Direction: &direction,
	}

	groupFilter := models.GroupFilterType{
		ContainingGroups: &models.HierarchicalMultiCriterionInput{
			Value:    intslice.IntSliceToStringSlice([]int{groupIdxWithChild}),
			Modifier: models.CriterionModifierIncludes,
		},
	}

	withTxn(func(ctx context.Context) error {
		// just ensure there are no errors
		_, _, err := db.Group.Query(ctx, &groupFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying group: %s", err.Error())
		}

		_, _, err = db.Group.Query(ctx, nil, &findFilter)
		if err != nil {
			t.Errorf("Error querying group: %s", err.Error())
		}

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

func TestGroupQueryContainingGroups(t *testing.T) {
	const nameField = "Name"

	type criterion struct {
		valueIdxs []int
		modifier  models.CriterionModifier
		depth     int
	}

	tests := []struct {
		name        string
		c           criterion
		q           string
		includeIdxs []int
	}{
		{
			"includes",
			criterion{
				[]int{groupIdxWithChild},
				models.CriterionModifierIncludes,
				0,
			},
			"",
			[]int{groupIdxWithParent},
		},
		{
			"excludes",
			criterion{
				[]int{groupIdxWithChild},
				models.CriterionModifierExcludes,
				0,
			},
			getGroupStringValue(groupIdxWithParent, nameField),
			nil,
		},
		{
			"includes (all levels)",
			criterion{
				[]int{groupIdxWithGrandChild},
				models.CriterionModifierIncludes,
				-1,
			},
			"",
			[]int{groupIdxWithParentAndChild, groupIdxWithGrandParent},
		},
		{
			"includes (1 level)",
			criterion{
				[]int{groupIdxWithGrandChild},
				models.CriterionModifierIncludes,
				1,
			},
			"",
			[]int{groupIdxWithParentAndChild, groupIdxWithGrandParent},
		},
		{
			"is null",
			criterion{
				nil,
				models.CriterionModifierIsNull,
				0,
			},
			getGroupStringValue(groupIdxWithParent, nameField),
			nil,
		},
		{
			"not null",
			criterion{
				nil,
				models.CriterionModifierNotNull,
				0,
			},
			"",
			[]int{groupIdxWithParentAndChild, groupIdxWithParent, groupIdxWithGrandParent, groupIdxWithParentAndScene},
		},
	}

	qb := db.Group

	for _, tt := range tests {
		valueIDs := indexesToIDs(groupIDs, tt.c.valueIdxs)
		expectedIDs := indexesToIDs(groupIDs, tt.includeIdxs)

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			groupFilter := &models.GroupFilterType{
				ContainingGroups: &models.HierarchicalMultiCriterionInput{
					Value:    intslice.IntSliceToStringSlice(valueIDs),
					Modifier: tt.c.modifier,
				},
			}

			if tt.c.depth != 0 {
				groupFilter.ContainingGroups.Depth = &tt.c.depth
			}

			findFilter := models.FindFilterType{}
			if tt.q != "" {
				findFilter.Q = &tt.q
			}

			groups, _, err := qb.Query(ctx, groupFilter, &findFilter)
			if err != nil {
				t.Errorf("GroupStore.Query() error = %v", err)
				return
			}

			// get ids of groups
			groupIDs := sliceutil.Map(groups, func(g *models.Group) int { return g.ID })
			assert.ElementsMatch(t, expectedIDs, groupIDs)
		})
	}
}

func TestGroupQuerySubGroups(t *testing.T) {
	const nameField = "Name"

	type criterion struct {
		valueIdxs []int
		modifier  models.CriterionModifier
		depth     int
	}

	tests := []struct {
		name         string
		c            criterion
		q            string
		expectedIdxs []int
	}{
		{
			"includes",
			criterion{
				[]int{groupIdxWithParent},
				models.CriterionModifierIncludes,
				0,
			},
			"",
			[]int{groupIdxWithChild},
		},
		{
			"excludes",
			criterion{
				[]int{groupIdxWithParent},
				models.CriterionModifierExcludes,
				0,
			},
			getGroupStringValue(groupIdxWithChild, nameField),
			nil,
		},
		{
			"includes (all levels)",
			criterion{
				[]int{groupIdxWithGrandParent},
				models.CriterionModifierIncludes,
				-1,
			},
			"",
			[]int{groupIdxWithGrandChild, groupIdxWithParentAndChild},
		},
		{
			"includes (1 level)",
			criterion{
				[]int{groupIdxWithGrandParent},
				models.CriterionModifierIncludes,
				1,
			},
			"",
			[]int{groupIdxWithGrandChild, groupIdxWithParentAndChild},
		},
		{
			"is null",
			criterion{
				nil,
				models.CriterionModifierIsNull,
				0,
			},
			getGroupStringValue(groupIdxWithChild, nameField),
			nil,
		},
		{
			"not null",
			criterion{
				nil,
				models.CriterionModifierNotNull,
				0,
			},
			"",
			[]int{groupIdxWithGrandChild, groupIdxWithChild, groupIdxWithParentAndChild, groupIdxWithChildWithScene},
		},
	}

	qb := db.Group

	for _, tt := range tests {
		valueIDs := indexesToIDs(groupIDs, tt.c.valueIdxs)
		expectedIDs := indexesToIDs(groupIDs, tt.expectedIdxs)

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			groupFilter := &models.GroupFilterType{
				SubGroups: &models.HierarchicalMultiCriterionInput{
					Value:    intslice.IntSliceToStringSlice(valueIDs),
					Modifier: tt.c.modifier,
				},
			}

			if tt.c.depth != 0 {
				groupFilter.SubGroups.Depth = &tt.c.depth
			}

			findFilter := models.FindFilterType{}
			if tt.q != "" {
				findFilter.Q = &tt.q
			}

			groups, _, err := qb.Query(ctx, groupFilter, &findFilter)
			if err != nil {
				t.Errorf("GroupStore.Query() error = %v", err)
				return
			}

			// get ids of groups
			groupIDs := sliceutil.Map(groups, func(g *models.Group) int { return g.ID })
			assert.ElementsMatch(t, expectedIDs, groupIDs)
		})
	}
}

func TestGroupQueryContainingGroupCount(t *testing.T) {
	const nameField = "Name"

	tests := []struct {
		name         string
		value        int
		modifier     models.CriterionModifier
		q            string
		expectedIdxs []int
	}{
		{
			"equals",
			1,
			models.CriterionModifierEquals,
			"",
			[]int{groupIdxWithParent, groupIdxWithGrandParent, groupIdxWithParentAndChild, groupIdxWithParentAndScene},
		},
		{
			"not equals",
			1,
			models.CriterionModifierNotEquals,
			getGroupStringValue(groupIdxWithParent, nameField),
			nil,
		},
		{
			"less than",
			1,
			models.CriterionModifierLessThan,
			getGroupStringValue(groupIdxWithParent, nameField),
			nil,
		},
		{
			"greater than",
			0,
			models.CriterionModifierGreaterThan,
			"",
			[]int{groupIdxWithParent, groupIdxWithGrandParent, groupIdxWithParentAndChild, groupIdxWithParentAndScene},
		},
	}

	qb := db.Group

	for _, tt := range tests {
		expectedIDs := indexesToIDs(groupIDs, tt.expectedIdxs)

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			groupFilter := &models.GroupFilterType{
				ContainingGroupCount: &models.IntCriterionInput{
					Value:    tt.value,
					Modifier: tt.modifier,
				},
			}

			findFilter := models.FindFilterType{}
			if tt.q != "" {
				findFilter.Q = &tt.q
			}

			groups, _, err := qb.Query(ctx, groupFilter, &findFilter)
			if err != nil {
				t.Errorf("GroupStore.Query() error = %v", err)
				return
			}

			// get ids of groups
			groupIDs := sliceutil.Map(groups, func(g *models.Group) int { return g.ID })
			assert.ElementsMatch(t, expectedIDs, groupIDs)
		})
	}
}

func TestGroupQuerySubGroupCount(t *testing.T) {
	const nameField = "Name"

	tests := []struct {
		name         string
		value        int
		modifier     models.CriterionModifier
		q            string
		expectedIdxs []int
	}{
		{
			"equals",
			1,
			models.CriterionModifierEquals,
			"",
			[]int{groupIdxWithChild, groupIdxWithGrandChild, groupIdxWithParentAndChild, groupIdxWithChildWithScene},
		},
		{
			"not equals",
			1,
			models.CriterionModifierNotEquals,
			getGroupStringValue(groupIdxWithChild, nameField),
			nil,
		},
		{
			"less than",
			1,
			models.CriterionModifierLessThan,
			getGroupStringValue(groupIdxWithChild, nameField),
			nil,
		},
		{
			"greater than",
			0,
			models.CriterionModifierGreaterThan,
			"",
			[]int{groupIdxWithChild, groupIdxWithGrandChild, groupIdxWithParentAndChild, groupIdxWithChildWithScene},
		},
	}

	qb := db.Group

	for _, tt := range tests {
		expectedIDs := indexesToIDs(groupIDs, tt.expectedIdxs)

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			groupFilter := &models.GroupFilterType{
				SubGroupCount: &models.IntCriterionInput{
					Value:    tt.value,
					Modifier: tt.modifier,
				},
			}

			findFilter := models.FindFilterType{}
			if tt.q != "" {
				findFilter.Q = &tt.q
			}

			groups, _, err := qb.Query(ctx, groupFilter, &findFilter)
			if err != nil {
				t.Errorf("GroupStore.Query() error = %v", err)
				return
			}

			// get ids of groups
			groupIDs := sliceutil.Map(groups, func(g *models.Group) int { return g.ID })
			assert.ElementsMatch(t, expectedIDs, groupIDs)
		})
	}
}

func TestGroupFindInAncestors(t *testing.T) {
	tests := []struct {
		name         string
		ancestorIdxs []int
		idxs         []int
		expectedIdxs []int
	}{
		{
			"basic",
			[]int{groupIdxWithGrandParent},
			[]int{groupIdxWithGrandChild},
			[]int{groupIdxWithGrandChild},
		},
		{
			"same",
			[]int{groupIdxWithScene},
			[]int{groupIdxWithScene},
			[]int{groupIdxWithScene},
		},
		{
			"no matches",
			[]int{groupIdxWithGrandParent},
			[]int{groupIdxWithScene},
			nil,
		},
	}

	qb := db.Group

	for _, tt := range tests {
		ancestorIDs := indexesToIDs(groupIDs, tt.ancestorIdxs)
		ids := indexesToIDs(groupIDs, tt.idxs)
		expectedIDs := indexesToIDs(groupIDs, tt.expectedIdxs)

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			found, err := qb.FindInAncestors(ctx, ancestorIDs, ids)
			if err != nil {
				t.Errorf("GroupStore.FindInAncestors() error = %v", err)
				return
			}

			// get ids of groups
			assert.ElementsMatch(t, found, expectedIDs)
		})
	}
}

func TestGroupReorderSubGroups(t *testing.T) {
	tests := []struct {
		name        string
		subGroupLen int
		idxsToMove  []int
		insertLoc   int
		insertAfter bool
		// order of elements, using original indexes
		expectedIdxs []int
	}{
		{
			"move single back before",
			5,
			[]int{2},
			1,
			false,
			[]int{0, 2, 1, 3, 4},
		},
		{
			"move single forward before",
			5,
			[]int{2},
			4,
			false,
			[]int{0, 1, 3, 2, 4},
		},
		{
			"move multiple back before",
			5,
			[]int{3, 2, 4},
			0,
			false,
			[]int{3, 2, 4, 0, 1},
		},
		{
			"move multiple forward before",
			5,
			[]int{2, 1, 0},
			4,
			false,
			[]int{3, 2, 1, 0, 4},
		},
		{
			"move single back after",
			5,
			[]int{2},
			0,
			true,
			[]int{0, 2, 1, 3, 4},
		},
		{
			"move single forward after",
			5,
			[]int{2},
			4,
			true,
			[]int{0, 1, 3, 4, 2},
		},
		{
			"move multiple back after",
			5,
			[]int{3, 2, 4},
			0,
			false,
			[]int{0, 3, 2, 4, 1},
		},
		{
			"move multiple forward after",
			5,
			[]int{2, 1, 0},
			4,
			false,
			[]int{3, 4, 2, 1, 0},
		},
	}

	qb := db.Group

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			// create the group
			group := models.Group{
				Name: "TestGroupReorderSubGroups",
			}

			if err := qb.Create(ctx, &group); err != nil {
				t.Errorf("GroupStore.Create() error = %v", err)
				return
			}

			// and sub-groups
			idxToId := make([]int, tt.subGroupLen)

			for i := 0; i < tt.subGroupLen; i++ {
				subGroup := models.Group{
					Name: fmt.Sprintf("SubGroup %d", i),
					ContainingGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{
						{GroupID: group.ID},
					}),
				}

				if err := qb.Create(ctx, &subGroup); err != nil {
					t.Errorf("GroupStore.Create() error = %v", err)
					return
				}

				idxToId[i] = subGroup.ID
			}

			// reorder
			idsToMove := indexesToIDs(idxToId, tt.idxsToMove)
			insertID := idxToId[tt.insertLoc]
			if err := qb.ReorderSubGroups(ctx, group.ID, idsToMove, insertID, tt.insertAfter); err != nil {
				t.Errorf("GroupStore.ReorderSubGroups() error = %v", err)
				return
			}

			// validate the new order
			gd, err := qb.GetSubGroupDescriptions(ctx, group.ID)
			if err != nil {
				t.Errorf("GroupStore.GetSubGroupDescriptions() error = %v", err)
				return
			}

			// get ids of groups
			newIDs := sliceutil.Map(gd, func(gd models.GroupIDDescription) int { return gd.GroupID })
			newIdxs := sliceutil.Map(newIDs, func(id int) int { return slices.Index(idxToId, id) })

			assert.ElementsMatch(t, tt.expectedIdxs, newIdxs)
		})
	}
}

func TestGroupAddSubGroups(t *testing.T) {
	tests := []struct {
		name                string
		existingSubGroupLen int
		insertGroupsLen     int
		insertLoc           int
		// order of elements, using original indexes
		expectedIdxs []int
	}{
		{
			"append single",
			4,
			1,
			999,
			[]int{0, 1, 2, 3, 4},
		},
		{
			"insert single middle",
			4,
			1,
			2,
			[]int{0, 1, 4, 2, 3},
		},
		{
			"insert single start",
			4,
			1,
			0,
			[]int{4, 0, 1, 2, 3},
		},
		{
			"append multiple",
			4,
			2,
			999,
			[]int{0, 1, 2, 3, 4, 5},
		},
		{
			"insert multiple middle",
			4,
			2,
			2,
			[]int{0, 1, 4, 5, 2, 3},
		},
		{
			"insert multiple start",
			4,
			2,
			0,
			[]int{4, 5, 0, 1, 2, 3},
		},
	}

	qb := db.Group

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			// create the group
			group := models.Group{
				Name: "TestGroupReorderSubGroups",
			}

			if err := qb.Create(ctx, &group); err != nil {
				t.Errorf("GroupStore.Create() error = %v", err)
				return
			}

			// and sub-groups
			idxToId := make([]int, tt.existingSubGroupLen+tt.insertGroupsLen)

			for i := 0; i < tt.existingSubGroupLen; i++ {
				subGroup := models.Group{
					Name: fmt.Sprintf("Existing SubGroup %d", i),
					ContainingGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{
						{GroupID: group.ID},
					}),
				}

				if err := qb.Create(ctx, &subGroup); err != nil {
					t.Errorf("GroupStore.Create() error = %v", err)
					return
				}

				idxToId[i] = subGroup.ID
			}

			// and sub-groups to insert
			for i := 0; i < tt.insertGroupsLen; i++ {
				subGroup := models.Group{
					Name: fmt.Sprintf("Inserted SubGroup %d", i),
				}

				if err := qb.Create(ctx, &subGroup); err != nil {
					t.Errorf("GroupStore.Create() error = %v", err)
					return
				}

				idxToId[i+tt.existingSubGroupLen] = subGroup.ID
			}

			// convert ids to description
			idDescriptions := make([]models.GroupIDDescription, tt.insertGroupsLen)
			for i, id := range idxToId[tt.existingSubGroupLen:] {
				idDescriptions[i] = models.GroupIDDescription{GroupID: id}
			}

			// add
			if err := qb.AddSubGroups(ctx, group.ID, idDescriptions, &tt.insertLoc); err != nil {
				t.Errorf("GroupStore.AddSubGroups() error = %v", err)
				return
			}

			// validate the new order
			gd, err := qb.GetSubGroupDescriptions(ctx, group.ID)
			if err != nil {
				t.Errorf("GroupStore.GetSubGroupDescriptions() error = %v", err)
				return
			}

			// get ids of groups
			newIDs := sliceutil.Map(gd, func(gd models.GroupIDDescription) int { return gd.GroupID })
			newIdxs := sliceutil.Map(newIDs, func(id int) int { return slices.Index(idxToId, id) })

			assert.ElementsMatch(t, tt.expectedIdxs, newIdxs)
		})
	}
}

func TestGroupRemoveSubGroups(t *testing.T) {
	tests := []struct {
		name        string
		subGroupLen int
		removeIdxs  []int
		// order of elements, using original indexes
		expectedIdxs []int
	}{
		{
			"remove last",
			4,
			[]int{3},
			[]int{0, 1, 2},
		},
		{
			"remove first",
			4,
			[]int{0},
			[]int{1, 2, 3},
		},
		{
			"remove middle",
			4,
			[]int{2},
			[]int{0, 1, 3},
		},
		{
			"remove multiple",
			4,
			[]int{1, 3},
			[]int{0, 2},
		},
		{
			"remove all",
			4,
			[]int{0, 1, 2, 3},
			[]int{},
		},
	}

	qb := db.Group

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			// create the group
			group := models.Group{
				Name: "TestGroupReorderSubGroups",
			}

			if err := qb.Create(ctx, &group); err != nil {
				t.Errorf("GroupStore.Create() error = %v", err)
				return
			}

			// and sub-groups
			idxToId := make([]int, tt.subGroupLen)

			for i := 0; i < tt.subGroupLen; i++ {
				subGroup := models.Group{
					Name: fmt.Sprintf("Existing SubGroup %d", i),
					ContainingGroups: models.NewRelatedGroupDescriptions([]models.GroupIDDescription{
						{GroupID: group.ID},
					}),
				}

				if err := qb.Create(ctx, &subGroup); err != nil {
					t.Errorf("GroupStore.Create() error = %v", err)
					return
				}

				idxToId[i] = subGroup.ID
			}

			idsToRemove := indexesToIDs(idxToId, tt.removeIdxs)
			if err := qb.RemoveSubGroups(ctx, group.ID, idsToRemove); err != nil {
				t.Errorf("GroupStore.RemoveSubGroups() error = %v", err)
				return
			}

			// validate the new order
			gd, err := qb.GetSubGroupDescriptions(ctx, group.ID)
			if err != nil {
				t.Errorf("GroupStore.GetSubGroupDescriptions() error = %v", err)
				return
			}

			// get ids of groups
			newIDs := sliceutil.Map(gd, func(gd models.GroupIDDescription) int { return gd.GroupID })
			newIdxs := sliceutil.Map(newIDs, func(id int) int { return slices.Index(idxToId, id) })

			assert.ElementsMatch(t, tt.expectedIdxs, newIdxs)
		})
	}
}

func TestGroupFindSubGroupIDs(t *testing.T) {
	tests := []struct {
		name               string
		containingGroupIdx int
		subIdxs            []int
		expectedIdxs       []int
	}{
		{
			"overlap",
			groupIdxWithGrandChild,
			[]int{groupIdxWithParentAndChild, groupIdxWithGrandParent},
			[]int{groupIdxWithParentAndChild},
		},
		{
			"non-overlap",
			groupIdxWithGrandChild,
			[]int{groupIdxWithGrandParent},
			[]int{},
		},
		{
			"none",
			groupIdxWithScene,
			[]int{groupIdxWithDupName},
			[]int{},
		},
		{
			"invalid",
			invalidID,
			[]int{invalidID},
			[]int{},
		},
	}

	qb := db.Group

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			subIDs := indexesToIDs(groupIDs, tt.subIdxs)

			id := indexToID(groupIDs, tt.containingGroupIdx)

			found, err := qb.FindSubGroupIDs(ctx, id, subIDs)
			if err != nil {
				t.Errorf("GroupStore.FindSubGroupIDs() error = %v", err)
				return
			}

			// get ids of groups
			foundIdxs := sliceutil.Map(found, func(id int) int { return slices.Index(groupIDs, id) })

			assert.ElementsMatch(t, tt.expectedIdxs, foundIdxs)
		})
	}
}

// TODO Update
// TODO Destroy - ensure image is destroyed
// TODO Find
// TODO Count
// TODO All
// TODO Query
