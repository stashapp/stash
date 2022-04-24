//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"math"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stretchr/testify/assert"
)

var invalidID = -1

func Test_galleryQueryBuilder_Create(t *testing.T) {
	var (
		path        = "path"
		title       = "title"
		checksum    = "checksum"
		url         = "url"
		rating      = 3
		details     = "details"
		fileModTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		createdAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	date := models.NewDate("2003-02-01")

	tests := []struct {
		name      string
		newObject models.Gallery
		wantErr   bool
	}{
		{
			"full",
			models.Gallery{
				Path:         &path,
				Checksum:     checksum,
				Zip:          false,
				Title:        title,
				URL:          url,
				Date:         &date,
				Details:      details,
				Rating:       &rating,
				Organized:    true,
				StudioID:     &studioIDs[studioIdxWithScene],
				FileModTime:  &fileModTime,
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				SceneIDs:     []int{sceneIDs[sceneIdx1WithPerformer], sceneIDs[sceneIdx1WithStudio]},
				TagIDs:       []int{tagIDs[tagIdx1WithScene], tagIDs[tagIdx1WithDupName]},
				PerformerIDs: []int{performerIDs[performerIdx1WithScene], performerIDs[performerIdx1WithDupName]},
			},
			false,
		},
		{
			"invalid studio id",
			models.Gallery{
				StudioID: &invalidID,
			},
			true,
		},
		{
			"invalid scene id",
			models.Gallery{
				SceneIDs: []int{invalidID},
			},
			true,
		},
		{
			"invalid tag id",
			models.Gallery{
				TagIDs: []int{invalidID},
			},
			true,
		},
		{
			"invalid performer id",
			models.Gallery{
				PerformerIDs: []int{invalidID},
			},
			true,
		},
	}

	qb := sqlite.GalleryReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			s := tt.newObject
			if err := qb.Create(ctx, &s); (err != nil) != tt.wantErr {
				t.Errorf("galleryQueryBuilder.Create() error = %v, wantErr = %v", err, tt.wantErr)
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
				t.Errorf("galleryQueryBuilder.Find() error = %v", err)
			}

			assert.Equal(copy, *found)

			return
		})
	}
}

func Test_galleryQueryBuilder_Update(t *testing.T) {
	var (
		path        = "path"
		title       = "title"
		checksum    = "checksum"
		url         = "url"
		rating      = 3
		details     = "details"
		fileModTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		createdAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	date := models.NewDate("2003-02-01")

	tests := []struct {
		name          string
		updatedObject *models.Gallery
		wantErr       bool
	}{
		{
			"full",
			&models.Gallery{
				ID:           galleryIDs[galleryIdxWithScene],
				Path:         &path,
				Checksum:     checksum,
				Zip:          false,
				Title:        title,
				URL:          url,
				Date:         &date,
				Details:      details,
				Rating:       &rating,
				Organized:    true,
				StudioID:     &studioIDs[studioIdxWithScene],
				FileModTime:  &fileModTime,
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				SceneIDs:     []int{sceneIDs[sceneIdx1WithPerformer], sceneIDs[sceneIdx1WithStudio]},
				TagIDs:       []int{tagIDs[tagIdx1WithScene], tagIDs[tagIdx1WithDupName]},
				PerformerIDs: []int{performerIDs[performerIdx1WithScene], performerIDs[performerIdx1WithDupName]},
			},
			false,
		},
		{
			"clear nullables",
			&models.Gallery{
				ID:        galleryIDs[galleryIdxWithImage],
				Checksum:  checksum,
				Organized: true,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"clear scene ids",
			&models.Gallery{
				ID:        galleryIDs[galleryIdxWithScene],
				Checksum:  checksum,
				Organized: true,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"clear tag ids",
			&models.Gallery{
				ID:        galleryIDs[galleryIdxWithTag],
				Checksum:  checksum,
				Organized: true,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"clear performer ids",
			&models.Gallery{
				ID:        galleryIDs[galleryIdxWithPerformer],
				Checksum:  checksum,
				Organized: true,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"invalid studio id",
			&models.Gallery{
				ID:        galleryIDs[galleryIdxWithImage],
				Checksum:  checksum,
				Organized: true,
				StudioID:  &invalidID,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			true,
		},
		{
			"invalid scene id",
			&models.Gallery{
				ID:        galleryIDs[galleryIdxWithImage],
				Checksum:  checksum,
				Organized: true,
				SceneIDs:  []int{invalidID},
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			true,
		},
		{
			"invalid tag id",
			&models.Gallery{
				ID:        galleryIDs[galleryIdxWithImage],
				Checksum:  checksum,
				Organized: true,
				TagIDs:    []int{invalidID},
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			true,
		},
		{
			"invalid performer id",
			&models.Gallery{
				ID:           galleryIDs[galleryIdxWithImage],
				Checksum:     checksum,
				Organized:    true,
				PerformerIDs: []int{invalidID},
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
			},
			true,
		},
	}

	qb := sqlite.GalleryReaderWriter
	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			copy := *tt.updatedObject

			if err := qb.Update(ctx, tt.updatedObject); (err != nil) != tt.wantErr {
				t.Errorf("galleryQueryBuilder.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			s, err := qb.Find(ctx, tt.updatedObject.ID)
			if err != nil {
				t.Errorf("galleryQueryBuilder.Find() error = %v", err)
			}

			assert.Equal(copy, *s)

			return
		})
	}
}

func clearGalleryPartial() models.GalleryPartial {
	// leave mandatory fields
	return models.GalleryPartial{
		Title:        models.OptionalString{Set: true, Null: true},
		Details:      models.OptionalString{Set: true, Null: true},
		URL:          models.OptionalString{Set: true, Null: true},
		Date:         models.OptionalDate{Set: true, Null: true},
		Rating:       models.OptionalInt{Set: true, Null: true},
		StudioID:     models.OptionalInt{Set: true, Null: true},
		FileModTime:  models.OptionalTime{Set: true, Null: true},
		TagIDs:       &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
		PerformerIDs: &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
	}
}

func Test_galleryQueryBuilder_UpdatePartial(t *testing.T) {
	var (
		path        = "path"
		title       = "title"
		checksum    = "checksum"
		zip         = true
		details     = "details"
		url         = "url"
		rating      = 3
		fileModTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		createdAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)

		date = models.NewDate("2003-02-01")
	)

	tests := []struct {
		name    string
		id      int
		partial models.GalleryPartial
		want    models.Gallery
		wantErr bool
	}{
		{
			"full",
			galleryIDs[galleryIdxWithImage],
			models.GalleryPartial{
				Path:        models.NewOptionalString(path),
				Checksum:    models.NewOptionalString(checksum),
				Zip:         models.NewOptionalBool(zip),
				Title:       models.NewOptionalString(title),
				Details:     models.NewOptionalString(details),
				URL:         models.NewOptionalString(url),
				Date:        models.NewOptionalDate(date),
				Rating:      models.NewOptionalInt(rating),
				Organized:   models.NewOptionalBool(true),
				StudioID:    models.NewOptionalInt(studioIDs[studioIdxWithGallery]),
				FileModTime: models.NewOptionalTime(fileModTime),
				CreatedAt:   models.NewOptionalTime(createdAt),
				UpdatedAt:   models.NewOptionalTime(updatedAt),

				SceneIDs: &models.UpdateIDs{
					IDs:  []int{sceneIDs[sceneIdxWithGallery]},
					Mode: models.RelationshipUpdateModeSet,
				},
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithGallery], tagIDs[tagIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeSet,
				},
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithGallery], performerIDs[performerIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeSet,
				},
			},
			models.Gallery{
				ID:           galleryIDs[galleryIdxWithImage],
				Path:         &path,
				Checksum:     checksum,
				Zip:          zip,
				Title:        title,
				Details:      details,
				URL:          url,
				Date:         &date,
				Rating:       &rating,
				Organized:    true,
				StudioID:     &studioIDs[studioIdxWithGallery],
				FileModTime:  &fileModTime,
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				SceneIDs:     []int{sceneIDs[sceneIdxWithGallery]},
				TagIDs:       []int{tagIDs[tagIdx1WithGallery], tagIDs[tagIdx1WithDupName]},
				PerformerIDs: []int{performerIDs[performerIdx1WithGallery], performerIDs[performerIdx1WithDupName]},
			},
			false,
		},
		{
			"clear all",
			galleryIDs[galleryIdxWithImage],
			clearGalleryPartial(),
			models.Gallery{
				ID:       galleryIDs[galleryIdxWithImage],
				Path:     getGalleryNullStringPtr(galleryIdxWithImage, pathField),
				Checksum: getGalleryStringValue(galleryIdxWithImage, checksumField),
			},
			false,
		},
		{
			"invalid id",
			invalidID,
			models.GalleryPartial{},
			models.Gallery{},
			true,
		},
	}
	for _, tt := range tests {
		qb := sqlite.GalleryReaderWriter

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			got, err := qb.UpdatePartial(ctx, tt.id, tt.partial)
			if (err != nil) != tt.wantErr {
				t.Errorf("galleryQueryBuilder.UpdatePartial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			assert.Equal(tt.want, *got)

			s, err := qb.Find(ctx, tt.id)
			if err != nil {
				t.Errorf("galleryQueryBuilder.Find() error = %v", err)
			}

			assert.Equal(tt.want, *s)
		})
	}
}

func Test_galleryQueryBuilder_UpdatePartialRelationships(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		partial models.GalleryPartial
		want    models.Gallery
		wantErr bool
	}{
		{
			"add scenes",
			galleryIDs[galleryIdx1WithImage],
			models.GalleryPartial{
				SceneIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[sceneIdx1WithStudio], tagIDs[sceneIdx1WithPerformer]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Gallery{
				SceneIDs: append(indexesToIDs(sceneIDs, sceneGalleries.reverseLookup(galleryIdx1WithImage)),
					sceneIDs[sceneIdx1WithStudio],
					sceneIDs[sceneIdx1WithPerformer],
				),
			},
			false,
		},
		{
			"add tags",
			galleryIDs[galleryIdxWithTwoTags],
			models.GalleryPartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithImage]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Gallery{
				TagIDs: append(indexesToIDs(tagIDs, galleryTags[galleryIdxWithTwoTags]),
					tagIDs[tagIdx1WithDupName],
					tagIDs[tagIdx1WithImage],
				),
			},
			false,
		},
		{
			"add performers",
			galleryIDs[galleryIdxWithTwoPerformers],
			models.GalleryPartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithDupName], performerIDs[performerIdx1WithImage]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Gallery{
				PerformerIDs: append(indexesToIDs(performerIDs, galleryPerformers[galleryIdxWithTwoPerformers]),
					performerIDs[performerIdx1WithDupName],
					performerIDs[performerIdx1WithImage],
				),
			},
			false,
		},
		{
			"add duplicate scenes",
			galleryIDs[galleryIdxWithScene],
			models.GalleryPartial{
				SceneIDs: &models.UpdateIDs{
					IDs:  []int{sceneIDs[sceneIdxWithGallery], sceneIDs[sceneIdx1WithPerformer]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Gallery{
				SceneIDs: append(indexesToIDs(sceneIDs, sceneGalleries.reverseLookup(galleryIdxWithScene)),
					sceneIDs[sceneIdx1WithPerformer],
				),
			},
			false,
		},
		{
			"add duplicate tags",
			galleryIDs[galleryIdxWithTwoTags],
			models.GalleryPartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithGallery], tagIDs[tagIdx1WithScene]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Gallery{
				TagIDs: append(indexesToIDs(tagIDs, galleryTags[galleryIdxWithTwoTags]),
					tagIDs[tagIdx1WithScene],
				),
			},
			false,
		},
		{
			"add duplicate performers",
			galleryIDs[galleryIdxWithTwoPerformers],
			models.GalleryPartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithGallery], performerIDs[performerIdx1WithScene]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Gallery{
				PerformerIDs: append(indexesToIDs(performerIDs, galleryPerformers[galleryIdxWithTwoPerformers]),
					performerIDs[performerIdx1WithScene],
				),
			},
			false,
		},
		{
			"add invalid scenes",
			galleryIDs[galleryIdxWithScene],
			models.GalleryPartial{
				SceneIDs: &models.UpdateIDs{
					IDs:  []int{invalidID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Gallery{},
			true,
		},
		{
			"add invalid tags",
			galleryIDs[galleryIdxWithTwoTags],
			models.GalleryPartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{invalidID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Gallery{},
			true,
		},
		{
			"add invalid performers",
			galleryIDs[galleryIdxWithTwoPerformers],
			models.GalleryPartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{invalidID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Gallery{},
			true,
		},
		{
			"remove scenes",
			galleryIDs[galleryIdxWithScene],
			models.GalleryPartial{
				SceneIDs: &models.UpdateIDs{
					IDs:  []int{sceneIDs[sceneIdxWithGallery]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Gallery{},
			false,
		},
		{
			"remove tags",
			galleryIDs[galleryIdxWithTwoTags],
			models.GalleryPartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithGallery]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Gallery{
				TagIDs: []int{tagIDs[tagIdx2WithGallery]},
			},
			false,
		},
		{
			"remove performers",
			galleryIDs[galleryIdxWithTwoPerformers],
			models.GalleryPartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithGallery]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Gallery{
				PerformerIDs: []int{performerIDs[performerIdx2WithGallery]},
			},
			false,
		},
		{
			"remove unrelated scenes",
			galleryIDs[galleryIdxWithScene],
			models.GalleryPartial{
				SceneIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[sceneIdx1WithPerformer]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Gallery{
				SceneIDs: []int{sceneIDs[sceneIdxWithGallery]},
			},
			false,
		},
		{
			"remove unrelated tags",
			galleryIDs[galleryIdxWithTwoTags],
			models.GalleryPartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithPerformer]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Gallery{
				TagIDs: indexesToIDs(tagIDs, galleryTags[galleryIdxWithTwoTags]),
			},
			false,
		},
		{
			"remove unrelated performers",
			galleryIDs[galleryIdxWithTwoPerformers],
			models.GalleryPartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Gallery{
				PerformerIDs: indexesToIDs(performerIDs, galleryPerformers[galleryIdxWithTwoPerformers]),
			},
			false,
		},
	}

	for _, tt := range tests {
		qb := sqlite.GalleryReaderWriter

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			got, err := qb.UpdatePartial(ctx, tt.id, tt.partial)
			if (err != nil) != tt.wantErr {
				t.Errorf("galleryQueryBuilder.UpdatePartial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			s, err := qb.Find(ctx, tt.id)
			if err != nil {
				t.Errorf("galleryQueryBuilder.Find() error = %v", err)
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
			if tt.partial.SceneIDs != nil {
				assert.Equal(tt.want.SceneIDs, got.SceneIDs)
				assert.Equal(tt.want.SceneIDs, s.SceneIDs)
			}
		})
	}
}

func Test_galleryQueryBuilder_Destroy(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			"valid",
			galleryIDs[galleryIdxWithScene],
			false,
		},
		{
			"invalid",
			invalidID,
			true,
		},
	}

	qb := sqlite.GalleryReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			if err := qb.Destroy(ctx, tt.id); (err != nil) != tt.wantErr {
				t.Errorf("galleryQueryBuilder.Destroy() error = %v, wantErr %v", err, tt.wantErr)
			}

			// ensure cannot be found
			i, err := qb.Find(ctx, tt.id)

			assert.NotNil(err)
			assert.Nil(i)
			return

		})
	}
}

func makeGalleryWithID(index int) *models.Gallery {
	const includeScenes = true
	ret := makeGallery(index, includeScenes)
	ret.ID = galleryIDs[index]

	if ret.Date != nil && ret.Date.IsZero() {
		ret.Date = nil
	}

	return ret
}

func Test_galleryQueryBuilder_Find(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    *models.Gallery
		wantErr bool
	}{
		{
			"valid",
			galleryIDs[galleryIdxWithImage],
			makeGalleryWithID(galleryIdxWithImage),
			false,
		},
		{
			"invalid",
			invalidID,
			nil,
			true,
		},
		{
			"with performers",
			galleryIDs[galleryIdxWithTwoPerformers],
			makeGalleryWithID(galleryIdxWithTwoPerformers),
			false,
		},
		{
			"with tags",
			galleryIDs[galleryIdxWithTwoTags],
			makeGalleryWithID(galleryIdxWithTwoTags),
			false,
		},
	}

	qb := sqlite.GalleryReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.Find(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("galleryQueryBuilder.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func Test_galleryQueryBuilder_FindMany(t *testing.T) {
	tests := []struct {
		name    string
		ids     []int
		want    []*models.Gallery
		wantErr bool
	}{
		{
			"valid with relationships",
			[]int{galleryIDs[galleryIdxWithImage], galleryIDs[galleryIdxWithTwoPerformers], galleryIDs[galleryIdxWithTwoTags]},
			[]*models.Gallery{
				makeGalleryWithID(galleryIdxWithImage),
				makeGalleryWithID(galleryIdxWithTwoPerformers),
				makeGalleryWithID(galleryIdxWithTwoTags),
			},
			false,
		},
		{
			"invalid",
			[]int{galleryIDs[galleryIdxWithImage], galleryIDs[galleryIdxWithTwoPerformers], invalidID},
			nil,
			true,
		},
	}

	qb := sqlite.GalleryReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.FindMany(ctx, tt.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("galleryQueryBuilder.FindMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("galleryQueryBuilder.FindMany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_galleryQueryBuilder_FindByChecksum(t *testing.T) {
	getChecksum := func(index int) string {
		return getGalleryStringValue(index, checksumField)
	}

	tests := []struct {
		name     string
		checksum string
		want     *models.Gallery
		wantErr  bool
	}{
		{
			"valid",
			getChecksum(galleryIdxWithImage),
			makeGalleryWithID(galleryIdxWithImage),
			false,
		},
		{
			"invalid",
			"invalid checksum",
			nil,
			false,
		},
		{
			"with performers",
			getChecksum(galleryIdxWithTwoPerformers),
			makeGalleryWithID(galleryIdxWithTwoPerformers),
			false,
		},
		{
			"with tags",
			getChecksum(galleryIdxWithTwoTags),
			makeGalleryWithID(galleryIdxWithTwoTags),
			false,
		},
	}

	qb := sqlite.GalleryReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.FindByChecksum(ctx, tt.checksum)
			if (err != nil) != tt.wantErr {
				t.Errorf("galleryQueryBuilder.FindByChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("galleryQueryBuilder.FindByChecksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_galleryQueryBuilder_FindByChecksums(t *testing.T) {
	getChecksum := func(index int) string {
		return getGalleryStringValue(index, checksumField)
	}

	tests := []struct {
		name      string
		checksums []string
		want      []*models.Gallery
		wantErr   bool
	}{
		{
			"valid with relationships",
			[]string{
				getChecksum(galleryIdxWithImage),
				getChecksum(galleryIdxWithTwoPerformers),
				getChecksum(galleryIdxWithTwoTags),
			},
			[]*models.Gallery{
				makeGalleryWithID(galleryIdxWithImage),
				makeGalleryWithID(galleryIdxWithTwoPerformers),
				makeGalleryWithID(galleryIdxWithTwoTags),
			},
			false,
		},
		{
			"with invalid",
			[]string{
				getChecksum(galleryIdxWithImage),
				getChecksum(galleryIdxWithTwoPerformers),
				"invalid checksum",
				getChecksum(galleryIdxWithTwoTags),
			},
			[]*models.Gallery{
				makeGalleryWithID(galleryIdxWithImage),
				makeGalleryWithID(galleryIdxWithTwoPerformers),
				makeGalleryWithID(galleryIdxWithTwoTags),
			},
			false,
		},
	}

	qb := sqlite.GalleryReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.FindByChecksums(ctx, tt.checksums)
			if (err != nil) != tt.wantErr {
				t.Errorf("galleryQueryBuilder.FindByChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("galleryQueryBuilder.FindByChecksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_galleryQueryBuilder_FindByPath(t *testing.T) {
	getPath := func(index int) string {
		return getGalleryStringValue(index, pathField)
	}

	tests := []struct {
		name    string
		path    string
		want    *models.Gallery
		wantErr bool
	}{
		{
			"valid",
			getPath(galleryIdxWithImage),
			makeGalleryWithID(galleryIdxWithImage),
			false,
		},
		{
			"invalid",
			"invalid path",
			nil,
			false,
		},
		{
			"with performers",
			getPath(galleryIdxWithTwoPerformers),
			makeGalleryWithID(galleryIdxWithTwoPerformers),
			false,
		},
		{
			"with tags",
			getPath(galleryIdxWithTwoTags),
			makeGalleryWithID(galleryIdxWithTwoTags),
			false,
		},
	}

	qb := sqlite.GalleryReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.FindByPath(ctx, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("galleryQueryBuilder.FindByPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("galleryQueryBuilder.FindByPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_galleryQueryBuilder_FindBySceneID(t *testing.T) {
	tests := []struct {
		name    string
		sceneID int
		want    []*models.Gallery
		wantErr bool
	}{
		{
			"valid",
			sceneIDs[sceneIdxWithGallery],
			[]*models.Gallery{makeGalleryWithID(galleryIdxWithScene)},
			false,
		},
		{
			"none",
			sceneIDs[sceneIdx1WithPerformer],
			nil,
			false,
		},
	}

	qb := sqlite.GalleryReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindBySceneID(ctx, tt.sceneID)
			if (err != nil) != tt.wantErr {
				t.Errorf("galleryQueryBuilder.FindBySceneID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func Test_galleryQueryBuilder_FindByImageID(t *testing.T) {
	tests := []struct {
		name    string
		imageID int
		want    []*models.Gallery
		wantErr bool
	}{
		{
			"valid",
			imageIDs[imageIdxWithTwoGalleries],
			[]*models.Gallery{
				makeGalleryWithID(galleryIdx1WithImage),
				makeGalleryWithID(galleryIdx2WithImage),
			},
			false,
		},
		{
			"none",
			imageIDs[imageIdx1WithPerformer],
			nil,
			false,
		},
	}

	qb := sqlite.GalleryReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.FindByImageID(ctx, tt.imageID)
			if (err != nil) != tt.wantErr {
				t.Errorf("galleryQueryBuilder.FindByImageID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("galleryQueryBuilder.FindByImageID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_galleryQueryBuilder_CountByImageID(t *testing.T) {
	tests := []struct {
		name    string
		imageID int
		want    int
		wantErr bool
	}{
		{
			"valid",
			imageIDs[imageIdxWithTwoGalleries],
			2,
			false,
		},
		{
			"none",
			imageIDs[imageIdx1WithPerformer],
			0,
			false,
		},
	}

	qb := sqlite.GalleryReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.CountByImageID(ctx, tt.imageID)
			if (err != nil) != tt.wantErr {
				t.Errorf("galleryQueryBuilder.CountByImageID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("galleryQueryBuilder.CountByImageID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGalleryQueryQ(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		const galleryIdx = 0

		q := getGalleryStringValue(galleryIdx, pathField)

		sqb := sqlite.GalleryReaderWriter

		galleryQueryQ(ctx, t, sqb, q, galleryIdx)

		return nil
	})
}

func galleryQueryQ(ctx context.Context, t *testing.T, qb models.GalleryReader, q string, expectedGalleryIdx int) {
	filter := models.FindFilterType{
		Q: &q,
	}
	galleries, _, err := qb.Query(ctx, nil, &filter)
	if err != nil {
		t.Errorf("Error querying gallery: %s", err.Error())
	}

	assert.Len(t, galleries, 1)
	gallery := galleries[0]
	assert.Equal(t, galleryIDs[expectedGalleryIdx], gallery.ID)

	// no Q should return all results
	filter.Q = nil
	galleries, _, err = qb.Query(ctx, nil, &filter)
	if err != nil {
		t.Errorf("Error querying gallery: %s", err.Error())
	}

	assert.Len(t, galleries, totalGalleries)
}

func TestGalleryQueryPath(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		const galleryIdx = 1
		galleryPath := getGalleryStringValue(galleryIdx, "Path")

		pathCriterion := models.StringCriterionInput{
			Value:    galleryPath,
			Modifier: models.CriterionModifierEquals,
		}

		verifyGalleriesPath(ctx, t, sqlite.GalleryReaderWriter, pathCriterion)

		pathCriterion.Modifier = models.CriterionModifierNotEquals
		verifyGalleriesPath(ctx, t, sqlite.GalleryReaderWriter, pathCriterion)

		pathCriterion.Modifier = models.CriterionModifierMatchesRegex
		pathCriterion.Value = "gallery.*1_Path"
		verifyGalleriesPath(ctx, t, sqlite.GalleryReaderWriter, pathCriterion)

		pathCriterion.Modifier = models.CriterionModifierNotMatchesRegex
		verifyGalleriesPath(ctx, t, sqlite.GalleryReaderWriter, pathCriterion)

		return nil
	})
}

func verifyGalleriesPath(ctx context.Context, t *testing.T, sqb models.GalleryReader, pathCriterion models.StringCriterionInput) {
	galleryFilter := models.GalleryFilterType{
		Path: &pathCriterion,
	}

	galleries, _, err := sqb.Query(ctx, &galleryFilter, nil)
	if err != nil {
		t.Errorf("Error querying gallery: %s", err.Error())
	}

	for _, gallery := range galleries {
		verifyStringPtr(t, gallery.Path, pathCriterion)
	}
}

func TestGalleryQueryPathOr(t *testing.T) {
	const gallery1Idx = 1
	const gallery2Idx = 2

	gallery1Path := getGalleryStringValue(gallery1Idx, "Path")
	gallery2Path := getGalleryStringValue(gallery2Idx, "Path")

	galleryFilter := models.GalleryFilterType{
		Path: &models.StringCriterionInput{
			Value:    gallery1Path,
			Modifier: models.CriterionModifierEquals,
		},
		Or: &models.GalleryFilterType{
			Path: &models.StringCriterionInput{
				Value:    gallery2Path,
				Modifier: models.CriterionModifierEquals,
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter

		galleries := queryGallery(ctx, t, sqb, &galleryFilter, nil)

		assert.Len(t, galleries, 2)
		assert.Equal(t, gallery1Path, *galleries[0].Path)
		assert.Equal(t, gallery2Path, *galleries[1].Path)

		return nil
	})
}

func TestGalleryQueryPathAndRating(t *testing.T) {
	const galleryIdx = 1
	galleryPath := getGalleryStringValue(galleryIdx, "Path")
	galleryRating := getIntPtr(getRating(galleryIdx))

	galleryFilter := models.GalleryFilterType{
		Path: &models.StringCriterionInput{
			Value:    galleryPath,
			Modifier: models.CriterionModifierEquals,
		},
		And: &models.GalleryFilterType{
			Rating: &models.IntCriterionInput{
				Value:    *galleryRating,
				Modifier: models.CriterionModifierEquals,
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter

		galleries := queryGallery(ctx, t, sqb, &galleryFilter, nil)

		assert.Len(t, galleries, 1)
		assert.Equal(t, galleryPath, *galleries[0].Path)
		assert.Equal(t, *galleryRating, *galleries[0].Rating)

		return nil
	})
}

func TestGalleryQueryPathNotRating(t *testing.T) {
	const galleryIdx = 1

	galleryRating := getRating(galleryIdx)

	pathCriterion := models.StringCriterionInput{
		Value:    "gallery_.*1_Path",
		Modifier: models.CriterionModifierMatchesRegex,
	}

	ratingCriterion := models.IntCriterionInput{
		Value:    int(galleryRating.Int64),
		Modifier: models.CriterionModifierEquals,
	}

	galleryFilter := models.GalleryFilterType{
		Path: &pathCriterion,
		Not: &models.GalleryFilterType{
			Rating: &ratingCriterion,
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter

		galleries := queryGallery(ctx, t, sqb, &galleryFilter, nil)

		for _, gallery := range galleries {
			verifyStringPtr(t, gallery.Path, pathCriterion)
			ratingCriterion.Modifier = models.CriterionModifierNotEquals
			verifyIntPtr(t, gallery.Rating, ratingCriterion)
		}

		return nil
	})
}

func TestGalleryIllegalQuery(t *testing.T) {
	assert := assert.New(t)

	const galleryIdx = 1
	subFilter := models.GalleryFilterType{
		Path: &models.StringCriterionInput{
			Value:    getGalleryStringValue(galleryIdx, "Path"),
			Modifier: models.CriterionModifierEquals,
		},
	}

	galleryFilter := &models.GalleryFilterType{
		And: &subFilter,
		Or:  &subFilter,
	}

	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter

		_, _, err := sqb.Query(ctx, galleryFilter, nil)
		assert.NotNil(err)

		galleryFilter.Or = nil
		galleryFilter.Not = &subFilter
		_, _, err = sqb.Query(ctx, galleryFilter, nil)
		assert.NotNil(err)

		galleryFilter.And = nil
		galleryFilter.Or = &subFilter
		_, _, err = sqb.Query(ctx, galleryFilter, nil)
		assert.NotNil(err)

		return nil
	})
}

func TestGalleryQueryURL(t *testing.T) {
	const sceneIdx = 1
	galleryURL := getGalleryStringValue(sceneIdx, urlField)

	urlCriterion := models.StringCriterionInput{
		Value:    galleryURL,
		Modifier: models.CriterionModifierEquals,
	}

	filter := models.GalleryFilterType{
		URL: &urlCriterion,
	}

	verifyFn := func(g *models.Gallery) {
		t.Helper()
		verifyString(t, g.URL, urlCriterion)
	}

	verifyGalleryQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotEquals
	verifyGalleryQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierMatchesRegex
	urlCriterion.Value = "gallery_.*1_URL"
	verifyGalleryQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotMatchesRegex
	verifyGalleryQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierIsNull
	urlCriterion.Value = ""
	verifyGalleryQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotNull
	verifyGalleryQuery(t, filter, verifyFn)
}

func verifyGalleryQuery(t *testing.T, filter models.GalleryFilterType, verifyFn func(s *models.Gallery)) {
	withTxn(func(ctx context.Context) error {
		t.Helper()
		sqb := sqlite.GalleryReaderWriter

		galleries := queryGallery(ctx, t, sqb, &filter, nil)

		// assume it should find at least one
		assert.Greater(t, len(galleries), 0)

		for _, gallery := range galleries {
			verifyFn(gallery)
		}

		return nil
	})
}

func TestGalleryQueryRating(t *testing.T) {
	const rating = 3
	ratingCriterion := models.IntCriterionInput{
		Value:    rating,
		Modifier: models.CriterionModifierEquals,
	}

	verifyGalleriesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotEquals
	verifyGalleriesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyGalleriesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierLessThan
	verifyGalleriesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierIsNull
	verifyGalleriesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotNull
	verifyGalleriesRating(t, ratingCriterion)
}

func verifyGalleriesRating(t *testing.T, ratingCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter
		galleryFilter := models.GalleryFilterType{
			Rating: &ratingCriterion,
		}

		galleries, _, err := sqb.Query(ctx, &galleryFilter, nil)
		if err != nil {
			t.Errorf("Error querying gallery: %s", err.Error())
		}

		for _, gallery := range galleries {
			verifyIntPtr(t, gallery.Rating, ratingCriterion)
		}

		return nil
	})
}

func TestGalleryQueryIsMissingScene(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		qb := sqlite.GalleryReaderWriter
		isMissing := "scenes"
		galleryFilter := models.GalleryFilterType{
			IsMissing: &isMissing,
		}

		q := getGalleryStringValue(galleryIdxWithScene, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		galleries, _, err := qb.Query(ctx, &galleryFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying gallery: %s", err.Error())
		}

		assert.Len(t, galleries, 0)

		findFilter.Q = nil
		galleries, _, err = qb.Query(ctx, &galleryFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying gallery: %s", err.Error())
		}

		// ensure non of the ids equal the one with gallery
		for _, gallery := range galleries {
			assert.NotEqual(t, galleryIDs[galleryIdxWithScene], gallery.ID)
		}

		return nil
	})
}

func queryGallery(ctx context.Context, t *testing.T, sqb models.GalleryReader, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) []*models.Gallery {
	galleries, _, err := sqb.Query(ctx, galleryFilter, findFilter)
	if err != nil {
		t.Errorf("Error querying gallery: %s", err.Error())
	}

	return galleries
}

func TestGalleryQueryIsMissingStudio(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter
		isMissing := "studio"
		galleryFilter := models.GalleryFilterType{
			IsMissing: &isMissing,
		}

		q := getGalleryStringValue(galleryIdxWithStudio, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		galleries := queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)

		assert.Len(t, galleries, 0)

		findFilter.Q = nil
		galleries = queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)

		// ensure non of the ids equal the one with studio
		for _, gallery := range galleries {
			assert.NotEqual(t, galleryIDs[galleryIdxWithStudio], gallery.ID)
		}

		return nil
	})
}

func TestGalleryQueryIsMissingPerformers(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter
		isMissing := "performers"
		galleryFilter := models.GalleryFilterType{
			IsMissing: &isMissing,
		}

		q := getGalleryStringValue(galleryIdxWithPerformer, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		galleries := queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)

		assert.Len(t, galleries, 0)

		findFilter.Q = nil
		galleries = queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)

		assert.True(t, len(galleries) > 0)

		// ensure non of the ids equal the one with movies
		for _, gallery := range galleries {
			assert.NotEqual(t, galleryIDs[galleryIdxWithPerformer], gallery.ID)
		}

		return nil
	})
}

func TestGalleryQueryIsMissingTags(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter
		isMissing := "tags"
		galleryFilter := models.GalleryFilterType{
			IsMissing: &isMissing,
		}

		q := getGalleryStringValue(galleryIdxWithTwoTags, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		galleries := queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)

		assert.Len(t, galleries, 0)

		findFilter.Q = nil
		galleries = queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)

		assert.True(t, len(galleries) > 0)

		return nil
	})
}

func TestGalleryQueryIsMissingDate(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter
		isMissing := "date"
		galleryFilter := models.GalleryFilterType{
			IsMissing: &isMissing,
		}

		galleries := queryGallery(ctx, t, sqb, &galleryFilter, nil)

		// three in four scenes have no date
		assert.Len(t, galleries, int(math.Ceil(float64(totalGalleries)/4*3)))

		// ensure date is null, empty or "0001-01-01"
		for _, g := range galleries {
			assert.True(t, g.Date == nil || g.Date.Time == time.Time{})
		}

		return nil
	})
}

func TestGalleryQueryPerformers(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter
		performerCriterion := models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(performerIDs[performerIdxWithGallery]),
				strconv.Itoa(performerIDs[performerIdx1WithGallery]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		galleryFilter := models.GalleryFilterType{
			Performers: &performerCriterion,
		}

		galleries := queryGallery(ctx, t, sqb, &galleryFilter, nil)

		assert.Len(t, galleries, 2)

		// ensure ids are correct
		for _, gallery := range galleries {
			assert.True(t, gallery.ID == galleryIDs[galleryIdxWithPerformer] || gallery.ID == galleryIDs[galleryIdxWithTwoPerformers])
		}

		performerCriterion = models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(performerIDs[performerIdx1WithGallery]),
				strconv.Itoa(performerIDs[performerIdx2WithGallery]),
			},
			Modifier: models.CriterionModifierIncludesAll,
		}

		galleries = queryGallery(ctx, t, sqb, &galleryFilter, nil)

		assert.Len(t, galleries, 1)
		assert.Equal(t, galleryIDs[galleryIdxWithTwoPerformers], galleries[0].ID)

		performerCriterion = models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(performerIDs[performerIdx1WithGallery]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getGalleryStringValue(galleryIdxWithTwoPerformers, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		galleries = queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)
		assert.Len(t, galleries, 0)

		return nil
	})
}

func TestGalleryQueryTags(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter
		tagCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdxWithGallery]),
				strconv.Itoa(tagIDs[tagIdx1WithGallery]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		galleryFilter := models.GalleryFilterType{
			Tags: &tagCriterion,
		}

		galleries := queryGallery(ctx, t, sqb, &galleryFilter, nil)
		assert.Len(t, galleries, 2)

		// ensure ids are correct
		for _, gallery := range galleries {
			assert.True(t, gallery.ID == galleryIDs[galleryIdxWithTag] || gallery.ID == galleryIDs[galleryIdxWithTwoTags])
		}

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithGallery]),
				strconv.Itoa(tagIDs[tagIdx2WithGallery]),
			},
			Modifier: models.CriterionModifierIncludesAll,
		}

		galleries = queryGallery(ctx, t, sqb, &galleryFilter, nil)

		assert.Len(t, galleries, 1)
		assert.Equal(t, galleryIDs[galleryIdxWithTwoTags], galleries[0].ID)

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithGallery]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getGalleryStringValue(galleryIdxWithTwoTags, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		galleries = queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)
		assert.Len(t, galleries, 0)

		return nil
	})
}

func TestGalleryQueryStudio(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter
		studioCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithGallery]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		galleryFilter := models.GalleryFilterType{
			Studios: &studioCriterion,
		}

		galleries := queryGallery(ctx, t, sqb, &galleryFilter, nil)

		assert.Len(t, galleries, 1)

		// ensure id is correct
		assert.Equal(t, galleryIDs[galleryIdxWithStudio], galleries[0].ID)

		studioCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithGallery]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getGalleryStringValue(galleryIdxWithStudio, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		galleries = queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)
		assert.Len(t, galleries, 0)

		return nil
	})
}

func TestGalleryQueryStudioDepth(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter
		depth := 2
		studioCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithGrandChild]),
			},
			Modifier: models.CriterionModifierIncludes,
			Depth:    &depth,
		}

		galleryFilter := models.GalleryFilterType{
			Studios: &studioCriterion,
		}

		galleries := queryGallery(ctx, t, sqb, &galleryFilter, nil)
		assert.Len(t, galleries, 1)

		depth = 1

		galleries = queryGallery(ctx, t, sqb, &galleryFilter, nil)
		assert.Len(t, galleries, 0)

		studioCriterion.Value = []string{strconv.Itoa(studioIDs[studioIdxWithParentAndChild])}
		galleries = queryGallery(ctx, t, sqb, &galleryFilter, nil)
		assert.Len(t, galleries, 1)

		// ensure id is correct
		assert.Equal(t, galleryIDs[galleryIdxWithGrandChildStudio], galleries[0].ID)

		depth = 2

		studioCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithGrandChild]),
			},
			Modifier: models.CriterionModifierExcludes,
			Depth:    &depth,
		}

		q := getGalleryStringValue(galleryIdxWithGrandChildStudio, pathField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		galleries = queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)
		assert.Len(t, galleries, 0)

		depth = 1
		galleries = queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)
		assert.Len(t, galleries, 1)

		studioCriterion.Value = []string{strconv.Itoa(studioIDs[studioIdxWithParentAndChild])}
		galleries = queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)
		assert.Len(t, galleries, 0)

		return nil
	})
}

func TestGalleryQueryPerformerTags(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter
		tagCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdxWithPerformer]),
				strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		galleryFilter := models.GalleryFilterType{
			PerformerTags: &tagCriterion,
		}

		galleries := queryGallery(ctx, t, sqb, &galleryFilter, nil)
		assert.Len(t, galleries, 2)

		// ensure ids are correct
		for _, gallery := range galleries {
			assert.True(t, gallery.ID == galleryIDs[galleryIdxWithPerformerTag] || gallery.ID == galleryIDs[galleryIdxWithPerformerTwoTags])
		}

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
				strconv.Itoa(tagIDs[tagIdx2WithPerformer]),
			},
			Modifier: models.CriterionModifierIncludesAll,
		}

		galleries = queryGallery(ctx, t, sqb, &galleryFilter, nil)

		assert.Len(t, galleries, 1)
		assert.Equal(t, galleryIDs[galleryIdxWithPerformerTwoTags], galleries[0].ID)

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getGalleryStringValue(galleryIdxWithPerformerTwoTags, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		galleries = queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)
		assert.Len(t, galleries, 0)

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Modifier: models.CriterionModifierIsNull,
		}
		q = getGalleryStringValue(galleryIdx1WithImage, titleField)

		galleries = queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)
		assert.Len(t, galleries, 1)
		assert.Equal(t, galleryIDs[galleryIdx1WithImage], galleries[0].ID)

		q = getGalleryStringValue(galleryIdxWithPerformerTag, titleField)
		galleries = queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)
		assert.Len(t, galleries, 0)

		tagCriterion.Modifier = models.CriterionModifierNotNull

		galleries = queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)
		assert.Len(t, galleries, 1)
		assert.Equal(t, galleryIDs[galleryIdxWithPerformerTag], galleries[0].ID)

		q = getGalleryStringValue(galleryIdx1WithImage, titleField)
		galleries = queryGallery(ctx, t, sqb, &galleryFilter, &findFilter)
		assert.Len(t, galleries, 0)

		return nil
	})
}

func TestGalleryQueryTagCount(t *testing.T) {
	const tagCount = 1
	tagCountCriterion := models.IntCriterionInput{
		Value:    tagCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyGalleriesTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyGalleriesTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyGalleriesTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyGalleriesTagCount(t, tagCountCriterion)
}

func verifyGalleriesTagCount(t *testing.T, tagCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter
		galleryFilter := models.GalleryFilterType{
			TagCount: &tagCountCriterion,
		}

		galleries := queryGallery(ctx, t, sqb, &galleryFilter, nil)
		assert.Greater(t, len(galleries), 0)

		for _, gallery := range galleries {
			verifyInt(t, len(gallery.TagIDs), tagCountCriterion)
		}

		return nil
	})
}

func TestGalleryQueryPerformerCount(t *testing.T) {
	const performerCount = 1
	performerCountCriterion := models.IntCriterionInput{
		Value:    performerCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyGalleriesPerformerCount(t, performerCountCriterion)

	performerCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyGalleriesPerformerCount(t, performerCountCriterion)

	performerCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyGalleriesPerformerCount(t, performerCountCriterion)

	performerCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyGalleriesPerformerCount(t, performerCountCriterion)
}

func verifyGalleriesPerformerCount(t *testing.T, performerCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter
		galleryFilter := models.GalleryFilterType{
			PerformerCount: &performerCountCriterion,
		}

		galleries := queryGallery(ctx, t, sqb, &galleryFilter, nil)
		assert.Greater(t, len(galleries), 0)

		for _, gallery := range galleries {
			verifyInt(t, len(gallery.PerformerIDs), performerCountCriterion)
		}

		return nil
	})
}

func TestGalleryQueryAverageResolution(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		qb := sqlite.GalleryReaderWriter
		resolution := models.ResolutionEnumLow
		galleryFilter := models.GalleryFilterType{
			AverageResolution: &models.ResolutionCriterionInput{
				Value:    resolution,
				Modifier: models.CriterionModifierEquals,
			},
		}

		// not verifying average - just ensure we get at least one
		galleries := queryGallery(ctx, t, qb, &galleryFilter, nil)
		assert.Greater(t, len(galleries), 0)

		return nil
	})
}

func TestGalleryQueryImageCount(t *testing.T) {
	const imageCount = 0
	imageCountCriterion := models.IntCriterionInput{
		Value:    imageCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyGalleriesImageCount(t, imageCountCriterion)

	imageCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyGalleriesImageCount(t, imageCountCriterion)

	imageCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyGalleriesImageCount(t, imageCountCriterion)

	imageCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyGalleriesImageCount(t, imageCountCriterion)
}

func verifyGalleriesImageCount(t *testing.T, imageCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.GalleryReaderWriter
		galleryFilter := models.GalleryFilterType{
			ImageCount: &imageCountCriterion,
		}

		galleries := queryGallery(ctx, t, sqb, &galleryFilter, nil)
		assert.Greater(t, len(galleries), -1)

		for _, gallery := range galleries {
			pp := 0

			result, err := sqlite.ImageReaderWriter.Query(ctx, models.ImageQueryOptions{
				QueryOptions: models.QueryOptions{
					FindFilter: &models.FindFilterType{
						PerPage: &pp,
					},
					Count: true,
				},
				ImageFilter: &models.ImageFilterType{
					Galleries: &models.MultiCriterionInput{
						Value:    []string{strconv.Itoa(gallery.ID)},
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

// TODO Count
// TODO All
// TODO Query
// TODO Destroy
