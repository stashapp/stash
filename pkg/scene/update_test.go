package scene

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdater_IsEmpty(t *testing.T) {
	organized := true
	ids := []int{1}
	stashIDs := []models.StashID{
		{},
	}
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
				Partial: models.ScenePartial{
					Organized: models.NewOptionalBool(organized),
				},
			},
			false,
		},
		{
			"performer set",
			&UpdateSet{
				Partial: models.ScenePartial{
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
				Partial: models.ScenePartial{
					TagIDs: &models.UpdateIDs{
						IDs:  ids,
						Mode: models.RelationshipUpdateModeSet,
					},
				},
			},
			false,
		},
		{
			"performer set",
			&UpdateSet{
				Partial: models.ScenePartial{
					StashIDs: &models.UpdateStashIDs{
						StashIDs: stashIDs,
						Mode:     models.RelationshipUpdateModeSet,
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
		sceneID = iota + 1
		badUpdateID
		badPerformersID
		badTagsID
		badStashIDsID
		badCoverID
		performerID
		tagID
	)

	performerIDs := []int{performerID}
	tagIDs := []int{tagID}
	stashID := "stashID"
	endpoint := "endpoint"

	title := "title"
	cover := []byte("cover")

	validScene := &models.Scene{}

	updateErr := errors.New("error updating")

	db := mocks.NewDatabase()

	db.Scene.On("UpdatePartial", testCtx, mock.MatchedBy(func(id int) bool {
		return id != badUpdateID
	}), mock.Anything).Return(validScene, nil)
	db.Scene.On("UpdatePartial", testCtx, badUpdateID, mock.Anything).Return(nil, updateErr)

	db.Scene.On("UpdateCover", testCtx, sceneID, cover).Return(nil).Once()
	db.Scene.On("UpdateCover", testCtx, badCoverID, cover).Return(updateErr).Once()

	tests := []struct {
		name    string
		u       *UpdateSet
		wantNil bool
		wantErr bool
	}{
		{
			"empty",
			&UpdateSet{
				ID: sceneID,
			},
			true,
			true,
		},
		{
			"update all",
			&UpdateSet{
				ID: sceneID,
				Partial: models.ScenePartial{
					PerformerIDs: &models.UpdateIDs{
						IDs:  performerIDs,
						Mode: models.RelationshipUpdateModeSet,
					},
					TagIDs: &models.UpdateIDs{
						IDs:  tagIDs,
						Mode: models.RelationshipUpdateModeSet,
					},
					StashIDs: &models.UpdateStashIDs{
						StashIDs: []models.StashID{
							{
								StashID:  stashID,
								Endpoint: endpoint,
							},
						},
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
				ID: sceneID,
				Partial: models.ScenePartial{
					Title: models.NewOptionalString(title),
				},
			},
			false,
			false,
		},
		{
			"error updating scene",
			&UpdateSet{
				ID: badUpdateID,
				Partial: models.ScenePartial{
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
			got, err := tt.u.Update(testCtx, db.Scene)
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
		sceneID = iota + 1
		badUpdateID
		badPerformersID
		badTagsID
		badStashIDsID
		badCoverID
		performerID
		tagID
	)

	sceneIDStr := strconv.Itoa(sceneID)

	performerIDs := []int{performerID}
	performerIDStrs := intslice.IntSliceToStringSlice(performerIDs)
	tagIDs := []int{tagID}
	tagIDStrs := intslice.IntSliceToStringSlice(tagIDs)
	stashID := "stashID"
	endpoint := "endpoint"
	updatedAt := time.Now()
	stashIDs := []models.StashID{
		{
			StashID:   stashID,
			Endpoint:  endpoint,
			UpdatedAt: updatedAt,
		},
	}
	stashIDInputs := []models.StashIDInput{
		{
			StashID:   stashID,
			Endpoint:  endpoint,
			UpdatedAt: &updatedAt,
		},
	}

	title := "title"
	cover := []byte("cover")
	coverB64 := "Y292ZXI="

	tests := []struct {
		name string
		u    UpdateSet
		want models.SceneUpdateInput
	}{
		{
			"empty",
			UpdateSet{
				ID: sceneID,
			},
			models.SceneUpdateInput{
				ID: sceneIDStr,
			},
		},
		{
			"update all",
			UpdateSet{
				ID: sceneID,
				Partial: models.ScenePartial{
					PerformerIDs: &models.UpdateIDs{
						IDs:  performerIDs,
						Mode: models.RelationshipUpdateModeSet,
					},
					TagIDs: &models.UpdateIDs{
						IDs:  tagIDs,
						Mode: models.RelationshipUpdateModeSet,
					},
					StashIDs: &models.UpdateStashIDs{
						StashIDs: stashIDs,
						Mode:     models.RelationshipUpdateModeSet,
					},
				},
				CoverImage: cover,
			},
			models.SceneUpdateInput{
				ID:           sceneIDStr,
				PerformerIds: performerIDStrs,
				TagIds:       tagIDStrs,
				StashIds:     stashIDInputs,
				CoverImage:   &coverB64,
			},
		},
		{
			"update fields only",
			UpdateSet{
				ID: sceneID,
				Partial: models.ScenePartial{
					Title: models.NewOptionalString(title),
				},
			},
			models.SceneUpdateInput{
				ID:    sceneIDStr,
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
