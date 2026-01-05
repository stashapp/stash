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
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

var testCustomFields = map[string]interface{}{
	"string": "aaa",
	"int":    int64(123), // int64 to match the type of the field in the database
	"real":   1.23,
}

func loadPerformerRelationships(ctx context.Context, expected models.Performer, actual *models.Performer) error {
	if expected.Aliases.Loaded() {
		if err := actual.LoadAliases(ctx, db.Performer); err != nil {
			return err
		}
	}
	if expected.URLs.Loaded() {
		if err := actual.LoadURLs(ctx, db.Performer); err != nil {
			return err
		}
	}
	if expected.TagIDs.Loaded() {
		if err := actual.LoadTagIDs(ctx, db.Performer); err != nil {
			return err
		}
	}
	if expected.StashIDs.Loaded() {
		if err := actual.LoadStashIDs(ctx, db.Performer); err != nil {
			return err
		}
	}

	return nil
}

func Test_PerformerStore_Create(t *testing.T) {
	var (
		name           = "name"
		disambiguation = "disambiguation"
		gender         = models.GenderEnumFemale
		details        = "details"
		url            = "url"
		twitter        = "twitter"
		instagram      = "instagram"
		urls           = []string{url, twitter, instagram}
		rating         = 3
		ethnicity      = "ethnicity"
		country        = "country"
		eyeColor       = "eyeColor"
		height         = 134
		measurements   = "measurements"
		fakeTits       = "fakeTits"
		penisLength    = 1.23
		circumcised    = models.CircumisedEnumCut
		careerLength   = "careerLength"
		tattoos        = "tattoos"
		piercings      = "piercings"
		aliases        = []string{"alias1", "alias2"}
		hairColor      = "hairColor"
		weight         = 123
		ignoreAutoTag  = true
		favorite       = true
		endpoint1      = "endpoint1"
		endpoint2      = "endpoint2"
		stashID1       = "stashid1"
		stashID2       = "stashid2"
		createdAt      = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt      = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)

		birthdate, _ = models.ParseDate("2003-02-01")
		deathdate, _ = models.ParseDate("2023-02-01")
	)

	tests := []struct {
		name      string
		newObject models.CreatePerformerInput
		wantErr   bool
	}{
		{
			"full",
			models.CreatePerformerInput{
				Performer: &models.Performer{
					Name:           name,
					Disambiguation: disambiguation,
					Gender:         &gender,
					URLs:           models.NewRelatedStrings(urls),
					Birthdate:      &birthdate,
					Ethnicity:      ethnicity,
					Country:        country,
					EyeColor:       eyeColor,
					Height:         &height,
					Measurements:   measurements,
					FakeTits:       fakeTits,
					PenisLength:    &penisLength,
					Circumcised:    &circumcised,
					CareerLength:   careerLength,
					Tattoos:        tattoos,
					Piercings:      piercings,
					Favorite:       favorite,
					Rating:         &rating,
					Details:        details,
					DeathDate:      &deathdate,
					HairColor:      hairColor,
					Weight:         &weight,
					IgnoreAutoTag:  ignoreAutoTag,
					TagIDs:         models.NewRelatedIDs([]int{tagIDs[tagIdx1WithPerformer], tagIDs[tagIdx1WithDupName]}),
					Aliases:        models.NewRelatedStrings(aliases),
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
			models.CreatePerformerInput{
				Performer: &models.Performer{
					Name:   name,
					TagIDs: models.NewRelatedIDs([]int{invalidID}),
				},
			},
			true,
		},
	}

	qb := db.Performer

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			p := tt.newObject
			if err := qb.Create(ctx, &p); (err != nil) != tt.wantErr {
				t.Errorf("PerformerStore.Create() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				assert.Zero(p.ID)
				return
			}

			assert.NotZero(p.ID)

			copy := *tt.newObject.Performer
			copy.ID = p.ID

			// load relationships
			if err := loadPerformerRelationships(ctx, copy, p.Performer); err != nil {
				t.Errorf("loadPerformerRelationships() error = %v", err)
				return
			}

			assert.Equal(copy, *p.Performer)

			// ensure can find the performer
			found, err := qb.Find(ctx, p.ID)
			if err != nil {
				t.Errorf("PerformerStore.Find() error = %v", err)
			}

			if !assert.NotNil(found) {
				return
			}

			// load relationships
			if err := loadPerformerRelationships(ctx, copy, found); err != nil {
				t.Errorf("loadPerformerRelationships() error = %v", err)
				return
			}
			assert.Equal(copy, *found)

			// ensure custom fields are set
			cf, err := qb.GetCustomFields(ctx, p.ID)
			if err != nil {
				t.Errorf("PerformerStore.GetCustomFields() error = %v", err)
				return
			}

			assert.Equal(tt.newObject.CustomFields, cf)

			return
		})
	}
}

func Test_PerformerStore_Update(t *testing.T) {
	var (
		name           = "name"
		disambiguation = "disambiguation"
		gender         = models.GenderEnumFemale
		details        = "details"
		url            = "url"
		twitter        = "twitter"
		instagram      = "instagram"
		urls           = []string{url, twitter, instagram}
		rating         = 3
		ethnicity      = "ethnicity"
		country        = "country"
		eyeColor       = "eyeColor"
		height         = 134
		measurements   = "measurements"
		fakeTits       = "fakeTits"
		penisLength    = 1.23
		circumcised    = models.CircumisedEnumCut
		careerLength   = "careerLength"
		tattoos        = "tattoos"
		piercings      = "piercings"
		aliases        = []string{"alias1", "alias2"}
		hairColor      = "hairColor"
		weight         = 123
		ignoreAutoTag  = true
		favorite       = true
		endpoint1      = "endpoint1"
		endpoint2      = "endpoint2"
		stashID1       = "stashid1"
		stashID2       = "stashid2"
		createdAt      = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt      = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)

		birthdate, _ = models.ParseDate("2003-02-01")
		deathdate, _ = models.ParseDate("2023-02-01")
	)

	tests := []struct {
		name          string
		updatedObject models.UpdatePerformerInput
		wantErr       bool
	}{
		{
			"full",
			models.UpdatePerformerInput{
				Performer: &models.Performer{
					ID:             performerIDs[performerIdxWithGallery],
					Name:           name,
					Disambiguation: disambiguation,
					Gender:         &gender,
					URLs:           models.NewRelatedStrings(urls),
					Birthdate:      &birthdate,
					Ethnicity:      ethnicity,
					Country:        country,
					EyeColor:       eyeColor,
					Height:         &height,
					Measurements:   measurements,
					FakeTits:       fakeTits,
					PenisLength:    &penisLength,
					Circumcised:    &circumcised,
					CareerLength:   careerLength,
					Tattoos:        tattoos,
					Piercings:      piercings,
					Favorite:       favorite,
					Rating:         &rating,
					Details:        details,
					DeathDate:      &deathdate,
					HairColor:      hairColor,
					Weight:         &weight,
					IgnoreAutoTag:  ignoreAutoTag,
					Aliases:        models.NewRelatedStrings(aliases),
					TagIDs:         models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithPerformer]}),
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
			models.UpdatePerformerInput{
				Performer: &models.Performer{
					ID:       performerIDs[performerIdxWithGallery],
					Aliases:  models.NewRelatedStrings([]string{}),
					URLs:     models.NewRelatedStrings([]string{}),
					TagIDs:   models.NewRelatedIDs([]int{}),
					StashIDs: models.NewRelatedStashIDs([]models.StashID{}),
				},
			},
			false,
		},
		{
			"clear tag ids",
			models.UpdatePerformerInput{
				Performer: &models.Performer{
					ID:     performerIDs[sceneIdxWithTag],
					TagIDs: models.NewRelatedIDs([]int{}),
				},
			},
			false,
		},
		{
			"set custom fields",
			models.UpdatePerformerInput{
				Performer: &models.Performer{
					ID: performerIDs[performerIdxWithGallery],
				},
				CustomFields: models.CustomFieldsInput{
					Full: testCustomFields,
				},
			},
			false,
		},
		{
			"clear custom fields",
			models.UpdatePerformerInput{
				Performer: &models.Performer{
					ID: performerIDs[performerIdxWithGallery],
				},
				CustomFields: models.CustomFieldsInput{
					Full: map[string]interface{}{},
				},
			},
			false,
		},
		{
			"invalid tag id",
			models.UpdatePerformerInput{
				Performer: &models.Performer{
					ID:     performerIDs[sceneIdxWithGallery],
					TagIDs: models.NewRelatedIDs([]int{invalidID}),
				},
			},
			true,
		},
	}

	qb := db.Performer
	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			copy := *tt.updatedObject.Performer

			if err := qb.Update(ctx, &tt.updatedObject); (err != nil) != tt.wantErr {
				t.Errorf("PerformerStore.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			s, err := qb.Find(ctx, tt.updatedObject.ID)
			if err != nil {
				t.Errorf("PerformerStore.Find() error = %v", err)
			}

			// load relationships
			if err := loadPerformerRelationships(ctx, copy, s); err != nil {
				t.Errorf("loadPerformerRelationships() error = %v", err)
				return
			}

			assert.Equal(copy, *s)

			// ensure custom fields are correct
			if tt.updatedObject.CustomFields.Full != nil {
				cf, err := qb.GetCustomFields(ctx, tt.updatedObject.ID)
				if err != nil {
					t.Errorf("PerformerStore.GetCustomFields() error = %v", err)
					return
				}

				assert.Equal(tt.updatedObject.CustomFields.Full, cf)
			}
		})
	}
}

func clearPerformerPartial() models.PerformerPartial {
	nullString := models.OptionalString{Set: true, Null: true}
	nullDate := models.OptionalDate{Set: true, Null: true}
	nullInt := models.OptionalInt{Set: true, Null: true}
	nullFloat := models.OptionalFloat64{Set: true, Null: true}

	// leave mandatory fields
	return models.PerformerPartial{
		Disambiguation: nullString,
		Gender:         nullString,
		URLs:           &models.UpdateStrings{Mode: models.RelationshipUpdateModeSet},
		Birthdate:      nullDate,
		Ethnicity:      nullString,
		Country:        nullString,
		EyeColor:       nullString,
		Height:         nullInt,
		Measurements:   nullString,
		FakeTits:       nullString,
		PenisLength:    nullFloat,
		Circumcised:    nullString,
		CareerLength:   nullString,
		Tattoos:        nullString,
		Piercings:      nullString,
		Aliases:        &models.UpdateStrings{Mode: models.RelationshipUpdateModeSet},
		Rating:         nullInt,
		Details:        nullString,
		DeathDate:      nullDate,
		HairColor:      nullString,
		Weight:         nullInt,
		TagIDs:         &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
		StashIDs:       &models.UpdateStashIDs{Mode: models.RelationshipUpdateModeSet},
	}
}

func Test_PerformerStore_UpdatePartial(t *testing.T) {
	var (
		name           = "name"
		disambiguation = "disambiguation"
		gender         = models.GenderEnumFemale
		details        = "details"
		url            = "url"
		twitter        = "twitter"
		instagram      = "instagram"
		urls           = []string{url, twitter, instagram}
		rating         = 3
		ethnicity      = "ethnicity"
		country        = "country"
		eyeColor       = "eyeColor"
		height         = 143
		measurements   = "measurements"
		fakeTits       = "fakeTits"
		penisLength    = 1.23
		circumcised    = models.CircumisedEnumCut
		careerLength   = "careerLength"
		tattoos        = "tattoos"
		piercings      = "piercings"
		aliases        = []string{"alias1", "alias2"}
		hairColor      = "hairColor"
		weight         = 123
		ignoreAutoTag  = true
		favorite       = true
		endpoint1      = "endpoint1"
		endpoint2      = "endpoint2"
		stashID1       = "stashid1"
		stashID2       = "stashid2"
		createdAt      = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt      = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)

		birthdate, _ = models.ParseDate("2003-02-01")
		deathdate, _ = models.ParseDate("2023-02-01")
	)

	tests := []struct {
		name    string
		id      int
		partial models.PerformerPartial
		want    models.Performer
		wantErr bool
	}{
		{
			"full",
			performerIDs[performerIdxWithDupName],
			models.PerformerPartial{
				Name:           models.NewOptionalString(name),
				Disambiguation: models.NewOptionalString(disambiguation),
				Gender:         models.NewOptionalString(gender.String()),
				URLs: &models.UpdateStrings{
					Values: urls,
					Mode:   models.RelationshipUpdateModeSet,
				},
				Birthdate:    models.NewOptionalDate(birthdate),
				Ethnicity:    models.NewOptionalString(ethnicity),
				Country:      models.NewOptionalString(country),
				EyeColor:     models.NewOptionalString(eyeColor),
				Height:       models.NewOptionalInt(height),
				Measurements: models.NewOptionalString(measurements),
				FakeTits:     models.NewOptionalString(fakeTits),
				PenisLength:  models.NewOptionalFloat64(penisLength),
				Circumcised:  models.NewOptionalString(circumcised.String()),
				CareerLength: models.NewOptionalString(careerLength),
				Tattoos:      models.NewOptionalString(tattoos),
				Piercings:    models.NewOptionalString(piercings),
				Aliases: &models.UpdateStrings{
					Values: aliases,
					Mode:   models.RelationshipUpdateModeSet,
				},
				Favorite:      models.NewOptionalBool(favorite),
				Rating:        models.NewOptionalInt(rating),
				Details:       models.NewOptionalString(details),
				DeathDate:     models.NewOptionalDate(deathdate),
				HairColor:     models.NewOptionalString(hairColor),
				Weight:        models.NewOptionalInt(weight),
				IgnoreAutoTag: models.NewOptionalBool(ignoreAutoTag),
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithPerformer]},
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
			models.Performer{
				ID:             performerIDs[performerIdxWithDupName],
				Name:           name,
				Disambiguation: disambiguation,
				Gender:         &gender,
				URLs:           models.NewRelatedStrings(urls),
				Birthdate:      &birthdate,
				Ethnicity:      ethnicity,
				Country:        country,
				EyeColor:       eyeColor,
				Height:         &height,
				Measurements:   measurements,
				FakeTits:       fakeTits,
				PenisLength:    &penisLength,
				Circumcised:    &circumcised,
				CareerLength:   careerLength,
				Tattoos:        tattoos,
				Piercings:      piercings,
				Aliases:        models.NewRelatedStrings(aliases),
				Favorite:       favorite,
				Rating:         &rating,
				Details:        details,
				DeathDate:      &deathdate,
				HairColor:      hairColor,
				Weight:         &weight,
				IgnoreAutoTag:  ignoreAutoTag,
				TagIDs:         models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithPerformer]}),
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
			performerIDs[performerIdxWithTwoTags],
			clearPerformerPartial(),
			models.Performer{
				ID:            performerIDs[performerIdxWithTwoTags],
				Name:          getPerformerStringValue(performerIdxWithTwoTags, "Name"),
				Favorite:      getPerformerBoolValue(performerIdxWithTwoTags),
				URLs:          models.NewRelatedStrings([]string{}),
				Aliases:       models.NewRelatedStrings([]string{}),
				TagIDs:        models.NewRelatedIDs([]int{}),
				StashIDs:      models.NewRelatedStashIDs([]models.StashID{}),
				IgnoreAutoTag: getIgnoreAutoTag(performerIdxWithTwoTags),
			},
			false,
		},
		{
			"invalid id",
			invalidID,
			models.PerformerPartial{},
			models.Performer{},
			true,
		},
	}
	for _, tt := range tests {
		qb := db.Performer

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			got, err := qb.UpdatePartial(ctx, tt.id, tt.partial)
			if (err != nil) != tt.wantErr {
				t.Errorf("PerformerStore.UpdatePartial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if err := loadPerformerRelationships(ctx, tt.want, got); err != nil {
				t.Errorf("loadPerformerRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, *got)

			s, err := qb.Find(ctx, tt.id)
			if err != nil {
				t.Errorf("PerformerStore.Find() error = %v", err)
			}

			// load relationships
			if err := loadPerformerRelationships(ctx, tt.want, s); err != nil {
				t.Errorf("loadPerformerRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, *s)
		})
	}
}

func Test_PerformerStore_UpdatePartialCustomFields(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		partial  models.PerformerPartial
		expected map[string]interface{} // nil to use the partial
	}{
		{
			"set custom fields",
			performerIDs[performerIdxWithGallery],
			models.PerformerPartial{
				CustomFields: models.CustomFieldsInput{
					Full: testCustomFields,
				},
			},
			nil,
		},
		{
			"clear custom fields",
			performerIDs[performerIdxWithGallery],
			models.PerformerPartial{
				CustomFields: models.CustomFieldsInput{
					Full: map[string]interface{}{},
				},
			},
			nil,
		},
		{
			"partial custom fields",
			performerIDs[performerIdxWithGallery],
			models.PerformerPartial{
				CustomFields: models.CustomFieldsInput{
					Partial: map[string]interface{}{
						"string":    "bbb",
						"new_field": "new",
					},
				},
			},
			map[string]interface{}{
				"int":       int64(3),
				"real":      1.3,
				"string":    "bbb",
				"new_field": "new",
			},
		},
	}
	for _, tt := range tests {
		qb := db.Performer

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			_, err := qb.UpdatePartial(ctx, tt.id, tt.partial)
			if err != nil {
				t.Errorf("PerformerStore.UpdatePartial() error = %v", err)
				return
			}

			// ensure custom fields are correct
			cf, err := qb.GetCustomFields(ctx, tt.id)
			if err != nil {
				t.Errorf("PerformerStore.GetCustomFields() error = %v", err)
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

func TestPerformerFindBySceneID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		pqb := db.Performer
		sceneID := sceneIDs[sceneIdxWithPerformer]

		performers, err := pqb.FindBySceneID(ctx, sceneID)

		if err != nil {
			t.Errorf("Error finding performer: %s", err.Error())
		}

		if !assert.Equal(t, 1, len(performers)) {
			return nil
		}

		performer := performers[0]

		assert.Equal(t, getPerformerStringValue(performerIdxWithScene, "Name"), performer.Name)

		performers, err = pqb.FindBySceneID(ctx, 0)

		if err != nil {
			t.Errorf("Error finding performer: %s", err.Error())
		}

		assert.Equal(t, 0, len(performers))

		return nil
	})
}

func TestPerformerFindByImageID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		pqb := db.Performer
		imageID := imageIDs[imageIdxWithPerformer]

		performers, err := pqb.FindByImageID(ctx, imageID)

		if err != nil {
			t.Errorf("Error finding performer: %s", err.Error())
		}

		if !assert.Equal(t, 1, len(performers)) {
			return nil
		}

		performer := performers[0]

		assert.Equal(t, getPerformerStringValue(performerIdxWithImage, "Name"), performer.Name)

		performers, err = pqb.FindByImageID(ctx, 0)

		if err != nil {
			t.Errorf("Error finding performer: %s", err.Error())
		}

		assert.Equal(t, 0, len(performers))

		return nil
	})
}

func TestPerformerFindByGalleryID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		pqb := db.Performer
		galleryID := galleryIDs[galleryIdxWithPerformer]

		performers, err := pqb.FindByGalleryID(ctx, galleryID)

		if err != nil {
			t.Errorf("Error finding performer: %s", err.Error())
		}

		if !assert.Equal(t, 1, len(performers)) {
			return nil
		}

		performer := performers[0]

		assert.Equal(t, getPerformerStringValue(performerIdxWithGallery, "Name"), performer.Name)

		performers, err = pqb.FindByGalleryID(ctx, 0)

		if err != nil {
			t.Errorf("Error finding performer: %s", err.Error())
		}

		assert.Equal(t, 0, len(performers))

		return nil
	})
}

func TestPerformerFindByNames(t *testing.T) {
	getNames := func(p []*models.Performer) []string {
		var ret []string
		for _, pp := range p {
			ret = append(ret, pp.Name)
		}
		return ret
	}

	withTxn(func(ctx context.Context) error {
		var names []string

		pqb := db.Performer

		names = append(names, performerNames[performerIdxWithScene]) // find performers by names

		performers, err := pqb.FindByNames(ctx, names, false)
		if err != nil {
			t.Errorf("Error finding performers: %s", err.Error())
		}
		assert.Len(t, performers, 1)
		assert.Equal(t, performerNames[performerIdxWithScene], performers[0].Name)

		performers, err = pqb.FindByNames(ctx, names, true) // find performers by names nocase
		if err != nil {
			t.Errorf("Error finding performers: %s", err.Error())
		}
		assert.Len(t, performers, 2) // performerIdxWithScene and performerIdxWithDupName
		assert.Equal(t, strings.ToLower(performerNames[performerIdxWithScene]), strings.ToLower(performers[0].Name))
		assert.Equal(t, strings.ToLower(performerNames[performerIdxWithScene]), strings.ToLower(performers[1].Name))

		names = append(names, performerNames[performerIdx1WithScene]) // find performers by names ( 2 names )

		performers, err = pqb.FindByNames(ctx, names, false)
		if err != nil {
			t.Errorf("Error finding performers: %s", err.Error())
		}
		retNames := getNames(performers)
		assert.Equal(t, names, retNames)

		performers, err = pqb.FindByNames(ctx, names, true) // find performers by names ( 2 names nocase)
		if err != nil {
			t.Errorf("Error finding performers: %s", err.Error())
		}
		retNames = getNames(performers)
		assert.Equal(t, []string{
			performerNames[performerIdxWithScene],
			performerNames[performerIdx1WithScene],
			performerNames[performerIdx1WithDupName],
			performerNames[performerIdxWithDupName],
		}, retNames)

		return nil
	})
}

func TestPerformerQueryEthnicityOr(t *testing.T) {
	const performer1Idx = 1
	const performer2Idx = 2

	performer1Eth := getPerformerStringValue(performer1Idx, "Ethnicity")
	performer2Eth := getPerformerStringValue(performer2Idx, "Ethnicity")

	performerFilter := models.PerformerFilterType{
		Ethnicity: &models.StringCriterionInput{
			Value:    performer1Eth,
			Modifier: models.CriterionModifierEquals,
		},
		OperatorFilter: models.OperatorFilter[models.PerformerFilterType]{
			Or: &models.PerformerFilterType{
				Ethnicity: &models.StringCriterionInput{
					Value:    performer2Eth,
					Modifier: models.CriterionModifierEquals,
				},
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		performers := queryPerformers(ctx, t, &performerFilter, nil)

		assert.Len(t, performers, 2)
		assert.Equal(t, performer1Eth, performers[0].Ethnicity)
		assert.Equal(t, performer2Eth, performers[1].Ethnicity)

		return nil
	})
}

func TestPerformerQueryEthnicityAndRating(t *testing.T) {
	const performerIdx = 1
	performerEth := getPerformerStringValue(performerIdx, "Ethnicity")
	performerRating := int(getRating(performerIdx).Int64)

	performerFilter := models.PerformerFilterType{
		Ethnicity: &models.StringCriterionInput{
			Value:    performerEth,
			Modifier: models.CriterionModifierEquals,
		},
		OperatorFilter: models.OperatorFilter[models.PerformerFilterType]{
			And: &models.PerformerFilterType{
				Rating100: &models.IntCriterionInput{
					Value:    performerRating,
					Modifier: models.CriterionModifierEquals,
				},
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		performers := queryPerformers(ctx, t, &performerFilter, nil)

		if !assert.Len(t, performers, 1) {
			return nil
		}

		assert.Equal(t, performerEth, performers[0].Ethnicity)
		if assert.NotNil(t, performers[0].Rating) {
			assert.Equal(t, performerRating, *performers[0].Rating)
		}

		return nil
	})
}

func TestPerformerQueryEthnicityNotRating(t *testing.T) {
	const performerIdx = 1

	performerRating := getRating(performerIdx)

	ethCriterion := models.StringCriterionInput{
		Value:    "performer_.*1_Ethnicity",
		Modifier: models.CriterionModifierMatchesRegex,
	}

	ratingCriterion := models.IntCriterionInput{
		Value:    int(performerRating.Int64),
		Modifier: models.CriterionModifierEquals,
	}

	performerFilter := models.PerformerFilterType{
		Ethnicity: &ethCriterion,
		OperatorFilter: models.OperatorFilter[models.PerformerFilterType]{
			Not: &models.PerformerFilterType{
				Rating100: &ratingCriterion,
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		performers := queryPerformers(ctx, t, &performerFilter, nil)

		for _, performer := range performers {
			verifyString(t, performer.Ethnicity, ethCriterion)
			ratingCriterion.Modifier = models.CriterionModifierNotEquals
			verifyIntPtr(t, performer.Rating, ratingCriterion)
		}

		return nil
	})
}

func TestPerformerIllegalQuery(t *testing.T) {
	assert := assert.New(t)

	const performerIdx = 1
	subFilter := models.PerformerFilterType{
		Ethnicity: &models.StringCriterionInput{
			Value:    getPerformerStringValue(performerIdx, "Ethnicity"),
			Modifier: models.CriterionModifierEquals,
		},
	}

	tests := []struct {
		name   string
		filter models.PerformerFilterType
	}{
		{
			// And and Or in the same filter
			"AndOr",
			models.PerformerFilterType{
				OperatorFilter: models.OperatorFilter[models.PerformerFilterType]{
					And: &subFilter,
					Or:  &subFilter,
				},
			},
		},
		{
			// And and Not in the same filter
			"AndNot",
			models.PerformerFilterType{
				OperatorFilter: models.OperatorFilter[models.PerformerFilterType]{
					And: &subFilter,
					Not: &subFilter,
				},
			},
		},
		{
			// Or and Not in the same filter
			"OrNot",
			models.PerformerFilterType{
				OperatorFilter: models.OperatorFilter[models.PerformerFilterType]{
					Or:  &subFilter,
					Not: &subFilter,
				},
			},
		},
		{
			"invalid height modifier",
			models.PerformerFilterType{
				Height: &models.StringCriterionInput{
					Modifier: models.CriterionModifierMatchesRegex,
					Value:    "123",
				},
			},
		},
		{
			"invalid height value",
			models.PerformerFilterType{
				Height: &models.StringCriterionInput{
					Modifier: models.CriterionModifierEquals,
					Value:    "foo",
				},
			},
		},
	}

	sqb := db.Performer

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			_, _, err := sqb.Query(ctx, &tt.filter, nil)
			assert.NotNil(err)
		})
	}
}

func TestPerformerQueryIgnoreAutoTag(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		ignoreAutoTag := true
		performerFilter := models.PerformerFilterType{
			IgnoreAutoTag: &ignoreAutoTag,
		}

		performers := queryPerformers(ctx, t, &performerFilter, nil)

		assert.Len(t, performers, int(math.Ceil(float64(totalPerformers)/5)))
		for _, p := range performers {
			assert.True(t, p.IgnoreAutoTag)
		}

		return nil
	})
}

func TestPerformerQuery(t *testing.T) {
	var (
		endpoint = performerStashID(performerIdxWithGallery).Endpoint
		stashID  = performerStashID(performerIdxWithGallery).StashID
	)

	tests := []struct {
		name        string
		findFilter  *models.FindFilterType
		filter      *models.PerformerFilterType
		includeIdxs []int
		excludeIdxs []int
		wantErr     bool
	}{
		{
			"stash id with endpoint",
			nil,
			&models.PerformerFilterType{
				StashIDEndpoint: &models.StashIDCriterionInput{
					Endpoint: &endpoint,
					StashID:  &stashID,
					Modifier: models.CriterionModifierEquals,
				},
			},
			[]int{performerIdxWithGallery},
			nil,
			false,
		},
		{
			"exclude stash id with endpoint",
			nil,
			&models.PerformerFilterType{
				StashIDEndpoint: &models.StashIDCriterionInput{
					Endpoint: &endpoint,
					StashID:  &stashID,
					Modifier: models.CriterionModifierNotEquals,
				},
			},
			nil,
			[]int{performerIdxWithGallery},
			false,
		},
		{
			"null stash id with endpoint",
			nil,
			&models.PerformerFilterType{
				StashIDEndpoint: &models.StashIDCriterionInput{
					Endpoint: &endpoint,
					Modifier: models.CriterionModifierIsNull,
				},
			},
			nil,
			[]int{performerIdxWithGallery},
			false,
		},
		{
			"not null stash id with endpoint",
			nil,
			&models.PerformerFilterType{
				StashIDEndpoint: &models.StashIDCriterionInput{
					Endpoint: &endpoint,
					Modifier: models.CriterionModifierNotNull,
				},
			},
			[]int{performerIdxWithGallery},
			nil,
			false,
		},
		{
			"circumcised (cut)",
			nil,
			&models.PerformerFilterType{
				Circumcised: &models.CircumcisionCriterionInput{
					Value:    []models.CircumisedEnum{models.CircumisedEnumCut},
					Modifier: models.CriterionModifierIncludes,
				},
			},
			[]int{performerIdx1WithScene},
			[]int{performerIdxWithScene, performerIdx2WithScene},
			false,
		},
		{
			"circumcised (excludes cut)",
			nil,
			&models.PerformerFilterType{
				Circumcised: &models.CircumcisionCriterionInput{
					Value:    []models.CircumisedEnum{models.CircumisedEnumCut},
					Modifier: models.CriterionModifierExcludes,
				},
			},
			[]int{performerIdx2WithScene},
			// performerIdxWithScene has null value
			[]int{performerIdx1WithScene, performerIdxWithScene},
			false,
		},
		{
			"include scene studio",
			nil,
			&models.PerformerFilterType{
				Studios: &models.HierarchicalMultiCriterionInput{
					Value:    []string{strconv.Itoa(studioIDs[studioIdxWithScenePerformer])},
					Modifier: models.CriterionModifierIncludes,
				},
			},
			[]int{performerIdxWithSceneStudio},
			nil,
			false,
		},
		{
			"include image studio",
			nil,
			&models.PerformerFilterType{
				Studios: &models.HierarchicalMultiCriterionInput{
					Value:    []string{strconv.Itoa(studioIDs[studioIdxWithImagePerformer])},
					Modifier: models.CriterionModifierIncludes,
				},
			},
			[]int{performerIdxWithImageStudio},
			nil,
			false,
		},
		{
			"include gallery studio",
			nil,
			&models.PerformerFilterType{
				Studios: &models.HierarchicalMultiCriterionInput{
					Value:    []string{strconv.Itoa(studioIDs[studioIdxWithGalleryPerformer])},
					Modifier: models.CriterionModifierIncludes,
				},
			},
			[]int{performerIdxWithGalleryStudio},
			nil,
			false,
		},
		{
			"exclude scene studio",
			nil,
			&models.PerformerFilterType{
				Studios: &models.HierarchicalMultiCriterionInput{
					Value:    []string{strconv.Itoa(studioIDs[studioIdxWithScenePerformer])},
					Modifier: models.CriterionModifierExcludes,
				},
			},
			nil,
			[]int{performerIdxWithSceneStudio},
			false,
		},
		{
			"exclude image studio",
			nil,
			&models.PerformerFilterType{
				Studios: &models.HierarchicalMultiCriterionInput{
					Value:    []string{strconv.Itoa(studioIDs[studioIdxWithImagePerformer])},
					Modifier: models.CriterionModifierExcludes,
				},
			},
			nil,
			[]int{performerIdxWithImageStudio},
			false,
		},
		{
			"exclude gallery studio",
			nil,
			&models.PerformerFilterType{
				Studios: &models.HierarchicalMultiCriterionInput{
					Value:    []string{strconv.Itoa(studioIDs[studioIdxWithGalleryPerformer])},
					Modifier: models.CriterionModifierExcludes,
				},
			},
			nil,
			[]int{performerIdxWithGalleryStudio},
			false,
		},
		{
			"include and exclude scene studio",
			nil,
			&models.PerformerFilterType{
				Studios: &models.HierarchicalMultiCriterionInput{
					Value:    []string{strconv.Itoa(studioIDs[studioIdx1WithTwoScenePerformer])},
					Modifier: models.CriterionModifierIncludes,
					Excludes: []string{strconv.Itoa(studioIDs[studioIdx2WithTwoScenePerformer])},
				},
			},
			nil,
			[]int{performerIdxWithTwoSceneStudio},
			false,
		},
	}

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			performers, _, err := db.Performer.Query(ctx, tt.filter, tt.findFilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("PerformerStore.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ids := performersToIDs(performers)
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

func TestPerformerQueryCustomFields(t *testing.T) {
	tests := []struct {
		name        string
		filter      *models.PerformerFilterType
		includeIdxs []int
		excludeIdxs []int
		wantErr     bool
	}{
		{
			"equals",
			&models.PerformerFilterType{
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierEquals,
						Value:    []any{getPerformerStringValue(performerIdxWithGallery, "custom")},
					},
				},
			},
			[]int{performerIdxWithGallery},
			nil,
			false,
		},
		{
			"not equals",
			&models.PerformerFilterType{
				Name: &models.StringCriterionInput{
					Value:    getPerformerStringValue(performerIdxWithGallery, "Name"),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierNotEquals,
						Value:    []any{getPerformerStringValue(performerIdxWithGallery, "custom")},
					},
				},
			},
			nil,
			[]int{performerIdxWithGallery},
			false,
		},
		{
			"includes",
			&models.PerformerFilterType{
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierIncludes,
						Value:    []any{getPerformerStringValue(performerIdxWithGallery, "custom")[9:]},
					},
				},
			},
			[]int{performerIdxWithGallery},
			nil,
			false,
		},
		{
			"excludes",
			&models.PerformerFilterType{
				Name: &models.StringCriterionInput{
					Value:    getPerformerStringValue(performerIdxWithGallery, "Name"),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierExcludes,
						Value:    []any{getPerformerStringValue(performerIdxWithGallery, "custom")[9:]},
					},
				},
			},
			nil,
			[]int{performerIdxWithGallery},
			false,
		},
		{
			"regex",
			&models.PerformerFilterType{
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierMatchesRegex,
						Value:    []any{".*13_custom"},
					},
				},
			},
			[]int{performerIdxWithGallery},
			nil,
			false,
		},
		{
			"invalid regex",
			&models.PerformerFilterType{
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
			&models.PerformerFilterType{
				Name: &models.StringCriterionInput{
					Value:    getPerformerStringValue(performerIdxWithGallery, "Name"),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierNotMatchesRegex,
						Value:    []any{".*13_custom"},
					},
				},
			},
			nil,
			[]int{performerIdxWithGallery},
			false,
		},
		{
			"invalid not matches regex",
			&models.PerformerFilterType{
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
			&models.PerformerFilterType{
				Name: &models.StringCriterionInput{
					Value:    getPerformerStringValue(performerIdxWithGallery, "Name"),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "not existing",
						Modifier: models.CriterionModifierIsNull,
					},
				},
			},
			[]int{performerIdxWithGallery},
			nil,
			false,
		},
		{
			"null",
			&models.PerformerFilterType{
				Name: &models.StringCriterionInput{
					Value:    getPerformerStringValue(performerIdxWithGallery, "Name"),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierNotNull,
					},
				},
			},
			[]int{performerIdxWithGallery},
			nil,
			false,
		},
		{
			"between",
			&models.PerformerFilterType{
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "real",
						Modifier: models.CriterionModifierBetween,
						Value:    []any{0.05, 0.15},
					},
				},
			},
			[]int{performerIdx1WithScene},
			nil,
			false,
		},
		{
			"not between",
			&models.PerformerFilterType{
				Name: &models.StringCriterionInput{
					Value:    getPerformerStringValue(performerIdx1WithScene, "Name"),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "real",
						Modifier: models.CriterionModifierNotBetween,
						Value:    []any{0.05, 0.15},
					},
				},
			},
			nil,
			[]int{performerIdx1WithScene},
			false,
		},
	}

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			performers, _, err := db.Performer.Query(ctx, tt.filter, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("PerformerStore.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ids := performersToIDs(performers)
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

func TestPerformerQueryPenisLength(t *testing.T) {
	var upper = 4.0

	tests := []struct {
		name     string
		modifier models.CriterionModifier
		value    float64
		value2   *float64
	}{
		{
			"equals",
			models.CriterionModifierEquals,
			1,
			nil,
		},
		{
			"not equals",
			models.CriterionModifierNotEquals,
			1,
			nil,
		},
		{
			"greater than",
			models.CriterionModifierGreaterThan,
			1,
			nil,
		},
		{
			"between",
			models.CriterionModifierBetween,
			2,
			&upper,
		},
		{
			"greater than",
			models.CriterionModifierNotBetween,
			2,
			&upper,
		},
		{
			"null",
			models.CriterionModifierIsNull,
			0,
			nil,
		},
		{
			"not null",
			models.CriterionModifierNotNull,
			0,
			nil,
		},
	}

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			filter := &models.PerformerFilterType{
				PenisLength: &models.FloatCriterionInput{
					Modifier: tt.modifier,
					Value:    tt.value,
					Value2:   tt.value2,
				},
			}

			performers, _, err := db.Performer.Query(ctx, filter, nil)
			if err != nil {
				t.Errorf("PerformerStore.Query() error = %v", err)
				return
			}

			for _, p := range performers {
				verifyFloat(t, p.PenisLength, *filter.PenisLength)
			}
		})
	}
}

func verifyFloat(t *testing.T, value *float64, criterion models.FloatCriterionInput) bool {
	t.Helper()
	assert := assert.New(t)
	switch criterion.Modifier {
	case models.CriterionModifierEquals:
		return assert.NotNil(value) && assert.Equal(criterion.Value, *value)
	case models.CriterionModifierNotEquals:
		return assert.NotNil(value) && assert.NotEqual(criterion.Value, *value)
	case models.CriterionModifierGreaterThan:
		return assert.NotNil(value) && assert.Greater(*value, criterion.Value)
	case models.CriterionModifierLessThan:
		return assert.NotNil(value) && assert.Less(*value, criterion.Value)
	case models.CriterionModifierBetween:
		return assert.NotNil(value) && assert.GreaterOrEqual(*value, criterion.Value) && assert.LessOrEqual(*value, *criterion.Value2)
	case models.CriterionModifierNotBetween:
		return assert.NotNil(value) && assert.True(*value < criterion.Value || *value > *criterion.Value2)
	case models.CriterionModifierIsNull:
		return assert.Nil(value)
	case models.CriterionModifierNotNull:
		return assert.NotNil(value)
	}

	return false
}

func TestPerformerQueryForAutoTag(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		tqb := db.Performer

		name := performerNames[performerIdx1WithScene] // find a performer by name

		performers, err := tqb.QueryForAutoTag(ctx, []string{name})

		if err != nil {
			t.Errorf("Error finding performers: %s", err.Error())
		}

		assert.Len(t, performers, 2)
		assert.Equal(t, strings.ToLower(performerNames[performerIdx1WithScene]), strings.ToLower(performers[0].Name))
		assert.Equal(t, strings.ToLower(performerNames[performerIdx1WithScene]), strings.ToLower(performers[1].Name))

		return nil
	})
}

func TestPerformerUpdatePerformerImage(t *testing.T) {
	if err := withRollbackTxn(func(ctx context.Context) error {
		qb := db.Performer

		// create performer to test against
		const name = "TestPerformerUpdatePerformerImage"
		performer := models.Performer{
			Name: name,
		}
		err := qb.Create(ctx, &models.CreatePerformerInput{Performer: &performer})
		if err != nil {
			return fmt.Errorf("Error creating performer: %s", err.Error())
		}

		return testUpdateImage(t, ctx, performer.ID, qb.UpdateImage, qb.GetImage)
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestPerformerQueryAge(t *testing.T) {
	const age = 19
	ageCriterion := models.IntCriterionInput{
		Value:    age,
		Modifier: models.CriterionModifierEquals,
	}

	verifyPerformerAge(t, ageCriterion)

	ageCriterion.Modifier = models.CriterionModifierNotEquals
	verifyPerformerAge(t, ageCriterion)

	ageCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyPerformerAge(t, ageCriterion)

	ageCriterion.Modifier = models.CriterionModifierLessThan
	verifyPerformerAge(t, ageCriterion)
}

func verifyPerformerAge(t *testing.T, ageCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		qb := db.Performer
		performerFilter := models.PerformerFilterType{
			Age: &ageCriterion,
		}

		performers, _, err := qb.Query(ctx, &performerFilter, nil)
		if err != nil {
			t.Errorf("Error querying performer: %s", err.Error())
		}

		now := time.Now()
		for _, performer := range performers {
			cd := now

			if performer.DeathDate != nil {
				cd = performer.DeathDate.Time
			}

			d := performer.Birthdate.Time
			age := cd.Year() - d.Year()
			// using YearDay screws up on leap years
			if cd.Month() < d.Month() || (cd.Month() == d.Month() && cd.Day() < d.Day()) {
				age = age - 1
			}

			if !verifyInt(t, age, ageCriterion) {
				t.Errorf("Performer birthdate: %s, deathdate: %s", performer.Birthdate.String(), performer.DeathDate.String())
			}
		}

		return nil
	})
}

func TestPerformerQueryCareerLength(t *testing.T) {
	const value = "2005"
	careerLengthCriterion := models.StringCriterionInput{
		Value:    value,
		Modifier: models.CriterionModifierEquals,
	}

	verifyPerformerCareerLength(t, careerLengthCriterion)

	careerLengthCriterion.Modifier = models.CriterionModifierNotEquals
	verifyPerformerCareerLength(t, careerLengthCriterion)

	careerLengthCriterion.Modifier = models.CriterionModifierMatchesRegex
	verifyPerformerCareerLength(t, careerLengthCriterion)

	careerLengthCriterion.Modifier = models.CriterionModifierNotMatchesRegex
	verifyPerformerCareerLength(t, careerLengthCriterion)
}

func verifyPerformerCareerLength(t *testing.T, criterion models.StringCriterionInput) {
	withTxn(func(ctx context.Context) error {
		qb := db.Performer
		performerFilter := models.PerformerFilterType{
			CareerLength: &criterion,
		}

		performers, _, err := qb.Query(ctx, &performerFilter, nil)
		if err != nil {
			t.Errorf("Error querying performer: %s", err.Error())
		}

		for _, performer := range performers {
			cl := performer.CareerLength
			verifyString(t, cl, criterion)
		}

		return nil
	})
}

func TestPerformerQueryURL(t *testing.T) {
	const sceneIdx = 1
	performerURL := getPerformerStringValue(sceneIdx, urlField)

	urlCriterion := models.StringCriterionInput{
		Value:    performerURL,
		Modifier: models.CriterionModifierEquals,
	}

	filter := models.PerformerFilterType{
		URL: &urlCriterion,
	}

	verifyFn := func(g *models.Performer) {
		t.Helper()

		urls := g.URLs.List()
		var url string
		if len(urls) > 0 {
			url = urls[0]
		}

		verifyString(t, url, urlCriterion)
	}

	verifyPerformerQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotEquals
	verifyPerformerQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierMatchesRegex
	urlCriterion.Value = "performer_.*1_URL"
	verifyPerformerQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotMatchesRegex
	verifyPerformerQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierIsNull
	urlCriterion.Value = ""
	verifyPerformerQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotNull
	verifyPerformerQuery(t, filter, verifyFn)
}

func verifyPerformerQuery(t *testing.T, filter models.PerformerFilterType, verifyFn func(s *models.Performer)) {
	withTxn(func(ctx context.Context) error {
		t.Helper()
		performers := queryPerformers(ctx, t, &filter, nil)

		for _, performer := range performers {
			if err := performer.LoadURLs(ctx, db.Performer); err != nil {
				t.Errorf("Error loading url relationships: %v", err)
			}
		}

		// assume it should find at least one
		assert.Greater(t, len(performers), 0)

		for _, p := range performers {
			verifyFn(p)
		}

		return nil
	})
}

func queryPerformers(ctx context.Context, t *testing.T, performerFilter *models.PerformerFilterType, findFilter *models.FindFilterType) []*models.Performer {
	t.Helper()
	performers, _, err := db.Performer.Query(ctx, performerFilter, findFilter)
	if err != nil {
		t.Errorf("Error querying performers: %s", err.Error())
	}

	return performers
}

func TestPerformerQueryTags(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		tagCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdxWithPerformer]),
				strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		performerFilter := models.PerformerFilterType{
			Tags: &tagCriterion,
		}

		// ensure ids are correct
		performers := queryPerformers(ctx, t, &performerFilter, nil)
		assert.Len(t, performers, 2)
		for _, performer := range performers {
			assert.True(t, performer.ID == performerIDs[performerIdxWithTag] || performer.ID == performerIDs[performerIdxWithTwoTags])
		}

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
				strconv.Itoa(tagIDs[tagIdx2WithPerformer]),
			},
			Modifier: models.CriterionModifierIncludesAll,
		}

		performers = queryPerformers(ctx, t, &performerFilter, nil)

		assert.Len(t, performers, 1)
		assert.Equal(t, sceneIDs[performerIdxWithTwoTags], performers[0].ID)

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getSceneStringValue(performerIdxWithTwoTags, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		performers = queryPerformers(ctx, t, &performerFilter, &findFilter)
		assert.Len(t, performers, 0)

		return nil
	})
}

func TestPerformerQueryTagCount(t *testing.T) {
	const tagCount = 1
	tagCountCriterion := models.IntCriterionInput{
		Value:    tagCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyPerformersTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyPerformersTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyPerformersTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyPerformersTagCount(t, tagCountCriterion)
}

func verifyPerformersTagCount(t *testing.T, tagCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Performer
		performerFilter := models.PerformerFilterType{
			TagCount: &tagCountCriterion,
		}

		performers := queryPerformers(ctx, t, &performerFilter, nil)
		assert.Greater(t, len(performers), 0)

		for _, performer := range performers {
			ids, err := sqb.GetTagIDs(ctx, performer.ID)
			if err != nil {
				return err
			}
			verifyInt(t, len(ids), tagCountCriterion)
		}

		return nil
	})
}

func TestPerformerQuerySceneCount(t *testing.T) {
	const sceneCount = 1
	sceneCountCriterion := models.IntCriterionInput{
		Value:    sceneCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyPerformersSceneCount(t, sceneCountCriterion)

	sceneCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyPerformersSceneCount(t, sceneCountCriterion)

	sceneCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyPerformersSceneCount(t, sceneCountCriterion)

	sceneCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyPerformersSceneCount(t, sceneCountCriterion)
}

func verifyPerformersSceneCount(t *testing.T, sceneCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		performerFilter := models.PerformerFilterType{
			SceneCount: &sceneCountCriterion,
		}

		performers := queryPerformers(ctx, t, &performerFilter, nil)
		assert.Greater(t, len(performers), 0)

		for _, performer := range performers {
			ids, err := db.Scene.FindByPerformerID(ctx, performer.ID)
			if err != nil {
				return err
			}
			verifyInt(t, len(ids), sceneCountCriterion)
		}

		return nil
	})
}

func TestPerformerQueryImageCount(t *testing.T) {
	const imageCount = 1
	imageCountCriterion := models.IntCriterionInput{
		Value:    imageCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyPerformersImageCount(t, imageCountCriterion)

	imageCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyPerformersImageCount(t, imageCountCriterion)

	imageCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyPerformersImageCount(t, imageCountCriterion)

	imageCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyPerformersImageCount(t, imageCountCriterion)
}

func verifyPerformersImageCount(t *testing.T, imageCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		performerFilter := models.PerformerFilterType{
			ImageCount: &imageCountCriterion,
		}

		performers := queryPerformers(ctx, t, &performerFilter, nil)
		assert.Greater(t, len(performers), 0)

		for _, performer := range performers {
			pp := 0

			result, err := db.Image.Query(ctx, models.ImageQueryOptions{
				QueryOptions: models.QueryOptions{
					FindFilter: &models.FindFilterType{
						PerPage: &pp,
					},
					Count: true,
				},
				ImageFilter: &models.ImageFilterType{
					Performers: &models.MultiCriterionInput{
						Value:    []string{strconv.Itoa(performer.ID)},
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

func TestPerformerQueryGalleryCount(t *testing.T) {
	const galleryCount = 1
	galleryCountCriterion := models.IntCriterionInput{
		Value:    galleryCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyPerformersGalleryCount(t, galleryCountCriterion)

	galleryCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyPerformersGalleryCount(t, galleryCountCriterion)

	galleryCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyPerformersGalleryCount(t, galleryCountCriterion)

	galleryCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyPerformersGalleryCount(t, galleryCountCriterion)
}

func verifyPerformersGalleryCount(t *testing.T, galleryCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		performerFilter := models.PerformerFilterType{
			GalleryCount: &galleryCountCriterion,
		}

		performers := queryPerformers(ctx, t, &performerFilter, nil)
		assert.Greater(t, len(performers), 0)

		for _, performer := range performers {
			pp := 0

			_, count, err := db.Gallery.Query(ctx, &models.GalleryFilterType{
				Performers: &models.MultiCriterionInput{
					Value:    []string{strconv.Itoa(performer.ID)},
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

func TestPerformerQueryStudio(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		testCases := []struct {
			studioIndex    int
			performerIndex int
		}{
			{studioIndex: studioIdxWithScenePerformer, performerIndex: performerIdxWithSceneStudio},
			{studioIndex: studioIdxWithImagePerformer, performerIndex: performerIdxWithImageStudio},
			{studioIndex: studioIdxWithGalleryPerformer, performerIndex: performerIdxWithGalleryStudio},
		}

		for _, tc := range testCases {
			studioCriterion := models.HierarchicalMultiCriterionInput{
				Value: []string{
					strconv.Itoa(studioIDs[tc.studioIndex]),
				},
				Modifier: models.CriterionModifierIncludes,
			}

			performerFilter := models.PerformerFilterType{
				Studios: &studioCriterion,
			}

			performers := queryPerformers(ctx, t, &performerFilter, nil)

			assert.Len(t, performers, 1)

			// ensure id is correct
			assert.Equal(t, performerIDs[tc.performerIndex], performers[0].ID)

			studioCriterion = models.HierarchicalMultiCriterionInput{
				Value: []string{
					strconv.Itoa(studioIDs[tc.studioIndex]),
				},
				Modifier: models.CriterionModifierExcludes,
			}

			q := getPerformerStringValue(tc.performerIndex, "Name")
			findFilter := models.FindFilterType{
				Q: &q,
			}

			performers = queryPerformers(ctx, t, &performerFilter, &findFilter)
			assert.Len(t, performers, 0)
		}

		// test NULL/not NULL
		q := getPerformerStringValue(performerIdx1WithImage, "Name")
		performerFilter := &models.PerformerFilterType{
			Studios: &models.HierarchicalMultiCriterionInput{
				Modifier: models.CriterionModifierIsNull,
			},
		}
		findFilter := &models.FindFilterType{
			Q: &q,
		}

		performers := queryPerformers(ctx, t, performerFilter, findFilter)
		assert.Len(t, performers, 1)
		assert.Equal(t, imageIDs[performerIdx1WithImage], performers[0].ID)

		q = getPerformerStringValue(performerIdxWithSceneStudio, "Name")
		performers = queryPerformers(ctx, t, performerFilter, findFilter)
		assert.Len(t, performers, 0)

		performerFilter.Studios.Modifier = models.CriterionModifierNotNull
		performers = queryPerformers(ctx, t, performerFilter, findFilter)
		assert.Len(t, performers, 1)
		assert.Equal(t, imageIDs[performerIdxWithSceneStudio], performers[0].ID)

		q = getPerformerStringValue(performerIdx1WithImage, "Name")
		performers = queryPerformers(ctx, t, performerFilter, findFilter)
		assert.Len(t, performers, 0)

		return nil
	})
}

func TestPerformerStashIDs(t *testing.T) {
	if err := withRollbackTxn(func(ctx context.Context) error {
		qb := db.Performer

		// create scene to test against
		const name = "TestPerformerStashIDs"
		performer := &models.Performer{
			Name: name,
		}
		if err := qb.Create(ctx, &models.CreatePerformerInput{Performer: performer}); err != nil {
			return fmt.Errorf("Error creating performer: %s", err.Error())
		}

		if err := performer.LoadStashIDs(ctx, qb); err != nil {
			return err
		}

		testPerformerStashIDs(ctx, t, performer)
		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func testPerformerStashIDs(ctx context.Context, t *testing.T, s *models.Performer) {
	// ensure no stash IDs to begin with
	assert.Len(t, s.StashIDs.List(), 0)

	// add stash ids
	const stashIDStr = "stashID"
	const endpoint = "endpoint"
	stashID := models.StashID{
		StashID:   stashIDStr,
		Endpoint:  endpoint,
		UpdatedAt: epochTime,
	}

	qb := db.Performer

	// update stash ids and ensure was updated
	var err error
	s, err = qb.UpdatePartial(ctx, s.ID, models.PerformerPartial{
		StashIDs: &models.UpdateStashIDs{
			StashIDs: []models.StashID{stashID},
			Mode:     models.RelationshipUpdateModeSet,
		},
	})
	if err != nil {
		t.Error(err.Error())
	}

	if err := s.LoadStashIDs(ctx, qb); err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, []models.StashID{stashID}, s.StashIDs.List())

	// remove stash ids and ensure was updated
	s, err = qb.UpdatePartial(ctx, s.ID, models.PerformerPartial{
		StashIDs: &models.UpdateStashIDs{
			StashIDs: []models.StashID{stashID},
			Mode:     models.RelationshipUpdateModeRemove,
		},
	})
	if err != nil {
		t.Error(err.Error())
	}

	if err := s.LoadStashIDs(ctx, qb); err != nil {
		t.Error(err.Error())
		return
	}

	assert.Len(t, s.StashIDs.List(), 0)
}

func TestPerformerQueryRating100(t *testing.T) {
	const rating = 60
	ratingCriterion := models.IntCriterionInput{
		Value:    rating,
		Modifier: models.CriterionModifierEquals,
	}

	verifyPerformersRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotEquals
	verifyPerformersRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyPerformersRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierLessThan
	verifyPerformersRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierIsNull
	verifyPerformersRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotNull
	verifyPerformersRating100(t, ratingCriterion)
}

func verifyPerformersRating100(t *testing.T, ratingCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		performerFilter := models.PerformerFilterType{
			Rating100: &ratingCriterion,
		}

		performers := queryPerformers(ctx, t, &performerFilter, nil)

		for _, performer := range performers {
			verifyIntPtr(t, performer.Rating, ratingCriterion)
		}

		return nil
	})
}

func performerQueryIsMissing(ctx context.Context, t *testing.T, m string) []*models.Performer {
	performerFilter := models.PerformerFilterType{
		IsMissing: &m,
	}

	return queryPerformers(ctx, t, &performerFilter, nil)
}

func TestPerformerQueryIsMissingRating(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		performers := performerQueryIsMissing(ctx, t, "rating")

		assert.True(t, len(performers) > 0)

		for _, performer := range performers {
			assert.Nil(t, performer.Rating)
		}

		return nil
	})
}

func TestPerformerQueryIsMissingImage(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		performers := performerQueryIsMissing(ctx, t, "image")

		assert.True(t, len(performers) > 0)

		for _, performer := range performers {
			img, err := db.Performer.GetImage(ctx, performer.ID)
			if err != nil {
				t.Errorf("error getting performer image: %s", err.Error())
			}
			assert.Nil(t, img)
		}

		return nil
	})
}

func TestPerformerQueryIsMissingAlias(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		performers := performerQueryIsMissing(ctx, t, "aliases")

		assert.True(t, len(performers) > 0)

		for _, performer := range performers {
			a, err := db.Performer.GetAliases(ctx, performer.ID)
			if err != nil {
				t.Errorf("error getting performer aliases: %s", err.Error())
			}
			assert.Nil(t, a)
		}

		return nil
	})
}

func TestPerformerQuerySortScenesCount(t *testing.T) {
	sort := "scenes_count"
	direction := models.SortDirectionEnumDesc
	findFilter := &models.FindFilterType{
		Sort:      &sort,
		Direction: &direction,
	}

	withTxn(func(ctx context.Context) error {
		// just ensure it queries without error
		performers, _, err := db.Performer.Query(ctx, nil, findFilter)
		if err != nil {
			t.Errorf("Error querying performers: %s", err.Error())
		}

		assert.True(t, len(performers) > 0)

		// first performer should be performerIdx1WithScene
		firstPerformer := performers[0]

		assert.Equal(t, performerIDs[performerIdx1WithScene], firstPerformer.ID)

		// sort in ascending order
		direction = models.SortDirectionEnumAsc

		performers, _, err = db.Performer.Query(ctx, nil, findFilter)
		if err != nil {
			t.Errorf("Error querying performers: %s", err.Error())
		}

		assert.True(t, len(performers) > 0)
		lastPerformer := performers[len(performers)-1]

		assert.Equal(t, performerIDs[performerIdxWithTwoSceneStudio], lastPerformer.ID)

		return nil
	})
}

func TestPerformerCountByTagID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Performer
		count, err := sqb.CountByTagID(ctx, tagIDs[tagIdxWithPerformer])

		if err != nil {
			t.Errorf("Error counting performers: %s", err.Error())
		}

		assert.Equal(t, 1, count)

		count, err = sqb.CountByTagID(ctx, 0)

		if err != nil {
			t.Errorf("Error counting performers: %s", err.Error())
		}

		assert.Equal(t, 0, count)

		return nil
	})
}

func TestPerformerCount(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Performer
		count, err := sqb.Count(ctx)

		if err != nil {
			t.Errorf("Error counting performers: %s", err.Error())
		}

		assert.Equal(t, totalPerformers, count)

		return nil
	})
}

func TestPerformerAll(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Performer
		all, err := sqb.All(ctx)

		if err != nil {
			t.Errorf("Error counting performers: %s", err.Error())
		}

		assert.Len(t, all, totalPerformers)

		return nil
	})
}

func performersToIDs(i []*models.Performer) []int {
	ret := make([]int, len(i))
	for i, v := range i {
		ret[i] = v.ID
	}

	return ret
}

func TestPerformerStore_FindByStashID(t *testing.T) {
	type args struct {
		stashID models.StashID
	}
	tests := []struct {
		name        string
		stashID     models.StashID
		expectedIDs []int
		wantErr     bool
	}{
		{
			name:        "existing",
			stashID:     performerStashID(performerIdxWithScene),
			expectedIDs: []int{performerIDs[performerIdxWithScene]},
			wantErr:     false,
		},
		{
			name: "non-existing",
			stashID: models.StashID{
				StashID:  getPerformerStringValue(performerIdxWithScene, "stashid"),
				Endpoint: "non-existing",
			},
			expectedIDs: []int{},
			wantErr:     false,
		},
	}

	qb := db.Performer

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.FindByStashID(ctx, tt.stashID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PerformerStore.FindByStashID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.ElementsMatch(t, performersToIDs(got), tt.expectedIDs)
		})
	}
}

func TestPerformerStore_FindByStashIDStatus(t *testing.T) {
	type args struct {
		stashID models.StashID
	}
	tests := []struct {
		name             string
		hasStashID       bool
		stashboxEndpoint string
		include          []int
		exclude          []int
		wantErr          bool
	}{
		{
			name:             "existing",
			hasStashID:       true,
			stashboxEndpoint: getPerformerStringValue(performerIdxWithScene, "endpoint"),
			include:          []int{performerIdxWithScene},
			wantErr:          false,
		},
		{
			name:             "non-existing",
			hasStashID:       true,
			stashboxEndpoint: getPerformerStringValue(performerIdxWithScene, "non-existing"),
			exclude:          []int{performerIdxWithScene},
			wantErr:          false,
		},
		{
			name:             "!hasStashID",
			hasStashID:       false,
			stashboxEndpoint: getPerformerStringValue(performerIdxWithScene, "endpoint"),
			include:          []int{performerIdxWithTwoScenes},
			exclude:          []int{performerIdx2WithScene},
			wantErr:          false,
		},
	}

	qb := db.Performer

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.FindByStashIDStatus(ctx, tt.hasStashID, tt.stashboxEndpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("PerformerStore.FindByStashIDStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			include := indexesToIDs(performerIDs, tt.include)
			exclude := indexesToIDs(performerIDs, tt.exclude)

			ids := performersToIDs(got)

			assert := assert.New(t)
			for _, i := range include {
				assert.Contains(ids, i)
			}
			for _, e := range exclude {
				assert.NotContains(ids, e)
			}
		})
	}
}

func TestPerformerMerge(t *testing.T) {
	tests := []struct {
		name    string
		srcIdxs []int
		destIdx int
		wantErr bool
	}{
		{
			name:    "merge into self",
			srcIdxs: []int{performerIdx1WithDupName},
			destIdx: performerIdx1WithDupName,
			wantErr: true,
		},
		{
			name: "merge multiple",
			srcIdxs: []int{
				performerIdx2WithScene,
				performerIdxWithTwoScenes,
				performerIdx1WithImage,
				performerIdxWithTwoImages,
				performerIdxWithGallery,
				performerIdxWithTwoGalleries,
				performerIdxWithTag,
				performerIdxWithTwoTags,
			},
			destIdx: tagIdxWithPerformer,
			wantErr: false,
		},
	}

	qb := db.Performer

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			// load src tag ids to compare after merge
			performerTagIds := make(map[int][]int)
			for _, srcIdx := range tt.srcIdxs {
				srcPerformer, err := qb.Find(ctx, performerIDs[srcIdx])
				if err != nil {
					t.Errorf("Error finding performer: %s", err.Error())
				}
				if err := srcPerformer.LoadTagIDs(ctx, qb); err != nil {
					t.Errorf("Error loading performer tag IDs: %s", err.Error())
				}
				srcTagIDs := srcPerformer.TagIDs.List()
				performerTagIds[srcIdx] = srcTagIDs
			}

			err := qb.Merge(ctx, indexesToIDs(tagIDs, tt.srcIdxs), tagIDs[tt.destIdx])

			if (err != nil) != tt.wantErr {
				t.Errorf("PerformerStore.Merge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// ensure source performers are destroyed
			for _, srcIdx := range tt.srcIdxs {
				p, err := qb.Find(ctx, performerIDs[srcIdx])

				// not found returns nil performer and nil error
				if err != nil {
					t.Errorf("Error finding performer: %s", err.Error())
					continue
				}
				assert.Nil(p)
			}

			// ensure items point to new performer
			for _, srcIdx := range tt.srcIdxs {
				sceneIdxs := scenePerformers.reverseLookup(srcIdx)
				for _, sceneIdx := range sceneIdxs {
					s, err := db.Scene.Find(ctx, sceneIDs[sceneIdx])
					if err != nil {
						t.Errorf("Error finding scene: %s", err.Error())
					}
					if err := s.LoadPerformerIDs(ctx, db.Scene); err != nil {
						t.Errorf("Error loading scene performer IDs: %s", err.Error())
					}
					scenePerformerIDs := s.PerformerIDs.List()

					assert.Contains(scenePerformerIDs, performerIDs[tt.destIdx])
					assert.NotContains(scenePerformerIDs, performerIDs[srcIdx])
				}

				imageIdxs := imagePerformers.reverseLookup(srcIdx)
				for _, imageIdx := range imageIdxs {
					i, err := db.Image.Find(ctx, imageIDs[imageIdx])
					if err != nil {
						t.Errorf("Error finding image: %s", err.Error())
					}
					if err := i.LoadPerformerIDs(ctx, db.Image); err != nil {
						t.Errorf("Error loading image performer IDs: %s", err.Error())
					}
					imagePerformerIDs := i.PerformerIDs.List()

					assert.Contains(imagePerformerIDs, performerIDs[tt.destIdx])
					assert.NotContains(imagePerformerIDs, performerIDs[srcIdx])
				}

				galleryIdxs := galleryPerformers.reverseLookup(srcIdx)
				for _, galleryIdx := range galleryIdxs {
					g, err := db.Gallery.Find(ctx, galleryIDs[galleryIdx])
					if err != nil {
						t.Errorf("Error finding gallery: %s", err.Error())
					}
					if err := g.LoadPerformerIDs(ctx, db.Gallery); err != nil {
						t.Errorf("Error loading gallery performer IDs: %s", err.Error())
					}
					galleryPerformerIDs := g.PerformerIDs.List()

					assert.Contains(galleryPerformerIDs, performerIDs[tt.destIdx])
					assert.NotContains(galleryPerformerIDs, performerIDs[srcIdx])
				}
			}

			// ensure tags were merged
			destPerformer, err := qb.Find(ctx, performerIDs[tt.destIdx])
			if err != nil {
				t.Errorf("Error finding performer: %s", err.Error())
			}
			if err := destPerformer.LoadTagIDs(ctx, qb); err != nil {
				t.Errorf("Error loading performer tag IDs: %s", err.Error())
			}
			destTagIDs := destPerformer.TagIDs.List()

			for _, srcIdx := range tt.srcIdxs {
				for _, tagID := range performerTagIds[srcIdx] {
					assert.Contains(destTagIDs, tagID)
				}
			}
		})
	}
}

// TODO Update
// TODO Destroy
// TODO Find
// TODO Query
