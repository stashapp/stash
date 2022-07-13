//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stretchr/testify/assert"
)

func Test_sceneQueryBuilder_Create(t *testing.T) {
	var (
		title       = "title"
		details     = "details"
		url         = "url"
		rating      = 3
		ocounter    = 5
		createdAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		sceneIndex  = 123
		sceneIndex2 = 234
		endpoint1   = "endpoint1"
		endpoint2   = "endpoint2"
		stashID1    = "stashid1"
		stashID2    = "stashid2"

		date = models.NewDate("2003-02-01")

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
				Details:      details,
				URL:          url,
				Date:         &date,
				Rating:       &rating,
				Organized:    true,
				OCounter:     ocounter,
				StudioID:     &studioIDs[studioIdxWithScene],
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				GalleryIDs:   []int{galleryIDs[galleryIdxWithScene]},
				TagIDs:       []int{tagIDs[tagIdx1WithScene], tagIDs[tagIdx1WithDupName]},
				PerformerIDs: []int{performerIDs[performerIdx1WithScene], performerIDs[performerIdx1WithDupName]},
				Movies: []models.MoviesScenes{
					{
						MovieID:    movieIDs[movieIdxWithScene],
						SceneIndex: &sceneIndex,
					},
					{
						MovieID:    movieIDs[movieIdxWithStudio],
						SceneIndex: &sceneIndex2,
					},
				},
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
			},
			false,
		},
		{
			"with file",
			models.Scene{
				Title:     title,
				Details:   details,
				URL:       url,
				Date:      &date,
				Rating:    &rating,
				Organized: true,
				OCounter:  ocounter,
				StudioID:  &studioIDs[studioIdxWithScene],
				Files: []*file.VideoFile{
					videoFile.(*file.VideoFile),
				},
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				GalleryIDs:   []int{galleryIDs[galleryIdxWithScene]},
				TagIDs:       []int{tagIDs[tagIdx1WithScene], tagIDs[tagIdx1WithDupName]},
				PerformerIDs: []int{performerIDs[performerIdx1WithScene], performerIDs[performerIdx1WithDupName]},
				Movies: []models.MoviesScenes{
					{
						MovieID:    movieIDs[movieIdxWithScene],
						SceneIndex: &sceneIndex,
					},
					{
						MovieID:    movieIDs[movieIdxWithStudio],
						SceneIndex: &sceneIndex2,
					},
				},
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
				GalleryIDs: []int{invalidID},
			},
			true,
		},
		{
			"invalid tag id",
			models.Scene{
				TagIDs: []int{invalidID},
			},
			true,
		},
		{
			"invalid performer id",
			models.Scene{
				PerformerIDs: []int{invalidID},
			},
			true,
		},
		{
			"invalid movie id",
			models.Scene{
				Movies: []models.MoviesScenes{
					{
						MovieID:    invalidID,
						SceneIndex: &sceneIndex,
					},
				},
			},
			true,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			var fileIDs []file.ID
			for _, f := range tt.newObject.Files {
				fileIDs = append(fileIDs, f.ID)
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

			assert.Equal(copy, s)

			// ensure can find the scene
			found, err := qb.Find(ctx, s.ID)
			if err != nil {
				t.Errorf("sceneQueryBuilder.Find() error = %v", err)
			}

			if !assert.NotNil(found) {
				return
			}
			assert.Equal(copy, *found)

			return
		})
	}
}

func clearSceneFileIDs(scene *models.Scene) {
	for _, f := range scene.Files {
		f.Base().ID = 0
	}
}

func makeSceneFileWithID(i int) *file.VideoFile {
	ret := makeSceneFile(i)
	ret.ID = sceneFileIDs[i]
	return ret
}

func Test_sceneQueryBuilder_Update(t *testing.T) {
	var (
		title       = "title"
		details     = "details"
		url         = "url"
		rating      = 3
		ocounter    = 5
		createdAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		sceneIndex  = 123
		sceneIndex2 = 234
		endpoint1   = "endpoint1"
		endpoint2   = "endpoint2"
		stashID1    = "stashid1"
		stashID2    = "stashid2"

		date = models.NewDate("2003-02-01")
	)

	tests := []struct {
		name          string
		updatedObject *models.Scene
		wantErr       bool
	}{
		{
			"full",
			&models.Scene{
				ID: sceneIDs[sceneIdxWithGallery],
				Files: []*file.VideoFile{
					makeSceneFileWithID(sceneIdxWithGallery),
				},
				Title:        title,
				Details:      details,
				URL:          url,
				Date:         &date,
				Rating:       &rating,
				Organized:    true,
				OCounter:     ocounter,
				StudioID:     &studioIDs[studioIdxWithScene],
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				GalleryIDs:   []int{galleryIDs[galleryIdxWithScene]},
				TagIDs:       []int{tagIDs[tagIdx1WithScene], tagIDs[tagIdx1WithDupName]},
				PerformerIDs: []int{performerIDs[performerIdx1WithScene], performerIDs[performerIdx1WithDupName]},
				Movies: []models.MoviesScenes{
					{
						MovieID:    movieIDs[movieIdxWithScene],
						SceneIndex: &sceneIndex,
					},
					{
						MovieID:    movieIDs[movieIdxWithStudio],
						SceneIndex: &sceneIndex2,
					},
				},
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
			},
			false,
		},
		{
			"clear nullables",
			&models.Scene{
				ID: sceneIDs[sceneIdxWithSpacedName],
				Files: []*file.VideoFile{
					makeSceneFileWithID(sceneIdxWithSpacedName),
				},
			},
			false,
		},
		{
			"clear gallery ids",
			&models.Scene{
				ID: sceneIDs[sceneIdxWithGallery],
				Files: []*file.VideoFile{
					makeSceneFileWithID(sceneIdxWithGallery),
				},
			},
			false,
		},
		{
			"clear tag ids",
			&models.Scene{
				ID: sceneIDs[sceneIdxWithTag],
				Files: []*file.VideoFile{
					makeSceneFileWithID(sceneIdxWithTag),
				},
			},
			false,
		},
		{
			"clear performer ids",
			&models.Scene{
				ID: sceneIDs[sceneIdxWithPerformer],
				Files: []*file.VideoFile{
					makeSceneFileWithID(sceneIdxWithPerformer),
				},
			},
			false,
		},
		{
			"clear movies",
			&models.Scene{
				ID: sceneIDs[sceneIdxWithMovie],
				Files: []*file.VideoFile{
					makeSceneFileWithID(sceneIdxWithMovie),
				},
			},
			false,
		},
		{
			"invalid studio id",
			&models.Scene{
				ID: sceneIDs[sceneIdxWithGallery],
				Files: []*file.VideoFile{
					makeSceneFileWithID(sceneIdxWithGallery),
				},
				StudioID: &invalidID,
			},
			true,
		},
		{
			"invalid gallery id",
			&models.Scene{
				ID: sceneIDs[sceneIdxWithGallery],
				Files: []*file.VideoFile{
					makeSceneFileWithID(sceneIdxWithGallery),
				},
				GalleryIDs: []int{invalidID},
			},
			true,
		},
		{
			"invalid tag id",
			&models.Scene{
				ID: sceneIDs[sceneIdxWithGallery],
				Files: []*file.VideoFile{
					makeSceneFileWithID(sceneIdxWithGallery),
				},
				TagIDs: []int{invalidID},
			},
			true,
		},
		{
			"invalid performer id",
			&models.Scene{
				ID: sceneIDs[sceneIdxWithGallery],
				Files: []*file.VideoFile{
					makeSceneFileWithID(sceneIdxWithGallery),
				},
				PerformerIDs: []int{invalidID},
			},
			true,
		},
		{
			"invalid movie id",
			&models.Scene{
				ID: sceneIDs[sceneIdxWithSpacedName],
				Files: []*file.VideoFile{
					makeSceneFileWithID(sceneIdxWithSpacedName),
				},
				Movies: []models.MoviesScenes{
					{
						MovieID:    invalidID,
						SceneIndex: &sceneIndex,
					},
				},
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

			assert.Equal(copy, *s)
		})
	}
}

func clearScenePartial() models.ScenePartial {
	// leave mandatory fields
	return models.ScenePartial{
		Title:        models.OptionalString{Set: true, Null: true},
		Details:      models.OptionalString{Set: true, Null: true},
		URL:          models.OptionalString{Set: true, Null: true},
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
		title       = "title"
		details     = "details"
		url         = "url"
		rating      = 3
		ocounter    = 5
		createdAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		sceneIndex  = 123
		sceneIndex2 = 234
		endpoint1   = "endpoint1"
		endpoint2   = "endpoint2"
		stashID1    = "stashid1"
		stashID2    = "stashid2"

		date = models.NewDate("2003-02-01")
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
				Title:     models.NewOptionalString(title),
				Details:   models.NewOptionalString(details),
				URL:       models.NewOptionalString(url),
				Date:      models.NewOptionalDate(date),
				Rating:    models.NewOptionalInt(rating),
				Organized: models.NewOptionalBool(true),
				OCounter:  models.NewOptionalInt(ocounter),
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
				MovieIDs: &models.UpdateMovieIDs{
					Movies: []models.MoviesScenes{
						{
							MovieID:    movieIDs[movieIdxWithScene],
							SceneIndex: &sceneIndex,
						},
						{
							MovieID:    movieIDs[movieIdxWithStudio],
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
			},
			models.Scene{
				ID: sceneIDs[sceneIdxWithSpacedName],
				Files: []*file.VideoFile{
					makeSceneFile(sceneIdxWithSpacedName),
				},
				Title:        title,
				Details:      details,
				URL:          url,
				Date:         &date,
				Rating:       &rating,
				Organized:    true,
				OCounter:     ocounter,
				StudioID:     &studioIDs[studioIdxWithScene],
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				GalleryIDs:   []int{galleryIDs[galleryIdxWithScene]},
				TagIDs:       []int{tagIDs[tagIdx1WithScene], tagIDs[tagIdx1WithDupName]},
				PerformerIDs: []int{performerIDs[performerIdx1WithScene], performerIDs[performerIdx1WithDupName]},
				Movies: []models.MoviesScenes{
					{
						MovieID:    movieIDs[movieIdxWithScene],
						SceneIndex: &sceneIndex,
					},
					{
						MovieID:    movieIDs[movieIdxWithStudio],
						SceneIndex: &sceneIndex2,
					},
				},
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
			},
			false,
		},
		{
			"clear all",
			sceneIDs[sceneIdxWithSpacedName],
			clearScenePartial(),
			models.Scene{
				ID: sceneIDs[sceneIdxWithSpacedName],
				Files: []*file.VideoFile{
					makeSceneFile(sceneIdxWithSpacedName),
				},
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

			// ignore file ids
			clearSceneFileIDs(got)

			assert.Equal(tt.want, *got)

			s, err := qb.Find(ctx, tt.id)
			if err != nil {
				t.Errorf("sceneQueryBuilder.Find() error = %v", err)
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

		movieScenes = []models.MoviesScenes{
			{
				MovieID:    movieIDs[movieIdxWithDupName],
				SceneIndex: &sceneIndex,
			},
			{
				MovieID:    movieIDs[movieIdxWithStudio],
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
				GalleryIDs: append(indexesToIDs(galleryIDs, sceneGalleries[sceneIdxWithGallery]),
					galleryIDs[galleryIdx1WithImage],
					galleryIDs[galleryIdx1WithPerformer],
				),
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
				TagIDs: append(indexesToIDs(tagIDs, sceneTags[sceneIdxWithTwoTags]),
					tagIDs[tagIdx1WithDupName],
					tagIDs[tagIdx1WithGallery],
				),
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
				PerformerIDs: append(indexesToIDs(performerIDs, scenePerformers[sceneIdxWithTwoPerformers]),
					performerIDs[performerIdx1WithDupName],
					performerIDs[performerIdx1WithGallery],
				),
			},
			false,
		},
		{
			"add movies",
			sceneIDs[sceneIdxWithMovie],
			models.ScenePartial{
				MovieIDs: &models.UpdateMovieIDs{
					Movies: movieScenes,
					Mode:   models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				Movies: append([]models.MoviesScenes{
					{
						MovieID: indexesToIDs(movieIDs, sceneMovies[sceneIdxWithMovie])[0],
					},
				}, movieScenes...),
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
				StashIDs: append(stashIDs, []models.StashID{sceneStashID(sceneIdxWithSpacedName)}...),
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
				GalleryIDs: append(indexesToIDs(galleryIDs, sceneGalleries[sceneIdxWithGallery]),
					galleryIDs[galleryIdx1WithPerformer],
				),
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
				TagIDs: append(indexesToIDs(tagIDs, sceneTags[sceneIdxWithTwoTags]),
					tagIDs[tagIdx1WithGallery],
				),
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
				PerformerIDs: append(indexesToIDs(performerIDs, scenePerformers[sceneIdxWithTwoPerformers]),
					performerIDs[performerIdx1WithGallery],
				),
			},
			false,
		},
		{
			"add duplicate movies",
			sceneIDs[sceneIdxWithMovie],
			models.ScenePartial{
				MovieIDs: &models.UpdateMovieIDs{
					Movies: append([]models.MoviesScenes{
						{
							MovieID:    movieIDs[movieIdxWithScene],
							SceneIndex: &sceneIndex,
						},
					},
						movieScenes...,
					),
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Scene{
				Movies: append([]models.MoviesScenes{
					{
						MovieID: indexesToIDs(movieIDs, sceneMovies[sceneIdxWithMovie])[0],
					},
				}, movieScenes...),
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
				StashIDs: []models.StashID{sceneStashID(sceneIdxWithSpacedName)},
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
			"add invalid movies",
			sceneIDs[sceneIdxWithMovie],
			models.ScenePartial{
				MovieIDs: &models.UpdateMovieIDs{
					Movies: []models.MoviesScenes{
						{
							MovieID: invalidID,
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
			models.Scene{},
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
				TagIDs: []int{tagIDs[tagIdx2WithScene]},
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
				PerformerIDs: []int{performerIDs[performerIdx2WithScene]},
			},
			false,
		},
		{
			"remove movies",
			sceneIDs[sceneIdxWithMovie],
			models.ScenePartial{
				MovieIDs: &models.UpdateMovieIDs{
					Movies: []models.MoviesScenes{
						{
							MovieID: movieIDs[movieIdxWithScene],
						},
					},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Scene{},
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
			models.Scene{},
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
				GalleryIDs: []int{galleryIDs[galleryIdxWithScene]},
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
				TagIDs: indexesToIDs(tagIDs, sceneTags[sceneIdxWithTwoTags]),
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
				PerformerIDs: indexesToIDs(performerIDs, scenePerformers[sceneIdxWithTwoPerformers]),
			},
			false,
		},
		{
			"remove unrelated movies",
			sceneIDs[sceneIdxWithMovie],
			models.ScenePartial{
				MovieIDs: &models.UpdateMovieIDs{
					Movies: []models.MoviesScenes{
						{
							MovieID: movieIDs[movieIdxWithDupName],
						},
					},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Scene{
				Movies: []models.MoviesScenes{
					{
						MovieID: indexesToIDs(movieIDs, sceneMovies[sceneIdxWithMovie])[0],
					},
				},
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
				StashIDs: []models.StashID{sceneStashID(sceneIdxWithGallery)},
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

			// only compare fields that were in the partial
			if tt.partial.PerformerIDs != nil {
				assert.Equal(tt.want.PerformerIDs, got.PerformerIDs)
				assert.Equal(tt.want.PerformerIDs, s.PerformerIDs)
			}
			if tt.partial.TagIDs != nil {
				assert.Equal(tt.want.TagIDs, got.TagIDs)
				assert.Equal(tt.want.TagIDs, s.TagIDs)
			}
			if tt.partial.GalleryIDs != nil {
				assert.Equal(tt.want.GalleryIDs, got.GalleryIDs)
				assert.Equal(tt.want.GalleryIDs, s.GalleryIDs)
			}
			if tt.partial.MovieIDs != nil {
				assert.Equal(tt.want.Movies, got.Movies)
				assert.Equal(tt.want.Movies, s.Movies)
			}
			if tt.partial.StashIDs != nil {
				assert.Equal(tt.want.StashIDs, got.StashIDs)
				assert.Equal(tt.want.StashIDs, s.StashIDs)
			}
		})
	}
}

func Test_sceneQueryBuilder_IncrementOCounter(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    int
		wantErr bool
	}{
		{
			"increment",
			sceneIDs[1],
			2,
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
			got, err := qb.IncrementOCounter(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.IncrementOCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sceneQueryBuilder.IncrementOCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sceneQueryBuilder_DecrementOCounter(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    int
		wantErr bool
	}{
		{
			"decrement",
			sceneIDs[2],
			1,
			false,
		},
		{
			"zero",
			sceneIDs[0],
			0,
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
			got, err := qb.DecrementOCounter(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.DecrementOCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sceneQueryBuilder.DecrementOCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sceneQueryBuilder_ResetOCounter(t *testing.T) {
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
			got, err := qb.ResetOCounter(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneQueryBuilder.ResetOCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sceneQueryBuilder.ResetOCounter() = %v, want %v", got, tt.want)
			}
		})
	}
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
			withRollbackTxn(func(ctx context.Context) error {
				if err := qb.Destroy(ctx, tt.id); (err != nil) != tt.wantErr {
					t.Errorf("sceneQueryBuilder.Destroy() error = %v, wantErr %v", err, tt.wantErr)
				}

				// ensure cannot be found
				i, err := qb.Find(ctx, tt.id)

				assert.NotNil(err)
				assert.Nil(i)
				return nil
			})
		})
	}
}

func makeSceneWithID(index int) *models.Scene {
	ret := makeScene(index)
	ret.ID = sceneIDs[index]

	if ret.Date != nil && ret.Date.IsZero() {
		ret.Date = nil
	}

	ret.Files = []*file.VideoFile{makeSceneFile(index)}

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
			true,
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
			"with movies",
			sceneIDs[sceneIdxWithMovie],
			makeSceneWithID(sceneIdxWithMovie),
			false,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			withTxn(func(ctx context.Context) error {
				got, err := qb.Find(ctx, tt.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("sceneQueryBuilder.Find() error = %v, wantErr %v", err, tt.wantErr)
					return nil
				}

				if got != nil {
					clearSceneFileIDs(got)
				}

				assert.Equal(tt.want, got)
				return nil
			})
		})
	}
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
				sceneIDs[sceneIdxWithMovie],
			},
			[]*models.Scene{
				makeSceneWithID(sceneIdxWithGallery),
				makeSceneWithID(sceneIdxWithTwoPerformers),
				makeSceneWithID(sceneIdxWithTwoTags),
				makeSceneWithID(sceneIdxWithMovie),
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

			for _, s := range got {
				clearSceneFileIDs(s)
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
			"with movies",
			getChecksum(sceneIdxWithMovie),
			[]*models.Scene{makeSceneWithID(sceneIdxWithMovie)},
			false,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			withTxn(func(ctx context.Context) error {
				assert := assert.New(t)
				got, err := qb.FindByChecksum(ctx, tt.checksum)
				if (err != nil) != tt.wantErr {
					t.Errorf("sceneQueryBuilder.FindByChecksum() error = %v, wantErr %v", err, tt.wantErr)
					return nil
				}

				for _, s := range got {
					clearSceneFileIDs(s)
				}

				assert.Equal(tt.want, got)

				return nil
			})
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
			"with movies",
			getOSHash(sceneIdxWithMovie),
			[]*models.Scene{makeSceneWithID(sceneIdxWithMovie)},
			false,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			withTxn(func(ctx context.Context) error {
				got, err := qb.FindByOSHash(ctx, tt.oshash)
				if (err != nil) != tt.wantErr {
					t.Errorf("sceneQueryBuilder.FindByOSHash() error = %v, wantErr %v", err, tt.wantErr)
					return nil
				}

				for _, s := range got {
					clearSceneFileIDs(s)
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("sceneQueryBuilder.FindByOSHash() = %v, want %v", got, tt.want)
				}
				return nil
			})
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
			"with movies",
			getPath(sceneIdxWithMovie),
			[]*models.Scene{makeSceneWithID(sceneIdxWithMovie)},
			false,
		},
	}

	qb := db.Scene

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			withTxn(func(ctx context.Context) error {
				assert := assert.New(t)
				got, err := qb.FindByPath(ctx, tt.path)
				if (err != nil) != tt.wantErr {
					t.Errorf("sceneQueryBuilder.FindByPath() error = %v, wantErr %v", err, tt.wantErr)
					return nil
				}

				for _, s := range got {
					clearSceneFileIDs(s)
				}

				assert.Equal(tt.want, got)

				return nil
			})
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

			for _, s := range got {
				clearSceneFileIDs(s)
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
		assert.Equal(t, scenePath, scene.Path())

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
		},
		SceneFilter: sceneFilter,
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

	assert.Len(t, scenes, 1)
	scene := scenes[0]
	assert.Equal(t, sceneIDs[expectedSceneIdx], scene.ID)

	// no Q should return all results
	filter.Q = nil
	scenes = queryScene(ctx, t, sqb, nil, &filter)

	assert.Len(t, scenes, totalScenes)
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
			"equals folder name",
			models.StringCriterionInput{
				Value:    folder,
				Modifier: models.CriterionModifierEquals,
			},
			[]int{sceneIdx},
			nil,
		},
		{
			"equals folder name trailing slash",
			models.StringCriterionInput{
				Value:    folder + string(filepath.Separator),
				Modifier: models.CriterionModifierEquals,
			},
			[]int{sceneIdx},
			nil,
		},
		{
			"equals base name",
			models.StringCriterionInput{
				Value:    basename,
				Modifier: models.CriterionModifierEquals,
			},
			[]int{sceneIdx},
			nil,
		},
		{
			"equals base name leading slash",
			models.StringCriterionInput{
				Value:    string(filepath.Separator) + basename,
				Modifier: models.CriterionModifierEquals,
			},
			[]int{sceneIdx},
			nil,
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
			"not equals folder name",
			models.StringCriterionInput{
				Value:    folder,
				Modifier: models.CriterionModifierNotEquals,
			},
			nil,
			[]int{sceneIdx},
		},
		{
			"not equals basename",
			models.StringCriterionInput{
				Value:    basename,
				Modifier: models.CriterionModifierNotEquals,
			},
			nil,
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

			missing := intslice.IntExclude(mustInclude, got.IDs)
			if len(missing) > 0 {
				t.Errorf("SceneStore.TestSceneQueryPath() missing expected IDs: %v", missing)
			}

			notExcluded := intslice.IntIntercect(mustExclude, got.IDs)
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
		verifyString(t, s.URL, urlCriterion)
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
		Or: &models.SceneFilterType{
			Path: &models.StringCriterionInput{
				Value:    scene2Path,
				Modifier: models.CriterionModifierEquals,
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Scene

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		if !assert.Len(t, scenes, 2) {
			return nil
		}
		assert.Equal(t, scene1Path, scenes[0].Path())
		assert.Equal(t, scene2Path, scenes[1].Path())

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
		And: &models.SceneFilterType{
			Rating: &models.IntCriterionInput{
				Value:    sceneRating,
				Modifier: models.CriterionModifierEquals,
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Scene

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		if !assert.Len(t, scenes, 1) {
			return nil
		}
		assert.Equal(t, scenePath, scenes[0].Path())
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
		Not: &models.SceneFilterType{
			Rating: &ratingCriterion,
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Scene

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		for _, scene := range scenes {
			verifyString(t, scene.Path(), pathCriterion)
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
		And: &subFilter,
		Or:  &subFilter,
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
			verifyString(t, scene.Path(), pathCriterion)
		}

		return nil
	})
}

func verifyNullString(t *testing.T, value sql.NullString, criterion models.StringCriterionInput) {
	t.Helper()
	assert := assert.New(t)
	if criterion.Modifier == models.CriterionModifierIsNull {
		if value.Valid && value.String == "" {
			// correct
			return
		}
		assert.False(value.Valid, "expect is null values to be null")
	}
	if criterion.Modifier == models.CriterionModifierNotNull {
		assert.True(value.Valid, "expect is null values to be null")
		assert.Greater(len(value.String), 0)
	}
	if criterion.Modifier == models.CriterionModifierEquals {
		assert.Equal(criterion.Value, value.String)
	}
	if criterion.Modifier == models.CriterionModifierNotEquals {
		assert.NotEqual(criterion.Value, value.String)
	}
	if criterion.Modifier == models.CriterionModifierMatchesRegex {
		assert.True(value.Valid)
		assert.Regexp(regexp.MustCompile(criterion.Value), value)
	}
	if criterion.Modifier == models.CriterionModifierNotMatchesRegex {
		if !value.Valid {
			// correct
			return
		}
		assert.NotRegexp(regexp.MustCompile(criterion.Value), value)
	}
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
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		sceneFilter := models.SceneFilterType{
			Rating: &ratingCriterion,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		for _, scene := range scenes {
			verifyIntPtr(t, scene.Rating, ratingCriterion)
		}

		return nil
	})
}

func verifyInt64(t *testing.T, value sql.NullInt64, criterion models.IntCriterionInput) {
	t.Helper()
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
			verifyInt(t, scene.OCounter, oCounterCriterion)
		}

		return nil
	})
}

func verifyInt(t *testing.T, value int, criterion models.IntCriterionInput) {
	t.Helper()
	assert := assert.New(t)
	if criterion.Modifier == models.CriterionModifierEquals {
		assert.Equal(criterion.Value, value)
	}
	if criterion.Modifier == models.CriterionModifierNotEquals {
		assert.NotEqual(criterion.Value, value)
	}
	if criterion.Modifier == models.CriterionModifierGreaterThan {
		assert.Greater(value, criterion.Value)
	}
	if criterion.Modifier == models.CriterionModifierLessThan {
		assert.Less(value, criterion.Value)
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
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		sceneFilter := models.SceneFilterType{
			Duration: &durationCriterion,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		for _, scene := range scenes {
			duration := scene.Duration()
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
			f := scene.PrimaryFile()
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

	return queryScene(ctx, t, queryBuilder, &sceneFilter, nil)
}

func createScene(ctx context.Context, width int, height int) (*models.Scene, error) {
	name := fmt.Sprintf("TestSceneQueryResolutionModifiers %d %d", width, height)

	sceneFile := &file.VideoFile{
		BaseFile: &file.BaseFile{
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

	if err := db.Scene.Create(ctx, scene, []file.ID{sceneFile.ID}); err != nil {
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

		q := getSceneStringValue(sceneIdxWithMovie, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, &findFilter)

		assert.Len(t, scenes, 0)

		findFilter.Q = nil
		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)

		// ensure non of the ids equal the one with movies
		for _, scene := range scenes {
			assert.NotEqual(t, sceneIDs[sceneIdxWithMovie], scene.ID)
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

		// three in four scenes have no date
		assert.Len(t, scenes, int(math.Ceil(float64(totalScenes)/4*3)))

		// ensure date is null, empty or "0001-01-01"
		for _, scene := range scenes {
			assert.True(t, scene.Date == nil || scene.Date.Time == time.Time{})
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

		// ensure date is null, empty or "0001-01-01"
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
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
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

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

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

		scenes = queryScene(ctx, t, sqb, &sceneFilter, nil)

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

		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)
		assert.Len(t, scenes, 0)

		return nil
	})
}

func TestSceneQueryTags(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		tagCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdxWithScene]),
				strconv.Itoa(tagIDs[tagIdx1WithScene]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		sceneFilter := models.SceneFilterType{
			Tags: &tagCriterion,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)
		assert.Len(t, scenes, 2)

		// ensure ids are correct
		for _, scene := range scenes {
			assert.True(t, scene.ID == sceneIDs[sceneIdxWithTag] || scene.ID == sceneIDs[sceneIdxWithTwoTags])
		}

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithScene]),
				strconv.Itoa(tagIDs[tagIdx2WithScene]),
			},
			Modifier: models.CriterionModifierIncludesAll,
		}

		scenes = queryScene(ctx, t, sqb, &sceneFilter, nil)

		assert.Len(t, scenes, 1)
		assert.Equal(t, sceneIDs[sceneIdxWithTwoTags], scenes[0].ID)

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithScene]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getSceneStringValue(sceneIdxWithTwoTags, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)
		assert.Len(t, scenes, 0)

		return nil
	})
}

func TestSceneQueryPerformerTags(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		tagCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdxWithPerformer]),
				strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		sceneFilter := models.SceneFilterType{
			PerformerTags: &tagCriterion,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)
		assert.Len(t, scenes, 2)

		// ensure ids are correct
		for _, scene := range scenes {
			assert.True(t, scene.ID == sceneIDs[sceneIdxWithPerformerTag] || scene.ID == sceneIDs[sceneIdxWithPerformerTwoTags])
		}

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
				strconv.Itoa(tagIDs[tagIdx2WithPerformer]),
			},
			Modifier: models.CriterionModifierIncludesAll,
		}

		scenes = queryScene(ctx, t, sqb, &sceneFilter, nil)

		assert.Len(t, scenes, 1)
		assert.Equal(t, sceneIDs[sceneIdxWithPerformerTwoTags], scenes[0].ID)

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getSceneStringValue(sceneIdxWithPerformerTwoTags, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)
		assert.Len(t, scenes, 0)

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Modifier: models.CriterionModifierIsNull,
		}
		q = getSceneStringValue(sceneIdx1WithPerformer, titleField)

		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)
		assert.Len(t, scenes, 1)
		assert.Equal(t, sceneIDs[sceneIdx1WithPerformer], scenes[0].ID)

		q = getSceneStringValue(sceneIdxWithPerformerTag, titleField)
		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)
		assert.Len(t, scenes, 0)

		tagCriterion.Modifier = models.CriterionModifierNotNull

		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)
		assert.Len(t, scenes, 1)
		assert.Equal(t, sceneIDs[sceneIdxWithPerformerTag], scenes[0].ID)

		q = getSceneStringValue(sceneIdx1WithPerformer, titleField)
		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)
		assert.Len(t, scenes, 0)

		return nil
	})
}

func TestSceneQueryStudio(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene
		studioCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithScene]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		sceneFilter := models.SceneFilterType{
			Studios: &studioCriterion,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

		assert.Len(t, scenes, 1)

		// ensure id is correct
		assert.Equal(t, sceneIDs[sceneIdxWithStudio], scenes[0].ID)

		studioCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithScene]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getSceneStringValue(sceneIdxWithStudio, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		scenes = queryScene(ctx, t, sqb, &sceneFilter, &findFilter)
		assert.Len(t, scenes, 0)

		return nil
	})
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
				strconv.Itoa(movieIDs[movieIdxWithScene]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		sceneFilter := models.SceneFilterType{
			Movies: &movieCriterion,
		}

		scenes := queryScene(ctx, t, sqb, &sceneFilter, nil)

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
			"size",
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
			verifyInt(t, len(scene.TagIDs), tagCountCriterion)
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
			verifyInt(t, len(scene.PerformerIDs), performerCountCriterion)
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

func TestSceneCountByMovieID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Scene

		sceneCount, err := sqb.CountByMovieID(ctx, movieIDs[movieIdxWithScene])

		if err != nil {
			t.Errorf("error calling CountByMovieID: %s", err.Error())
		}

		assert.Equal(t, 1, sceneCount)

		sceneCount, err = sqb.CountByMovieID(ctx, 0)

		if err != nil {
			t.Errorf("error calling CountByMovieID: %s", err.Error())
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

		scenes, err := sqb.FindByMovieID(ctx, movieIDs[movieIdxWithScene])

		if err != nil {
			t.Errorf("error calling FindByMovieID: %s", err.Error())
		}

		assert.Len(t, scenes, 1)
		assert.Equal(t, sceneIDs[sceneIdxWithMovie], scenes[0].ID)

		scenes, err = sqb.FindByMovieID(ctx, 0)

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

		image := []byte("image")
		if err := qb.UpdateCover(ctx, sceneID, image); err != nil {
			return fmt.Errorf("Error updating scene cover: %s", err.Error())
		}

		// ensure image set
		storedImage, err := qb.GetCover(ctx, sceneID)
		if err != nil {
			return fmt.Errorf("Error getting image: %s", err.Error())
		}
		assert.Equal(t, storedImage, image)

		// set nil image
		err = qb.UpdateCover(ctx, sceneID, nil)
		if err == nil {
			return fmt.Errorf("Expected error setting nil image")
		}

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestSceneDestroySceneCover(t *testing.T) {
	if err := withTxn(func(ctx context.Context) error {
		qb := db.Scene

		sceneID := sceneIDs[sceneIdxWithGallery]

		image := []byte("image")
		if err := qb.UpdateCover(ctx, sceneID, image); err != nil {
			return fmt.Errorf("Error updating scene image: %s", err.Error())
		}

		if err := qb.DestroyCover(ctx, sceneID); err != nil {
			return fmt.Errorf("Error destroying scene cover: %s", err.Error())
		}

		// image should be nil
		storedImage, err := qb.GetCover(ctx, sceneID)
		if err != nil {
			return fmt.Errorf("Error getting image: %s", err.Error())
		}
		assert.Nil(t, storedImage)

		return nil
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

		testSceneStashIDs(ctx, t, scene)
		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func testSceneStashIDs(ctx context.Context, t *testing.T, s *models.Scene) {
	// ensure no stash IDs to begin with
	assert.Len(t, s.StashIDs, 0)

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

	assert.Equal(t, []models.StashID{stashID}, s.StashIDs)

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

	assert.Len(t, s.StashIDs, 0)
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
		got, err := qb.FindDuplicates(ctx, distance)
		if err != nil {
			t.Errorf("SceneStore.FindDuplicates() error = %v", err)
			return nil
		}

		assert.Len(t, got, dupeScenePhashes)

		distance = 1
		got, err = qb.FindDuplicates(ctx, distance)
		if err != nil {
			t.Errorf("SceneStore.FindDuplicates() error = %v", err)
			return nil
		}

		assert.Len(t, got, dupeScenePhashes)

		return nil
	})
}

// TODO Count
// TODO SizeCount
