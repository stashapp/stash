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
	"github.com/stretchr/testify/assert"
)

func loadImageRelationships(ctx context.Context, expected models.Image, actual *models.Image) error {
	if expected.URLs.Loaded() {
		if err := actual.LoadURLs(ctx, db.Image); err != nil {
			return err
		}
	}
	if expected.GalleryIDs.Loaded() {
		if err := actual.LoadGalleryIDs(ctx, db.Image); err != nil {
			return err
		}
	}
	if expected.TagIDs.Loaded() {
		if err := actual.LoadTagIDs(ctx, db.Image); err != nil {
			return err
		}
	}
	if expected.PerformerIDs.Loaded() {
		if err := actual.LoadPerformerIDs(ctx, db.Image); err != nil {
			return err
		}
	}
	if expected.Files.Loaded() {
		if err := actual.LoadFiles(ctx, db.Image); err != nil {
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
	if expected.PrimaryFileID == nil {
		actual.PrimaryFileID = nil
	}

	return nil
}

func Test_imageQueryBuilder_Create(t *testing.T) {
	var (
		title        = "title"
		code         = "code"
		rating       = 60
		details      = "details"
		photographer = "photographer"
		ocounter     = 5
		url          = "url"
		date, _      = models.ParseDate("2003-02-01")
		createdAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)

		imageFile = makeFileWithID(fileIdxStartImageFiles)
	)

	tests := []struct {
		name      string
		newObject models.CreateImageInput
		wantErr   bool
	}{
		{
			"full",
			models.CreateImageInput{
				Image: &models.Image{
					Title:        title,
					Code:         code,
					Rating:       &rating,
					Date:         &date,
					Details:      details,
					Photographer: photographer,
					URLs:         models.NewRelatedStrings([]string{url}),
					Organized:    true,
					OCounter:     ocounter,
					StudioID:     &studioIDs[studioIdxWithImage],
					CreatedAt:    createdAt,
					UpdatedAt:    updatedAt,
					GalleryIDs:   models.NewRelatedIDs([]int{galleryIDs[galleryIdxWithImage]}),
					TagIDs:       models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithImage]}),
					PerformerIDs: models.NewRelatedIDs([]int{performerIDs[performerIdx1WithImage], performerIDs[performerIdx1WithDupName]}),
				},
				CustomFields: testCustomFields,
			},
			false,
		},
		{
			"with file",
			models.CreateImageInput{
				Image: &models.Image{
					Title:        title,
					Code:         code,
					Rating:       &rating,
					Date:         &date,
					Details:      details,
					Photographer: photographer,
					URLs:         models.NewRelatedStrings([]string{url}),
					Organized:    true,
					OCounter:     ocounter,
					StudioID:     &studioIDs[studioIdxWithImage],
					Files: models.NewRelatedFiles([]models.File{
						imageFile.(*models.ImageFile),
					}),
					PrimaryFileID: &imageFile.Base().ID,
					Path:          imageFile.Base().Path,
					CreatedAt:     createdAt,
					UpdatedAt:     updatedAt,
					GalleryIDs:    models.NewRelatedIDs([]int{galleryIDs[galleryIdxWithImage]}),
					TagIDs:        models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithImage]}),
					PerformerIDs:  models.NewRelatedIDs([]int{performerIDs[performerIdx1WithImage], performerIDs[performerIdx1WithDupName]}),
				},
			},
			false,
		},
		{
			"invalid studio id",
			models.CreateImageInput{
				Image: &models.Image{
					StudioID: &invalidID,
				},
			},
			true,
		},
		{
			"invalid gallery id",
			models.CreateImageInput{
				Image: &models.Image{
					GalleryIDs: models.NewRelatedIDs([]int{invalidID}),
				},
			},
			true,
		},
		{
			"invalid tag id",
			models.CreateImageInput{
				Image: &models.Image{
					TagIDs: models.NewRelatedIDs([]int{invalidID}),
				},
			},
			true,
		},
		{
			"invalid performer id",
			models.CreateImageInput{
				Image: &models.Image{
					PerformerIDs: models.NewRelatedIDs([]int{invalidID}),
				},
			},
			true,
		},
	}

	qb := db.Image

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			var fileIDs []models.FileID
			if tt.newObject.Files.Loaded() {
				for _, f := range tt.newObject.Files.List() {
					fileIDs = append(fileIDs, f.Base().ID)
				}
			}
			s := *tt.newObject.Image
			if err := qb.Create(ctx, &models.CreateImageInput{
				Image:   &s,
				FileIDs: fileIDs,
			}); (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.Create() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				assert.Zero(s.ID)
				return
			}

			assert.NotZero(s.ID)

			copy := *tt.newObject.Image
			copy.ID = s.ID

			// load relationships
			if err := loadImageRelationships(ctx, copy, &s); err != nil {
				t.Errorf("loadImageRelationships() error = %v", err)
				return
			}

			assert.Equal(copy, s)

			// ensure can find the image
			found, err := qb.Find(ctx, s.ID)
			if err != nil {
				t.Errorf("imageQueryBuilder.Find() error = %v", err)
			}

			// load relationships
			if err := loadImageRelationships(ctx, copy, found); err != nil {
				t.Errorf("loadImageRelationships() error = %v", err)
				return
			}

			assert.Equal(copy, *found)
		})
	}
}

func clearImageFileIDs(image *models.Image) {
	if image.Files.Loaded() {
		for _, f := range image.Files.List() {
			f.Base().ID = 0
		}
	}
}

func makeImageFileWithID(i int) *models.ImageFile {
	ret := makeImageFile(i)
	ret.ID = imageFileIDs[i]
	return ret
}

func Test_imageQueryBuilder_Update(t *testing.T) {
	var (
		title        = "title"
		code         = "code"
		rating       = 60
		url          = "url"
		details      = "details"
		photographer = "photographer"
		date, _      = models.ParseDate("2003-02-01")
		ocounter     = 5
		createdAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
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
				Title:        title,
				Code:         code,
				Rating:       &rating,
				URLs:         models.NewRelatedStrings([]string{url}),
				Date:         &date,
				Details:      details,
				Photographer: photographer,
				Organized:    true,
				OCounter:     ocounter,
				StudioID:     &studioIDs[studioIdxWithImage],
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				GalleryIDs:   models.NewRelatedIDs([]int{galleryIDs[galleryIdxWithImage]}),
				TagIDs:       models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithImage]}),
				PerformerIDs: models.NewRelatedIDs([]int{performerIDs[performerIdx1WithImage], performerIDs[performerIdx1WithDupName]}),
			},
			false,
		},
		{
			"clear nullables",
			&models.Image{
				ID:           imageIDs[imageIdxWithGallery],
				GalleryIDs:   models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				Organized:    true,
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
			},
			false,
		},
		{
			"clear gallery ids",
			&models.Image{
				ID:           imageIDs[imageIdxWithGallery],
				GalleryIDs:   models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				Organized:    true,
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
			},
			false,
		},
		{
			"clear tag ids",
			&models.Image{
				ID:           imageIDs[imageIdxWithTag],
				GalleryIDs:   models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				Organized:    true,
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
			},
			false,
		},
		{
			"clear performer ids",
			&models.Image{
				ID:           imageIDs[imageIdxWithPerformer],
				GalleryIDs:   models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				Organized:    true,
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
			},
			false,
		},
		{
			"invalid studio id",
			&models.Image{
				ID:        imageIDs[imageIdxWithGallery],
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
				Organized:  true,
				GalleryIDs: models.NewRelatedIDs([]int{invalidID}),
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
			},
			true,
		},
		{
			"invalid tag id",
			&models.Image{
				ID:        imageIDs[imageIdxWithGallery],
				Organized: true,
				TagIDs:    models.NewRelatedIDs([]int{invalidID}),
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			true,
		},
		{
			"invalid performer id",
			&models.Image{
				ID:           imageIDs[imageIdxWithGallery],
				Organized:    true,
				PerformerIDs: models.NewRelatedIDs([]int{invalidID}),
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
			},
			true,
		},
	}

	qb := db.Image
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

			// load relationships
			if err := loadImageRelationships(ctx, copy, s); err != nil {
				t.Errorf("loadImageRelationships() error = %v", err)
				return
			}

			assert.Equal(copy, *s)
		})
	}
}

func clearImagePartial() models.ImagePartial {
	// leave mandatory fields
	return models.ImagePartial{
		Title:        models.OptionalString{Set: true, Null: true},
		Code:         models.OptionalString{Set: true, Null: true},
		Details:      models.OptionalString{Set: true, Null: true},
		Photographer: models.OptionalString{Set: true, Null: true},
		Rating:       models.OptionalInt{Set: true, Null: true},
		URLs:         &models.UpdateStrings{Mode: models.RelationshipUpdateModeSet},
		Date:         models.OptionalDate{Set: true, Null: true},
		StudioID:     models.OptionalInt{Set: true, Null: true},
		GalleryIDs:   &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
		TagIDs:       &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
		PerformerIDs: &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
	}
}

func Test_imageQueryBuilder_UpdatePartial(t *testing.T) {
	var (
		title        = "title"
		code         = "code"
		details      = "details"
		photographer = "photographer"
		rating       = 60
		url          = "url"
		date, _      = models.ParseDate("2003-02-01")
		ocounter     = 5
		createdAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
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
				Title:        models.NewOptionalString(title),
				Code:         models.NewOptionalString(code),
				Details:      models.NewOptionalString(details),
				Photographer: models.NewOptionalString(photographer),
				Rating:       models.NewOptionalInt(rating),
				URLs: &models.UpdateStrings{
					Values: []string{url},
					Mode:   models.RelationshipUpdateModeSet,
				},
				Date:      models.NewOptionalDate(date),
				Organized: models.NewOptionalBool(true),
				OCounter:  models.NewOptionalInt(ocounter),
				StudioID:  models.NewOptionalInt(studioIDs[studioIdxWithImage]),
				CreatedAt: models.NewOptionalTime(createdAt),
				UpdatedAt: models.NewOptionalTime(updatedAt),
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
				Title:        title,
				Code:         code,
				Details:      details,
				Photographer: photographer,
				Rating:       &rating,
				URLs:         models.NewRelatedStrings([]string{url}),
				Date:         &date,
				Organized:    true,
				OCounter:     ocounter,
				StudioID:     &studioIDs[studioIdxWithImage],
				Files: models.NewRelatedFiles([]models.File{
					makeImageFile(imageIdx1WithGallery),
				}),
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				GalleryIDs:   models.NewRelatedIDs([]int{galleryIDs[galleryIdxWithImage]}),
				TagIDs:       models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithImage]}),
				PerformerIDs: models.NewRelatedIDs([]int{performerIDs[performerIdx1WithImage], performerIDs[performerIdx1WithDupName]}),
			},
			false,
		},
		{
			"clear all",
			imageIDs[imageIdx1WithGallery],
			clearImagePartial(),
			models.Image{
				ID:       imageIDs[imageIdx1WithGallery],
				OCounter: getOCounter(imageIdx1WithGallery),
				Files: models.NewRelatedFiles([]models.File{
					makeImageFile(imageIdx1WithGallery),
				}),
				GalleryIDs:   models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
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
		qb := db.Image

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

			// load relationships
			if err := loadImageRelationships(ctx, tt.want, got); err != nil {
				t.Errorf("loadImageRelationships() error = %v", err)
				return
			}
			clearImageFileIDs(got)

			assert.Equal(tt.want, *got)

			s, err := qb.Find(ctx, tt.id)
			if err != nil {
				t.Errorf("imageQueryBuilder.Find() error = %v", err)
			}

			// load relationships
			if err := loadImageRelationships(ctx, tt.want, s); err != nil {
				t.Errorf("loadImageRelationships() error = %v", err)
				return
			}
			clearImageFileIDs(s)
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
				GalleryIDs: models.NewRelatedIDs(append(indexesToIDs(galleryIDs, imageGalleries[imageIdxWithGallery]),
					galleryIDs[galleryIdx1WithImage],
					galleryIDs[galleryIdx1WithPerformer],
				)),
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
				TagIDs: models.NewRelatedIDs(append(
					[]int{
						tagIDs[tagIdx1WithGallery],
						tagIDs[tagIdx1WithDupName],
					},
					indexesToIDs(tagIDs, imageTags[imageIdxWithTwoTags])...,
				)),
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
				PerformerIDs: models.NewRelatedIDs(append(indexesToIDs(performerIDs, imagePerformers[imageIdxWithTwoPerformers]),
					performerIDs[performerIdx1WithDupName],
					performerIDs[performerIdx1WithGallery],
				)),
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
				GalleryIDs: models.NewRelatedIDs(append(indexesToIDs(galleryIDs, imageGalleries[imageIdxWithGallery]),
					galleryIDs[galleryIdx1WithPerformer],
				)),
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
				TagIDs: models.NewRelatedIDs(append(
					[]int{tagIDs[tagIdx1WithGallery]},
					indexesToIDs(tagIDs, imageTags[imageIdxWithTwoTags])...,
				)),
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
				PerformerIDs: models.NewRelatedIDs(append(indexesToIDs(performerIDs, imagePerformers[imageIdxWithTwoPerformers]),
					performerIDs[performerIdx1WithGallery],
				)),
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
			models.Image{
				GalleryIDs: models.NewRelatedIDs([]int{}),
			},
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
				TagIDs: models.NewRelatedIDs([]int{tagIDs[tagIdx2WithImage]}),
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
				PerformerIDs: models.NewRelatedIDs([]int{performerIDs[performerIdx2WithImage]}),
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
				GalleryIDs: models.NewRelatedIDs([]int{galleryIDs[galleryIdxWithImage]}),
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
				TagIDs: models.NewRelatedIDs(indexesToIDs(tagIDs, imageTags[imageIdxWithTwoTags])),
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
				PerformerIDs: models.NewRelatedIDs(indexesToIDs(performerIDs, imagePerformers[imageIdxWithTwoPerformers])),
			},
			false,
		},
	}

	for _, tt := range tests {
		qb := db.Image

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

			// load relationships
			if err := loadImageRelationships(ctx, tt.want, got); err != nil {
				t.Errorf("loadImageRelationships() error = %v", err)
				return
			}
			if err := loadImageRelationships(ctx, tt.want, s); err != nil {
				t.Errorf("loadImageRelationships() error = %v", err)
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
		})
	}
}

func Test_ImageStore_UpdatePartialCustomFields(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		partial  models.ImagePartial
		expected map[string]interface{} // nil to use the partial
	}{
		{
			"set custom fields",
			imageIDs[imageIdx1WithGallery],
			models.ImagePartial{
				CustomFields: models.CustomFieldsInput{
					Full: testCustomFields,
				},
			},
			nil,
		},
		{
			"clear custom fields",
			imageIDs[imageIdx1WithGallery],
			models.ImagePartial{
				CustomFields: models.CustomFieldsInput{
					Full: map[string]interface{}{},
				},
			},
			nil,
		},
		{
			"partial custom fields",
			imageIDs[imageIdxWithStudio],
			models.ImagePartial{
				CustomFields: models.CustomFieldsInput{
					Partial: map[string]interface{}{
						"string":    "bbb",
						"new_field": "new",
					},
				},
			},
			map[string]interface{}{
				"int":       int64(2),
				"real":      1.2,
				"string":    "bbb",
				"new_field": "new",
			},
		},
	}
	for _, tt := range tests {
		qb := db.Image

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			_, err := qb.UpdatePartial(ctx, tt.id, tt.partial)
			if err != nil {
				t.Errorf("ImageStore.UpdatePartial() error = %v", err)
				return
			}

			// ensure custom fields are correct
			cf, err := qb.GetCustomFields(ctx, tt.id)
			if err != nil {
				t.Errorf("ImageStore.GetCustomFields() error = %v", err)
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

	qb := db.Image

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

	qb := db.Image

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

	qb := db.Image

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

	qb := db.Image

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			if err := qb.Destroy(ctx, tt.id); (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.Destroy() error = %v, wantErr %v", err, tt.wantErr)
			}

			// ensure cannot be found
			i, err := qb.Find(ctx, tt.id)

			assert.Nil(err)
			assert.Nil(i)
		})
	}
}

func makeImageWithID(index int) *models.Image {
	const fromDB = true
	ret := makeImage(index)
	ret.ID = imageIDs[index]

	ret.Files = models.NewRelatedFiles([]models.File{makeImageFile(index)})

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
			false,
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

	qb := db.Image

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.Find(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				// load relationships
				if err := loadImageRelationships(ctx, *tt.want, got); err != nil {
					t.Errorf("loadImageRelationships() error = %v", err)
					return
				}
				clearImageFileIDs(got)
			}
			assert.Equal(tt.want, got)
		})
	}
}

func postFindImages(ctx context.Context, want []*models.Image, got []*models.Image) error {
	for i, s := range got {
		// load relationships
		if i < len(want) {
			if err := loadImageRelationships(ctx, *want[i], s); err != nil {
				return err
			}
		}
		clearImageFileIDs(s)
	}

	return nil
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

	qb := db.Image

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.FindMany(ctx, tt.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.FindMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := postFindImages(ctx, tt.want, got); err != nil {
				t.Errorf("loadImageRelationships() error = %v", err)
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
		want     []*models.Image
		wantErr  bool
	}{
		{
			"valid",
			getChecksum(imageIdxWithGallery),
			[]*models.Image{makeImageWithID(imageIdxWithGallery)},
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
			[]*models.Image{makeImageWithID(imageIdxWithTwoPerformers)},
			false,
		},
		{
			"with tags",
			getChecksum(imageIdxWithTwoTags),
			[]*models.Image{makeImageWithID(imageIdxWithTwoTags)},
			false,
		},
	}

	qb := db.Image

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByChecksum(ctx, tt.checksum)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.FindByChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := postFindImages(ctx, tt.want, got); err != nil {
				t.Errorf("loadImageRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func Test_imageQueryBuilder_FindByFingerprints(t *testing.T) {
	getChecksum := func(index int) string {
		return getImageStringValue(index, checksumField)
	}

	tests := []struct {
		name         string
		fingerprints []models.Fingerprint
		want         []*models.Image
		wantErr      bool
	}{
		{
			"valid",
			[]models.Fingerprint{
				{
					Type:        models.FingerprintTypeMD5,
					Fingerprint: getChecksum(imageIdxWithGallery),
				},
			},
			[]*models.Image{makeImageWithID(imageIdxWithGallery)},
			false,
		},
		{
			"invalid",
			[]models.Fingerprint{
				{
					Type:        models.FingerprintTypeMD5,
					Fingerprint: "invalid checksum",
				},
			},
			nil,
			false,
		},
		{
			"with performers",
			[]models.Fingerprint{
				{
					Type:        models.FingerprintTypeMD5,
					Fingerprint: getChecksum(imageIdxWithTwoPerformers),
				},
			},
			[]*models.Image{makeImageWithID(imageIdxWithTwoPerformers)},
			false,
		},
		{
			"with tags",
			[]models.Fingerprint{
				{
					Type:        models.FingerprintTypeMD5,
					Fingerprint: getChecksum(imageIdxWithTwoTags),
				},
			},
			[]*models.Image{makeImageWithID(imageIdxWithTwoTags)},
			false,
		},
	}

	qb := db.Image

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByFingerprints(ctx, tt.fingerprints)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.FindByChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := postFindImages(ctx, tt.want, got); err != nil {
				t.Errorf("loadImageRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, got)
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

	qb := db.Image

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByGalleryID(ctx, tt.galleryID)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageQueryBuilder.FindByGalleryID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := postFindImages(ctx, tt.want, got); err != nil {
				t.Errorf("loadImageRelationships() error = %v", err)
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

	qb := db.Image

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

func imagesToIDs(i []*models.Image) []int {
	var ret []int
	for _, ii := range i {
		ret = append(ret, ii.ID)
	}

	return ret
}

func Test_imageStore_FindByFileID(t *testing.T) {
	tests := []struct {
		name    string
		fileID  models.FileID
		include []int
		exclude []int
	}{
		{
			"valid",
			imageFileIDs[imageIdxWithGallery],
			[]int{imageIdxWithGallery},
			nil,
		},
		{
			"invalid",
			invalidFileID,
			nil,
			[]int{imageIdxWithGallery},
		},
	}

	qb := db.Image

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByFileID(ctx, tt.fileID)
			if err != nil {
				t.Errorf("ImageStore.FindByFileID() error = %v", err)
				return
			}
			for _, f := range got {
				clearImageFileIDs(f)
			}

			ids := imagesToIDs(got)
			include := indexesToIDs(imageIDs, tt.include)
			exclude := indexesToIDs(imageIDs, tt.exclude)

			for _, i := range include {
				assert.Contains(ids, i)
			}
			for _, e := range exclude {
				assert.NotContains(ids, e)
			}
		})
	}
}

func Test_imageStore_FindByFolderID(t *testing.T) {
	tests := []struct {
		name     string
		folderID models.FolderID
		include  []int
		exclude  []int
	}{
		{
			"valid",
			folderIDs[folderIdxWithImageFiles],
			[]int{imageIdxWithGallery},
			nil,
		},
		{
			"invalid",
			invalidFolderID,
			nil,
			[]int{imageIdxWithGallery},
		},
		{
			"parent folder",
			folderIDs[folderIdxForObjectFiles],
			nil,
			[]int{imageIdxWithGallery},
		},
	}

	qb := db.Image

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByFolderID(ctx, tt.folderID)
			if err != nil {
				t.Errorf("ImageStore.FindByFolderID() error = %v", err)
				return
			}
			for _, f := range got {
				clearImageFileIDs(f)
			}

			ids := imagesToIDs(got)
			include := indexesToIDs(imageIDs, tt.include)
			exclude := indexesToIDs(imageIDs, tt.exclude)

			for _, i := range include {
				assert.Contains(ids, i)
			}
			for _, e := range exclude {
				assert.NotContains(ids, e)
			}
		})
	}
}

func Test_imageStore_FindByZipFileID(t *testing.T) {
	tests := []struct {
		name      string
		zipFileID models.FileID
		include   []int
		exclude   []int
	}{
		{
			"valid",
			fileIDs[fileIdxZip],
			[]int{imageIdxInZip},
			nil,
		},
		{
			"invalid",
			invalidFileID,
			nil,
			[]int{imageIdxInZip},
		},
	}

	qb := db.Image

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByZipFileID(ctx, tt.zipFileID)
			if err != nil {
				t.Errorf("ImageStore.FindByZipFileID() error = %v", err)
				return
			}
			for _, f := range got {
				clearImageFileIDs(f)
			}

			ids := imagesToIDs(got)
			include := indexesToIDs(imageIDs, tt.include)
			exclude := indexesToIDs(imageIDs, tt.exclude)

			for _, i := range include {
				assert.Contains(ids, i)
			}
			for _, e := range exclude {
				assert.NotContains(ids, e)
			}
		})
	}
}

func TestImageQueryQ(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		const imageIdx = 2

		q := getImageStringValue(imageIdx, titleField)

		sqb := db.Image

		imageQueryQ(ctx, t, sqb, q, imageIdx)

		return nil
	})
}

func TestImageQueryQ_Details(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		const imageIdx = 3

		q := getImageStringValue(imageIdx, "Details")

		sqb := db.Image

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

func verifyImageQuery(t *testing.T, filter models.ImageFilterType, verifyFn func(ctx context.Context, s *models.Image)) {
	t.Helper()
	withTxn(func(ctx context.Context) error {
		t.Helper()
		sqb := db.Image

		images := queryImages(ctx, t, sqb, &filter, nil)

		// assume it should find at least one
		assert.Greater(t, len(images), 0)

		for _, image := range images {
			verifyFn(ctx, image)
		}

		return nil
	})
}

func TestImageQueryURL(t *testing.T) {
	const imageIdx = 1
	imageURL := getImageStringValue(imageIdx, urlField)
	urlCriterion := models.StringCriterionInput{
		Value:    imageURL,
		Modifier: models.CriterionModifierEquals,
	}
	filter := models.ImageFilterType{
		URL: &urlCriterion,
	}

	verifyFn := func(ctx context.Context, o *models.Image) {
		t.Helper()

		if err := o.LoadURLs(ctx, db.Image); err != nil {
			t.Errorf("Error loading scene URLs: %v", err)
		}

		urls := o.URLs.List()
		var url string
		if len(urls) > 0 {
			url = urls[0]
		}

		verifyString(t, url, urlCriterion)
	}

	verifyImageQuery(t, filter, verifyFn)
	urlCriterion.Modifier = models.CriterionModifierNotEquals
	verifyImageQuery(t, filter, verifyFn)
	urlCriterion.Modifier = models.CriterionModifierMatchesRegex
	urlCriterion.Value = "image_.*1_URL"
	verifyImageQuery(t, filter, verifyFn)
	urlCriterion.Modifier = models.CriterionModifierNotMatchesRegex
	verifyImageQuery(t, filter, verifyFn)
	urlCriterion.Modifier = models.CriterionModifierIsNull
	urlCriterion.Value = ""
	verifyImageQuery(t, filter, verifyFn)
	urlCriterion.Modifier = models.CriterionModifierNotNull
	verifyImageQuery(t, filter, verifyFn)
}

func TestImageQueryPath(t *testing.T) {
	const imageIdx = 1
	imagePath := getFilePath(folderIdxWithImageFiles, getImageBasename(imageIdx))

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
		sqb := db.Image
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

	image1Path := getFilePath(folderIdxWithImageFiles, getImageBasename(image1Idx))
	image2Path := getFilePath(folderIdxWithImageFiles, getImageBasename(image2Idx))

	imageFilter := models.ImageFilterType{
		Path: &models.StringCriterionInput{
			Value:    image1Path,
			Modifier: models.CriterionModifierEquals,
		},
		OperatorFilter: models.OperatorFilter[models.ImageFilterType]{
			Or: &models.ImageFilterType{
				Path: &models.StringCriterionInput{
					Value:    image2Path,
					Modifier: models.CriterionModifierEquals,
				},
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Image

		images := queryImages(ctx, t, sqb, &imageFilter, nil)

		if !assert.Len(t, images, 2) {
			return nil
		}

		assert.Equal(t, image1Path, images[0].Path)
		assert.Equal(t, image2Path, images[1].Path)

		return nil
	})
}

func TestImageQueryPathAndRating(t *testing.T) {
	const imageIdx = 1
	imagePath := getFilePath(folderIdxWithImageFiles, getImageBasename(imageIdx))
	imageRating := getRating(imageIdx)

	imageFilter := models.ImageFilterType{
		Path: &models.StringCriterionInput{
			Value:    imagePath,
			Modifier: models.CriterionModifierEquals,
		},
		OperatorFilter: models.OperatorFilter[models.ImageFilterType]{
			And: &models.ImageFilterType{
				Rating100: &models.IntCriterionInput{
					Value:    int(imageRating.Int64),
					Modifier: models.CriterionModifierEquals,
				},
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Image

		images := queryImages(ctx, t, sqb, &imageFilter, nil)

		if !assert.Len(t, images, 1) {
			return nil
		}

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
		OperatorFilter: models.OperatorFilter[models.ImageFilterType]{
			Not: &models.ImageFilterType{
				Rating100: &ratingCriterion,
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Image

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
		OperatorFilter: models.OperatorFilter[models.ImageFilterType]{
			And: &subFilter,
			Or:  &subFilter,
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Image

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

func TestImageQueryRating100(t *testing.T) {
	const rating = 60
	ratingCriterion := models.IntCriterionInput{
		Value:    rating,
		Modifier: models.CriterionModifierEquals,
	}

	verifyImagesRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotEquals
	verifyImagesRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyImagesRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierLessThan
	verifyImagesRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierIsNull
	verifyImagesRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotNull
	verifyImagesRating100(t, ratingCriterion)
}

func verifyImagesRating100(t *testing.T, ratingCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Image
		imageFilter := models.ImageFilterType{
			Rating100: &ratingCriterion,
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
		sqb := db.Image
		imageFilter := models.ImageFilterType{
			OCounter: &oCounterCriterion,
		}

		images, _, err := queryImagesWithCount(ctx, sqb, &imageFilter, nil)
		if err != nil {
			t.Errorf("Error querying image: %s", err.Error())
		}

		for _, image := range images {
			verifyIntPtr(t, image.OCounter, oCounterCriterion)
		}

		return nil
	})
}
