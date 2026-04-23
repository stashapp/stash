// TODO(audio): update this file

package audio

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdater_IsEmpty(t *testing.T) {
	organized := true
	ids := []int{1}
	cover := []byte{1}

	tests := []struct {
		name string
		u    *UpdateSet
		want bool
	}{
		{
			"empty",
			&UpdateSet{},
			true,
		},
		{
			"partial set",
			&UpdateSet{
				Partial: models.AudioPartial{
					Organized: models.NewOptionalBool(organized),
				},
			},
			false,
		},
		{
			"performer set",
			&UpdateSet{
				Partial: models.AudioPartial{
					PerformerIDs: &models.UpdateIDs{
						IDs:  ids,
						Mode: models.RelationshipUpdateModeSet,
					},
				},
			},
			false,
		},
		{
			"tags set",
			&UpdateSet{
				Partial: models.AudioPartial{
					TagIDs: &models.UpdateIDs{
						IDs:  ids,
						Mode: models.RelationshipUpdateModeSet,
					},
				},
			},
			false,
		},
		{
			"cover set",
			&UpdateSet{
				CoverImage: cover,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.IsEmpty(); got != tt.want {
				t.Errorf("Updater.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdater_Update(t *testing.T) {
	const (
		audioID = iota + 1
		badUpdateID
		badPerformersID
		badTagsID
		badCoverID
		performerID
		tagID
	)

	performerIDs := []int{performerID}
	tagIDs := []int{tagID}

	title := "title"
	cover := []byte("cover")

	validAudio := &models.Audio{}

	updateErr := errors.New("error updating")

	db := mocks.NewDatabase()

	db.Audio.On("UpdatePartial", testCtx, mock.MatchedBy(func(id int) bool {
		return id != badUpdateID
	}), mock.Anything).Return(validAudio, nil)
	db.Audio.On("UpdatePartial", testCtx, badUpdateID, mock.Anything).Return(nil, updateErr)

	db.Audio.On("UpdateCover", testCtx, audioID, cover).Return(nil).Once()
	db.Audio.On("UpdateCover", testCtx, badCoverID, cover).Return(updateErr).Once()

	tests := []struct {
		name    string
		u       *UpdateSet
		wantNil bool
		wantErr bool
	}{
		{
			"empty",
			&UpdateSet{
				ID: audioID,
			},
			true,
			true,
		},
		{
			"update all",
			&UpdateSet{
				ID: audioID,
				Partial: models.AudioPartial{
					PerformerIDs: &models.UpdateIDs{
						IDs:  performerIDs,
						Mode: models.RelationshipUpdateModeSet,
					},
					TagIDs: &models.UpdateIDs{
						IDs:  tagIDs,
						Mode: models.RelationshipUpdateModeSet,
					},
				},
				CoverImage: cover,
			},
			false,
			false,
		},
		{
			"update fields only",
			&UpdateSet{
				ID: audioID,
				Partial: models.AudioPartial{
					Title: models.NewOptionalString(title),
				},
			},
			false,
			false,
		},
		{
			"error updating audio",
			&UpdateSet{
				ID: badUpdateID,
				Partial: models.AudioPartial{
					Title: models.NewOptionalString(title),
				},
			},
			true,
			true,
		},
		{
			"error updating cover",
			&UpdateSet{
				ID:         badCoverID,
				CoverImage: cover,
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.Update(testCtx, db.Audio)
			if (err != nil) != tt.wantErr {
				t.Errorf("Updater.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("Updater.Update() = %v, want %v", got, tt.wantNil)
			}
		})
	}

	db.AssertExpectations(t)
}

func TestUpdateSet_UpdateInput(t *testing.T) {
	const (
		audioID = iota + 1
		badUpdateID
		badPerformersID
		badTagsID
		badCoverID
		performerID
		tagID
	)

	audioIDStr := strconv.Itoa(audioID)

	performerIDs := []int{performerID}
	performerIDStrs := intslice.IntSliceToStringSlice(performerIDs)
	tagIDs := []int{tagID}
	tagIDStrs := intslice.IntSliceToStringSlice(tagIDs)

	title := "title"
	cover := []byte("cover")
	coverB64 := "Y292ZXI="

	tests := []struct {
		name string
		u    UpdateSet
		want models.AudioUpdateInput
	}{
		{
			"empty",
			UpdateSet{
				ID: audioID,
			},
			models.AudioUpdateInput{
				ID: audioIDStr,
			},
		},
		{
			"update all",
			UpdateSet{
				ID: audioID,
				Partial: models.AudioPartial{
					PerformerIDs: &models.UpdateIDs{
						IDs:  performerIDs,
						Mode: models.RelationshipUpdateModeSet,
					},
					TagIDs: &models.UpdateIDs{
						IDs:  tagIDs,
						Mode: models.RelationshipUpdateModeSet,
					},
				},
				CoverImage: cover,
			},
			models.AudioUpdateInput{
				ID:           audioIDStr,
				PerformerIds: performerIDStrs,
				TagIds:       tagIDStrs,
				CoverImage:   &coverB64,
			},
		},
		{
			"update fields only",
			UpdateSet{
				ID: audioID,
				Partial: models.AudioPartial{
					Title: models.NewOptionalString(title),
				},
			},
			models.AudioUpdateInput{
				ID:    audioIDStr,
				Title: &title,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.u.UpdateInput()
			assert.Equal(t, tt.want, got)
		})
	}
}
