//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestStudioFindByName(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Studio

		name := studioNames[studioIdxWithScene] // find a studio by name

		studio, err := sqb.FindByName(ctx, name, false)

		if err != nil {
			t.Errorf("Error finding studios: %s", err.Error())
		}

		assert.Equal(t, studioNames[studioIdxWithScene], studio.Name)

		name = studioNames[studioIdxWithDupName] // find a studio by name nocase

		studio, err = sqb.FindByName(ctx, name, true)

		if err != nil {
			t.Errorf("Error finding studios: %s", err.Error())
		}
		// studioIdxWithDupName and studioIdxWithScene should have similar names ( only diff should be Name vs NaMe)
		//studio.Name should match with studioIdxWithScene since its ID is before studioIdxWithDupName
		assert.Equal(t, studioNames[studioIdxWithScene], studio.Name)
		//studio.Name should match with studioIdxWithDupName if the check is not case sensitive
		assert.Equal(t, strings.ToLower(studioNames[studioIdxWithDupName]), strings.ToLower(studio.Name))

		return nil
	})
}

func loadStudioRelationships(ctx context.Context, expected models.Studio, actual *models.Studio) error {
	if expected.Aliases.Loaded() {
		if err := actual.LoadAliases(ctx, db.Studio); err != nil {
			return err
		}
	}
	if expected.URLs.Loaded() {
		if err := actual.LoadURLs(ctx, db.Studio); err != nil {
			return err
		}
	}
	if expected.TagIDs.Loaded() {
		if err := actual.LoadTagIDs(ctx, db.Studio); err != nil {
			return err
		}
	}
	if expected.StashIDs.Loaded() {
		if err := actual.LoadStashIDs(ctx, db.Studio); err != nil {
			return err
		}
	}

	return nil
}

func Test_StudioStore_Create(t *testing.T) {
	var (
		name          = "name"
		details       = "details"
		url           = "url"
		rating        = 3
		aliases       = []string{"alias1", "alias2"}
		ignoreAutoTag = true
		favorite      = true
		endpoint1     = "endpoint1"
		endpoint2     = "endpoint2"
		stashID1      = "stashid1"
		stashID2      = "stashid2"
		createdAt     = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt     = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name      string
		newObject models.CreateStudioInput
		wantErr   bool
	}{
		{
			"full",
			models.CreateStudioInput{
				Studio: &models.Studio{
					Name:          name,
					URLs:          models.NewRelatedStrings([]string{url}),
					Favorite:      favorite,
					Rating:        &rating,
					Details:       details,
					IgnoreAutoTag: ignoreAutoTag,
					TagIDs:        models.NewRelatedIDs([]int{tagIDs[tagIdx1WithStudio], tagIDs[tagIdx1WithDupName]}),
					Aliases:       models.NewRelatedStrings(aliases),
					StashIDs: models.NewRelatedStashIDs([]models.StashID{
						{
							StashID:   stashID1,
							Endpoint:  endpoint1,
							UpdatedAt: epochTime,
						},
						{
							StashID:   stashID2,
							Endpoint:  endpoint2,
							UpdatedAt: epochTime,
						},
					}),
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				CustomFields: testCustomFields,
			},
			false,
		},
		{
			"invalid tag id",
			models.CreateStudioInput{
				Studio: &models.Studio{
					Name:   name,
					TagIDs: models.NewRelatedIDs([]int{invalidID}),
				},
			},
			true,
		},
	}

	qb := db.Studio

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			p := tt.newObject
			if err := qb.Create(ctx, &p); (err != nil) != tt.wantErr {
				t.Errorf("StudioStore.Create() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				assert.Zero(p.ID)
				return
			}

			assert.NotZero(p.ID)

			copy := *tt.newObject.Studio
			copy.ID = p.ID

			// load relationships
			if err := loadStudioRelationships(ctx, copy, p.Studio); err != nil {
				t.Errorf("loadStudioRelationships() error = %v", err)
				return
			}

			assert.Equal(copy, *p.Studio)

			// ensure can find the Studio
			found, err := qb.Find(ctx, p.ID)
			if err != nil {
				t.Errorf("StudioStore.Find() error = %v", err)
			}

			if !assert.NotNil(found) {
				return
			}

			// load relationships
			if err := loadStudioRelationships(ctx, copy, found); err != nil {
				t.Errorf("loadStudioRelationships() error = %v", err)
				return
			}
			assert.Equal(copy, *found)

			// ensure custom fields are set
			cf, err := qb.GetCustomFields(ctx, p.ID)
			if err != nil {
				t.Errorf("StudioStore.GetCustomFields() error = %v", err)
				return
			}

			assert.Equal(tt.newObject.CustomFields, cf)

			return
		})
	}
}

func Test_StudioStore_Update(t *testing.T) {
	var (
		name          = "name"
		details       = "details"
		url           = "url"
		rating        = 3
		aliases       = []string{"aliasX", "aliasY"}
		ignoreAutoTag = true
		favorite      = true
		endpoint1     = "endpoint1"
		endpoint2     = "endpoint2"
		stashID1      = "stashid1"
		stashID2      = "stashid2"
		createdAt     = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt     = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name          string
		updatedObject models.UpdateStudioInput
		wantErr       bool
	}{
		{
			"full",
			models.UpdateStudioInput{
				Studio: &models.Studio{
					ID:            studioIDs[studioIdxWithGallery],
					Name:          name,
					URLs:          models.NewRelatedStrings([]string{url}),
					Favorite:      favorite,
					Rating:        &rating,
					Details:       details,
					IgnoreAutoTag: ignoreAutoTag,
					Aliases:       models.NewRelatedStrings(aliases),
					TagIDs:        models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithStudio]}),
					StashIDs: models.NewRelatedStashIDs([]models.StashID{
						{
							StashID:   stashID1,
							Endpoint:  endpoint1,
							UpdatedAt: epochTime,
						},
						{
							StashID:   stashID2,
							Endpoint:  endpoint2,
							UpdatedAt: epochTime,
						},
					}),
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
			},
			false,
		},
		{
			"clear nullables",
			models.UpdateStudioInput{
				Studio: &models.Studio{
					ID:       studioIDs[studioIdxWithGallery],
					Name:     name, // name is mandatory
					URLs:     models.NewRelatedStrings([]string{}),
					Aliases:  models.NewRelatedStrings([]string{}),
					TagIDs:   models.NewRelatedIDs([]int{}),
					StashIDs: models.NewRelatedStashIDs([]models.StashID{}),
				},
			},
			false,
		},
		{
			"clear tag ids",
			models.UpdateStudioInput{
				Studio: &models.Studio{
					ID:     studioIDs[sceneIdxWithTag],
					Name:   name, // name is mandatory
					TagIDs: models.NewRelatedIDs([]int{}),
				},
			},
			false,
		},
		{
			"set custom fields",
			models.UpdateStudioInput{
				Studio: &models.Studio{
					ID:   studioIDs[studioIdxWithGallery],
					Name: name, // name is mandatory
				},
				CustomFields: models.CustomFieldsInput{
					Full: testCustomFields,
				},
			},
			false,
		},
		{
			"clear custom fields",
			models.UpdateStudioInput{
				Studio: &models.Studio{
					ID:   studioIDs[studioIdxWithGallery],
					Name: name, // name is mandatory
				},
				CustomFields: models.CustomFieldsInput{
					Full: map[string]interface{}{},
				},
			},
			false,
		},
		{
			"invalid tag id",
			models.UpdateStudioInput{
				Studio: &models.Studio{
					ID:     studioIDs[sceneIdxWithGallery],
					Name:   name, // name is mandatory
					TagIDs: models.NewRelatedIDs([]int{invalidID}),
				},
			},
			true,
		},
	}

	qb := db.Studio
	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			copy := *tt.updatedObject.Studio

			if err := qb.Update(ctx, &tt.updatedObject); (err != nil) != tt.wantErr {
				t.Errorf("StudioStore.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			s, err := qb.Find(ctx, tt.updatedObject.ID)
			if err != nil {
				t.Errorf("StudioStore.Find() error = %v", err)
			}

			// load relationships
			if err := loadStudioRelationships(ctx, copy, s); err != nil {
				t.Errorf("loadStudioRelationships() error = %v", err)
				return
			}

			assert.Equal(copy, *s)

			// ensure custom fields are correct
			if tt.updatedObject.CustomFields.Full != nil {
				cf, err := qb.GetCustomFields(ctx, tt.updatedObject.ID)
				if err != nil {
					t.Errorf("StudioStore.GetCustomFields() error = %v", err)
					return
				}

				assert.Equal(tt.updatedObject.CustomFields.Full, cf)
			}
		})
	}
}

func clearStudioPartial() models.StudioPartial {
	nullString := models.OptionalString{Set: true, Null: true}
	nullInt := models.OptionalInt{Set: true, Null: true}

	// leave mandatory fields
	return models.StudioPartial{
		URLs:     &models.UpdateStrings{Mode: models.RelationshipUpdateModeSet},
		Aliases:  &models.UpdateStrings{Mode: models.RelationshipUpdateModeSet},
		Rating:   nullInt,
		Details:  nullString,
		TagIDs:   &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
		StashIDs: &models.UpdateStashIDs{Mode: models.RelationshipUpdateModeSet},
	}
}

func Test_StudioStore_UpdatePartial(t *testing.T) {
	var (
		name          = "name"
		details       = "details"
		url           = "url"
		aliases       = []string{"aliasX", "aliasY"}
		rating        = 3
		ignoreAutoTag = true
		favorite      = true
		endpoint1     = "endpoint1"
		endpoint2     = "endpoint2"
		stashID1      = "stashid1"
		stashID2      = "stashid2"
		createdAt     = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt     = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name    string
		id      int
		partial models.StudioPartial
		want    models.Studio
		wantErr bool
	}{
		{
			"full",
			studioIDs[studioIdxWithDupName],
			models.StudioPartial{
				Name: models.NewOptionalString(name),
				URLs: &models.UpdateStrings{
					Values: []string{url},
					Mode:   models.RelationshipUpdateModeSet,
				},
				Aliases: &models.UpdateStrings{
					Values: aliases,
					Mode:   models.RelationshipUpdateModeSet,
				},
				Favorite:      models.NewOptionalBool(favorite),
				Rating:        models.NewOptionalInt(rating),
				Details:       models.NewOptionalString(details),
				IgnoreAutoTag: models.NewOptionalBool(ignoreAutoTag),
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithStudio], tagIDs[tagIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeSet,
				},
				StashIDs: &models.UpdateStashIDs{
					StashIDs: []models.StashID{
						{
							StashID:   stashID1,
							Endpoint:  endpoint1,
							UpdatedAt: epochTime,
						},
						{
							StashID:   stashID2,
							Endpoint:  endpoint2,
							UpdatedAt: epochTime,
						},
					},
					Mode: models.RelationshipUpdateModeSet,
				},
				CreatedAt: models.NewOptionalTime(createdAt),
				UpdatedAt: models.NewOptionalTime(updatedAt),
			},
			models.Studio{
				ID:            studioIDs[studioIdxWithDupName],
				Name:          name,
				URLs:          models.NewRelatedStrings([]string{url}),
				Aliases:       models.NewRelatedStrings(aliases),
				Favorite:      favorite,
				Rating:        &rating,
				Details:       details,
				IgnoreAutoTag: ignoreAutoTag,
				TagIDs:        models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithStudio]}),
				StashIDs: models.NewRelatedStashIDs([]models.StashID{
					{
						StashID:   stashID1,
						Endpoint:  endpoint1,
						UpdatedAt: epochTime,
					},
					{
						StashID:   stashID2,
						Endpoint:  endpoint2,
						UpdatedAt: epochTime,
					},
				}),
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"clear all",
			studioIDs[studioIdxWithTwoTags],
			clearStudioPartial(),
			models.Studio{
				ID:            studioIDs[studioIdxWithTwoTags],
				Name:          getStudioStringValue(studioIdxWithTwoTags, "Name"),
				Favorite:      getStudioBoolValue(studioIdxWithTwoTags),
				Aliases:       models.NewRelatedStrings([]string{}),
				TagIDs:        models.NewRelatedIDs([]int{}),
				StashIDs:      models.NewRelatedStashIDs([]models.StashID{}),
				IgnoreAutoTag: getIgnoreAutoTag(studioIdxWithTwoTags),
			},
			false,
		},
		{
			"invalid id",
			invalidID,
			models.StudioPartial{Name: models.NewOptionalString(name)},
			models.Studio{},
			true,
		},
	}
	for _, tt := range tests {
		qb := db.Studio

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			tt.partial.ID = tt.id

			got, err := qb.UpdatePartial(ctx, tt.partial)
			if (err != nil) != tt.wantErr {
				t.Errorf("StudioStore.UpdatePartial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if err := loadStudioRelationships(ctx, tt.want, got); err != nil {
				t.Errorf("loadStudioRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, *got)

			s, err := qb.Find(ctx, tt.id)
			if err != nil {
				t.Errorf("StudioStore.Find() error = %v", err)
			}

			// load relationships
			if err := loadStudioRelationships(ctx, tt.want, s); err != nil {
				t.Errorf("loadStudioRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, *s)
		})
	}
}

func Test_StudioStore_UpdatePartialCustomFields(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		partial  models.StudioPartial
		expected map[string]interface{} // nil to use the partial
	}{
		{
			"set custom fields",
			studioIDs[studioIdxWithGallery],
			models.StudioPartial{
				CustomFields: models.CustomFieldsInput{
					Full: testCustomFields,
				},
			},
			nil,
		},
		{
			"clear custom fields",
			studioIDs[studioIdxWithGallery],
			models.StudioPartial{
				CustomFields: models.CustomFieldsInput{
					Full: map[string]interface{}{},
				},
			},
			nil,
		},
		{
			"partial custom fields",
			studioIDs[studioIdxWithGallery],
			models.StudioPartial{
				CustomFields: models.CustomFieldsInput{
					Partial: map[string]interface{}{
						"string":    "bbb",
						"new_field": "new",
					},
				},
			},
			map[string]interface{}{
				"int":       int64(2),
				"real":      0.7,
				"string":    "bbb",
				"new_field": "new",
			},
		},
	}
	for _, tt := range tests {
		qb := db.Studio

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			tt.partial.ID = tt.id

			_, err := qb.UpdatePartial(ctx, tt.partial)
			if err != nil {
				t.Errorf("StudioStore.UpdatePartial() error = %v", err)
				return
			}

			// ensure custom fields are correct
			cf, err := qb.GetCustomFields(ctx, tt.id)
			if err != nil {
				t.Errorf("StudioStore.GetCustomFields() error = %v", err)
				return
			}
			if tt.expected == nil {
				assert.Equal(tt.partial.CustomFields.Full, cf)
			} else {
				assert.Equal(tt.expected, cf)
			}
		})
	}
}

func TestStudioQueryNameOr(t *testing.T) {
	const studio1Idx = 1
	const studio2Idx = 2

	studio1Name := getStudioStringValue(studio1Idx, "Name")
	studio2Name := getStudioStringValue(studio2Idx, "Name")

	studioFilter := models.StudioFilterType{
		Name: &models.StringCriterionInput{
			Value:    studio1Name,
			Modifier: models.CriterionModifierEquals,
		},
		OperatorFilter: models.OperatorFilter[models.StudioFilterType]{
			Or: &models.StudioFilterType{
				Name: &models.StringCriterionInput{
					Value:    studio2Name,
					Modifier: models.CriterionModifierEquals,
				},
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Studio

		studios := queryStudio(ctx, t, sqb, &studioFilter, nil)

		assert.Len(t, studios, 2)
		assert.Equal(t, studio1Name, studios[0].Name)
		assert.Equal(t, studio2Name, studios[1].Name)

		return nil
	})
}

func TestStudioQueryNameAndUrl(t *testing.T) {
	const studioIdx = 1
	studioName := getStudioStringValue(studioIdx, "Name")
	studioUrl := getStudioNullStringValue(studioIdx, urlField)

	studioFilter := models.StudioFilterType{
		Name: &models.StringCriterionInput{
			Value:    studioName,
			Modifier: models.CriterionModifierEquals,
		},
		OperatorFilter: models.OperatorFilter[models.StudioFilterType]{
			And: &models.StudioFilterType{
				URL: &models.StringCriterionInput{
					Value:    studioUrl,
					Modifier: models.CriterionModifierEquals,
				},
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Studio

		studios := queryStudio(ctx, t, sqb, &studioFilter, nil)

		if !assert.Len(t, studios, 1) {
			return nil
		}

		if err := studios[0].LoadURLs(ctx, db.Studio); err != nil {
			t.Errorf("Error loading studio relationships: %v", err)
		}

		assert.Equal(t, studioName, studios[0].Name)
		assert.Equal(t, []string{studioUrl}, studios[0].URLs.List())

		return nil
	})
}

func TestStudioQueryNameNotUrl(t *testing.T) {
	const studioIdx = 1

	studioUrl := getStudioNullStringValue(studioIdx, urlField)

	nameCriterion := models.StringCriterionInput{
		Value:    "studio_.*1_Name",
		Modifier: models.CriterionModifierMatchesRegex,
	}

	urlCriterion := models.StringCriterionInput{
		Value:    studioUrl,
		Modifier: models.CriterionModifierEquals,
	}

	studioFilter := models.StudioFilterType{
		Name: &nameCriterion,
		OperatorFilter: models.OperatorFilter[models.StudioFilterType]{
			Not: &models.StudioFilterType{
				URL: &urlCriterion,
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Studio

		studios := queryStudio(ctx, t, sqb, &studioFilter, nil)

		for _, studio := range studios {
			if err := studio.LoadURLs(ctx, db.Studio); err != nil {
				t.Errorf("Error loading studio relationships: %v", err)
			}

			verifyString(t, studio.Name, nameCriterion)
			urlCriterion.Modifier = models.CriterionModifierNotEquals
			verifyStringList(t, studio.URLs.List(), urlCriterion)
		}

		return nil
	})
}

func TestStudioIllegalQuery(t *testing.T) {
	assert := assert.New(t)

	const studioIdx = 1
	subFilter := models.StudioFilterType{
		Name: &models.StringCriterionInput{
			Value:    getStudioStringValue(studioIdx, "Name"),
			Modifier: models.CriterionModifierEquals,
		},
	}

	studioFilter := &models.StudioFilterType{
		OperatorFilter: models.OperatorFilter[models.StudioFilterType]{
			And: &subFilter,
			Or:  &subFilter,
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Studio

		_, _, err := sqb.Query(ctx, studioFilter, nil)
		assert.NotNil(err)

		studioFilter.Or = nil
		studioFilter.Not = &subFilter
		_, _, err = sqb.Query(ctx, studioFilter, nil)
		assert.NotNil(err)

		studioFilter.And = nil
		studioFilter.Or = &subFilter
		_, _, err = sqb.Query(ctx, studioFilter, nil)
		assert.NotNil(err)

		return nil
	})
}

func TestStudioQueryIgnoreAutoTag(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		ignoreAutoTag := true
		studioFilter := models.StudioFilterType{
			IgnoreAutoTag: &ignoreAutoTag,
		}

		sqb := db.Studio

		studios := queryStudio(ctx, t, sqb, &studioFilter, nil)

		assert.Len(t, studios, int(math.Ceil(float64(totalStudios)/5)))
		for _, s := range studios {
			assert.True(t, s.IgnoreAutoTag)
		}

		return nil
	})
}

func TestStudioQueryForAutoTag(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		tqb := db.Studio

		name := studioNames[studioIdxWithGroup] // find a studio by name

		studios, err := tqb.QueryForAutoTag(ctx, []string{name})

		if err != nil {
			t.Errorf("Error finding studios: %s", err.Error())
		}

		assert.Len(t, studios, 1)
		assert.Equal(t, strings.ToLower(studioNames[studioIdxWithGroup]), strings.ToLower(studios[0].Name))

		name = getStudioStringValue(studioIdxWithGroup, "Alias")
		studios, err = tqb.QueryForAutoTag(ctx, []string{name})

		if err != nil {
			t.Errorf("Error finding studios: %s", err.Error())
		}
		if assert.Len(t, studios, 1) {
			assert.Equal(t, studioIDs[studioIdxWithGroup], studios[0].ID)
		}
		return nil
	})
}

func TestStudioQueryParent(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Studio
		studioCriterion := models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithChildStudio]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		studioFilter := models.StudioFilterType{
			Parents: &studioCriterion,
		}

		studios, _, err := sqb.Query(ctx, &studioFilter, nil)
		if err != nil {
			t.Errorf("Error querying studio: %s", err.Error())
		}

		assert.Len(t, studios, 1)

		// ensure id is correct
		assert.Equal(t, sceneIDs[studioIdxWithParentStudio], studios[0].ID)

		studioCriterion = models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithChildStudio]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getStudioStringValue(studioIdxWithParentStudio, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		studios, _, err = sqb.Query(ctx, &studioFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying studio: %s", err.Error())
		}
		assert.Len(t, studios, 0)

		return nil
	})
}

func TestStudioDestroyParent(t *testing.T) {
	const parentName = "parent"
	const childName = "child"

	// create parent and child studios
	if err := withTxn(func(ctx context.Context) error {
		createdParent, err := createStudio(ctx, db.Studio, parentName, nil, nil)
		if err != nil {
			return fmt.Errorf("Error creating parent studio: %s", err.Error())
		}

		parentID := createdParent.ID
		createdChild, err := createStudio(ctx, db.Studio, childName, &parentID, nil)
		if err != nil {
			return fmt.Errorf("Error creating child studio: %s", err.Error())
		}

		sqb := db.Studio

		// destroy the parent
		err = sqb.Destroy(ctx, createdParent.ID)
		if err != nil {
			return fmt.Errorf("Error destroying parent studio: %s", err.Error())
		}

		// destroy the child
		err = sqb.Destroy(ctx, createdChild.ID)
		if err != nil {
			return fmt.Errorf("Error destroying child studio: %s", err.Error())
		}

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestStudioFindChildren(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Studio

		studios, err := sqb.FindChildren(ctx, studioIDs[studioIdxWithChildStudio])

		if err != nil {
			t.Errorf("error calling FindChildren: %s", err.Error())
		}

		assert.Len(t, studios, 1)
		assert.Equal(t, studioIDs[studioIdxWithParentStudio], studios[0].ID)

		studios, err = sqb.FindChildren(ctx, 0)

		if err != nil {
			t.Errorf("error calling FindChildren: %s", err.Error())
		}

		assert.Len(t, studios, 0)

		return nil
	})
}

func TestStudioUpdateClearParent(t *testing.T) {
	const parentName = "clearParent_parent"
	const childName = "clearParent_child"

	// create parent and child studios
	if err := withTxn(func(ctx context.Context) error {
		createdParent, err := createStudio(ctx, db.Studio, parentName, nil, nil)
		if err != nil {
			return fmt.Errorf("Error creating parent studio: %s", err.Error())
		}

		parentID := createdParent.ID
		createdChild, err := createStudio(ctx, db.Studio, childName, &parentID, nil)
		if err != nil {
			return fmt.Errorf("Error creating child studio: %s", err.Error())
		}

		sqb := db.Studio

		// clear the parent id from the child
		input := models.StudioPartial{
			ID:       createdChild.ID,
			ParentID: models.NewOptionalIntPtr(nil),
		}

		updatedStudio, err := sqb.UpdatePartial(ctx, input)

		if err != nil {
			return fmt.Errorf("Error updated studio: %s", err.Error())
		}

		if updatedStudio.ParentID != nil {
			return errors.New("updated studio has parent ID set")
		}

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestStudioUpdateStudioImage(t *testing.T) {
	if err := withTxn(func(ctx context.Context) error {
		qb := db.Studio

		// create studio to test against
		const name = "TestStudioUpdateStudioImage"
		created, err := createStudio(ctx, db.Studio, name, nil, nil)
		if err != nil {
			return fmt.Errorf("Error creating studio: %s", err.Error())
		}

		return testUpdateImage(t, ctx, created.ID, qb.UpdateImage, qb.GetImage)
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestStudioQuerySceneCount(t *testing.T) {
	const sceneCount = 1
	sceneCountCriterion := models.IntCriterionInput{
		Value:    sceneCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyStudiosSceneCount(t, sceneCountCriterion)

	sceneCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyStudiosSceneCount(t, sceneCountCriterion)

	sceneCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyStudiosSceneCount(t, sceneCountCriterion)

	sceneCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyStudiosSceneCount(t, sceneCountCriterion)
}

func verifyStudiosSceneCount(t *testing.T, sceneCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Studio
		studioFilter := models.StudioFilterType{
			SceneCount: &sceneCountCriterion,
		}

		studios := queryStudio(ctx, t, sqb, &studioFilter, nil)
		assert.Greater(t, len(studios), 0)

		for _, studio := range studios {
			sceneCount, err := db.Scene.CountByStudioID(ctx, studio.ID)
			if err != nil {
				return err
			}
			verifyInt(t, sceneCount, sceneCountCriterion)
		}

		return nil
	})
}

func TestStudioQueryImageCount(t *testing.T) {
	const imageCount = 1
	imageCountCriterion := models.IntCriterionInput{
		Value:    imageCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyStudiosImageCount(t, imageCountCriterion)

	imageCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyStudiosImageCount(t, imageCountCriterion)

	imageCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyStudiosImageCount(t, imageCountCriterion)

	imageCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyStudiosImageCount(t, imageCountCriterion)
}

func verifyStudiosImageCount(t *testing.T, imageCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Studio
		studioFilter := models.StudioFilterType{
			ImageCount: &imageCountCriterion,
		}

		studios := queryStudio(ctx, t, sqb, &studioFilter, nil)
		assert.Greater(t, len(studios), 0)

		for _, studio := range studios {
			pp := 0

			result, err := db.Image.Query(ctx, models.ImageQueryOptions{
				QueryOptions: models.QueryOptions{
					FindFilter: &models.FindFilterType{
						PerPage: &pp,
					},
					Count: true,
				},
				ImageFilter: &models.ImageFilterType{
					Studios: &models.HierarchicalMultiCriterionInput{
						Value:    []string{strconv.Itoa(studio.ID)},
						Modifier: models.CriterionModifierIncludes,
					},
				},
			})
			if err != nil {
				return err
			}
			verifyInt(t, result.Count, imageCountCriterion)
		}

		return nil
	})
}

func TestStudioQueryGalleryCount(t *testing.T) {
	const galleryCount = 1
	galleryCountCriterion := models.IntCriterionInput{
		Value:    galleryCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyStudiosGalleryCount(t, galleryCountCriterion)

	galleryCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyStudiosGalleryCount(t, galleryCountCriterion)

	galleryCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyStudiosGalleryCount(t, galleryCountCriterion)

	galleryCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyStudiosGalleryCount(t, galleryCountCriterion)
}

func verifyStudiosGalleryCount(t *testing.T, galleryCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Studio
		studioFilter := models.StudioFilterType{
			GalleryCount: &galleryCountCriterion,
		}

		studios := queryStudio(ctx, t, sqb, &studioFilter, nil)
		assert.Greater(t, len(studios), 0)

		for _, studio := range studios {
			pp := 0

			_, count, err := db.Gallery.Query(ctx, &models.GalleryFilterType{
				Studios: &models.HierarchicalMultiCriterionInput{
					Value:    []string{strconv.Itoa(studio.ID)},
					Modifier: models.CriterionModifierIncludes,
				},
			}, &models.FindFilterType{
				PerPage: &pp,
			})
			if err != nil {
				return err
			}
			verifyInt(t, count, galleryCountCriterion)
		}

		return nil
	})
}

func TestStudioStashIDs(t *testing.T) {
	if err := withRollbackTxn(func(ctx context.Context) error {
		qb := db.Studio

		// create studio to test against
		const name = "TestStudioStashIDs"
		created, err := createStudio(ctx, db.Studio, name, nil, nil)
		if err != nil {
			return fmt.Errorf("Error creating studio: %s", err.Error())
		}

		studio, err := qb.Find(ctx, created.ID)
		if err != nil {
			return fmt.Errorf("Error getting studio: %s", err.Error())
		}

		if err := studio.LoadStashIDs(ctx, qb); err != nil {
			return err
		}

		testStudioStashIDs(ctx, t, studio)
		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func testStudioStashIDs(ctx context.Context, t *testing.T, s *models.Studio) {
	qb := db.Studio

	if err := s.LoadStashIDs(ctx, qb); err != nil {
		t.Error(err.Error())
		return
	}

	// ensure no stash IDs to begin with
	assert.Len(t, s.StashIDs.List(), 0)

	// add stash ids
	const stashIDStr = "stashID"
	const endpoint = "endpoint"
	stashID := models.StashID{
		StashID:  stashIDStr,
		Endpoint: endpoint,
	}

	// update stash ids and ensure was updated
	input := models.StudioPartial{
		ID: s.ID,
		StashIDs: &models.UpdateStashIDs{
			StashIDs: []models.StashID{stashID},
			Mode:     models.RelationshipUpdateModeSet,
		},
	}
	var err error
	s, err = qb.UpdatePartial(ctx, input)
	if err != nil {
		t.Error(err.Error())
	}

	if err := s.LoadStashIDs(ctx, qb); err != nil {
		t.Error(err.Error())
		return
	}

	// #5563 - set the UpdatedAt field to epoch
	stashID.UpdatedAt = epochTime

	assert.Equal(t, []models.StashID{stashID}, s.StashIDs.List())

	// remove stash ids and ensure was updated
	input = models.StudioPartial{
		ID: s.ID,
		StashIDs: &models.UpdateStashIDs{
			StashIDs: []models.StashID{stashID},
			Mode:     models.RelationshipUpdateModeRemove,
		},
	}
	s, err = qb.UpdatePartial(ctx, input)
	if err != nil {
		t.Error(err.Error())
	}

	if err := s.LoadStashIDs(ctx, qb); err != nil {
		t.Error(err.Error())
		return
	}

	assert.Len(t, s.StashIDs.List(), 0)
}

func TestStudioQueryURL(t *testing.T) {
	const sceneIdx = 1
	studioURL := getStudioStringValue(sceneIdx, urlField)

	urlCriterion := models.StringCriterionInput{
		Value:    studioURL,
		Modifier: models.CriterionModifierEquals,
	}

	filter := models.StudioFilterType{
		URL: &urlCriterion,
	}

	verifyFn := func(ctx context.Context, g *models.Studio) {
		t.Helper()
		if err := g.LoadURLs(ctx, db.Studio); err != nil {
			t.Errorf("Error loading studio relationships: %v", err)
			return
		}
		verifyStringList(t, g.URLs.List(), urlCriterion)
	}

	verifyStudioQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotEquals
	verifyStudioQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierMatchesRegex
	urlCriterion.Value = "studio_.*1_URL"
	verifyStudioQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotMatchesRegex
	verifyStudioQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierIsNull
	urlCriterion.Value = ""
	verifyStudioQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotNull
	verifyStudioQuery(t, filter, verifyFn)
}

func TestStudioQueryRating(t *testing.T) {
	const rating = 60
	ratingCriterion := models.IntCriterionInput{
		Value:    rating,
		Modifier: models.CriterionModifierEquals,
	}

	verifyStudiosRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotEquals
	verifyStudiosRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyStudiosRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierLessThan
	verifyStudiosRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierIsNull
	verifyStudiosRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotNull
	verifyStudiosRating(t, ratingCriterion)
}

func queryStudios(ctx context.Context, t *testing.T, studioFilter *models.StudioFilterType, findFilter *models.FindFilterType) []*models.Studio {
	t.Helper()
	studios, _, err := db.Studio.Query(ctx, studioFilter, findFilter)
	if err != nil {
		t.Errorf("Error querying studio: %s", err.Error())
	}

	return studios
}

func TestStudioQueryTags(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		tagCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdxWithStudio]),
				strconv.Itoa(tagIDs[tagIdx1WithStudio]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		studioFilter := models.StudioFilterType{
			Tags: &tagCriterion,
		}

		// ensure ids are correct
		studios := queryStudios(ctx, t, &studioFilter, nil)
		assert.Len(t, studios, 2)
		for _, studio := range studios {
			assert.True(t, studio.ID == studioIDs[studioIdxWithTag] || studio.ID == studioIDs[studioIdxWithTwoTags])
		}

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithStudio]),
				strconv.Itoa(tagIDs[tagIdx2WithStudio]),
			},
			Modifier: models.CriterionModifierIncludesAll,
		}

		studios = queryStudios(ctx, t, &studioFilter, nil)

		assert.Len(t, studios, 1)
		assert.Equal(t, sceneIDs[studioIdxWithTwoTags], studios[0].ID)

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithStudio]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getSceneStringValue(studioIdxWithTwoTags, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		studios = queryStudios(ctx, t, &studioFilter, &findFilter)
		assert.Len(t, studios, 0)

		return nil
	})
}

func TestStudioQueryTagCount(t *testing.T) {
	const tagCount = 1
	tagCountCriterion := models.IntCriterionInput{
		Value:    tagCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyStudiosTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyStudiosTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyStudiosTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyStudiosTagCount(t, tagCountCriterion)
}

func verifyStudiosTagCount(t *testing.T, tagCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Studio
		studioFilter := models.StudioFilterType{
			TagCount: &tagCountCriterion,
		}

		studios := queryStudios(ctx, t, &studioFilter, nil)
		assert.Greater(t, len(studios), 0)

		for _, studio := range studios {
			ids, err := sqb.GetTagIDs(ctx, studio.ID)
			if err != nil {
				return err
			}
			verifyInt(t, len(ids), tagCountCriterion)
		}

		return nil
	})
}

func verifyStudioQuery(t *testing.T, filter models.StudioFilterType, verifyFn func(ctx context.Context, s *models.Studio)) {
	withTxn(func(ctx context.Context) error {
		t.Helper()
		sqb := db.Studio

		studios := queryStudio(ctx, t, sqb, &filter, nil)

		// assume it should find at least one
		assert.Greater(t, len(studios), 0)

		for _, studio := range studios {
			verifyFn(ctx, studio)
		}

		return nil
	})
}

func verifyStudiosRating(t *testing.T, ratingCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Studio
		studioFilter := models.StudioFilterType{
			Rating100: &ratingCriterion,
		}

		studios, _, err := sqb.Query(ctx, &studioFilter, nil)

		if err != nil {
			t.Errorf("Error querying studio: %s", err.Error())
		}

		for _, studio := range studios {
			verifyIntPtr(t, studio.Rating, ratingCriterion)
		}

		return nil
	})
}

func TestStudioQueryIsMissingRating(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Studio
		isMissing := "rating"
		studioFilter := models.StudioFilterType{
			IsMissing: &isMissing,
		}

		studios, _, err := sqb.Query(ctx, &studioFilter, nil)

		if err != nil {
			t.Errorf("Error querying studio: %s", err.Error())
		}

		assert.True(t, len(studios) > 0)

		for _, studio := range studios {
			assert.Nil(t, studio.Rating)
		}

		return nil
	})
}

func queryStudio(ctx context.Context, t *testing.T, sqb models.StudioReader, studioFilter *models.StudioFilterType, findFilter *models.FindFilterType) []*models.Studio {
	studios, _, err := sqb.Query(ctx, studioFilter, findFilter)
	if err != nil {
		t.Errorf("Error querying studio: %s", err.Error())
	}

	return studios
}

func TestStudioQueryName(t *testing.T) {
	const studioIdx = 1
	studioName := getStudioStringValue(studioIdx, "Name")

	nameCriterion := &models.StringCriterionInput{
		Value:    studioName,
		Modifier: models.CriterionModifierEquals,
	}

	studioFilter := models.StudioFilterType{
		Name: nameCriterion,
	}

	verifyFn := func(ctx context.Context, studio *models.Studio) {
		verifyString(t, studio.Name, *nameCriterion)
	}

	verifyStudioQuery(t, studioFilter, verifyFn)

	nameCriterion.Modifier = models.CriterionModifierNotEquals
	verifyStudioQuery(t, studioFilter, verifyFn)

	nameCriterion.Modifier = models.CriterionModifierMatchesRegex
	nameCriterion.Value = "studio_.*1_Name"
	verifyStudioQuery(t, studioFilter, verifyFn)

	nameCriterion.Modifier = models.CriterionModifierNotMatchesRegex
	verifyStudioQuery(t, studioFilter, verifyFn)
}

func TestStudioQueryAlias(t *testing.T) {
	const studioIdx = studioIdxWithGroup
	studioName := getStudioStringValue(studioIdx, "Alias")

	aliasCriterion := &models.StringCriterionInput{
		Value:    studioName,
		Modifier: models.CriterionModifierEquals,
	}

	studioFilter := models.StudioFilterType{
		Aliases: aliasCriterion,
	}

	verifyFn := func(ctx context.Context, studio *models.Studio) {
		t.Helper()
		aliases, err := db.Studio.GetAliases(ctx, studio.ID)
		if err != nil {
			t.Errorf("Error querying studios: %s", err.Error())
		}

		var alias string
		if len(aliases) > 0 {
			alias = aliases[0]
		}

		verifyString(t, alias, *aliasCriterion)
	}

	verifyStudioQuery(t, studioFilter, verifyFn)

	aliasCriterion.Modifier = models.CriterionModifierNotEquals
	verifyStudioQuery(t, studioFilter, verifyFn)

	aliasCriterion.Modifier = models.CriterionModifierMatchesRegex
	aliasCriterion.Value = "studio_.*2_Alias"
	verifyStudioQuery(t, studioFilter, verifyFn)

	aliasCriterion.Modifier = models.CriterionModifierNotMatchesRegex
	verifyStudioQuery(t, studioFilter, verifyFn)

	aliasCriterion.Modifier = models.CriterionModifierIsNull
	aliasCriterion.Value = ""
	verifyStudioQuery(t, studioFilter, verifyFn)

	aliasCriterion.Modifier = models.CriterionModifierNotNull
	verifyStudioQuery(t, studioFilter, verifyFn)
}

func TestStudioAlias(t *testing.T) {
	if err := withRollbackTxn(func(ctx context.Context) error {
		qb := db.Studio

		// create studio to test against
		const name = "TestStudioAlias"
		created, err := createStudio(ctx, db.Studio, name, nil, nil)
		if err != nil {
			return fmt.Errorf("Error creating studio: %s", err.Error())
		}

		studio, err := qb.Find(ctx, created.ID)
		if err != nil {
			return fmt.Errorf("Error getting studio: %s", err.Error())
		}

		if err := studio.LoadStashIDs(ctx, qb); err != nil {
			return err
		}

		testStudioAlias(ctx, t, studio)
		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func testStudioAlias(ctx context.Context, t *testing.T, s *models.Studio) {
	qb := db.Studio
	if err := s.LoadAliases(ctx, qb); err != nil {
		t.Error(err.Error())
		return
	}

	// ensure no alias to begin with
	assert.Len(t, s.Aliases.List(), 0)

	aliases := []string{"alias1", "alias2"}

	// update alias and ensure was updated
	input := models.StudioPartial{
		ID: s.ID,
		Aliases: &models.UpdateStrings{
			Values: aliases,
			Mode:   models.RelationshipUpdateModeSet,
		},
	}
	var err error
	s, err = qb.UpdatePartial(ctx, input)
	if err != nil {
		t.Error(err.Error())
	}

	if err := s.LoadAliases(ctx, qb); err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, aliases, s.Aliases.List())

	// remove alias and ensure was updated
	input = models.StudioPartial{
		ID: s.ID,
		Aliases: &models.UpdateStrings{
			Values: aliases,
			Mode:   models.RelationshipUpdateModeRemove,
		},
	}
	s, err = qb.UpdatePartial(ctx, input)
	if err != nil {
		t.Error(err.Error())
	}

	if err := s.LoadAliases(ctx, qb); err != nil {
		t.Error(err.Error())
		return
	}

	assert.Len(t, s.Aliases.List(), 0)
}

// TestStudioQueryFast does a quick test for major errors, no result verification
func TestStudioQueryFast(t *testing.T) {

	tsString := "test"
	tsInt := 1

	testStringCriterion := models.StringCriterionInput{
		Value:    tsString,
		Modifier: models.CriterionModifierEquals,
	}
	testIncludesMultiCriterion := models.MultiCriterionInput{
		Value:    []string{tsString},
		Modifier: models.CriterionModifierIncludes,
	}
	testIntCriterion := models.IntCriterionInput{
		Value:    tsInt,
		Modifier: models.CriterionModifierEquals,
	}

	nameFilter := models.StudioFilterType{
		Name: &testStringCriterion,
	}
	aliasesFilter := models.StudioFilterType{
		Aliases: &testStringCriterion,
	}
	stashIDFilter := models.StudioFilterType{
		StashID: &testStringCriterion,
	}
	urlFilter := models.StudioFilterType{
		URL: &testStringCriterion,
	}
	ratingFilter := models.StudioFilterType{
		Rating100: &testIntCriterion,
	}
	sceneCountFilter := models.StudioFilterType{
		SceneCount: &testIntCriterion,
	}
	imageCountFilter := models.StudioFilterType{
		SceneCount: &testIntCriterion,
	}
	parentsFilter := models.StudioFilterType{
		Parents: &testIncludesMultiCriterion,
	}

	filters := []models.StudioFilterType{nameFilter, aliasesFilter, stashIDFilter, urlFilter, ratingFilter, sceneCountFilter, imageCountFilter, parentsFilter}

	missingStrings := []string{"image", "stash_id", "details"}

	for _, m := range missingStrings {
		filters = append(filters, models.StudioFilterType{
			IsMissing: &m,
		})
	}

	sortbyStrings := []string{"scenes_count", "images_count", "galleries_count", "created_at", "updated_at", "name", "random_26819649", "rating"}

	var findFilters []models.FindFilterType

	for _, sb := range sortbyStrings {
		findFilters = append(findFilters, models.FindFilterType{
			Q:       &tsString,
			Page:    &tsInt,
			PerPage: &tsInt,
			Sort:    &sb,
		})

	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Studio
		for _, f := range filters {
			for _, ff := range findFilters {
				_, _, err := sqb.Query(ctx, &f, &ff)
				if err != nil {
					t.Errorf("Error querying studio: %s", err.Error())
				}
			}
		}

		return nil
	})
}

func studiesToIDs(i []*models.Studio) []int {
	ret := make([]int, len(i))
	for i, v := range i {
		ret[i] = v.ID
	}

	return ret
}

func TestStudioQueryCustomFields(t *testing.T) {
	tests := []struct {
		name        string
		filter      *models.StudioFilterType
		includeIdxs []int
		excludeIdxs []int
		wantErr     bool
	}{
		{
			"equals",
			&models.StudioFilterType{
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierEquals,
						Value:    []any{getStudioStringValue(studioIdxWithTwoScenes, "custom")},
					},
				},
			},
			[]int{studioIdxWithTwoScenes},
			nil,
			false,
		},
		{
			"not equals",
			&models.StudioFilterType{
				Name: &models.StringCriterionInput{
					Value:    getStudioStringValue(studioIdxWithTwoScenes, "Name"),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierNotEquals,
						Value:    []any{getStudioStringValue(studioIdxWithTwoScenes, "custom")},
					},
				},
			},
			nil,
			[]int{studioIdxWithTwoScenes},
			false,
		},
		{
			"includes",
			&models.StudioFilterType{
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierIncludes,
						Value:    []any{getStudioStringValue(studioIdxWithTwoScenes, "custom")[9:]},
					},
				},
			},
			[]int{studioIdxWithTwoScenes},
			nil,
			false,
		},
		{
			"excludes",
			&models.StudioFilterType{
				Name: &models.StringCriterionInput{
					Value:    getStudioStringValue(studioIdxWithTwoScenes, "Name"),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierExcludes,
						Value:    []any{getStudioStringValue(studioIdxWithTwoScenes, "custom")[9:]},
					},
				},
			},
			nil,
			[]int{studioIdxWithTwoScenes},
			false,
		},
		{
			"regex",
			&models.StudioFilterType{
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierMatchesRegex,
						Value:    []any{".*1_custom"},
					},
				},
			},
			[]int{studioIdxWithTwoScenes},
			nil,
			false,
		},
		{
			"invalid regex",
			&models.StudioFilterType{
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierMatchesRegex,
						Value:    []any{"["},
					},
				},
			},
			nil,
			nil,
			true,
		},
		{
			"not matches regex",
			&models.StudioFilterType{
				Name: &models.StringCriterionInput{
					Value:    getStudioStringValue(studioIdxWithTwoScenes, "Name"),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierNotMatchesRegex,
						Value:    []any{".*1_custom"},
					},
				},
			},
			nil,
			[]int{studioIdxWithTwoScenes},
			false,
		},
		{
			"invalid not matches regex",
			&models.StudioFilterType{
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierNotMatchesRegex,
						Value:    []any{"["},
					},
				},
			},
			nil,
			nil,
			true,
		},
		{
			"null",
			&models.StudioFilterType{
				Name: &models.StringCriterionInput{
					Value:    getStudioStringValue(studioIdxWithTwoScenes, "Name"),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "not existing",
						Modifier: models.CriterionModifierIsNull,
					},
				},
			},
			[]int{studioIdxWithTwoScenes},
			nil,
			false,
		},
		{
			"not null",
			&models.StudioFilterType{
				Name: &models.StringCriterionInput{
					Value:    getStudioStringValue(studioIdxWithTwoScenes, "Name"),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierNotNull,
					},
				},
			},
			[]int{studioIdxWithTwoScenes},
			nil,
			false,
		},
		{
			"between",
			&models.StudioFilterType{
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "real",
						Modifier: models.CriterionModifierBetween,
						Value:    []any{0.15, 0.25},
					},
				},
			},
			[]int{studioIdxWithGroup},
			nil,
			false,
		},
		{
			"not between",
			&models.StudioFilterType{
				Name: &models.StringCriterionInput{
					Value:    getStudioStringValue(studioIdxWithGroup, "Name"),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "real",
						Modifier: models.CriterionModifierNotBetween,
						Value:    []any{0.15, 0.25},
					},
				},
			},
			nil,
			[]int{studioIdxWithGroup},
			false,
		},
	}

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			studios, _, err := db.Studio.Query(ctx, tt.filter, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("StudioStore.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ids := studiesToIDs(studios)
			include := indexesToIDs(studioIDs, tt.includeIdxs)
			exclude := indexesToIDs(studioIDs, tt.excludeIdxs)

			for _, i := range include {
				assert.Contains(ids, i)
			}
			for _, e := range exclude {
				assert.NotContains(ids, e)
			}
		})
	}
}

// TODO Create
// TODO Update
// TODO Destroy
// TODO Find
// TODO FindBySceneID
// TODO Count
// TODO All
// TODO AllSlim
// TODO Query
