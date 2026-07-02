//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"math"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stretchr/testify/assert"
)

func loadAudioRelationships(ctx context.Context, expected models.Audio, actual *models.Audio) error {
	if expected.URLs.Loaded() {
		if err := actual.LoadURLs(ctx, db.Audio); err != nil {
			return err
		}
	}

	if expected.TagIDs.Loaded() {
		if err := actual.LoadTagIDs(ctx, db.Audio); err != nil {
			return err
		}
	}
	if expected.PerformerIDs.Loaded() {
		if err := actual.LoadPerformerIDs(ctx, db.Audio); err != nil {
			return err
		}
	}
	if expected.Groups.Loaded() {
		if err := actual.LoadGroups(ctx, db.Audio); err != nil {
			return err
		}
	}
	if expected.Files.Loaded() {
		if err := actual.LoadFiles(ctx, db.Audio); err != nil {
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

func Test_audioQueryBuilder_Create(t *testing.T) {
	var (
		title        = "title"
		code         = "1337"
		details      = "details"
		url          = "url"
		rating       = 60
		resumeTime   = 10.0
		playDuration = 34.0
		createdAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		audioIndex   = 123
		audioIndex2  = 234

		date, _ = models.ParseDate("2003-02-01")

		audioFile = makeFileWithID(fileIdxStartAudioFiles)
	)

	tests := []struct {
		name      string
		newObject models.Audio
		wantErr   bool
	}{
		{
			"full",
			models.Audio{
				Title:        title,
				Code:         code,
				Details:      details,
				URLs:         models.NewRelatedStrings([]string{url}),
				Date:         &date,
				Rating:       &rating,
				Organized:    true,
				StudioID:     &studioIDs[studioIdxWithAudio],
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				TagIDs:       models.NewRelatedIDs([]int{tagIDs[tagIdx1WithAudio], tagIDs[tagIdx1WithNothing]}),
				PerformerIDs: models.NewRelatedIDs([]int{performerIDs[performerIdx1WithAudio], performerIDs[performerIdx1WithDupName]}),
				Groups: models.NewRelatedGroupsAudio([]models.GroupsAudios{
					{
						GroupID:    groupIDs[groupIdxWithAudio],
						AudioIndex: &audioIndex,
					},
					{
						GroupID:    groupIDs[groupIdxWithStudio],
						AudioIndex: &audioIndex2,
					},
				}),
				ResumeTime:   float64(resumeTime),
				PlayDuration: playDuration,
			},
			false,
		},
		{
			"with file",
			models.Audio{
				Title:     title,
				Code:      code,
				Details:   details,
				URLs:      models.NewRelatedStrings([]string{url}),
				Date:      &date,
				Rating:    &rating,
				Organized: true,
				StudioID:  &studioIDs[studioIdxWithAudio],
				Files: models.NewRelatedAudioFiles([]*models.AudioFile{
					audioFile.(*models.AudioFile),
				}),
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				TagIDs:       models.NewRelatedIDs([]int{tagIDs[tagIdx1WithAudio], tagIDs[tagIdx1WithNothing]}),
				PerformerIDs: models.NewRelatedIDs([]int{performerIDs[performerIdx1WithAudio], performerIDs[performerIdx1WithDupName]}),
				Groups: models.NewRelatedGroupsAudio([]models.GroupsAudios{
					{
						GroupID:    groupIDs[groupIdxWithAudio],
						AudioIndex: &audioIndex,
					},
					{
						GroupID:    groupIDs[groupIdxWithStudio],
						AudioIndex: &audioIndex2,
					},
				}),
				ResumeTime:   resumeTime,
				PlayDuration: playDuration,
			},
			false,
		},
		{
			"invalid studio id",
			models.Audio{
				StudioID: &invalidID,
			},
			true,
		},
		{
			"invalid tag id",
			models.Audio{
				TagIDs: models.NewRelatedIDs([]int{invalidID}),
			},
			true,
		},
		{
			"invalid performer id",
			models.Audio{
				PerformerIDs: models.NewRelatedIDs([]int{invalidID}),
			},
			true,
		},
		{
			"invalid group id",
			models.Audio{
				Groups: models.NewRelatedGroupsAudio([]models.GroupsAudios{
					{
						GroupID:    invalidID,
						AudioIndex: &audioIndex,
					},
				}),
			},
			true,
		},
	}

	qb := db.Audio

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
				t.Errorf("audioQueryBuilder.Create() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				assert.Zero(s.ID)
				return
			}

			assert.NotZero(s.ID)

			copy := tt.newObject
			copy.ID = s.ID

			// load relationships
			if err := loadAudioRelationships(ctx, copy, &s); err != nil {
				t.Errorf("loadAudioRelationships() error = %v", err)
				return
			}

			assert.Equal(copy, s)

			// ensure can find the audio
			found, err := qb.Find(ctx, s.ID)
			if err != nil {
				t.Errorf("audioQueryBuilder.Find() error = %v", err)
			}

			if !assert.NotNil(found) {
				return
			}

			// load relationships
			if err := loadAudioRelationships(ctx, copy, found); err != nil {
				t.Errorf("loadAudioRelationships() error = %v", err)
				return
			}
			assert.Equal(copy, *found)

			return
		})
	}
}

func clearAudioFileIDs(audio *models.Audio) {
	if audio.Files.Loaded() {
		for _, f := range audio.Files.List() {
			f.Base().ID = 0
		}
	}
}

func makeAudioFileWithID(i int) *models.AudioFile {
	ret := makeAudioFile(i)
	ret.ID = audioFileIDs[i]
	return ret
}

func Test_audioQueryBuilder_Update(t *testing.T) {
	var (
		title        = "title"
		code         = "1337"
		details      = "details"
		url          = "url"
		rating       = 60
		resumeTime   = 10.0
		playDuration = 34.0
		createdAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		audioIndex   = 123
		audioIndex2  = 234

		date, _ = models.ParseDate("2003-02-01")
	)

	tests := []struct {
		name          string
		updatedObject *models.Audio
		wantErr       bool
	}{
		{
			"full",
			&models.Audio{
				ID:           audioIDs[audioIdxWithGallery],
				Title:        title,
				Code:         code,
				Details:      details,
				URLs:         models.NewRelatedStrings([]string{url}),
				Date:         &date,
				Rating:       &rating,
				Organized:    true,
				StudioID:     &studioIDs[studioIdxWithAudio],
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				TagIDs:       models.NewRelatedIDs([]int{tagIDs[tagIdx1WithAudio], tagIDs[tagIdx1WithNothing]}),
				PerformerIDs: models.NewRelatedIDs([]int{performerIDs[performerIdx1WithAudio], performerIDs[performerIdx1WithDupName]}),
				Groups: models.NewRelatedGroupsAudio([]models.GroupsAudios{
					{
						GroupID:    groupIDs[groupIdxWithAudio],
						AudioIndex: &audioIndex,
					},
					{
						GroupID:    groupIDs[groupIdxWithStudio],
						AudioIndex: &audioIndex2,
					},
				}),
				ResumeTime:   resumeTime,
				PlayDuration: playDuration,
			},
			false,
		},
		{
			"clear nullables",
			&models.Audio{
				ID:           audioIDs[audioIdxWithSpacedName],
				TagIDs:       models.NewRelatedIDs([]int{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				Groups:       models.NewRelatedGroupsAudio([]models.GroupsAudios{}),
			},
			false,
		},
		{
			"clear tag ids",
			&models.Audio{
				ID:     audioIDs[audioIdxWithTag],
				TagIDs: models.NewRelatedIDs([]int{}),
			},
			false,
		},
		{
			"clear performer ids",
			&models.Audio{
				ID:           audioIDs[audioIdxWithPerformer],
				PerformerIDs: models.NewRelatedIDs([]int{}),
			},
			false,
		},
		{
			"clear groups",
			&models.Audio{
				ID:     audioIDs[audioIdxWithGroup],
				Groups: models.NewRelatedGroupsAudio([]models.GroupsAudios{}),
			},
			false,
		},
		{
			"invalid studio id",
			&models.Audio{
				ID:       audioIDs[audioIdxWithGallery],
				StudioID: &invalidID,
			},
			true,
		},
		{
			"invalid tag id",
			&models.Audio{
				ID:     audioIDs[audioIdxWithGallery],
				TagIDs: models.NewRelatedIDs([]int{invalidID}),
			},
			true,
		},
		{
			"invalid performer id",
			&models.Audio{
				ID:           audioIDs[audioIdxWithGallery],
				PerformerIDs: models.NewRelatedIDs([]int{invalidID}),
			},
			true,
		},
		{
			"invalid group id",
			&models.Audio{
				ID: audioIDs[audioIdxWithSpacedName],
				Groups: models.NewRelatedGroupsAudio([]models.GroupsAudios{
					{
						GroupID:    invalidID,
						AudioIndex: &audioIndex,
					},
				}),
			},
			true,
		},
	}

	qb := db.Audio
	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			copy := *tt.updatedObject

			if err := qb.Update(ctx, tt.updatedObject); (err != nil) != tt.wantErr {
				t.Errorf("audioQueryBuilder.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			s, err := qb.Find(ctx, tt.updatedObject.ID)
			if err != nil {
				t.Errorf("audioQueryBuilder.Find() error = %v", err)
			}

			// load relationships
			if err := loadAudioRelationships(ctx, copy, s); err != nil {
				t.Errorf("loadAudioRelationships() error = %v", err)
				return
			}

			assert.Equal(copy, *s)
		})
	}
}

func clearAudioPartial() models.AudioPartial {
	// leave mandatory fields
	return models.AudioPartial{
		Title:        models.OptionalString{Set: true, Null: true},
		Code:         models.OptionalString{Set: true, Null: true},
		Details:      models.OptionalString{Set: true, Null: true},
		URLs:         &models.UpdateStrings{Mode: models.RelationshipUpdateModeSet},
		Date:         models.OptionalDate{Set: true, Null: true},
		Rating:       models.OptionalInt{Set: true, Null: true},
		StudioID:     models.OptionalInt{Set: true, Null: true},
		TagIDs:       &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
		PerformerIDs: &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
	}
}

func Test_audioQueryBuilder_UpdatePartial(t *testing.T) {
	var (
		title        = "title"
		code         = "1337"
		details      = "details"
		url          = "url"
		rating       = 60
		resumeTime   = 10.0
		playDuration = 34.0
		createdAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt    = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		audioIndex   = 123
		audioIndex2  = 234

		date, _ = models.ParseDate("2003-02-01")
	)

	tests := []struct {
		name    string
		id      int
		partial models.AudioPartial
		want    models.Audio
		wantErr bool
	}{
		{
			"full",
			audioIDs[audioIdxWithSpacedName],
			models.AudioPartial{
				Title:   models.NewOptionalString(title),
				Code:    models.NewOptionalString(code),
				Details: models.NewOptionalString(details),
				URLs: &models.UpdateStrings{
					Values: []string{url},
					Mode:   models.RelationshipUpdateModeSet,
				},
				Date:      models.NewOptionalDate(date),
				Rating:    models.NewOptionalInt(rating),
				Organized: models.NewOptionalBool(true),
				StudioID:  models.NewOptionalInt(studioIDs[studioIdxWithAudio]),
				CreatedAt: models.NewOptionalTime(createdAt),
				UpdatedAt: models.NewOptionalTime(updatedAt),
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithAudio], tagIDs[tagIdx1WithNothing]},
					Mode: models.RelationshipUpdateModeSet,
				},
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithAudio], performerIDs[performerIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeSet,
				},
				GroupIDs: &models.UpdateGroupIDsAudio{
					Groups: []models.GroupsAudios{
						{
							GroupID:    groupIDs[groupIdxWithAudio],
							AudioIndex: &audioIndex,
						},
						{
							GroupID:    groupIDs[groupIdxWithStudio],
							AudioIndex: &audioIndex2,
						},
					},
					Mode: models.RelationshipUpdateModeSet,
				},
				ResumeTime:   models.NewOptionalFloat64(resumeTime),
				PlayDuration: models.NewOptionalFloat64(playDuration),
			},
			models.Audio{
				ID: audioIDs[audioIdxWithSpacedName],
				Files: models.NewRelatedAudioFiles([]*models.AudioFile{
					makeAudioFile(audioIdxWithSpacedName),
				}),
				Title:        title,
				Code:         code,
				Details:      details,
				URLs:         models.NewRelatedStrings([]string{url}),
				Date:         &date,
				Rating:       &rating,
				Organized:    true,
				StudioID:     &studioIDs[studioIdxWithAudio],
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
				TagIDs:       models.NewRelatedIDs([]int{tagIDs[tagIdx1WithAudio], tagIDs[tagIdx1WithNothing]}),
				PerformerIDs: models.NewRelatedIDs([]int{performerIDs[performerIdx1WithAudio], performerIDs[performerIdx1WithDupName]}),
				Groups: models.NewRelatedGroupsAudio([]models.GroupsAudios{
					{
						GroupID:    groupIDs[groupIdxWithAudio],
						AudioIndex: &audioIndex,
					},
					{
						GroupID:    groupIDs[groupIdxWithStudio],
						AudioIndex: &audioIndex2,
					},
				}),
				ResumeTime:   resumeTime,
				PlayDuration: playDuration,
			},
			false,
		},
		{
			"clear all",
			audioIDs[audioIdxWithSpacedName],
			clearAudioPartial(),
			models.Audio{
				ID: audioIDs[audioIdxWithSpacedName],
				Files: models.NewRelatedAudioFiles([]*models.AudioFile{
					makeAudioFile(audioIdxWithSpacedName),
				}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				Groups:       models.NewRelatedGroupsAudio([]models.GroupsAudios{}),
				PlayDuration: getPlayDuration(audioIdxWithSpacedName),
				ResumeTime:   getResumeTime(audioIdxWithSpacedName),
			},
			false,
		},
		{
			"invalid id",
			invalidID,
			models.AudioPartial{},
			models.Audio{},
			true,
		},
	}
	for _, tt := range tests {
		qb := db.Audio

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			got, err := qb.UpdatePartial(ctx, tt.id, tt.partial)
			if (err != nil) != tt.wantErr {
				t.Errorf("audioQueryBuilder.UpdatePartial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// load relationships
			if err := loadAudioRelationships(ctx, tt.want, got); err != nil {
				t.Errorf("loadAudioRelationships() error = %v", err)
				return
			}

			// ignore file ids
			clearAudioFileIDs(got)

			assert.Equal(tt.want, *got)

			s, err := qb.Find(ctx, tt.id)
			if err != nil {
				t.Errorf("audioQueryBuilder.Find() error = %v", err)
			}

			// load relationships
			if err := loadAudioRelationships(ctx, tt.want, s); err != nil {
				t.Errorf("loadAudioRelationships() error = %v", err)
				return
			}
			// ignore file ids
			clearAudioFileIDs(s)

			assert.Equal(tt.want, *s)
		})
	}
}

func Test_audioQueryBuilder_UpdatePartialRelationships(t *testing.T) {
	var (
		audioIndex  = 123
		audioIndex2 = 234

		groupAudios = []models.GroupsAudios{
			{
				GroupID:    groupIDs[groupIdxWithDupName],
				AudioIndex: &audioIndex,
			},
			{
				GroupID:    groupIDs[groupIdxWithStudio],
				AudioIndex: &audioIndex2,
			},
		}
	)

	tests := []struct {
		name    string
		id      int
		partial models.AudioPartial
		want    models.Audio
		wantErr bool
	}{
		{
			"add tags",
			audioIDs[audioIdxWithTwoTags],
			models.AudioPartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithGallery], tagIDs[tagIdx1WithNothing]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Audio{
				TagIDs: models.NewRelatedIDs(append(
					[]int{
						tagIDs[tagIdx1WithGallery],
						tagIDs[tagIdx1WithNothing],
					},
					indexesToIDs(tagIDs, audioTags[audioIdxWithTwoTags])...,
				)),
			},
			false,
		},
		{
			"add identical tags",
			audioIDs[audioIdxWithTwoTags],
			models.AudioPartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithNothing], tagIDs[tagIdx1WithNothing]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Audio{
				TagIDs: models.NewRelatedIDs(append(
					[]int{
						tagIDs[tagIdx1WithNothing],
					},
					indexesToIDs(tagIDs, audioTags[audioIdxWithTwoTags])...,
				)),
			},
			false,
		},
		{
			"add performers",
			audioIDs[audioIdxWithTwoPerformers],
			models.AudioPartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithDupName], performerIDs[performerIdx1WithGallery]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Audio{
				PerformerIDs: models.NewRelatedIDs(append(indexesToIDs(performerIDs, audioPerformers[audioIdxWithTwoPerformers]),
					performerIDs[performerIdx1WithDupName],
					performerIDs[performerIdx1WithGallery],
				)),
			},
			false,
		},
		{
			"add identical performers",
			audioIDs[audioIdxWithTwoPerformers],
			models.AudioPartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithDupName], performerIDs[performerIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Audio{
				PerformerIDs: models.NewRelatedIDs(append(indexesToIDs(performerIDs, audioPerformers[audioIdxWithTwoPerformers]),
					performerIDs[performerIdx1WithDupName],
				)),
			},
			false,
		},
		{
			"add groups",
			audioIDs[audioIdxWithGroup],
			models.AudioPartial{
				GroupIDs: &models.UpdateGroupIDsAudio{
					Groups: groupAudios,
					Mode:   models.RelationshipUpdateModeAdd,
				},
			},
			models.Audio{
				Groups: models.NewRelatedGroupsAudio(append([]models.GroupsAudios{
					{
						GroupID: indexesToIDs(groupIDs, audioGroups[audioIdxWithGroup])[0],
					},
				}, groupAudios...)),
			},
			false,
		},
		{
			"add groups to empty",
			audioIDs[audioIdx1WithPerformer],
			models.AudioPartial{
				GroupIDs: &models.UpdateGroupIDsAudio{
					Groups: groupAudios,
					Mode:   models.RelationshipUpdateModeAdd,
				},
			},
			models.Audio{
				Groups: models.NewRelatedGroupsAudio([]models.GroupsAudios{
					{
						GroupID:    groupIDs[groupIdxWithDupName],
						AudioIndex: &audioIndex,
					},
					{
						GroupID:    groupIDs[groupIdxWithStudio],
						AudioIndex: &audioIndex2,
					},
				}),
			},
			false,
		},
		{
			"add duplicate tags",
			audioIDs[audioIdxWithTwoTags],
			models.AudioPartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithAudio], tagIDs[tagIdx1WithGallery]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Audio{
				TagIDs: models.NewRelatedIDs(append(
					[]int{tagIDs[tagIdx1WithGallery]},
					indexesToIDs(tagIDs, audioTags[audioIdxWithTwoTags])...,
				)),
			},
			false,
		},
		{
			"add duplicate performers",
			audioIDs[audioIdxWithTwoPerformers],
			models.AudioPartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithAudio], performerIDs[performerIdx1WithGallery]},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Audio{
				PerformerIDs: models.NewRelatedIDs(append(indexesToIDs(performerIDs, audioPerformers[audioIdxWithTwoPerformers]),
					performerIDs[performerIdx1WithGallery],
				)),
			},
			false,
		},
		{
			"add duplicate groups",
			audioIDs[audioIdxWithGroup],
			models.AudioPartial{
				GroupIDs: &models.UpdateGroupIDsAudio{
					Groups: append([]models.GroupsAudios{
						{
							GroupID:    groupIDs[groupIdxWithAudio],
							AudioIndex: &audioIndex,
						},
					},
						groupAudios...,
					),
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Audio{
				Groups: models.NewRelatedGroupsAudio(append([]models.GroupsAudios{
					{
						GroupID: indexesToIDs(groupIDs, audioGroups[audioIdxWithGroup])[0],
					},
				}, groupAudios...)),
			},
			false,
		},
		{
			"add invalid tags",
			audioIDs[audioIdxWithTwoTags],
			models.AudioPartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{invalidID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Audio{},
			true,
		},
		{
			"add invalid performers",
			audioIDs[audioIdxWithTwoPerformers],
			models.AudioPartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{invalidID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Audio{},
			true,
		},
		{
			"add invalid groups",
			audioIDs[audioIdxWithGroup],
			models.AudioPartial{
				GroupIDs: &models.UpdateGroupIDsAudio{
					Groups: []models.GroupsAudios{
						{
							GroupID: invalidID,
						},
					},
					Mode: models.RelationshipUpdateModeAdd,
				},
			},
			models.Audio{},
			true,
		},
		{
			"remove tags",
			audioIDs[audioIdxWithTwoTags],
			models.AudioPartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithAudio]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Audio{
				TagIDs: models.NewRelatedIDs([]int{tagIDs[tagIdx2WithAudio]}),
			},
			false,
		},
		{
			"remove performers",
			audioIDs[audioIdxWithTwoPerformers],
			models.AudioPartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithAudio]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Audio{
				PerformerIDs: models.NewRelatedIDs([]int{performerIDs[performerIdx2WithAudio]}),
			},
			false,
		},
		{
			"remove groups",
			audioIDs[audioIdxWithGroup],
			models.AudioPartial{
				GroupIDs: &models.UpdateGroupIDsAudio{
					Groups: []models.GroupsAudios{
						{
							GroupID: groupIDs[groupIdxWithAudio],
						},
					},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Audio{
				Groups: models.NewRelatedGroupsAudio([]models.GroupsAudios{}),
			},
			false,
		},
		{
			"remove unrelated tags",
			audioIDs[audioIdxWithTwoTags],
			models.AudioPartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithPerformer]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Audio{
				TagIDs: models.NewRelatedIDs(indexesToIDs(tagIDs, audioTags[audioIdxWithTwoTags])),
			},
			false,
		},
		{
			"remove unrelated performers",
			audioIDs[audioIdxWithTwoPerformers],
			models.AudioPartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerIDs[performerIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Audio{
				PerformerIDs: models.NewRelatedIDs(indexesToIDs(performerIDs, audioPerformers[audioIdxWithTwoPerformers])),
			},
			false,
		},
		{
			"remove unrelated groups",
			audioIDs[audioIdxWithGroup],
			models.AudioPartial{
				GroupIDs: &models.UpdateGroupIDsAudio{
					Groups: []models.GroupsAudios{
						{
							GroupID: groupIDs[groupIdxWithDupName],
						},
					},
					Mode: models.RelationshipUpdateModeRemove,
				},
			},
			models.Audio{
				Groups: models.NewRelatedGroupsAudio([]models.GroupsAudios{
					{
						GroupID: indexesToIDs(groupIDs, audioGroups[audioIdxWithGroup])[0],
					},
				}),
			},
			false,
		},
	}

	for _, tt := range tests {
		qb := db.Audio

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			got, err := qb.UpdatePartial(ctx, tt.id, tt.partial)
			if (err != nil) != tt.wantErr {
				t.Errorf("audioQueryBuilder.UpdatePartial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			s, err := qb.Find(ctx, tt.id)
			if err != nil {
				t.Errorf("audioQueryBuilder.Find() error = %v", err)
			}

			// load relationships
			if err := loadAudioRelationships(ctx, tt.want, got); err != nil {
				t.Errorf("loadAudioRelationships() error = %v", err)
				return
			}
			if err := loadAudioRelationships(ctx, tt.want, s); err != nil {
				t.Errorf("loadAudioRelationships() error = %v", err)
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
			if tt.partial.GroupIDs != nil {
				assert.ElementsMatch(tt.want.Groups.List(), got.Groups.List())
				assert.ElementsMatch(tt.want.Groups.List(), s.Groups.List())
			}
		})
	}
}

func Test_audioQueryBuilder_AddO(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    int
		wantErr bool
	}{
		{
			"increment",
			audioIDs[1],
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

	qb := db.Audio

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.AddO(ctx, tt.id, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("audioQueryBuilder.AddO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("audioQueryBuilder.AddO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_audioQueryBuilder_DeleteO(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    int
		wantErr bool
	}{
		{
			"decrement",
			audioIDs[2],
			0,
			false,
		},
		{
			"zero",
			audioIDs[0],
			0,
			false,
		},
	}

	qb := db.Audio

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.DeleteO(ctx, tt.id, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("audioQueryBuilder.DeleteO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("audioQueryBuilder.DeleteO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_audioQueryBuilder_ResetO(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    int
		wantErr bool
	}{
		{
			"decrement",
			audioIDs[2],
			0,
			false,
		},
		{
			"zero",
			audioIDs[0],
			0,
			false,
		},
	}

	qb := db.Audio

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.ResetO(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("audioQueryBuilder.ResetO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("audioQueryBuilder.ResetOCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_audioQueryBuilder_ResetWatchCount(t *testing.T) {
	return
}

func Test_audioQueryBuilder_Destroy(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			"valid",
			audioIDs[audioIdxWithGallery],
			false,
		},
		{
			"invalid",
			invalidID,
			true,
		},
	}

	qb := db.Audio

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			if err := qb.Destroy(ctx, tt.id); (err != nil) != tt.wantErr {
				t.Errorf("audioQueryBuilder.Destroy() error = %v, wantErr %v", err, tt.wantErr)
			}

			// ensure cannot be found
			i, err := qb.Find(ctx, tt.id)

			assert.Nil(err)
			assert.Nil(i)
		})
	}
}

func makeAudioWithID(index int) *models.Audio {
	ret := makeAudio(index)
	ret.ID = audioIDs[index]

	ret.Files = models.NewRelatedAudioFiles([]*models.AudioFile{makeAudioFile(index)})

	return ret
}

func Test_audioQueryBuilder_Find(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    *models.Audio
		wantErr bool
	}{
		{
			"valid",
			audioIDs[audioIdxWithSpacedName],
			makeAudioWithID(audioIdxWithSpacedName),
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
			audioIDs[audioIdxWithGallery],
			makeAudioWithID(audioIdxWithGallery),
			false,
		},
		{
			"with performers",
			audioIDs[audioIdxWithTwoPerformers],
			makeAudioWithID(audioIdxWithTwoPerformers),
			false,
		},
		{
			"with tags",
			audioIDs[audioIdxWithTwoTags],
			makeAudioWithID(audioIdxWithTwoTags),
			false,
		},
		{
			"with groups",
			audioIDs[audioIdxWithGroup],
			makeAudioWithID(audioIdxWithGroup),
			false,
		},
	}

	qb := db.Audio

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.Find(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("audioQueryBuilder.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				// load relationships
				if err := loadAudioRelationships(ctx, *tt.want, got); err != nil {
					t.Errorf("loadAudioRelationships() error = %v", err)
					return
				}

				clearAudioFileIDs(got)
			}

			assert.Equal(tt.want, got)
		})
	}
}

func postFindAudios(ctx context.Context, want []*models.Audio, got []*models.Audio) error {
	for i, s := range got {
		// load relationships
		if i < len(want) {
			if err := loadAudioRelationships(ctx, *want[i], s); err != nil {
				return err
			}
		}
		clearAudioFileIDs(s)
	}

	return nil
}

func Test_audioQueryBuilder_FindMany(t *testing.T) {
	tests := []struct {
		name    string
		ids     []int
		want    []*models.Audio
		wantErr bool
	}{
		{
			"valid with relationships",
			[]int{
				audioIDs[audioIdxWithGallery],
				audioIDs[audioIdxWithTwoPerformers],
				audioIDs[audioIdxWithTwoTags],
				audioIDs[audioIdxWithGroup],
			},
			[]*models.Audio{
				makeAudioWithID(audioIdxWithGallery),
				makeAudioWithID(audioIdxWithTwoPerformers),
				makeAudioWithID(audioIdxWithTwoTags),
				makeAudioWithID(audioIdxWithGroup),
			},
			false,
		},
		{
			"invalid",
			[]int{audioIDs[audioIdxWithGallery], audioIDs[audioIdxWithTwoPerformers], invalidID},
			nil,
			true,
		},
	}

	qb := db.Audio

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindMany(ctx, tt.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("audioQueryBuilder.FindMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := postFindAudios(ctx, tt.want, got); err != nil {
				t.Errorf("loadAudioRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func Test_audioQueryBuilder_FindByChecksum(t *testing.T) {
	getChecksum := func(index int) string {
		return getAudioStringValue(index, checksumField)
	}

	tests := []struct {
		name     string
		checksum string
		want     []*models.Audio
		wantErr  bool
	}{
		{
			"valid",
			getChecksum(audioIdxWithSpacedName),
			[]*models.Audio{makeAudioWithID(audioIdxWithSpacedName)},
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
			getChecksum(audioIdxWithGallery),
			[]*models.Audio{makeAudioWithID(audioIdxWithGallery)},
			false,
		},
		{
			"with performers",
			getChecksum(audioIdxWithTwoPerformers),
			[]*models.Audio{makeAudioWithID(audioIdxWithTwoPerformers)},
			false,
		},
		{
			"with tags",
			getChecksum(audioIdxWithTwoTags),
			[]*models.Audio{makeAudioWithID(audioIdxWithTwoTags)},
			false,
		},
		{
			"with groups",
			getChecksum(audioIdxWithGroup),
			[]*models.Audio{makeAudioWithID(audioIdxWithGroup)},
			false,
		},
	}

	qb := db.Audio

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByChecksum(ctx, tt.checksum)
			if (err != nil) != tt.wantErr {
				t.Errorf("audioQueryBuilder.FindByChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := postFindAudios(ctx, tt.want, got); err != nil {
				t.Errorf("loadAudioRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func Test_audioQueryBuilder_FindByPath(t *testing.T) {
	getPath := func(index int) string {
		return getFilePath(folderIdxWithAudioFiles, getAudioBasename(index))
	}

	tests := []struct {
		name    string
		path    string
		want    []*models.Audio
		wantErr bool
	}{
		{
			"valid",
			getPath(audioIdxWithSpacedName),
			[]*models.Audio{makeAudioWithID(audioIdxWithSpacedName)},
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
			getPath(audioIdxWithGallery),
			[]*models.Audio{makeAudioWithID(audioIdxWithGallery)},
			false,
		},
		{
			"with performers",
			getPath(audioIdxWithTwoPerformers),
			[]*models.Audio{makeAudioWithID(audioIdxWithTwoPerformers)},
			false,
		},
		{
			"with tags",
			getPath(audioIdxWithTwoTags),
			[]*models.Audio{makeAudioWithID(audioIdxWithTwoTags)},
			false,
		},
		{
			"with groups",
			getPath(audioIdxWithGroup),
			[]*models.Audio{makeAudioWithID(audioIdxWithGroup)},
			false,
		},
	}

	qb := db.Audio

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByPath(ctx, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("audioQueryBuilder.FindByPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := postFindAudios(ctx, tt.want, got); err != nil {
				t.Errorf("loadAudioRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func TestAudioCountByPerformerID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio
		count, err := sqb.CountByPerformerID(ctx, performerIDs[performerIdxWithAudio])

		if err != nil {
			t.Errorf("Error counting audios: %s", err.Error())
		}

		assert.Equal(t, 1, count)

		count, err = sqb.CountByPerformerID(ctx, 0)

		if err != nil {
			t.Errorf("Error counting audios: %s", err.Error())
		}

		assert.Equal(t, 0, count)

		return nil
	})
}

func audiosToIDs(i []*models.Audio) []int {
	ret := make([]int, len(i))
	for i, v := range i {
		ret[i] = v.ID
	}

	return ret
}

func Test_audioStore_FindByFileID(t *testing.T) {
	tests := []struct {
		name    string
		fileID  models.FileID
		include []int
		exclude []int
	}{
		{
			"valid",
			audioFileIDs[audioIdx1WithPerformer],
			[]int{audioIdx1WithPerformer},
			nil,
		},
		{
			"invalid",
			invalidFileID,
			nil,
			[]int{audioIdx1WithPerformer},
		},
	}

	qb := db.Audio

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByFileID(ctx, tt.fileID)
			if err != nil {
				t.Errorf("AudioStore.FindByFileID() error = %v", err)
				return
			}
			for _, f := range got {
				clearAudioFileIDs(f)
			}

			ids := audiosToIDs(got)
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

func Test_audioStore_CountByFileID(t *testing.T) {
	tests := []struct {
		name   string
		fileID models.FileID
		want   int
	}{
		{
			"valid",
			audioFileIDs[audioIdxWithTwoPerformers],
			1,
		},
		{
			"invalid",
			invalidFileID,
			0,
		},
	}

	qb := db.Audio

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.CountByFileID(ctx, tt.fileID)
			if err != nil {
				t.Errorf("AudioStore.CountByFileID() error = %v", err)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func TestAudioQueryQ(t *testing.T) {
	const audioIdx = 2

	q := getAudioStringValue(audioIdx, titleField)

	withTxn(func(ctx context.Context) error {
		sqb := db.Audio

		audioQueryQ(ctx, t, sqb, q, audioIdx)

		return nil
	})
}

func queryAudio(ctx context.Context, t *testing.T, sqb models.AudioReader, audioFilter *models.AudioFilterType, findFilter *models.FindFilterType) []*models.Audio {
	t.Helper()
	result, err := sqb.Query(ctx, models.AudioQueryOptions{
		QueryOptions: models.QueryOptions{
			FindFilter: findFilter,
			Count:      true,
		},
		AudioFilter:   audioFilter,
		TotalDuration: true,
		TotalSize:     true,
	})
	if err != nil {
		t.Errorf("Error querying audio: %v", err)
		return nil
	}

	audios, err := result.Resolve(ctx)
	if err != nil {
		t.Errorf("Error resolving audios: %v", err)
	}

	return audios
}

func audioQueryQ(ctx context.Context, t *testing.T, sqb models.AudioReader, q string, expectedAudioIdx int) {
	filter := models.FindFilterType{
		Q: &q,
	}
	audios := queryAudio(ctx, t, sqb, nil, &filter)

	if !assert.Len(t, audios, 1) {
		return
	}
	audio := audios[0]
	assert.Equal(t, audioIDs[expectedAudioIdx], audio.ID)

	// no Q should return all results
	filter.Q = nil
	pp := totalAudios
	filter.PerPage = &pp
	audios = queryAudio(ctx, t, sqb, nil, &filter)

	assert.Len(t, audios, totalAudios)
}

func TestAudioQuery(t *testing.T) {
	var (
		depth = -1
	)

	tests := []struct {
		name        string
		findFilter  *models.FindFilterType
		filter      *models.AudioFilterType
		includeIdxs []int
		excludeIdxs []int
		wantErr     bool
	}{
		{
			"specific resume time",
			nil,
			&models.AudioFilterType{
				ResumeTime: &models.IntCriterionInput{
					Modifier: models.CriterionModifierEquals,
					Value:    int(getResumeTime(audioIdxWithGallery)),
				},
			},
			[]int{audioIdxWithGallery},
			[]int{audioIdxWithGroup},
			false,
		},
		{
			"specific play duration",
			nil,
			&models.AudioFilterType{
				PlayDuration: &models.IntCriterionInput{
					Modifier: models.CriterionModifierEquals,
					Value:    int(getPlayDuration(audioIdxWithGallery)),
				},
			},
			[]int{audioIdxWithGallery},
			[]int{audioIdxWithGroup},
			false,
		},
		// {
		// 	"specific play count",
		// 	nil,
		// 	&models.AudioFilterType{
		// 		PlayCount: &models.IntCriterionInput{
		// 			Modifier: models.CriterionModifierEquals,
		// 			Value:    getAudioPlayCount(audioIdxWithGallery),
		// 		},
		// 	},
		// 	[]int{audioIdxWithGallery},
		// 	[]int{audioIdxWithGroup},
		// 	false,
		// },
		{
			"with studio id 0 including child studios",
			nil,
			&models.AudioFilterType{
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

			results, err := db.Audio.Query(ctx, models.AudioQueryOptions{
				AudioFilter: tt.filter,
				QueryOptions: models.QueryOptions{
					FindFilter: tt.findFilter,
				},
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("AudioStore.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			include := indexesToIDs(audioIDs, tt.includeIdxs)
			exclude := indexesToIDs(audioIDs, tt.excludeIdxs)

			for _, i := range include {
				assert.Contains(results.IDs, i)
			}
			for _, e := range exclude {
				assert.NotContains(results.IDs, e)
			}
		})
	}
}

func TestAudioQueryPath(t *testing.T) {
	const (
		audioIdx      = 1
		otherAudioIdx = 2
	)
	folder := folderPaths[folderIdxWithAudioFiles]
	basename := getAudioBasename(audioIdx)
	audioPath := getFilePath(folderIdxWithAudioFiles, getAudioBasename(audioIdx))

	tests := []struct {
		name        string
		input       models.StringCriterionInput
		mustInclude []int
		mustExclude []int
	}{
		{
			"equals full path",
			models.StringCriterionInput{
				Value:    audioPath,
				Modifier: models.CriterionModifierEquals,
			},
			[]int{audioIdx},
			[]int{otherAudioIdx},
		},
		{
			"equals full path wildcard",
			models.StringCriterionInput{
				Value:    filepath.Join(folder, "audio_0001_%"),
				Modifier: models.CriterionModifierEquals,
			},
			[]int{audioIdx},
			[]int{otherAudioIdx},
		},
		{
			"not equals full path",
			models.StringCriterionInput{
				Value:    audioPath,
				Modifier: models.CriterionModifierNotEquals,
			},
			[]int{otherAudioIdx},
			[]int{audioIdx},
		},
		{
			"includes folder name",
			models.StringCriterionInput{
				Value:    folder,
				Modifier: models.CriterionModifierIncludes,
			},
			[]int{audioIdx},
			nil,
		},
		{
			"includes base name",
			models.StringCriterionInput{
				Value:    basename,
				Modifier: models.CriterionModifierIncludes,
			},
			[]int{audioIdx},
			nil,
		},
		{
			"includes full path",
			models.StringCriterionInput{
				Value:    audioPath,
				Modifier: models.CriterionModifierIncludes,
			},
			[]int{audioIdx},
			[]int{otherAudioIdx},
		},
		{
			"matches regex",
			models.StringCriterionInput{
				Value:    "audio_.*1_Path",
				Modifier: models.CriterionModifierMatchesRegex,
			},
			[]int{audioIdx},
			nil,
		},
		{
			"not matches regex",
			models.StringCriterionInput{
				Value:    "audio_.*1_Path",
				Modifier: models.CriterionModifierNotMatchesRegex,
			},
			nil,
			[]int{audioIdx},
		},
	}

	qb := db.Audio

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.Query(ctx, models.AudioQueryOptions{
				AudioFilter: &models.AudioFilterType{
					Path: &tt.input,
				},
			})

			if err != nil {
				t.Errorf("audioQueryBuilder.TestAudioQueryPath() error = %v", err)
				return
			}

			mustInclude := indexesToIDs(audioIDs, tt.mustInclude)
			mustExclude := indexesToIDs(audioIDs, tt.mustExclude)

			missing := sliceutil.Exclude(mustInclude, got.IDs)
			if len(missing) > 0 {
				t.Errorf("AudioStore.TestAudioQueryPath() missing expected IDs: %v", missing)
			}

			notExcluded := sliceutil.Intersect(mustExclude, got.IDs)
			if len(notExcluded) > 0 {
				t.Errorf("AudioStore.TestAudioQueryPath() expected IDs to be excluded: %v", notExcluded)
			}
		})
	}
}

func TestAudioQueryURL(t *testing.T) {
	const audioIdx = 1
	audioURL := getAudioStringValue(audioIdx, urlField)

	urlCriterion := models.StringCriterionInput{
		Value:    audioURL,
		Modifier: models.CriterionModifierEquals,
	}

	filter := models.AudioFilterType{
		URL: &urlCriterion,
	}

	verifyFn := func(s *models.Audio) {
		t.Helper()

		urls := s.URLs.List()
		var url string
		if len(urls) > 0 {
			url = urls[0]
		}

		verifyString(t, url, urlCriterion)
	}

	verifyAudioQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotEquals
	verifyAudioQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierMatchesRegex
	urlCriterion.Value = "audio_.*1_URL"
	verifyAudioQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotMatchesRegex
	verifyAudioQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierIsNull
	urlCriterion.Value = ""
	verifyAudioQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotNull
	verifyAudioQuery(t, filter, verifyFn)
}

func TestAudioQueryPathOr(t *testing.T) {
	const audio1Idx = 1
	const audio2Idx = 2

	audio1Path := getFilePath(folderIdxWithAudioFiles, getAudioBasename(audio1Idx))
	audio2Path := getFilePath(folderIdxWithAudioFiles, getAudioBasename(audio2Idx))

	audioFilter := models.AudioFilterType{
		Path: &models.StringCriterionInput{
			Value:    audio1Path,
			Modifier: models.CriterionModifierEquals,
		},
		OperatorFilter: models.OperatorFilter[models.AudioFilterType]{
			Or: &models.AudioFilterType{
				Path: &models.StringCriterionInput{
					Value:    audio2Path,
					Modifier: models.CriterionModifierEquals,
				},
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Audio

		audios := queryAudio(ctx, t, sqb, &audioFilter, nil)

		if !assert.Len(t, audios, 2) {
			return nil
		}
		assert.Equal(t, audio1Path, audios[0].Path)
		assert.Equal(t, audio2Path, audios[1].Path)

		return nil
	})
}

func TestAudioQueryPathAndRating(t *testing.T) {
	const audioIdx = 1
	audioPath := getFilePath(folderIdxWithAudioFiles, getAudioBasename(audioIdx))
	audioRating := int(getRating(audioIdx).Int64)

	audioFilter := models.AudioFilterType{
		Path: &models.StringCriterionInput{
			Value:    audioPath,
			Modifier: models.CriterionModifierEquals,
		},
		OperatorFilter: models.OperatorFilter[models.AudioFilterType]{
			And: &models.AudioFilterType{
				Rating100: &models.IntCriterionInput{
					Value:    audioRating,
					Modifier: models.CriterionModifierEquals,
				},
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Audio

		audios := queryAudio(ctx, t, sqb, &audioFilter, nil)

		if !assert.Len(t, audios, 1) {
			return nil
		}
		assert.Equal(t, audioPath, audios[0].Path)
		assert.Equal(t, audioRating, *audios[0].Rating)

		return nil
	})
}

func TestAudioQueryPathNotRating(t *testing.T) {
	const audioIdx = 1

	audioRating := getRating(audioIdx)

	pathCriterion := models.StringCriterionInput{
		Value:    "audio_.*1_Path",
		Modifier: models.CriterionModifierMatchesRegex,
	}

	ratingCriterion := models.IntCriterionInput{
		Value:    int(audioRating.Int64),
		Modifier: models.CriterionModifierEquals,
	}

	audioFilter := models.AudioFilterType{
		Path: &pathCriterion,
		OperatorFilter: models.OperatorFilter[models.AudioFilterType]{
			Not: &models.AudioFilterType{
				Rating100: &ratingCriterion,
			},
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Audio

		audios := queryAudio(ctx, t, sqb, &audioFilter, nil)

		for _, audio := range audios {
			verifyString(t, audio.Path, pathCriterion)
			ratingCriterion.Modifier = models.CriterionModifierNotEquals
			verifyIntPtr(t, audio.Rating, ratingCriterion)
		}

		return nil
	})
}

func TestAudioIllegalQuery(t *testing.T) {
	assert := assert.New(t)

	const audioIdx = 1
	subFilter := models.AudioFilterType{
		Path: &models.StringCriterionInput{
			Value:    getAudioStringValue(audioIdx, "Path"),
			Modifier: models.CriterionModifierEquals,
		},
	}

	audioFilter := &models.AudioFilterType{
		OperatorFilter: models.OperatorFilter[models.AudioFilterType]{
			And: &subFilter,
			Or:  &subFilter,
		},
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Audio

		queryOptions := models.AudioQueryOptions{
			AudioFilter: audioFilter,
		}

		_, err := sqb.Query(ctx, queryOptions)
		assert.NotNil(err)

		audioFilter.Or = nil
		audioFilter.Not = &subFilter
		_, err = sqb.Query(ctx, queryOptions)
		assert.NotNil(err)

		audioFilter.And = nil
		audioFilter.Or = &subFilter
		_, err = sqb.Query(ctx, queryOptions)
		assert.NotNil(err)

		return nil
	})
}

func verifyAudioQuery(t *testing.T, filter models.AudioFilterType, verifyFn func(s *models.Audio)) {
	t.Helper()
	withTxn(func(ctx context.Context) error {
		t.Helper()
		sqb := db.Audio

		audios := queryAudio(ctx, t, sqb, &filter, nil)

		for _, audio := range audios {
			if err := audio.LoadRelationships(ctx, sqb); err != nil {
				t.Errorf("Error loading audio relationships: %v", err)
			}
		}

		// assume it should find at least one
		assert.Greater(t, len(audios), 0)

		for _, audio := range audios {
			verifyFn(audio)
		}

		return nil
	})
}

func verifyAudiosPath(t *testing.T, pathCriterion models.StringCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio
		audioFilter := models.AudioFilterType{
			Path: &pathCriterion,
		}

		audios := queryAudio(ctx, t, sqb, &audioFilter, nil)

		for _, audio := range audios {
			verifyString(t, audio.Path, pathCriterion)
		}

		return nil
	})
}

func TestAudioQueryRating100(t *testing.T) {
	const rating = 60
	ratingCriterion := models.IntCriterionInput{
		Value:    rating,
		Modifier: models.CriterionModifierEquals,
	}

	verifyAudiosRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotEquals
	verifyAudiosRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyAudiosRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierLessThan
	verifyAudiosRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierIsNull
	verifyAudiosRating100(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotNull
	verifyAudiosRating100(t, ratingCriterion)
}

func verifyAudiosRating100(t *testing.T, ratingCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio
		audioFilter := models.AudioFilterType{
			Rating100: &ratingCriterion,
		}

		audios := queryAudio(ctx, t, sqb, &audioFilter, nil)

		for _, audio := range audios {
			verifyIntPtr(t, audio.Rating, ratingCriterion)
		}

		return nil
	})
}

func TestAudioQueryOCounter(t *testing.T) {
	const oCounter = 1
	oCounterCriterion := models.IntCriterionInput{
		Value:    oCounter,
		Modifier: models.CriterionModifierEquals,
	}

	verifyAudiosOCounter(t, oCounterCriterion)

	oCounterCriterion.Modifier = models.CriterionModifierNotEquals
	verifyAudiosOCounter(t, oCounterCriterion)

	oCounterCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyAudiosOCounter(t, oCounterCriterion)

	oCounterCriterion.Modifier = models.CriterionModifierLessThan
	verifyAudiosOCounter(t, oCounterCriterion)
}

func verifyAudiosOCounter(t *testing.T, oCounterCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio
		audioFilter := models.AudioFilterType{
			OCounter: &oCounterCriterion,
		}

		audios := queryAudio(ctx, t, sqb, &audioFilter, nil)

		for _, audio := range audios {
			count, err := sqb.GetOCount(ctx, audio.ID)
			if err != nil {
				t.Errorf("Error getting ocounter: %v", err)
			}
			verifyInt(t, count, oCounterCriterion)
		}

		return nil
	})
}

func TestAudioQueryDuration(t *testing.T) {
	duration := 200.432

	durationCriterion := models.IntCriterionInput{
		Value:    int(duration),
		Modifier: models.CriterionModifierEquals,
	}
	verifyAudiosDuration(t, durationCriterion)

	durationCriterion.Modifier = models.CriterionModifierNotEquals
	verifyAudiosDuration(t, durationCriterion)

	durationCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyAudiosDuration(t, durationCriterion)

	durationCriterion.Modifier = models.CriterionModifierLessThan
	verifyAudiosDuration(t, durationCriterion)

	durationCriterion.Modifier = models.CriterionModifierIsNull
	verifyAudiosDuration(t, durationCriterion)

	durationCriterion.Modifier = models.CriterionModifierNotNull
	verifyAudiosDuration(t, durationCriterion)
}

func verifyAudiosDuration(t *testing.T, durationCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio
		audioFilter := models.AudioFilterType{
			Duration: &durationCriterion,
		}

		audios := queryAudio(ctx, t, sqb, &audioFilter, nil)

		for _, audio := range audios {
			if err := audio.LoadPrimaryFile(ctx, db.File); err != nil {
				t.Errorf("Error querying audio files: %v", err)
				return nil
			}

			duration := audio.Files.Primary().Duration
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

func TestAudioQueryIsMissingStudio(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio
		isMissing := "studio"
		audioFilter := models.AudioFilterType{
			IsMissing: &isMissing,
		}

		q := getAudioStringValue(audioIdxWithStudio, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		audios := queryAudio(ctx, t, sqb, &audioFilter, &findFilter)

		assert.Len(t, audios, 0)

		findFilter.Q = nil
		audios = queryAudio(ctx, t, sqb, &audioFilter, &findFilter)

		// ensure non of the ids equal the one with studio
		for _, audio := range audios {
			assert.NotEqual(t, audioIDs[audioIdxWithStudio], audio.ID)
		}

		return nil
	})
}

func TestAudioQueryIsMissingMovies(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio
		isMissing := "movie"
		audioFilter := models.AudioFilterType{
			IsMissing: &isMissing,
		}

		q := getAudioStringValue(audioIdxWithGroup, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		audios := queryAudio(ctx, t, sqb, &audioFilter, &findFilter)

		assert.Len(t, audios, 0)

		findFilter.Q = nil
		audios = queryAudio(ctx, t, sqb, &audioFilter, &findFilter)

		// ensure non of the ids equal the one with movies
		for _, audio := range audios {
			assert.NotEqual(t, audioIDs[audioIdxWithGroup], audio.ID)
		}

		return nil
	})
}

func TestAudioQueryIsMissingPerformers(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio
		isMissing := "performers"
		audioFilter := models.AudioFilterType{
			IsMissing: &isMissing,
		}

		q := getAudioStringValue(audioIdxWithPerformer, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		audios := queryAudio(ctx, t, sqb, &audioFilter, &findFilter)

		assert.Len(t, audios, 0)

		findFilter.Q = nil
		audios = queryAudio(ctx, t, sqb, &audioFilter, &findFilter)

		assert.True(t, len(audios) > 0)

		// ensure non of the ids equal the one with movies
		for _, audio := range audios {
			assert.NotEqual(t, audioIDs[audioIdxWithPerformer], audio.ID)
		}

		return nil
	})
}

func TestAudioQueryIsMissingDate(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio
		isMissing := "date"
		audioFilter := models.AudioFilterType{
			IsMissing: &isMissing,
		}

		audios := queryAudio(ctx, t, sqb, &audioFilter, nil)

		// one in four audios have no date
		assert.Len(t, audios, int(math.Ceil(float64(totalAudios)/4)))

		// ensure date is null
		for _, audio := range audios {
			assert.Nil(t, audio.Date)
		}

		return nil
	})
}

func TestAudioQueryIsMissingTags(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio
		isMissing := "tags"
		audioFilter := models.AudioFilterType{
			IsMissing: &isMissing,
		}

		q := getAudioStringValue(audioIdxWithTwoTags, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		audios := queryAudio(ctx, t, sqb, &audioFilter, &findFilter)

		assert.Len(t, audios, 0)

		findFilter.Q = nil
		audios = queryAudio(ctx, t, sqb, &audioFilter, &findFilter)

		assert.True(t, len(audios) > 0)

		return nil
	})
}

func TestAudioQueryIsMissingRating(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio
		isMissing := "rating"
		audioFilter := models.AudioFilterType{
			IsMissing: &isMissing,
		}

		audios := queryAudio(ctx, t, sqb, &audioFilter, nil)

		assert.True(t, len(audios) > 0)

		// ensure rating is null
		for _, audio := range audios {
			assert.Nil(t, audio.Rating)
		}

		return nil
	})
}

func TestAudioQueryPerformers(t *testing.T) {
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
					strconv.Itoa(performerIDs[performerIdxWithAudio]),
					strconv.Itoa(performerIDs[performerIdx1WithAudio]),
				},
				Modifier: models.CriterionModifierIncludes,
			},
			[]int{
				audioIdxWithPerformer,
				audioIdxWithTwoPerformers,
			},
			[]int{
				audioIdxWithGallery,
			},
			false,
		},
		{
			"includes all",
			models.MultiCriterionInput{
				Value: []string{
					strconv.Itoa(performerIDs[performerIdx1WithAudio]),
					strconv.Itoa(performerIDs[performerIdx2WithAudio]),
				},
				Modifier: models.CriterionModifierIncludesAll,
			},
			[]int{
				audioIdxWithTwoPerformers,
			},
			[]int{
				audioIdxWithPerformer,
			},
			false,
		},
		{
			"excludes",
			models.MultiCriterionInput{
				Modifier: models.CriterionModifierExcludes,
				Value:    []string{strconv.Itoa(tagIDs[performerIdx1WithAudio])},
			},
			nil,
			[]int{audioIdxWithTwoPerformers},
			false,
		},
		{
			"is null",
			models.MultiCriterionInput{
				Modifier: models.CriterionModifierIsNull,
			},
			[]int{audioIdxWithTag},
			[]int{
				audioIdxWithPerformer,
				audioIdxWithTwoPerformers,
				audioIdxWithPerformerTwoTags,
			},
			false,
		},
		{
			"not null",
			models.MultiCriterionInput{
				Modifier: models.CriterionModifierNotNull,
			},
			[]int{
				audioIdxWithPerformer,
				audioIdxWithTwoPerformers,
				audioIdxWithPerformerTwoTags,
			},
			[]int{audioIdxWithTag},
			false,
		},
		{
			"equals",
			models.MultiCriterionInput{
				Modifier: models.CriterionModifierEquals,
				Value: []string{
					strconv.Itoa(tagIDs[performerIdx1WithAudio]),
					strconv.Itoa(tagIDs[performerIdx2WithAudio]),
				},
			},
			[]int{audioIdxWithTwoPerformers},
			[]int{
				audioIdxWithThreePerformers,
			},
			false,
		},
		{
			"not equals",
			models.MultiCriterionInput{
				Modifier: models.CriterionModifierNotEquals,
				Value: []string{
					strconv.Itoa(tagIDs[performerIdx1WithAudio]),
					strconv.Itoa(tagIDs[performerIdx2WithAudio]),
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

			results, err := db.Audio.Query(ctx, models.AudioQueryOptions{
				AudioFilter: &models.AudioFilterType{
					Performers: &tt.filter,
				},
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("AudioStore.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			include := indexesToIDs(audioIDs, tt.includeIdxs)
			exclude := indexesToIDs(audioIDs, tt.excludeIdxs)

			for _, i := range include {
				assert.Contains(results.IDs, i)
			}
			for _, e := range exclude {
				assert.NotContains(results.IDs, e)
			}
		})
	}
}

func TestAudioQueryTags(t *testing.T) {
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
					strconv.Itoa(tagIDs[tagIdxWithAudio]),
					strconv.Itoa(tagIDs[tagIdx1WithAudio]),
				},
				Modifier: models.CriterionModifierIncludes,
			},
			[]int{
				audioIdxWithTag,
				audioIdxWithTwoTags,
			},
			[]int{
				audioIdxWithGallery,
			},
			false,
		},
		{
			"includes all",
			models.HierarchicalMultiCriterionInput{
				Value: []string{
					strconv.Itoa(tagIDs[tagIdx1WithAudio]),
					strconv.Itoa(tagIDs[tagIdx2WithAudio]),
				},
				Modifier: models.CriterionModifierIncludesAll,
			},
			[]int{
				audioIdxWithTwoTags,
			},
			[]int{
				audioIdxWithTag,
			},
			false,
		},
		{
			"excludes",
			models.HierarchicalMultiCriterionInput{
				Modifier: models.CriterionModifierExcludes,
				Value:    []string{strconv.Itoa(tagIDs[tagIdx1WithAudio])},
			},
			nil,
			[]int{audioIdxWithTwoTags},
			false,
		},
		{
			"is null",
			models.HierarchicalMultiCriterionInput{
				Modifier: models.CriterionModifierIsNull,
			},
			[]int{audioIdx1WithPerformer},
			[]int{
				audioIdxWithTag,
				audioIdxWithTwoTags,
			},
			false,
		},
		{
			"not null",
			models.HierarchicalMultiCriterionInput{
				Modifier: models.CriterionModifierNotNull,
			},
			[]int{
				audioIdxWithTag,
				audioIdxWithTwoTags,
			},
			[]int{audioIdx1WithPerformer},
			false,
		},
		{
			"equals",
			models.HierarchicalMultiCriterionInput{
				Modifier: models.CriterionModifierEquals,
				Value: []string{
					strconv.Itoa(tagIDs[tagIdx1WithAudio]),
					strconv.Itoa(tagIDs[tagIdx2WithAudio]),
				},
			},
			[]int{audioIdxWithTwoTags},
			[]int{
				audioIdxWithThreeTags,
			},
			false,
		},
		{
			"not equals",
			models.HierarchicalMultiCriterionInput{
				Modifier: models.CriterionModifierNotEquals,
				Value: []string{
					strconv.Itoa(tagIDs[tagIdx1WithAudio]),
					strconv.Itoa(tagIDs[tagIdx2WithAudio]),
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

			results, err := db.Audio.Query(ctx, models.AudioQueryOptions{
				AudioFilter: &models.AudioFilterType{
					Tags: &tt.filter,
				},
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("AudioStore.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			include := indexesToIDs(audioIDs, tt.includeIdxs)
			exclude := indexesToIDs(audioIDs, tt.excludeIdxs)

			for _, i := range include {
				assert.Contains(results.IDs, i)
			}
			for _, e := range exclude {
				assert.NotContains(results.IDs, e)
			}
		})
	}
}

func TestAudioQueryPerformerTags(t *testing.T) {
	allDepth := -1

	tests := []struct {
		name        string
		findFilter  *models.FindFilterType
		filter      *models.AudioFilterType
		includeIdxs []int
		excludeIdxs []int
		wantErr     bool
	}{
		{
			"includes",
			nil,
			&models.AudioFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Value: []string{
						strconv.Itoa(tagIDs[tagIdxWithPerformer]),
						strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
					},
					Modifier: models.CriterionModifierIncludes,
				},
			},
			[]int{
				audioIdxWithPerformerTag,
				audioIdxWithPerformerTwoTags,
				audioIdxWithTwoPerformerTag,
			},
			[]int{
				audioIdxWithPerformer,
			},
			false,
		},
		{
			"includes sub-tags",
			nil,
			&models.AudioFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Value: []string{
						strconv.Itoa(tagIDs[tagIdxWithParentAndChild]),
					},
					Depth:    &allDepth,
					Modifier: models.CriterionModifierIncludes,
				},
			},
			[]int{
				audioIdxWithPerformerParentTag,
			},
			[]int{
				audioIdxWithPerformer,
				audioIdxWithPerformerTag,
				audioIdxWithPerformerTwoTags,
				audioIdxWithTwoPerformerTag,
			},
			false,
		},
		{
			"includes all",
			nil,
			&models.AudioFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Value: []string{
						strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
						strconv.Itoa(tagIDs[tagIdx2WithPerformer]),
					},
					Modifier: models.CriterionModifierIncludesAll,
				},
			},
			[]int{
				audioIdxWithPerformerTwoTags,
			},
			[]int{
				audioIdxWithPerformer,
				audioIdxWithPerformerTag,
				audioIdxWithTwoPerformerTag,
			},
			false,
		},
		{
			"excludes performer tag tagIdx2WithPerformer",
			nil,
			&models.AudioFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Modifier: models.CriterionModifierExcludes,
					Value:    []string{strconv.Itoa(tagIDs[tagIdx2WithPerformer])},
				},
			},
			nil,
			[]int{audioIdxWithTwoPerformerTag},
			false,
		},
		{
			"excludes sub-tags",
			nil,
			&models.AudioFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Value: []string{
						strconv.Itoa(tagIDs[tagIdxWithParentAndChild]),
					},
					Depth:    &allDepth,
					Modifier: models.CriterionModifierExcludes,
				},
			},
			[]int{
				audioIdxWithPerformer,
				audioIdxWithPerformerTag,
				audioIdxWithPerformerTwoTags,
				audioIdxWithTwoPerformerTag,
			},
			[]int{
				audioIdxWithPerformerParentTag,
			},
			false,
		},
		{
			"is null",
			nil,
			&models.AudioFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Modifier: models.CriterionModifierIsNull,
				},
			},
			[]int{audioIdx1WithPerformer},
			[]int{audioIdxWithPerformerTag},
			false,
		},
		{
			"not null",
			nil,
			&models.AudioFilterType{
				PerformerTags: &models.HierarchicalMultiCriterionInput{
					Modifier: models.CriterionModifierNotNull,
				},
			},
			[]int{audioIdxWithPerformerTag},
			[]int{audioIdx1WithPerformer},
			false,
		},
		{
			"equals",
			nil,
			&models.AudioFilterType{
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
			&models.AudioFilterType{
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

			results, err := db.Audio.Query(ctx, models.AudioQueryOptions{
				AudioFilter: tt.filter,
				QueryOptions: models.QueryOptions{
					FindFilter: tt.findFilter,
				},
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("AudioStore.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			include := indexesToIDs(audioIDs, tt.includeIdxs)
			exclude := indexesToIDs(audioIDs, tt.excludeIdxs)

			for _, i := range include {
				assert.Contains(results.IDs, i)
			}
			for _, e := range exclude {
				assert.NotContains(results.IDs, e)
			}
		})
	}
}

func TestAudioQueryStudio(t *testing.T) {
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
					strconv.Itoa(studioIDs[studioIdxWithAudio]),
				},
				Modifier: models.CriterionModifierIncludes,
			},
			[]int{audioIDs[audioIdxWithStudio]},
			false,
		},
		{
			"excludes",
			getAudioStringValue(audioIdxWithStudio, titleField),
			models.HierarchicalMultiCriterionInput{
				Value: []string{
					strconv.Itoa(studioIDs[studioIdxWithAudio]),
				},
				Modifier: models.CriterionModifierExcludes,
			},
			[]int{},
			false,
		},
		{
			"excludes includes null",
			getAudioStringValue(audioIdxWithGallery, titleField),
			models.HierarchicalMultiCriterionInput{
				Value: []string{
					strconv.Itoa(studioIDs[studioIdxWithAudio]),
				},
				Modifier: models.CriterionModifierExcludes,
			},
			[]int{audioIDs[audioIdxWithGallery]},
			false,
		},
		{
			"equals",
			"",
			models.HierarchicalMultiCriterionInput{
				Value: []string{
					strconv.Itoa(studioIDs[studioIdxWithAudio]),
				},
				Modifier: models.CriterionModifierEquals,
			},
			[]int{audioIDs[audioIdxWithStudio]},
			false,
		},
		{
			"not equals",
			getAudioStringValue(audioIdxWithStudio, titleField),
			models.HierarchicalMultiCriterionInput{
				Value: []string{
					strconv.Itoa(studioIDs[studioIdxWithAudio]),
				},
				Modifier: models.CriterionModifierNotEquals,
			},
			[]int{},
			false,
		},
	}

	qb := db.Audio

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			studioCriterion := tt.studioCriterion

			audioFilter := models.AudioFilterType{
				Studios: &studioCriterion,
			}

			var findFilter *models.FindFilterType
			if tt.q != "" {
				findFilter = &models.FindFilterType{
					Q: &tt.q,
				}
			}

			audios := queryAudio(ctx, t, qb, &audioFilter, findFilter)

			assert.ElementsMatch(t, audiosToIDs(audios), tt.expectedIDs)
		})
	}
}

func TestAudioQueryStudioDepth(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio
		depth := 2
		studioCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithGrandChild]),
			},
			Modifier: models.CriterionModifierIncludes,
			Depth:    &depth,
		}

		audioFilter := models.AudioFilterType{
			Studios: &studioCriterion,
		}

		audios := queryAudio(ctx, t, sqb, &audioFilter, nil)
		assert.Len(t, audios, 1)

		depth = 1

		audios = queryAudio(ctx, t, sqb, &audioFilter, nil)
		assert.Len(t, audios, 0)

		studioCriterion.Value = []string{strconv.Itoa(studioIDs[studioIdxWithParentAndChild])}
		audios = queryAudio(ctx, t, sqb, &audioFilter, nil)
		assert.Len(t, audios, 1)

		// ensure id is correct
		assert.Equal(t, audioIDs[audioIdxWithGrandChildStudio], audios[0].ID)
		depth = 2

		studioCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithGrandChild]),
			},
			Modifier: models.CriterionModifierExcludes,
			Depth:    &depth,
		}

		q := getAudioStringValue(audioIdxWithGrandChildStudio, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		audios = queryAudio(ctx, t, sqb, &audioFilter, &findFilter)
		assert.Len(t, audios, 0)

		depth = 1
		audios = queryAudio(ctx, t, sqb, &audioFilter, &findFilter)
		assert.Len(t, audios, 1)

		studioCriterion.Value = []string{strconv.Itoa(studioIDs[studioIdxWithParentAndChild])}
		audios = queryAudio(ctx, t, sqb, &audioFilter, &findFilter)
		assert.Len(t, audios, 0)

		return nil
	})
}

func TestAudioGroups(t *testing.T) {
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
		excludeIdxs []int
	}{
		{
			"includes",
			criterion{
				[]int{groupIdxWithAudio},
				models.CriterionModifierIncludes,
				0,
			},
			"",
			[]int{audioIdxWithGroup},
			nil,
		},
		{
			"excludes",
			criterion{
				[]int{groupIdxWithAudio},
				models.CriterionModifierExcludes,
				0,
			},
			getAudioStringValue(audioIdxWithGroup, titleField),
			nil,
			[]int{audioIdxWithGroup},
		},
		{
			"includes (depth = 1)",
			criterion{
				[]int{groupIdxWithChildWithAudio},
				models.CriterionModifierIncludes,
				1,
			},
			"",
			[]int{audioIdxWithGroupWithParent},
			nil,
		},
	}

	for _, tt := range tests {
		valueIDs := indexesToIDs(groupIDs, tt.c.valueIdxs)

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			audioFilter := &models.AudioFilterType{
				Groups: &models.HierarchicalMultiCriterionInput{
					Value:    intslice.IntSliceToStringSlice(valueIDs),
					Modifier: tt.c.modifier,
				},
			}

			if tt.c.depth != 0 {
				audioFilter.Groups.Depth = &tt.c.depth
			}

			findFilter := &models.FindFilterType{}
			if tt.q != "" {
				findFilter.Q = &tt.q
			}

			results, err := db.Audio.Query(ctx, models.AudioQueryOptions{
				AudioFilter: audioFilter,
				QueryOptions: models.QueryOptions{
					FindFilter: findFilter,
				},
			})
			if err != nil {
				t.Errorf("AudioStore.Query() error = %v", err)
				return
			}

			include := indexesToIDs(audioIDs, tt.includeIdxs)
			exclude := indexesToIDs(audioIDs, tt.excludeIdxs)

			assert.Subset(results.IDs, include)

			for _, e := range exclude {
				assert.NotContains(results.IDs, e)
			}
		})
	}
}

func TestAudioQuerySorting(t *testing.T) {
	tests := []struct {
		name          string
		sortBy        string
		dir           models.SortDirectionEnum
		firstAudioIdx int // -1 to ignore
		lastAudioIdx  int
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
			"sample rate",
			"sample_rate",
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
			audioIDs[audioIdx1WithPerformer],
			-1,
		},
		{
			"play_duration",
			"play_duration",
			models.SortDirectionEnumDesc,
			audioIDs[audioIdx1WithPerformer],
			-1,
		},
		{
			"performer_age",
			"performer_age",
			models.SortDirectionEnumDesc,
			-1,
			-1,
		},
	}

	qb := db.Audio

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.Query(ctx, models.AudioQueryOptions{
				QueryOptions: models.QueryOptions{
					FindFilter: &models.FindFilterType{
						Sort:      &tt.sortBy,
						Direction: &tt.dir,
					},
				},
			})

			if err != nil {
				t.Errorf("audioQueryBuilder.TestAudioQuerySorting() error = %v", err)
				return
			}

			audios, err := got.Resolve(ctx)
			if err != nil {
				t.Errorf("audioQueryBuilder.TestAudioQuerySorting() error = %v", err)
				return
			}

			if !assert.Greater(len(audios), 0) {
				return
			}

			// audios should be in same order as indexes
			firstAudio := audios[0]
			lastAudio := audios[len(audios)-1]

			if tt.firstAudioIdx != -1 {
				firstAudioID := audioIDs[tt.firstAudioIdx]
				assert.Equal(firstAudioID, firstAudio.ID)
			}
			if tt.lastAudioIdx != -1 {
				lastAudioID := audioIDs[tt.lastAudioIdx]
				assert.Equal(lastAudioID, lastAudio.ID)
			}
		})
	}
}

func TestAudioQueryPagination(t *testing.T) {
	perPage := 1
	findFilter := models.FindFilterType{
		PerPage: &perPage,
	}

	withTxn(func(ctx context.Context) error {
		sqb := db.Audio
		audios := queryAudio(ctx, t, sqb, nil, &findFilter)

		assert.Len(t, audios, 1)

		firstID := audios[0].ID

		page := 2
		findFilter.Page = &page
		audios = queryAudio(ctx, t, sqb, nil, &findFilter)

		assert.Len(t, audios, 1)
		secondID := audios[0].ID
		assert.NotEqual(t, firstID, secondID)

		perPage = 2
		page = 1

		audios = queryAudio(ctx, t, sqb, nil, &findFilter)
		assert.Len(t, audios, 2)
		assert.Equal(t, firstID, audios[0].ID)
		assert.Equal(t, secondID, audios[1].ID)

		return nil
	})
}

func TestAudioQueryTagCount(t *testing.T) {
	const tagCount = 1
	tagCountCriterion := models.IntCriterionInput{
		Value:    tagCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyAudiosTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyAudiosTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyAudiosTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyAudiosTagCount(t, tagCountCriterion)
}

func verifyAudiosTagCount(t *testing.T, tagCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio
		audioFilter := models.AudioFilterType{
			TagCount: &tagCountCriterion,
		}

		audios := queryAudio(ctx, t, sqb, &audioFilter, nil)
		assert.Greater(t, len(audios), 0)

		for _, audio := range audios {
			if err := audio.LoadTagIDs(ctx, sqb); err != nil {
				t.Errorf("audio.LoadTagIDs() error = %v", err)
				return nil
			}
			verifyInt(t, len(audio.TagIDs.List()), tagCountCriterion)
		}

		return nil
	})
}

func TestAudioQueryPerformerCount(t *testing.T) {
	const performerCount = 1
	performerCountCriterion := models.IntCriterionInput{
		Value:    performerCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyAudiosPerformerCount(t, performerCountCriterion)

	performerCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyAudiosPerformerCount(t, performerCountCriterion)

	performerCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyAudiosPerformerCount(t, performerCountCriterion)

	performerCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyAudiosPerformerCount(t, performerCountCriterion)
}

func verifyAudiosPerformerCount(t *testing.T, performerCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio
		audioFilter := models.AudioFilterType{
			PerformerCount: &performerCountCriterion,
		}

		audios := queryAudio(ctx, t, sqb, &audioFilter, nil)
		assert.Greater(t, len(audios), 0)

		for _, audio := range audios {
			if err := audio.LoadPerformerIDs(ctx, sqb); err != nil {
				t.Errorf("audio.LoadPerformerIDs() error = %v", err)
				return nil
			}

			verifyInt(t, len(audio.PerformerIDs.List()), performerCountCriterion)
		}

		return nil
	})
}

func TestAudioFindByMovieID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio

		audios, err := sqb.FindByGroupID(ctx, groupIDs[groupIdxWithAudio])

		if err != nil {
			t.Errorf("error calling FindByMovieID: %s", err.Error())
		}

		assert.Len(t, audios, 1)
		assert.Equal(t, audioIDs[audioIdxWithGroup], audios[0].ID)

		audios, err = sqb.FindByGroupID(ctx, 0)

		if err != nil {
			t.Errorf("error calling FindByMovieID: %s", err.Error())
		}

		assert.Len(t, audios, 0)

		return nil
	})
}

func TestAudioFindByPerformerID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Audio

		audios, err := sqb.FindByPerformerID(ctx, performerIDs[performerIdxWithAudio])

		if err != nil {
			t.Errorf("error calling FindByPerformerID: %s", err.Error())
		}

		assert.Len(t, audios, 1)
		assert.Equal(t, audioIDs[audioIdxWithPerformer], audios[0].ID)

		audios, err = sqb.FindByPerformerID(ctx, 0)

		if err != nil {
			t.Errorf("error calling FindByPerformerID: %s", err.Error())
		}

		assert.Len(t, audios, 0)

		return nil
	})
}

func TestAudioUpdateAudioCover(t *testing.T) {
	if err := withTxn(func(ctx context.Context) error {
		qb := db.Audio

		audioID := audioIDs[audioIdxWithGallery]

		return testUpdateImage(t, ctx, audioID, qb.UpdateCover, qb.GetCover)
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestAudioQueryQTrim(t *testing.T) {
	if err := withTxn(func(ctx context.Context) error {
		qb := db.Audio

		expectedID := audioIDs[audioIdxWithSpacedName]

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
			audios := queryAudio(ctx, t, qb, nil, &f)

			assert.Len(t, audios, tst.count)
			if len(audios) > 0 {
				assert.Equal(t, tst.id, audios[0].ID)
			}
		}

		findFilter := models.FindFilterType{}
		audios := queryAudio(ctx, t, qb, nil, &findFilter)
		assert.NotEqual(t, 0, len(audios))

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestAudioStore_All(t *testing.T) {
	qb := db.Audio

	withRollbackTxn(func(ctx context.Context) error {
		got, err := qb.All(ctx)
		if err != nil {
			t.Errorf("AudioStore.All() error = %v", err)
			return nil
		}

		// it's possible that other tests have created audios
		assert.GreaterOrEqual(t, len(got), len(audioIDs))

		return nil
	})
}

func TestAudioStore_AssignFiles(t *testing.T) {
	tests := []struct {
		name    string
		audioID int
		fileID  models.FileID
		wantErr bool
	}{
		{
			"valid",
			audioIDs[audioIdx1WithPerformer],
			audioFileIDs[audioIdx1WithStudio],
			false,
		},
		{
			"invalid file id",
			audioIDs[audioIdx1WithPerformer],
			invalidFileID,
			true,
		},
		{
			"invalid audio id",
			invalidID,
			audioFileIDs[audioIdx1WithStudio],
			true,
		},
	}

	qb := db.Audio

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withRollbackTxn(func(ctx context.Context) error {
				if err := qb.AssignFiles(ctx, tt.audioID, []models.FileID{tt.fileID}); (err != nil) != tt.wantErr {
					t.Errorf("AudioStore.AssignFiles() error = %v, wantErr %v", err, tt.wantErr)
				}

				return nil
			})
		})
	}
}

func TestAudioStore_AddView(t *testing.T) {
	tests := []struct {
		name          string
		audioID       int
		expectedCount int
		wantErr       bool
	}{
		{
			"valid",
			audioIDs[audioIdx1WithPerformer],
			1, //getAudioPlayCount(audioIdx1WithPerformer) + 1,
			false,
		},
		{
			"invalid audio id",
			invalidID,
			0,
			true,
		},
	}

	qb := db.Audio

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withRollbackTxn(func(ctx context.Context) error {
				views, err := qb.AddViews(ctx, tt.audioID, nil)
				if (err != nil) != tt.wantErr {
					t.Errorf("AudioStore.AddView() error = %v, wantErr %v", err, tt.wantErr)
				}

				if err != nil {
					return nil
				}

				assert := assert.New(t)
				assert.Equal(tt.expectedCount, len(views))

				// find the audio and check the count
				count, err := qb.CountViews(ctx, tt.audioID)
				if err != nil {
					t.Errorf("AudioStore.CountViews() error = %v", err)
				}

				lastView, err := qb.LastView(ctx, tt.audioID)
				if err != nil {
					t.Errorf("AudioStore.LastView() error = %v", err)
				}

				assert.Equal(tt.expectedCount, count)
				assert.True(lastView.After(time.Now().Add(-1 * time.Minute)))

				return nil
			})
		})
	}
}

func TestAudioStore_DecrementWatchCount(t *testing.T) {
	return
}

func TestAudioStore_SaveActivity(t *testing.T) {
	var (
		resumeTime   = 111.2
		playDuration = 98.7
	)

	tests := []struct {
		name         string
		audioIdx     int
		resumeTime   *float64
		playDuration *float64
		wantErr      bool
	}{
		{
			"both",
			audioIdx1WithPerformer,
			&resumeTime,
			&playDuration,
			false,
		},
		{
			"resumeTime only",
			audioIdx1WithPerformer,
			&resumeTime,
			nil,
			false,
		},
		{
			"playDuration only",
			audioIdx1WithPerformer,
			nil,
			&playDuration,
			false,
		},
		{
			"none",
			audioIdx1WithPerformer,
			nil,
			nil,
			false,
		},
		{
			"invalid audio id",
			-1,
			&resumeTime,
			&playDuration,
			true,
		},
	}

	qb := db.Audio

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withRollbackTxn(func(ctx context.Context) error {
				id := -1
				if tt.audioIdx != -1 {
					id = audioIDs[tt.audioIdx]
				}

				_, err := qb.SaveActivity(ctx, id, tt.resumeTime, tt.playDuration)
				if (err != nil) != tt.wantErr {
					t.Errorf("AudioStore.SaveActivity() error = %v, wantErr %v", err, tt.wantErr)
				}

				if err != nil {
					return nil
				}

				assert := assert.New(t)

				// find the audio and check the values
				audio, err := qb.Find(ctx, id)
				if err != nil {
					t.Errorf("AudioStore.Find() error = %v", err)
				}

				expectedResumeTime := getResumeTime(tt.audioIdx)
				expectedPlayDuration := getPlayDuration(tt.audioIdx)

				if tt.resumeTime != nil {
					expectedResumeTime = *tt.resumeTime
				}
				if tt.playDuration != nil {
					expectedPlayDuration += *tt.playDuration
				}

				assert.Equal(expectedResumeTime, audio.ResumeTime)
				assert.Equal(expectedPlayDuration, audio.PlayDuration)

				return nil
			})
		})
	}
}

func TestAudioQueryCustomFields(t *testing.T) {
	tests := []struct {
		name        string
		filter      *models.AudioFilterType
		includeIdxs []int
		excludeIdxs []int
		wantErr     bool
	}{
		{
			"equals",
			&models.AudioFilterType{
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierEquals,
						Value:    []any{getAudioStringValue(audioIdxWithGallery, "custom")},
					},
				},
			},
			[]int{audioIdxWithGallery},
			nil,
			false,
		},
		{
			"not equals",
			&models.AudioFilterType{
				Title: &models.StringCriterionInput{
					Value:    getAudioTitle(audioIdxWithGallery),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierNotEquals,
						Value:    []any{getAudioStringValue(audioIdxWithGallery, "custom")},
					},
				},
			},
			nil,
			[]int{audioIdxWithGallery},
			false,
		},
		{
			"includes",
			&models.AudioFilterType{
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierIncludes,
						Value:    []any{getAudioStringValue(audioIdxWithGallery, "custom")[9:]},
					},
				},
			},
			[]int{audioIdxWithGallery},
			nil,
			false,
		},
		{
			"excludes",
			&models.AudioFilterType{
				Title: &models.StringCriterionInput{
					Value:    getAudioTitle(audioIdxWithGallery),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierExcludes,
						Value:    []any{getAudioStringValue(audioIdxWithGallery, "custom")[9:]},
					},
				},
			},
			nil,
			[]int{audioIdxWithGallery},
			false,
		},
		{
			"regex",
			&models.AudioFilterType{
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierMatchesRegex,
						Value:    []any{".*" + getAudioStringValue(audioIdxWithPerformerTag, "custom")[6:]},
					},
				},
			},
			[]int{audioIdxWithPerformerTag},
			nil,
			false,
		},
		{
			"invalid regex",
			&models.AudioFilterType{
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
			&models.AudioFilterType{
				Title: &models.StringCriterionInput{
					Value:    getAudioTitle(audioIdxWithPerformerTag),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierNotMatchesRegex,
						Value:    []any{".*" + getAudioStringValue(audioIdxWithPerformerTag, "custom")[6:]},
					},
				},
			},
			nil,
			[]int{audioIdxWithPerformerTag},
			false,
		},
		{
			"invalid not matches regex",
			&models.AudioFilterType{
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
			&models.AudioFilterType{
				Title: &models.StringCriterionInput{
					Value:    getAudioTitle(audioIdxWithGallery),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "not existing",
						Modifier: models.CriterionModifierIsNull,
					},
				},
			},
			[]int{audioIdxWithGallery},
			nil,
			false,
		},
		{
			"not null",
			&models.AudioFilterType{
				Title: &models.StringCriterionInput{
					Value:    getAudioTitle(audioIdxWithGallery),
					Modifier: models.CriterionModifierEquals,
				},
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "string",
						Modifier: models.CriterionModifierNotNull,
					},
				},
			},
			[]int{audioIdxWithGallery},
			nil,
			false,
		},
		{
			"between",
			&models.AudioFilterType{
				CustomFields: []models.CustomFieldCriterionInput{
					{
						Field:    "real",
						Modifier: models.CriterionModifierBetween,
						Value:    []any{0.15, 0.25},
					},
				},
			},
			[]int{audioIdxWithPerformer},
			nil,
			false,
		},
		{
			"not between",
			&models.AudioFilterType{
				Title: &models.StringCriterionInput{
					Value:    getAudioTitle(audioIdxWithPerformer),
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
			[]int{audioIdxWithPerformer},
			false,
		},
	}

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			result, err := db.Audio.Query(ctx, models.AudioQueryOptions{
				AudioFilter: tt.filter,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("AudioStore.Query() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				return
			}

			audios, err := result.Resolve(ctx)
			if err != nil {
				t.Errorf("AudioStore.Query().Resolve() error = %v", err)
				return
			}

			ids := audiosToIDs(audios)
			include := indexesToIDs(audioIDs, tt.includeIdxs)
			exclude := indexesToIDs(audioIDs, tt.excludeIdxs)

			for _, i := range include {
				assert.Contains(ids, i)
			}
			for _, e := range exclude {
				assert.NotContains(ids, e)
			}
		})
	}
}

// TODO Count
// TODO SizeCount

// TODO - this should be in history_test and generalised
func TestAudioStore_CountAllViews(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		qb := db.Audio

		audioID := audioIDs[audioIdx1WithPerformer]

		// get the current play count
		currentCount, err := qb.CountAllViews(ctx)
		if err != nil {
			t.Errorf("AudioStore.CountAllViews() error = %v", err)
			return nil
		}

		// add a view
		_, err = qb.AddViews(ctx, audioID, nil)
		if err != nil {
			t.Errorf("AudioStore.AddViews() error = %v", err)
			return nil
		}

		// get the new play count
		newCount, err := qb.CountAllViews(ctx)
		if err != nil {
			t.Errorf("AudioStore.CountAllViews() error = %v", err)
			return nil
		}

		assert.Equal(t, currentCount+1, newCount)

		return nil
	})
}

func TestAudioStore_CountUniqueViews(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		qb := db.Audio

		audioID := audioIDs[audioIdx1WithPerformer]

		// get the current play count
		currentCount, err := qb.CountUniqueViews(ctx)
		if err != nil {
			t.Errorf("AudioStore.CountUniqueViews() error = %v", err)
			return nil
		}

		// add a view
		_, err = qb.AddViews(ctx, audioID, nil)
		if err != nil {
			t.Errorf("AudioStore.AddViews() error = %v", err)
			return nil
		}

		// add a second view
		_, err = qb.AddViews(ctx, audioID, nil)
		if err != nil {
			t.Errorf("AudioStore.AddViews() error = %v", err)
			return nil
		}

		// get the new play count
		newCount, err := qb.CountUniqueViews(ctx)
		if err != nil {
			t.Errorf("AudioStore.CountUniqueViews() error = %v", err)
			return nil
		}

		assert.Equal(t, currentCount+1, newCount)

		return nil
	})
}
