package scene

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
			"id only",
			&UpdateSet{
				Partial: models.ScenePartial{
					ID: 1,
				},
			},
			true,
		},
		{
			"partial set",
			&UpdateSet{
				Partial: models.ScenePartial{
					Organized: &organized,
				},
			},
			false,
		},
		{
			"performer set",
			&UpdateSet{
				PerformerIDs: ids,
			},
			false,
		},
		{
			"tags set",
			&UpdateSet{
				TagIDs: ids,
			},
			false,
		},
		{
			"performer set",
			&UpdateSet{
				StashIDs: stashIDs,
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

type mockCoverSetter struct {
	mock.Mock
}

func (_m *mockCoverSetter) SetCover(sceneID int, imageData []byte) error {
	ret := _m.Called(sceneID, imageData)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, []byte) error); ok {
		r0 = rf(sceneID, imageData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
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
	stashIDs := []models.StashID{
		{
			StashID:  stashID,
			Endpoint: endpoint,
		},
	}

	title := "title"
	cover := []byte("cover")

	validScene := &models.Scene{}

	updateErr := errors.New("error updating")

	qb := mocks.SceneReaderWriter{}
	qb.On("Update", mock.MatchedBy(func(s models.ScenePartial) bool {
		return s.ID != badUpdateID
	})).Return(validScene, nil)
	qb.On("Update", mock.MatchedBy(func(s models.ScenePartial) bool {
		return s.ID == badUpdateID
	})).Return(nil, updateErr)

	qb.On("UpdatePerformers", sceneID, performerIDs).Return(nil).Once()
	qb.On("UpdateTags", sceneID, tagIDs).Return(nil).Once()
	qb.On("UpdateStashIDs", sceneID, stashIDs).Return(nil).Once()

	qb.On("UpdatePerformers", badPerformersID, performerIDs).Return(updateErr).Once()
	qb.On("UpdateTags", badTagsID, tagIDs).Return(updateErr).Once()
	qb.On("UpdateStashIDs", badStashIDsID, stashIDs).Return(updateErr).Once()

	coverSetter := &mockCoverSetter{}
	coverSetter.On("SetCover", sceneID, cover).Return(nil).Once()
	coverSetter.On("SetCover", badCoverID, cover).Return(updateErr).Once()

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
				ID:           sceneID,
				PerformerIDs: performerIDs,
				TagIDs:       tagIDs,
				StashIDs: []models.StashID{
					{
						StashID:  stashID,
						Endpoint: endpoint,
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
					Title: models.NullStringPtr(title),
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
					Title: models.NullStringPtr(title),
				},
			},
			true,
			true,
		},
		{
			"error updating performers",
			&UpdateSet{
				ID:           badPerformersID,
				PerformerIDs: performerIDs,
			},
			true,
			true,
		},
		{
			"error updating tags",
			&UpdateSet{
				ID:     badTagsID,
				TagIDs: tagIDs,
			},
			true,
			true,
		},
		{
			"error updating stash IDs",
			&UpdateSet{
				ID:       badStashIDsID,
				StashIDs: stashIDs,
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
			got, err := tt.u.Update(&qb, coverSetter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Updater.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("Updater.Update() = %v, want %v", got, tt.wantNil)
			}
		})
	}

	qb.AssertExpectations(t)
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
	stashIDs := []models.StashID{
		{
			StashID:  stashID,
			Endpoint: endpoint,
		},
	}
	stashIDInputs := []*models.StashIDInput{
		{
			StashID:  stashID,
			Endpoint: endpoint,
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
				ID:           sceneID,
				PerformerIDs: performerIDs,
				TagIDs:       tagIDs,
				StashIDs:     stashIDs,
				CoverImage:   cover,
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
					Title: models.NullStringPtr(title),
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
