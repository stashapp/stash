//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"fmt"
	"math"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stretchr/testify/assert"
)

func loadSceneRelationships(ctx context.Context, expected models.Scene, actual *models.Scene) error {
	if expected.URLs.Loaded() {
		if err := actual.LoadURLs(ctx, db.Scene); err != nil {
			return err
		}
	}

	if expected.GalleryIDs.Loaded() {
		if err := actual.LoadGalleryIDs(ctx, db.Scene); err != nil {
			return err
		}
	}
	if expected.TagIDs.Loaded() {
		if err := actual.LoadTagIDs(ctx, db.Scene); err != nil {
			return err
		}
	}
	if expected.PerformerIDs.Loaded() {
		if err := actual.LoadPerformerIDs(ctx, db.Scene); err != nil {
			return err
		}
	}
	if expected.Groups.Loaded() {
		if err := actual.LoadGroups(ctx, db.Scene); err != nil {
			return err
		}
	}
	if expected.StashIDs.Loaded() {
		if err := actual.LoadStashIDs(ctx, db.Scene); err != nil {
			return err
		}
	}
	if expected.Files.Loaded() {
		if err := actual.LoadFiles(ctx, db.Scene); err != nil {
			return err
		}
	}

	// clear Path, Checksum, PrimaryFileID
	if expected.Path == "" {
		actual.Path = ""
	}
	if expected.Checksum == "" {
		actual.Checksum = ""
	}
	if expected.OSHash == "" {
		actual.OSHash = ""
	}
	if expected.PrimaryFileID == nil {
		actual.PrimaryFileID = nil
	}

	return nil
}

func Test_sceneQueryBuilder_Create(t *testing.T) {
	var (
		title        = "title"
		code         = "1337"
		details      = "details"
		director     = "director"
		url          = "url"
		rating       = 60
		resumeTime   = 10.0
		playDuration = 34.0
		createdAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		sceneIndex   = 123
		sceneIndex2  = 234
		endpoint1    = "endpoint1"
		endpoint2    = "endpoint2"
		stashID1     = "stashid1"
		stashID2     = "stashid2"

		date, _ = models.ParseDate("2003-02-01")

		videoFile = makeFileWithID(fileIdxStartVideoFiles)
	)

	tests := []struct {
		name      string
		newObject models.Scene
		wantErr   bool
	}{
		{
			"full",
			models.Scene{
				Title:        title,
				Code:         code,
				Details:      details,
				Director:     director,
				URLs:         models.NewRelatedStrings([]string{url}),
				Date:         &date,
				Rating:       &rating,
				Organized:    true,
				StudioID:     &studioIDs[studioIdxWithScene],
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				GalleryIDs:   models.NewRelatedIDs([]int{galleryIDs[galleryIdxWithScene]}),
				TagIDs:       models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithScene]}),
				PerformerIDs: models.NewRelatedIDs([]int{performerIDs[performerIdx1WithScene], performerIDs[performerIdx1WithDupName]}),
				Groups: models.NewRelatedGroups([]models.GroupsScenes{
					{
						GroupID:    groupIDs[groupIdxWithScene],
						SceneIndex: &sceneIndex,
					},
					{
						GroupID:    groupIDs[groupIdxWithStudio],
						SceneIndex: &sceneIndex2,
					},
				}),
				StashIDs: models.NewRelatedStashIDs([]models.StashID{
					{
						StashID:  stashID1,
						Endpoint: endpoint1,
					},
					{
						StashID:  stashID2,
						Endpoint: endpoint2,
					},
				}),
				ResumeTime:   float64(resumeTime),
				PlayDuration: playDuration,
			},
			false,
		},
		{
			"with file",
			models.Scene{
				Title:     title,
				Code:      code,
				Details:   details,
				Director:  director,
				URLs:      models.NewRelatedStrings([]string{url}),
				Date:      &date,
				Rating:    &rating,
				Organized: true,
				StudioID:  &studioIDs[studioIdxWithScene],
				Files: models.NewRelatedVideoFiles([]*models.VideoFile{
					videoFile.(*models.VideoFile),
				}),
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				GalleryIDs:   models.NewRelatedIDs([]int{galleryIDs[galleryIdxWithScene]}),
				TagIDs:       models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithScene]}),
				PerformerIDs: models.NewRelatedIDs([]int{performerIDs[performerIdx1WithScene], performerIDs[performerIdx1WithDupName]}),
				Groups: models.NewRelatedGroups([]models.GroupsScenes{
					{
						GroupID:    groupIDs[groupIdxWithScene],
						SceneIndex: &sceneIndex,
					},
					{
						GroupID:    groupIDs[groupIdxWithStudio],
						SceneIndex: &sceneIndex2,
					},
				}),
				StashIDs: models.NewRelatedStashIDs([]models.StashID{
					{
						StashID:  stashID1,
						Endpoint: endpoint1,
					},
					{
						StashID:  stashID2,
						Endpoint: endpoint2,
					},
				}),
				ResumeTime:   resumeTime,
				PlayDuration: playDuration,
			},
			false,
		},
		{
			"invalid studio id",
			models.Scene{
				StudioID: &invalidID,
			},
			true,
		},
		{
			"invalid gallery id",
			models.Scene{
				GalleryIDs: models.NewRelatedIDs([]int{invalidID}),
			},
			true,
		},
		{
			"invalid tag id",
			models.Scene{
				TagIDs: models.NewRelatedIDs([]int{invalidID}),
			},
			true,
		},
		{
			"invalid performer id",
			models.Scene{
				PerformerIDs: models.NewRelatedIDs([]int{invalidID}),
			},
			true,
		},
		{
			"invalid group id",
			models.Scene{
				Groups: models.NewRelatedGroups([]models.GroupsScenes{
					{
						GroupID:    invalidID,
						SceneIndex: &sceneIndex,
					},
				}),
			},
			true,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			var fileIDs []models.FileID
			if tt.newObject.Files.Loaded() {
				for _, f := range tt.newObject.Files.List() {
					fileIDs = append(fileIDs, f.ID)
				}
			}

			s := tt.newObject
			if err := qb.Create(ctx, &s, fileIDs); (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.Create() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				assert.Zero(s.ID)
				return
			}

			assert.NotZero(s.ID)

			copy := tt.newObject
			copy.ID = s.ID

			// load relationships
			if err := loadSceneRelationships(ctx, copy, &s); err != nil {
				t.Errorf("loadSceneRelationships() error = %v", err)
				return
			}

			assert.Equal(copy, s)

			// ensure can find the scene
			found, err := qb.Find(ctx, s.ID)
			if err != nil {
				t.Errorf("sceneQueryBuilder.Find() error = %v", err)
			}

			if !assert.NotNil(found) {
				return
			}

			// load relationships
			if err := loadSceneRelationships(ctx, copy, found); err != nil {
				t.Errorf("loadSceneRelationships() error = %v", err)
				return
			}
			assert.Equal(copy, *found)

			return
		})
	}
}

func clearSceneFileIDs(scene *models.Scene) {
	if scene.Files.Loaded() {
		for _, f := range scene.Files.List() {
			f.Base().ID = 0
		}
	}
}

func makeSceneFileWithID(i int) *models.VideoFile {
	ret := makeSceneFile(i)
	ret.ID = sceneFileIDs[i]
	return ret
}

func Test_sceneQueryBuilder_Update(t *testing.T) {
	var (
		title        = "title"
		code         = "1337"
		details      = "details"
		director     = "director"
		url          = "url"
		rating       = 60
		resumeTime   = 10.0
		playDuration = 34.0
		createdAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		sceneIndex   = 123
		sceneIndex2  = 234
		endpoint1    = "endpoint1"
		endpoint2    = "endpoint2"
		stashID1     = "stashid1"
		stashID2     = "stashid2"

		date, _ = models.ParseDate("2003-02-01")
	)

	tests := []struct {
		name          string
		updatedObject *models.Scene
		wantErr       bool
	}{
		{
			"full",
			&models.Scene{
				ID:           sceneIDs[sceneIdxWithGallery],
				Title:        title,
				Code:         code,
				Details:      details,
				Director:     director,
				URLs:         models.NewRelatedStrings([]string{url}),
				Date:         &date,
				Rating:       &rating,
				Organized:    true,
				StudioID:     &studioIDs[studioIdxWithScene],
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				GalleryIDs:   models.NewRelatedIDs([]int{galleryIDs[galleryIdxWithScene]}),
				TagIDs:       models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithScene]}),
				PerformerIDs: models.NewRelatedIDs([]int{performerIDs[performerIdx1WithScene], performerIDs[performerIdx1WithDupName]}),
				Groups: models.NewRelatedGroups([]models.GroupsScenes{
					{
						GroupID:    groupIDs[groupIdxWithScene],
						SceneIndex: &sceneIndex,
					},
					{
						GroupID:    groupIDs[groupIdxWithStudio],
						SceneIndex: &sceneIndex2,
					},
				}),
				StashIDs: models.NewRelatedStashIDs([]models.StashID{
					{
						StashID:  stashID1,
						Endpoint: endpoint1,
					},
					{
						StashID:  stashID2,
						Endpoint: endpoint2,
					},
				}),
				ResumeTime:   resumeTime,
				PlayDuration: playDuration,
			},
			false,
		},
		{
			"clear nullables",
			&models.Scene{
				ID:           sceneIDs[sceneIdxWithSpacedName],
				GalleryIDs:   models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				Groups:       models.NewRelatedGroups([]models.GroupsScenes{}),
				StashIDs:     models.NewRelatedStashIDs([]models.StashID{}),
			},
			false,
		},
		{
			"clear gallery ids",
			&models.Scene{
				ID:         sceneIDs[sceneIdxWithGallery],
				GalleryIDs: models.NewRelatedIDs([]int{}),
			},
			false,
		},
		{
			"clear tag ids",
			&models.Scene{
				ID:     sceneIDs[sceneIdxWithTag],
				TagIDs: models.NewRelatedIDs([]int{}),
			},
			false,
		},
		{
			"clear performer ids",
			&models.Scene{
				ID:           sceneIDs[sceneIdxWithPerformer],
				PerformerIDs: models.NewRelatedIDs([]int{}),
			},
			false,
		},
		{
			"clear groups",
			&models.Scene{
				ID:     sceneIDs[sceneIdxWithGroup],
				Groups: models.NewRelatedGroups([]models.GroupsScenes{}),
			},
			false,
		},
		{
			"invalid studio id",
			&models.Scene{
				ID:       sceneIDs[sceneIdxWithGallery],
				StudioID: &invalidID,
			},
			true,
		},
		{
			"invalid gallery id",
			&models.Scene{
				ID:         sceneIDs[sceneIdxWithGallery],
				GalleryIDs: models.NewRelatedIDs([]int{invalidID}),
			},
			true,
		},
		{
			"invalid tag id",
			&models.Scene{
				ID:     sceneIDs[sceneIdxWithGallery],
				TagIDs: models.NewRelatedIDs([]int{invalidID}),
			},
			true,
		},
		{
			"invalid performer id",
			&models.Scene{
				ID:           sceneIDs[sceneIdxWithGallery],
				PerformerIDs: models.NewRelatedIDs([]int{invalidID}),
			},
			true,
		},
		{
			"invalid group id",
			&models.Scene{
				ID: sceneIDs[sceneIdxWithSpacedName],
				Groups: models.NewRelatedGroups([]models.GroupsScenes{
					{
						GroupID:    invalidID,
						SceneIndex: &sceneIndex,
					},
				}),
			},
			true,
		},
	}

	qb := db.Scene
	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			copy := *tt.updatedObject

			if err := qb.Update(ctx, tt.updatedObject); (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			s, err := qb.Find(ctx, tt.updatedObject.ID)
			if err != nil {
				t.Errorf("sceneQueryBuilder.Find() error = %v", err)
			}

			// load relationships
			if err := loadSceneRelationships(ctx, copy, s); err != nil {
				t.Errorf("loadSceneRelationships() error = %v", err)
				return
			}

			assert.Equal(copy, *s)
		})
	}
}

func clearScenePartial() models.ScenePartial {
	// leave mandatory fields
	return models.ScenePartial{
		Title:        models.OptionalString{Set: true, Null: true},
		Code:         models.OptionalString{Set: true, Null: true},
		Details:      models.OptionalString{Set: true, Null: true},
		Director:     models.OptionalString{Set: true, Null: true},
		URLs:         &models.UpdateStrings{Mode: models.RelationshipUpdateModeSet},
		Date:         models.OptionalDate{Set: true, Null: true},
		Rating:       models.OptionalInt{Set: true, Null: true},
		StudioID:     models.OptionalInt{Set: true, Null: true},
		GalleryIDs:   &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
		TagIDs:       &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
		PerformerIDs: &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
		StashIDs:     &models.UpdateStashIDs{Mode: models.RelationshipUpdateModeSet},
	}
}

func Test_sceneQueryBuilder_UpdatePartial(t *testing.T) {
	var (
		title        = "title"
		code         = "1337"
		details      = "details"
		director     = "director"
		url          = "url"
		rating       = 60
		resumeTime   = 10.0
		playDuration = 34.0
		createdAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		sceneIndex   = 123
		sceneIndex2  = 234
		endpoint1    = "endpoint1"
		endpoint2    = "endpoint2"
		stashID1     = "stashid1"
		stashID2     = "stashid2"

		date, _ = models.ParseDate("2003-02-01")
	)

	tests := []struct {
		name    string
		id      int
		partial models.ScenePartial
		want    models.Scene
		wantErr bool
	}{
		{
			"full",
			sceneIDs[sceneIdxWithSpacedName],
			models.ScenePartial{
				Title:    models.NewOptionalString(title),
				Code:     models.NewOptionalString(code),
				Details:  models.NewOptionalString(details),
				Director: models.NewOptionalString(director),
				URLs: &models.UpdateStrings{
					Values: []string{url},
					Mode:   models.RelationshipUpdateModeSet,
				},
				Date:      models.NewOptionalDate(date),
				Rating:    models.NewOptionalInt(rating),
				Organized: models.NewOptionalBool(true),
				StudioID:  models.NewOptionalInt(studioIDs[studioIdxWithScene]),
				CreatedAt: models.NewOptionalTime(createdAt),
				UpdatedAt: models.NewOptionalTime(updatedAt),
				GalleryIDs: &models.UpdateIDs{
					IDs:  []int{galleryIDs[galleryIdxWithScene]},
					Mode: models.RelationshipUpdateModeSet,
				},
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithScene], tagIDs[tagIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeSet,
				},
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithScene], performerIDs[performerIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeSet,
				},
				GroupIDs: &models.UpdateGroupIDs{
					Groups: []models.GroupsScenes{
						{
							GroupID:    groupIDs[groupIdxWithScene],
							SceneIndex: &sceneIndex,
						},
						{
							GroupID:    groupIDs[groupIdxWithStudio],
							SceneIndex: &sceneIndex2,
						},
					},
					Mode: models.RelationshipUpdateModeSet,
				},
				StashIDs: &models.UpdateStashIDs{
					StashIDs: []models.StashID{
						{
							StashID:  stashID1,
							Endpoint: endpoint1,
						},
						{
							StashID:  stashID2,
							Endpoint: endpoint2,
						},
					},
					Mode: models.RelationshipUpdateModeSet,
				},
				ResumeTime:   models.NewOptionalFloat64(resumeTime),
				PlayDuration: models.NewOptionalFloat64(playDuration),
			},
			models.Scene{
				ID: sceneIDs[sceneIdxWithSpacedName],
				Files: models.NewRelatedVideoFiles([]*models.VideoFile{
					makeSceneFile(sceneIdxWithSpacedName),
				}),
				Title:        title,
				Code:         code,
				Details:      details,
				Director:     director,
				URLs:         models.NewRelatedStrings([]string{url}),
				Date:         &date,
				Rating:       &rating,
				Organized:    true,
				StudioID:     &studioIDs[studioIdxWithScene],
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				GalleryIDs:   models.NewRelatedIDs([]int{galleryIDs[galleryIdxWithScene]}),
				TagIDs:       models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithScene]}),
				PerformerIDs: models.NewRelatedIDs([]int{performerIDs[performerIdx1WithScene], performerIDs[performerIdx1WithDupName]}),
				Groups: models.NewRelatedGroups([]models.GroupsScenes{
					{
						GroupID:    groupIDs[groupIdxWithScene],
						SceneIndex: &sceneIndex,
					},
					{
						GroupID:    groupIDs[groupIdxWithStudio],
						SceneIndex: &sceneIndex2,
					},
				}),
				StashIDs: models.NewRelatedStashIDs([]models.StashID{
					{
						StashID:  stashID1,
						Endpoint: endpoint1,
					},
					{
						StashID:  stashID2,
						Endpoint: endpoint2,
					},
				}),
				ResumeTime:   resumeTime,
				PlayDuration: playDuration,
			},
			false,
		},
		{
			"clear all",
			sceneIDs[sceneIdxWithSpacedName],
			clearScenePartial(),
			models.Scene{
				ID: sceneIDs[sceneIdxWithSpacedName],
				Files: models.NewRelatedVideoFiles([]*models.VideoFile{
					makeSceneFile(sceneIdxWithSpacedName),
				}),
				GalleryIDs:   models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				Groups:       models.NewRelatedGroups([]models.GroupsScenes{}),
				StashIDs:     models.NewRelatedStashIDs([]models.StashID{}),
				PlayDuration: getScenePlayDuration(sceneIdxWithSpacedName),
				ResumeTime:   getSceneResumeTime(sceneIdxWithSpacedName),
			},
			false,
		},
		{
			"invalid id",
			invalidID,
			models.ScenePartial{},
			models.Scene{},
			true,
		},
	}
	for _, tt := range tests {
		qb := db.Scene

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			got, err := qb.UpdatePartial(ctx, tt.id, tt.partial)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.UpdatePartial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// load relationships
			if err := loadSceneRelationships(ctx, tt.want, got); err != nil {
				t.Errorf("loadSceneRelationships() error = %v", err)
				return
			}

			// ignore file ids
			clearSceneFileIDs(got)

			assert.Equal(tt.want, *got)

			s, err := qb.Find(ctx, tt.id)
			if err != nil {
				t.Errorf("sceneQueryBuilder.Find() error = %v", err)
			}

			// load relationships
			if err := loadSceneRelationships(ctx, tt.want, s); err != nil {
				t.Errorf("loadSceneRelationships() error = %v", err)
				return
			}
			// ignore file ids
			clearSceneFileIDs(s)

			assert.Equal(tt.want, *s)
		})
	}
}

func Test_sceneQueryBuilder_UpdatePartialRelationships(t *testing.T) {
	var (
		sceneIndex  = 123
		sceneIndex2 = 234
		endpoint1   = "endpoint1"
		endpoint2   = "endpoint2"
		stashID1    = "stashid1"
		stashID2    = "stashid2"

		groupScenes = []models.GroupsScenes{
			{
				GroupID:    groupIDs[groupIdxWithDupName],
				SceneIndex: &sceneIndex,
			},
			{
				GroupID:    groupIDs[groupIdxWithStudio],
				SceneIndex: &sceneIndex2,
			},
		}

		stashIDs = []models.StashID{
			{
				StashID:  stashID1,
				Endpoint: endpoint1,
			},
			{
				StashID:  stashID2,
				Endpoint: endpoint2,
			},
		}
	)

	tests := []struct {
		name    string
		id      int
		partial models.ScenePartial
		want    models.Scene
		wantErr bool
	}{
		{
			"add galleries",
			sceneIDs[sceneIdxWithGallery],
			models.ScenePartial{
				GalleryIDs: &models.UpdateIDs{
					IDs:  []int{galleryIDs[galleryIdx1WithImage], galleryIDs[galleryIdx1WithPerformer]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				GalleryIDs: models.NewRelatedIDs(append(indexesToIDs(galleryIDs, sceneGalleries[sceneIdxWithGallery]),
					galleryIDs[galleryIdx1WithImage],
					galleryIDs[galleryIdx1WithPerformer],
				)),
			},
			false,
		},
		{
			"add identical galleries",
			sceneIDs[sceneIdxWithGallery],
			models.ScenePartial{
				GalleryIDs: &models.UpdateIDs{
					IDs:  []int{galleryIDs[galleryIdx1WithImage], galleryIDs[galleryIdx1WithImage]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				GalleryIDs: models.NewRelatedIDs(append(indexesToIDs(galleryIDs, sceneGalleries[sceneIdxWithGallery]),
					galleryIDs[galleryIdx1WithImage],
				)),
			},
			false,
		},
		{
			"add tags",
			sceneIDs[sceneIdxWithTwoTags],
			models.ScenePartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithGallery]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				TagIDs: models.NewRelatedIDs(append(
					[]int{
						tagIDs[tagIdx1WithGallery],
						tagIDs[tagIdx1WithDupName],
					},
					indexesToIDs(tagIDs, sceneTags[sceneIdxWithTwoTags])...,
				)),
			},
			false,
		},
		{
			"add identical tags",
			sceneIDs[sceneIdxWithTwoTags],
			models.ScenePartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				TagIDs: models.NewRelatedIDs(append(
					[]int{
						tagIDs[tagIdx1WithDupName],
					},
					indexesToIDs(tagIDs, sceneTags[sceneIdxWithTwoTags])...,
				)),
			},
			false,
		},
		{
			"add performers",
			sceneIDs[sceneIdxWithTwoPerformers],
			models.ScenePartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithDupName], performerIDs[performerIdx1WithGallery]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				PerformerIDs: models.NewRelatedIDs(append(indexesToIDs(performerIDs, scenePerformers[sceneIdxWithTwoPerformers]),
					performerIDs[performerIdx1WithDupName],
					performerIDs[performerIdx1WithGallery],
				)),
			},
			false,
		},
		{
			"add identical performers",
			sceneIDs[sceneIdxWithTwoPerformers],
			models.ScenePartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithDupName], performerIDs[performerIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				PerformerIDs: models.NewRelatedIDs(append(indexesToIDs(performerIDs, scenePerformers[sceneIdxWithTwoPerformers]),
					performerIDs[performerIdx1WithDupName],
				)),
			},
			false,
		},
		{
			"add groups",
			sceneIDs[sceneIdxWithGroup],
			models.ScenePartial{
				GroupIDs: &models.UpdateGroupIDs{
					Groups: groupScenes,
					Mode:   models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				Groups: models.NewRelatedGroups(append([]models.GroupsScenes{
					{
						GroupID: indexesToIDs(groupIDs, sceneGroups[sceneIdxWithGroup])[0],
					},
				}, groupScenes...)),
			},
			false,
		},
		{
			"add groups to empty",
			sceneIDs[sceneIdx1WithPerformer],
			models.ScenePartial{
				GroupIDs: &models.UpdateGroupIDs{
					Groups: groupScenes,
					Mode:   models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				Groups: models.NewRelatedGroups([]models.GroupsScenes{
					{
						GroupID:    groupIDs[groupIdxWithDupName],
						SceneIndex: &sceneIndex,
					},
					{
						GroupID:    groupIDs[groupIdxWithStudio],
						SceneIndex: &sceneIndex2,
					},
				}),
			},
			false,
		},
		{
			"add stash ids",
			sceneIDs[sceneIdxWithSpacedName],
			models.ScenePartial{
				StashIDs: &models.UpdateStashIDs{
					StashIDs: stashIDs,
					Mode:     models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				StashIDs: models.NewRelatedStashIDs(append([]models.StashID{sceneStashID(sceneIdxWithSpacedName)}, stashIDs...)),
			},
			false,
		},
		{
			"add duplicate galleries",
			sceneIDs[sceneIdxWithGallery],
			models.ScenePartial{
				GalleryIDs: &models.UpdateIDs{
					IDs:  []int{galleryIDs[galleryIdxWithScene], galleryIDs[galleryIdx1WithPerformer]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				GalleryIDs: models.NewRelatedIDs(append(indexesToIDs(galleryIDs, sceneGalleries[sceneIdxWithGallery]),
					galleryIDs[galleryIdx1WithPerformer],
				)),
			},
			false,
		},
		{
			"add duplicate tags",
			sceneIDs[sceneIdxWithTwoTags],
			models.ScenePartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithScene], tagIDs[tagIdx1WithGallery]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				TagIDs: models.NewRelatedIDs(append(
					[]int{tagIDs[tagIdx1WithGallery]},
					indexesToIDs(tagIDs, sceneTags[sceneIdxWithTwoTags])...,
				)),
			},
			false,
		},
		{
			"add duplicate performers",
			sceneIDs[sceneIdxWithTwoPerformers],
			models.ScenePartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithScene], performerIDs[performerIdx1WithGallery]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				PerformerIDs: models.NewRelatedIDs(append(indexesToIDs(performerIDs, scenePerformers[sceneIdxWithTwoPerformers]),
					performerIDs[performerIdx1WithGallery],
				)),
			},
			false,
		},
		{
			"add duplicate groups",
			sceneIDs[sceneIdxWithGroup],
			models.ScenePartial{
				GroupIDs: &models.UpdateGroupIDs{
					Groups: append([]models.GroupsScenes{
						{
							GroupID:    groupIDs[groupIdxWithScene],
							SceneIndex: &sceneIndex,
						},
					},
						groupScenes...,
					),
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				Groups: models.NewRelatedGroups(append([]models.GroupsScenes{
					{
						GroupID: indexesToIDs(groupIDs, sceneGroups[sceneIdxWithGroup])[0],
					},
				}, groupScenes...)),
			},
			false,
		},
		{
			"add duplicate stash ids",
			sceneIDs[sceneIdxWithSpacedName],
			models.ScenePartial{
				StashIDs: &models.UpdateStashIDs{
					StashIDs: []models.StashID{
						sceneStashID(sceneIdxWithSpacedName),
					},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				StashIDs: models.NewRelatedStashIDs([]models.StashID{sceneStashID(sceneIdxWithSpacedName)}),
			},
			false,
		},
		{
			"add invalid galleries",
			sceneIDs[sceneIdxWithGallery],
			models.ScenePartial{
				GalleryIDs: &models.UpdateIDs{
					IDs:  []int{invalidID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{},
			true,
		},
		{
			"add invalid tags",
			sceneIDs[sceneIdxWithTwoTags],
			models.ScenePartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{invalidID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{},
			true,
		},
		{
			"add invalid performers",
			sceneIDs[sceneIdxWithTwoPerformers],
			models.ScenePartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{invalidID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{},
			true,
		},
		{
			"add invalid groups",
			sceneIDs[sceneIdxWithGroup],
			models.ScenePartial{
				GroupIDs: &models.UpdateGroupIDs{
					Groups: []models.GroupsScenes{
						{
							GroupID: invalidID,
						},
					},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{},
			true,
		},
		{
			"remove galleries",
			sceneIDs[sceneIdxWithGallery],
			models.ScenePartial{
				GalleryIDs: &models.UpdateIDs{
					IDs:  []int{galleryIDs[galleryIdxWithScene]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Scene{
				GalleryIDs: models.NewRelatedIDs([]int{}),
			},
			false,
		},
		{
			"remove tags",
			sceneIDs[sceneIdxWithTwoTags],
			models.ScenePartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithScene]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Scene{
				TagIDs: models.NewRelatedIDs([]int{tagIDs[tagIdx2WithScene]}),
			},
			false,
		},
		{
			"remove performers",
			sceneIDs[sceneIdxWithTwoPerformers],
			models.ScenePartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithScene]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Scene{
				PerformerIDs: models.NewRelatedIDs([]int{performerIDs[performerIdx2WithScene]}),
			},
			false,
		},
		{
			"remove groups",
			sceneIDs[sceneIdxWithGroup],
			models.ScenePartial{
				GroupIDs: &models.UpdateGroupIDs{
					Groups: []models.GroupsScenes{
						{
							GroupID: groupIDs[groupIdxWithScene],
						},
					},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Scene{
				Groups: models.NewRelatedGroups([]models.GroupsScenes{}),
			},
			false,
		},
		{
			"remove stash ids",
			sceneIDs[sceneIdxWithSpacedName],
			models.ScenePartial{
				StashIDs: &models.UpdateStashIDs{
					StashIDs: []models.StashID{sceneStashID(sceneIdxWithSpacedName)},
					Mode:     models.RelationshipUpdateModeRemove,
				},
			},
			models.Scene{
				StashIDs: models.NewRelatedStashIDs([]models.StashID{}),
			},
			false,
		},
		{
			"remove unrelated galleries",
			sceneIDs[sceneIdxWithGallery],
			models.ScenePartial{
				GalleryIDs: &models.UpdateIDs{
					IDs:  []int{galleryIDs[galleryIdx1WithImage]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Scene{
				GalleryIDs: models.NewRelatedIDs([]int{galleryIDs[galleryIdxWithScene]}),
			},
			false,
		},
		{
			"remove unrelated tags",
			sceneIDs[sceneIdxWithTwoTags],
			models.ScenePartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithPerformer]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Scene{
				TagIDs: models.NewRelatedIDs(indexesToIDs(tagIDs, sceneTags[sceneIdxWithTwoTags])),
			},
			false,
		},
		{
			"remove unrelated performers",
			sceneIDs[sceneIdxWithTwoPerformers],
			models.ScenePartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Scene{
				PerformerIDs: models.NewRelatedIDs(indexesToIDs(performerIDs, scenePerformers[sceneIdxWithTwoPerformers])),
			},
			false,
		},
		{
			"remove unrelated groups",
			sceneIDs[sceneIdxWithGroup],
			models.ScenePartial{
				GroupIDs: &models.UpdateGroupIDs{
					Groups: []models.GroupsScenes{
						{
							GroupID: groupIDs[groupIdxWithDupName],
						},
					},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Scene{
				Groups: models.NewRelatedGroups([]models.GroupsScenes{
					{
						GroupID: indexesToIDs(groupIDs, sceneGroups[sceneIdxWithGroup])[0],
					},
				}),
			},
			false,
		},
		{
			"remove unrelated stash ids",
			sceneIDs[sceneIdxWithGallery],
			models.ScenePartial{
				StashIDs: &models.UpdateStashIDs{
					StashIDs: stashIDs,
					Mode:     models.RelationshipUpdateModeRemove,
				},
			},
			models.Scene{
				StashIDs: models.NewRelatedStashIDs([]models.StashID{sceneStashID(sceneIdxWithGallery)}),
			},
			false,
		},
	}

	for _, tt := range tests {
		qb := db.Scene

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			got, err := qb.UpdatePartial(ctx, tt.id, tt.partial)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.UpdatePartial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			s, err := qb.Find(ctx, tt.id)
			if err != nil {
				t.Errorf("sceneQueryBuilder.Find() error = %v", err)
			}

			// load relationships
			if err := loadSceneRelationships(ctx, tt.want, got); err != nil {
				t.Errorf("loadSceneRelationships() error = %v", err)
				return
			}
			if err := loadSceneRelationships(ctx, tt.want, s); err != nil {
				t.Errorf("loadSceneRelationships() error = %v", err)
				return
			}

			// only compare fields that were in the partial
			if tt.partial.PerformerIDs != nil {
				assert.ElementsMatch(tt.want.PerformerIDs.List(), got.PerformerIDs.List())
				assert.ElementsMatch(tt.want.PerformerIDs.List(), s.PerformerIDs.List())
			}
			if tt.partial.TagIDs != nil {
				assert.ElementsMatch(tt.want.TagIDs.List(), got.TagIDs.List())
				assert.ElementsMatch(tt.want.TagIDs.List(), s.TagIDs.List())
			}
			if tt.partial.GalleryIDs != nil {
				assert.ElementsMatch(tt.want.GalleryIDs.List(), got.GalleryIDs.List())
				assert.ElementsMatch(tt.want.GalleryIDs.List(), s.GalleryIDs.List())
			}
			if tt.partial.GroupIDs != nil {
				assert.ElementsMatch(tt.want.Groups.List(), got.Groups.List())
				assert.ElementsMatch(tt.want.Groups.List(), s.Groups.List())
			}
			if tt.partial.StashIDs != nil {
				assert.ElementsMatch(tt.want.StashIDs.List(), got.StashIDs.List())
				assert.ElementsMatch(tt.want.StashIDs.List(), s.StashIDs.List())
			}
		})
	}
}

func Test_sceneQueryBuilder_AddO(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    int
		wantErr bool
	}{
		{
			"increment",
			sceneIDs[1],
			1,
			false,
		},
		{
			"invalid",
			invalidID,
			0,
			true,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.AddO(ctx, tt.id, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.AddO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("sceneQueryBuilder.AddO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sceneQueryBuilder_DeleteO(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    int
		wantErr bool
	}{
		{
			"decrement",
			sceneIDs[2],
			0,
			false,
		},
		{
			"zero",
			sceneIDs[0],
			0,
			false,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.DeleteO(ctx, tt.id, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.DeleteO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("sceneQueryBuilder.DeleteO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sceneQueryBuilder_ResetO(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    int
		wantErr bool
	}{
		{
			"decrement",
			sceneIDs[2],
			0,
			false,
		},
		{
			"zero",
			sceneIDs[0],
			0,
			false,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.ResetO(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.ResetO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sceneQueryBuilder.ResetOCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sceneQueryBuilder_ResetWatchCount(t *testing.T) {
	return
}

func Test_sceneQueryBuilder_Destroy(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			"valid",
			sceneIDs[sceneIdxWithGallery],
			false,
		},
		{
			"invalid",
			invalidID,
			true,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			if err := qb.Destroy(ctx, tt.id); (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.Destroy() error = %v, wantErr %v", err, tt.wantErr)
			}

			// ensure cannot be found
			i, err := qb.Find(ctx, tt.id)

			assert.Nil(err)
			assert.Nil(i)
		})
	}
}

func makeSceneWithID(index int) *models.Scene {
	ret := makeScene(index)
	ret.ID = sceneIDs[index]

	ret.Files = models.NewRelatedVideoFiles([]*models.VideoFile{makeSceneFile(index)})

	return ret
}

func Test_sceneQueryBuilder_Find(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    *models.Scene
		wantErr bool
	}{
		{
			"valid",
			sceneIDs[sceneIdxWithSpacedName],
			makeSceneWithID(sceneIdxWithSpacedName),
			false,
		},
		{
			"invalid",
			invalidID,
			nil,
			false,
		},
		{
			"with galleries",
			sceneIDs[sceneIdxWithGallery],
			makeSceneWithID(sceneIdxWithGallery),
			false,
		},
		{
			"with performers",
			sceneIDs[sceneIdxWithTwoPerformers],
			makeSceneWithID(sceneIdxWithTwoPerformers),
			false,
		},
		{
			"with tags",
			sceneIDs[sceneIdxWithTwoTags],
			makeSceneWithID(sceneIdxWithTwoTags),
			false,
		},
		{
			"with groups",
			sceneIDs[sceneIdxWithGroup],
			makeSceneWithID(sceneIdxWithGroup),
			false,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.Find(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				// load relationships
				if err := loadSceneRelationships(ctx, *tt.want, got); err != nil {
					t.Errorf("loadSceneRelationships() error = %v", err)
					return
				}

				clearSceneFileIDs(got)
			}

			assert.Equal(tt.want, got)
		})
	}
}

func postFindScenes(ctx context.Context, want []*models.Scene, got []*models.Scene) error {
	for i, s := range got {
		// load relationships
		if i < len(want) {
			if err := loadSceneRelationships(ctx, *want[i], s); err != nil {
				return err
			}
		}
		clearSceneFileIDs(s)
	}

	return nil
}

func Test_sceneQueryBuilder_FindMany(t *testing.T) {
	tests := []struct {
		name    string
		ids     []int
		want    []*models.Scene
		wantErr bool
	}{
		{
			"valid with relationships",
			[]int{
				sceneIDs[sceneIdxWithGallery],
				sceneIDs[sceneIdxWithTwoPerformers],
				sceneIDs[sceneIdxWithTwoTags],
				sceneIDs[sceneIdxWithGroup],
			},
			[]*models.Scene{
				makeSceneWithID(sceneIdxWithGallery),
				makeSceneWithID(sceneIdxWithTwoPerformers),
				makeSceneWithID(sceneIdxWithTwoTags),
				makeSceneWithID(sceneIdxWithGroup),
			},
			false,
		},
		{
			"invalid",
			[]int{sceneIDs[sceneIdxWithGallery], sceneIDs[sceneIdxWithTwoPerformers], invalidID},
			nil,
			true,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindMany(ctx, tt.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.FindMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := postFindScenes(ctx, tt.want, got); err != nil {
				t.Errorf("loadSceneRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func Test_sceneQueryBuilder_FindByChecksum(t *testing.T) {
	getChecksum := func(index int) string {
		return getSceneStringValue(index, checksumField)
	}

	tests := []struct {
		name     string
		checksum string
		want     []*models.Scene
		wantErr  bool
	}{
		{
			"valid",
			getChecksum(sceneIdxWithSpacedName),
			[]*models.Scene{makeSceneWithID(sceneIdxWithSpacedName)},
			false,
		},
		{
			"invalid",
			"invalid checksum",
			nil,
			false,
		},
		{
			"with galleries",
			getChecksum(sceneIdxWithGallery),
			[]*models.Scene{makeSceneWithID(sceneIdxWithGallery)},
			false,
		},
		{
			"with performers",
			getChecksum(sceneIdxWithTwoPerformers),
			[]*models.Scene{makeSceneWithID(sceneIdxWithTwoPerformers)},
			false,
		},
		{
			"with tags",
			getChecksum(sceneIdxWithTwoTags),
			[]*models.Scene{makeSceneWithID(sceneIdxWithTwoTags)},
			false,
		},
		{
			"with groups",
			getChecksum(sceneIdxWithGroup),
			[]*models.Scene{makeSceneWithID(sceneIdxWithGroup)},
			false,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByChecksum(ctx, tt.checksum)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.FindByChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := postFindScenes(ctx, tt.want, got); err != nil {
				t.Errorf("loadSceneRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func Test_sceneQueryBuilder_FindByOSHash(t *testing.T) {
	getOSHash := func(index int) string {
		return getSceneStringValue(index, "oshash")
	}

	tests := []struct {
		name    string
		oshash  string
		want    []*models.Scene
		wantErr bool
	}{
		{
			"valid",
			getOSHash(sceneIdxWithSpacedName),
			[]*models.Scene{makeSceneWithID(sceneIdxWithSpacedName)},
			false,
		},
		{
			"invalid",
			"invalid oshash",
			nil,
			false,
		},
		{
			"with galleries",
			getOSHash(sceneIdxWithGallery),
			[]*models.Scene{makeSceneWithID(sceneIdxWithGallery)},
			false,
		},
		{
			"with performers",
			getOSHash(sceneIdxWithTwoPerformers),
			[]*models.Scene{makeSceneWithID(sceneIdxWithTwoPerformers)},
			false,
		},
		{
			"with tags",
			getOSHash(sceneIdxWithTwoTags),
			[]*models.Scene{makeSceneWithID(sceneIdxWithTwoTags)},
			false,
		},
		{
			"with groups",
			getOSHash(sceneIdxWithGroup),
			[]*models.Scene{makeSceneWithID(sceneIdxWithGroup)},
			false,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.FindByOSHash(ctx, tt.oshash)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.FindByOSHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := postFindScenes(ctx, tt.want, got); err != nil {
				t.Errorf("loadSceneRelationships() error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sceneQueryBuilder.FindByOSHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sceneQueryBuilder_FindByPath(t *testing.T) {
	getPath := func(index int) string {
		return getFilePath(folderIdxWithSceneFiles, getSceneBasename(index))
	}

	tests := []struct {
		name    string
		path    string
		want    []*models.Scene
		wantErr bool
	}{
		{
			"valid",
			getPath(sceneIdxWithSpacedName),
			[]*models.Scene{makeSceneWithID(sceneIdxWithSpacedName)},
			false,
		},
		{
			"invalid",
			"invalid path",
			nil,
			false,
		},
		{
			"with galleries",
			getPath(sceneIdxWithGallery),
			[]*models.Scene{makeSceneWithID(sceneIdxWithGallery)},
			false,
		},
		{
			"with performers",
			getPath(sceneIdxWithTwoPerformers),
			[]*models.Scene{makeSceneWithID(sceneIdxWithTwoPerformers)},
			false,
		},
		{
			"with tags",
			getPath(sceneIdxWithTwoTags),
			[]*models.Scene{makeSceneWithID(sceneIdxWithTwoTags)},
			false,
		},
		{
			"with groups",
			getPath(sceneIdxWithGroup),
			[]*models.Scene{makeSceneWithID(sceneIdxWithGroup)},
			false,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByPath(ctx, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.FindByPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := postFindScenes(ctx, tt.want, got); err != nil {
				t.Errorf("loadSceneRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func Test_sceneQueryBuilder_FindByGalleryID(t *testing.T) {
	tests := []struct {
		name      string
		galleryID int
		want      []*models.Scene
		wantErr   bool
	}{
		{
			"valid",
			galleryIDs[galleryIdxWithScene],
			[]*models.Scene{makeSceneWithID(sceneIdxWithGallery)},
			false,
		},
		{
			"none",
			galleryIDs[galleryIdx1WithPerformer],
			nil,
			false,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByGalleryID(ctx, tt.galleryID)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.FindByGalleryID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := postFindScenes(ctx, tt.want, got); err != nil {
				t.Errorf("loadSceneRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, got)
			return
		})
	}
}

func TestSceneCountByPerformerID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		count, err := sqb.CountByPerformerID(ctx, performerIDs[performerIdxWithScene])

		if err != nil {
			t.Errorf("Error counting scenes: %s", err.Error())
		}

		assert.Equal(t, 1, count)

		count, err = sqb.CountByPerformerID(ctx, 0)

		if err != nil {
			t.Errorf("Error counting scenes: %s", err.Error())
		}

		assert.Equal(t, 0, count)

		return nil
	})
}

func scenesToIDs(i []*models.Scene) []int {
	ret := make([]int, len(i))
	for i, v := range i {
		ret[i] = v.ID
	}

	return ret
}

func Test_sceneStore_FindByFileID(t *testing.T) {
	tests := []struct {
		name    string
		fileID  models.FileID
		include []int
		exclude []int
	}{
		{
			"valid",
			sceneFileIDs[sceneIdx1WithPerformer],
			[]int{sceneIdx1WithPerformer},
			nil,
		},
		{
			"invalid",
			invalidFileID,
			nil,
			[]int{sceneIdx1WithPerformer},
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByFileID(ctx, tt.fileID)
			if err != nil {
				t.Errorf("SceneStore.FindByFileID() error = %v", err)
				return
			}
			for _, f := range got {
				clearSceneFileIDs(f)
			}

			ids := scenesToIDs(got)
			include := indexesToIDs(galleryIDs, tt.include)
			exclude := indexesToIDs(galleryIDs, tt.exclude)

			for _, i := range include {
				assert.Contains(ids, i)
			}
			for _, e := range exclude {
				assert.NotContains(ids, e)
			}
		})
	}
}

func Test_sceneStore_CountByFileID(t *testing.T) {
	tests := []struct {
		name   string
		fileID models.FileID
		want   int
	}{
		{
			"valid",
			sceneFileIDs[sceneIdxWithTwoPerformers],
			1,
		},
		{
			"invalid",
			invalidFileID,
			0,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.CountByFileID(ctx, tt.fileID)
			if err != nil {
				t.Errorf("SceneStore.CountByFileID() error = %v", err)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func Test_sceneStore_CountMissingChecksum(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			"valid",
			0,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.CountMissingChecksum(ctx)
			if err != nil {
				t.Errorf("SceneStore.CountMissingChecksum() error = %v", err)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func Test_sceneStore_CountMissingOshash(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			"valid",
			0,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.CountMissingOSHash(ctx)
			if err != nil {
				t.Errorf("SceneStore.CountMissingOSHash() error = %v", err)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func TestSceneWall(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene

		const sceneIdx = 2
		wallQuery := getSceneStringValue(sceneIdx, "Details")
		scenes, err := sqb.Wall(ctx, &wallQuery)

		if err != nil {
			t.Errorf("Error finding scenes: %s", err.Error())
			return nil
		}

		assert.Len(t, scenes, 1)
		scene := scenes[0]
		assert.Equal(t, sceneIDs[sceneIdx], scene.ID)
		scenePath := getFilePath(folderIdxWithSceneFiles, getSceneBasename(sceneIdx))
		assert.Equal(t, scenePath, scene.Path)

		wallQuery = "not exist"
		scenes, err = sqb.Wall(ctx, &wallQuery)

		if err != nil {
			t.Errorf("Error finding scene: %s", err.Error())
			return nil
		}

		assert.Len(t, scenes, 0)

		return nil
	})
}

func TestSceneQueryQ(t *testing.T) {
	const sceneIdx = 2

	q := getSceneStringValue(sceneIdx, titleField)

	withTxn(func(ctx context.Context) error {
		sqb := db.Scene

		sceneQueryQ(ctx, t, sqb, q, sceneIdx)

		return nil
	})
}

func queryScene(ctx context.Context, t *testing.T, sqb models.SceneReader, sceneFilter *models.SceneFilterType, findFilter *models.FindFilterType) []*models.Scene {
	t.Helper()
	result, err := sqb.Query(ctx, models.SceneQueryOptions{
		QueryOptions: models.QueryOptions{
			FindFilter: findFilter,
			Count:      true,
		},
		SceneFilter:   sceneFilter,
		TotalDuration: true,
		TotalSize:     true,
	})
	if err != nil {
		t.Errorf("Error querying scene: %v", err)
		return nil
	}

	scenes, err := result.Resolve(ctx)
	if err != nil {
		t.Errorf("Error resolving scenes: %v", err)
	}

	return scenes
}

func sceneQueryQ(ctx context.Context, t *testing.T, sqb models.SceneReader, q string, expectedSceneIdx int) {
	filter := models.FindFilterType{
		Q: &q,
	}
	scenes := queryScene(ctx, t, sqb, nil, &filter)

	if !assert.Len(t, scenes, 1) {
		return
	}
	scene := scenes[0]
	assert.Equal(t, sceneIDs[expectedSceneIdx], scene.ID)

	// no Q should return all results
	filter.Q = nil
	pp := totalScenes
	filter.PerPage = &pp
	scenes = queryScene(ctx, t, sqb, nil, &filter)

	assert.Len(t, scenes, totalScenes)
}

func TestSceneQuery(t *testing.T) {
	var (
		endpoint = sceneStashID(sceneIdxWithGallery).Endpoint
		stashID  = sceneStashID(sceneIdxWithGallery).StashID

		depth = -1
	)

	tests := []struct {
		name        string
		findFilter  *models.FindFilterType
		filter      *models.SceneFilterType
		includeIdxs []int
		excludeIdxs []int
		wantErr     bool
	}{
		{
			"specific resume time",
			nil,
			&models.SceneFilterType{
				ResumeTime: &models.IntCriterionInput{
					Modifier: models.CriterionModifierEquals,
					Value:    int(getSceneResumeTime(sceneIdxWithGallery)),
				},
			},
			[]int{sceneIdxWithGallery},
			[]int{sceneIdxWithGroup},
			false,
		},
		{
			"specific play duration",
			nil,
			&models.SceneFilterType{
				PlayDuration: &models.IntCriterionInput{
					Modifier: models.CriterionModifierEquals,
					Value:    int(getScenePlayDuration(sceneIdxWithGallery)),
				},
			},
			[]int{sceneIdxWithGallery},
			[]int{sceneIdxWithGroup},
			false,
		},
		// {
		// 	"specific play count",
		// 	nil,
		// 	&models.SceneFilterType{
		// 		PlayCount: &models.IntCriterionInput{
		// 			Modifier: models.CriterionModifierEquals,
		// 			Value:    getScenePlayCount(sceneIdxWithGallery),
		// 		},
		// 	},
		// 	[]int{sceneIdxWithGallery},
		// 	[]int{sceneIdxWithGroup},
		// 	false,
		// },
		{
			"stash id with endpoint",
			nil,
			&models.SceneFilterType{
				StashIDEndpoint: &models.StashIDCriterionInput{
					Endpoint: &endpoint,
					StashID:  &stashID,
					Modifier: models.CriterionModifierEquals,
				},
			},
			[]int{sceneIdxWithGallery},
			nil,
			false,
		},
		{
			"exclude stash id with endpoint",
			nil,
			&models.SceneFilterType{
				StashIDEndpoint: &models.StashIDCriterionInput{
					Endpoint: &endpoint,
					StashID:  &stashID,
					Modifier: models.CriterionModifierNotEquals,
				},
			},
			nil,
			[]int{sceneIdxWithGallery},
			false,
		},
		{
			"null stash id with endpoint",
			nil,
			&models.SceneFilterType{
				StashIDEndpoint: &models.StashIDCriterionInput{
					Endpoint: &endpoint,
					Modifier: models.CriterionModifierIsNull,
				},
			},
			nil,
			[]int{sceneIdxWithGallery},
			false,
		},
		{
			"not null stash id with endpoint",
			nil,
			&models.SceneFilterType{
				StashIDEndpoint: &models.StashIDCriterionInput{
					Endpoint: &endpoint,
					Modifier: models.CriterionModifierNotNull,
				},
			},
			[]int{sceneIdxWithGallery},
			nil,
			false,
		},
		{
			"with studio id 0 including child studios",
			nil,
			&models.SceneFilterType{
				Studios: &models.HierarchicalMultiCriterionInput{
					Value:    []string{"0"},
					Modifier: models.CriterionModifierIncludes,
					Depth:    &depth,
				},
			},
			nil,
			nil,
			false,
		},
	}

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			results, err := db.Scene.Query(ctx, models.SceneQueryOptions{
				SceneFilter: tt.filter,
				QueryOptions: models.QueryOptions{
					FindFilter: tt.findFilter,
				},
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("PerformerStore.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			include := indexesToIDs(sceneIDs, tt.includeIdxs)
			exclude := indexesToIDs(sceneIDs, tt.excludeIdxs)

			for _, i := range include {
				assert.Contains(results.IDs, i)
			}
			for _, e := range exclude {
				assert.NotContains(results.IDs, e)
			}
		})
	}
}

func TestSceneQueryPath(t *testing.T) {
	const (
		sceneIdx      = 1
		otherSceneIdx = 2
	)
	folder := folderPaths[folderIdxWithSceneFiles]
	basename := getSceneBasename(sceneIdx)
	scenePath := getFilePath(folderIdxWithSceneFiles, getSceneBasename(sceneIdx))

	tests := []struct {
		name        string
		input       models.StringCriterionInput
		mustInclude []int
		mustExclude []int
	}{
		{
			"equals full path",
			models.StringCriterionInput{
				Value:    scenePath,
				Modifier: models.CriterionModifierEquals,
			},
			[]int{sceneIdx},
			[]int{otherSceneIdx},
		},
		{
			"equals full path wildcard",
			models.StringCriterionInput{
				Value:    filepath.Join(folder, "scene_0001_%"),
				Modifier: models.CriterionModifierEquals,
			},
			[]int{sceneIdx},
			[]int{otherSceneIdx},
		},
		{
			"not equals full path",
			models.StringCriterionInput{
				Value:    scenePath,
				Modifier: models.CriterionModifierNotEquals,
			},
			[]int{otherSceneIdx},
			[]int{sceneIdx},
		},
		{
			"includes folder name",
			models.StringCriterionInput{
				Value:    folder,
				Modifier: models.CriterionModifierIncludes,
			},
			[]int{sceneIdx},
			nil,
		},
		{
			"includes base name",
			models.StringCriterionInput{
				Value:    basename,
				Modifier: models.CriterionModifierIncludes,
			},
			[]int{sceneIdx},
			nil,
		},
		{
			"includes full path",
			models.StringCriterionInput{
				Value:    scenePath,
				Modifier: models.CriterionModifierIncludes,
			},
			[]int{sceneIdx},
			[]int{otherSceneIdx},
		},
		{
			"matches regex",
			models.StringCriterionInput{
				Value:    "scene_.*1_Path",
				Modifier: models.CriterionModifierMatchesRegex,
			},
			[]int{sceneIdx},
			nil,
		},
		{
			"not matches regex",
			models.StringCriterionInput{
				Value:    "scene_.*1_Path",
				Modifier: models.CriterionModifierNotMatchesRegex,
			},
			nil,
			[]int{sceneIdx},
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.Query(ctx, models.SceneQueryOptions{
				SceneFilter: &models.SceneFilterType{
					Path: &tt.input,
				},
			})

			if err != nil {
				t.Errorf("sceneQueryBuilder.TestSceneQueryPath() error = %v", err)
				return
			}

			mustInclude := indexesToIDs(sceneIDs, tt.mustInclude)
			mustExclude := indexesToIDs(sceneIDs, tt.mustExclude)

			missing := sliceutil.Exclude(mustInclude, got.IDs)
			if len(missing) > 0 {
				t.Errorf("SceneStore.TestSceneQueryPath() missing expected IDs: %v", missing)
			}

			notExcluded := sliceutil.Intersect(mustExclude, got.IDs)
			if len(notExcluded) > 0 {
				t.Errorf("SceneStore.TestSceneQueryPath() expected IDs to be excluded: %v", notExcluded)
			}
		})
	}
}

func TestSceneQueryURL(t *testing.T) {
	const sceneIdx = 1
	sceneURL := getSceneStringValue(sceneIdx, urlField)

	urlCriterion := models.StringCriterionInput{
		Value:    sceneURL,
		Modifier: models.CriterionModifierEquals,
	}

	filter := models.SceneFilterType{
		URL: &urlCriterion,
	}

	verifyFn := func(s *models.Scene) {
		t.Helper()

		urls := s.URLs.List()
		var url string
		if len(urls) > 0 {
			url = urls[0]
		}

		verifyString(t, url, urlCriterion)
	}

	verifySceneQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotEquals
	verifySceneQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierMatchesRegex
	urlCriterion.Value = "scene_.*1_URL"
	verifySceneQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotMatchesRegex
	verifySceneQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierIsNull
	urlCriterion.Value = ""
	verifySceneQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotNull
	verifySceneQuery(t, filter, verifyFn)
}

func TestSceneQueryPathOr(t *testing.T) {
	const scene1Idx = 1
	const scene2Idx = 2

	scene1Path := getFilePath(folderIdxWithSceneFiles, getSceneBasename(scene1Idx))
	scene2Path := getFilePath(folderIdxWithSceneFiles, getSceneBasename(scene2Idx))

	sceneFilter := models.SceneFilterType{
		Path: &models.StringCriterionInput{
			Value:    scene1Path,
			Modifier: models.CriterionModifierEquals,
		},
		OperatorFilter: models.OperatorFilter[models.SceneFilterType]{
			Or: &models.SceneFilterType{
				Path: &models.StringCriterionInput{
					Value:    scene2Path,
					Modifier: models.CriterionModifierEquals,
				},
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Scene

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		if !assert.Len(t, scenes, 2) {
			return nil
		}
		assert.Equal(t, scene1Path, scenes[0].Path)
		assert.Equal(t, scene2Path, scenes[1].Path)

		return nil
	})
}

func TestSceneQueryPathAndRating(t *testing.T) {
	const sceneIdx = 1
	scenePath := getFilePath(folderIdxWithSceneFiles, getSceneBasename(sceneIdx))
	sceneRating := int(getRating(sceneIdx).Int64)

	sceneFilter := models.SceneFilterType{
		Path: &models.StringCriterionInput{
			Value:    scenePath,
			Modifier: models.CriterionModifierEquals,
		},
		OperatorFilter: models.OperatorFilter[models.SceneFilterType]{
			And: &models.SceneFilterType{
				Rating100: &models.IntCriterionInput{
					Value:    sceneRating,
					Modifier: models.CriterionModifierEquals,
				},
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Scene

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		if !assert.Len(t, scenes, 1) {
			return nil
		}
		assert.Equal(t, scenePath, scenes[0].Path)
		assert.Equal(t, sceneRating, *scenes[0].Rating)

		return nil
	})
}

func TestSceneQueryPathNotRating(t *testing.T) {
	const sceneIdx = 1

	sceneRating := getRating(sceneIdx)

	pathCriterion := models.StringCriterionInput{
		Value:    "scene_.*1_Path",
		Modifier: models.CriterionModifierMatchesRegex,
	}

	ratingCriterion := models.IntCriterionInput{
		Value:    int(sceneRating.Int64),
		Modifier: models.CriterionModifierEquals,
	}

	sceneFilter := models.SceneFilterType{
		Path: &pathCriterion,
		OperatorFilter: models.OperatorFilter[models.SceneFilterType]{
			Not: &models.SceneFilterType{
				Rating100: &ratingCriterion,
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Scene

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		for _, scene := range scenes {
			verifyString(t, scene.Path, pathCriterion)
			ratingCriterion.Modifier = models.CriterionModifierNotEquals
			verifyIntPtr(t, scene.Rating, ratingCriterion)
		}

		return nil
	})
}

func TestSceneIllegalQuery(t *testing.T) {
	assert := assert.New(t)

	const sceneIdx = 1
	subFilter := models.SceneFilterType{
		Path: &models.StringCriterionInput{
			Value:    getSceneStringValue(sceneIdx, "Path"),
			Modifier: models.CriterionModifierEquals,
		},
	}

	sceneFilter := &models.SceneFilterType{
		OperatorFilter: models.OperatorFilter[models.SceneFilterType]{
			And: &subFilter,
			Or:  &subFilter,
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Scene

		queryOptions := models.SceneQueryOptions{
			SceneFilter: sceneFilter,
		}

		_, err := sqb.Query(ctx, queryOptions)
		assert.NotNil(err)

		sceneFilter.Or = nil
		sceneFilter.Not = &subFilter
		_, err = sqb.Query(ctx, queryOptions)
		assert.NotNil(err)

		sceneFilter.And = nil
		sceneFilter.Or = &subFilter
		_, err = sqb.Query(ctx, queryOptions)
		assert.NotNil(err)

		return nil
	})
}

func verifySceneQuery(t *testing.T, filter models.SceneFilterType, verifyFn func(s *models.Scene)) {
	t.Helper()
	withTxn(func(ctx context.Context) error {
		t.Helper()
		sqb := db.Scene

		scenes := queryScene(ctx, t, sqb, &filter, nil)

		for _, scene := range scenes {
			if err := scene.LoadRelationships(ctx, sqb); err != nil {
				t.Errorf("Error loading scene relationships: %v", err)
			}
		}

		// assume it should find at least one
		assert.Greater(t, len(scenes), 0)

		for _, scene := range scenes {
			verifyFn(scene)
		}

		return nil
	})
}

func verifyScenesPath(t *testing.T, pathCriterion models.StringCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		sceneFilter := models.SceneFilterType{
			Path: &pathCriterion,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		for _, scene := range scenes {
			verifyString(t, scene.Path, pathCriterion)
		}

		return nil
	})
}

func verifyStringPtr(t *testing.T, value *string, criterion models.StringCriterionInput) {
	t.Helper()
	assert := assert.New(t)
	if criterion.Modifier == models.CriterionModifierIsNull {
		if value != nil && *value == "" {
			// correct
			return
		}
		assert.Nil(value, "expect is null values to be null")
	}
	if criterion.Modifier == models.CriterionModifierNotNull {
		assert.NotNil(value, "expect is null values to be null")
		assert.Greater(len(*value), 0)
	}
	if criterion.Modifier == models.CriterionModifierEquals {
		assert.Equal(criterion.Value, *value)
	}
	if criterion.Modifier == models.CriterionModifierNotEquals {
		assert.NotEqual(criterion.Value, *value)
	}
	if criterion.Modifier == models.CriterionModifierMatchesRegex {
		assert.NotNil(value)
		assert.Regexp(regexp.MustCompile(criterion.Value), *value)
	}
	if criterion.Modifier == models.CriterionModifierNotMatchesRegex {
		if value == nil {
			// correct
			return
		}
		assert.NotRegexp(regexp.MustCompile(criterion.Value), value)
	}
}

func verifyString(t *testing.T, value string, criterion models.StringCriterionInput) {
	t.Helper()
	assert := assert.New(t)
	switch criterion.Modifier {
	case models.CriterionModifierEquals:
		assert.Equal(criterion.Value, value)
	case models.CriterionModifierNotEquals:
		assert.NotEqual(criterion.Value, value)
	case models.CriterionModifierMatchesRegex:
		assert.Regexp(regexp.MustCompile(criterion.Value), value)
	case models.CriterionModifierNotMatchesRegex:
		assert.NotRegexp(regexp.MustCompile(criterion.Value), value)
	case models.CriterionModifierIsNull:
		assert.Equal("", value)
	case models.CriterionModifierNotNull:
		assert.NotEqual("", value)
	}
}

func TestSceneQueryRating100(t *testing.T) {
	const rating = 60
	ratingCriterion := models.IntCriterionInput{
		Value:    rating,
		Modifier: models.CriterionModifierEquals,
	}

	verifyScenesRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotEquals
	verifyScenesRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyScenesRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierLessThan
	verifyScenesRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierIsNull
	verifyScenesRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotNull
	verifyScenesRating100(t, ratingCriterion)
}

func verifyScenesRating100(t *testing.T, ratingCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		sceneFilter := models.SceneFilterType{
			Rating100: &ratingCriterion,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		for _, scene := range scenes {
			verifyIntPtr(t, scene.Rating, ratingCriterion)
		}

		return nil
	})
}

func verifyIntPtr(t *testing.T, value *int, criterion models.IntCriterionInput) {
	t.Helper()
	assert := assert.New(t)
	if criterion.Modifier == models.CriterionModifierIsNull {
		assert.Nil(value, "expect is null values to be null")
	}
	if criterion.Modifier == models.CriterionModifierNotNull {
		assert.NotNil(value, "expect is null values to be null")
	}
	if criterion.Modifier == models.CriterionModifierEquals {
		assert.Equal(criterion.Value, *value)
	}
	if criterion.Modifier == models.CriterionModifierNotEquals {
		assert.NotEqual(criterion.Value, *value)
	}
	if criterion.Modifier == models.CriterionModifierGreaterThan {
		assert.True(*value > criterion.Value)
	}
	if criterion.Modifier == models.CriterionModifierLessThan {
		assert.True(*value < criterion.Value)
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
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		sceneFilter := models.SceneFilterType{
			OCounter: &oCounterCriterion,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		for _, scene := range scenes {
			count, err := sqb.GetOCount(ctx, scene.ID)
			if err != nil {
				t.Errorf("Error getting ocounter: %v", err)
			}
			verifyInt(t, count, oCounterCriterion)
		}

		return nil
	})
}

func verifyInt(t *testing.T, value int, criterion models.IntCriterionInput) bool {
	t.Helper()
	assert := assert.New(t)
	if criterion.Modifier == models.CriterionModifierEquals {
		return assert.Equal(criterion.Value, value)
	}
	if criterion.Modifier == models.CriterionModifierNotEquals {
		return assert.NotEqual(criterion.Value, value)
	}
	if criterion.Modifier == models.CriterionModifierGreaterThan {
		return assert.Greater(value, criterion.Value)
	}
	if criterion.Modifier == models.CriterionModifierLessThan {
		return assert.Less(value, criterion.Value)
	}

	return true
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
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		sceneFilter := models.SceneFilterType{
			Duration: &durationCriterion,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		for _, scene := range scenes {
			if err := scene.LoadPrimaryFile(ctx, db.File); err != nil {
				t.Errorf("Error querying scene files: %v", err)
				return nil
			}

			duration := scene.Files.Primary().Duration
			if durationCriterion.Modifier == models.CriterionModifierEquals {
				assert.True(t, duration >= float64(durationCriterion.Value) && duration < float64(durationCriterion.Value+1))
			} else if durationCriterion.Modifier == models.CriterionModifierNotEquals {
				assert.True(t, duration < float64(durationCriterion.Value) || duration >= float64(durationCriterion.Value+1))
			} else {
				verifyFloat64(t, duration, durationCriterion)
			}
		}

		return nil
	})
}

func verifyFloat64(t *testing.T, value float64, criterion models.IntCriterionInput) {
	assert := assert.New(t)
	if criterion.Modifier == models.CriterionModifierEquals {
		assert.Equal(float64(criterion.Value), value)
	}
	if criterion.Modifier == models.CriterionModifierNotEquals {
		assert.NotEqual(float64(criterion.Value), value)
	}
	if criterion.Modifier == models.CriterionModifierGreaterThan {
		assert.True(value > float64(criterion.Value))
	}
	if criterion.Modifier == models.CriterionModifierLessThan {
		assert.True(value < float64(criterion.Value))
	}
}

func verifyFloat64Ptr(t *testing.T, value *float64, criterion models.IntCriterionInput) {
	assert := assert.New(t)
	switch criterion.Modifier {
	case models.CriterionModifierIsNull:
		assert.Nil(value, "expect is null values to be null")
	case models.CriterionModifierNotNull:
		assert.NotNil(value, "expect is not null values to not be null")
	case models.CriterionModifierEquals:
		assert.EqualValues(float64(criterion.Value), value)
	case models.CriterionModifierNotEquals:
		assert.NotEqualValues(float64(criterion.Value), value)
	case models.CriterionModifierGreaterThan:
		assert.True(value != nil && *value > float64(criterion.Value))
	case models.CriterionModifierLessThan:
		assert.True(value != nil && *value < float64(criterion.Value))
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
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		sceneFilter := models.SceneFilterType{
			Resolution: &models.ResolutionCriterionInput{
				Value:    resolution,
				Modifier: models.CriterionModifierEquals,
			},
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		for _, scene := range scenes {
			if err := scene.LoadPrimaryFile(ctx, db.File); err != nil {
				t.Errorf("Error querying scene files: %v", err)
				return nil
			}
			f := scene.Files.Primary()
			height := 0
			if f != nil {
				height = f.Height
			}
			verifySceneResolution(t, &height, resolution)
		}

		return nil
	})
}

func verifySceneResolution(t *testing.T, height *int, resolution models.ResolutionEnum) {
	if !resolution.IsValid() {
		return
	}

	assert := assert.New(t)
	assert.NotNil(height)
	if t.Failed() {
		return
	}

	h := *height

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

func TestAllResolutionsHaveResolutionRange(t *testing.T) {
	for _, resolution := range models.AllResolutionEnum {
		assert.NotZero(t, resolution.GetMinResolution(), "Define resolution range for %s in extension_resolution.go", resolution)
		assert.NotZero(t, resolution.GetMaxResolution(), "Define resolution range for %s in extension_resolution.go", resolution)
	}
}

func TestSceneQueryResolutionModifiers(t *testing.T) {
	if err := withRollbackTxn(func(ctx context.Context) error {
		qb := db.Scene
		sceneNoResolution, _ := createScene(ctx, 0, 0)
		firstScene540P, _ := createScene(ctx, 960, 540)
		secondScene540P, _ := createScene(ctx, 1280, 719)
		firstScene720P, _ := createScene(ctx, 1280, 720)
		secondScene720P, _ := createScene(ctx, 1280, 721)
		thirdScene720P, _ := createScene(ctx, 1920, 1079)
		scene1080P, _ := createScene(ctx, 1920, 1080)

		scenesEqualTo720P := queryScenes(ctx, t, qb, models.ResolutionEnumStandardHd, models.CriterionModifierEquals)
		scenesNotEqualTo720P := queryScenes(ctx, t, qb, models.ResolutionEnumStandardHd, models.CriterionModifierNotEquals)
		scenesGreaterThan720P := queryScenes(ctx, t, qb, models.ResolutionEnumStandardHd, models.CriterionModifierGreaterThan)
		scenesLessThan720P := queryScenes(ctx, t, qb, models.ResolutionEnumStandardHd, models.CriterionModifierLessThan)

		assert.Subset(t, scenesEqualTo720P, []*models.Scene{firstScene720P, secondScene720P, thirdScene720P})
		assert.NotSubset(t, scenesEqualTo720P, []*models.Scene{sceneNoResolution, firstScene540P, secondScene540P, scene1080P})

		assert.Subset(t, scenesNotEqualTo720P, []*models.Scene{sceneNoResolution, firstScene540P, secondScene540P, scene1080P})
		assert.NotSubset(t, scenesNotEqualTo720P, []*models.Scene{firstScene720P, secondScene720P, thirdScene720P})

		assert.Subset(t, scenesGreaterThan720P, []*models.Scene{scene1080P})
		assert.NotSubset(t, scenesGreaterThan720P, []*models.Scene{sceneNoResolution, firstScene540P, secondScene540P, firstScene720P, secondScene720P, thirdScene720P})

		assert.Subset(t, scenesLessThan720P, []*models.Scene{sceneNoResolution, firstScene540P, secondScene540P})
		assert.NotSubset(t, scenesLessThan720P, []*models.Scene{scene1080P, firstScene720P, secondScene720P, thirdScene720P})

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func queryScenes(ctx context.Context, t *testing.T, queryBuilder models.SceneReaderWriter, resolution models.ResolutionEnum, modifier models.CriterionModifier) []*models.Scene {
	sceneFilter := models.SceneFilterType{
		Resolution: &models.ResolutionCriterionInput{
			Value:    resolution,
			Modifier: modifier,
		},
	}

	// needed so that we don't hit the default limit of 25 scenes
	pp := 1000
	findFilter := &models.FindFilterType{
		PerPage: &pp,
	}

	return queryScene(ctx, t, queryBuilder, &sceneFilter, findFilter)
}

func createScene(ctx context.Context, width int, height int) (*models.Scene, error) {
	name := fmt.Sprintf("TestSceneQueryResolutionModifiers %d %d", width, height)

	sceneFile := &models.VideoFile{
		BaseFile: &models.BaseFile{
			Basename:       name,
			ParentFolderID: folderIDs[folderIdxWithSceneFiles],
		},
		Width:  width,
		Height: height,
	}

	if err := db.File.Create(ctx, sceneFile); err != nil {
		return nil, err
	}

	scene := &models.Scene{}

	if err := db.Scene.Create(ctx, scene, []models.FileID{sceneFile.ID}); err != nil {
		return nil, err
	}

	return scene, nil
}

func TestSceneQueryHasMarkers(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		hasMarkers := "true"
		sceneFilter := models.SceneFilterType{
			HasMarkers: &hasMarkers,
		}

		q := getSceneStringValue(sceneIdxWithMarkers, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, &findFilter)

		assert.Len(t, scenes, 1)
		assert.Equal(t, sceneIDs[sceneIdxWithMarkers], scenes[0].ID)

		hasMarkers = "false"
		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)
		assert.Len(t, scenes, 0)

		findFilter.Q = nil
		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)

		assert.NotEqual(t, 0, len(scenes))

		// ensure non of the ids equal the one with gallery
		for _, scene := range scenes {
			assert.NotEqual(t, sceneIDs[sceneIdxWithMarkers], scene.ID)
		}

		return nil
	})
}

func TestSceneQueryIsMissingGallery(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		isMissing := "galleries"
		sceneFilter := models.SceneFilterType{
			IsMissing: &isMissing,
		}

		q := getSceneStringValue(sceneIdxWithGallery, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, &findFilter)

		assert.Len(t, scenes, 0)

		findFilter.Q = nil
		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)

		// ensure non of the ids equal the one with gallery
		for _, scene := range scenes {
			assert.NotEqual(t, sceneIDs[sceneIdxWithGallery], scene.ID)
		}

		return nil
	})
}

func TestSceneQueryIsMissingStudio(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		isMissing := "studio"
		sceneFilter := models.SceneFilterType{
			IsMissing: &isMissing,
		}

		q := getSceneStringValue(sceneIdxWithStudio, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, &findFilter)

		assert.Len(t, scenes, 0)

		findFilter.Q = nil
		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)

		// ensure non of the ids equal the one with studio
		for _, scene := range scenes {
			assert.NotEqual(t, sceneIDs[sceneIdxWithStudio], scene.ID)
		}

		return nil
	})
}

func TestSceneQueryIsMissingMovies(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		isMissing := "movie"
		sceneFilter := models.SceneFilterType{
			IsMissing: &isMissing,
		}

		q := getSceneStringValue(sceneIdxWithGroup, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, &findFilter)

		assert.Len(t, scenes, 0)

		findFilter.Q = nil
		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)

		// ensure non of the ids equal the one with movies
		for _, scene := range scenes {
			assert.NotEqual(t, sceneIDs[sceneIdxWithGroup], scene.ID)
		}

		return nil
	})
}

func TestSceneQueryIsMissingPerformers(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		isMissing := "performers"
		sceneFilter := models.SceneFilterType{
			IsMissing: &isMissing,
		}

		q := getSceneStringValue(sceneIdxWithPerformer, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, &findFilter)

		assert.Len(t, scenes, 0)

		findFilter.Q = nil
		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)

		assert.True(t, len(scenes) > 0)

		// ensure non of the ids equal the one with movies
		for _, scene := range scenes {
			assert.NotEqual(t, sceneIDs[sceneIdxWithPerformer], scene.ID)
		}

		return nil
	})
}

func TestSceneQueryIsMissingDate(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		isMissing := "date"
		sceneFilter := models.SceneFilterType{
			IsMissing: &isMissing,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		// one in four scenes have no date
		assert.Len(t, scenes, int(math.Ceil(float64(totalScenes)/4)))

		// ensure date is null
		for _, scene := range scenes {
			assert.Nil(t, scene.Date)
		}

		return nil
	})
}

func TestSceneQueryIsMissingTags(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		isMissing := "tags"
		sceneFilter := models.SceneFilterType{
			IsMissing: &isMissing,
		}

		q := getSceneStringValue(sceneIdxWithTwoTags, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, &findFilter)

		assert.Len(t, scenes, 0)

		findFilter.Q = nil
		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)

		assert.True(t, len(scenes) > 0)

		return nil
	})
}

func TestSceneQueryIsMissingRating(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		isMissing := "rating"
		sceneFilter := models.SceneFilterType{
			IsMissing: &isMissing,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		assert.True(t, len(scenes) > 0)

		// ensure rating is null
		for _, scene := range scenes {
			assert.Nil(t, scene.Rating)
		}

		return nil
	})
}

func TestSceneQueryIsMissingPhash(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		isMissing := "phash"
		sceneFilter := models.SceneFilterType{
			IsMissing: &isMissing,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		if !assert.Len(t, scenes, 1) {
			return nil
		}

		assert.Equal(t, sceneIDs[sceneIdxMissingPhash], scenes[0].ID)

		return nil
	})
}

func TestSceneQueryPerformers(t *testing.T) {
	tests := []struct {
		name        string
		filter      models.MultiCriterionInput
		includeIdxs []int
		excludeIdxs []int
		wantErr     bool
	}{
		{
			"includes",
			models.MultiCriterionInput{
				Value: []string{
					strconv.Itoa(performerIDs[performerIdxWithScene]),
					strconv.Itoa(performerIDs[performerIdx1WithScene]),
				},
				Modifier: models.CriterionModifierIncludes,
			},
			[]int{
				sceneIdxWithPerformer,
				sceneIdxWithTwoPerformers,
			},
			[]int{
				sceneIdxWithGallery,
			},
			false,
		},
		{
			"includes all",
			models.MultiCriterionInput{
				Value: []string{
					strconv.Itoa(performerIDs[performerIdx1WithScene]),
					strconv.Itoa(performerIDs[performerIdx2WithScene]),
				},
				Modifier: models.CriterionModifierIncludesAll,
			},
			[]int{
				sceneIdxWithTwoPerformers,
			},
			[]int{
				sceneIdxWithPerformer,
			},
			false,
		},
		{
			"excludes",
			models.MultiCriterionInput{
				Modifier: models.CriterionModifierExcludes,
				Value:    []string{strconv.Itoa(tagIDs[performerIdx1WithScene])},
			},
			nil,
			[]int{sceneIdxWithTwoPerformers},
			false,
		},
		{
			"is null",
			models.MultiCriterionInput{
				Modifier: models.CriterionModifierIsNull,
			},
			[]int{sceneIdxWithTag},
			[]int{
				sceneIdxWithPerformer,
				sceneIdxWithTwoPerformers,
				sceneIdxWithPerformerTwoTags,
			},
			false,
		},
		{
			"not null",
			models.MultiCriterionInput{
				Modifier: models.CriterionModifierNotNull,
			},
			[]int{
				sceneIdxWithPerformer,
				sceneIdxWithTwoPerformers,
				sceneIdxWithPerformerTwoTags,
			},
			[]int{sceneIdxWithTag},
			false,
		},
		{
			"equals",
			models.MultiCriterionInput{
				Modifier: models.CriterionModifierEquals,
				Value: []string{
					strconv.Itoa(tagIDs[performerIdx1WithScene]),
					strconv.Itoa(tagIDs[performerIdx2WithScene]),
				},
			},
			[]int{sceneIdxWithTwoPerformers},
			[]int{
				sceneIdxWithThreePerformers,
			},
			false,
		},
		{
			"not equals",
			models.MultiCriterionInput{
				Modifier: models.CriterionModifierNotEquals,
				Value: []string{
					strconv.Itoa(tagIDs[performerIdx1WithScene]),
					strconv.Itoa(tagIDs[performerIdx2WithScene]),
				},
			},
			nil,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			results, err := db.Scene.Query(ctx, models.SceneQueryOptions{
				SceneFilter: &models.SceneFilterType{
					Performers: &tt.filter,
				},
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("SceneStore.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			include := indexesToIDs(sceneIDs, tt.includeIdxs)
			exclude := indexesToIDs(sceneIDs, tt.excludeIdxs)

			for _, i := range include {
				assert.Contains(results.IDs, i)
			}
			for _, e := range exclude {
				assert.NotContains(results.IDs, e)
			}
		})
	}
}

func TestSceneQueryTags(t *testing.T) {
	tests := []struct {
		name        string
		filter      models.HierarchicalMultiCriterionInput
		includeIdxs []int
		excludeIdxs []int
		wantErr     bool
	}{
		{
			"includes",
			models.HierarchicalMultiCriterionInput{
				Value: []string{
					strconv.Itoa(tagIDs[tagIdxWithScene]),
					strconv.Itoa(tagIDs[tagIdx1WithScene]),
				},
				Modifier: models.CriterionModifierIncludes,
			},
			[]int{
				sceneIdxWithTag,
				sceneIdxWithTwoTags,
			},
			[]int{
				sceneIdxWithGallery,
			},
			false,
		},
		{
			"includes all",
			models.HierarchicalMultiCriterionInput{
				Value: []string{
					strconv.Itoa(tagIDs[tagIdx1WithScene]),
					strconv.Itoa(tagIDs[tagIdx2WithScene]),
				},
				Modifier: models.CriterionModifierIncludesAll,
			},
			[]int{
				sceneIdxWithTwoTags,
			},
			[]int{
				sceneIdxWithTag,
			},
			false,
		},
		{
			"excludes",
			models.HierarchicalMultiCriterionInput{
				Modifier: models.CriterionModifierExcludes,
				Value:    []string{strconv.Itoa(tagIDs[tagIdx1WithScene])},
			},
			nil,
			[]int{sceneIdxWithTwoTags},
			false,
		},
		{
			"is null",
			models.HierarchicalMultiCriterionInput{
				Modifier: models.CriterionModifierIsNull,
			},
			[]int{sceneIdx1WithPerformer},
			[]int{
				sceneIdxWithTag,
				sceneIdxWithTwoTags,
				sceneIdxWithMarkerAndTag,
			},
			false,
		},
		{
			"not null",
			models.HierarchicalMultiCriterionInput{
				Modifier: models.CriterionModifierNotNull,
			},
			[]int{
				sceneIdxWithTag,
				sceneIdxWithTwoTags,
				sceneIdxWithMarkerAndTag,
			},
			[]int{sceneIdx1WithPerformer},
			false,
		},
		{
			"equals",
			models.HierarchicalMultiCriterionInput{
				Modifier: models.CriterionModifierEquals,
				Value: []string{
					strconv.Itoa(tagIDs[tagIdx1WithScene]),
					strconv.Itoa(tagIDs[tagIdx2WithScene]),
				},
			},
			[]int{sceneIdxWithTwoTags},
			[]int{
				sceneIdxWithThreeTags,
			},
			false,
		},
		{
			"not equals",
			models.HierarchicalMultiCriterionInput{
				Modifier: models.CriterionModifierNotEquals,
				Value: []string{
					strconv.Itoa(tagIDs[tagIdx1WithScene]),
					strconv.Itoa(tagIDs[tagIdx2WithScene]),
				},
			},
			nil,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			results, err := db.Scene.Query(ctx, models.SceneQueryOptions{
				SceneFilter: &models.SceneFilterType{
					Tags: &tt.filter,
				},
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("SceneStore.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			include := indexesToIDs(sceneIDs, tt.includeIdxs)
			exclude := indexesToIDs(sceneIDs, tt.excludeIdxs)

			for _, i := range include {
				assert.Contains(results.IDs, i)
			}
			for _, e := range exclude {
				assert.NotContains(results.IDs, e)
			}
		})
	}
}

func TestSceneQueryPerformerTags(t *testing.T) {
	allDepth := -1

	tests := []struct {
		name        string
		findFilter  *models.FindFilterType
		filter      *models.SceneFilterType
		includeIdxs []int
		excludeIdxs []int
		wantErr     bool
	}{
		{
			"includes",
			nil,
			&models.SceneFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Value: []string{
						strconv.Itoa(tagIDs[tagIdxWithPerformer]),
						strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
					},
					Modifier: models.CriterionModifierIncludes,
				},
			},
			[]int{
				sceneIdxWithPerformerTag,
				sceneIdxWithPerformerTwoTags,
				sceneIdxWithTwoPerformerTag,
			},
			[]int{
				sceneIdxWithPerformer,
			},
			false,
		},
		{
			"includes sub-tags",
			nil,
			&models.SceneFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Value: []string{
						strconv.Itoa(tagIDs[tagIdxWithParentAndChild]),
					},
					Depth:    &allDepth,
					Modifier: models.CriterionModifierIncludes,
				},
			},
			[]int{
				sceneIdxWithPerformerParentTag,
			},
			[]int{
				sceneIdxWithPerformer,
				sceneIdxWithPerformerTag,
				sceneIdxWithPerformerTwoTags,
				sceneIdxWithTwoPerformerTag,
			},
			false,
		},
		{
			"includes all",
			nil,
			&models.SceneFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Value: []string{
						strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
						strconv.Itoa(tagIDs[tagIdx2WithPerformer]),
					},
					Modifier: models.CriterionModifierIncludesAll,
				},
			},
			[]int{
				sceneIdxWithPerformerTwoTags,
			},
			[]int{
				sceneIdxWithPerformer,
				sceneIdxWithPerformerTag,
				sceneIdxWithTwoPerformerTag,
			},
			false,
		},
		{
			"excludes performer tag tagIdx2WithPerformer",
			nil,
			&models.SceneFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Modifier: models.CriterionModifierExcludes,
					Value:    []string{strconv.Itoa(tagIDs[tagIdx2WithPerformer])},
				},
			},
			nil,
			[]int{sceneIdxWithTwoPerformerTag},
			false,
		},
		{
			"excludes sub-tags",
			nil,
			&models.SceneFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Value: []string{
						strconv.Itoa(tagIDs[tagIdxWithParentAndChild]),
					},
					Depth:    &allDepth,
					Modifier: models.CriterionModifierExcludes,
				},
			},
			[]int{
				sceneIdxWithPerformer,
				sceneIdxWithPerformerTag,
				sceneIdxWithPerformerTwoTags,
				sceneIdxWithTwoPerformerTag,
			},
			[]int{
				sceneIdxWithPerformerParentTag,
			},
			false,
		},
		{
			"is null",
			nil,
			&models.SceneFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Modifier: models.CriterionModifierIsNull,
				},
			},
			[]int{sceneIdx1WithPerformer},
			[]int{sceneIdxWithPerformerTag},
			false,
		},
		{
			"not null",
			nil,
			&models.SceneFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Modifier: models.CriterionModifierNotNull,
				},
			},
			[]int{sceneIdxWithPerformerTag},
			[]int{sceneIdx1WithPerformer},
			false,
		},
		{
			"equals",
			nil,
			&models.SceneFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Modifier: models.CriterionModifierEquals,
					Value: []string{
						strconv.Itoa(tagIDs[tagIdx2WithPerformer]),
					},
				},
			},
			nil,
			nil,
			true,
		},
		{
			"not equals",
			nil,
			&models.SceneFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Modifier: models.CriterionModifierNotEquals,
					Value: []string{
						strconv.Itoa(tagIDs[tagIdx2WithPerformer]),
					},
				},
			},
			nil,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			results, err := db.Scene.Query(ctx, models.SceneQueryOptions{
				SceneFilter: tt.filter,
				QueryOptions: models.QueryOptions{
					FindFilter: tt.findFilter,
				},
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("SceneStore.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			include := indexesToIDs(sceneIDs, tt.includeIdxs)
			exclude := indexesToIDs(sceneIDs, tt.excludeIdxs)

			for _, i := range include {
				assert.Contains(results.IDs, i)
			}
			for _, e := range exclude {
				assert.NotContains(results.IDs, e)
			}
		})
	}
}

func TestSceneQueryStudio(t *testing.T) {
	tests := []struct {
		name            string
		q               string
		studioCriterion models.HierarchicalMultiCriterionInput
		expectedIDs     []int
		wantErr         bool
	}{
		{
			"includes",
			"",
			models.HierarchicalMultiCriterionInput{
				Value: []string{
					strconv.Itoa(studioIDs[studioIdxWithScene]),
				},
				Modifier: models.CriterionModifierIncludes,
			},
			[]int{sceneIDs[sceneIdxWithStudio]},
			false,
		},
		{
			"excludes",
			getSceneStringValue(sceneIdxWithStudio, titleField),
			models.HierarchicalMultiCriterionInput{
				Value: []string{
					strconv.Itoa(studioIDs[studioIdxWithScene]),
				},
				Modifier: models.CriterionModifierExcludes,
			},
			[]int{},
			false,
		},
		{
			"excludes includes null",
			getSceneStringValue(sceneIdxWithGallery, titleField),
			models.HierarchicalMultiCriterionInput{
				Value: []string{
					strconv.Itoa(studioIDs[studioIdxWithScene]),
				},
				Modifier: models.CriterionModifierExcludes,
			},
			[]int{sceneIDs[sceneIdxWithGallery]},
			false,
		},
		{
			"equals",
			"",
			models.HierarchicalMultiCriterionInput{
				Value: []string{
					strconv.Itoa(studioIDs[studioIdxWithScene]),
				},
				Modifier: models.CriterionModifierEquals,
			},
			[]int{sceneIDs[sceneIdxWithStudio]},
			false,
		},
		{
			"not equals",
			getSceneStringValue(sceneIdxWithStudio, titleField),
			models.HierarchicalMultiCriterionInput{
				Value: []string{
					strconv.Itoa(studioIDs[studioIdxWithScene]),
				},
				Modifier: models.CriterionModifierNotEquals,
			},
			[]int{},
			false,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			studioCriterion := tt.studioCriterion

			sceneFilter := models.SceneFilterType{
				Studios: &studioCriterion,
			}

			var findFilter *models.FindFilterType
			if tt.q != "" {
				findFilter = &models.FindFilterType{
					Q: &tt.q,
				}
			}

			scenes := queryScene(ctx, t, qb, &sceneFilter, findFilter)

			assert.ElementsMatch(t, scenesToIDs(scenes), tt.expectedIDs)
		})
	}
}

func TestSceneQueryStudioDepth(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		depth := 2
		studioCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithGrandChild]),
			},
			Modifier: models.CriterionModifierIncludes,
			Depth:    &depth,
		}

		sceneFilter := models.SceneFilterType{
			Studios: &studioCriterion,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)
		assert.Len(t, scenes, 1)

		depth = 1

		scenes = queryScene(ctx, t, sqb, &sceneFilter, nil)
		assert.Len(t, scenes, 0)

		studioCriterion.Value = []string{strconv.Itoa(studioIDs[studioIdxWithParentAndChild])}
		scenes = queryScene(ctx, t, sqb, &sceneFilter, nil)
		assert.Len(t, scenes, 1)

		// ensure id is correct
		assert.Equal(t, sceneIDs[sceneIdxWithGrandChildStudio], scenes[0].ID)
		depth = 2

		studioCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithGrandChild]),
			},
			Modifier: models.CriterionModifierExcludes,
			Depth:    &depth,
		}

		q := getSceneStringValue(sceneIdxWithGrandChildStudio, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)
		assert.Len(t, scenes, 0)

		depth = 1
		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)
		assert.Len(t, scenes, 1)

		studioCriterion.Value = []string{strconv.Itoa(studioIDs[studioIdxWithParentAndChild])}
		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)
		assert.Len(t, scenes, 0)

		return nil
	})
}

func TestSceneQueryMovies(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		movieCriterion := models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(groupIDs[groupIdxWithScene]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		sceneFilter := models.SceneFilterType{
			Movies: &movieCriterion,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		assert.Len(t, scenes, 1)

		// ensure id is correct
		assert.Equal(t, sceneIDs[sceneIdxWithGroup], scenes[0].ID)

		movieCriterion = models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(groupIDs[groupIdxWithScene]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getSceneStringValue(sceneIdxWithGroup, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)
		assert.Len(t, scenes, 0)

		return nil
	})
}

func TestSceneQueryPhashDuplicated(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		duplicated := true
		phashCriterion := models.PHashDuplicationCriterionInput{
			Duplicated: &duplicated,
		}

		sceneFilter := models.SceneFilterType{
			Duplicated: &phashCriterion,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		assert.Len(t, scenes, dupeScenePhashes*2)

		duplicated = false

		scenes = queryScene(ctx, t, sqb, &sceneFilter, nil)
		// -1 for missing phash
		assert.Len(t, scenes, totalScenes-(dupeScenePhashes*2)-1)

		return nil
	})
}

func TestSceneQuerySorting(t *testing.T) {
	tests := []struct {
		name          string
		sortBy        string
		dir           models.SortDirectionEnum
		firstSceneIdx int // -1 to ignore
		lastSceneIdx  int
	}{
		{
			"bitrate",
			"bitrate",
			models.SortDirectionEnumAsc,
			-1,
			-1,
		},
		{
			"duration",
			"duration",
			models.SortDirectionEnumDesc,
			-1,
			-1,
		},
		{
			"file mod time",
			"file_mod_time",
			models.SortDirectionEnumDesc,
			-1,
			-1,
		},
		{
			"file size",
			"filesize",
			models.SortDirectionEnumDesc,
			-1,
			-1,
		},
		{
			"frame rate",
			"framerate",
			models.SortDirectionEnumDesc,
			-1,
			-1,
		},
		{
			"path",
			"path",
			models.SortDirectionEnumDesc,
			-1,
			-1,
		},
		{
			"perceptual_similarity",
			"perceptual_similarity",
			models.SortDirectionEnumDesc,
			-1,
			-1,
		},
		{
			"play_count",
			"play_count",
			models.SortDirectionEnumDesc,
			-1,
			-1,
		},
		{
			"last_played_at",
			"last_played_at",
			models.SortDirectionEnumDesc,
			-1,
			-1,
		},
		{
			"resume_time",
			"resume_time",
			models.SortDirectionEnumDesc,
			sceneIDs[sceneIdx1WithPerformer],
			-1,
		},
		{
			"play_duration",
			"play_duration",
			models.SortDirectionEnumDesc,
			sceneIDs[sceneIdx1WithPerformer],
			-1,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.Query(ctx, models.SceneQueryOptions{
				QueryOptions: models.QueryOptions{
					FindFilter: &models.FindFilterType{
						Sort:      &tt.sortBy,
						Direction: &tt.dir,
					},
				},
			})

			if err != nil {
				t.Errorf("sceneQueryBuilder.TestSceneQuerySorting() error = %v", err)
				return
			}

			scenes, err := got.Resolve(ctx)
			if err != nil {
				t.Errorf("sceneQueryBuilder.TestSceneQuerySorting() error = %v", err)
				return
			}

			if !assert.Greater(len(scenes), 0) {
				return
			}

			// scenes should be in same order as indexes
			firstScene := scenes[0]
			lastScene := scenes[len(scenes)-1]

			if tt.firstSceneIdx != -1 {
				firstSceneID := sceneIDs[tt.firstSceneIdx]
				assert.Equal(firstSceneID, firstScene.ID)
			}
			if tt.lastSceneIdx != -1 {
				lastSceneID := sceneIDs[tt.lastSceneIdx]
				assert.Equal(lastSceneID, lastScene.ID)
			}
		})
	}
}

func TestSceneQueryPagination(t *testing.T) {
	perPage := 1
	findFilter := models.FindFilterType{
		PerPage: &perPage,
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		scenes := queryScene(ctx, t, sqb, nil, &findFilter)

		assert.Len(t, scenes, 1)

		firstID := scenes[0].ID

		page := 2
		findFilter.Page = &page
		scenes = queryScene(ctx, t, sqb, nil, &findFilter)

		assert.Len(t, scenes, 1)
		secondID := scenes[0].ID
		assert.NotEqual(t, firstID, secondID)

		perPage = 2
		page = 1

		scenes = queryScene(ctx, t, sqb, nil, &findFilter)
		assert.Len(t, scenes, 2)
		assert.Equal(t, firstID, scenes[0].ID)
		assert.Equal(t, secondID, scenes[1].ID)

		return nil
	})
}

func TestSceneQueryTagCount(t *testing.T) {
	const tagCount = 1
	tagCountCriterion := models.IntCriterionInput{
		Value:    tagCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyScenesTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyScenesTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyScenesTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyScenesTagCount(t, tagCountCriterion)
}

func verifyScenesTagCount(t *testing.T, tagCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		sceneFilter := models.SceneFilterType{
			TagCount: &tagCountCriterion,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)
		assert.Greater(t, len(scenes), 0)

		for _, scene := range scenes {
			if err := scene.LoadTagIDs(ctx, sqb); err != nil {
				t.Errorf("scene.LoadTagIDs() error = %v", err)
				return nil
			}
			verifyInt(t, len(scene.TagIDs.List()), tagCountCriterion)
		}

		return nil
	})
}

func TestSceneQueryPerformerCount(t *testing.T) {
	const performerCount = 1
	performerCountCriterion := models.IntCriterionInput{
		Value:    performerCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyScenesPerformerCount(t, performerCountCriterion)

	performerCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyScenesPerformerCount(t, performerCountCriterion)

	performerCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyScenesPerformerCount(t, performerCountCriterion)

	performerCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyScenesPerformerCount(t, performerCountCriterion)
}

func verifyScenesPerformerCount(t *testing.T, performerCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		sceneFilter := models.SceneFilterType{
			PerformerCount: &performerCountCriterion,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)
		assert.Greater(t, len(scenes), 0)

		for _, scene := range scenes {
			if err := scene.LoadPerformerIDs(ctx, sqb); err != nil {
				t.Errorf("scene.LoadPerformerIDs() error = %v", err)
				return nil
			}

			verifyInt(t, len(scene.PerformerIDs.List()), performerCountCriterion)
		}

		return nil
	})
}

func TestSceneCountByTagID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene

		sceneCount, err := sqb.CountByTagID(ctx, tagIDs[tagIdxWithScene])

		if err != nil {
			t.Errorf("error calling CountByTagID: %s", err.Error())
		}

		assert.Equal(t, 1, sceneCount)

		sceneCount, err = sqb.CountByTagID(ctx, 0)

		if err != nil {
			t.Errorf("error calling CountByTagID: %s", err.Error())
		}

		assert.Equal(t, 0, sceneCount)

		return nil
	})
}

func TestSceneCountByGroupID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene

		sceneCount, err := sqb.CountByGroupID(ctx, groupIDs[groupIdxWithScene])

		if err != nil {
			t.Errorf("error calling CountByGroupID: %s", err.Error())
		}

		assert.Equal(t, 1, sceneCount)

		sceneCount, err = sqb.CountByGroupID(ctx, 0)

		if err != nil {
			t.Errorf("error calling CountByGroupID: %s", err.Error())
		}

		assert.Equal(t, 0, sceneCount)

		return nil
	})
}

func TestSceneCountByStudioID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene

		sceneCount, err := sqb.CountByStudioID(ctx, studioIDs[studioIdxWithScene])

		if err != nil {
			t.Errorf("error calling CountByStudioID: %s", err.Error())
		}

		assert.Equal(t, 1, sceneCount)

		sceneCount, err = sqb.CountByStudioID(ctx, 0)

		if err != nil {
			t.Errorf("error calling CountByStudioID: %s", err.Error())
		}

		assert.Equal(t, 0, sceneCount)

		return nil
	})
}

func TestFindByMovieID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene

		scenes, err := sqb.FindByGroupID(ctx, groupIDs[groupIdxWithScene])

		if err != nil {
			t.Errorf("error calling FindByMovieID: %s", err.Error())
		}

		assert.Len(t, scenes, 1)
		assert.Equal(t, sceneIDs[sceneIdxWithGroup], scenes[0].ID)

		scenes, err = sqb.FindByGroupID(ctx, 0)

		if err != nil {
			t.Errorf("error calling FindByMovieID: %s", err.Error())
		}

		assert.Len(t, scenes, 0)

		return nil
	})
}

func TestFindByPerformerID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene

		scenes, err := sqb.FindByPerformerID(ctx, performerIDs[performerIdxWithScene])

		if err != nil {
			t.Errorf("error calling FindByPerformerID: %s", err.Error())
		}

		assert.Len(t, scenes, 1)
		assert.Equal(t, sceneIDs[sceneIdxWithPerformer], scenes[0].ID)

		scenes, err = sqb.FindByPerformerID(ctx, 0)

		if err != nil {
			t.Errorf("error calling FindByPerformerID: %s", err.Error())
		}

		assert.Len(t, scenes, 0)

		return nil
	})
}

func TestSceneUpdateSceneCover(t *testing.T) {
	if err := withTxn(func(ctx context.Context) error {
		qb := db.Scene

		sceneID := sceneIDs[sceneIdxWithGallery]

		return testUpdateImage(t, ctx, sceneID, qb.UpdateCover, qb.GetCover)
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestSceneStashIDs(t *testing.T) {
	if err := withTxn(func(ctx context.Context) error {
		qb := db.Scene

		// create scene to test against
		const name = "TestSceneStashIDs"
		scene := &models.Scene{
			Title: name,
		}
		if err := qb.Create(ctx, scene, nil); err != nil {
			return fmt.Errorf("Error creating scene: %s", err.Error())
		}

		if err := scene.LoadStashIDs(ctx, qb); err != nil {
			return err
		}

		testSceneStashIDs(ctx, t, scene)
		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func testSceneStashIDs(ctx context.Context, t *testing.T, s *models.Scene) {
	// ensure no stash IDs to begin with
	assert.Len(t, s.StashIDs.List(), 0)

	// add stash ids
	const stashIDStr = "stashID"
	const endpoint = "endpoint"
	stashID := models.StashID{
		StashID:  stashIDStr,
		Endpoint: endpoint,
	}

	qb := db.Scene

	// update stash ids and ensure was updated
	var err error
	s, err = qb.UpdatePartial(ctx, s.ID, models.ScenePartial{
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
	s, err = qb.UpdatePartial(ctx, s.ID, models.ScenePartial{
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

func TestSceneQueryQTrim(t *testing.T) {
	if err := withTxn(func(ctx context.Context) error {
		qb := db.Scene

		expectedID := sceneIDs[sceneIdxWithSpacedName]

		type test struct {
			query string
			id    int
			count int
		}
		tests := []test{
			{query: " zzz    yyy    ", id: expectedID, count: 1},
			{query: "   \"zzz yyy xxx\" ", id: expectedID, count: 1},
			{query: "zzz", id: expectedID, count: 1},
			{query: "\" zzz    yyy    \"", count: 0},
			{query: "\"zzz    yyy\"", count: 0},
			{query: "\" zzz yyy\"", count: 0},
			{query: "\"zzz yyy  \"", count: 0},
		}

		for _, tst := range tests {
			f := models.FindFilterType{
				Q: &tst.query,
			}
			scenes := queryScene(ctx, t, qb, nil, &f)

			assert.Len(t, scenes, tst.count)
			if len(scenes) > 0 {
				assert.Equal(t, tst.id, scenes[0].ID)
			}
		}

		findFilter := models.FindFilterType{}
		scenes := queryScene(ctx, t, qb, nil, &findFilter)
		assert.NotEqual(t, 0, len(scenes))

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestSceneStore_All(t *testing.T) {
	qb := db.Scene

	withRollbackTxn(func(ctx context.Context) error {
		got, err := qb.All(ctx)
		if err != nil {
			t.Errorf("SceneStore.All() error = %v", err)
			return nil
		}

		// it's possible that other tests have created scenes
		assert.GreaterOrEqual(t, len(got), len(sceneIDs))

		return nil
	})
}

func TestSceneStore_FindDuplicates(t *testing.T) {
	qb := db.Scene

	withRollbackTxn(func(ctx context.Context) error {
		distance := 0
		durationDiff := -1.
		got, err := qb.FindDuplicates(ctx, distance, durationDiff)
		if err != nil {
			t.Errorf("SceneStore.FindDuplicates() error = %v", err)
			return nil
		}

		assert.Len(t, got, dupeScenePhashes)

		distance = 1
		durationDiff = -1.
		got, err = qb.FindDuplicates(ctx, distance, durationDiff)
		if err != nil {
			t.Errorf("SceneStore.FindDuplicates() error = %v", err)
			return nil
		}

		assert.Len(t, got, dupeScenePhashes)

		return nil
	})
}

func TestSceneStore_AssignFiles(t *testing.T) {
	tests := []struct {
		name    string
		sceneID int
		fileID  models.FileID
		wantErr bool
	}{
		{
			"valid",
			sceneIDs[sceneIdx1WithPerformer],
			sceneFileIDs[sceneIdx1WithStudio],
			false,
		},
		{
			"invalid file id",
			sceneIDs[sceneIdx1WithPerformer],
			invalidFileID,
			true,
		},
		{
			"invalid scene id",
			invalidID,
			sceneFileIDs[sceneIdx1WithStudio],
			true,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withRollbackTxn(func(ctx context.Context) error {
				if err := qb.AssignFiles(ctx, tt.sceneID, []models.FileID{tt.fileID}); (err != nil) != tt.wantErr {
					t.Errorf("SceneStore.AssignFiles() error = %v, wantErr %v", err, tt.wantErr)
				}

				return nil
			})
		})
	}
}

func TestSceneStore_AddView(t *testing.T) {
	tests := []struct {
		name          string
		sceneID       int
		expectedCount int
		wantErr       bool
	}{
		{
			"valid",
			sceneIDs[sceneIdx1WithPerformer],
			1, //getScenePlayCount(sceneIdx1WithPerformer) + 1,
			false,
		},
		{
			"invalid scene id",
			invalidID,
			0,
			true,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withRollbackTxn(func(ctx context.Context) error {
				views, err := qb.AddViews(ctx, tt.sceneID, nil)
				if (err != nil) != tt.wantErr {
					t.Errorf("SceneStore.AddView() error = %v, wantErr %v", err, tt.wantErr)
				}

				if err != nil {
					return nil
				}

				assert := assert.New(t)
				assert.Equal(tt.expectedCount, len(views))

				// find the scene and check the count
				count, err := qb.CountViews(ctx, tt.sceneID)
				if err != nil {
					t.Errorf("SceneStore.CountViews() error = %v", err)
				}

				lastView, err := qb.LastView(ctx, tt.sceneID)
				if err != nil {
					t.Errorf("SceneStore.LastView() error = %v", err)
				}

				assert.Equal(tt.expectedCount, count)
				assert.True(lastView.After(time.Now().Add(-1 * time.Minute)))

				return nil
			})
		})
	}
}

func TestSceneStore_DecrementWatchCount(t *testing.T) {
	return
}

func TestSceneStore_SaveActivity(t *testing.T) {
	var (
		resumeTime   = 111.2
		playDuration = 98.7
	)

	tests := []struct {
		name         string
		sceneIdx     int
		resumeTime   *float64
		playDuration *float64
		wantErr      bool
	}{
		{
			"both",
			sceneIdx1WithPerformer,
			&resumeTime,
			&playDuration,
			false,
		},
		{
			"resumeTime only",
			sceneIdx1WithPerformer,
			&resumeTime,
			nil,
			false,
		},
		{
			"playDuration only",
			sceneIdx1WithPerformer,
			nil,
			&playDuration,
			false,
		},
		{
			"none",
			sceneIdx1WithPerformer,
			nil,
			nil,
			false,
		},
		{
			"invalid scene id",
			-1,
			&resumeTime,
			&playDuration,
			true,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withRollbackTxn(func(ctx context.Context) error {
				id := -1
				if tt.sceneIdx != -1 {
					id = sceneIDs[tt.sceneIdx]
				}

				_, err := qb.SaveActivity(ctx, id, tt.resumeTime, tt.playDuration)
				if (err != nil) != tt.wantErr {
					t.Errorf("SceneStore.SaveActivity() error = %v, wantErr %v", err, tt.wantErr)
				}

				if err != nil {
					return nil
				}

				assert := assert.New(t)

				// find the scene and check the values
				scene, err := qb.Find(ctx, id)
				if err != nil {
					t.Errorf("SceneStore.Find() error = %v", err)
				}

				expectedResumeTime := getSceneResumeTime(tt.sceneIdx)
				expectedPlayDuration := getScenePlayDuration(tt.sceneIdx)

				if tt.resumeTime != nil {
					expectedResumeTime = *tt.resumeTime
				}
				if tt.playDuration != nil {
					expectedPlayDuration += *tt.playDuration
				}

				assert.Equal(expectedResumeTime, scene.ResumeTime)
				assert.Equal(expectedPlayDuration, scene.PlayDuration)

				return nil
			})
		})
	}
}

// TODO Count
// TODO SizeCount

// TODO - this should be in history_test and generalised
func TestSceneStore_CountAllViews(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		qb := db.Scene

		sceneID := sceneIDs[sceneIdx1WithPerformer]

		// get the current play count
		currentCount, err := qb.CountAllViews(ctx)
		if err != nil {
			t.Errorf("SceneStore.CountAllViews() error = %v", err)
			return nil
		}

		// add a view
		_, err = qb.AddViews(ctx, sceneID, nil)
		if err != nil {
			t.Errorf("SceneStore.AddViews() error = %v", err)
			return nil
		}

		// get the new play count
		newCount, err := qb.CountAllViews(ctx)
		if err != nil {
			t.Errorf("SceneStore.CountAllViews() error = %v", err)
			return nil
		}

		assert.Equal(t, currentCount+1, newCount)

		return nil
	})
}

func TestSceneStore_CountUniqueViews(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		qb := db.Scene

		sceneID := sceneIDs[sceneIdx1WithPerformer]

		// get the current play count
		currentCount, err := qb.CountUniqueViews(ctx)
		if err != nil {
			t.Errorf("SceneStore.CountUniqueViews() error = %v", err)
			return nil
		}

		// add a view
		_, err = qb.AddViews(ctx, sceneID, nil)
		if err != nil {
			t.Errorf("SceneStore.AddViews() error = %v", err)
			return nil
		}

		// add a second view
		_, err = qb.AddViews(ctx, sceneID, nil)
		if err != nil {
			t.Errorf("SceneStore.AddViews() error = %v", err)
			return nil
		}

		// get the new play count
		newCount, err := qb.CountUniqueViews(ctx)
		if err != nil {
			t.Errorf("SceneStore.CountUniqueViews() error = %v", err)
			return nil
		}

		assert.Equal(t, currentCount+1, newCount)

		return nil
	})
}
