//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stretchr/testify/assert"
)

func Test_imageQueryBuilder_Create(t *testing.T) {
	var (
		path              = "path"
		title             = "title"
		checksum          = "checksum"
		rating            = 3
		ocounter          = 5
		size        int64 = 1234
		width             = 640
		height            = 480
		fileModTime       = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		createdAt         = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt         = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name      string
		newObject models.Image
		wantErr   bool
	}{
		{
			"full",
			models.Image{
				Path:         path,
				Checksum:     checksum,
				Title:        title,
				Rating:       &rating,
				Organized:    true,
				OCounter:     ocounter,
				Size:         &size,
				Width:        &width,
				Height:       &height,
				StudioID:     &studioIDs[studioIdxWithImage],
				FileModTime:  &fileModTime,
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				GalleryIDs:   []int{galleryIDs[galleryIdxWithImage]},
				TagIDs:       []int{tagIDs[tagIdx1WithImage], tagIDs[tagIdx1WithDupName]},
				PerformerIDs: []int{performerIDs[performerIdx1WithImage], performerIDs[performerIdx1WithDupName]},
			},
			false,
		},
		{
			"invalid studio id",
			models.Image{
				StudioID: &invalidID,
			},
			true,
		},
		{
			"invalid gallery id",
			models.Image{
				GalleryIDs: []int{invalidID},
			},
			true,
		},
		{
			"invalid tag id",
			models.Image{
				TagIDs: []int{invalidID},
			},
			true,
		},
		{
			"invalid performer id",
			models.Image{
				PerformerIDs: []int{invalidID},
			},
			true,
		},
	}

	qb := sqlite.ImageReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			s := tt.newObject
			if err := qb.Create(ctx, &s); (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.Create() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				assert.Zero(s.ID)
				return
			}

			assert.NotZero(s.ID)

			copy := tt.newObject
			copy.ID = s.ID

			assert.Equal(copy, s)

			// ensure can find the image
			found, err := qb.Find(ctx, s.ID)
			if err != nil {
				t.Errorf("imageQueryBuilder.Find() error = %v", err)
			}

			assert.Equal(copy, *found)

			return
		})
	}
}

func Test_imageQueryBuilder_Update(t *testing.T) {
	var (
		title             = "title"
		checksum          = "checksum"
		rating            = 3
		ocounter          = 5
		size        int64 = 1234
		width             = 640
		height            = 480
		fileModTime       = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		createdAt         = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt         = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name          string
		updatedObject *models.Image
		wantErr       bool
	}{
		{
			"full",
			&models.Image{
				ID:           imageIDs[imageIdxWithGallery],
				Checksum:     checksum,
				Title:        title,
				Rating:       &rating,
				Organized:    true,
				OCounter:     ocounter,
				Size:         &size,
				Width:        &width,
				Height:       &height,
				StudioID:     &studioIDs[studioIdxWithImage],
				FileModTime:  &fileModTime,
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				GalleryIDs:   []int{galleryIDs[galleryIdxWithImage]},
				TagIDs:       []int{tagIDs[tagIdx1WithImage], tagIDs[tagIdx1WithDupName]},
				PerformerIDs: []int{performerIDs[performerIdx1WithImage], performerIDs[performerIdx1WithDupName]},
			},
			false,
		},
		{
			"clear nullables",
			&models.Image{
				ID:        imageIDs[imageIdxWithGallery],
				Checksum:  checksum,
				Organized: true,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"clear gallery ids",
			&models.Image{
				ID:        imageIDs[imageIdxWithGallery],
				Checksum:  checksum,
				Organized: true,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"clear tag ids",
			&models.Image{
				ID:        imageIDs[imageIdxWithTag],
				Checksum:  checksum,
				Organized: true,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"clear performer ids",
			&models.Image{
				ID:        imageIDs[imageIdxWithPerformer],
				Checksum:  checksum,
				Organized: true,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"invalid studio id",
			&models.Image{
				ID:        imageIDs[imageIdxWithGallery],
				Checksum:  checksum,
				Organized: true,
				StudioID:  &invalidID,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			true,
		},
		{
			"invalid gallery id",
			&models.Image{
				ID:         imageIDs[imageIdxWithGallery],
				Checksum:   checksum,
				Organized:  true,
				GalleryIDs: []int{invalidID},
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
			},
			true,
		},
		{
			"invalid tag id",
			&models.Image{
				ID:        imageIDs[imageIdxWithGallery],
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
			&models.Image{
				ID:           imageIDs[imageIdxWithGallery],
				Checksum:     checksum,
				Organized:    true,
				PerformerIDs: []int{invalidID},
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
			},
			true,
		},
	}

	qb := sqlite.ImageReaderWriter
	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			copy := *tt.updatedObject

			if err := qb.Update(ctx, tt.updatedObject); (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			s, err := qb.Find(ctx, tt.updatedObject.ID)
			if err != nil {
				t.Errorf("imageQueryBuilder.Find() error = %v", err)
			}

			assert.Equal(copy, *s)

			return
		})
	}
}

func clearImagePartial() models.ImagePartial {
	// leave mandatory fields
	return models.ImagePartial{
		Title:        models.OptionalString{Set: true, Null: true},
		Rating:       models.OptionalInt{Set: true, Null: true},
		Size:         models.OptionalInt64{Set: true, Null: true},
		Width:        models.OptionalInt{Set: true, Null: true},
		Height:       models.OptionalInt{Set: true, Null: true},
		StudioID:     models.OptionalInt{Set: true, Null: true},
		FileModTime:  models.OptionalTime{Set: true, Null: true},
		GalleryIDs:   &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
		TagIDs:       &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
		PerformerIDs: &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
	}
}

func Test_imageQueryBuilder_UpdatePartial(t *testing.T) {
	var (
		path              = "path"
		title             = "title"
		checksum          = "checksum"
		rating            = 3
		ocounter          = 5
		size        int64 = 1234
		width             = 640
		height            = 480
		fileModTime       = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		createdAt         = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt         = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name    string
		id      int
		partial models.ImagePartial
		want    models.Image
		wantErr bool
	}{
		{
			"full",
			imageIDs[imageIdx1WithGallery],
			models.ImagePartial{
				Path:        models.NewOptionalString(path),
				Checksum:    models.NewOptionalString(checksum),
				Title:       models.NewOptionalString(title),
				Rating:      models.NewOptionalInt(rating),
				Organized:   models.NewOptionalBool(true),
				OCounter:    models.NewOptionalInt(ocounter),
				Size:        models.NewOptionalInt64(size),
				Width:       models.NewOptionalInt(width),
				Height:      models.NewOptionalInt(height),
				StudioID:    models.NewOptionalInt(studioIDs[studioIdxWithImage]),
				FileModTime: models.NewOptionalTime(fileModTime),
				CreatedAt:   models.NewOptionalTime(createdAt),
				UpdatedAt:   models.NewOptionalTime(updatedAt),
				GalleryIDs: &models.UpdateIDs{
					IDs:  []int{galleryIDs[galleryIdxWithImage]},
					Mode: models.RelationshipUpdateModeSet,
				},
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithImage], tagIDs[tagIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeSet,
				},
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithImage], performerIDs[performerIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeSet,
				},
			},
			models.Image{
				ID:           imageIDs[imageIdx1WithGallery],
				Path:         path,
				Checksum:     checksum,
				Title:        title,
				Rating:       &rating,
				Organized:    true,
				OCounter:     ocounter,
				Size:         &size,
				Width:        &width,
				Height:       &height,
				StudioID:     &studioIDs[studioIdxWithImage],
				FileModTime:  &fileModTime,
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				GalleryIDs:   []int{galleryIDs[galleryIdxWithImage]},
				TagIDs:       []int{tagIDs[tagIdx1WithImage], tagIDs[tagIdx1WithDupName]},
				PerformerIDs: []int{performerIDs[performerIdx1WithImage], performerIDs[performerIdx1WithDupName]},
			},
			false,
		},
		{
			"clear all",
			imageIDs[imageIdx1WithGallery],
			clearImagePartial(),
			models.Image{
				ID:       imageIDs[imageIdx1WithGallery],
				Path:     getImageStringValue(imageIdx1WithGallery, pathField),
				Checksum: getImageStringValue(imageIdx1WithGallery, checksumField),
				OCounter: getOCounter(imageIdx1WithGallery),
			},
			false,
		},
		{
			"invalid id",
			invalidID,
			models.ImagePartial{},
			models.Image{},
			true,
		},
	}
	for _, tt := range tests {
		qb := sqlite.ImageReaderWriter

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			got, err := qb.UpdatePartial(ctx, tt.id, tt.partial)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.UpdatePartial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			assert.Equal(tt.want, *got)

			s, err := qb.Find(ctx, tt.id)
			if err != nil {
				t.Errorf("imageQueryBuilder.Find() error = %v", err)
			}

			assert.Equal(tt.want, *s)
		})
	}
}

func Test_imageQueryBuilder_UpdatePartialRelationships(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		partial models.ImagePartial
		want    models.Image
		wantErr bool
	}{
		{
			"add galleries",
			imageIDs[imageIdxWithGallery],
			models.ImagePartial{
				GalleryIDs: &models.UpdateIDs{
					IDs:  []int{galleryIDs[galleryIdx1WithImage], galleryIDs[galleryIdx1WithPerformer]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Image{
				GalleryIDs: append(indexesToIDs(galleryIDs, imageGalleries[imageIdxWithGallery]),
					galleryIDs[galleryIdx1WithImage],
					galleryIDs[galleryIdx1WithPerformer],
				),
			},
			false,
		},
		{
			"add tags",
			imageIDs[imageIdxWithTwoTags],
			models.ImagePartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithGallery]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Image{
				TagIDs: append(indexesToIDs(tagIDs, imageTags[imageIdxWithTwoTags]),
					tagIDs[tagIdx1WithDupName],
					tagIDs[tagIdx1WithGallery],
				),
			},
			false,
		},
		{
			"add performers",
			imageIDs[imageIdxWithTwoPerformers],
			models.ImagePartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithDupName], performerIDs[performerIdx1WithGallery]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Image{
				PerformerIDs: append(indexesToIDs(performerIDs, imagePerformers[imageIdxWithTwoPerformers]),
					performerIDs[performerIdx1WithDupName],
					performerIDs[performerIdx1WithGallery],
				),
			},
			false,
		},
		{
			"add duplicate galleries",
			imageIDs[imageIdxWithGallery],
			models.ImagePartial{
				GalleryIDs: &models.UpdateIDs{
					IDs:  []int{galleryIDs[galleryIdxWithImage], galleryIDs[galleryIdx1WithPerformer]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Image{
				GalleryIDs: append(indexesToIDs(galleryIDs, imageGalleries[imageIdxWithGallery]),
					galleryIDs[galleryIdx1WithPerformer],
				),
			},
			false,
		},
		{
			"add duplicate tags",
			imageIDs[imageIdxWithTwoTags],
			models.ImagePartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithImage], tagIDs[tagIdx1WithGallery]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Image{
				TagIDs: append(indexesToIDs(tagIDs, imageTags[imageIdxWithTwoTags]),
					tagIDs[tagIdx1WithGallery],
				),
			},
			false,
		},
		{
			"add duplicate performers",
			imageIDs[imageIdxWithTwoPerformers],
			models.ImagePartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithImage], performerIDs[performerIdx1WithGallery]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Image{
				PerformerIDs: append(indexesToIDs(performerIDs, imagePerformers[imageIdxWithTwoPerformers]),
					performerIDs[performerIdx1WithGallery],
				),
			},
			false,
		},
		{
			"add invalid galleries",
			imageIDs[imageIdxWithGallery],
			models.ImagePartial{
				GalleryIDs: &models.UpdateIDs{
					IDs:  []int{invalidID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Image{},
			true,
		},
		{
			"add invalid tags",
			imageIDs[imageIdxWithTwoTags],
			models.ImagePartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{invalidID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Image{},
			true,
		},
		{
			"add invalid performers",
			imageIDs[imageIdxWithTwoPerformers],
			models.ImagePartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{invalidID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Image{},
			true,
		},
		{
			"remove galleries",
			imageIDs[imageIdxWithGallery],
			models.ImagePartial{
				GalleryIDs: &models.UpdateIDs{
					IDs:  []int{galleryIDs[galleryIdxWithImage]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Image{},
			false,
		},
		{
			"remove tags",
			imageIDs[imageIdxWithTwoTags],
			models.ImagePartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithImage]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Image{
				TagIDs: []int{tagIDs[tagIdx2WithImage]},
			},
			false,
		},
		{
			"remove performers",
			imageIDs[imageIdxWithTwoPerformers],
			models.ImagePartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithImage]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Image{
				PerformerIDs: []int{performerIDs[performerIdx2WithImage]},
			},
			false,
		},
		{
			"remove unrelated galleries",
			imageIDs[imageIdxWithGallery],
			models.ImagePartial{
				GalleryIDs: &models.UpdateIDs{
					IDs:  []int{galleryIDs[galleryIdx1WithImage]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Image{
				GalleryIDs: []int{galleryIDs[galleryIdxWithImage]},
			},
			false,
		},
		{
			"remove unrelated tags",
			imageIDs[imageIdxWithTwoTags],
			models.ImagePartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithPerformer]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Image{
				TagIDs: indexesToIDs(tagIDs, imageTags[imageIdxWithTwoTags]),
			},
			false,
		},
		{
			"remove unrelated performers",
			imageIDs[imageIdxWithTwoPerformers],
			models.ImagePartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Image{
				PerformerIDs: indexesToIDs(performerIDs, imagePerformers[imageIdxWithTwoPerformers]),
			},
			false,
		},
	}

	for _, tt := range tests {
		qb := sqlite.ImageReaderWriter

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			got, err := qb.UpdatePartial(ctx, tt.id, tt.partial)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.UpdatePartial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			s, err := qb.Find(ctx, tt.id)
			if err != nil {
				t.Errorf("imageQueryBuilder.Find() error = %v", err)
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
		})
	}
}

func Test_imageQueryBuilder_IncrementOCounter(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    int
		wantErr bool
	}{
		{
			"increment",
			imageIDs[1],
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

	qb := sqlite.ImageReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.IncrementOCounter(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.IncrementOCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("imageQueryBuilder.IncrementOCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_imageQueryBuilder_DecrementOCounter(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    int
		wantErr bool
	}{
		{
			"decrement",
			imageIDs[2],
			1,
			false,
		},
		{
			"zero",
			imageIDs[0],
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

	qb := sqlite.ImageReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.DecrementOCounter(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.DecrementOCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("imageQueryBuilder.DecrementOCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_imageQueryBuilder_ResetOCounter(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    int
		wantErr bool
	}{
		{
			"decrement",
			imageIDs[2],
			0,
			false,
		},
		{
			"zero",
			imageIDs[0],
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

	qb := sqlite.ImageReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.ResetOCounter(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.ResetOCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("imageQueryBuilder.ResetOCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_imageQueryBuilder_Destroy(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			"valid",
			imageIDs[imageIdxWithGallery],
			false,
		},
		{
			"invalid",
			invalidID,
			true,
		},
	}

	qb := sqlite.ImageReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			withRollbackTxn(func(ctx context.Context) error {
				if err := qb.Destroy(ctx, tt.id); (err != nil) != tt.wantErr {
					t.Errorf("imageQueryBuilder.Destroy() error = %v, wantErr %v", err, tt.wantErr)
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

func makeImageWithID(index int) *models.Image {
	ret := makeImage(index)
	ret.ID = imageIDs[index]
	return ret
}

func Test_imageQueryBuilder_Find(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    *models.Image
		wantErr bool
	}{
		{
			"valid",
			imageIDs[imageIdxWithGallery],
			makeImageWithID(imageIdxWithGallery),
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
			imageIDs[imageIdxWithTwoPerformers],
			makeImageWithID(imageIdxWithTwoPerformers),
			false,
		},
		{
			"with tags",
			imageIDs[imageIdxWithTwoTags],
			makeImageWithID(imageIdxWithTwoTags),
			false,
		},
	}

	qb := sqlite.ImageReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.Find(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func Test_imageQueryBuilder_FindMany(t *testing.T) {
	tests := []struct {
		name    string
		ids     []int
		want    []*models.Image
		wantErr bool
	}{
		{
			"valid with relationships",
			[]int{imageIDs[imageIdxWithGallery], imageIDs[imageIdxWithTwoPerformers], imageIDs[imageIdxWithTwoTags]},
			[]*models.Image{
				makeImageWithID(imageIdxWithGallery),
				makeImageWithID(imageIdxWithTwoPerformers),
				makeImageWithID(imageIdxWithTwoTags),
			},
			false,
		},
		{
			"invalid",
			[]int{imageIDs[imageIdxWithGallery], imageIDs[imageIdxWithTwoPerformers], invalidID},
			nil,
			true,
		},
	}

	qb := sqlite.ImageReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.FindMany(ctx, tt.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.FindMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("imageQueryBuilder.FindMany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_imageQueryBuilder_FindByChecksum(t *testing.T) {
	getChecksum := func(index int) string {
		return getImageStringValue(index, checksumField)
	}

	tests := []struct {
		name     string
		checksum string
		want     *models.Image
		wantErr  bool
	}{
		{
			"valid",
			getChecksum(imageIdxWithGallery),
			makeImageWithID(imageIdxWithGallery),
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
			getChecksum(imageIdxWithTwoPerformers),
			makeImageWithID(imageIdxWithTwoPerformers),
			false,
		},
		{
			"with tags",
			getChecksum(imageIdxWithTwoTags),
			makeImageWithID(imageIdxWithTwoTags),
			false,
		},
	}

	qb := sqlite.ImageReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.FindByChecksum(ctx, tt.checksum)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.FindByChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("imageQueryBuilder.FindByChecksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_imageQueryBuilder_FindByPath(t *testing.T) {
	getPath := func(index int) string {
		return getImageStringValue(index, pathField)
	}

	tests := []struct {
		name    string
		path    string
		want    *models.Image
		wantErr bool
	}{
		{
			"valid",
			getPath(imageIdxWithGallery),
			makeImageWithID(imageIdxWithGallery),
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
			getPath(imageIdxWithTwoPerformers),
			makeImageWithID(imageIdxWithTwoPerformers),
			false,
		},
		{
			"with tags",
			getPath(imageIdxWithTwoTags),
			makeImageWithID(imageIdxWithTwoTags),
			false,
		},
	}

	qb := sqlite.ImageReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.FindByPath(ctx, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.FindByPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("imageQueryBuilder.FindByPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_imageQueryBuilder_FindByGalleryID(t *testing.T) {
	tests := []struct {
		name      string
		galleryID int
		want      []*models.Image
		wantErr   bool
	}{
		{
			"valid",
			galleryIDs[galleryIdxWithTwoImages],
			[]*models.Image{makeImageWithID(imageIdx1WithGallery), makeImageWithID(imageIdx2WithGallery)},
			false,
		},
		{
			"none",
			galleryIDs[galleryIdx1WithPerformer],
			nil,
			false,
		},
	}

	qb := sqlite.ImageReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByGalleryID(ctx, tt.galleryID)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.FindByGalleryID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(tt.want, got)
			return
		})
	}
}

func Test_imageQueryBuilder_CountByGalleryID(t *testing.T) {
	tests := []struct {
		name      string
		galleryID int
		want      int
		wantErr   bool
	}{
		{
			"valid",
			galleryIDs[galleryIdxWithTwoImages],
			2,
			false,
		},
		{
			"none",
			galleryIDs[galleryIdx1WithPerformer],
			0,
			false,
		},
	}

	qb := sqlite.ImageReaderWriter

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.CountByGalleryID(ctx, tt.galleryID)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.CountByGalleryID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("imageQueryBuilder.CountByGalleryID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageQueryQ(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		const imageIdx = 2

		q := getImageStringValue(imageIdx, titleField)

		sqb := sqlite.ImageReaderWriter

		imageQueryQ(ctx, t, sqb, q, imageIdx)

		return nil
	})
}

func queryImagesWithCount(ctx context.Context, sqb models.ImageReader, imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) ([]*models.Image, int, error) {
	result, err := sqb.Query(ctx, models.ImageQueryOptions{
		QueryOptions: models.QueryOptions{
			FindFilter: findFilter,
			Count:      true,
		},
		ImageFilter: imageFilter,
	})
	if err != nil {
		return nil, 0, err
	}

	images, err := result.Resolve(ctx)
	if err != nil {
		return nil, 0, err
	}

	return images, result.Count, nil
}

func imageQueryQ(ctx context.Context, t *testing.T, sqb models.ImageReader, q string, expectedImageIdx int) {
	filter := models.FindFilterType{
		Q: &q,
	}
	images := queryImages(ctx, t, sqb, nil, &filter)

	assert.Len(t, images, 1)
	image := images[0]
	assert.Equal(t, imageIDs[expectedImageIdx], image.ID)

	count, err := sqb.QueryCount(ctx, nil, &filter)
	if err != nil {
		t.Errorf("Error querying image: %s", err.Error())
	}
	assert.Equal(t, len(images), count)

	// no Q should return all results
	filter.Q = nil
	images = queryImages(ctx, t, sqb, nil, &filter)

	assert.Len(t, images, totalImages)
}

func TestImageQueryPath(t *testing.T) {
	const imageIdx = 1
	imagePath := getImageStringValue(imageIdx, "Path")

	pathCriterion := models.StringCriterionInput{
		Value:    imagePath,
		Modifier: models.CriterionModifierEquals,
	}

	verifyImagePath(t, pathCriterion, 1)

	pathCriterion.Modifier = models.CriterionModifierNotEquals
	verifyImagePath(t, pathCriterion, totalImages-1)

	pathCriterion.Modifier = models.CriterionModifierMatchesRegex
	pathCriterion.Value = "image_.*01_Path"
	verifyImagePath(t, pathCriterion, 1) // TODO - 2 if zip path is included

	pathCriterion.Modifier = models.CriterionModifierNotMatchesRegex
	verifyImagePath(t, pathCriterion, totalImages-1) // TODO - -2 if zip path is included
}

func verifyImagePath(t *testing.T, pathCriterion models.StringCriterionInput, expected int) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		imageFilter := models.ImageFilterType{
			Path: &pathCriterion,
		}

		images := queryImages(ctx, t, sqb, &imageFilter, nil)

		assert.Equal(t, expected, len(images), "number of returned images")

		for _, image := range images {
			verifyString(t, image.Path, pathCriterion)
		}

		return nil
	})
}

func TestImageQueryPathOr(t *testing.T) {
	const image1Idx = 1
	const image2Idx = 2

	image1Path := getImageStringValue(image1Idx, "Path")
	image2Path := getImageStringValue(image2Idx, "Path")

	imageFilter := models.ImageFilterType{
		Path: &models.StringCriterionInput{
			Value:    image1Path,
			Modifier: models.CriterionModifierEquals,
		},
		Or: &models.ImageFilterType{
			Path: &models.StringCriterionInput{
				Value:    image2Path,
				Modifier: models.CriterionModifierEquals,
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter

		images := queryImages(ctx, t, sqb, &imageFilter, nil)

		assert.Len(t, images, 2)
		assert.Equal(t, image1Path, images[0].Path)
		assert.Equal(t, image2Path, images[1].Path)

		return nil
	})
}

func TestImageQueryPathAndRating(t *testing.T) {
	const imageIdx = 1
	imagePath := getImageStringValue(imageIdx, "Path")
	imageRating := getRating(imageIdx)

	imageFilter := models.ImageFilterType{
		Path: &models.StringCriterionInput{
			Value:    imagePath,
			Modifier: models.CriterionModifierEquals,
		},
		And: &models.ImageFilterType{
			Rating: &models.IntCriterionInput{
				Value:    int(imageRating.Int64),
				Modifier: models.CriterionModifierEquals,
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter

		images := queryImages(ctx, t, sqb, &imageFilter, nil)

		assert.Len(t, images, 1)
		assert.Equal(t, imagePath, images[0].Path)
		assert.Equal(t, int(imageRating.Int64), *images[0].Rating)

		return nil
	})
}

func TestImageQueryPathNotRating(t *testing.T) {
	const imageIdx = 1

	imageRating := getRating(imageIdx)

	pathCriterion := models.StringCriterionInput{
		Value:    "image_.*1_Path",
		Modifier: models.CriterionModifierMatchesRegex,
	}

	ratingCriterion := models.IntCriterionInput{
		Value:    int(imageRating.Int64),
		Modifier: models.CriterionModifierEquals,
	}

	imageFilter := models.ImageFilterType{
		Path: &pathCriterion,
		Not: &models.ImageFilterType{
			Rating: &ratingCriterion,
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter

		images := queryImages(ctx, t, sqb, &imageFilter, nil)

		for _, image := range images {
			verifyString(t, image.Path, pathCriterion)
			ratingCriterion.Modifier = models.CriterionModifierNotEquals
			verifyIntPtr(t, image.Rating, ratingCriterion)
		}

		return nil
	})
}

func TestImageIllegalQuery(t *testing.T) {
	assert := assert.New(t)

	const imageIdx = 1
	subFilter := models.ImageFilterType{
		Path: &models.StringCriterionInput{
			Value:    getImageStringValue(imageIdx, "Path"),
			Modifier: models.CriterionModifierEquals,
		},
	}

	imageFilter := &models.ImageFilterType{
		And: &subFilter,
		Or:  &subFilter,
	}

	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter

		_, _, err := queryImagesWithCount(ctx, sqb, imageFilter, nil)
		assert.NotNil(err)

		imageFilter.Or = nil
		imageFilter.Not = &subFilter
		_, _, err = queryImagesWithCount(ctx, sqb, imageFilter, nil)
		assert.NotNil(err)

		imageFilter.And = nil
		imageFilter.Or = &subFilter
		_, _, err = queryImagesWithCount(ctx, sqb, imageFilter, nil)
		assert.NotNil(err)

		return nil
	})
}

func TestImageQueryRating(t *testing.T) {
	const rating = 3
	ratingCriterion := models.IntCriterionInput{
		Value:    rating,
		Modifier: models.CriterionModifierEquals,
	}

	verifyImagesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotEquals
	verifyImagesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyImagesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierLessThan
	verifyImagesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierIsNull
	verifyImagesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotNull
	verifyImagesRating(t, ratingCriterion)
}

func verifyImagesRating(t *testing.T, ratingCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		imageFilter := models.ImageFilterType{
			Rating: &ratingCriterion,
		}

		images, _, err := queryImagesWithCount(ctx, sqb, &imageFilter, nil)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		for _, image := range images {
			verifyIntPtr(t, image.Rating, ratingCriterion)
		}

		return nil
	})
}

func TestImageQueryOCounter(t *testing.T) {
	const oCounter = 1
	oCounterCriterion := models.IntCriterionInput{
		Value:    oCounter,
		Modifier: models.CriterionModifierEquals,
	}

	verifyImagesOCounter(t, oCounterCriterion)

	oCounterCriterion.Modifier = models.CriterionModifierNotEquals
	verifyImagesOCounter(t, oCounterCriterion)

	oCounterCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyImagesOCounter(t, oCounterCriterion)

	oCounterCriterion.Modifier = models.CriterionModifierLessThan
	verifyImagesOCounter(t, oCounterCriterion)
}

func verifyImagesOCounter(t *testing.T, oCounterCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		imageFilter := models.ImageFilterType{
			OCounter: &oCounterCriterion,
		}

		images, _, err := queryImagesWithCount(ctx, sqb, &imageFilter, nil)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		for _, image := range images {
			verifyInt(t, image.OCounter, oCounterCriterion)
		}

		return nil
	})
}

func TestImageQueryResolution(t *testing.T) {
	verifyImagesResolution(t, models.ResolutionEnumLow)
	verifyImagesResolution(t, models.ResolutionEnumStandard)
	verifyImagesResolution(t, models.ResolutionEnumStandardHd)
	verifyImagesResolution(t, models.ResolutionEnumFullHd)
	verifyImagesResolution(t, models.ResolutionEnumFourK)
	verifyImagesResolution(t, models.ResolutionEnum("unknown"))
}

func verifyImagesResolution(t *testing.T, resolution models.ResolutionEnum) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		imageFilter := models.ImageFilterType{
			Resolution: &models.ResolutionCriterionInput{
				Value:    resolution,
				Modifier: models.CriterionModifierEquals,
			},
		}

		images, _, err := queryImagesWithCount(ctx, sqb, &imageFilter, nil)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		for _, image := range images {
			verifyImageResolution(t, image.Height, resolution)
		}

		return nil
	})
}

func verifyImageResolution(t *testing.T, height *int, resolution models.ResolutionEnum) {
	if !resolution.IsValid() {
		return
	}

	assert := assert.New(t)
	assert.NotNil(height)
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

func TestImageQueryIsMissingGalleries(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		isMissing := "galleries"
		imageFilter := models.ImageFilterType{
			IsMissing: &isMissing,
		}

		q := getImageStringValue(imageIdxWithGallery, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		images, _, err := queryImagesWithCount(ctx, sqb, &imageFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		assert.Len(t, images, 0)

		findFilter.Q = nil
		images, _, err = queryImagesWithCount(ctx, sqb, &imageFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		assert.Greater(t, len(images), 0)

		// ensure non of the ids equal the one with gallery
		for _, image := range images {
			assert.NotEqual(t, imageIDs[imageIdxWithGallery], image.ID)
		}

		return nil
	})
}

func TestImageQueryIsMissingStudio(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		isMissing := "studio"
		imageFilter := models.ImageFilterType{
			IsMissing: &isMissing,
		}

		q := getImageStringValue(imageIdxWithStudio, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		images, _, err := queryImagesWithCount(ctx, sqb, &imageFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		assert.Len(t, images, 0)

		findFilter.Q = nil
		images, _, err = queryImagesWithCount(ctx, sqb, &imageFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		// ensure non of the ids equal the one with studio
		for _, image := range images {
			assert.NotEqual(t, imageIDs[imageIdxWithStudio], image.ID)
		}

		return nil
	})
}

func TestImageQueryIsMissingPerformers(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		isMissing := "performers"
		imageFilter := models.ImageFilterType{
			IsMissing: &isMissing,
		}

		q := getImageStringValue(imageIdxWithPerformer, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		images, _, err := queryImagesWithCount(ctx, sqb, &imageFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		assert.Len(t, images, 0)

		findFilter.Q = nil
		images, _, err = queryImagesWithCount(ctx, sqb, &imageFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		assert.True(t, len(images) > 0)

		// ensure non of the ids equal the one with movies
		for _, image := range images {
			assert.NotEqual(t, imageIDs[imageIdxWithPerformer], image.ID)
		}

		return nil
	})
}

func TestImageQueryIsMissingTags(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		isMissing := "tags"
		imageFilter := models.ImageFilterType{
			IsMissing: &isMissing,
		}

		q := getImageStringValue(imageIdxWithTwoTags, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		images, _, err := queryImagesWithCount(ctx, sqb, &imageFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		assert.Len(t, images, 0)

		findFilter.Q = nil
		images, _, err = queryImagesWithCount(ctx, sqb, &imageFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		assert.True(t, len(images) > 0)

		return nil
	})
}

func TestImageQueryIsMissingRating(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		isMissing := "rating"
		imageFilter := models.ImageFilterType{
			IsMissing: &isMissing,
		}

		images, _, err := queryImagesWithCount(ctx, sqb, &imageFilter, nil)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		assert.True(t, len(images) > 0)

		// ensure rating is null
		for _, image := range images {
			assert.Nil(t, image.Rating)
		}

		return nil
	})
}

func TestImageQueryGallery(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		galleryCriterion := models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(galleryIDs[galleryIdxWithImage]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		imageFilter := models.ImageFilterType{
			Galleries: &galleryCriterion,
		}

		images := queryImages(ctx, t, sqb, &imageFilter, nil)
		assert.Len(t, images, 1)

		// ensure ids are correct
		for _, image := range images {
			assert.True(t, image.ID == imageIDs[imageIdxWithGallery])
		}

		galleryCriterion = models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(galleryIDs[galleryIdx1WithImage]),
				strconv.Itoa(galleryIDs[galleryIdx2WithImage]),
			},
			Modifier: models.CriterionModifierIncludesAll,
		}

		images = queryImages(ctx, t, sqb, &imageFilter, nil)

		assert.Len(t, images, 1)
		assert.Equal(t, imageIDs[imageIdxWithTwoGalleries], images[0].ID)

		galleryCriterion = models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(performerIDs[galleryIdx1WithImage]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getImageStringValue(imageIdxWithTwoGalleries, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 0)

		q = getImageStringValue(imageIdxWithPerformer, titleField)
		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 1)

		return nil
	})
}

func TestImageQueryPerformers(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		performerCriterion := models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(performerIDs[performerIdxWithImage]),
				strconv.Itoa(performerIDs[performerIdx1WithImage]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		imageFilter := models.ImageFilterType{
			Performers: &performerCriterion,
		}

		images := queryImages(ctx, t, sqb, &imageFilter, nil)
		assert.Len(t, images, 2)

		// ensure ids are correct
		for _, image := range images {
			assert.True(t, image.ID == imageIDs[imageIdxWithPerformer] || image.ID == imageIDs[imageIdxWithTwoPerformers])
		}

		performerCriterion = models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(performerIDs[performerIdx1WithImage]),
				strconv.Itoa(performerIDs[performerIdx2WithImage]),
			},
			Modifier: models.CriterionModifierIncludesAll,
		}

		images = queryImages(ctx, t, sqb, &imageFilter, nil)
		assert.Len(t, images, 1)
		assert.Equal(t, imageIDs[imageIdxWithTwoPerformers], images[0].ID)

		performerCriterion = models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(performerIDs[performerIdx1WithImage]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getImageStringValue(imageIdxWithTwoPerformers, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 0)

		performerCriterion = models.MultiCriterionInput{
			Modifier: models.CriterionModifierIsNull,
		}
		q = getImageStringValue(imageIdxWithGallery, titleField)

		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 1)
		assert.Equal(t, imageIDs[imageIdxWithGallery], images[0].ID)

		q = getImageStringValue(imageIdxWithPerformerTag, titleField)
		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 0)

		performerCriterion.Modifier = models.CriterionModifierNotNull

		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 1)
		assert.Equal(t, imageIDs[imageIdxWithPerformerTag], images[0].ID)

		q = getImageStringValue(imageIdxWithGallery, titleField)
		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 0)

		return nil
	})
}

func TestImageQueryTags(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		tagCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdxWithImage]),
				strconv.Itoa(tagIDs[tagIdx1WithImage]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		imageFilter := models.ImageFilterType{
			Tags: &tagCriterion,
		}

		images := queryImages(ctx, t, sqb, &imageFilter, nil)
		assert.Len(t, images, 2)

		// ensure ids are correct
		for _, image := range images {
			assert.True(t, image.ID == imageIDs[imageIdxWithTag] || image.ID == imageIDs[imageIdxWithTwoTags])
		}

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithImage]),
				strconv.Itoa(tagIDs[tagIdx2WithImage]),
			},
			Modifier: models.CriterionModifierIncludesAll,
		}

		images = queryImages(ctx, t, sqb, &imageFilter, nil)
		assert.Len(t, images, 1)
		assert.Equal(t, imageIDs[imageIdxWithTwoTags], images[0].ID)

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithImage]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getImageStringValue(imageIdxWithTwoTags, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 0)

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Modifier: models.CriterionModifierIsNull,
		}
		q = getImageStringValue(imageIdxWithGallery, titleField)

		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 1)
		assert.Equal(t, imageIDs[imageIdxWithGallery], images[0].ID)

		q = getImageStringValue(imageIdxWithTag, titleField)
		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 0)

		tagCriterion.Modifier = models.CriterionModifierNotNull

		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 1)
		assert.Equal(t, imageIDs[imageIdxWithTag], images[0].ID)

		q = getImageStringValue(imageIdxWithGallery, titleField)
		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 0)

		return nil
	})
}

func TestImageQueryStudio(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		studioCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithImage]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		imageFilter := models.ImageFilterType{
			Studios: &studioCriterion,
		}

		images, _, err := queryImagesWithCount(ctx, sqb, &imageFilter, nil)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		assert.Len(t, images, 1)

		// ensure id is correct
		assert.Equal(t, imageIDs[imageIdxWithStudio], images[0].ID)

		studioCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithImage]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getImageStringValue(imageIdxWithStudio, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		images, _, err = queryImagesWithCount(ctx, sqb, &imageFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}
		assert.Len(t, images, 0)

		return nil
	})
}

func TestImageQueryStudioDepth(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		depth := 2
		studioCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithGrandChild]),
			},
			Modifier: models.CriterionModifierIncludes,
			Depth:    &depth,
		}

		imageFilter := models.ImageFilterType{
			Studios: &studioCriterion,
		}

		images := queryImages(ctx, t, sqb, &imageFilter, nil)
		assert.Len(t, images, 1)

		depth = 1

		images = queryImages(ctx, t, sqb, &imageFilter, nil)
		assert.Len(t, images, 0)

		studioCriterion.Value = []string{strconv.Itoa(studioIDs[studioIdxWithParentAndChild])}
		images = queryImages(ctx, t, sqb, &imageFilter, nil)
		assert.Len(t, images, 1)

		// ensure id is correct
		assert.Equal(t, imageIDs[imageIdxWithGrandChildStudio], images[0].ID)

		depth = 2

		studioCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithGrandChild]),
			},
			Modifier: models.CriterionModifierExcludes,
			Depth:    &depth,
		}

		q := getImageStringValue(imageIdxWithGrandChildStudio, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 0)

		depth = 1
		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 1)

		studioCriterion.Value = []string{strconv.Itoa(studioIDs[studioIdxWithParentAndChild])}
		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 0)

		return nil
	})
}

func queryImages(ctx context.Context, t *testing.T, sqb models.ImageReader, imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) []*models.Image {
	images, _, err := queryImagesWithCount(ctx, sqb, imageFilter, findFilter)
	if err != nil {
		t.Errorf("Error querying images: %s", err.Error())
	}

	return images
}

func TestImageQueryPerformerTags(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		tagCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdxWithPerformer]),
				strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		imageFilter := models.ImageFilterType{
			PerformerTags: &tagCriterion,
		}

		images := queryImages(ctx, t, sqb, &imageFilter, nil)
		assert.Len(t, images, 2)

		// ensure ids are correct
		for _, image := range images {
			assert.True(t, image.ID == imageIDs[imageIdxWithPerformerTag] || image.ID == imageIDs[imageIdxWithPerformerTwoTags])
		}

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
				strconv.Itoa(tagIDs[tagIdx2WithPerformer]),
			},
			Modifier: models.CriterionModifierIncludesAll,
		}

		images = queryImages(ctx, t, sqb, &imageFilter, nil)

		assert.Len(t, images, 1)
		assert.Equal(t, imageIDs[imageIdxWithPerformerTwoTags], images[0].ID)

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getImageStringValue(imageIdxWithPerformerTwoTags, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 0)

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Modifier: models.CriterionModifierIsNull,
		}
		q = getImageStringValue(imageIdxWithGallery, titleField)

		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 1)
		assert.Equal(t, imageIDs[imageIdxWithGallery], images[0].ID)

		q = getImageStringValue(imageIdxWithPerformerTag, titleField)
		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 0)

		tagCriterion.Modifier = models.CriterionModifierNotNull

		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 1)
		assert.Equal(t, imageIDs[imageIdxWithPerformerTag], images[0].ID)

		q = getImageStringValue(imageIdxWithGallery, titleField)
		images = queryImages(ctx, t, sqb, &imageFilter, &findFilter)
		assert.Len(t, images, 0)

		return nil
	})
}

func TestImageQueryTagCount(t *testing.T) {
	const tagCount = 1
	tagCountCriterion := models.IntCriterionInput{
		Value:    tagCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyImagesTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyImagesTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyImagesTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyImagesTagCount(t, tagCountCriterion)
}

func verifyImagesTagCount(t *testing.T, tagCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		imageFilter := models.ImageFilterType{
			TagCount: &tagCountCriterion,
		}

		images := queryImages(ctx, t, sqb, &imageFilter, nil)
		assert.Greater(t, len(images), 0)

		for _, image := range images {
			ids, err := sqb.GetTagIDs(ctx, image.ID)
			if err != nil {
				return err
			}
			verifyInt(t, len(ids), tagCountCriterion)
		}

		return nil
	})
}

func TestImageQueryPerformerCount(t *testing.T) {
	const performerCount = 1
	performerCountCriterion := models.IntCriterionInput{
		Value:    performerCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyImagesPerformerCount(t, performerCountCriterion)

	performerCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyImagesPerformerCount(t, performerCountCriterion)

	performerCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyImagesPerformerCount(t, performerCountCriterion)

	performerCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyImagesPerformerCount(t, performerCountCriterion)
}

func verifyImagesPerformerCount(t *testing.T, performerCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := sqlite.ImageReaderWriter
		imageFilter := models.ImageFilterType{
			PerformerCount: &performerCountCriterion,
		}

		images := queryImages(ctx, t, sqb, &imageFilter, nil)
		assert.Greater(t, len(images), 0)

		for _, image := range images {
			ids, err := sqb.GetPerformerIDs(ctx, image.ID)
			if err != nil {
				return err
			}
			verifyInt(t, len(ids), performerCountCriterion)
		}

		return nil
	})
}

func TestImageQuerySorting(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sort := titleField
		direction := models.SortDirectionEnumAsc
		findFilter := models.FindFilterType{
			Sort:      &sort,
			Direction: &direction,
		}

		sqb := sqlite.ImageReaderWriter
		images, _, err := queryImagesWithCount(ctx, sqb, nil, &findFilter)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		// images should be in same order as indexes
		firstImage := images[0]
		lastImage := images[len(images)-1]

		assert.Equal(t, imageIDs[0], firstImage.ID)
		assert.Equal(t, imageIDs[len(imageIDs)-1], lastImage.ID)

		// sort in descending order
		direction = models.SortDirectionEnumDesc

		images, _, err = queryImagesWithCount(ctx, sqb, nil, &findFilter)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}
		firstImage = images[0]
		lastImage = images[len(images)-1]

		assert.Equal(t, imageIDs[len(imageIDs)-1], firstImage.ID)
		assert.Equal(t, imageIDs[0], lastImage.ID)

		return nil
	})
}

func TestImageQueryPagination(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		perPage := 1
		findFilter := models.FindFilterType{
			PerPage: &perPage,
		}

		sqb := sqlite.ImageReaderWriter
		images, _, err := queryImagesWithCount(ctx, sqb, nil, &findFilter)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		assert.Len(t, images, 1)

		firstID := images[0].ID

		page := 2
		findFilter.Page = &page
		images, _, err = queryImagesWithCount(ctx, sqb, nil, &findFilter)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		assert.Len(t, images, 1)
		secondID := images[0].ID
		assert.NotEqual(t, firstID, secondID)

		perPage = 2
		page = 1

		images, _, err = queryImagesWithCount(ctx, sqb, nil, &findFilter)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}
		assert.Len(t, images, 2)
		assert.Equal(t, firstID, images[0].ID)
		assert.Equal(t, secondID, images[1].ID)

		return nil
	})
}

// TODO Count
// TODO SizeCount
// TODO All
